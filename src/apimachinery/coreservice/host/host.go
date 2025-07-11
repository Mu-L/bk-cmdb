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

// Package host TODO
package host

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// HostClientInterface TODO
type HostClientInterface interface {
	TransferToInnerModule(ctx context.Context, h http.Header, input *metadata.TransferHostToInnerModule) (
		[]metadata.ExceptionResult, errors.CCErrorCoder)
	TransferToNormalModule(ctx context.Context, header http.Header, input *metadata.HostsModuleRelation) (
		[]metadata.ExceptionResult, errors.CCErrorCoder)
	TransferToAnotherBusiness(ctx context.Context, header http.Header,
		input *metadata.TransferHostsCrossBusinessRequest) (resp *metadata.OperaterException, err error)

	RemoveFromModule(ctx context.Context, header http.Header, input *metadata.RemoveHostsFromModuleOption) (
		resp *metadata.OperaterException, err error)
	DeleteHostFromSystem(ctx context.Context, header http.Header, input *metadata.DeleteHostRequest) errors.CCErrorCoder

	GetHostModuleRelation(ctx context.Context, header http.Header, input *metadata.HostModuleRelationRequest) (
		*metadata.HostConfigData, error)
	FindIdentifier(ctx context.Context, header http.Header, input *metadata.SearchHostIdentifierParam) (
		*metadata.SearchHostIdentifierData, error)

	GetHostByID(ctx context.Context, header http.Header, hostID int64) (resp *metadata.HostInstanceResult, err error)
	GetHosts(ctx context.Context, header http.Header, opt *metadata.QueryInput) (*metadata.HostInfo, error)
	LockHost(ctx context.Context, header http.Header, input *metadata.HostLockRequest) (resp *metadata.HostLockResponse,
		err error)
	UnlockHost(ctx context.Context, header http.Header, input *metadata.HostLockRequest) (
		resp *metadata.HostLockResponse, err error)
	QueryHostLock(ctx context.Context, header http.Header, input *metadata.QueryHostLockRequest) (
		resp *metadata.HostLockQueryResponse, err error)

	// CreateDynamicGroup TODO
	// dynamic grouping interfaces.
	CreateDynamicGroup(ctx context.Context, header http.Header, data *metadata.DynamicGroup) (resp *metadata.IDResult,
		err error)
	UpdateDynamicGroup(ctx context.Context, bizID, id string, header http.Header, data map[string]interface{}) (
		resp *metadata.BaseResp, err error)
	DeleteDynamicGroup(ctx context.Context, bizID, id string, header http.Header) (resp *metadata.BaseResp, err error)
	GetDynamicGroup(ctx context.Context, bizID, id string, header http.Header) (resp *metadata.GetDynamicGroupResult,
		err error)
	SearchDynamicGroup(ctx context.Context, header http.Header, opt *metadata.QueryCondition) (
		resp *metadata.SearchDynamicGroupResult, err error)

	AddUserCustom(ctx context.Context, user string, h http.Header, dat map[string]interface{}) (resp *metadata.BaseResp,
		err error)
	UpdateUserCustomByID(ctx context.Context, user string, id string, h http.Header, dat map[string]interface{}) (
		resp *metadata.BaseResp, err error)
	GetUserCustomByUser(ctx context.Context, user string, h http.Header) (resp *metadata.GetUserCustomResult, err error)
	GetDefaultUserCustom(ctx context.Context, header http.Header) (resp *metadata.GetUserCustomResult, err error)
	UpdateDefaultUserCustom(ctx context.Context, header http.Header, dat map[string]interface{}) (
		resp *metadata.BaseResp, err error)

	AddHostFavourite(ctx context.Context, user string, h http.Header, dat *metadata.FavouriteParms) (
		resp *metadata.IDResult, err error)
	UpdateHostFavouriteByID(ctx context.Context, user string, id string, h http.Header, dat map[string]interface{}) (
		resp *metadata.BaseResp, err error)
	DeleteHostFavouriteByID(ctx context.Context, user string, id string, h http.Header) (resp *metadata.BaseResp,
		err error)
	ListHostFavourites(ctx context.Context, user string, h http.Header, dat *metadata.QueryInput) (
		resp *metadata.GetHostFavoriteResult, err error)
	GetHostFavouriteByID(ctx context.Context, user string, id string, h http.Header) (
		resp *metadata.GetHostFavoriteWithIDResult, err error)

	GetHostModulesIDs(ctx context.Context, h http.Header, dat *metadata.ModuleHostConfigParams) (
		resp *metadata.GetHostModuleIDsResult, err error)

	// ListHosts TODO
	// search host
	ListHosts(ctx context.Context, header http.Header, option *metadata.ListHosts) (resp *metadata.ListHostResult,
		err error)

	// UpdateHostCloudAreaField TODO
	// update host's cloud area field
	UpdateHostCloudAreaField(ctx context.Context, header http.Header,
		option metadata.UpdateHostCloudAreaFieldOption) errors.CCErrorCoder

	// FindCloudAreaHostCount find host count in every cloudarea
	FindCloudAreaHostCount(ctx context.Context, header http.Header, option metadata.CloudAreaHostCount) (
		resp *metadata.CloudAreaHostCountResult, err error)

	GetDistinctHostIDByTopology(ctx context.Context, header http.Header,
		input *metadata.DistinctHostIDByTopoRelationRequest) ([]int64, errors.CCErrorCoder)

	TransferHostResourceDirectory(ctx context.Context, header http.Header,
		option *metadata.TransferHostResourceDirectory) errors.CCErrorCoder
}

// NewHostClientInterface TODO
func NewHostClientInterface(client rest.ClientInterface) HostClientInterface {
	return &host{client: client}
}

type host struct {
	client rest.ClientInterface
}
