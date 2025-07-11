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
	"encoding/json"
	"net/http"

	"github.com/emicklei/go-restful/v3"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
)

// SearchCollector TODO
func (s *Service) SearchCollector(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(pheader))

	cond := metadata.ParamNetcollectorSearch{}
	if err := json.NewDecoder(req.Request.Body).Decode(&cond); err != nil {
		blog.Errorf("[NetDevice][SearchCollector] decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	blog.Infof("[NetDevice][SearchCollector] search by %+v", cond)

	count, result, err := s.logics.SearchCollector(pheader, cond)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCollectNetCollectorSearchFail)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(metadata.RspNetcollectorSearch{
		Count: count,
		Info:  result,
	}))
	return
}

// UpdateCollector TODO
func (s *Service) UpdateCollector(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(pheader))

	cond := metadata.Netcollector{}
	if err := json.NewDecoder(req.Request.Body).Decode(&cond); err != nil {
		blog.Errorf("[NetDevice][UpdateCollector] decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	blog.Infof("[NetDevice][UpdateCollector] update by %+v", cond)

	if cond.BizID <= 0 {
		resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsLostField, common.BKAppIDField)})
		return
	}
	if len(cond.InnerIP) <= 0 {
		resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsLostField, common.BKHostInnerIPField)})
		return
	}

	err := s.logics.UpdateCollector(pheader, cond)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCollectNetCollectorUpdateFail)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
	return
}

// DiscoverNetDevice TODO
func (s *Service) DiscoverNetDevice(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(pheader))

	cond := metadata.ParamNetcollectDiscover{}
	if err := json.NewDecoder(req.Request.Body).Decode(&cond); err != nil {
		blog.Errorf("[NetDevice][DiscoverNetDevice] decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	blog.Infof("[NetDevice][DiscoverNetDevice] discover by %+v", cond)

	err := s.logics.DiscoverNetDevice(pheader, cond.Collectors)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCollectNetCollectorDiscoverFail)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
	return
}
