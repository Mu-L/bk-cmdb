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

package parser

import (
	"context"
	"fmt"
	"sync"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// objectBase is subset of metadata.Object
type objectBase struct {
	ID       int64  `field:"id" json:"id" bson:"id"`
	ObjCls   string `field:"bk_classification_id" json:"bk_classification_id" bson:"bk_classification_id"`
	ObjectID string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
}

func (ps *parseStream) getOneModel(cond mapstr.MapStr) (objectBase, error) {
	model := objectBase{}

	models, err := ps.searchModels(cond)
	if err != nil {
		return model, err
	}
	if len(models) == 0 {
		return model, fmt.Errorf("model [%+v] not found", cond)
	}
	if len(models) > 1 {
		return model, fmt.Errorf("model [%+v] not found", cond)
	}

	model = objectBase{
		ID:       models[0].ID,
		ObjCls:   models[0].ObjCls,
		ObjectID: models[0].ObjectID,
	}
	return model, nil
}

func (ps *parseStream) getClassIDWithObject(obj string) (int64, error) {
	filter := mapstr.MapStr{
		common.BKObjIDField: obj,
	}

	model, err := ps.getOneModel(filter)
	if err != nil {
		return 0, err
	}

	filter = map[string]interface{}{
		common.BKClassificationIDField: model.ObjCls,
	}

	class, err := ps.getOneClassification(filter)
	if err != nil {
		return 0, err
	}

	return class.ID, nil
}

func (ps *parseStream) searchModels(cond mapstr.MapStr) ([]metadata.Object, error) {

	model, err := ps.engine.CoreAPI.CoreService().Model().ReadModel(context.Background(), ps.RequestCtx.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return nil, err
	}

	if len(model.Info) <= 0 {
		return nil, fmt.Errorf("model [%+v] not found", cond)
	}

	return model.Info, nil
}

func (ps *parseStream) getModelAttribute(bizID int64, cond mapstr.MapStr) ([]metadata.Attribute, error) {

	attr, err := ps.engine.CoreAPI.CoreService().Model().ReadModelAttrsWithTableByCondition(context.Background(),
		ps.RequestCtx.Header, bizID, &metadata.QueryCondition{Condition: cond})
	if err != nil {
		return nil, err
	}

	if len(attr.Info) <= 0 {
		return nil, fmt.Errorf("attribute [%+v] not found", cond)
	}

	attrs := make([]metadata.Attribute, 0)
	for _, info := range attr.Info {
		attrs = append(attrs, info)
	}
	return attrs, nil
}

func (ps *parseStream) getCls(clsID string) (metadata.Classification, error) {

	model, err := ps.engine.CoreAPI.CoreService().Model().ReadModelClassification(context.Background(),
		ps.RequestCtx.Header, &metadata.QueryCondition{Condition: mapstr.MapStr{common.BKClassificationIDField: clsID}})
	if err != nil {
		return metadata.Classification{}, err
	}

	if len(model.Info) <= 0 {
		return metadata.Classification{}, fmt.Errorf("model classification [%s] not found", clsID)
	}
	return model.Info[0], nil
}

func (ps *parseStream) getAttributeGroup(cond interface{}) ([]metadata.Group, error) {

	mspstrCond, err := mapstr.NewFromInterface(cond)
	if err != nil {
		return nil, err
	}

	groups, err := ps.engine.CoreAPI.CoreService().Model().ReadAttributeGroupByCondition(context.Background(),
		ps.RequestCtx.Header, metadata.QueryCondition{Condition: mspstrCond})
	if err != nil {
		return nil, err
	}

	return groups.Info, nil
}

func (ps *parseStream) isMainlineModel(modelID string) (bool, error) {
	asst, err := ps.getModelAssociation(mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline})
	if err != nil {
		return false, err
	}
	for _, mainline := range asst {
		if mainline.ObjectID == modelID {
			return true, nil
		}
	}

	return false, nil
}

func (ps *parseStream) getModelAssociation(cond mapstr.MapStr) ([]metadata.Association, error) {

	asst, err := ps.engine.CoreAPI.CoreService().Association().ReadModelAssociation(context.Background(),
		ps.RequestCtx.Header, &metadata.QueryCondition{Condition: cond})
	if err != nil {
		return nil, err
	}

	if len(asst.Info) <= 0 {
		return nil, fmt.Errorf("model association [%+v] not found", cond)
	}

	return asst.Info, nil
}

func (ps *parseStream) getInstAssociation(objID string, cond mapstr.MapStr) (metadata.InstAsst, error) {

	asst, err := ps.engine.CoreAPI.CoreService().Association().ReadInstAssociation(context.Background(),
		ps.RequestCtx.Header, &metadata.InstAsstQueryCondition{
			Cond:  metadata.QueryCondition{Condition: cond},
			ObjID: objID,
		})
	if err != nil {
		return metadata.InstAsst{}, err
	}

	if len(asst.Info) <= 0 {
		return metadata.InstAsst{}, fmt.Errorf("model association [%+v] not found", cond)
	}

	return asst.Info[0], nil
}

func (ps *parseStream) getInstanceTypeByObject(objID string, ID int64) (meta.ResourceType, error) {
	switch objID {
	case common.BKInnerObjIDPlat:
		return meta.CloudAreaInstance, nil
	case common.BKInnerObjIDHost:
		return meta.HostInstance, nil
	case common.BKInnerObjIDModule:
		return meta.ModelModule, nil
	case common.BKInnerObjIDSet:
		return meta.ModelSet, nil
	case common.BKInnerObjIDApp:
		return meta.Business, nil
	case common.BKInnerObjIDProc:
		return meta.Process, nil
	case common.BKInnerObjIDBizSet:
		return meta.BizSet, nil
	case common.BKInnerObjIDProject:
		return meta.Project, nil

	}
	isMainline, err := ps.isMainlineModel(objID)
	if err != nil {
		return "", err
	}
	if isMainline {
		return meta.MainlineInstance, nil
	}
	return iam.GenCMDBDynamicResType(ID), nil
}

func (ps *parseStream) getBizIDByHostID(hostID int64) (int64, error) {

	result, err := ps.engine.CoreAPI.CoreService().Host().GetHostModuleRelation(context.Background(),
		ps.RequestCtx.Header, &metadata.HostModuleRelationRequest{
			HostIDArr: []int64{hostID},
			Fields:    []string{common.BKAppIDField},
		})
	if err != nil {
		return 0, err
	}

	for _, relation := range result.Info {
		return relation.AppID, nil
	}
	return 0, nil
}

func (ps *parseStream) getOneClassification(cond mapstr.MapStr) (metadata.Classification, error) {
	classification := metadata.Classification{}
	classifications, err := ps.getClassification(cond)
	if err != nil {
		return classification, err
	}
	if len(classifications) <= 0 {
		return classification, fmt.Errorf("classification [%+v] not found", cond)
	}
	if len(classifications) > 1 {
		return classification, fmt.Errorf("get multiple classification with [%+v]", cond)
	}
	return classifications[0], nil
}

func (ps *parseStream) getClassification(cond mapstr.MapStr) ([]metadata.Classification, error) {
	classificationResult, err := ps.engine.CoreAPI.CoreService().Model().ReadModelClassification(context.Background(),
		ps.RequestCtx.Header, &metadata.QueryCondition{Condition: cond})
	if err != nil {
		return nil, err
	}

	if len(classificationResult.Info) <= 0 {
		return nil, fmt.Errorf("classification [%+v] not found", cond)
	}

	return classificationResult.Info, nil
}

func (ps *parseStream) getModelUnique(cond mapstr.MapStr) (metadata.ObjectUnique, error) {
	unique := metadata.ObjectUnique{}
	filter := metadata.QueryCondition{Condition: cond}
	modelUniqueResult, err := ps.engine.CoreAPI.CoreService().Model().ReadModelAttrUnique(context.Background(),
		ps.RequestCtx.Header, filter)
	if err != nil {
		return unique, err
	}

	if len(modelUniqueResult.Info) <= 0 {
		return unique, fmt.Errorf("model unique [%+v] not found", cond)
	}

	if len(modelUniqueResult.Info) > 1 {
		return unique, fmt.Errorf("get multiple model unique with [%+v]", cond)
	}
	return modelUniqueResult.Info[0], nil
}

// getRscPoolHostModuleRelation TODO
// get hosts relation which these hosts must be in the resource pool
func (ps *parseStream) getRscPoolHostModuleRelation(hostIDs []int64) (map[int64]int64, error) {
	opt := &metadata.HostModuleRelationRequest{
		HostIDArr: hostIDs,
		Fields:    []string{common.BKAppIDField, common.BKHostIDField, common.BKModuleIDField},
	}

	result, err := ps.engine.CoreAPI.CoreService().Host().GetHostModuleRelation(context.Background(),
		ps.RequestCtx.Header, opt)
	if err != nil {
		return nil, err
	}

	resourceBiz, err := ps.getResourcePoolBusinessID()
	if err != nil {
		return nil, err
	}

	relation := make(map[int64]int64)
	for _, rel := range result.Info {
		if rel.AppID != resourceBiz {
			return nil, errors.New(common.CCErrCommParamsInvalid, "host does not belongs to host pool")
		}

		relation[rel.HostID] = rel.ModuleID
	}

	return relation, nil
}

var resourcePoolBizIDMap = sync.Map{}

func (ps *parseStream) getResourcePoolBusinessID() (int64, error) {
	supplierAccount := httpheader.GetSupplierAccount(ps.RequestCtx.Header)

	if bizID, ok := resourcePoolBizIDMap.Load(supplierAccount); ok {
		return util.GetInt64ByInterface(bizID)
	}

	opt := &metadata.QueryCondition{
		Fields: []string{common.BKAppIDField, common.BkSupplierAccount},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{
			common.BkSupplierAccount: supplierAccount,
			"default":                1,
		},
	}

	result, err := ps.engine.CoreAPI.CoreService().Instance().ReadInstance(context.Background(), ps.RequestCtx.Header,
		common.BKInnerObjIDApp, opt)
	if err != nil {
		return 0, err
	}

	for _, biz := range result.Info {
		bizSupplierAccount, err := biz.String(common.BkSupplierAccount)
		if err != nil {
			return 0, err
		}

		if bizSupplierAccount == supplierAccount {
			id, err := util.GetInt64ByInterface(biz[common.BKAppIDField])
			if err != nil {
				return 0, errors.New(common.CCErrorUnknownOrUnrecognizedError, "invalid resource biz id")
			}

			resourcePoolBizIDMap.Store(supplierAccount, id)
			return id, nil
		}
	}

	return 0, errors.New(common.CCErrCommParamsIsInvalid, "biz with the supplier account does not exist")
}

type hostPool struct {
	Business int64   `json:"bk_biz_id"`
	HostID   []int64 `json:"bk_host_id"`
}

var resourcePoolDefaultDirIDMap = sync.Map{}

func (ps *parseStream) getResourcePoolDefaultDirID() (dirID int64, err error) {
	supplierAccount := httpheader.GetSupplierAccount(ps.RequestCtx.Header)

	if dirID, ok := resourcePoolDefaultDirIDMap.Load(supplierAccount); ok {
		return util.GetInt64ByInterface(dirID)
	}

	bizID, err := ps.getResourcePoolBusinessID()
	if err != nil {
		return 0, err
	}

	opt := &metadata.QueryCondition{
		Fields: []string{common.BKModuleIDField, common.BkSupplierAccount},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{
			common.BKDefaultField:    common.DefaultResModuleFlag,
			common.BKAppIDField:      bizID,
			common.BkSupplierAccount: supplierAccount,
		},
	}

	result, err := ps.engine.CoreAPI.CoreService().Instance().ReadInstance(context.Background(), ps.RequestCtx.Header,
		common.BKInnerObjIDModule, opt)
	if err != nil {
		return 0, err
	}

	for _, directory := range result.Info {
		dirSupplierAccount, err := directory.String(common.BkSupplierAccount)
		if err != nil {
			return 0, err
		}

		if dirSupplierAccount == supplierAccount {
			id, err := util.GetInt64ByInterface(directory[common.BKModuleIDField])
			if err != nil {
				return 0, errors.New(common.CCErrorUnknownOrUnrecognizedError,
					"invalid resource pool default directory id")
			}

			resourcePoolDefaultDirIDMap.Store(supplierAccount, id)
			return id, nil
		}
	}
	return 0, errors.New(common.CCErrCommParamsIsInvalid,
		"directory with the supplier account does not exist")
}

// generateUpdateInstanceResource generate update instance auth resource by their object type
func (ps *parseStream) generateUpdateInstanceResource(model *metadata.Object, instID int64) (*meta.ResourceAttribute,
	error) {

	instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
	if err != nil {
		return nil, err
	}

	switch instanceType {
	case meta.HostInstance:
		hostRelationReq := &metadata.HostModuleRelationRequest{
			HostIDArr: []int64{instID},
			Fields:    []string{common.BKAppIDField, common.BKModuleIDField},
		}

		relationRes, err := ps.engine.CoreAPI.CoreService().Host().GetHostModuleRelation(context.Background(),
			ps.RequestCtx.Header, hostRelationReq)
		if err != nil {
			return nil, err
		}

		if len(relationRes.Info) == 0 {
			return nil, errors.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("host %d has no relations", instID))
		}
		bizID := relationRes.Info[0].AppID

		rscPoolBizID, err := ps.getResourcePoolBusinessID()
		if err != nil {
			return nil, err
		}

		if bizID != rscPoolBizID {
			return &meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:       instanceType,
					Action:     meta.Update,
					InstanceID: instID,
				},
				Layers:     []meta.Item{{Type: meta.Business, InstanceID: bizID}},
				BusinessID: bizID,
			}, nil
		}

		if len(relationRes.Info) > 1 {
			return nil, errors.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("host %d is in many dirs", instID))
		}

		return &meta.ResourceAttribute{
			Basic: meta.Basic{
				Type:       instanceType,
				Action:     meta.Update,
				InstanceID: instID,
			},
			Layers: []meta.Item{{Type: meta.ResourcePoolDirectory, InstanceID: relationRes.Info[0].ModuleID}},
		}, nil
	case meta.Business, meta.BizSet, meta.CloudAreaInstance, meta.Project:
		return &meta.ResourceAttribute{
			Basic: meta.Basic{
				Type:       instanceType,
				Action:     meta.Update,
				InstanceID: instID,
			},
		}, nil
	case meta.ModelSet, meta.ModelModule, meta.MainlineInstance, meta.Process:
		bizIDReq := &metadata.QueryCondition{
			Fields:    []string{common.BKAppIDField},
			Page:      metadata.BasePage{Limit: 1},
			Condition: mapstr.MapStr{metadata.GetInstIDFieldByObjID(model.ObjectID): instID},
		}

		bizIDResult, err := ps.engine.CoreAPI.CoreService().Instance().ReadInstance(context.Background(),
			ps.RequestCtx.Header, model.ObjectID, bizIDReq)
		if err != nil {
			return nil, err
		}

		if len(bizIDResult.Info) == 0 {
			return nil, errors.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("host %d has no relations", instID))
		}

		bizID, err := util.GetInt64ByInterface(bizIDResult.Info[0][common.BKAppIDField])
		if err != nil {
			return nil, errors.New(common.CCErrCommParamsInvalid, fmt.Sprintf("inst %d has invalid biz id", instID))
		}

		return &meta.ResourceAttribute{
			BusinessID: bizID,
			Basic: meta.Basic{
				Type:   instanceType,
				Action: meta.Update,
			},
		}, nil
	default:
		if iam.IsCMDBSysInstance(instanceType) {
			return &meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:       instanceType,
					Action:     meta.Update,
					InstanceID: instID,
				},
			}, nil
		}
		return nil, errors.New(common.CCErrCommParamsIsInvalid,
			fmt.Sprintf("instance type %s is invalid", instanceType))
	}
}
