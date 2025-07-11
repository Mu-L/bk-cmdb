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
  <div class="service-wrapper">
    <div class="options clearfix">
      <template v-if="!readonly">
        <bk-checkbox class="options-checkall"
          :size="16"
          v-model="isCheckAll"
          :disabled="!instances.length"
          :title="$t('全选本页')"
          @change="handleCheckAll">
        </bk-checkbox>
        <cmdb-auth :auth="HOST_AUTH.D_SERVICE_INSTANCE">
          <bk-button slot-scope="{ disabled }"
            class="ml10"
            :disabled="disabled || !checked.length"
            @click="batchDelete(!checked.length)">
            {{$t('批量删除')}}
          </bk-button>
        </cmdb-auth>
      </template>
      <div class="option-right fr">
        <bk-checkbox class="options-checkbox"
          :size="16"
          :disabled="!instances.length"
          v-model="isExpandAll"
          @change="handleExpandAll">
          <span class="checkbox-label">{{$t('全部展开')}}</span>
        </bk-checkbox>
        <div class="options-search">
          <bk-search-select
            ref="searchSelect"
            :show-condition="false"
            :placeholder="$t('请输入实例名称或选择标签')"
            :data="searchSelect"
            v-model="searchSelectData"
            @menu-child-condition-select="handleConditionSelect"
            @change="handleSearch">
          </bk-search-select>
        </div>
      </div>
    </div>
    <div class="tables">
      <service-instance-table
        v-for="(instance, index) in instances"
        ref="serviceInstanceTable"
        :key="instance.id"
        :instance="instance"
        :expanded="index === 0"
        @delete-instance="handleDeleteInstance"
        @update-instance="handleUpdateInstance"
        @check-change="handleCheckChange"
        @edit-name="handleEditName(instance)"
        @edit-name-success="handleEditNameSuccess(instance, ...arguments)"
        @cancel-edit-name="handleCancelEditName(instance)">
      </service-instance-table>
    </div>
    <bk-table v-if="!instances.length" :data="[]" class="mb10">
      <cmdb-data-empty
        slot="empty"
        :stuff="emptyStuff"
        @create="handleGoAddInstance"
        @clear="handleFilterClear">
      </cmdb-data-empty>
    </bk-table>
    <bk-pagination v-if="instances.length"
      class="pagination"
      align="right"
      size="small"
      :current="pagination.current"
      :count="pagination.count"
      :limit="pagination.size"
      @change="handlePageChange"
      @limit-change="handleSizeChange">
    </bk-pagination>
  </div>
</template>

<script>
  import has from 'has'
  import {
    MENU_BUSINESS_HOST_AND_SERVICE
  } from '@/dictionary/menu-symbol'
  import { mapState } from 'vuex'
  import { t } from '@/i18n'
  import serviceInstanceTable from './service-instance-table.vue'
  import authMixin from '../mixin-auth'
  import { readonlyMixin } from '../mixin-readonly'
  import { historyLabelProxy, hostServiceInstancesProxy, topoPathProxy } from '../service-proxy'

  const defaultSearchSelect = () => ([
    {
      name: t('服务实例名'),
      id: 0
    },
    {
      name: t('标签值'),
      id: 1,
      children: [{
        id: '',
        name: ''
      }],
      conditions: []
    },
    {
      name: t('标签键'),
      id: 2,
      children: [{
        id: '',
        name: ''
      }]
    }
  ])

  export default {
    components: {
      serviceInstanceTable
    },
    mixins: [authMixin, readonlyMixin],
    data() {
      return {
        searchSelect: defaultSearchSelect(),
        searchSelectData: [],
        pagination: {
          current: 1,
          count: 0,
          size: 10
        },
        checked: [],
        isExpandAll: false,
        isCheckAll: false,
        filter: [],
        instances: [],
        currentView: 'path',
        historyLabels: {},
        emptyStuff: {
          type: 'default',
          payload: {
            path: '无服务实例提示',
            action: this.$t('去业务拓扑添加')
          }
        }
      }
    },
    computed: {
      ...mapState('hostDetails', ['info']),
      host() {
        return this.info.host || {}
      }
    },
    watch: {
      checked() {
        this.isCheckAll = (this.checked.length === this.instances.length) && this.checked.length !== 0
      }
    },
    created() {
      this.getHostSeriveInstances()
      this.getHistoryLabel()
    },
    methods: {
      async getHostSeriveInstances() {
        try {
          const searchKey = this.searchSelectData.find(item => (item.id === 0 && has(item, 'values'))
            || (![0, 1].includes(item.id) && !has(item, 'values')))
          const data = await hostServiceInstancesProxy({
            page: {
              start: (this.pagination.current - 1) * this.pagination.size,
              limit: this.pagination.size
            },
            bk_host_id: this.host.bk_host_id,
            bk_biz_id: this.info.biz[0].bk_biz_id,
            // eslint-disable-next-line no-nested-ternary
            search_key: searchKey
              ? has(searchKey, 'values') ? searchKey.values[0].name : searchKey.name
              : '',
            selectors: this.getSelectorParams()
          })
          if (data.count && !data.info.length) {
            this.pagination.current -= 1
            this.getHostSeriveInstances()
            return
          }

          // 获取所属模块的拓扑路径
          if (data.info?.length) {
            const topopath = await topoPathProxy(this.info.biz[0].bk_biz_id, {
              topo_nodes: data.info.map(item => ({ bk_obj_id: 'module', bk_inst_id: item.bk_module_id }))
            })
            // 将拓扑路径加入到服务列表数据中
            data.info.forEach((item) => {
              const moduleNode = topopath.nodes.find(node => node.topo_node.bk_inst_id === item.bk_module_id)
              item.topo_path = moduleNode.topo_path
            })
          }

          this.checked = []
          this.isCheckAll = false
          this.isExpandAll = false
          this.pagination.count = data.count
          this.instances = data.info.map(instance => ({ ...instance, editing: { name: false } }))
        } catch (e) {
          console.error(e)
          this.instances = []
          this.pagination.count = 0
        } finally {
          this.emptyStuff.type = this.searchSelectData.length === 0 ? 'default' : 'search'
        }
      },
      getSelectorParams() {
        try {
          const labels = this.searchSelectData.filter(item => item.id === 1 && has(item, 'values'))
          const labelsKey = this.searchSelectData.filter(item => item.id === 2 && has(item, 'values'))
          const submitLabel = {}
          const submitLabelKey = {}
          labels.forEach((label) => {
            const conditionId = label.condition.id
            if (!submitLabel[conditionId]) {
              submitLabel[conditionId] = [label.values[0].id]
            } else {
              if (submitLabel[conditionId].indexOf(label.values[0].id) < 0) {
                submitLabel[conditionId].push(label.values[0].id)
              }
            }
          })
          labelsKey.forEach((label) => {
            // eslint-disable-next-line prefer-destructuring
            const { id } = label.values[0]
            if (!submitLabelKey[id]) {
              submitLabelKey[id] = id
            }
          })
          const selectors = Object.keys(submitLabel).map(key => ({
            key,
            operator: 'in',
            values: submitLabel[key]
          }))
          const selectorsKey = Object.keys(submitLabelKey).map(key => ({
            key,
            operator: 'exists',
            values: []
          }))
          return selectors.concat(selectorsKey)
        } catch (e) {
          console.error(e)
          return []
        }
      },
      async getHistoryLabel() {
        const historyLabels = await historyLabelProxy(
          {
            bk_biz_id: this.info.biz[0].bk_biz_id
          },
          {
            requestId: 'getHistoryLabel',
            cancelPrevious: true
          }
        )
        this.historyLabels = historyLabels
        const keys = Object.keys(historyLabels)
        const valueOption = keys.map(key => ({
          name: `${key} : `,
          id: key
        }))
        const keyOption = keys.map(key => ({
          name: key,
          id: key
        }))
        if (!valueOption.length) {
          this.$set(this.searchSelect[1], 'disabled', true)
        }
        if (!keyOption.length) {
          this.$set(this.searchSelect[2], 'disabled', true)
        }
        this.$set(this.searchSelect[1], 'conditions', valueOption)
        this.$set(this.searchSelect[2], 'children', keyOption)
      },
      handleDeleteInstance() {
        this.getHostSeriveInstances()
      },
      handleUpdateInstance() {
        this.getHostSeriveInstances()
      },
      handleEditName(instance) {
        this.instances.forEach(instance => (instance.editing.name = false))
        instance.editing.name = true
      },
      handleEditNameSuccess(instance, value) {
        instance.name = value
        instance.editing.name = false
      },
      handleCancelEditName(instance) {
        instance.editing.name = false
      },
      handleCheckAll(checked) {
        this.searchSelectData = []
        this.isCheckAll = checked
        this.$refs.serviceInstanceTable.forEach((table) => {
          table.checked = checked
        })
      },
      handleExpandAll(expanded) {
        this.searchSelectData = []
        this.isExpandAll = expanded
        this.$refs.serviceInstanceTable.forEach((table) => {
          table.localExpanded = expanded
        })
      },
      batchDelete(disabled) {
        if (disabled) {
          return false
        }
        this.$bkInfo({
          title: this.$t('确定删除N个实例', { count: this.checked.length }),
          confirmLoading: true,
          confirmFn: async () => {
            const serviceInstanceIds = this.checked.map(instance => instance.id)
            try {
              await this.$store.dispatch('serviceInstance/deleteServiceInstance', {
                config: {
                  data: {
                    service_instance_ids: serviceInstanceIds,
                    bk_biz_id: this.info.biz[0].bk_biz_id
                  }
                }
              })
              this.$success(this.$t('删除成功'))
              this.getHostSeriveInstances()
              return true
            } catch (e) {
              console.error(e)
              return false
            }
          }
        })
      },
      handleCheckChange(checked, instance) {
        if (checked) {
          this.checked.push(instance)
        } else {
          this.checked = this.checked.filter(target => target.id !== instance.id)
        }
      },
      handleClearFilter() {
        this.handleSearch()
      },
      handleSearch() {
        const instanceName = this.searchSelectData.filter(item => (item.id === 0 && has(item, 'values'))
          || (![0, 1].includes(item.id) && !has(item, 'values')))
        if (instanceName.length) {
          this.searchSelect[0].id === 0 && this.searchSelect.shift()
        } else {
          this.searchSelect[0].id !== 0 && this.searchSelect.unshift({
            name: this.$t('服务实例名'),
            id: 0
          })
        }
        if (instanceName.length >= 2) {
          this.searchSelectData.pop()
          this.$bkMessage({
            message: this.$t('服务实例名重复'),
            theme: 'warning'
          })
          return
        }

        this.handlePageChange(1)
      },
      handlePageChange(page) {
        this.pagination.current = page
        this.getHostSeriveInstances()
      },
      handleSizeChange(size) {
        this.pagination.current = 1
        this.pagination.size = size
        this.getHostSeriveInstances()
      },
      handleConditionSelect(cur, index) {
        const values = this.historyLabels[cur.id]
        const children = values.map(item => ({
          id: item,
          name: item
        }))
        const el = this.$refs.searchSelect
        el.curItem.children = children
        el.updateChildMenu(cur, index, false)
        el.showChildMenu(children)
      },
      handleGoAddInstance() {
        const [biz] = this.info.biz
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_AND_SERVICE,
          params: {
            bizId: biz.bk_biz_id
          },
          query: {
            node: `biz-${biz.bk_biz_id}`,
            ip: this.info.host.bk_host_innerip
          }
        })
      },
      handleFilterClear() {
        this.searchSelectData = []
        this.searchSelect = defaultSearchSelect()
        this.getHistoryLabel()
        this.handlePageChange(1)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .service-wrapper {
        height: calc(100% - 30px);
        padding: 14px 0 0 0;
    }
    .options-checkall {
        width: 36px;
        height: 32px;
        line-height: 30px;
        padding: 0 9px;
        text-align: center;
        border: 1px solid #f0f1f5;
        border-radius: 2px;
    }
    .options-checkbox {
        margin: 0 15px 0 0;
        .checkbox-label {
            padding: 0 0 0 4px;
        }
    }
    .options-search {
        @include inlineBlock;
        position: relative;
        min-width: 240px;
        max-width: 500px;
        z-index: 99;
    }
    .popover-main {
        display: flex;
        justify-content: space-between;
        align-items: center;
        .bk-icon {
            margin: 0 0 0 10px;
            cursor: pointer;
        }
    }
    .options-check-view {
        @include inlineBlock;
        height: 32px;
        line-height: 30px;
        font-size: 0;
        border: 1px solid #c4c6cc;
        padding: 0 10px;
        border-radius: 2px;
        .view-btn {
            color: #c4c6cc;
            font-size: 14px;
            height: 100%;
            line-height: 30px;
            cursor: pointer;
            &:hover, &.active {
                color: #3a84ff;
            }
        }
        .dividing-line {
            @include inlineBlock;
            width: 1px;
            height: 14px;
            background-color: #dcdee5;
        }
    }
    .tables {
        @include scrollbar-y;
        max-height: calc(100% - 80px);
        margin: 16px 0 10px;
    }
    .empty-text {
        text-align: center;
        p {
            font-size: 14px;
            margin: -10px 0 0;
        }
        span {
            color: #3a84ff;
            cursor: pointer;
        }
    }
</style>

<style lang="scss">
    .check-view-color-theme {
        padding: 10px !important;
        background-color: #699df4 !important;
        .tippy-arrow {
            border-bottom-color: #699df4 !important;
        }
    }
</style>
