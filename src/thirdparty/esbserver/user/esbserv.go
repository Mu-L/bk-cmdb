// Package user TODO
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
package user

import (
	"context"
	"net/http"
	"net/url"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/esbserver/esbutil"
)

// UserClientInterface TODO
type UserClientInterface interface {
	GetAllUsers(ctx context.Context, h http.Header) (resp *metadata.EsbUserListResponse, err error)
	ListUsers(ctx context.Context, h http.Header, params map[string]string) (resp *metadata.EsbListUserResponse,
		err error)
	GetDepartment(ctx context.Context, h http.Header, u *url.URL) (resp *metadata.EsbDepartmentResponse,
		err error)
	GetAllDepartment(ctx context.Context, h http.Header,
		params map[string]string) (resp *metadata.EsbDepartmentResponse,
		err error)
	GetDepartmentProfile(ctx context.Context, q *http.Request) (resp *metadata.EsbDepartmentProfileResponse, err error)
}

// NewUserClientInterface TODO
func NewUserClientInterface(client rest.ClientInterface, config *esbutil.EsbConfigSrv) UserClientInterface {
	return &user{
		client: client,
		config: config,
	}
}

type user struct {
	config *esbutil.EsbConfigSrv
	client rest.ClientInterface
}
