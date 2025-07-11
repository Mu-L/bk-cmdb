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

// Package process TODO
package process

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// ProcessInterface TODO
type ProcessInterface interface {
	// CreateServiceCategory TODO
	// service category
	CreateServiceCategory(ctx context.Context, h http.Header, category *metadata.ServiceCategory) (
		*metadata.ServiceCategory, errors.CCErrorCoder)
	GetServiceCategory(ctx context.Context, h http.Header, categoryID int64) (*metadata.ServiceCategory,
		errors.CCErrorCoder)
	UpdateServiceCategory(ctx context.Context, h http.Header, categoryID int64, category *metadata.ServiceCategory) (
		*metadata.ServiceCategory, errors.CCErrorCoder)
	ListServiceCategories(ctx context.Context, h http.Header, option metadata.ListServiceCategoriesOption) (
		*metadata.MultipleServiceCategoryWithStatistics, errors.CCErrorCoder)
	DeleteServiceCategory(ctx context.Context, h http.Header, categoryID int64) errors.CCErrorCoder
	GetDefaultServiceCategory(ctx context.Context, h http.Header) (*metadata.ServiceCategory, errors.CCErrorCoder)

	// CreateServiceTemplate TODO
	// service template
	CreateServiceTemplate(ctx context.Context, h http.Header, template *metadata.ServiceTemplate) (
		*metadata.ServiceTemplate, errors.CCErrorCoder)
	GetServiceTemplate(ctx context.Context, h http.Header, templateID int64) (*metadata.ServiceTemplate,
		errors.CCErrorCoder)
	GetServiceTemplateWithStatistics(ctx context.Context, h http.Header, templateID int64) (
		*metadata.ServiceTemplateWithStatistics, errors.CCErrorCoder)
	ListServiceTemplateDetail(ctx context.Context, h http.Header, bizID int64, templateIDs ...int64) (
		metadata.MultipleServiceTemplateDetail, errors.CCErrorCoder)
	UpdateServiceTemplate(ctx context.Context, h http.Header, templateID int64, template *metadata.ServiceTemplate) (
		*metadata.ServiceTemplate, errors.CCErrorCoder)
	UpdateBatchServiceTemplate(ctx context.Context, h http.Header, option *metadata.UpdateOption) errors.CCErrorCoder

	ListServiceTemplates(ctx context.Context, h http.Header, option *metadata.ListServiceTemplateOption) (
		*metadata.MultipleServiceTemplate, errors.CCErrorCoder)
	DeleteServiceTemplate(ctx context.Context, h http.Header, serviceTemplateID int64) errors.CCErrorCoder

	// CreateProcessTemplate TODO
	// process template
	CreateProcessTemplate(ctx context.Context, h http.Header, template *metadata.ProcessTemplate) (
		*metadata.ProcessTemplate, errors.CCErrorCoder)
	GetProcessTemplate(ctx context.Context, h http.Header, templateID int64) (*metadata.ProcessTemplate,
		errors.CCErrorCoder)
	UpdateProcessTemplate(ctx context.Context, h http.Header, templateID int64, property map[string]interface{}) (
		*metadata.ProcessTemplate, errors.CCErrorCoder)
	ListProcessTemplates(ctx context.Context, h http.Header, option *metadata.ListProcessTemplatesOption) (
		*metadata.MultipleProcessTemplate, errors.CCErrorCoder)
	DeleteProcessTemplate(ctx context.Context, h http.Header, processTemplateID int64) errors.CCErrorCoder
	DeleteProcessTemplateBatch(ctx context.Context, h http.Header, processTemplateIDs []int64) errors.CCErrorCoder

	// CreateServiceInstance TODO
	// service instance
	CreateServiceInstance(ctx context.Context, h http.Header, template *metadata.ServiceInstance) (
		*metadata.ServiceInstance, errors.CCErrorCoder)
	CreateServiceInstances(ctx context.Context, h http.Header, instances []*metadata.ServiceInstance) (
		[]*metadata.ServiceInstance, errors.CCErrorCoder)
	GetServiceInstance(ctx context.Context, h http.Header, serviceInstanceID int64) (*metadata.ServiceInstance,
		errors.CCErrorCoder)
	UpdateServiceInstances(ctx context.Context, h http.Header, bizID int64,
		option *metadata.UpdateServiceInstanceOption) errors.CCErrorCoder
	ListServiceInstance(ctx context.Context, h http.Header, option *metadata.ListServiceInstanceOption) (
		*metadata.MultipleServiceInstance, errors.CCErrorCoder)
	DeleteServiceInstance(ctx context.Context, h http.Header,
		option *metadata.CoreDeleteServiceInstanceOption) errors.CCErrorCoder
	GetBusinessDefaultSetModuleInfo(ctx context.Context, h http.Header, bizID int64) (
		metadata.BusinessDefaultSetModuleInfo, errors.CCErrorCoder)
	ListServiceInstanceDetail(ctx context.Context, h http.Header, option *metadata.ListServiceInstanceDetailOption) (
		*metadata.MultipleServiceInstanceDetail, errors.CCErrorCoder)

	// CreateProcessInstanceRelation TODO
	// process instance relation
	CreateProcessInstanceRelation(ctx context.Context, h http.Header, relation *metadata.ProcessInstanceRelation) (
		*metadata.ProcessInstanceRelation, errors.CCErrorCoder)
	CreateProcessInstanceRelations(ctx context.Context, h http.Header, relations []*metadata.ProcessInstanceRelation) (
		[]*metadata.ProcessInstanceRelation, errors.CCErrorCoder)
	GetProcessInstanceRelation(ctx context.Context, h http.Header, processID int64) (*metadata.ProcessInstanceRelation,
		errors.CCErrorCoder)
	UpdateProcessInstanceRelation(ctx context.Context, h http.Header, processID int64,
		template *metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder)
	ListProcessInstanceRelation(ctx context.Context, h http.Header,
		option *metadata.ListProcessInstanceRelationOption) (*metadata.MultipleProcessInstanceRelation,
		errors.CCErrorCoder)
	DeleteProcessInstanceRelation(ctx context.Context, h http.Header,
		option metadata.DeleteProcessInstanceRelationOption) errors.CCErrorCoder

	RemoveTemplateBindingOnModule(ctx context.Context, h http.Header, moduleID int64) (
		*metadata.RemoveTemplateBoundOnModuleResult, errors.CCErrorCoder)
	ConstructServiceInstanceName(ctx context.Context, h http.Header,
		params *metadata.SrvInstNameParams) errors.CCErrorCoder
	ReconstructServiceInstanceName(ctx context.Context, h http.Header, instanceID int64) errors.CCErrorCoder

	// UpdateServiceTemplateAttribute TODO
	// service template attribute
	UpdateServiceTemplateAttribute(ctx context.Context, h http.Header,
		option *metadata.UpdateServTempAttrOption) errors.CCErrorCoder
	DeleteServiceTemplateAttribute(ctx context.Context, h http.Header,
		option *metadata.DeleteServTempAttrOption) errors.CCErrorCoder
	ListServiceTemplateAttribute(ctx context.Context, h http.Header, option *metadata.ListServTempAttrOption) (
		*metadata.ServTempAttrData, errors.CCErrorCoder)
	CreateServiceTemplateAttrs(ctx context.Context, h http.Header, option *metadata.CreateSvcTempAttrsOption) (
		[]int64, errors.CCErrorCoder)
}

// NewProcessInterfaceClient TODO
func NewProcessInterfaceClient(client rest.ClientInterface) ProcessInterface {
	return &process{client: client}
}

type process struct {
	client rest.ClientInterface
}
