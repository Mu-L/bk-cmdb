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

package x19_08_26_02

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// Idgen TODO
type Idgen struct {
	ID         string `bson:"_id"`
	SequenceID uint64 `bson:"SequenceID"`
}

// FixServiceInstanceMaxID 将 cc_idgenerator 中 cc_ProcessTemplate 和 cc_ServiceInstance 两个key的值升级到他们的最大值
func FixServiceInstanceMaxID(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		"_id": common.BKTableNameServiceInstance,
	}
	var maxID uint64
	idGenerator := Idgen{}
	if err := db.Table(common.BKTableNameIDgenerator).Find(filter).One(ctx, &idGenerator); err != nil {
		if db.IsNotFoundError(err) == false {
			return fmt.Errorf("upgrade x19_08_26_02, get service instance id generator failed, err: %v", err)
		}
	} else {
		if idGenerator.SequenceID > maxID {
			maxID = idGenerator.SequenceID
		}
	}

	processTemplateFilter := map[string]interface{}{
		"_id": common.BKTableNameProcessTemplate,
	}
	if err := db.Table(common.BKTableNameIDgenerator).Find(processTemplateFilter).One(ctx, &idGenerator); err != nil {
		if db.IsNotFoundError(err) == false {
			return fmt.Errorf("upgrade x19_08_26_02, get process template id generator failed, err: %v", err)
		}
	} else {
		if idGenerator.SequenceID > maxID {
			maxID = idGenerator.SequenceID
		}
	}

	updateFiler := map[string]interface{}{
		"_id": map[string]interface{}{
			common.BKDBIN: []string{common.BKTableNameProcessTemplate, common.BKTableNameServiceInstance},
		},
	}
	doc := map[string]interface{}{
		"SequenceID": maxID,
	}
	if err := db.Table(common.BKTableNameIDgenerator).Update(ctx, updateFiler, doc); err != nil {
		if db.IsNotFoundError(err) == false {
			return fmt.Errorf("upgrade x19_08_26_02, update max id failed, err: %v", err)
		}
	}
	return nil
}
