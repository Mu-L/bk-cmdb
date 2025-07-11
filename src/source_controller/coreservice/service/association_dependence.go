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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
)

// IsInstanceExist TODO
func (s *coreService) IsInstanceExist(kit *rest.Kit, objID string, instID uint64) (exists bool, err error) {
	instIDFieldName := common.GetInstIDField(objID)
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: instIDFieldName, Val: instID})
	countCond := &metadata.Condition{Condition: cond.ToMapStr()}
	result, err := s.core.InstanceOperation().CountModelInstances(kit, objID, countCond)
	if err != nil {
		blog.Errorf("search model instance error: %v, rid: %s", err, kit.Rid)
		return false, err
	}
	if result.Count == 0 {
		return false, nil
	}
	return true, nil
}
