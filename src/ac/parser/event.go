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

package parser

import (
	"fmt"
	"net/http"
	"regexp"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/watch"

	"github.com/tidwall/gjson"
)

func (ps *parseStream) eventRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.watch().
		syncHostIdentifier().
		pushHostIdentifier().
		findHostIdentifierPushResult()
	return ps
}

var (
	watchResourceRegexp = regexp.MustCompile(`^/api/v3/event/watch/resource/\S+/?$`)
)

func (ps *parseStream) watch() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// watch resource.
	if ps.hitRegexp(watchResourceRegexp, http.MethodPost) {
		resource := ps.RequestCtx.Elements[5]
		if len(resource) == 0 {
			ps.err = fmt.Errorf("watch event resource, but got empty resource: %s", ps.RequestCtx.Elements[5])
			return ps
		}

		if resource == string(watch.HostIdentifier) {
			// redirect host identity resource to host resource in iam.
			resource = string(watch.Host)
		}

		if resource == string(watch.BizSetRelation) {
			// redirect biz set relation resource to biz set resource in iam.
			resource = string(watch.BizSet)
		}

		authResource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Type:   meta.EventWatch,
				Action: meta.Action(resource),
			},
		}

		switch watch.CursorType(resource) {
		case watch.ObjectBase, watch.MainlineInstance, watch.InstAsst:
			body, err := ps.RequestCtx.getRequestBody()
			if err != nil {
				ps.err = err
				return ps
			}

			// use sub resource(corresponding to the bk_obj_id of the object) for authorization if it is set
			// if sub resource is not set, verify authorization of the resource(which means all sub resources)
			subResource := gjson.GetBytes(body, "bk_filter."+common.BKSubResourceField)
			if subResource.Exists() {
				model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: subResource.String()})
				if err != nil {
					ps.err = err
					return ps
				}
				authResource.InstanceID = model.ID
			}
		case watch.KubeWorkload:
			body, err := ps.RequestCtx.getRequestBody()
			if err != nil {
				ps.err = err
				return ps
			}

			// use sub resource(corresponding to the kind of the workload) for authorization if it is set
			// if sub resource is not set, verify authorization of the resource(which means all sub resources)
			subResource := gjson.GetBytes(body, "bk_filter."+common.BKSubResourceField)
			if subResource.Exists() {
				authResource.InstanceIDEx = subResource.String()
			}
		}

		ps.Attribute.Resources = append(ps.Attribute.Resources, authResource)
		return ps
	}

	return ps
}

const (
	syncHostIdentifierPattern           = "/api/v3/event/sync/host_identifier"
	pushHostIdentifierPattern           = "/api/v3/event/push/host_identifier"
	findHostIdentifierPushResultPattern = "/api/v3/event/find/host_identifier_push_result"
)

func (ps *parseStream) syncHostIdentifier() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitPattern(syncHostIdentifierPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

func (ps *parseStream) pushHostIdentifier() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitPattern(pushHostIdentifierPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

func (ps *parseStream) findHostIdentifierPushResult() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitPattern(findHostIdentifierPushResultPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}
