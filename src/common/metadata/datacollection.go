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

const (
	// HostSnapDataSourcesDelayQueue host snap data comes from the delay queue.
	HostSnapDataSourcesDelayQueue = "delay_queue"
	// HostSnapDataSourcesChannel the source of host snap data is channels such as redis or kafka.
	HostSnapDataSourcesChannel = "channel"
)

// AddDeviceResult TODO
type AddDeviceResult struct {
	DeviceID uint64 `json:"device_id"`
}

// BatchAddDevice TODO
type BatchAddDevice struct {
	Data []NetcollectDevice `json:"data"`
}

// BatchAddDeviceResult TODO
type BatchAddDeviceResult struct {
	Result   bool   `json:"result"`
	ErrMsg   string `json:"error_msg"`
	DeviceID uint64 `json:"device_id"`
}

// SearchNetDevice TODO
type SearchNetDevice struct {
	Count uint64             `json:"count"`
	Info  []NetcollectDevice `json:"info"`
}

// SearchNetDeviceResult TODO
type SearchNetDeviceResult struct {
	BaseResp `json:",inline"`
	Data     SearchNetDevice `json:"data"`
}

// NetCollSearchParams TODO
type NetCollSearchParams struct {
	Page      BasePage        `json:"page,omitempty"`
	Fields    []string        `json:"fields,omitempty"`
	Condition []ConditionItem `json:"condition,omitempty"`
}

// DeleteNetDeviceBatchOpt TODO
type DeleteNetDeviceBatchOpt struct {
	DeviceIDs []uint64 `json:"device_id"`
}

// AddNetPropertyResult TODO
type AddNetPropertyResult struct {
	NetcollectPropertyID uint64 `json:"netcollect_property_id"`
}

// BatchAddNetPropertyResult TODO
type BatchAddNetPropertyResult struct {
	Result               bool   `json:"result"`
	ErrMsg               string `json:"error_msg"`
	NetcollectPropertyID uint64 `json:"netcollect_property_id"`
}

// BatchAddNetProperty TODO
type BatchAddNetProperty struct {
	Data []NetcollectProperty `json:"data"`
}

// SearchNetProperty TODO
type SearchNetProperty struct {
	Count uint64               `json:"count"`
	Info  []NetcollectProperty `json:"info"`
}

// SearchNetPropertyResult TODO
type SearchNetPropertyResult struct {
	BaseResp `json:",inline"`
	Data     SearchNetProperty `json:"data"`
}

// DeleteNetPropertyBatchOpt TODO
type DeleteNetPropertyBatchOpt struct {
	NetcollectPropertyIDs []uint64 `json:"netcollect_property_id"`
}
