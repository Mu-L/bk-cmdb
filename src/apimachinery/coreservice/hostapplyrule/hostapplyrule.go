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

// Package hostapplyrule TODO
package hostapplyrule

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// HostApplyRuleInterface TODO
type HostApplyRuleInterface interface {
	CreateHostApplyRule(ctx context.Context, header http.Header, bizID int64,
		option metadata.CreateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder)
	UpdateHostApplyRule(ctx context.Context, header http.Header, bizID int64, ruleID int64,
		option metadata.UpdateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder)
	GetHostApplyRule(ctx context.Context, header http.Header, bizID int64, ruleID int64) (metadata.HostApplyRule,
		errors.CCErrorCoder)
	ListHostApplyRule(ctx context.Context, header http.Header, bizID int64, option metadata.ListHostApplyRuleOption) (
		metadata.MultipleHostApplyRuleResult, errors.CCErrorCoder)
	DeleteHostApplyRule(ctx context.Context, header http.Header, bizID int64,
		option metadata.DeleteHostApplyRuleOption) errors.CCErrorCoder
	GenerateApplyPlan(ctx context.Context, header http.Header, bizID int64, option metadata.HostApplyPlanOption) (
		metadata.HostApplyPlanResult, errors.CCErrorCoder)
	SearchRuleRelatedModules(ctx context.Context, header http.Header, bizID int64,
		option metadata.SearchRuleRelatedModulesOption) ([]metadata.Module, errors.CCErrorCoder)
	BatchUpdateHostApplyRule(ctx context.Context, header http.Header, bizID int64,
		option metadata.BatchCreateOrUpdateApplyRuleOption) (metadata.BatchCreateOrUpdateHostApplyRuleResult,
		errors.CCErrorCoder)

	RunHostApplyOnHosts(ctx context.Context, header http.Header, bizID int64,
		option metadata.UpdateHostByHostApplyRuleOption) (metadata.MultipleHostApplyResult, errors.CCErrorCoder)
	SearchRuleRelatedServiceTemplates(ctx context.Context, header http.Header,
		option *metadata.RuleRelatedServiceTemplateOption) ([]metadata.SrvTemplate, errors.CCErrorCoder)
}

// NewHostApplyRuleClient TODO
func NewHostApplyRuleClient(client rest.ClientInterface) HostApplyRuleInterface {
	return &hostApplyRule{client: client}
}

type hostApplyRule struct {
	client rest.ClientInterface
}
