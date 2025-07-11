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

package y3_9_202008172134

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

// reconcileAuditTableIndexes update indexes for common audit log query params
func reconcileAuditTableIndexes(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	indexes := []types.Index{
		{Name: "index_id", Keys: bson.D{{common.BKFieldID, 1}}, Background: true},
		{Name: "index_operationTime", Keys: bson.D{{common.BKOperationTimeField, 1}}, Background: true},
		{Name: "index_user", Keys: bson.D{{common.BKUser, 1}}, Background: true},
		{Name: "index_resourceName", Keys: bson.D{{common.BKResourceNameField, 1}}, Background: true},
		{Name: "index_operationTime_auditType_resourceType_action", Keys: bson.D{
			{common.BKOperationTimeField, 1},
			{common.BKAuditTypeField, 1},
			{common.BKResourceTypeField, 1},
			{common.BKOperationDetailField + "." + common.BKObjIDField, 1},
			{common.BKActionField, 1},
		}, Background: true},
	}

	existIndexArr, err := db.Table(common.BKTableNameAuditLog).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for audit table failed, err: %s", err.Error())
		return err
	}

	existIdxMap := make(map[string]bool)
	for _, index := range existIndexArr {
		existIdxMap[index.Name] = true
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}

		if err = db.Table(common.BKTableNameAuditLog).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for audit table failed, err: %s, index: %+v", err.Error(), index)
			return err
		}
	}

	removeIndexes := []string{"index_bk_supplier_account", "index_audit_type", "index_action"}
	for _, removeIndex := range removeIndexes {
		if _, exist := existIdxMap[removeIndex]; !exist {
			continue
		}
		if err = db.Table(common.BKTableNameAuditLog).DropIndex(ctx, removeIndex); err != nil {
			blog.Errorf("remove index for audit table failed, err: %s, index: %s", err.Error(), removeIndex)
			return err
		}

	}
	return nil
}
