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

package iam

import "configcenter/src/common/metadata"

var (
	businessResource = RelateResourceType{
		SystemID:    SystemIDCMDB,
		ID:          Business,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       BusinessSelection,
		}},
	}

	resourcePoolDirResource = RelateResourceType{
		SystemID:    SystemIDCMDB,
		ID:          SysResourcePoolDirectory,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       SysResourcePoolDirectorySelection,
		}},
	}
)

// ActionIDNameMap TODO
var ActionIDNameMap = map[ActionID]string{
	EditBusinessHost:                    "业务主机编辑",
	BusinessHostTransferToResourcePool:  "主机归还主机池",
	HostTransferAcrossBusiness:          "主机转移到其他业务",
	CreateBusinessCustomQuery:           "动态分组新建",
	EditBusinessCustomQuery:             "动态分组编辑",
	DeleteBusinessCustomQuery:           "动态分组删除",
	EditBusinessCustomField:             "业务自定义字段编辑",
	CreateBusinessServiceCategory:       "服务分类新建",
	EditBusinessServiceCategory:         "服务分类编辑",
	DeleteBusinessServiceCategory:       "服务分类删除",
	CreateBusinessServiceInstance:       "服务实例新建",
	EditBusinessServiceInstance:         "服务实例编辑",
	DeleteBusinessServiceInstance:       "服务实例删除",
	CreateBusinessServiceTemplate:       "服务模板新建",
	EditBusinessServiceTemplate:         "服务模板编辑",
	DeleteBusinessServiceTemplate:       "服务模板删除",
	CreateBusinessSetTemplate:           "集群模板新建",
	EditBusinessSetTemplate:             "集群模板编辑",
	DeleteBusinessSetTemplate:           "集群模板删除",
	CreateBusinessTopology:              "业务拓扑新建",
	EditBusinessTopology:                "业务拓扑编辑",
	DeleteBusinessTopology:              "业务拓扑删除",
	EditBusinessHostApply:               "主机自动应用编辑",
	ViewResourcePoolHost:                "主机池主机查看",
	CreateResourcePoolHost:              "主机池主机创建",
	EditResourcePoolHost:                "主机池主机编辑",
	DeleteResourcePoolHost:              "主机池主机删除",
	ResourcePoolHostTransferToBusiness:  "主机池主机分配到业务",
	ResourcePoolHostTransferToDirectory: "主机池主机分配到目录",
	CreateResourcePoolDirectory:         "主机池目录创建",
	EditResourcePoolDirectory:           "主机池目录编辑",
	DeleteResourcePoolDirectory:         "主机池目录删除",
	CreateBusiness:                      "业务创建",
	EditBusiness:                        "业务编辑",
	ArchiveBusiness:                     "业务归档",
	FindBusiness:                        "业务查询",
	ViewBusinessResource:                "业务访问",
	CreateBizSet:                        "业务集新增",
	EditBizSet:                          "业务集编辑",
	DeleteBizSet:                        "业务集删除",
	ViewBizSet:                          "业务集查看",
	AccessBizSet:                        "业务集访问",
	CreateProject:                       "项目新建",
	EditProject:                         "项目编辑",
	DeleteProject:                       "项目删除",
	ViewProject:                         "项目查看",
	ViewCloudArea:                       "管控区域查看",
	CreateCloudArea:                     "管控区域创建",
	EditCloudArea:                       "管控区域编辑",
	DeleteCloudArea:                     "管控区域删除",
	CreateCloudAccount:                  "云账户新建",
	EditCloudAccount:                    "云账户编辑",
	DeleteCloudAccount:                  "云账户删除",
	FindCloudAccount:                    "云账户查询",
	CreateCloudResourceTask:             "云资源任务新建",
	EditCloudResourceTask:               "云资源任务编辑",
	DeleteCloudResourceTask:             "云资源任务删除",
	FindCloudResourceTask:               "云资源任务查询",
	ViewSysModel:                        "模型查看",
	CreateSysModel:                      "模型新建",
	EditSysModel:                        "模型编辑",
	DeleteSysModel:                      "模型删除",
	CreateAssociationType:               "关联类型新建",
	EditAssociationType:                 "关联类型编辑",
	DeleteAssociationType:               "关联类型删除",
	CreateModelGroup:                    "模型分组新建",
	EditModelGroup:                      "模型分组编辑",
	DeleteModelGroup:                    "模型分组删除",
	ViewModelTopo:                       "模型拓扑查看",
	EditBusinessLayer:                   "业务层级编辑",
	EditModelTopologyView:               "模型拓扑视图编辑",
	FindOperationStatistic:              "运营统计查询",
	EditOperationStatistic:              "运营统计编辑",
	FindAuditLog:                        "操作审计查询",
	WatchHostEvent:                      "主机事件监听",
	WatchHostRelationEvent:              "主机关系事件监听",
	WatchBizEvent:                       "业务事件监听",
	WatchSetEvent:                       "集群事件监听",
	WatchModuleEvent:                    "模块数据监听",
	WatchProcessEvent:                   "进程数据监听",
	WatchCommonInstanceEvent:            "模型实例事件监听",
	WatchMainlineInstanceEvent:          "自定义拓扑层级事件监听",
	WatchInstAsstEvent:                  "实例关联事件监听",
	WatchBizSetEvent:                    "业务集事件监听",
	WatchPlatEvent:                      "管控区域事件监听",
	WatchKubeClusterEvent:               "容器集群事件监听",
	WatchKubeNodeEvent:                  "容器节点事件监听",
	WatchKubeNamespaceEvent:             "容器命名空间事件监听",
	WatchKubeWorkloadEvent:              "容器工作负载事件监听",
	WatchKubePodEvent:                   "容器Pod事件监听",
	WatchProjectEvent:                   "项目事件监听",
	GlobalSettings:                      "全局设置",
	ManageHostAgentID:                   "主机AgentID管理",
	CreateContainerCluster:              "容器集群新建",
	EditContainerCluster:                "容器集群编辑",
	DeleteContainerCluster:              "容器集群删除",
	CreateContainerNode:                 "容器集群节点新建",
	EditContainerNode:                   "容器集群节点编辑",
	DeleteContainerNode:                 "容器集群节点删除",
	CreateContainerNamespace:            "容器命名空间新建",
	EditContainerNamespace:              "容器命名空间编辑",
	DeleteContainerNamespace:            "容器命名空间删除",
	CreateContainerWorkload:             "容器工作负载新建",
	EditContainerWorkload:               "容器工作负载编辑",
	DeleteContainerWorkload:             "容器工作负载删除",
	CreateContainerPod:                  "容器Pod新建",
	DeleteContainerPod:                  "容器Pod删除",
	UseFulltextSearch:                   "全文检索",
	CreateFieldGroupingTemplate:         "字段组合模板新建",
	ViewFieldGroupingTemplate:           "字段组合模板查看",
	EditFieldGroupingTemplate:           "字段组合模板编辑",
	DeleteFieldGroupingTemplate:         "字段组合模板删除",
	EditIDRuleIncrID:                    "ID规则自增ID编辑",
	CreateFullSyncCond:                  "全量同步缓存条件新建",
	ViewFullSyncCond:                    "全量同步缓存条件查看",
	EditFullSyncCond:                    "全量同步缓存条件编辑",
	DeleteFullSyncCond:                  "全量同步缓存条件删除",
	ViewGeneralCache:                    "通用缓存查询",
}

// GenerateActions generate all the actions registered to IAM.
func GenerateActions(objects []metadata.Object) []ResourceAction {
	resourceActionList := GenerateStaticActions()
	resourceActionList = append(resourceActionList, genDynamicActions(objects)...)
	return resourceActionList
}

// GenerateStaticActions TODO
func GenerateStaticActions() []ResourceAction {
	resourceActionList := make([]ResourceAction, 0)
	// add business resource actions
	resourceActionList = append(resourceActionList, genBusinessHostActions()...)
	resourceActionList = append(resourceActionList, genBusinessCustomQueryActions()...)
	resourceActionList = append(resourceActionList, genBusinessCustomFieldActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceCategoryActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceInstanceActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceTemplateActions()...)
	resourceActionList = append(resourceActionList, genBusinessSetTemplateActions()...)
	resourceActionList = append(resourceActionList, genBusinessTopologyActions()...)
	resourceActionList = append(resourceActionList, genBusinessHostApplyActions()...)

	// add public resource actions
	resourceActionList = append(resourceActionList, genResourcePoolHostActions()...)
	resourceActionList = append(resourceActionList, genResourcePoolDirectoryActions()...)
	resourceActionList = append(resourceActionList, genBusinessActions()...)
	resourceActionList = append(resourceActionList, genBizSetActions()...)
	resourceActionList = append(resourceActionList, genProjectActions()...)
	resourceActionList = append(resourceActionList, genCloudAreaActions()...)
	resourceActionList = append(resourceActionList, genCloudAccountActions()...)
	resourceActionList = append(resourceActionList, genCloudResourceTaskActions()...)
	resourceActionList = append(resourceActionList, genModelActions()...)
	resourceActionList = append(resourceActionList, genAssociationTypeActions()...)
	resourceActionList = append(resourceActionList, genModelGroupActions()...)
	resourceActionList = append(resourceActionList, genBusinessLayerActions()...)
	resourceActionList = append(resourceActionList, genModelTopologyViewActions()...)
	resourceActionList = append(resourceActionList, genOperationStatisticActions()...)
	resourceActionList = append(resourceActionList, genAuditLogActions()...)
	resourceActionList = append(resourceActionList, genEventWatchActions()...)
	resourceActionList = append(resourceActionList, genKubeEventWatchActions()...)
	resourceActionList = append(resourceActionList, genConfigAdminActions()...)
	resourceActionList = append(resourceActionList, genContainerManagementActions()...)
	resourceActionList = append(resourceActionList, genFulltextSearchActions()...)
	resourceActionList = append(resourceActionList, genFieldGroupingTemplateActions()...)
	resourceActionList = append(resourceActionList, genIDRuleActions()...)
	resourceActionList = append(resourceActionList, genFullSyncCondActions()...)
	resourceActionList = append(resourceActionList, genCacheActions()...)

	return resourceActionList
}

func genBusinessHostActions() []ResourceAction {
	hostSelection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       BizHostInstanceSelection,
	}}

	relatedResource := []RelateResourceType{{
		SystemID:    SystemIDCMDB,
		ID:          Host,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		// 配置权限时可选择实例和配置属性, 后者用于属性鉴权
		SelectionMode:      modeAll,
		InstanceSelections: hostSelection,
	}}

	actions := make([]ResourceAction, 0)

	// edit business's host actions
	actions = append(actions, ResourceAction{
		ID:                   EditBusinessHost,
		Name:                 ActionIDNameMap[EditBusinessHost],
		NameEn:               "Edit Business Hosts",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	// business host transfer to resource pool actions
	actions = append(actions, ResourceAction{
		ID:                   BusinessHostTransferToResourcePool,
		Name:                 ActionIDNameMap[BusinessHostTransferToResourcePool],
		NameEn:               "Return Hosts To Pool",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource, resourcePoolDirResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	businessHostResource := RelateResourceType{
		SystemID:    SystemIDCMDB,
		ID:          BusinessForHostTrans,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       BusinessHostTransferSelection,
		}},
	}

	// business host transfer to another business actions
	actions = append(actions, ResourceAction{
		ID:                   HostTransferAcrossBusiness,
		Name:                 ActionIDNameMap[HostTransferAcrossBusiness],
		NameEn:               "Assigned Host To Other Business",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessHostResource, businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessCustomQueryActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       BizCustomQuerySelection,
	}}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 BizCustomQuery,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessCustomQuery,
		Name:                 ActionIDNameMap[CreateBusinessCustomQuery],
		NameEn:               "Create Dynamic Grouping",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessCustomQuery,
		Name:                 ActionIDNameMap[EditBusinessCustomQuery],
		NameEn:               "Edit Dynamic Grouping",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessCustomQuery,
		Name:                 ActionIDNameMap[DeleteBusinessCustomQuery],
		NameEn:               "Delete Dynamic Grouping",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessCustomFieldActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessCustomField,
		Name:                 ActionIDNameMap[EditBusinessCustomField],
		NameEn:               "Edit Custom Field",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessServiceCategoryActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessServiceCategory,
		Name:                 ActionIDNameMap[CreateBusinessServiceCategory],
		NameEn:               "Create Service Category",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceCategory,
		Name:                 ActionIDNameMap[EditBusinessServiceCategory],
		NameEn:               "Edit Service Category",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceCategory,
		Name:                 ActionIDNameMap[DeleteBusinessServiceCategory],
		NameEn:               "Delete Service Category",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessServiceInstanceActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessServiceInstance,
		Name:                 ActionIDNameMap[CreateBusinessServiceInstance],
		NameEn:               "Create Service Instance",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceInstance,
		Name:                 ActionIDNameMap[EditBusinessServiceInstance],
		NameEn:               "Edit Service Instance",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceInstance,
		Name:                 ActionIDNameMap[DeleteBusinessServiceInstance],
		NameEn:               "Delete Service Instance",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessServiceTemplateActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID:       SystemIDCMDB,
		ID:             BizProcessServiceTemplateSelection,
		IgnoreAuthPath: true,
	}}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 BizProcessServiceTemplate,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessServiceTemplate,
		Name:                 ActionIDNameMap[CreateBusinessServiceTemplate],
		NameEn:               "Create Service Template",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceTemplate,
		Name:                 ActionIDNameMap[EditBusinessServiceTemplate],
		NameEn:               "Edit Service Template",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceTemplate,
		Name:                 ActionIDNameMap[DeleteBusinessServiceTemplate],
		NameEn:               "Delete Service Template",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessSetTemplateActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID:       SystemIDCMDB,
		ID:             BizSetTemplateSelection,
		IgnoreAuthPath: true,
	}}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 BizSetTemplate,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessSetTemplate,
		Name:                 ActionIDNameMap[CreateBusinessSetTemplate],
		NameEn:               "Create Set Template",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessSetTemplate,
		Name:                 ActionIDNameMap[EditBusinessSetTemplate],
		NameEn:               "Edit Set Template",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessSetTemplate,
		Name:                 ActionIDNameMap[DeleteBusinessSetTemplate],
		NameEn:               "Delete Set Template",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessTopologyActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessTopology,
		Name:                 ActionIDNameMap[CreateBusinessTopology],
		NameEn:               "Create Business Topo",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessTopology,
		Name:                 ActionIDNameMap[EditBusinessTopology],
		NameEn:               "Edit Business Topo",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessTopology,
		Name:                 ActionIDNameMap[DeleteBusinessTopology],
		NameEn:               "Delete Business Topo",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessHostApplyActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessHostApply,
		Name:                 ActionIDNameMap[EditBusinessHostApply],
		NameEn:               "Edit Host Apply",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genResourcePoolHostActions() []ResourceAction {
	hostSelection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysHostInstanceSelection,
	}}

	relatedResource := []RelateResourceType{{
		SystemID:    SystemIDCMDB,
		ID:          Host,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		// 配置权限时可选择实例和配置属性, 后者用于属性鉴权
		SelectionMode:      modeAll,
		InstanceSelections: hostSelection,
	}}

	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   ViewResourcePoolHost,
		Name:                 ActionIDNameMap[ViewResourcePoolHost],
		NameEn:               "View Resource Pool Hosts",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   CreateResourcePoolHost,
		Name:                 ActionIDNameMap[CreateResourcePoolHost],
		NameEn:               "Create Pool Hosts",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditResourcePoolHost,
		Name:                 ActionIDNameMap[EditResourcePoolHost],
		NameEn:               "Edit Pool Hosts",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewResourcePoolHost},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteResourcePoolHost,
		Name:                 ActionIDNameMap[DeleteResourcePoolHost],
		NameEn:               "Delete Pool Hosts",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewResourcePoolHost},
		Version:              1,
	})

	relatedHostResource := []RelateResourceType{{
		SystemID:    SystemIDCMDB,
		ID:          SysHostRscPoolDirectory,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       SysHostRscPoolDirectorySelection,
		}},
	}}

	transferToBusinessRelatedResource := append(relatedHostResource, businessResource)
	actions = append(actions, ResourceAction{
		ID:                   ResourcePoolHostTransferToBusiness,
		Name:                 ActionIDNameMap[ResourcePoolHostTransferToBusiness],
		NameEn:               "Assigned Pool Hosts To Business",
		Type:                 Edit,
		RelatedResourceTypes: transferToBusinessRelatedResource,
		RelatedActions:       []ActionID{ViewResourcePoolHost},
		Version:              1,
	})

	transferToDirectoryRelatedResource := append(relatedHostResource, resourcePoolDirResource)
	actions = append(actions, ResourceAction{
		ID:                   ResourcePoolHostTransferToDirectory,
		Name:                 ActionIDNameMap[ResourcePoolHostTransferToDirectory],
		NameEn:               "Assigned Pool Hosts To Directory",
		Type:                 Edit,
		RelatedResourceTypes: transferToDirectoryRelatedResource,
		RelatedActions:       []ActionID{ViewResourcePoolHost},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:      ManageHostAgentID,
		Name:    ActionIDNameMap[ManageHostAgentID],
		NameEn:  "Manage Host AgentID",
		Type:    Edit,
		Version: 1,
	})

	return actions
}

func genResourcePoolDirectoryActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateResourcePoolDirectory,
		Name:                 ActionIDNameMap[CreateResourcePoolDirectory],
		NameEn:               "Create Pool Directory",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditResourcePoolDirectory,
		Name:                 ActionIDNameMap[EditResourcePoolDirectory],
		NameEn:               "Edit Pool Directory",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteResourcePoolDirectory,
		Name:                 ActionIDNameMap[DeleteResourcePoolDirectory],
		NameEn:               "Delete Pool Directory",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusiness,
		Name:                 ActionIDNameMap[CreateBusiness],
		NameEn:               "Create Business",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusiness,
		Name:                 ActionIDNameMap[EditBusiness],
		NameEn:               "Edit Business",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{FindBusiness},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   ArchiveBusiness,
		Name:                 ActionIDNameMap[ArchiveBusiness],
		NameEn:               "Archive Business",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{FindBusiness},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindBusiness,
		Name:                 ActionIDNameMap[FindBusiness],
		NameEn:               "View Business",
		Type:                 View,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:     ViewBusinessResource,
		Name:   ActionIDNameMap[ViewBusinessResource],
		NameEn: "View Business Resource",
		Type:   View,
		// TODO add business collection resource
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBizSetActions() []ResourceAction {
	bizSetResource := RelateResourceType{
		SystemID: SystemIDCMDB,
		ID:       BizSet,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       BizSetSelection,
		}},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBizSet,
		Name:                 ActionIDNameMap[CreateBizSet],
		NameEn:               "Create Business Set",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       []ActionID{ViewBizSet},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBizSet,
		Name:                 ActionIDNameMap[EditBizSet],
		NameEn:               "Edit Business Set",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{bizSetResource},
		RelatedActions:       []ActionID{ViewBizSet},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBizSet,
		Name:                 ActionIDNameMap[DeleteBizSet],
		NameEn:               "Delete Business Set",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{bizSetResource},
		RelatedActions:       []ActionID{ViewBizSet},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   ViewBizSet,
		Name:                 ActionIDNameMap[ViewBizSet],
		NameEn:               "View Business Set",
		Type:                 View,
		RelatedResourceTypes: []RelateResourceType{bizSetResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   AccessBizSet,
		Name:                 ActionIDNameMap[AccessBizSet],
		NameEn:               "Access Business Set",
		Type:                 View,
		RelatedResourceTypes: []RelateResourceType{bizSetResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genProjectActions() []ResourceAction {
	projectResource := RelateResourceType{
		SystemID: SystemIDCMDB,
		ID:       Project,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       ProjectSelection,
		}},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateProject,
		Name:                 ActionIDNameMap[CreateProject],
		NameEn:               "Create Project",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditProject,
		Name:                 ActionIDNameMap[EditProject],
		NameEn:               "Edit Project",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{projectResource},
		RelatedActions:       []ActionID{ViewProject},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteProject,
		Name:                 ActionIDNameMap[DeleteProject],
		NameEn:               "Delete Project",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{projectResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   ViewProject,
		Name:                 ActionIDNameMap[ViewProject],
		NameEn:               "View Project",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genCloudAreaActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysCloudAreaSelection,
	}}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysCloudArea,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   ViewCloudArea,
		Name:                 ActionIDNameMap[ViewCloudArea],
		NameEn:               "View Cloud Area",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   CreateCloudArea,
		Name:                 ActionIDNameMap[CreateCloudArea],
		NameEn:               "Create Cloud Area",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditCloudArea,
		Name:                 ActionIDNameMap[EditCloudArea],
		NameEn:               "Edit Cloud Area",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewCloudArea},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteCloudArea,
		Name:                 ActionIDNameMap[DeleteCloudArea],
		NameEn:               "Delete Cloud Area",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewCloudArea},
		Version:              1,
	})

	return actions
}

func genCloudAccountActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysCloudAccountSelection,
	}}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysCloudAccount,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateCloudAccount,
		Name:                 ActionIDNameMap[CreateCloudAccount],
		NameEn:               "Create Cloud Account",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditCloudAccount,
		Name:                 ActionIDNameMap[EditCloudAccount],
		NameEn:               "Edit Cloud Account",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{FindCloudAccount},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteCloudAccount,
		Name:                 ActionIDNameMap[DeleteCloudAccount],
		NameEn:               "Delete Cloud Account",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{FindCloudAccount},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindCloudAccount,
		Name:                 ActionIDNameMap[FindCloudAccount],
		NameEn:               "View Cloud Account",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genCloudResourceTaskActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysCloudResourceTaskSelection,
	}}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysCloudResourceTask,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateCloudResourceTask,
		Name:                 ActionIDNameMap[CreateCloudResourceTask],
		NameEn:               "Create Cloud Resource Task",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditCloudResourceTask,
		Name:                 ActionIDNameMap[EditCloudResourceTask],
		NameEn:               "Edit Cloud Resource Task",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{FindCloudResourceTask},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteCloudResourceTask,
		Name:                 ActionIDNameMap[DeleteCloudResourceTask],
		NameEn:               "Delete Cloud Resource Task",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{FindCloudResourceTask},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindCloudResourceTask,
		Name:                 ActionIDNameMap[FindCloudResourceTask],
		NameEn:               "View Cloud Resource Task",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genModelActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysModelSelection,
	}}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysModel,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   ViewSysModel,
		Name:                 ActionIDNameMap[ViewSysModel],
		NameEn:               "View Model",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:     CreateSysModel,
		Name:   ActionIDNameMap[CreateSysModel],
		NameEn: "Create Model",
		Type:   Create,
		RelatedResourceTypes: []RelateResourceType{
			{
				SystemID:    SystemIDCMDB,
				ID:          SysModelGroup,
				NameAlias:   "",
				NameAliasEn: "",
				Scope:       nil,
				InstanceSelections: []RelatedInstanceSelection{{
					SystemID: SystemIDCMDB,
					ID:       SysModelGroupSelection,
				}},
			},
		},
		RelatedActions: nil,
		Version:        1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditSysModel,
		Name:                 ActionIDNameMap[EditSysModel],
		NameEn:               "Edit Model",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewSysModel},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteSysModel,
		Name:                 ActionIDNameMap[DeleteSysModel],
		NameEn:               "Delete Model",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewSysModel},
		Version:              1,
	})

	return actions
}

func genAssociationTypeActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysAssociationTypeSelection,
	}}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysAssociationType,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateAssociationType,
		Name:                 ActionIDNameMap[CreateAssociationType],
		NameEn:               "Create Association Type",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditAssociationType,
		Name:                 ActionIDNameMap[EditAssociationType],
		NameEn:               "Edit Association Type",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteAssociationType,
		Name:                 ActionIDNameMap[DeleteAssociationType],
		NameEn:               "Delete Association Type",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genModelGroupActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysModelGroupSelection,
	}}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysModelGroup,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateModelGroup,
		Name:                 ActionIDNameMap[CreateModelGroup],
		NameEn:               "Create Model Group",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditModelGroup,
		Name:                 ActionIDNameMap[EditModelGroup],
		NameEn:               "Edit Model Group",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteModelGroup,
		Name:                 ActionIDNameMap[DeleteModelGroup],
		NameEn:               "Delete Model Group",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessLayerActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   EditBusinessLayer,
		Name:                 ActionIDNameMap[EditBusinessLayer],
		NameEn:               "Edit Business Level",
		Type:                 Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       []ActionID{ViewModelTopo},
		Version:              1,
	})
	return actions
}

func genModelTopologyViewActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:      ViewModelTopo,
		Name:    ActionIDNameMap[ViewModelTopo],
		NameEn:  "View Model Topo",
		Type:    View,
		Version: 1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditModelTopologyView,
		Name:                 ActionIDNameMap[EditModelTopologyView],
		NameEn:               "Edit Model Topo View",
		Type:                 Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       []ActionID{ViewModelTopo},
		Version:              1,
	})
	return actions
}

func genOperationStatisticActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   FindOperationStatistic,
		Name:                 ActionIDNameMap[FindOperationStatistic],
		NameEn:               "View Operational Statistics",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditOperationStatistic,
		Name:                 ActionIDNameMap[EditOperationStatistic],
		NameEn:               "Edit Operational Statistics",
		Type:                 Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       []ActionID{FindOperationStatistic},
		Version:              1,
	})

	return actions
}

func genAuditLogActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   FindAuditLog,
		Name:                 ActionIDNameMap[FindAuditLog],
		NameEn:               "View Operation Audit",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})
	return actions
}

func genEventWatchActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   WatchHostEvent,
		Name:                 ActionIDNameMap[WatchHostEvent],
		NameEn:               "Host Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   WatchHostRelationEvent,
		Name:                 ActionIDNameMap[WatchHostRelationEvent],
		NameEn:               "Host Relation Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   WatchBizEvent,
		Name:                 ActionIDNameMap[WatchBizEvent],
		NameEn:               "Business Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   WatchSetEvent,
		Name:                 ActionIDNameMap[WatchSetEvent],
		NameEn:               "Set Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   WatchModuleEvent,
		Name:                 ActionIDNameMap[WatchModuleEvent],
		NameEn:               "Module Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   WatchProcessEvent,
		Name:                 ActionIDNameMap[WatchProcessEvent],
		NameEn:               "Process Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:      WatchBizSetEvent,
		Name:    ActionIDNameMap[WatchBizSetEvent],
		NameEn:  "Business Set Event Listen",
		Type:    View,
		Version: 1,
	})

	actions = append(actions, ResourceAction{
		ID:      WatchPlatEvent,
		Name:    ActionIDNameMap[WatchPlatEvent],
		NameEn:  "Cloud Area Event Listen",
		Type:    View,
		Version: 1,
	})

	actions = append(actions, ResourceAction{
		ID:      WatchProjectEvent,
		Name:    ActionIDNameMap[WatchProjectEvent],
		NameEn:  "Project Event Listen",
		Type:    View,
		Version: 1,
	})

	modelSelection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysModelEventSelection,
	}}

	modelResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysModelEvent,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: modelSelection,
		},
	}

	actions = append(actions, ResourceAction{
		ID:                   WatchCommonInstanceEvent,
		Name:                 ActionIDNameMap[WatchCommonInstanceEvent],
		NameEn:               "Common Model Instance Event Listen",
		Type:                 View,
		RelatedResourceTypes: modelResource,
		RelatedActions:       nil,
		Version:              1,
	})

	mainlineModelSelection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       MainlineModelEventSelection,
	}}

	mainlineModelResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 MainlineModelEvent,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: mainlineModelSelection,
		},
	}

	actions = append(actions, ResourceAction{
		ID:                   WatchMainlineInstanceEvent,
		Name:                 ActionIDNameMap[WatchMainlineInstanceEvent],
		NameEn:               "Custom Topo Layer Event Listen",
		Type:                 View,
		RelatedResourceTypes: mainlineModelResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:     WatchInstAsstEvent,
		Name:   ActionIDNameMap[WatchInstAsstEvent],
		NameEn: "Instance Association Event Listen",
		Type:   View,
		RelatedResourceTypes: []RelateResourceType{
			{
				SystemID:    SystemIDCMDB,
				ID:          InstAsstEvent,
				NameAlias:   "",
				NameAliasEn: "",
				Scope:       nil,
				InstanceSelections: []RelatedInstanceSelection{{
					SystemID: SystemIDCMDB,
					ID:       InstAsstEventSelection,
				}},
			},
		},
		RelatedActions: nil,
		Version:        1,
	})
	return actions
}

func genKubeEventWatchActions() []ResourceAction {
	return []ResourceAction{
		{
			ID:      WatchKubeClusterEvent,
			Name:    ActionIDNameMap[WatchKubeClusterEvent],
			NameEn:  "Kube Cluster Event Listen",
			Type:    View,
			Version: 1,
		},
		{
			ID:      WatchKubeNodeEvent,
			Name:    ActionIDNameMap[WatchKubeNodeEvent],
			NameEn:  "Kube Node Event Listen",
			Type:    View,
			Version: 1,
		},
		{
			ID:      WatchKubeNamespaceEvent,
			Name:    ActionIDNameMap[WatchKubeNamespaceEvent],
			NameEn:  "Kube Namespace Event Listen",
			Type:    View,
			Version: 1,
		},
		{
			ID:     WatchKubeWorkloadEvent,
			Name:   ActionIDNameMap[WatchKubeWorkloadEvent],
			NameEn: "Kube Workload Event Listen",
			Type:   View,
			RelatedResourceTypes: []RelateResourceType{
				{
					SystemID: SystemIDCMDB,
					ID:       KubeWorkloadEvent,
					InstanceSelections: []RelatedInstanceSelection{{
						SystemID: SystemIDCMDB,
						ID:       KubeWorkloadEventSelection,
					}},
				},
			},
			Version: 1,
		},
		{
			ID:      WatchKubePodEvent,
			Name:    ActionIDNameMap[WatchKubePodEvent],
			NameEn:  "Kube Pod Event Listen",
			Type:    View,
			Version: 1,
		},
	}
}

func genConfigAdminActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   GlobalSettings,
		Name:                 ActionIDNameMap[GlobalSettings],
		NameEn:               "Global Settings",
		Type:                 Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})
	return actions
}

func genContainerManagementActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, genContainerClusterActions()...)
	actions = append(actions, genContainerNodeActions()...)
	actions = append(actions, genContainerNamespaceActions()...)
	actions = append(actions, genContainerWorkloadActions()...)
	actions = append(actions, genContainerPodActions()...)

	return actions
}

func genContainerClusterActions() []ResourceAction {
	return []ResourceAction{
		{
			ID:      CreateContainerCluster,
			Name:    ActionIDNameMap[CreateContainerCluster],
			NameEn:  "Create Container Cluster",
			Type:    Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      EditContainerCluster,
			Name:    ActionIDNameMap[EditContainerCluster],
			NameEn:  "Edit Container Cluster",
			Type:    Edit,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      DeleteContainerCluster,
			Name:    ActionIDNameMap[DeleteContainerCluster],
			NameEn:  "Delete Container Cluster",
			Type:    Delete,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genContainerNodeActions() []ResourceAction {
	return []ResourceAction{
		{
			ID:      CreateContainerNode,
			Name:    ActionIDNameMap[CreateContainerNode],
			NameEn:  "Create Container Node",
			Type:    Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      EditContainerNode,
			Name:    ActionIDNameMap[EditContainerNode],
			NameEn:  "Edit Container Node",
			Type:    Edit,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      DeleteContainerNode,
			Name:    ActionIDNameMap[DeleteContainerNode],
			NameEn:  "Delete Container Node",
			Type:    Delete,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genContainerNamespaceActions() []ResourceAction {
	return []ResourceAction{
		{
			ID:      CreateContainerNamespace,
			Name:    ActionIDNameMap[CreateContainerNamespace],
			NameEn:  "Create Container Namespace",
			Type:    Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      EditContainerNamespace,
			Name:    ActionIDNameMap[EditContainerNamespace],
			NameEn:  "Edit Container Namespace",
			Type:    Edit,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      DeleteContainerNamespace,
			Name:    ActionIDNameMap[DeleteContainerNamespace],
			NameEn:  "Delete Container Namespace",
			Type:    Delete,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genContainerWorkloadActions() []ResourceAction {
	return []ResourceAction{
		{
			ID:      CreateContainerWorkload,
			Name:    ActionIDNameMap[CreateContainerWorkload],
			NameEn:  "Create Container Workload",
			Type:    Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      EditContainerWorkload,
			Name:    ActionIDNameMap[EditContainerWorkload],
			NameEn:  "Edit Container Workload",
			Type:    Edit,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      DeleteContainerWorkload,
			Name:    ActionIDNameMap[DeleteContainerWorkload],
			NameEn:  "Delete Container Workload",
			Type:    Delete,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genContainerPodActions() []ResourceAction {
	return []ResourceAction{
		{
			ID:      CreateContainerPod,
			Name:    ActionIDNameMap[CreateContainerPod],
			NameEn:  "Create Container Pod",
			Type:    Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      DeleteContainerPod,
			Name:    ActionIDNameMap[DeleteContainerPod],
			NameEn:  "Delete Container Pod",
			Type:    Delete,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genFulltextSearchActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:      UseFulltextSearch,
		Name:    ActionIDNameMap[UseFulltextSearch],
		NameEn:  "Fulltext Search",
		Type:    View,
		Version: 1,
	})
	return actions
}

func genFieldGroupingTemplateActions() []ResourceAction {
	templateResource := RelateResourceType{
		SystemID: SystemIDCMDB,
		ID:       FieldGroupingTemplate,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       FieldGroupingTemplateSelection,
		}},
	}

	return []ResourceAction{
		{
			ID:      CreateFieldGroupingTemplate,
			Name:    ActionIDNameMap[CreateFieldGroupingTemplate],
			NameEn:  "Create Field Grouping Template",
			Type:    Create,
			Version: 1,
		},
		{
			ID:                   ViewFieldGroupingTemplate,
			Name:                 ActionIDNameMap[ViewFieldGroupingTemplate],
			NameEn:               "View Field Grouping Template",
			Type:                 View,
			RelatedResourceTypes: []RelateResourceType{templateResource},
			Version:              1,
		},
		{
			ID:                   EditFieldGroupingTemplate,
			Name:                 ActionIDNameMap[EditFieldGroupingTemplate],
			NameEn:               "Edit Field Grouping Template",
			Type:                 Edit,
			RelatedResourceTypes: []RelateResourceType{templateResource},
			RelatedActions:       []ActionID{ViewFieldGroupingTemplate},
			Version:              1,
		},
		{
			ID:                   DeleteFieldGroupingTemplate,
			Name:                 ActionIDNameMap[DeleteFieldGroupingTemplate],
			NameEn:               "Delete Field Grouping Template",
			Type:                 Delete,
			RelatedResourceTypes: []RelateResourceType{templateResource},
			RelatedActions:       []ActionID{ViewFieldGroupingTemplate},
			Version:              1,
		},
	}
}

func genIDRuleActions() []ResourceAction {
	return []ResourceAction{
		{
			ID:      EditIDRuleIncrID,
			Name:    ActionIDNameMap[EditIDRuleIncrID],
			NameEn:  "Edit ID Rule",
			Type:    Edit,
			Version: 1,
		},
	}
}

func genFullSyncCondActions() []ResourceAction {
	return []ResourceAction{
		{
			ID:      CreateFullSyncCond,
			Name:    ActionIDNameMap[CreateFullSyncCond],
			NameEn:  "Create Full Sync Cond",
			Type:    Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      EditFullSyncCond,
			Name:    ActionIDNameMap[EditFullSyncCond],
			NameEn:  "Edit Full Sync Cond",
			Type:    Edit,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      DeleteFullSyncCond,
			Name:    ActionIDNameMap[DeleteFullSyncCond],
			NameEn:  "Delete Full Sync Cond",
			Type:    Delete,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      ViewFullSyncCond,
			Name:    ActionIDNameMap[ViewFullSyncCond],
			NameEn:  "View Full Sync Cond",
			Type:    View,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genCacheActions() []ResourceAction {
	return []ResourceAction{
		{
			ID:     ViewGeneralCache,
			Name:   ActionIDNameMap[ViewGeneralCache],
			NameEn: "View General Resource Cache",
			Type:   View,
			RelatedResourceTypes: []RelateResourceType{{
				SystemID: SystemIDCMDB,
				ID:       GeneralCache,
				InstanceSelections: []RelatedInstanceSelection{{
					SystemID: SystemIDCMDB,
					ID:       GeneralCacheSelection,
				}},
			}},
			Version: 1,
			Hidden:  true,
		},
	}
}
