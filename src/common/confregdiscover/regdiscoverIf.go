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

package confregdiscover

// ConfRegDiscvIf define the configurate register and discover function interface
type ConfRegDiscvIf interface {
	// Ping to ping server
	Ping() error
	// Write the config data into configure register-discover service
	Write(key string, data []byte) error
	// Read the config data from configure register-discover service
	Read(key string) (string, error)
	// Discover the config change
	Discover(key string) (<-chan *DiscoverEvent, error)
}
