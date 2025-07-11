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

// DeleteOption common delete condition options
type DeleteOption struct {
	Condition mapstr.MapStr `json:"condition"`
}

// DeleteTableOption delete table field option.
type DeleteTableOption struct {
	ObjID      string `json:"bk_obj_id"`
	ID         int64  `json:"id"`
	PropertyID string `json:"bk_property_id"`
}

// DeletedOptionResult delete  api http response return result struct
type DeletedOptionResult struct {
	BaseResp `json:",inline"`
	Data     DeletedCount `json:"data"`
}
