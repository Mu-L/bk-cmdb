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

import (
	"time"

	"configcenter/src/common"
)

// ChartConfig TODO
type ChartConfig struct {
	ConfigID   uint64 `json:"config_id" bson:"config_id"`
	ReportType string `json:"report_type" bson:"report_type"`
	Name       string `json:"name" bson:"name"`
	CreateTime Time   `json:"create_time" bson:"create_time"`
	OwnerID    string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjID      string `json:"bk_obj_id" bson:"bk_obj_id"`
	Width      string `json:"width" bson:"width"`
	ChartType  string `json:"chart_type" bson:"chart_type"`
	Field      string `json:"field" bson:"field"`
	XAxisCount int64  `json:"x_axis_count" bson:"x_axis_count"`
}

// ChartPosition TODO
type ChartPosition struct {
	BizID    int64        `json:"bk_biz_id" bson:"bk_biz_id"`
	Position PositionInfo `json:"position" bson:"position"`
	OwnerID  string       `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// PositionInfo TODO
type PositionInfo struct {
	Host []uint64 `json:"host" bson:"host"`
	Inst []uint64 `json:"inst" bson:"inst"`
}

// ModelInstChange TODO
type ModelInstChange map[string]*InstChangeCount

// InstChangeCount TODO
type InstChangeCount struct {
	Create int64 `json:"create" bson:"create"`
	Update int64 `json:"update" bson:"update"`
	Delete int64 `json:"delete" bson:"delete"`
}

// AggregateIntResponse TODO
type AggregateIntResponse struct {
	BaseResp `json:",inline"`
	Data     []IntIDCount `json:"data"`
}

// IntIDCount int类型字段做mongoDB聚合时使用
type IntIDCount struct {
	ID    int64 `json:"id" bson:"_id"`
	Count int64 `json:"count" bson:"count"`
}

// IntIDArrayCount int类型字段做mongoDB聚合，且结果为数组时使用
type IntIDArrayCount struct {
	ID    int64   `json:"id" bson:"_id"`
	Count []int64 `json:"count" bson:"count"`
}

// AggregateStringResponse TODO
type AggregateStringResponse struct {
	BaseResp `json:",inline"`
	Data     []StringIDCount `json:"data"`
}

// StringIDCount string类型字段做mongoDB聚合时使用
type StringIDCount struct {
	ID    string `json:"id" bson:"_id"`
	Count int64  `json:"count" bson:"count"`
}

// ObjectIDCount object count statistics information used for
// group aggregate operation.
type ObjectIDCount struct {
	// ObjID object id.
	ObjID string `bson:"_id" json:"bk_obj_id"`

	// Count targets count.
	Count int64 `bson:"count" json:"instance_count"`
}

// UpdateInstCount TODO
type UpdateInstCount struct {
	ID    UpdateID `json:"id" bson:"_id"`
	Count int64    `json:"count" bson:"count"`
}

// UpdateID TODO
type UpdateID struct {
	ObjID  string `json:"bk_obj_id" bson:"bk_obj_id"`
	InstID int64  `json:"bk_inst_id" bson:"bk_inst_id"`
}

// HostChangeChartData TODO
type HostChangeChartData struct {
	ReportType string          `json:"report_type" bson:"report_type"`
	Data       []StringIDCount `json:"data" bson:"data"`
	OwnerID    string          `json:"bk_supplier_account" bson:"bk_supplier_account"`
	CreateTime string          `json:"create_time" bson:"create_time"`
}

// ChartData TODO
type ChartData struct {
	ReportType string      `json:"report_type" bson:"report_type"`
	Data       interface{} `json:"data" data:"data"`
	OwnerID    string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime   time.Time   `json:"last_time" bson:"last_time"`
}

// ModelInstChartData TODO
type ModelInstChartData struct {
	ReportType string          `json:"report_type" bson:"report_type"`
	Data       []StringIDCount `json:"data" data:"data"`
	OwnerID    string          `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime   time.Time       `json:"last_time" bson:"last_time"`
}

// SearchChartResponse TODO
type SearchChartResponse struct {
	BaseResp `json:",inline"`
	Data     SearchChartConfig `json:"data"`
}

// SearchChartCommon TODO
type SearchChartCommon struct {
	BaseResp `json:",inline"`
	Data     CommonSearchChart `json:"data"`
}

// CommonSearchChart TODO
type CommonSearchChart struct {
	Count uint64      `json:"count"`
	Info  ChartConfig `json:"info"`
}

// SearchChartConfig TODO
type SearchChartConfig struct {
	Count uint64                   `json:"count"`
	Info  map[string][]ChartConfig `json:"info"`
}

// CloudMapping TODO
type CloudMapping struct {
	CreateTime Time   `json:"create_time" bson:"create_time"`
	LastTime   Time   `json:"last_time" bson:"lsat_time"`
	CloudName  string `json:"bk_cloud_name" bson:"bk_cloud_name"`
	OwnerID    string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	CloudID    int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
}

// ChartClassification TODO
type ChartClassification struct {
	Host []ChartConfig `json:"host"`
	Inst []ChartConfig `json:"inst"`
	Nav  []ChartConfig `json:"nav"`
}

// ObjectIDName TODO
type ObjectIDName struct {
	ObjectID   string `json:"bk_object_id"`
	ObjectName string `json:"bk_object_name"`
}

// StatisticInstOperation TODO
type StatisticInstOperation struct {
	Create []StringIDCount   `json:"create"`
	Delete []StringIDCount   `json:"delete"`
	Update []UpdateInstCount `json:"update"`
}

var (
	// BizModuleHostChart TODO
	BizModuleHostChart = ChartConfig{
		ReportType: common.BizModuleHostChart,
	}

	// HostOsChart TODO
	HostOsChart = ChartConfig{
		ReportType: common.HostOSChart,
		Name:       "按操作系统类型统计",
		ObjID:      "host",
		Width:      "50",
		ChartType:  "pie",
		Field:      "bk_os_type",
		XAxisCount: 10,
	}

	// HostBizChart TODO
	HostBizChart = ChartConfig{
		ReportType: common.HostBizChart,
		Name:       "按业务统计",
		ObjID:      "host",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	// HostCloudChart TODO
	HostCloudChart = ChartConfig{
		ReportType: common.HostCloudChart,
		Name:       "按管控区域统计",
		Width:      "100",
		ObjID:      "host",
		ChartType:  "bar",
		Field:      common.BKCloudIDField,
		XAxisCount: 20,
	}

	// HostChangeBizChart TODO
	HostChangeBizChart = ChartConfig{
		ReportType: common.HostChangeBizChart,
		Name:       "主机数量变化趋势",
		Width:      "100",
		XAxisCount: 20,
	}

	// ModelAndInstCountChart TODO
	ModelAndInstCountChart = ChartConfig{
		ReportType: common.ModelAndInstCount,
	}

	// ModelInstChart TODO
	ModelInstChart = ChartConfig{
		ReportType: common.ModelInstChart,
		Name:       "实例数量统计",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	// ModelInstChangeChart TODO
	ModelInstChangeChart = ChartConfig{
		ReportType: common.ModelInstChangeChart,
		Name:       "实例变更统计",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	// InnerChartsMap TODO
	InnerChartsMap = map[string]ChartConfig{
		common.BizModuleHostChart:   BizModuleHostChart,
		common.ModelAndInstCount:    ModelAndInstCountChart,
		common.HostOSChart:          HostOsChart,
		common.HostBizChart:         HostBizChart,
		common.HostCloudChart:       HostCloudChart,
		common.HostChangeBizChart:   HostChangeBizChart,
		common.ModelInstChart:       ModelInstChart,
		common.ModelInstChangeChart: ModelInstChangeChart,
	}

	// InnerChartsArr TODO
	InnerChartsArr = []string{
		common.BizModuleHostChart,
		common.ModelAndInstCount,
		common.HostOSChart,
		common.HostBizChart,
		common.HostCloudChart,
		common.HostChangeBizChart,
		common.ModelInstChart,
		common.ModelInstChangeChart,
	}
)
