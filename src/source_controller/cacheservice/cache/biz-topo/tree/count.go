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

package tree

import (
	"context"

	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
)

// TreeWithCount defines a topology tree type with resource count
type TreeWithCount struct{}

// RearrangeBizTopo rearrange business topology tree
func (t *TreeWithCount) RearrangeBizTopo(_ context.Context, topo *types.BizTopo, _ string) (*types.BizTopo, error) {
	cnt := int64(0)
	for _, node := range topo.Nodes {
		if node.Count != nil {
			cnt += *node.Count
		}
	}
	topo.Biz.Count = &cnt
	return topo, nil
}
