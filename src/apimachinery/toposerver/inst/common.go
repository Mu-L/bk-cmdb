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

package inst

import (
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

// SearchAuditDict TODO
func (t *instanceClient) SearchAuditDict(ctx context.Context, h http.Header) (*metadata.Response, error) {
	resp := new(metadata.Response)
	subPath := "/find/audit_dict"

	err := t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return nil, errors.New(resp.Code, resp.ErrMsg)
	}

	return resp, nil
}

// SearchAuditList TODO
func (t *instanceClient) SearchAuditList(ctx context.Context, h http.Header,
	input *metadata.AuditQueryInput) (*metadata.Response, error) {
	resp := new(metadata.Response)
	subPath := "/findmany/audit_list"

	err := t.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return nil, errors.New(resp.Code, resp.ErrMsg)
	}

	return resp, nil
}

// SearchAuditDetail TODO
func (t *instanceClient) SearchAuditDetail(ctx context.Context, h http.Header,
	input *metadata.AuditDetailQueryInput) (*metadata.Response, error) {
	resp := new(metadata.Response)
	subPath := "/find/audit"

	err := t.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return nil, errors.New(resp.Code, resp.ErrMsg)
	}

	return resp, nil
}

// GetInternalModule TODO
func (t *instanceClient) GetInternalModule(ctx context.Context, ownerID, appID string,
	h http.Header) (resp *metadata.SearchInnterAppTopoResult, err error) {
	resp = new(metadata.SearchInnterAppTopoResult)
	subPath := "/topo/internal/%s/%s"

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, ownerID, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchBriefBizTopo TODO
func (t *instanceClient) SearchBriefBizTopo(ctx context.Context, h http.Header, bizID int64,
	input map[string]interface{}) (resp *metadata.SearchBriefBizTopoResult, err error) {
	resp = new(metadata.SearchBriefBizTopoResult)
	subPath := "/find/topo/tree/brief/biz/%d"

	err = t.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
