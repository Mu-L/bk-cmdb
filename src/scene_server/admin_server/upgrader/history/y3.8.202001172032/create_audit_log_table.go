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

package y3_8_202001172032

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func createAuditLogTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	exists, err := db.HasTable(ctx, common.BKTableNameAuditLog)
	if err != nil {
		blog.ErrorJSON("search audit log table error, err:%s", err.Error())
		return err
	}
	if !exists {
		if err = db.CreateTable(ctx, common.BKTableNameAuditLog); err != nil && !db.IsDuplicatedError(err) {
			blog.ErrorJSON("create audit log table error, err:%s", err.Error())
			return err
		}
	}
	return nil
}

func addAuditLogTableIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	idxArr, err := db.Table(common.BKTableNameAuditLog).Indexes(ctx)
	if err != nil {
		blog.Errorf("get table %s index error. err:%s", common.BKTableNameAuditLog, err.Error())
		return err
	}

	createIdxArr := []types.Index{
		{Name: "index_bk_supplier_account", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
		{Name: "index_audit_type", Keys: bson.D{{common.BKAuditTypeField, 1}}, Background: true},
		{Name: "index_action", Keys: bson.D{{common.BKActionField, 1}}, Background: true},
	}
	for _, idx := range createIdxArr {
		exist := false
		for _, existIdx := range idxArr {
			if existIdx.Name == idx.Name {
				exist = true
				break
			}
		}
		if exist {
			continue
		}
		if err := db.Table(common.BKTableNameAuditLog).CreateIndex(ctx, idx); err != nil && !db.IsDuplicatedError(err) {
			blog.ErrorJSON("create index to BKTableNameAuditLog error, err:%s, current index:%s, all create index:%s", err.Error(), idx, createIdxArr)
			return err
		}

	}

	return nil
}
