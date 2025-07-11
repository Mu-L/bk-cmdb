/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package instances

import (
	"sort"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/common/valid"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/mongodb/instancemapping"
)

func (m *instanceManager) batchSave(kit *rest.Kit, objID string, params []mapstr.MapStr) ([]uint64, error) {
	instTableName := common.GetInstTableName(objID, kit.SupplierAccount)
	instIDFieldName := common.GetInstIDField(objID)
	ids, err := getSequences(kit, instTableName, len(params))
	if err != nil {
		return nil, err
	}
	mappings := make([]mapstr.MapStr, 0)
	ts := time.Now()

	for idx := range params {
		if objID == common.BKInnerObjIDHost {
			params[idx], err = metadata.ConvertHostSpecialStringToArray(params[idx])
			if err != nil {
				blog.Errorf("convert host special string to array failed, err: %v, rid: %s", err, kit.Rid)
				return nil, err
			}
		}

		// build new object instance data.
		if !valid.IsInnerObject(objID) {
			params[idx][common.BKObjIDField] = objID
		}
		params[idx].Set(instIDFieldName, ids[idx])
		params[idx].Set(common.BKOwnerIDField, kit.SupplierAccount)
		params[idx].Set(common.CreateTimeField, ts)
		params[idx].Set(common.LastTimeField, ts)

		params[idx].Set(common.BKCreatedBy, kit.User)
		params[idx].Set(common.BKCreatedAt, ts)
		params[idx].Set(common.BKUpdatedAt, ts)

		if !metadata.IsCommon(objID) {
			continue
		}
		// build new object mapping data for inner object instance.
		mapping := make(mapstr.MapStr, 0)
		mapping[instIDFieldName] = ids[idx]
		mapping[common.BKObjIDField] = objID
		mapping[common.BkSupplierAccount] = kit.SupplierAccount

		mappings = append(mappings, mapping)
	}

	if len(mappings) != 0 {
		// save new object mappings data for inner object instance.
		if err := instancemapping.Create(kit.Ctx, mappings); err != nil {
			return nil, err
		}
	}

	// save object instances.
	err = mongodb.Client().Table(instTableName).Insert(kit.Ctx, params)
	if err != nil {
		blog.Errorf("save instances failed, rid: %s, err: %v, objID: %s, instances: %v", kit.Rid, err, objID, params)
		if mongodb.Client().IsDuplicatedError(err) {
			return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err))
		}
		return nil, err
	}

	return ids, nil
}

func (m *instanceManager) save(kit *rest.Kit, objID string, inputParam mapstr.MapStr) (uint64, error) {
	if objID == common.BKInnerObjIDHost {
		var err error
		inputParam, err = metadata.ConvertHostSpecialStringToArray(inputParam)
		if err != nil {
			return 0, err
		}
	}

	instTableName := common.GetInstTableName(objID, kit.SupplierAccount)
	ids, err := getSequences(kit, instTableName, 1)
	if err != nil {
		return 0, err
	}

	// build new object instance data.
	instIDFieldName := common.GetInstIDField(objID)
	inputParam[instIDFieldName] = ids[0]
	if !valid.IsInnerObject(objID) {
		inputParam[common.BKObjIDField] = objID
	}
	ts := time.Now()
	inputParam.Set(common.BKOwnerIDField, kit.SupplierAccount)
	inputParam.Set(common.CreateTimeField, ts)
	inputParam.Set(common.LastTimeField, ts)

	inputParam.Set(common.BKCreatedBy, kit.User)
	inputParam.Set(common.BKCreatedAt, ts)
	inputParam.Set(common.BKUpdatedAt, ts)

	// build and save new object mapping data for inner object instance.
	if metadata.IsCommon(objID) {
		mapping := make(mapstr.MapStr, 0)
		mapping[instIDFieldName] = ids[0]
		mapping[common.BKObjIDField] = objID
		mapping[common.BkSupplierAccount] = kit.SupplierAccount

		// save instance object type mapping.
		if err := instancemapping.Create(kit.Ctx, mapping); err != nil {
			return 0, err
		}
	}

	// save object instance.
	err = mongodb.Client().Table(instTableName).Insert(kit.Ctx, inputParam)
	if err != nil {
		blog.ErrorJSON("save instance error. err: %s, objID: %s, instance: %s, rid: %s",
			err.Error(), objID, inputParam, kit.Rid)
		if mongodb.Client().IsDuplicatedError(err) {
			return ids[0], kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err))
		}
		return 0, err
	}

	return ids[0], nil
}

func getSequences(kit *rest.Kit, table string, count int) ([]uint64, error) {
	if count <= 0 {
		return nil, kit.CCError.CCError(common.CCErrCommHTTPInputInvalid)
	}

	ids, err := mongodb.Client().NextSequences(kit.Ctx, table, count)
	if err != nil {
		return nil, err
	}

	if table != common.BKTableNameBasePlat {
		return ids, nil
	}

	sort.Sort(util.Uint64Slice(ids))

	if ids[0] > common.ReservedCloudAreaEndID {
		return ids, nil
	}

	if ids[len(ids)-1] < common.ReservedCloudAreaStartID {
		return ids, nil
	}

	// 此处不要求那么准确，直接跳过保留的管控区域长度，然后再获取id即可
	if _, err = mongodb.Client().NextSequences(kit.Ctx, table,
		common.ReservedCloudAreaEndID-common.ReservedCloudAreaStartID); err != nil {
		return nil, err
	}

	ids, err = mongodb.Client().NextSequences(kit.Ctx, table, count)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (m *instanceManager) update(kit *rest.Kit, objID string, data mapstr.MapStr, cond mapstr.MapStr) errors.CCError {
	if objID == common.BKInnerObjIDHost {
		var err error
		data, err = metadata.ConvertHostSpecialStringToArray(data)
		if err != nil {
			return err
		}
	}
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)
	if !valid.IsInnerObject(objID) {
		cond.Set(common.BKObjIDField, objID)
	}
	ts := time.Now()
	data.Set(common.LastTimeField, ts)
	data.Set(common.BKUpdatedBy, kit.User)
	data.Set(common.BKUpdatedAt, ts)

	data.Remove(common.BKObjIDField)
	err := mongodb.Client().Table(tableName).Update(kit.Ctx, cond, data)
	if err != nil {
		blog.ErrorJSON("update instance error. err: %s, objID: %s, instance: %s, cond: %s, rid: %s",
			err.Error(), objID, data, cond, kit.Rid)
		if mongodb.Client().IsDuplicatedError(err) {
			return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err))
		}
		return kit.CCError.Error(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

func (m *instanceManager) getInsts(kit *rest.Kit, objID string, cond mapstr.MapStr) (origins []mapstr.MapStr,
	exists bool, err error) {
	origins = make([]mapstr.MapStr, 0)
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)
	if !valid.IsInnerObject(objID) {
		cond.Set(common.BKObjIDField, objID)
	}
	if objID == common.BKInnerObjIDHost {
		hosts := make([]metadata.HostMapStr, 0)
		err = mongodb.Client().Table(tableName).Find(cond).All(kit.Ctx, &hosts)
		for _, host := range hosts {
			origins = append(origins, mapstr.MapStr(host))
		}
	} else {
		err = mongodb.Client().Table(tableName).Find(cond).All(kit.Ctx, &origins)
	}
	return origins, !mongodb.Client().IsNotFoundError(err), err
}

func (m *instanceManager) getInstDataByID(kit *rest.Kit, objID string, instID int64) (origin mapstr.MapStr, err error) {
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.GetInstIDField(objID), Val: instID})

	if common.IsObjectInstShardingTable(common.GetInstTableName(objID, kit.SupplierAccount)) {
		cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID})
	}

	if objID == common.BKInnerObjIDHost {
		host := make(metadata.HostMapStr)
		err = mongodb.Client().Table(tableName).Find(cond.ToMapStr()).One(kit.Ctx, &host)
		origin = mapstr.MapStr(host)
	} else {
		err = mongodb.Client().Table(tableName).Find(cond.ToMapStr()).One(kit.Ctx, &origin)
	}
	if nil != err {
		return nil, err
	}
	return origin, nil
}

func (m *instanceManager) countInstance(kit *rest.Kit, objID string, cond mapstr.MapStr) (count uint64, err error) {
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)

	if cond == nil {
		cond = make(map[string]interface{})
	}

	if common.IsObjectInstShardingTable(tableName) {
		objIDCond, ok := cond[common.BKObjIDField]
		if ok && objIDCond != objID {
			blog.V(9).Infof("countInstance condition's bk_obj_id: %s not match objID: %s, rid: %s", objIDCond, objID,
				kit.Rid)
			return 0, nil
		}
		cond[common.BKObjIDField] = objID
	}

	cond = util.SetQueryOwner(cond, kit.SupplierAccount)
	count, err = mongodb.Client().Table(tableName).Find(cond).Count(kit.Ctx)

	return count, err
}
