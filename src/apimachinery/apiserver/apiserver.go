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

// Package apiserver TODO
package apiserver

import (
	"context"
	"fmt"
	"net/http"

	fieldtmpl "configcenter/src/apimachinery/apiserver/field_template"
	modelquote "configcenter/src/apimachinery/apiserver/model_quote"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/transaction"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// ApiServerClientInterface TODO
type ApiServerClientInterface interface {
	Client() rest.ClientInterface
	ModelQuote() modelquote.Interface
	FieldTemplate() fieldtmpl.Interface
	Txn() transaction.Interface

	AddDefaultApp(ctx context.Context, h http.Header, ownerID string, params mapstr.MapStr) (resp *metadata.Response,
		err error)
	SearchDefaultApp(ctx context.Context, h http.Header, ownerID string) (resp *metadata.QueryInstResult, err error)
	GetObjectData(ctx context.Context, h http.Header, cond *metadata.ExportObjectCondition) (
		resp *metadata.ObjectAttrBatchResult, err error)
	GetInstDetail(ctx context.Context, h http.Header, objID string, params mapstr.MapStr) (
		resp *metadata.QueryInstResult, err error)
	GetInstUniqueFields(ctx context.Context, h http.Header, objID string, uniqueID int64, params mapstr.MapStr) (
		resp metadata.QueryUniqueFieldsResult, err error)
	CreateObjectAtt(ctx context.Context, h http.Header, obj *metadata.ObjAttDes) (resp *metadata.Response, err error)
	UpdateObjectAtt(ctx context.Context, objID string, h http.Header,
		data map[string]interface{}) (resp *metadata.Response, err error)
	DeleteObjectAtt(ctx context.Context, objID string, h http.Header) (resp *metadata.Response, err error)
	GetObjectAttr(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ObjectAttrResult, err error)
	GetHostData(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error)
	ListHostWithoutApp(ctx context.Context, h http.Header,
		option metadata.ListHostsWithNoBizParameter) (resp *metadata.ListHostWithoutAppResponse, err error)
	GetObjectGroup(ctx context.Context, h http.Header, ownerID, objID string,
		params mapstr.MapStr) (resp *metadata.ObjectAttrGroupResult, err error)
	AddHost(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error)
	AddHostByExcel(ctx context.Context, h http.Header, params mapstr.MapStr) (*metadata.ImportInstRes, error)
	UpdateHost(ctx context.Context, h http.Header, params mapstr.MapStr) (*metadata.ImportInstRes, error)
	GetHostModuleRelation(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.HostModuleResp,
		err error)
	AddInst(ctx context.Context, h http.Header, ownerID, objID string,
		params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error)
	AddInstByImport(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (
		*metadata.ImportInstRes, error)
	AddObjectBatch(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.Response, err error)
	SearchAssociationInst(ctx context.Context, h http.Header,
		request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error)
	ImportAssociation(ctx context.Context, h http.Header, objID string,
		input *metadata.RequestImportAssociation) (resp *metadata.ResponeImportAssociation, err error)

	SearchNetCollectDevice(ctx context.Context, h http.Header,
		cond condition.Condition) (resp *metadata.ResponseInstData, err error)
	SearchNetDeviceProperty(ctx context.Context, h http.Header,
		cond condition.Condition) (resp *metadata.ResponseInstData, err error)
	SearchNetCollectDeviceBatch(ctx context.Context, h http.Header,
		cond mapstr.MapStr) (resp *metadata.ResponseInstData, err error)
	SearchNetDevicePropertyBatch(ctx context.Context, h http.Header,
		cond mapstr.MapStr) (resp *metadata.ResponseInstData, err error)

	CreateBiz(ctx context.Context, ownerID string, h http.Header, dat map[string]interface{}) (
		resp *metadata.CreateInstResult, err error)
	UpdateBiz(ctx context.Context, ownerID string, bizID string, h http.Header, data map[string]interface{}) (
		resp *metadata.Response, err error)
	UpdateBizDataStatus(ctx context.Context, ownerID string, flag common.DataStatusFlag, bizID int64,
		h http.Header) errors.CCErrorCoder
	UpdateBizPropertyBatch(ctx context.Context, h http.Header, param metadata.UpdateBizPropertyBatchParameter) (
		resp *metadata.Response, err error)
	DeleteBiz(ctx context.Context, h http.Header, param metadata.DeleteBizParam) error
	SearchBiz(ctx context.Context, ownerID string, h http.Header, param *metadata.QueryBusinessRequest) (
		resp *metadata.SearchInstResult, err error)

	ReadModuleAssociation(ctx context.Context, h http.Header, cond *metadata.QueryCondition) (*metadata.AsstResult,
		errors.CCErrorCoder)
	ReadModel(ctx context.Context, h http.Header, input *metadata.QueryCondition) (*metadata.QueryModelDataResult,
		errors.CCErrorCoder)
	ReadModelForUI(ctx context.Context, h http.Header, input *metadata.QueryCondition) (*metadata.QueryModelDataResult,
		errors.CCErrorCoder)
	ReadInstance(ctx context.Context, h http.Header, objID string, input *metadata.QueryCondition) (
		resp *metadata.QueryConditionResult, err error)
	SearchObjectUnique(ctx context.Context, objID string, h http.Header) (resp *metadata.SearchUniqueResult, err error)

	FindAssociationByObjectAssociationID(ctx context.Context, h http.Header, objID string,
		input metadata.FindAssociationByObjectAssociationIDRequest) (
		resp *metadata.FindAssociationByObjectAssociationIDResponse, err error)

	SearchObjectAssociation(ctx context.Context, h http.Header,
		request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error)

	SearchObjectWithTotalInfo(ctx context.Context, h http.Header, params *metadata.BatchExportObject) (
		*metadata.TotalObjectInfo, error)
	CreateManyObject(ctx context.Context, h http.Header, params metadata.ImportObjects) ([]metadata.Object, error)

	SearchCloudArea(ctx context.Context, h http.Header, params metadata.CloudAreaSearchParam) (
		*metadata.SearchDataResult, error)

	SearchPlatformSetting(ctx context.Context, h http.Header, status string) (resp *metadata.PlatformSettingResult,
		err error)

	CountObjectInstances(ctx context.Context, h http.Header, objID string, input *metadata.CommonCountFilter) (
		*metadata.CommonCountResult, errors.CCErrorCoder)
	CountObjInstByFilters(ctx context.Context, h http.Header, objID string, filters []map[string]interface{}) (
		[]int64, errors.CCErrorCoder)
	GroupRelResByIDs(ctx context.Context, h http.Header, kind metadata.GroupByResKind,
		opt *metadata.GroupRelResByIDsOption) (map[int64][]interface{}, errors.CCErrorCoder)

	HealthCheck() (bool, error)
	SearchProject(ctx context.Context, h http.Header, params *metadata.SearchProjectOption) (resp *metadata.InstResult,
		err error)
}

// NewApiServerClientInterface TODO
func NewApiServerClientInterface(c *util.Capability, version string) ApiServerClientInterface {
	base := fmt.Sprintf("/api/%s", version)
	return &apiServer{
		client: rest.NewRESTClient(c, base),
	}
}

// NewWrappedApiServerClientI new wrapped api server client interface by restful client
func NewWrappedApiServerClientI(client rest.ClientInterface, wrappers ...rest.RequestWrapper) ApiServerClientInterface {
	return &apiServer{
		client: rest.NewClientWrapper(client, wrappers...),
	}
}

type apiServer struct {
	client rest.ClientInterface
}

// ModelQuote return the model quote client
func (a *apiServer) ModelQuote() modelquote.Interface {
	return modelquote.New(a.client)
}

// FieldTemplate return the field template client
func (a *apiServer) FieldTemplate() fieldtmpl.Interface {
	return fieldtmpl.New(a.client)
}

// Txn returns transaction client
func (a *apiServer) Txn() transaction.Interface {
	return transaction.NewTxn(a.client)
}
