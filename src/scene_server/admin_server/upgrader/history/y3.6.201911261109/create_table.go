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

package y3_6_201911261109

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
)

// CreateTableOptions TODO
var CreateTableOptions = []struct {
	TableName  string
	TableIndex []types.Index
}{
	{
		TableName: common.BKTableNameChartConfig,
		TableIndex: []types.Index{
			{Keys: bson.D{{"config_id", 1}}, Name: "config_id", Unique: true, Background: true},
			{Name: common.BKObjIDField, Keys: bson.D{{"bk_obj_id", 1}}, Background: true},
		},
	},
	{
		TableName: common.BKTableNameChartPosition,
		TableIndex: []types.Index{
			{Name: "bk_biz_id", Keys: bson.D{{"bk_biz_id", 1}}, Background: true},
		},
	},
	{
		TableName:  common.BKTableNameChartData,
		TableIndex: []types.Index{},
	},
}

func upsertTable(ctx context.Context, db dal.RDB, conf *upgrader.Config, tableName string, indices []types.Index) error {
	exists, err := db.HasTable(ctx, tableName)
	if err != nil {
		return fmt.Errorf("check HasTable failed, tableName: %s, err: %+v", tableName, err)
	}
	if exists == false {
		if err = db.CreateTable(ctx, tableName); err != nil && !mgo.IsDup(err) {
			return fmt.Errorf("CreateTable failed, tableName: %s, err: %+v", tableName, err)
		}
	}

	existIndices, err := db.Table(tableName).Indexes(ctx)
	if err != nil {
		return fmt.Errorf("upsertTable failed, Indexes failed, tableName: %s, err:%+v", tableName, err)
	}
	existIdxMap := make(map[string]bool)
	for _, idx := range existIndices {
		existIdxMap[idx.Name] = true
	}
	for _, index := range indices {
		if _, ok := existIdxMap[index.Name]; ok == true {
			continue
		}
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return fmt.Errorf("CreateIndex failed, tableName: %s, err:%+v", tableName, err)
		}
	}
	return nil
}

// CreateTables TODO
func CreateTables(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	for _, item := range CreateTableOptions {
		err := upsertTable(ctx, db, conf, item.TableName, item.TableIndex)
		if err != nil {
			return fmt.Errorf("upsertTable failed, err: %s", err.Error())
		}
	}
	return nil
}
