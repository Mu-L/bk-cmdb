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

package mainline

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
)

// SearchMainlineModelTopo TODO
func (m *mainline) SearchMainlineModelTopo(ctx context.Context, header http.Header,
	withDetail bool) (*metadata.TopoModelNode, errors.CCErrorCoder) {
	rid := httpheader.GetRid(header)
	ret := new(metadata.SearchTopoModelNodeResult)
	// resp = new(metadata.TopoModelNode)
	subPath := "/read/mainline/model"

	input := map[string]bool{}
	input["with_detail"] = withDetail

	err := m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("SearchMainlineModelTopo failed, http failed, err: %s, rid: %s", err.Error(), rid)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

// SearchMainlineInstanceTopo TODO
func (m *mainline) SearchMainlineInstanceTopo(ctx context.Context, header http.Header, bkBizID int64,
	withDetail bool) (*metadata.TopoInstanceNode, errors.CCErrorCoder) {
	rid := httpheader.GetRid(header)
	input := map[string]bool{}
	input["with_detail"] = withDetail

	ret := new(metadata.SearchTopoInstanceNodeResult)
	subPath := "/read/mainline/instance/%d"
	err := m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, bkBizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("SearchMainlineInstanceTopo failed, http failed, err: %s, rid: %s", err.Error(), rid)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}
