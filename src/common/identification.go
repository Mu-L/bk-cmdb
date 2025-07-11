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

package common

import (
	"configcenter/src/common/types"
	"configcenter/src/common/version"
)

var identification = "unknown"
var server *types.ServerInfo

// SetIdentification TODO
func SetIdentification(id string) {
	if identification == "unknown" {
		version.ServiceName = id
		identification = id
	}
}

// GetIdentification TODO
func GetIdentification() string {
	return identification
}

// SetServerInfo Information about the current process in service governance
func SetServerInfo(srvInfo *types.ServerInfo) {
	server = srvInfo
}

// GetServerInfo Information about the current process in service governance
func GetServerInfo() *types.ServerInfo {
	return server
}
