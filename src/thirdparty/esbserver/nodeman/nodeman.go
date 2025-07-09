// Package nodeman TODO
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
package nodeman

import (
	"context"
	"net/http"

	"configcenter/src/thirdparty/esbserver/esbutil"
)

// SearchPackage TODO
func (p *nodeman) SearchPackage(ctx context.Context, h http.Header,
	processname string) (resp *SearchPluginPackageResult, err error) {
	resp = new(SearchPluginPackageResult)
	subPath := "/v2/nodeman/%s/package/"
	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, processname).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Peek().
		Do().
		Into(resp)
	return
}

// SearchProcess TODO
func (p *nodeman) SearchProcess(ctx context.Context, h http.Header,
	processname string) (resp *SearchPluginProcessResult, err error) {
	resp = new(SearchPluginProcessResult)
	subPath := "/v2/nodeman/process/%s/"
	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, processname).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Peek().
		Do().
		Into(resp)
	return
}

// SearchProcessInfo TODO
func (p *nodeman) SearchProcessInfo(ctx context.Context, h http.Header,
	processname string) (resp *SearchPluginProcessInfoResult, err error) {
	resp = new(SearchPluginProcessInfoResult)
	subPath := "/v2/nodeman/process_info/%s/"
	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, processname).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Peek().
		Do().
		Into(resp)
	return
}

// UpgradePlugin TODO
func (p *nodeman) UpgradePlugin(ctx context.Context, h http.Header, bizID string,
	data *UpgradePluginRequest) (resp *UpgradePluginResult, err error) {
	resp = new(UpgradePluginResult)
	subPath := "/v2/nodeman/%s/tasks/"

	params := esbUpgradePluginParams{
		EsbCommParams:        esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		UpgradePluginRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, bizID).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Peek().
		Do().
		Into(resp)
	return
}

// SearchTask TODO
func (p *nodeman) SearchTask(ctx context.Context, h http.Header, bizID int64, taskID int64) (resp *SearchTaskResult,
	err error) {
	resp = new(SearchTaskResult)
	subPath := "/v2/nodeman/%d/tasks/%d/"
	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, bizID, taskID).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Peek().
		Do().
		Into(resp)
	return
}

// SearchPluginHost TODO
func (p *nodeman) SearchPluginHost(ctx context.Context, h http.Header,
	processname string) (resp *SearchPluginHostResult, err error) {
	resp = new(SearchPluginHostResult)
	subPath := "/v2/nodeman/0/host_status/get_host/"
	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath).
		WithParam("name", processname).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Peek().
		Do().
		Into(resp)
	return
}
