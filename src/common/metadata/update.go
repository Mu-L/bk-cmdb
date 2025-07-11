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

package metadata

import "configcenter/src/common/mapstr"

// UpdateOption common update options
type UpdateOption struct {
	Data      mapstr.MapStr `json:"data" mapstructure:"data"`
	Condition mapstr.MapStr `json:"condition" mapstructure:"condition"`
	// can edit all fields, including not editable properties, used by collectors
	CanEditAll bool `json:"can_edit_all" mapstructure:"can_edit_all"`
	// IsSync in the update scene, it is distinguished whether it is
	// a synchronization scene of the field template
	IsSync bool `json:"is_sync" mapstructure:"is_sync"`
}

// CreatePartDataOption newly added headers and default values the user update scenario
type CreatePartDataOption struct {
	ObjID string      `json:"bk_obj_id" mapstructure:"bk_obj_id"`
	Data  []Attribute `json:"data" mapstructure:"data"`
}

// UpdateTableOption common update options
type UpdateTableOption struct {
	// CreateData is the table attribute header data to be created, used only for validation
	CreateData CreatePartDataOption `json:"create_data" mapstructure:"create_data"`
	// UpdateData is the complete table attribute data to be updated
	UpdateData mapstr.MapStr `json:"update_data" mapstructure:"update_data"`
	Condition  mapstr.MapStr `json:"condition" mapstructure:"condition"`
	// IsSync in the update scene, it is distinguished whether it is
	// a synchronization scene of the field template
	IsSync bool `json:"is_sync" mapstructure:"is_sync"`
}

// UpdatedOptionResult common update result
type UpdatedOptionResult struct {
	BaseResp `json:",inline"`
	Data     UpdatedCount `json:"data" mapstructure:"data"`
}

// UpdateAttrIndexInput update object attribute index input
type UpdateAttrIndexInput struct {
	BizID         int64  `json:"bk_biz_id"`
	PropertyGroup string `json:"bk_property_group"`
	PropertyIndex int64  `json:"bk_property_index"`
}
