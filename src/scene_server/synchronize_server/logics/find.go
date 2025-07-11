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

package logics

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

const (
	defaultLimit = 100
)

func (lgc *Logics) findInstance(ctx context.Context, objID string, input *metadata.QueryCondition) (
	*metadata.InstDataInfo, error) {
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, objID, input)
	if err != nil {
		blog.Errorf("FindInstance ReadInstance http do error, error: %s,input:  %#v,rid:%s", err.Error(), input,
			lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return result, nil
}

func (lgc *Logics) find(ctx context.Context, input *metadata.SynchronizeFindInfoParameter) (*metadata.InstDataInfo,
	errors.CCError) {
	result, err := lgc.CoreAPI.CoreService().Synchronize().SynchronizeFind(ctx, lgc.header, input)
	if err != nil {
		blog.Errorf("find http do error. err:%s,input:%#v,rid:%s", err.Error(), input, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("find http reply error. err code:%d,err msg:%s,input:%#v,rid:%s", result.Code, result.ErrMsg, input,
			lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	return &result.Data, nil
}

// Find TODO
func (lgc *Logics) Find(ctx context.Context, input *metadata.SynchronizeFindInfoParameter) (*metadata.InstDataInfo,
	errors.CCError) {
	switch input.DataType {
	case metadata.SynchronizeOperateDataTypeInstance:
		return lgc.findInstance(ctx, input.DataClassify, SynchronizeFindInfoParameterToQuerycondition(input))
	case metadata.SynchronizeOperateDataTypeAssociation:
		return lgc.find(ctx, input)
	case metadata.SynchronizeOperateDataTypeModel:
		// cancel limit
		// input.Limit = 0
		return lgc.find(ctx, input)
	}
	blog.Warnf("Find not found, input:%#v,rid:%s", input, lgc.rid)
	return nil, nil
}

// SynchronizeFindInfoParameterToQuerycondition  SynchronizeFindInfoParameter to Querycondition
func SynchronizeFindInfoParameterToQuerycondition(input *metadata.SynchronizeFindInfoParameter) *metadata.QueryCondition {
	ret := &metadata.QueryCondition{
		Page:      metadata.BasePage{Limit: int(input.Limit), Start: int(input.Start)},
		Condition: input.Condition,
	}
	if ret.Page.Limit <= 0 {
		//  limit
		ret.Page.Limit = defaultLimit
	}
	return ret
}
