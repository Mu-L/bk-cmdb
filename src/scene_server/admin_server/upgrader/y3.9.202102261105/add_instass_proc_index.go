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

package y3_9_202102261105

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

// addInstassProcIndex add index for table cc_InstAsst and cc_ProcessTemplate
func addInstassProcIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	index := types.Index{
		Keys: bson.D{
			{"bk_asst_obj_id", 1},
			{"bk_asst_inst_id", 1},
		},
		Name:       "bk_idx_bk_asst_obj_id_bk_asst_inst_id",
		Unique:     false,
		Background: true,
	}

	err := db.Table(common.BKTableNameInstAsst).CreateIndex(ctx, index)
	if err != nil && !db.IsDuplicatedError(err) {
		blog.ErrorJSON("add index %s for table %s failed, err:%s", index, common.BKTableNameInstAsst, err)
		return err
	}

	index = types.Index{
		Keys: bson.D{
			{"service_template_id", 1},
			{"bk_process_name", 1},
		},
		Name:       "bk_idx_service_template_id_bk_process_name",
		Unique:     true,
		Background: true,
	}

	err = db.Table(common.BKTableNameProcessTemplate).CreateIndex(ctx, index)
	if err != nil && !db.IsDuplicatedError(err) {
		blog.ErrorJSON("add index %s for table %s failed, err:%s", index, common.BKTableNameProcessTemplate, err)
		return err
	}

	return nil
}
