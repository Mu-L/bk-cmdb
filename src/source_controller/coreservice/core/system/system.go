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

// Package system TODO
package system

import (
	"encoding/json"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
)

var _ core.SystemOperation = (*systemManager)(nil)

type systemManager struct {
}

// New create a new instance manager instance
func New() core.SystemOperation {
	return &systemManager{}
}

// GetSystemUserConfig TODO
func (sm *systemManager) GetSystemUserConfig(kit *rest.Kit) (map[string]interface{}, errors.CCErrorCoder) {
	cond := map[string]string{"type": metadata.CCSystemUserConfigSwitch}
	result := make(map[string]interface{}, 0)
	err := mongodb.Client().Table(common.BKTableNameSystem).Find(cond).One(kit.Ctx, &result)
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		blog.ErrorJSON("GetSystemUserConfig find error. cond:%s, err:%s, rid:%s", cond, err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return result, nil
}

// SearchConfigAdmin TODO
func (sm *systemManager) SearchConfigAdmin(kit *rest.Kit) (*metadata.ConfigAdmin, errors.CCErrorCoder) {
	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}

	ret := struct {
		Config string `json:"config"`
	}{}
	err := mongodb.Client().Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(kit.Ctx, &ret)
	if err != nil {
		blog.Errorf("SearchConfigAdmin failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	conf := new(metadata.ConfigAdmin)
	if err := json.Unmarshal([]byte(ret.Config), conf); err != nil {
		blog.Errorf("SearchConfigAdmin failed, Unmarshal err: %v, config:%+v,rid:%s", err, ret.Config, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	return conf, nil
}

// SearchPlatformSettingConfig search platform setting.
func (sm *systemManager) SearchPlatformSettingConfig(kit *rest.Kit) (*metadata.PlatformSettingConfig,
	errors.CCErrorCoder) {

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}

	ret := make(map[string]interface{})

	err := mongodb.Client().Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).
		One(kit.Ctx, &ret)
	if err != nil {
		blog.Errorf("search platform setting failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if ret[common.ConfigAdminValueField] == nil {
		blog.Errorf("search platform setting failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if _, ok := ret[common.ConfigAdminValueField].(string); !ok {
		blog.Errorf("search platform setting failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	conf := new(metadata.PlatformSettingConfig)
	if err := json.Unmarshal([]byte(ret[common.ConfigAdminValueField].(string)), conf); err != nil {
		blog.Errorf("platform setting unmarshal err: %v, config: %v,rid: %s", err,
			ret[common.ConfigAdminValueField].(string), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	return conf, nil
}

// UpdatePlatformSettingConfig update platform setting.
func (sm *systemManager) UpdatePlatformSettingConfig(kit *rest.Kit,
	input *metadata.PlatformSettingConfig) errors.CCErrorCoder {

	bytes, err := json.Marshal(input)
	if err != nil {
		blog.Errorf("update config admin failed, Marshal err: %v, input: %v, rid: %s", err, *input, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	data := map[string]interface{}{
		common.ConfigAdminValueField: string(bytes),
		common.LastTimeField:         time.Now(),
	}
	err = mongodb.Client().Table(common.BKTableNameSystem).Update(kit.Ctx, cond, data)
	if err != nil {
		blog.Errorf("update config admin failed, update err: %v, rid: %s", err, kit.Rid)

		return kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed, err)
	}

	return nil
}
