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

package hostapplyrule

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CreateHostApplyRule TODO
func (p *hostApplyRule) CreateHostApplyRule(ctx context.Context, header http.Header, bizID int64,
	option metadata.CreateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp `json:",inline"`
		Data              metadata.HostApplyRule `json:"data"`
	}{}

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/create/host_apply_rule/bk_biz_id/%d/", bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("CreateHostApplyRule failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

// UpdateHostApplyRule TODO
func (p *hostApplyRule) UpdateHostApplyRule(ctx context.Context, header http.Header, bizID int64, ruleID int64,
	option metadata.UpdateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp `json:",inline"`
		Data              metadata.HostApplyRule `json:"data"`
	}{}

	err := p.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/update/host_apply_rule/%d/bk_biz_id/%d/", ruleID, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

// DeleteHostApplyRule TODO
func (p *hostApplyRule) DeleteHostApplyRule(ctx context.Context, header http.Header, bizID int64,
	option metadata.DeleteHostApplyRuleOption) errors.CCErrorCoder {
	ret := struct {
		metadata.BaseResp `json:",inline"`
	}{}

	err := p.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/deletemany/host_apply_rule/bk_biz_id/%d/", bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("DeleteHostApplyRule failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

// GetHostApplyRule TODO
func (p *hostApplyRule) GetHostApplyRule(ctx context.Context, header http.Header, bizID int64,
	ruleID int64) (metadata.HostApplyRule, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp `json:",inline"`
		Data              metadata.HostApplyRule `json:"data"`
	}{}

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef("/find/host_apply_rule/%d/bk_biz_id/%d/", ruleID, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("GetHostApplyRule failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

// ListHostApplyRule search host apply rule
func (p *hostApplyRule) ListHostApplyRule(ctx context.Context, header http.Header, bizID int64,
	option metadata.ListHostApplyRuleOption) (metadata.MultipleHostApplyRuleResult, errors.CCErrorCoder) {

	ret := struct {
		metadata.BaseResp
		Data metadata.MultipleHostApplyRuleResult `json:"data"`
	}{}

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/host_apply_rule/bk_biz_id/%d/", bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		return ret.Data, errors.CCHttpError
	}

	if ccErr := ret.CCError(); ccErr != nil {
		return ret.Data, ccErr
	}

	return ret.Data, nil
}

// BatchUpdateHostApplyRule TODO
func (p *hostApplyRule) BatchUpdateHostApplyRule(ctx context.Context, header http.Header, bizID int64,
	option metadata.BatchCreateOrUpdateApplyRuleOption) (metadata.BatchCreateOrUpdateHostApplyRuleResult,
	errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data metadata.BatchCreateOrUpdateHostApplyRuleResult `json:"data"`
	}{}

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/updatemany/host_apply_rule/bk_biz_id/%d/", bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("BatchUpdateHostApplyRule failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

// GenerateApplyPlan TODO
func (p *hostApplyRule) GenerateApplyPlan(ctx context.Context, header http.Header, bizID int64,
	option metadata.HostApplyPlanOption) (metadata.HostApplyPlanResult, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data metadata.HostApplyPlanResult `json:"data"`
	}{}

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/host_apply_plan/bk_biz_id/%d/", bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("GenerateApplyPlan failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

// SearchRuleRelatedModules TODO
func (p *hostApplyRule) SearchRuleRelatedModules(ctx context.Context, header http.Header, bizID int64,
	option metadata.SearchRuleRelatedModulesOption) ([]metadata.Module, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data []metadata.Module `json:"data"`
	}{
		BaseResp: metadata.BaseResp{},
		Data:     make([]metadata.Module, 0),
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/modules/bk_biz_id/%d/host_apply_rule_related", bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("SearchRuleRelatedModules failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

// RunHostApplyOnHosts TODO
func (p *hostApplyRule) RunHostApplyOnHosts(ctx context.Context, header http.Header, bizID int64,
	option metadata.UpdateHostByHostApplyRuleOption) (metadata.MultipleHostApplyResult, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data metadata.MultipleHostApplyResult `json:"data"`
	}{}

	err := p.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/updatemany/host/bk_biz_id/%d/update_by_host_apply", bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("RunHostApplyOnHosts failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}
	return ret.Data, nil
}

// SearchRuleRelatedServiceTemplates search rule related service templates
func (p *hostApplyRule) SearchRuleRelatedServiceTemplates(ctx context.Context, header http.Header,
	option *metadata.RuleRelatedServiceTemplateOption) ([]metadata.SrvTemplate, errors.CCErrorCoder) {

	resp := new(metadata.ServiceTemplatesResponse)

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/service_templates/host_apply_rule_related").
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return resp.Data, errors.CCHttpError
	}
	if resp.CCError() != nil {
		return resp.Data, resp.CCError()
	}

	return resp.Data, nil
}
