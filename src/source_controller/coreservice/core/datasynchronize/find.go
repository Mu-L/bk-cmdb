// Package datasynchronize TODO
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
package datasynchronize

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

type associationFindDataInterface interface {
	Find(kit *rest.Kit) ([]mapstr.MapStr, uint64, errors.CCError)
}

type associationFindData struct {
	dataClassify string
	dataType     metadata.SynchronizeOperateDataType
	start        uint64
	limit        uint64
	condition    mapstr.MapStr
}

// NewSynchronizeFindAdapter TODO
func NewSynchronizeFindAdapter(input *metadata.SynchronizeFindInfoParameter) associationFindDataInterface {

	return &associationFindData{
		dataClassify: input.DataClassify,
		dataType:     input.DataType,
		start:        input.Start,
		limit:        input.Limit,
		condition:    input.Condition,
	}
}

// Find TODO
func (a *associationFindData) Find(kit *rest.Kit) ([]mapstr.MapStr, uint64, errors.CCError) {
	switch a.dataType {
	case metadata.SynchronizeOperateDataTypeAssociation:
		return a.findAssociation(kit)
	case metadata.SynchronizeOperateDataTypeModel:
		return a.findModel(kit)
	}
	return nil, 0, nil
}

func (a *associationFindData) findModel(kit *rest.Kit) ([]mapstr.MapStr, uint64, errors.CCError) {

	switch a.dataClassify {
	case common.SynchronizeModelTypeBase:
		return a.dbQueryModel(kit, common.BKTableNameObjDes)
	case common.SynchronizeModelTypeAttribute:
		return a.dbQueryModel(kit, common.BKTableNameObjAttDes)
	case common.SynchronizeModelTypeAttributeGroup:
		return a.dbQueryModel(kit, common.BKTableNamePropertyGroup)
	case common.SynchronizeModelTypeClassification:
		return a.dbQueryModel(kit, common.BKTableNameObjClassification)
	}
	return nil, 0, nil
}

func (a *associationFindData) dbQueryModel(kit *rest.Kit, tableName string) ([]mapstr.MapStr, uint64, errors.CCError) {
	info := make([]mapstr.MapStr, 0)
	err := mongodb.Client().Table(tableName).Find(a.condition).Start(a.start).Limit(a.limit).All(kit.Ctx, &info)
	if err != nil {
		blog.Errorf("dbQueryModel info error. error:%s,rid:%s", err.Error(), kit.Rid)
		return nil, 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}
	cnt, err := mongodb.Client().Table(tableName).Find(nil).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("dbQueryModel count error. error:%s,rid:%s", err.Error(), kit.Rid)
		return nil, 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}
	return info, cnt, nil
}

func (a *associationFindData) findAssociation(kit *rest.Kit) ([]mapstr.MapStr, uint64, errors.CCError) {
	switch a.dataClassify {
	case common.SynchronizeAssociationTypeModelHost:
		return a.dbQueryAssociation(kit)
	}
	return nil, 0, nil
}

func (a *associationFindData) dbQueryAssociation(kit *rest.Kit) ([]mapstr.MapStr, uint64, errors.CCError) {
	info := make([]mapstr.MapStr, 0)
	err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(a.condition).Start(a.start).Limit(a.limit).All(kit.Ctx, &info)
	if err != nil {
		blog.Errorf("dbQueryAssociation info error. error:%s,rid:%s", err.Error(), kit.Rid)
		return nil, 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}
	cnt, err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(nil).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("dbQueryAssociation count error. error:%s,rid:%s", err.Error(), kit.Rid)
		return nil, 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	return info, cnt, nil

}
