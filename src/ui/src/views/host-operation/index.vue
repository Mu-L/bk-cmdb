<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="layout" v-bkloading="{
    isLoading: $loading(Object.values(request)) || loading
  }">
    <div v-show="!$loading(Object.values(request)) && !loading">
      <div class="info clearfix mb20">
        <label class="info-label fl">{{$t('已选主机')}}：</label>
        <div class="info-content">
          <i18n path="N台主机">
            <template #count><b class="info-count">{{resources.length}}</b></template>
          </i18n>
          <i class="edit-trigger icon icon-cc-edit" v-if="!changeHostDisabled" @click="handleChangeHost"></i>
        </div>
      </div>
      <div class="info clearfix mb10" v-if="type !== 'remove'">
        <label class="info-label fl">{{$t('转移到')}}：</label>
        <div class="info-content">
          <ul class="module-list">
            <li class="module-item" v-for="(id, index) in targetModules"
              :key="index"
              :class="{
                'is-business-module': type === 'business'
              }"
              v-bk-tooltips="getModulePath(id)">
              <span class="module-icon" v-if="type === 'business'">{{$i18n.locale === 'en' ? 'M' : '模'}}</span>
              {{getModuleName(id)}}
              <span class="module-mask"
                v-if="type === 'idle'"
                @click="handleChangeModule">
                {{$t('点击修改')}}
              </span>
            </li>
            <li class="module-item is-trigger"
              v-if="type === 'business'"
              @click="handleChangeModule">
              <i class="icon icon-cc-edit"></i>
            </li>
          </ul>
          <div class="module-grep"></div>
        </div>
      </div>
      <div class="info clearfix mb10" ref="changeInfo">
        <label class="info-label fl">{{$t('变更确认')}}：</label>
        <div class="info-content">
          <template v-if="availableTabList.length">
            <ul class="tab clearfix">
              <template v-for="(item, index) in availableTabList">
                <li class="tab-grep fl" v-if="index" :key="index"></li>
                <li class="tab-item fl"
                  :class="{ active: activeTab === item, 'hide-count': !item.showCount }"
                  :key="item.id"
                  @click="handleTabClick(item)">
                  <span class="tab-label">{{item.label}}</span>
                  <span v-if="item.showCount" :class="['tab-count', { 'unconfirmed': !item.confirmed }]">
                    {{item.props.info.length > 999 ? '999+' : item.props.info.length}}
                  </span>
                </li>
              </template>
            </ul>
            <component class="tab-component"
              v-for="item in availableTabList"
              v-bind="item.props"
              v-show="activeTab === item"
              :ref="item.id"
              :key="item.id"
              :is="item.component">
            </component>
          </template>
          <div class="tab-empty" v-else-if="isSameModule">
            {{$t('相同模块转移提示')}}
          </div>
          <div class="tab-empty" v-else-if="isEmptyChange">
            {{$t('无转移确认信息提示')}}
          </div>
          <div class="tab-empty" v-else>
            {{$t('无')}}
          </div>
        </div>
      </div>
      <div class="options" :class="{ 'is-sticky': hasScrollbar }">
        <bk-button theme="primary" :disabled="isSameModule" @click="handleConfrim">{{confirmText}}</bk-button>
        <bk-button class="ml10" theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </div>
    <cmdb-dialog v-model="dialog.show" :width="dialog.width" :height="dialog.height" :body-scroll="false">
      <component
        :is="dialog.component"
        :confirm-text="$t('确定')"
        v-bind="dialog.props"
        @cancel="handleDialogCancel"
        @confirm="handleDialogConfirm">
      </component>
    </cmdb-dialog>
  </div>
</template>

<script>
  import CreateServiceInstance from './children/create-service-instance.vue'
  import DeletedServiceInstance from './children/deleted-service-instance.vue'
  import MoveToIdleHost from './children/move-to-idle-host.vue'
  import ModuleSelector from '@/views/business-topology/host/module-selector.vue'
  import HostSelector from '@/views/business-topology/host/host-selector-new'
  import HostAttrsAutoApply from './children/host-attrs-auto-apply.vue'
  import {
    MENU_BUSINESS_TRANSFER_HOST
  } from '@/dictionary/menu-symbol'
  import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
  import { mapGetters } from 'vuex'
  import hostSearchService from '@/service/host/search'
  export default {
    components: {
      [CreateServiceInstance.name]: CreateServiceInstance,
      [DeletedServiceInstance.name]: DeletedServiceInstance,
      [MoveToIdleHost.name]: MoveToIdleHost,
      [ModuleSelector.name]: ModuleSelector,
      [HostSelector.name]: HostSelector,
      [HostAttrsAutoApply.name]: HostAttrsAutoApply
    },
    data() {
      return {
        hasScrollbar: false,
        hostInfo: [],
        tab: {
          active: null
        },
        dialog: {
          width: 830,
          height: 600,
          show: false,
          component: null,
          props: {}
        },
        tabList: [{
          id: 'createServiceInstance',
          label: this.$t('新增服务实例'),
          confirmed: false,
          showCount: false,
          component: CreateServiceInstance.name,
          props: {
            info: []
          }
        }, {
          id: 'deletedServiceInstance',
          label: this.$t('删除服务实例'),
          confirmed: false,
          showCount: true,
          component: DeletedServiceInstance.name,
          props: {
            info: []
          }
        }, {
          id: 'moveToIdleHost',
          label: this.$t('移动到空闲机的主机', { idleModule: this.$store.state.globalConfig.config.idlePool.idle }),
          confirmed: false,
          showCount: true,
          component: MoveToIdleHost.name,
          props: {
            info: []
          }
        }, {
          id: 'hostAttrsAutoApply',
          label: this.$t('属性自动应用'),
          confirmed: false,
          showCount: true,
          component: HostAttrsAutoApply.name,
          props: {
            info: []
          }
        }],
        request: {
          preview: Symbol('review'),
          module: Symbol('module'),
          confirm: Symbol('confirm'),
          mainline: Symbol('mainline'),
          host: Symbol('host'),
          internal: Symbol('internal')
        },
        targetModules: [],
        resources: [],
        type: this.$route.params.type,
        confirmParams: {},
        moduleMap: {},
        isSameModule: false,
        isEmptyChange: false,
        loading: true
      }
    },
    provide() {
      return {
        getModuleName: this.getModuleName,
        getModulePath: this.getModulePath
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId', 'currentBusiness']),
      ...mapGetters('businessHost', [
        'getDefaultSearchCondition'
      ]),
      confirmText() {
        const textMap = {
          remove: this.$t('确认移除'),
          idle: this.$t('确认转移'),
          business: this.$t('确认转移'),
          increment: this.$t('确认追加'),
          add: this.$t('确认添加'),
        }
        return textMap[this.type]
      },
      availableTabList() {
        const map = {
          remove: ['deletedServiceInstance', 'moveToIdleHost', 'hostAttrsAutoApply'],
          idle: ['deletedServiceInstance', 'hostAttrsAutoApply'],
          business: ['createServiceInstance', 'deletedServiceInstance', 'hostAttrsAutoApply'],
          increment: ['createServiceInstance', 'hostAttrsAutoApply'],
          add: ['createServiceInstance', 'hostAttrsAutoApply']
        }
        const available = map[this.type]
        return this.tabList.filter(tab => available.includes(tab.id) && tab.props.info.length > 0)
      },
      activeTab() {
        return this.tabList.find(tab => tab.id === this.tab.active) || this.availableTabList[0]
      },
      isRemoveModule() {
        const { type, module } = this.$route.params
        return type === 'remove' && module
      },
      isSingle() {
        return parseInt(this.$route.query.single, 10) === 1
      },
      changeHostDisabled() {
        return this.isSingle || this.isRemoveModule
      }
    },
    watch: {
      availableTabList(tabList) {
        tabList.forEach((tab) => {
          if (tab !== this.activeTab) {
            tab.confirmed = false
          }
        })
        const hasActiveTab = tabList.find(tab => tab === this.activeTab)
        if (!hasActiveTab) this.tab.active = null
      },
      activeTab(tab) {
        if (!tab) return
        tab.confirmed = true
      }
    },
    async created() {
      this.resolveData(this.$route)
      this.setBreadcrumbs()
      await Promise.all([
        this.getTopologyModels(),
        this.getHostInfo()
      ])
      this.getPreviewData()
    },
    mounted() {
      addResizeListener(this.$refs.changeInfo, this.resizeHandler)
    },
    beforeDestroy() {
      removeResizeListener(this.$refs.changeInfo, this.resizeHandler)
    },
    async beforeRouteUpdate(to, from, next) {
      this.resolveData(to)
      await this.getHostInfo()
      this.$nextTick(this.setBreadcrumbs)
      this.getPreviewData()
      next()
    },
    methods: {
      resolveData(route) {
        this.type = route.params.type
        const query = route.query || {}
        const { targetModules } = query
        if (!targetModules) {
          this.targetModules = []
        } else {
          this.targetModules = String(targetModules).split(',')
            .map(id => Number(id))
        }

        const { resources } = query
        if (!resources) {
          this.resources = []
        } else {
          this.resources = String(resources).split(',')
            .map(id => Number(id))
        }

        const isTransfer = ['idle', 'business'].includes(this.type)

        const params = {
          bk_host_ids: this.resources,
          // 是否从现有模块移除（清除掉主机现有的所有模块）
          is_remove_from_all: isTransfer
        }

        if (this.type === 'idle') {
          // eslint-disable-next-line prefer-destructuring
          params.default_internal_module = this.targetModules[0]
        } else if (this.type === 'remove') {
          params.remove_from_modules = [Number(query.sourceId)]
        } else if (this.targetModules.length) {
          params.add_to_modules = this.targetModules
        }
        this.confirmParams = params
      },
      setBreadcrumbs() {
        const titleMap = {
          idle: this.$t('转移到空闲模块', { idleSet: this.$store.state.globalConfig.config.set }),
          business: this.$t('转移到业务模块'),
          remove: this.$t('移除主机'),
          increment: this.$t('追加主机'),
          add: this.$t('添加主机')
        }
        this.$store.commit('setTitle', titleMap[this.type])
      },
      async getTopologyModels() {
        try {
          const topologyModels = await this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
            config: {
              requestId: this.request.mainline
            }
          })
          this.$store.commit('businessHost/setTopologyModels', topologyModels)
        } catch (e) {
          console.error(e)
        }
      },
      async getPreviewData() {
        try {
          this.loading = true
          const data = await this.$http.post(
            `host/transfer_with_auto_clear_service_instance/bk_biz_id/${this.bizId}/preview`,
            this.confirmParams,
            {
              requestId: this.request.preview,
              globalPermission: false
            }
          )
          this.setConfirmState(data)
          this.setModulePathInfo(data)
          this.setHostAttrsAutoApply(data)
          this.setCreateServiceInstance(data)
          this.setDeletedServiceInstance(data)
          if (this.type === 'remove') {
            this.setMoveToIdleHost(data)
          }
          this.loading = false
        } catch (e) {
          console.error(e)
          this.loading = false
          if (e.code === 9900403) {
            this.$route.meta.view = 'permission'
            this.$route.meta.auth.permission = e.permission
          }
        }
      },
      setConfirmState(data) {
        // 是否是相同的模块转换
        // eslint-disable-next-line max-len
        this.isSameModule = data.every(datum => !(datum.to_add_to_modules.length || datum.to_remove_from_modules.length))
        // 是否溢出的是空服务实例（前端流转不会创建空服务实例，但是ESB会）
        this.isEmptyChange = data.every((datum) => {
          const hasAdd = !datum.to_add_to_modules.length
          const hasRemoveInstance = datum.to_remove_from_modules.some(module => !module.service_instances.length)
          return !(hasAdd || hasRemoveInstance)
        })
      },
      async setModulePathInfo(data) {
        try {
          const moduleIds = [...this.targetModules]
          data.forEach((datum) => {
            moduleIds.push(...datum.to_add_to_modules.map(datum => datum.bk_module_id))
            moduleIds.push(...datum.to_remove_from_modules.map(datum => datum.bk_module_id))
          })
          const uniqueModules = [...new Set(moduleIds)]
          const result = await this.$store.dispatch('objectMainLineModule/getTopoPath', {
            bizId: this.bizId,
            params: {
              topo_nodes: uniqueModules.map(id => ({ bk_obj_id: 'module', bk_inst_id: id }))
            }
          })
          const moduleMap = {}
          result.nodes.forEach((node) => {
            moduleMap[node.topo_node.bk_inst_id] = node.topo_path
          })
          this.moduleMap = Object.freeze(moduleMap)
        } catch (e) {
          console.error(e)
        }
      },
      getModulePath(id) {
        const info = this.moduleMap[id] || []
        const path = info.map(node => node.bk_inst_name).reverse()
          .join(' / ')
        return path
      },
      setHostAttrsAutoApply(data) {
        const conflictInfo = (data || []).map(item => item.host_apply_plan)
        const conflictList = conflictInfo.filter(item => item.conflicts?.length || item.update_fields?.length)
        const tab = this.tabList.find(tab => tab.id === 'hostAttrsAutoApply')
        tab.props.info = Object.freeze(conflictList)
      },
      setCreateServiceInstance(data) {
        const instanceInfo = []
        data.forEach((item) => {
          item.to_add_to_modules.forEach((moduleInfo) => {
            instanceInfo.push({
              bk_host_id: item.bk_host_id,
              ...moduleInfo
            })
          })
        })
        const tab = this.tabList.find(tab => tab.id === 'createServiceInstance')
        tab.props.info = Object.freeze(instanceInfo)
      },
      async getHostInfo() {
        try {
          const result = await hostSearchService.getBizHosts({
            params: this.getSearchHostParams(),
            config: {
              requestId: this.request.host
            }
          })
          this.hostInfo = result.info
        } catch (e) {
          console.error(e)
        }
      },
      getSearchHostParams() {
        const params = {
          bk_biz_id: this.bizId,
          ip: { data: [], exact: 0, flag: 'bk_host_innerip|bk_host_outerip' },
          page: {},
          condition: this.getDefaultSearchCondition()
        }
        const hostCondition = params.condition.find(target => target.bk_obj_id === 'host')
        hostCondition.condition.push({
          field: 'bk_host_id',
          operator: '$in',
          value: this.resources
        })
        return params
      },
      setDeletedServiceInstance(data) {
        const deletedServiceInstance = []
        data.forEach((item) => {
          item.to_remove_from_modules.forEach((module) => {
            deletedServiceInstance.push(...module.service_instances)
          })
        })
        const tab = this.tabList.find(tab => tab.id === 'deletedServiceInstance')
        tab.props.info = Object.freeze(deletedServiceInstance)
      },
      getModuleName(id) {
        const topoInfo = this.moduleMap[id] || []
        const target = topoInfo.find(target => target.bk_obj_id === 'module' && target.bk_inst_id === id) || {}
        return target.bk_inst_name
      },
      async setMoveToIdleHost(data) {
        try {
          const internalTopology = await this.getInternalTopology()
          const internalModules = internalTopology.module
          const idleModule = internalModules.find(module => module.default === 1)
          const idleHost = []
          data.forEach((item) => {
            const finalModules = item.final_modules
            const isIdleModule = finalModules.length && finalModules[0] === idleModule.bk_module_id
            if (isIdleModule) {
              idleHost.push(item.bk_host_id)
            }
          })
          const tab = this.tabList.find(tab => tab.id === 'moveToIdleHost')
          tab.props.info = Object.freeze(idleHost)
        } catch (e) {
          console.error(e)
        }
      },
      getInternalTopology() {
        return this.$store.dispatch('objectMainLineModule/getInternalTopo', {
          bizId: this.bizId,
          config: {
            requestId: this.request.internal
          }
        })
      },
      handleRemoveModule(id) {
        const targetModules = this.targetModules.filter(exist => exist !== id)
        this.$routerActions.redirect({
          name: MENU_BUSINESS_TRANSFER_HOST,
          params: {
            type: 'business'
          },
          query: {
            ...this.$route.query,
            targetModules: targetModules.join(',')
          }
        })
      },
      resizeHandler() {
        this.$nextTick(() => {
          const scroller = this.$el.parentElement
          this.hasScrollbar = scroller.scrollHeight > scroller.offsetHeight
        })
      },
      handleTabClick(tab) {
        this.tab.active = tab.id
      },
      handleChangeModule() {
        const props = {
          moduleType: this.type,
          title: this.type === 'idle' ? this.$t('转移主机到空闲模块', { idleSet: this.$store.state.globalConfig.config.set }) : this.$t('转移主机到业务模块'),
          defaultChecked: this.targetModules,
          business: this.currentBusiness
        }
        const selection = this.hostInfo
        const firstSelectionModules = selection[0].module.map(module => module.bk_module_id).sort()
        const firstSelectionModulesStr = firstSelectionModules.join(',')
        const allSame = selection.slice(1).every((item) => {
          const modules = item.module.map(module => module.bk_module_id).sort()
            .join(',')
          return modules === firstSelectionModulesStr
        })
        if (allSame) {
          props.previousModules = firstSelectionModules
        }
        this.dialog.props = props
        this.dialog.width = 830
        this.dialog.height = 600
        this.dialog.component = ModuleSelector.name
        this.dialog.show = true
      },
      handleChangeHost() {
        const props = {
          exist: [...this.hostInfo]
        }
        if (this.type === 'remove') {
          props.displayNodes = [`${this.$route.query.sourceModel}-${this.$route.query.sourceId}`]
        }
        this.dialog.props = props
        this.dialog.width = 1280
        this.dialog.height = 750
        this.dialog.component = HostSelector.name
        this.dialog.show = true
      },
      handleDialogCancel() {
        this.dialog.show = false
      },
      handleDialogConfirm() {
        if (this.dialog.component === ModuleSelector.name) {
          // eslint-disable-next-line prefer-rest-params
          this.gotoTransferPage(...arguments)
          this.dialog.show = false
        } else if (this.dialog.component === HostSelector.name) {
          // eslint-disable-next-line prefer-rest-params
          this.refreshRemoveHost(...arguments)
          this.dialog.show = false
        }
      },
      refreshRemoveHost(hosts) {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_TRANSFER_HOST,
          params: {
            type: this.$route.params.type
          },
          query: {
            ...this.$route.query,
            resources: hosts.map(data => data.host.bk_host_id).join(',')
          }
        })
      },
      gotoTransferPage(modules) {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_TRANSFER_HOST,
          params: {
            type: this.dialog.props.moduleType
          },
          query: {
            ...this.$route.query,
            targetModules: modules.map(node => node.data.bk_inst_id).join(',')
          }
        })
      },
      refreshRetry(hosts) {
        this.$router.replace({
          name: MENU_BUSINESS_TRANSFER_HOST,
          params: {
            type: this.$route.params.type
          },
          query: {
            ...this.$route.query,
            resources: hosts.map(data => data.bk_host_id).join(','),
            retry: '1'
          }
        })
      },
      async handleConfrim() {
        try {
          const params = { ...this.confirmParams }
          const createComponent = this.$refs.createServiceInstance && this.$refs.createServiceInstance[0]
          const hostAttrsComponent = this.$refs.hostAttrsAutoApply && this.$refs.hostAttrsAutoApply[0]
          if (createComponent || hostAttrsComponent) {
            params.options = {}
            if (createComponent) {
              params.options.service_instance_options = createComponent.getServiceInstanceOptions()
            }
            if (hostAttrsComponent) {
              params.options.host_apply_trans_rule = hostAttrsComponent.getHostApplyConflictResolvers()
            }
          }

          await this.$http.post(`host/transfer_with_auto_clear_service_instance/bk_biz_id/${this.bizId}`, params, {
            requestId: this.request.confirm
          })

          const successText = ({ remove: '移除成功', add: '添加成功' })[this.type] || '转移成功'
          this.$success(this.$t(successText))

          this.redirect()
        } catch (e) {
          console.error(e)
        }
      },
      handleCancel() {
        this.redirect()
      },
      redirect() {
        this.$routerActions.back()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .layout {
        padding: 15px 0 0 0;

        .fail-detail-link {
            vertical-align: unset;
        }
    }
    .info {
        .info-label {
            width: 128px;
            font-size: 14px;
            font-weight: bold;
            color: $textColor;
            text-align: right;
            padding-top: 8px;
        }
        .info-content {
            overflow: hidden;
            padding: 8px 20px 0 8px;
            font-size: 14px;
            .info-count {
                font-weight: bold;
            }
            .module-grep {
                border-top: 1px solid $borderColor;
                margin-top: 10px;
            }
            .edit-trigger {
                @include inlineBlock;
                margin-left: 10px;
                color: $primaryColor;
                cursor: pointer;
                &:hover {
                    color: #1964E1;
                }
            }
        }
    }
    .module-list {
        font-size: 0;
        .module-item {
            position: relative;
            display: inline-block;
            vertical-align: middle;
            height: 26px;
            max-width: 150px;
            line-height: 24px;
            padding: 0 15px;
            margin: 0 10px 8px 0;
            border: 1px solid #C4C6CC;
            border-radius: 13px;
            color: $textColor;
            font-size: 12px;
            outline: none;
            cursor: default;
            @include ellipsis;
            &.is-business-module {
                padding: 0 12px 0 25px;
            }
            &.is-trigger {
                padding: 0;
                text-align: center;
                font-size: 0;
                cursor: pointer;
                border-color: transparent;
                color: $primaryColor;
                &:hover {
                    color: #1964E1;
                    border-color: transparent;
                }
                .icon-cc-edit {
                    font-size: 14px;
                }
            }
            &:hover {
                border-color: $primaryColor;
                color: $primaryColor;
                .module-mask {
                    display: block;
                }
                .module-icon {
                    background-color: $primaryColor;
                }
            }
            .module-mask {
                display: none;
                position: absolute;
                left: 0;
                top: 0;
                width: 100%;
                height: 100%;
                color: #fff;
                background-color: rgba(0, 0, 0, 0.53);
                text-align: center;
                cursor: pointer;
            }
            .module-icon {
                position: absolute;
                left: 2px;
                top: 2px;
                width: 20px;
                height: 20px;
                border-radius: 50%;
                line-height: 20px;
                text-align: center;
                color: #FFF;
                font-size: 12px;
                background-color: #C4C6CC;
            }
            .module-remove {
                position: absolute;
                right: 4px;
                top: 4px;
                width: 16px;
                height: 16px;
                border-radius: 50%;
                text-align: center;
                line-height: 16px;
                color: #FFF;
                font-size: 0px;
                background-color: #C4C6CC;
                cursor: pointer;
                &:before {
                    display: inline-block;
                    vertical-align: middle;
                    font-size: 20px;
                    transform: translateX(-2px) scale(.5);
                }
            }
        }
    }
    .tab {
        .tab-grep {
            width: 2px;
            height: 19px;
            margin: 0 15px;
            background-color: #C4C6CC;
        }
        .tab-item {
            position: relative;
            color: $textColor;
            font-size: 0;
            cursor: pointer;
            &.active {
                color: $primaryColor;
            }
            &.active:after {
                content: "";
                position: absolute;
                left: 0;
                top: 30px;
                width: 100%;
                height: 2px;
                background-color: $primaryColor;
            }
            &.hide-count {
              .tab-label {
                margin-right: 10px;
              }
            }
            .tab-label {
                display: inline-block;
                vertical-align: middle;
                margin-left: 10px;
                margin-right: 4px;
                font-size: 14px;
            }
            .tab-count {
                display: inline-block;
                vertical-align: middle;
                height: 16px;
                padding: 0 5px;
                border-radius: 8px;
                line-height: 14px;
                font-size: 12px;
                color: #FFF;
                background-color: #C4C6CC;
                text-align: center;
                border: 1px solid #fff;

                &.unconfirmed {
                    background-color: #FF5656;
                }
            }
        }
    }
    .tab-component {
        margin-top: 20px;
    }
    .tab-empty {
        height: 60px;
        padding: 0 28px;
        line-height: 60px;
        background-color: #F0F1F5;
        color: $textColor;
        &:before {
            content: "!";
            display: inline-block;
            width: 16px;
            height: 16px;
            line-height: 16px;
            border-radius: 50%;
            text-align: center;
            color: #FFF;
            font-size: 12px;
            background-color: #C4C6CC;
        }
    }
    .options {
        position: sticky;
        padding: 10px 0 10px 136px;
        font-size: 0;
        bottom: 0;
        left: 0;
        &.is-sticky {
            background-color: #FFF;
            border-top: 1px solid $borderColor;
            z-index: 100;
        }
    }
</style>
