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

package host

import (
	"context"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// SearchHostWithInnerIPForStatic Note: This function is only used to query the host through ip+cloud in the static IP
// scenario of the host snapshot !!!
func (b *baseCache) SearchHostWithInnerIPForStatic(ctx context.Context, h http.Header,
	opt *metadata.SearchHostWithInnerIPOption) (
	jsonString string, err error) {

	resp, err := b.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/find/cache/host/with_inner_ip").
		WithHeaders(h).
		Do().
		IntoJsonString()

	if err != nil {
		return "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data, nil
}

// SearchHostWithAgentID Note: find host information by agentID only for data collection api !!!.
func (b *baseCache) SearchHostWithAgentID(ctx context.Context, h http.Header, opt *metadata.SearchHostWithAgentID) (
	jsonString string, err error) {

	resp, err := b.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/find/cache/host/with_agent_id").
		WithHeaders(h).
		Do().
		IntoJsonString()

	if err != nil {
		return "", err
	}

	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data, nil
}

// SearchHostWithHostID TODO
func (b *baseCache) SearchHostWithHostID(ctx context.Context, h http.Header,
	opt *metadata.SearchHostWithIDOption) (jsonString string, err error) {

	resp, err := b.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/find/cache/host/with_host_id").
		WithHeaders(h).
		Do().
		IntoJsonString()

	if err != nil {
		return "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data, nil
}

// ListHostWithPage list hosts with page or id list, and returned with a json array string
func (b *baseCache) ListHostWithPage(ctx context.Context, h http.Header, opt *metadata.ListHostWithPage) (cnt int64,
	jsonArray string, err error) {

	resp, err := b.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/cache/host/with_page").
		WithHeaders(h).
		Do().
		IntoJsonCntInfoString()

	if err != nil {
		return 0, "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return 0, "", errors.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data.Count, resp.Data.Info, nil
}

// ListHostWithHostID TODO
func (b *baseCache) ListHostWithHostID(ctx context.Context, h http.Header,
	opt *metadata.ListWithIDOption) (jsonString string, err error) {

	resp, err := b.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/cache/host/with_host_id").
		WithHeaders(h).
		Do().
		IntoJsonString()

	if err != nil {
		return "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data, nil
}
