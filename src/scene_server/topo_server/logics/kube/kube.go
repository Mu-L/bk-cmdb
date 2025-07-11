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

package kube

import (
	"errors"
	"fmt"
	"sync"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/kube/types"
)

// KubeOperationInterface container cluster operation methods
type KubeOperationInterface interface {
	DeleteCluster(kit *rest.Kit, bizID int64, option *types.DeleteClusterOption) error
	BatchDeleteNode(kit *rest.Kit, bizID int64, option *types.BatchDeleteNodeOption) error
	BatchCreateNode(kit *rest.Kit, data *types.CreateNodesOption, bizID int64) ([]int64, error)
	BatchCreatePod(kit *rest.Kit, data *types.CreatePodsOption) ([]int64, error)

	GenSharedClusterListCond(kit *rest.Kit, bizID int64, cond *filter.Expression) (mapstr.MapStr, error)
	GenSharedNsListCond(kit *rest.Kit, objID string, bizID int64, cond *filter.Expression) (mapstr.MapStr, error)
	CheckPlatBizSharedNs(kit *rest.Kit, bizNsMap map[int64][]int64) error
}

// NewClusterOperation create a business instance
func NewClusterOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) KubeOperationInterface {
	return &kube{
		clientSet:   client,
		authManager: authManager,
	}
}

type kube struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	cluster     KubeOperationInterface
}

func (k *kube) getDeleteNodeInfo(kit *rest.Kit, ids []int64, bizID int64) ([]types.Node, error) {

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKFieldID: mapstr.MapStr{
				common.BKDBIN: ids,
			},
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	result, err := k.clientSet.CoreService().Kube().SearchNode(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("search node failed, filter: %+v, err: %v, rid: %s", query, err, kit.Rid)
		return nil, err
	}

	if len(result.Data) == 0 {
		return result.Data, nil
	}

	bizMap := make(map[int64]struct{})
	for _, node := range result.Data {
		bizMap[node.BizID] = struct{}{}
	}

	if len(bizMap) > 1 {
		blog.Errorf("node ids exist in different businesses, filter: %+v, rid: %s", query, kit.Rid)
		return nil, errors.New("node ids exist in different businesses")
	}

	if _, ok := bizMap[bizID]; !ok {
		blog.Errorf("node ids not in biz %d, filter: %+v, rid: %s", bizID, query, kit.Rid)
		return nil, fmt.Errorf("node ids not in biz %d", bizID)
	}

	return result.Data, nil
}

// BatchDeleteNode batch delete node.
func (k *kube) BatchDeleteNode(kit *rest.Kit, bizID int64, option *types.BatchDeleteNodeOption) error {

	nodes, err := k.getDeleteNodeInfo(kit, option.IDs, bizID)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return nil
	}

	// 1、check whether these nodes exist, they must all exist before they can be deleted,
	// otherwise an error will be returned.
	podCond := []map[string]interface{}{
		{
			types.BKNodeIDField: map[string]interface{}{common.BKDBIN: option.IDs},
		},
	}
	// 2、check if there is a pod on the node.
	counts, err := k.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameBasePod, podCond)
	if err != nil {
		blog.Errorf("count nodes failed, cond: %#v, err: %v, rid: %s", podCond, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed)
	}

	if counts[0] > 0 {
		blog.Errorf("count nodes failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errors.New("no pods can exist under the node"))
	}

	// 3、batch delete nodes
	err = k.clientSet.CoreService().Kube().BatchDeleteNode(kit.Ctx, kit.Header, &option.BatchDeleteNodeByIDsOption)
	if err != nil {
		blog.Errorf("delete node failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return err
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	audit := auditlog.NewKubeAudit(k.clientSet.CoreService())
	auditLog, err := audit.GenerateNodeAuditLog(generateAuditParameter, nodes)
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

func (k *kube) getDeleteClusterInfo(kit *rest.Kit, ids []int64, bizID int64) ([]types.Cluster, error) {

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKFieldID: mapstr.MapStr{
				common.BKDBIN: ids,
			},
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	result, err := k.clientSet.CoreService().Kube().SearchCluster(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("search cluster failed, filter: %+v, err: %v, rid: %s", query, err, kit.Rid)
		return nil, err
	}

	bizMap := make(map[int64]struct{})
	for _, node := range result.Data {
		bizMap[node.BizID] = struct{}{}
	}

	if len(bizMap) != 1 {
		blog.Errorf("cluster ids exist in different businesses, filter: %+v, rid: %s", query, kit.Rid)
		return nil, errors.New("node ids exist in different businesses")
	}

	if _, ok := bizMap[bizID]; !ok {
		blog.Errorf("cluster ids not in biz %d, filter: %+v, rid: %s", bizID, query, kit.Rid)
		return nil, fmt.Errorf("cluster ids not in biz %d", bizID)
	}

	return result.Data, nil
}

// DeleteCluster delete cluster.
func (k *kube) DeleteCluster(kit *rest.Kit, bizID int64, option *types.DeleteClusterOption) error {

	clusters, err := k.getDeleteClusterInfo(kit, option.IDs, bizID)
	if err != nil {
		return err
	}
	// whether the associated resources under the cluster have been deleted. such as namespace, node, deployment, pod.
	exist, cErr := k.isExistKubeResourceUnderCluster(kit, option, bizID)
	if cErr != nil {
		blog.Errorf("failed to obtain resources under the cluster, bizID: %d, cluster IDs: %+v, err: %v, rid: %s",
			bizID, option.IDs, cErr, kit.Rid)
		return cErr
	}
	if exist {
		blog.Errorf("the associated resources under the deleted cluster haven't been deleted, clusterIDs: %+v, rid: %s",
			option.IDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid)
	}

	err = k.clientSet.CoreService().Kube().DeleteCluster(kit.Ctx, kit.Header, &option.DeleteClusterByIDsOption)
	if err != nil {
		blog.Errorf("delete cluster failed, cluster IDs: %#v, err: %v, rid: %s", option.IDs, err, kit.Rid)
		return err
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	audit := auditlog.NewKubeAudit(k.clientSet.CoreService())
	auditLog, err := audit.GenerateClusterAuditLog(generateAuditParameter, clusters)
	if err != nil {
		blog.Errorf("generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return nil
}

func (k *kube) isExistKubeResourceUnderCluster(kit *rest.Kit, option *types.DeleteClusterOption, bizID int64) (
	bool, error) {

	if len(option.IDs) == 0 {
		return false, errors.New("ids must be set")
	}

	var (
		wg       sync.WaitGroup
		firstErr ccErr.CCErrorCoder
	)

	workLoads := types.GetWorkLoadTables()
	tables := []string{types.BKTableNameBaseNamespace, types.BKTableNameBaseNode, types.BKTableNameBasePod}
	tables = append(tables, workLoads...)

	filter := []map[string]interface{}{{
		types.BKClusterIDFiled: map[string]interface{}{common.BKDBIN: option.IDs},
	}}

	for _, table := range tables {
		wg.Add(1)
		go func(table string, bizID int64) {
			defer func() {
				wg.Done()
			}()

			counts, err := k.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, table, filter)
			if err != nil {
				blog.Errorf("count resource failed, cond: %#v, err: %v, rid: %s", filter, err, kit.Rid)
				firstErr = err
				return
			}
			if counts[0] > 0 {
				blog.Errorf("there are resources under the cluster that cannot be deleted, bizID: %d, filter: %+v, "+
					"table: %s, rid: %s", bizID, table, kit.Rid)
				firstErr = kit.CCError.CCErrorf(common.CCErrCommParamsInvalid)
				return
			}

		}(table, bizID)
	}
	wg.Wait()
	if firstErr != nil {
		return false, firstErr
	}

	return false, nil
}

// BatchCreatePod batch create pod.
func (k *kube) BatchCreatePod(kit *rest.Kit, data *types.CreatePodsOption) ([]int64, error) {

	// create pods and containers.
	result, err := k.clientSet.CoreService().Kube().BatchCreatePod(kit.Ctx, kit.Header, data)
	if err != nil {
		blog.Errorf("create pod failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewKubeAudit(k.clientSet.CoreService())

	podIDs := make([]int64, 0)
	for _, pod := range result {
		podIDs = append(podIDs, pod.ID)
	}
	auditLog, err := audit.GeneratePodAuditLog(generateAuditParameter, result)
	if err != nil {
		blog.Errorf("create cluster, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return podIDs, nil
}

// BatchCreateNode batch create node.
func (k *kube) BatchCreateNode(kit *rest.Kit, data *types.CreateNodesOption, bizID int64) ([]int64, error) {
	result, err := k.clientSet.CoreService().Kube().BatchCreateNode(kit.Ctx, kit.Header, data.Nodes)
	if err != nil {
		blog.Errorf("create nodes failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}
	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewKubeAudit(k.clientSet.CoreService())
	auditLog, err := audit.GenerateNodeAuditLog(generateAuditParameter, result.Info)

	if err != nil {
		blog.Errorf(" creat nodes, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create nodes, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrAuditSaveLogFailed)
	}

	ids := make([]int64, 0)
	for _, node := range result.Info {
		ids = append(ids, node.ID)
	}
	return ids, nil
}

// GenSharedClusterListCond generate cluster condition for user biz to get shared cluster of platform biz
func (k *kube) GenSharedClusterListCond(kit *rest.Kit, bizID int64, cond *filter.Expression) (mapstr.MapStr, error) {
	// get all shared clusters that are used by the user biz
	sharedClusterCond := &metadata.QueryCondition{
		Fields:    []string{types.BKClusterIDFiled},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{types.BKBizIDField: bizID},
	}

	relRes, err := k.clientSet.CoreService().Kube().ListNsSharedClusterRel(kit.Ctx, kit.Header, sharedClusterCond)
	if err != nil {
		blog.Errorf("failed to get topo path, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(relRes.Info) == 0 {
		clusterCond, err := filtertools.And(filtertools.GenAtomFilter(types.BKBizIDField, filter.Equal, bizID), cond)
		if err != nil {
			blog.Errorf("generate normal cluster cond with biz failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		mgoCond, rawErr := clusterCond.ToMgo()
		if rawErr != nil {
			blog.Errorf("parse normal cluster filter(%#v) failed, err: %v, rid: %s", clusterCond, rawErr, kit.Rid)
			return nil, rawErr
		}

		return mgoCond, nil
	}

	// add shared cluster condition to original condition
	sharedClusterIDs := make([]int64, len(relRes.Info))
	for i, rel := range relRes.Info {
		sharedClusterIDs[i] = rel.ClusterID
	}

	sharedCond := &filter.Expression{
		RuleFactory: &filter.CombinedRule{
			Condition: filter.Or,
			Rules: []filter.RuleFactory{
				&filter.AtomRule{
					Field:    types.BKBizIDField,
					Operator: filter.Equal.Factory(),
					Value:    bizID,
				},
				&filter.AtomRule{
					Field:    common.BKFieldID,
					Operator: filter.In.Factory(),
					Value:    sharedClusterIDs,
				},
			},
		},
	}

	clusterCond, rawErr := filtertools.And(sharedCond, cond)
	if rawErr != nil {
		blog.Errorf("generate shared cluster cond with biz failed, err: %v, rid: %s", rawErr, kit.Rid)
		return nil, rawErr
	}

	mgoCond, rawErr := clusterCond.ToMgo()
	if rawErr != nil {
		blog.Errorf("parse shared cluster filter(%#v) failed, err: %v, rid: %s", clusterCond, rawErr, kit.Rid)
		return nil, rawErr
	}

	return mgoCond, nil
}

// GenSharedNsListCond generate ns condition for platform biz to get resources of user biz under shared cluster
func (k *kube) GenSharedNsListCond(kit *rest.Kit, objID string, bizID int64, cond *filter.Expression) (mapstr.MapStr,
	error) {

	// get all shared namespaces that are related to the platform biz
	sharedNsCond := &metadata.QueryCondition{
		Fields:    []string{types.BKNamespaceIDField},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{types.BKAsstBizIDField: bizID},
	}

	relRes, err := k.clientSet.CoreService().Kube().ListNsSharedClusterRel(kit.Ctx, kit.Header, sharedNsCond)
	if err != nil {
		blog.Errorf("failed to get topo path, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(relRes.Info) == 0 {
		nsCond, err := filtertools.And(filtertools.GenAtomFilter(types.BKBizIDField, filter.Equal, bizID), cond)
		if err != nil {
			blog.Errorf("generate normal cluster cond with biz failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		mgoCond, rawErr := nsCond.ToMgo()
		if rawErr != nil {
			blog.Errorf("parse normal namespace filter(%#v) failed, err: %v, rid: %s", nsCond, rawErr, kit.Rid)
			return nil, rawErr
		}

		return mgoCond, nil
	}

	// add shared namespace condition to original condition
	sharedNsIDs := make([]int64, len(relRes.Info))
	for i, rel := range relRes.Info {
		sharedNsIDs[i] = rel.NamespaceID
	}

	nsField := types.BKNamespaceIDField
	if objID == types.KubeNamespace {
		nsField = types.BKIDField
	}

	sharedCond := &filter.Expression{
		RuleFactory: &filter.CombinedRule{
			Condition: filter.Or,
			Rules: []filter.RuleFactory{
				&filter.AtomRule{
					Field:    types.BKBizIDField,
					Operator: filter.Equal.Factory(),
					Value:    bizID,
				},
				&filter.AtomRule{
					Field:    nsField,
					Operator: filter.In.Factory(),
					Value:    sharedNsIDs,
				},
			},
		},
	}

	nsCond, rawErr := filtertools.And(sharedCond, cond)
	if rawErr != nil {
		blog.Errorf("generate shared namespace cond with biz failed, err: %v, rid: %s", rawErr, kit.Rid)
		return nil, rawErr
	}

	mgoCond, rawErr := nsCond.ToMgo()
	if rawErr != nil {
		blog.Errorf("parse shared namespace filter(%#v) failed, err: %v, rid: %s", nsCond, rawErr, kit.Rid)
		return nil, rawErr
	}

	return mgoCond, nil
}

// CheckPlatBizSharedNs check if platform biz's ns is in shared cluster and if its related biz matches the plat biz
func (k *kube) CheckPlatBizSharedNs(kit *rest.Kit, bizNsMap map[int64][]int64) error {
	if len(bizNsMap) == 0 {
		return nil
	}

	nsCnt := 0
	conds := make([]mapstr.MapStr, 0)

	for bizID, nsIDs := range bizNsMap {
		nsIDs = util.IntArrayUnique(nsIDs)
		nsCnt += len(nsIDs)
		conds = append(conds, mapstr.MapStr{
			types.BKAsstBizIDField:   bizID,
			types.BKNamespaceIDField: map[string]interface{}{common.BKDBIN: nsIDs},
		})
	}

	cond := []map[string]interface{}{{common.BKDBOR: conds}}

	counts, err := k.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameNsSharedClusterRel, cond)
	if err != nil {
		blog.Errorf("count shared ns failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	if len(counts) != 1 {
		blog.Errorf("shared ns count is invalid, cond: %+v, rid: %s", cond, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	if int(counts[0]) != nsCnt {
		blog.Errorf("shared ns count %d is invalid, cond: %+v, rid: %s", counts[0], cond, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	return nil
}
