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

package auth

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

// SearchAuthResource TODO
func (a *auth) SearchAuthResource(ctx context.Context, h http.Header,
	param metadata.PullResourceParam) (metadata.PullResourceResponse, error) {
	resp := metadata.PullResourceResponse{}
	subPath := "/search/auth/resource"

	err := a.client.Post().
		WithContext(ctx).
		Body(param).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return resp, err
}
