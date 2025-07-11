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

package process

import (
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// CreateServiceTemplate TODO
func (p *processOperation) CreateServiceTemplate(kit *rest.Kit, template metadata.ServiceTemplate) (*metadata.ServiceTemplate, errors.CCErrorCoder) {
	// base attribute validate
	if field, err := template.Validate(kit.CCError); err != nil {
		blog.Errorf("CreateServiceTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(kit, template.BizID); err != nil {
		blog.Errorf("CreateServiceTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	// keep metadata clean
	template.BizID = bizID

	// validate service category id field
	category, err := p.GetServiceCategory(kit, template.ServiceCategoryID)
	if err != nil {
		blog.Errorf("CreateServiceTemplate failed, category id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "service_category_id")
	}
	isLeafNode, ccErr := p.IsServiceCategoryLeafNode(kit, category.ID)
	if ccErr != nil {
		blog.Errorf("UpdateServiceTemplate failed, check leaf node failed, err: %+v, rid: %s", ccErr, kit.Rid)
		return nil, ccErr
	}
	if !isLeafNode {
		return nil, kit.CCError.CCError(common.CCErrCoreServiceOnlyNodeServiceCategoryAvailable)
	}

	// make sure biz id identical with category
	// categoryBizID 0 and 1 is default category
	if bizID != category.BizID && category.BizID != 0 {
		blog.Errorf("CreateServiceTemplate failed, validation failed, input bizID:%d not equal category bizID:%d, rid: %s", bizID, category.BizID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	// check name field unique under business
	nameUniqueFilter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKFieldName:  template.Name,
	}
	count, err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(nameUniqueFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("CreateServiceTemplate failed, count same name instance failed, filter: %+v, err: %+v, rid: %s", nameUniqueFilter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("CreateServiceTemplate failed, service instance name duplicated, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKFieldName)
	}

	// generate id field
	id, err := mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameServiceTemplate)
	if nil != err {
		blog.Errorf("CreateServiceTemplate failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	template.ID = int64(id)

	template.Creator = kit.User
	template.Modifier = kit.User
	template.CreateTime = time.Now()
	template.LastTime = time.Now()
	template.SupplierAccount = kit.SupplierAccount

	if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Insert(kit.Ctx, &template); nil != err {
		blog.Errorf("CreateServiceTemplate failed, mongodb failed, table: %s, template: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, template, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}
	return &template, nil
}

// GetServiceTemplate TODO
func (p *processOperation) GetServiceTemplate(kit *rest.Kit, templateID int64) (*metadata.ServiceTemplate, errors.CCErrorCoder) {
	template := metadata.ServiceTemplate{}

	filter := map[string]int64{common.BKFieldID: templateID}
	if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(filter).One(kit.Ctx, &template); nil != err {
		blog.Errorf("GetServiceTemplate failed, mongodb failed, table: %s, filter: %+v, template: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, template, err, kit.Rid)
		if mongodb.Client().IsNotFoundError(err) {
			return nil, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &template, nil
}

func (p *processOperation) validateServiceCategoryID(kit *rest.Kit,
	template *metadata.ServiceTemplate) errors.CCErrorCoder {

	// validate service category id field
	category, err := p.GetServiceCategory(kit, template.ServiceCategoryID)
	if err != nil {
		blog.Errorf("get category failed, serviceCategory id: %d, err: %v, rid: %s", template.ServiceCategoryID,
			err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
	}
	if category.BizID != 0 && category.BizID != template.BizID {
		blog.Errorf("category biz id and template not equal, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
	}
	isLeafNode, err := p.IsServiceCategoryLeafNode(kit, template.ServiceCategoryID)
	if err != nil {
		blog.Errorf("check leaf node failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	if !isLeafNode {
		return kit.CCError.CCError(common.CCErrCoreServiceOnlyNodeServiceCategoryAvailable)
	}
	return nil
}

func ifUpdateModuleName(kit *rest.Kit, template *metadata.ServiceTemplate) (bool, errors.CCErrorCoder) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var checkErr errors.CCErrorCoder
	// check service template name unique
	go func() {
		defer wg.Done()
		nameFilter := map[string]interface{}{
			common.BKAppIDField: template.BizID,
			common.BKFieldName:  template.Name,
			common.BKFieldID:    map[string]interface{}{common.BKDBNE: template.ID},
		}
		count, err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(nameFilter).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("count service template with same name failed, filter: %+v, err: %v, rid: %s", nameFilter,
				err, kit.Rid)
			checkErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)
			return
		}
		if count > 0 {
			blog.Errorf("service template name duplicated, count: %d, rid: %s", count, kit.Rid)
			checkErr = kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKFieldName)
			return
		}
	}()

	needUpdateModuleName := false

	// get modules using this service template
	go func() {
		defer wg.Done()
		moduleFilter := map[string]interface{}{
			common.BKServiceTemplateIDField: template.ID,
		}
		modules := make([]metadata.ModuleInst, 0)
		err := mongodb.Client().Table(common.BKTableNameBaseModule).Find(moduleFilter).All(kit.Ctx, &modules)
		if err != nil {
			blog.Errorf("count modules using this service template failed, filter: %+v, err: %v, rid: %s", moduleFilter,
				err, kit.Rid)
			checkErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)
			return
		}
		if len(modules) > 0 {
			parentIDs := make([]int64, len(modules))
			for _, module := range modules {
				parentIDs = append(parentIDs, module.ParentID)
			}

			// check if other modules has the same name with the service template name to be changed
			moduleNameFilter := map[string]interface{}{
				common.BKAppIDField:             template.BizID,
				common.BKModuleNameField:        template.Name,
				common.BKParentIDField:          map[string]interface{}{common.BKDBIN: parentIDs},
				common.BKServiceTemplateIDField: map[string]interface{}{common.BKDBNE: template.ID},
			}
			count, err := mongodb.Client().Table(common.BKTableNameBaseModule).Find(moduleNameFilter).Count(kit.Ctx)
			if err != nil {
				blog.Errorf("count modules with same name failed, filter: %+v, err: %v, rid: %s", moduleFilter, err,
					kit.Rid)
				checkErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)
				return
			}
			if count > 0 {
				blog.Errorf("service template has modules with same name, count: %d, rid: %s", count, kit.Rid)
				checkErr = kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKFieldName)
				return
			}
			needUpdateModuleName = true
		}
	}()

	wg.Wait()
	if checkErr != nil {
		return false, checkErr
	}
	return needUpdateModuleName, nil
}

// UpdateServiceTemplate not support update name field yet, so don't need validate name unique before update
func (p *processOperation) UpdateServiceTemplate(kit *rest.Kit, templateID int64,
	input metadata.ServiceTemplate) (*metadata.ServiceTemplate, errors.CCErrorCoder) {

	template, err := p.GetServiceTemplate(kit, templateID)
	if err != nil {
		return nil, err
	}

	// update fields to local object and do validation
	if input.ServiceCategoryID > 0 {
		template.ServiceCategoryID = input.ServiceCategoryID
		if err := p.validateServiceCategoryID(kit, template); err != nil {
			return nil, err
		}
	}

	needUpdateModuleName := false
	if len(input.Name) != 0 && template.Name != input.Name {
		template.Name = input.Name

		needUpdateModuleName, err = ifUpdateModuleName(kit, template)
		if err != nil {
			return nil, err
		}
	}

	if field, err := template.Validate(kit.CCError); err != nil {
		blog.Errorf("validate template(%+v) failed, err: %+v, rid: %s", template, err, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	// do update
	filter := map[string]int64{common.BKFieldID: templateID}
	if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Update(kit.Ctx, filter, template); err != nil {
		blog.Errorf("update service template failed, err: %+v, filter: %+v, data: %+v, rid: %s", err, filter,
			template, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
	}

	updateData := make(map[string]interface{})
	if needUpdateModuleName {
		updateData[common.BKModuleNameField] = template.Name
	}

	// update module category if there is a difference
	if input.ServiceCategoryID != 0 {
		option := map[string]interface{}{
			common.BKServiceTemplateIDField: template.ID,
			common.BKServiceCategoryIDField: mapstr.MapStr{
				common.BKDBNE: input.ServiceCategoryID,
			},
		}
		cnt, e := mongodb.Client().Table(common.BKTableNameBaseModule).Find(option).Count(kit.Ctx)
		if e != nil {
			blog.Errorf("count module failed, filter: %+v, err: %v, rid: %s", option, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}

		if cnt > 0 {
			updateData[common.BKServiceCategoryIDField] = template.ServiceCategoryID
		}
	}

	if len(updateData) == 0 {
		return template, nil
	}

	// update name & category of the modules using this service template
	moduleFilter := map[string]interface{}{common.BKServiceTemplateIDField: template.ID}
	rawErr := mongodb.Client().Table(common.BKTableNameBaseModule).Update(kit.Ctx, moduleFilter, updateData)
	if rawErr != nil {
		blog.Errorf("update modules using this service template failed, err: %v, filter: %+v, rid: %s", rawErr,
			moduleFilter, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return template, nil
}

// ListServiceTemplates TODO
func (p *processOperation) ListServiceTemplates(kit *rest.Kit, option metadata.ListServiceTemplateOption) (
	*metadata.MultipleServiceTemplate, errors.CCErrorCoder) {

	filter := map[string]interface{}{}
	
	if option.BusinessID != 0 {
		filter[common.BKAppIDField] = option.BusinessID
	}

	// filter with matching any sub category
	if option.ServiceCategoryID != nil && *option.ServiceCategoryID > 0 {
		categoriesWithSts, err := p.ListServiceCategories(kit, option.BusinessID, false)
		if err != nil {
			blog.Errorf("ListServiceTemplates failed, ListServiceCategories failed, err: %+v, rid: %s", err, kit.Rid)
			return nil, err
		}
		childrenIDs := make([]int64, 0)
		childrenIDs = append(childrenIDs, *option.ServiceCategoryID)
		for {
			pre := len(childrenIDs)
			for _, categoryWithSts := range categoriesWithSts.Info {
				category := categoryWithSts.ServiceCategory
				if category.ParentID == 0 {
					continue
				}
				if util.InArray(category.ParentID, childrenIDs) && !util.InArray(category.ID, childrenIDs) {
					childrenIDs = append(childrenIDs, category.ID)
				}
			}
			if pre == len(childrenIDs) {
				break
			}
		}
		filter[common.BKServiceCategoryIDField] = map[string][]int64{
			common.BKDBIN: childrenIDs,
		}
	}

	if option.ServiceTemplateIDs != nil && len(option.ServiceTemplateIDs) != 0 {
		filter[common.BKFieldID] = map[string][]int64{
			common.BKDBIN: option.ServiceTemplateIDs,
		}
	}

	if len(option.Search) > 0 {
		if option.IsExact {
			filter[common.BKFieldName] = option.Search
		} else {
			filter[common.BKFieldName] = map[string]interface{}{
				common.BKDBLIKE:    params.SpecialCharChange(option.Search),
				common.BKDBOPTIONS: "i",
			}
		}
	}

	var total uint64
	var err error
	if total, err = mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(filter).Count(kit.Ctx); nil != err {
		blog.Errorf("ListServiceTemplates failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	sort := "-id"
	if len(option.Page.Sort) > 0 {
		sort = option.Page.Sort
	}
	templates := make([]metadata.ServiceTemplate, 0)
	if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(filter).Start(uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).Sort(sort).All(kit.Ctx, &templates); nil != err {
		blog.Errorf("ListServiceTemplates failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	result := &metadata.MultipleServiceTemplate{
		Count: total,
		Info:  templates,
	}
	return result, nil
}

// DeleteServiceTemplate TODO
func (p *processOperation) DeleteServiceTemplate(kit *rest.Kit, serviceTemplateID int64) errors.CCErrorCoder {
	template, err := p.GetServiceTemplate(kit, serviceTemplateID)
	if err != nil {
		blog.Errorf("DeleteServiceTemplate failed, GetServiceTemplate failed, templateID: %d, err: %+v, rid: %s", serviceTemplateID, err, kit.Rid)
		return err
	}

	// service template that referenced by process template shouldn't be removed
	usageFilter := map[string]int64{
		common.BKServiceTemplateIDField: template.ID,
	}
	usageCount, e := mongodb.Client().Table(common.BKTableNameServiceInstance).Find(usageFilter).Count(kit.Ctx)
	if nil != e {
		blog.Errorf("DeleteServiceTemplate failed, mongodb failed, table: %s, process template usageFilter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, usageFilter, e, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if usageCount > 0 {
		blog.Errorf("DeleteServiceTemplate failed, forbidden delete service template be referenced, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, kit.Rid)
		err := kit.CCError.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	// service template that referenced by module shouldn't be removed
	usageCount, e = mongodb.Client().Table(common.BKTableNameBaseModule).Find(usageFilter).Count(kit.Ctx)
	if nil != e {
		blog.Errorf("DeleteServiceTemplate failed, mongodb failed, table: %s, module usageFilter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, usageFilter, e, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if usageCount > 0 {
		blog.Errorf("DeleteServiceTemplate failed, forbidden delete service template be referenced, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, kit.Rid)
		err := kit.CCError.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	delAttrFilter := map[string]int64{common.BKServiceTemplateIDField: template.ID}
	if err := mongodb.Client().Table(common.BKTableNameServiceTemplateAttr).Delete(kit.Ctx, delAttrFilter); err != nil {
		blog.Errorf("delete service template attr failed, filter: %+v, err: %v, rid: %s", delAttrFilter, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	deleteFilter := map[string]int64{common.BKFieldID: template.ID}
	if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Delete(kit.Ctx, deleteFilter); nil != err {
		blog.Errorf("DeleteServiceTemplate failed, mongodb failed, table: %s, deleteFilter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, deleteFilter, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}
	return nil
}
