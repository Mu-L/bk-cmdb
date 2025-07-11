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

package service

import (
	"net/http"

	"configcenter/src/common/http/rest"

	"github.com/emicklei/go-restful/v3"
)

func (s *Service) initSetTemplate(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/topo/set_template/bk_biz_id/{bk_biz_id}/",
		Handler: s.CreateSetTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/topo/set_template/all_info",
		Handler: s.CreateSetTemplateAllInfo})
	utility.AddHandler(rest.Action{Verb: http.MethodPut,
		Path: "/update/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/", Handler: s.UpdateSetTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/topo/set_template/all_info",
		Handler: s.UpdateSetTemplateAllInfo})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete,
		Path: "/deletemany/topo/set_template/bk_biz_id/{bk_biz_id}/", Handler: s.DeleteSetTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodGet,
		Path: "/find/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/", Handler: s.GetSetTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/topo/set_template/all_info",
		Handler: s.GetSetTemplateAllInfo})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/topo/set_template/bk_biz_id/{bk_biz_id}/",
		Handler: s.ListSetTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/findmany/topo/set_template/bk_biz_id/{bk_biz_id}/web/", Handler: s.ListSetTemplateWeb})
	utility.AddHandler(rest.Action{Verb: http.MethodGet,
		Path:    "/findmany/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/service_templates",
		Handler: s.ListSetTplRelatedSvcTpl})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/sets/web",
		Handler: s.ListSetTplRelatedSetsWeb})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/diff_with_instances",
		Handler: s.DiffSetTplWithInst})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/updatemany/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/sync_to_instances",
		Handler: s.SyncSetTplToInst})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/instances_sync_status",
		Handler: s.GetSetSyncDetails})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/findmany/topo/set_template_sync_status/bk_biz_id/{bk_biz_id}", Handler: s.ListSetTemplateSyncStatus})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/findmany/topo/set_template_sync_history/bk_biz_id/{bk_biz_id}", Handler: s.ListSetTemplateSyncHistory})
	utility.AddHandler(rest.Action{Verb: http.MethodGet,
		Path:    "/findmany/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/set_template_status",
		Handler: s.CheckSetInstUpdateToDateStatus})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/topo/set_template/bk_biz_id/{bk_biz_id}/set_template_status",
		Handler: s.BatchCheckSetInstUpdateToDateStatus})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/host_with_instances",
		Handler: s.SetWithHostFlag})

	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/topo/set_template/attribute",
		Handler: s.UpdateSetTemplateAttribute})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/topo/set_template/attribute",
		Handler: s.DeleteSetTemplateAttribute})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/topo/set_template/attribute",
		Handler: s.ListSetTemplateAttribute})

	utility.AddToRestfulWebService(web)
}
