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

package x19_10_22_02

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func isIndexExist(ctx context.Context, db dal.RDB, indexName string) (bool, error) {
	indexes, err := db.Table(common.BKTableNameBaseHost).Indexes(ctx)
	if err != nil {
		return false, fmt.Errorf("list indexes failed, table: %s, err: %s", common.BKTableNameBaseHost, err.Error())
	}
	for _, index := range indexes {
		if index.Name == indexName {
			return true, nil
		}
	}
	return false, nil
}

func addHostIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	indexName := "innerIP_platID"
	exist, err := isIndexExist(ctx, db, indexName)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	idx := types.Index{
		Keys: bson.D{
			{common.BKHostInnerIPField, 1},
			{common.BKCloudIDField, 1},
		},
		Name:       indexName,
		Unique:     false,
		Background: false,
	}

	if err := db.Table(common.BKTableNameBaseHost).CreateIndex(ctx, idx); err != nil && !db.IsDuplicatedError(err) {
		blog.Errorf("CreateIndex failed, err: %+v", err)
		return fmt.Errorf("CreateIndex failed, err: %v", err)
	}
	return nil
}
