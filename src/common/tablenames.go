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

package common

import (
	"fmt"
	"strings"
)

// table names
const (
	// BKTableNamePropertyGroup the table name of the property group
	BKTableNamePropertyGroup = "cc_PropertyGroup"

	// BKTableNameAsstDes the table name of the asst des
	BKTableNameAsstDes = "cc_AsstDes"

	// BKTableNameObjDes the table name of the object
	BKTableNameObjDes = "cc_ObjDes"

	// BKTableNameObjUnique the table name of the object
	BKTableNameObjUnique = "cc_ObjectUnique"

	// BKTableNameObjAttDes the table name of the object attribute
	BKTableNameObjAttDes = "cc_ObjAttDes"

	// BKTableNameObjClassification the table name of the object classification
	BKTableNameObjClassification = "cc_ObjClassification"

	// BKTableNameInstAsst the table name of the inst association
	BKTableNameInstAsst = "cc_InstAsst"

	BKTableNameBaseApp    = "cc_ApplicationBase"
	BKTableNameBaseBizSet = "cc_BizSetBase"

	// BKTableNameModelQuoteRelation model reference relationship table name.
	BKTableNameModelQuoteRelation = "cc_ModelQuoteRelation"

	BKTableNameBaseProject = "cc_ProjectBase"
	BKTableNameBaseHost    = "cc_HostBase"
	BKTableNameBaseModule  = "cc_ModuleBase"
	BKTableNameBaseInst    = "cc_ObjectBase"
	BKTableNameBasePlat    = "cc_PlatBase"
	BKTableNameBaseSet     = "cc_SetBase"
	BKTableNameBaseProcess = "cc_Process"

	BKTableNameDelArchive     = "cc_DelArchive"
	BKTableNameKubeDelArchive = "cc_KubeDelArchive"

	BKTableNameModuleHostConfig = "cc_ModuleHostConfig"
	BKTableNameSystem           = "cc_System"
	BKTableNameHistory          = "cc_History"
	BKTableNameHostFavorite     = "cc_HostFavourite"
	BKTableNameAuditLog         = "cc_AuditLog"
	BKTableNameUserAPI          = "cc_UserAPI"
	BKTableNameDynamicGroup     = "cc_DynamicGroup"
	BKTableNameUserCustom       = "cc_UserCustom"
	BKTableNameObjAsst          = "cc_ObjAsst"
	BKTableNameTopoGraphics     = "cc_TopoGraphics"
	BKTableNameTransaction      = "cc_Transaction"
	BKTableNameIDgenerator      = "cc_idgenerator"

	BKTableNameNetcollectDevice   = "cc_NetcollectDevice"
	BKTableNameNetcollectProperty = "cc_NetcollectProperty"

	BKTableNameNetcollectConfig  = "cc_NetcollectConfig"
	BKTableNameNetcollectReport  = "cc_NetcollectReport"
	BKTableNameNetcollectHistory = "cc_NetcollectHistory"

	BKTableNameHostLock = "cc_HostLock"

	// Operation tables
	BKTableNameChartConfig   = "cc_ChartConfig"
	BKTableNameChartPosition = "cc_ChartPosition"
	BKTableNameChartData     = "cc_ChartData"

	// process tables
	BKTableNameServiceCategory         = "cc_ServiceCategory"
	BKTableNameServiceTemplate         = "cc_ServiceTemplate"
	BKTableNameServiceTemplateAttr     = "cc_ServiceTemplateAttr"
	BKTableNameServiceInstance         = "cc_ServiceInstance"
	BKTableNameProcessTemplate         = "cc_ProcessTemplate"
	BKTableNameProcessInstanceRelation = "cc_ProcessInstanceRelation"

	BKTableNameSetTemplate                = "cc_SetTemplate"
	BKTableNameSetTemplateAttr            = "cc_SetTemplateAttr"
	BKTableNameSetServiceTemplateRelation = "cc_SetServiceTemplateRelation"
	BKTableNameAPITask                    = "cc_APITask"
	BKTableNameAPITaskSyncHistory         = "cc_APITaskSyncHistory"

	// BKTableNameHostApplyRule rule for host property auto apply
	BKTableNameHostApplyRule = "cc_HostApplyRule"

	// cloud sync tables
	BKTableNameCloudSyncTask    = "cc_CloudSyncTask"
	BKTableNameCloudAccount     = "cc_CloudAccount"
	BKTableNameCloudSyncHistory = "cc_CloudSyncHistory"

	// BKTableNameWatchToken the table to store the latest watch token for collections
	BKTableNameWatchToken = "cc_WatchToken"

	// BKTableNameMainlineInstance is a virtual collection name which represent for mainline instance events
	BKTableNameMainlineInstance = "cc_MainlineInstance"

	// BKTableNameFieldTemplate  field template table
	BKTableNameFieldTemplate = "cc_FieldTemplate"

	// BKTableNameObjAttDesTemplate  field template  attribute description table
	BKTableNameObjAttDesTemplate = "cc_ObjAttDesTemplate"

	// BKTableNameObjectUniqueTemplate  field template unique checklist table
	BKTableNameObjectUniqueTemplate = "cc_ObjectUniqueTemplate"

	// BKTableNameObjFieldTemplateRelation  object and field template relationship table
	BKTableNameObjFieldTemplateRelation = "cc_ObjFieldTemplateRelation"
)

// AllTables is all table names, not include the sharding tables which is created dynamically,
// such as object instance sharding table 'cc_ObjectBase_{supplierAccount}_pub_{objectID}'.
var AllTables = []string{
	BKTableNamePropertyGroup,
	BKTableNameObjDes,
	BKTableNameObjAttDes,
	BKTableNameObjClassification,
	BKTableNameInstAsst,
	BKTableNameBaseApp,
	BKTableNameBaseHost,
	BKTableNameBaseModule,
	BKTableNameBaseInst,
	BKTableNameBasePlat,
	BKTableNameBaseSet,
	BKTableNameBaseProcess,
	BKTableNameModuleHostConfig,
	BKTableNameSystem,
	BKTableNameHistory,
	BKTableNameHostFavorite,
	BKTableNameAuditLog,
	BKTableNameUserAPI,
	BKTableNameDynamicGroup,
	BKTableNameUserCustom,
	BKTableNameObjAsst,
	BKTableNameTopoGraphics,
	BKTableNameNetcollectConfig,
	BKTableNameNetcollectDevice,
	BKTableNameNetcollectProperty,
	BKTableNameNetcollectReport,
	BKTableNameNetcollectHistory,
	BKTableNameTransaction,
	BKTableNameIDgenerator,
	BKTableNameHostLock,
	BKTableNameObjUnique,
	BKTableNameAsstDes,
	BKTableNameServiceCategory,
	BKTableNameServiceTemplate,
	BKTableNameServiceInstance,
	BKTableNameProcessTemplate,
	BKTableNameProcessInstanceRelation,
	BKTableNameSetTemplate,
	BKTableNameSetServiceTemplateRelation,
	BKTableNameChartConfig,
	BKTableNameChartPosition,
	BKTableNameChartData,
	BKTableNameHostApplyRule,
	BKTableNameAPITask,
	BKTableNameAPITaskSyncHistory,
	BKTableNameCloudSyncTask,
	BKTableNameCloudAccount,
	BKTableNameCloudSyncHistory,
}

// TableSpecifier is table specifier type which describes the metadata
// access or classification level.
type TableSpecifier string

const (
	// TableSpecifierPublic is public specifier for table.
	TableSpecifierPublic TableSpecifier = "pub"
)

const (
	// BKObjectInstShardingTablePrefix is prefix of object instance sharding table.
	BKObjectInstShardingTablePrefix = BKTableNameBaseInst + "_"

	// BKObjectInstAsstShardingTablePrefix is prefix of object instance association sharding table.
	BKObjectInstAsstShardingTablePrefix = BKTableNameInstAsst + "_"
)

// GetObjectInstTableName return the object instance table name in sharding mode base on
// the object ID. Format: cc_ObjectBase_{supplierAccount}_{Specifier}_{ObjectID}, such as 'cc_ObjectBase_0_pub_switch'.
func GetObjectInstTableName(objID, supplierAccount string) string {
	return fmt.Sprintf("%s%s_%s_%s", BKObjectInstShardingTablePrefix, supplierAccount, TableSpecifierPublic, objID)
}

// GetObjectInstObjIDByTableName return the object id
// example: cc_ObjectBase_{supplierAccount}_{Specifier}_{ObjectID}, such as 'cc_ObjectBase_0_pub_switch',return switch.
func GetObjectInstObjIDByTableName(collectionName, supplierAccount string) (string, error) {
	prefix := fmt.Sprintf("%s%s_", BKObjectInstShardingTablePrefix, supplierAccount)
	suffix := strings.TrimPrefix(collectionName, prefix)
	suffixSlice := strings.Split(suffix, "_")
	if len(suffixSlice) <= 1 {
		return "", fmt.Errorf("collection name is error, collection name: %s", collectionName)
	}
	return strings.Join(suffixSlice[1:], "_"), nil
}

// GetObjectInstAsstTableName return the object instance association table name in sharding mode base on
// the object ID. Format: cc_InstAsst_{supplierAccount}_{Specifier}_{ObjectID}, such as 'cc_InstAsst_0_pub_switch'.
func GetObjectInstAsstTableName(objID, supplierAccount string) string {
	return fmt.Sprintf("%s%s_%s_%s", BKObjectInstAsstShardingTablePrefix, supplierAccount, TableSpecifierPublic, objID)
}

// IsObjectShardingTable returns if the target table is an object sharding table, include
// object instance and association.
func IsObjectShardingTable(tableName string) bool {
	if IsObjectInstShardingTable(tableName) {
		return true
	}
	return IsObjectInstAsstShardingTable(tableName)
}

// IsObjectInstShardingTable returns if the target table is an object instance sharding table.
func IsObjectInstShardingTable(tableName string) bool {
	// check object instance table, cc_ObjectBase_{Specifier}_{ObjectID}
	return strings.HasPrefix(tableName, BKObjectInstShardingTablePrefix)
}

// IsObjectInstAsstShardingTable returns if the target table is an object instance association sharding table.
func IsObjectInstAsstShardingTable(tableName string) bool {
	// check object instance association table, cc_InstAsst_{Specifier}_{ObjectID}
	return strings.HasPrefix(tableName, BKObjectInstAsstShardingTablePrefix)
}

// GetInstTableName returns inst data table name
func GetInstTableName(objID, supplierAccount string) string {
	switch objID {
	case BKInnerObjIDApp:
		return BKTableNameBaseApp
	case BKInnerObjIDBizSet:
		return BKTableNameBaseBizSet
	case BKInnerObjIDProject:
		return BKTableNameBaseProject
	case BKInnerObjIDSet:
		return BKTableNameBaseSet
	case BKInnerObjIDModule:
		return BKTableNameBaseModule
	case BKInnerObjIDHost:
		return BKTableNameBaseHost
	case BKInnerObjIDProc:
		return BKTableNameBaseProcess
	case BKInnerObjIDPlat:
		return BKTableNameBasePlat
	default:
		return GetObjectInstTableName(objID, supplierAccount)
	}
}

// GetInstObjIDByTableName get objID by table name
func GetInstObjIDByTableName(collectionName, supplierAccount string) (string, error) {
	switch collectionName {
	case BKTableNameBaseApp:
		return BKInnerObjIDApp, nil
	case BKTableNameBaseBizSet:
		return BKInnerObjIDBizSet, nil
	case BKTableNameBaseProject:
		return BKInnerObjIDProject, nil
	case BKTableNameBaseSet:
		return BKInnerObjIDSet, nil
	case BKTableNameBaseModule:
		return BKInnerObjIDModule, nil
	case BKTableNameBaseHost:
		return BKInnerObjIDHost, nil
	case BKTableNameBaseProcess:
		return BKInnerObjIDProc, nil
	case BKTableNameBasePlat:
		return BKInnerObjIDPlat, nil
	default:
		return GetObjectInstObjIDByTableName(collectionName, supplierAccount)
	}
}
