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

package y3_9_202011021415

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	tableNameSubscription = "cc_Subscription"
	subscriptionIDField   = "subscription_id"
	subscriptionNameField = "subscription_name"
)

var (
	sortFlag      = 1
	idUniqueIndex = types.Index{
		Keys:       bson.D{{common.BKFieldID, sortFlag}},
		Unique:     true,
		Background: true,
		Name:       "idx_unique_id",
	}
)

func createUniqueIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	tableIndexes := make(map[string][]types.Index, 0)
	buildTopoIndex(tableIndexes)
	buildTopoTemplateIndex(tableIndexes)
	buildModelIndex(tableIndexes)
	buildExtIndex(tableIndexes)
	tips := "If you have created an index for the same field in the table, you can delete " +
		"the existing index in the table and execute migrate again"
	for tableName, indexes := range tableIndexes {
		exists, err := db.HasTable(ctx, tableName)
		if err != nil {
			return err
		}
		if exists {
			for index := range indexes {
				if err = db.Table(tableName).
					CreateIndex(ctx, indexes[index]); err != nil && !db.IsDuplicatedError(err) {

					blog.ErrorJSON("create unique index error. err: %s, table: %s,  index: %s, tips: %s",
						err.Error(), tableName, index, tips)
					return err
				}
			}
		}

	}
	return nil
}

func buildTopoIndex(indexes map[string][]types.Index) {

	indexes[common.BKTableNameBaseApp] = []types.Index{
		{
			Keys:       bson.D{{common.BKAppIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_bizID",
		},
	}

	indexes[common.BKTableNameHostApplyRule] = []types.Index{
		idUniqueIndex,
		{
			Keys: bson.D{
				{common.BKAppIDField, sortFlag},
				{common.BKModuleIDField, sortFlag},
				{common.BKAttributeIDField, sortFlag},
			},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_bizID_moduleID_attrID",
		},
	}

	indexes[common.BKTableNameBaseHost] = []types.Index{
		{
			Keys:       bson.D{{common.BKHostIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_hostID",
		},
		/* 	{
			Keys: bson.D{{common.BKHostInnerIPField,sortFlag}, {common.BKCloudIDField, sortFlag}},
			Unique: true,
			Background: true,
		}, */
	}

	indexes[common.BKTableNameBaseModule] = []types.Index{
		{
			Keys:       bson.D{{common.BKModuleIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_moduleID",
		},
		{
			Keys: bson.D{
				{common.BKAppIDField, sortFlag},
				{common.BKSetIDField, sortFlag},
				{common.BKModuleNameField, sortFlag},
			},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_bizID_setID_moduleName",
		},
	}

	indexes[common.BKTableNameModuleHostConfig] = []types.Index{

		{
			Keys:       bson.D{{common.BKModuleIDField, sortFlag}, {common.BKHostIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_moduleID_hostID",
		},
	}
	indexes[common.BKTableNameBaseSet] = []types.Index{

		{
			Keys:       bson.D{{common.BKSetIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_setID",
		},
		{
			Keys:       bson.D{{common.BKAppIDField, sortFlag}, {common.BKSetNameField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_bizID_setName",
		},
	}
	indexes[common.BKTableNameBaseProcess] = []types.Index{

		{
			Keys:       bson.D{{common.BKProcessIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_procID",
		},
	}

	indexes[common.BKTableNameBasePlat] = []types.Index{

		{
			Keys:       bson.D{{common.BKCloudIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_cloudID",
		},
	}
	indexes[common.BKTableNameProcessInstanceRelation] = []types.Index{

		{
			Keys:       bson.D{{common.BKServiceInstanceIDField, sortFlag}, {common.BKProcessIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_serviceInstID_ProcID",
		},
		{
			Keys:       bson.D{{common.BKProcessIDField, sortFlag}, {common.BKHostIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_procID_hostID",
		},
	}

	indexes[common.BKTableNameServiceInstance] = []types.Index{
		idUniqueIndex,
	}

}

func buildTopoTemplateIndex(indexes map[string][]types.Index) {

	indexes[common.BKTableNameProcessTemplate] = []types.Index{
		idUniqueIndex,
	}
	indexes[common.BKTableNameServiceTemplate] = []types.Index{
		idUniqueIndex,
		{
			Keys:       bson.D{{common.BKAppIDField, sortFlag}, {common.BKFieldName, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_bizID_name",
		},
	}
	indexes[common.BKTableNameSetServiceTemplateRelation] = []types.Index{
		{
			Keys: bson.D{
				{common.BKSetTemplateIDField, sortFlag},
				{common.BKServiceTemplateIDField, sortFlag},
			},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_setTemplateID_serviceTemplateID",
		},
	}

	indexes[common.BKTableNameSetTemplate] = []types.Index{
		idUniqueIndex,
		{
			Keys:       bson.D{{common.BKAppIDField, sortFlag}, {common.BKFieldName, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_bizID_name",
		},
	}

}

func buildModelIndex(indexes map[string][]types.Index) {

	indexes[common.BKTableNameAsstDes] = []types.Index{
		idUniqueIndex,
		{
			Keys:       bson.D{{common.AssociationKindIDField, sortFlag}},
			Background: true,
			Name:       "idx_unique_asstID",
		},
	}

	indexes[common.BKTableNameInstAsst] = []types.Index{
		idUniqueIndex,
	}

	indexes[common.BKTableNameObjAsst] = []types.Index{
		idUniqueIndex,
	}

	indexes[common.BKTableNameObjAttDes] = []types.Index{
		idUniqueIndex,
		{
			Keys: bson.D{
				{common.BKObjIDField, sortFlag},
				{common.BKPropertyIDField, sortFlag},
				{common.BKAppIDField, sortFlag},
			},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_objID_propertyID_bizID",
		},
		{
			Keys: bson.D{
				{common.BKObjIDField, sortFlag},
				{common.BKPropertyNameField, sortFlag},
				{common.BKAppIDField, sortFlag},
			},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_objID_propertyName_bizID",
		},
	}

	indexes[common.BKTableNameObjClassification] = []types.Index{
		idUniqueIndex,
		{
			Keys:       bson.D{{common.BKClassificationIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_classificationID",
		},
		{
			Keys:       bson.D{{common.BKClassificationNameField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_classificationName",
		},
	}

	indexes[common.BKTableNameObjDes] = []types.Index{
		idUniqueIndex,
		{
			Keys:       bson.D{{common.BKObjIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_objID",
		},
	}

	indexes[common.BKTableNameBaseInst] = []types.Index{
		{
			Keys:       bson.D{{common.BKInstIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_instID",
		},
	}

	indexes[common.BKTableNameObjUnique] = []types.Index{
		idUniqueIndex,
	}

	indexes[common.BKTableNamePropertyGroup] = []types.Index{
		idUniqueIndex,
		{
			Keys: bson.D{{common.BKObjIDField, sortFlag}, {common.BKAppIDField, sortFlag},
				{common.BKPropertyGroupNameField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_objID_groupName",
		},
		{
			Keys: bson.D{{common.BKObjIDField, sortFlag}, {common.BKAppIDField, sortFlag},
				{common.BKPropertyGroupIndexField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_objID_groupIdx",
		},
	}
}

func buildExtIndex(indexes map[string][]types.Index) {
	indexes[common.BKTableNameServiceCategory] = []types.Index{
		idUniqueIndex,
		{
			Keys: bson.D{{common.BKFieldName, sortFlag},
				{common.BKParentIDField, sortFlag}, {common.BKAppIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_Name_parentID_bizID",
		},
	}

	indexes[tableNameSubscription] = []types.Index{
		{
			Keys:       bson.D{{subscriptionIDField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_subscriptionID",
		},
		{
			Keys:       bson.D{{subscriptionNameField, sortFlag}},
			Unique:     true,
			Background: true,
			Name:       "idx_unique_subscriptionName",
		},
	}

}
