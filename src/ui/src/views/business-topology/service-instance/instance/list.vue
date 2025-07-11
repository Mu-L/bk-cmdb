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
  <bk-table class="instance-table" ref="instanceTable" v-test-id.businessHostAndService="'svrInstList'"
    v-bkloading="{ isLoading: $loading(request.getList) }"
    :row-class-name="getRowClassName"
    :data="list"
    :pagination="pagination"
    :max-height="$APP.height - 250"
    @page-change="handlePageChange"
    @page-limit-change="handlePageLimitChange"
    @expand-change="handleExpandChange"
    @selection-change="handleSelectionChange"
    @row-click="handleRowClick">
    <bk-table-column type="selection" prop="id"></bk-table-column>
    <bk-table-column type="expand" prop="expand" width="15" :before-expand-change="beforeExpandChange">
      <div slot-scope="{ row }" v-bkloading="{ isLoading: row.pending }">
        <expand-list
          :service-instance="row"
          :list-request="processListRequest"
          @update-list="handleExpandResolved(row, ...arguments)"></expand-list>
      </div>
    </bk-table-column>
    <bk-table-column :label="$t('实例名称')" prop="name" width="340">
      <list-cell-name slot-scope="{ row }" :row="row"
        @edit="handleEditName(row)"
        @success="handleEditNameSuccess(row, ...arguments)"
        @cancel="handleCancelEditName(row)">
      </list-cell-name>
    </bk-table-column>
    <bk-table-column :label="$t('进程数量')" prop="process_count" width="240" :resizable="false">
      <list-cell-count slot-scope="{ row }" :row="row"
        @refresh-count="handleRefreshCount">
      </list-cell-count>
    </bk-table-column>
    <bk-table-column :label="$t('标签')" prop="tag" min-width="150">
      <list-cell-tag slot-scope="{ row }" :row="row"
        @update-labels="handleUpdateLabels(row, ...arguments)">
      </list-cell-tag>
    </bk-table-column>
    <bk-table-column :label="$t('操作')" :min-width="$i18n.locale === 'en' ? 200 : 150">
      <list-cell-operation slot-scope="{ row }" :row="row"
        @refresh-count="handleRefreshCount">
      </list-cell-operation>
    </bk-table-column>
    <cmdb-table-empty
      slot="empty"
      :stuff="table.stuff"
      @clear="handleClearFilter">
    </cmdb-table-empty>
  </bk-table>
</template>

<script>
  import { mapGetters } from 'vuex'
  import ListCellName from './list-cell-name'
  import ListCellCount from './list-cell-count'
  import ListCellTag from './list-cell-tag'
  import ListCellOperation from './list-cell-operation'
  import ExpandList from './expand-list'
  import RouterQuery from '@/router/query'
  import Bus from '../common/bus'
  import LabelBatchDialog from './dialog/label-batch-dialog.js'
  import has from 'has'
  import { rollReqByDataKey } from '@/service/utils'
  export default {
    components: {
      ListCellName,
      ListCellCount,
      ListCellTag,
      ListCellOperation,
      ExpandList
    },
    data() {
      return {
        list: [],
        selection: [],
        pagination: this.$tools.getDefaultPaginationConfig(),
        filters: [],
        request: {
          getList: Symbol('getList')
        },
        table: {
          stuff: {
            type: 'default',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            },
          }
        }
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('businessHost', ['selectedNode']),
      searchKey() {
        const nameFilter = this.filters.find(data => data.id === 'name')
        return nameFilter ? nameFilter.values[0].name : ''
      },
      searchTag() {
        const tagFilters = []
        this.filters.forEach((data) => {
          if (data.id === 'name') return
          if (has(data, 'condition')) {
            tagFilters.push({
              key: data.condition.id,
              operator: 'in',
              values: data.values.map(value => value.name)
            })
          } else {
            const [{ id }] = data.values
            tagFilters.push({
              key: id,
              operator: 'exists',
              values: []
            })
          }
        })
        return tagFilters
      }
    },
    watch: {
      selectedNode() {
        this.handlePageChange(1)
      }
    },
    created() {
      this.unwatch = RouterQuery.watch(['page', 'limit', '_t', 'view'], ({
        page = this.pagination.current,
        limit = this.pagination.limit,
        view = 'instance'
      }) => {
        if (view !== 'instance') {
          return false
        }
        this.pagination.page = parseInt(page, 10)
        this.pagination.limit = parseInt(limit, 10)
        this.getList()
      }, { immediate: true, throttle: true })
      Bus.$on('expand-all-change', this.handleExpandAllChange)
      Bus.$on('filter-change', this.handleFilterChange)
      Bus.$on('batch-edit-labels', this.handleBatchEditLabels)
    },
    beforeDestroy() {
      Bus.$off('expand-all-change', this.handleExpandAllChange)
      Bus.$off('filter-change', this.handleFilterChange)
      Bus.$off('batch-edit-labels', this.handleBatchEditLabels)
      this.unwatch()
    },
    methods: {
      handleFilterChange(filters) {
        this.filters = filters
        RouterQuery.set({
          node: this.selectedNode.id,
          page: 1,
          _t: Date.now()
        })
      },
      async getList() {
        try {
          const { count, info } = await this.$store.dispatch('serviceInstance/getModuleServiceInstances', {
            params: {
              bk_biz_id: this.bizId,
              bk_module_id: this.selectedNode.data.bk_inst_id,
              page: this.$tools.getPageParams(this.pagination),
              search_key: this.searchKey,
              selectors: this.searchTag,
              with_name: true
            },
            config: {
              requestId: this.request.getList,
              cancelPrevious: true
            }
          })
          this.list = info.map(data => ({ ...data, counting: true, pending: true, editing: { name: false } }))
          this.pagination.count = count
          this.table.stuff.type = this.filters.length === 0 ? 'default' : 'search'

          if (this.list?.length) {
            this.getProcessCounts()
          }
        } catch (error) {
          this.list = []
          this.pagination.count = 0
          console.error(error)
        }
      },
      async getProcessCounts() {
        const serviceInstIds = this.list.map(item => item.id)
        const processCounts = await rollReqByDataKey(`${window.API_HOST}count/service_instance/processes`, { ids: serviceInstIds }, { limit: 100 })
        this.list.forEach((serviceInst) => {
          const processCount = processCounts.find(item => item.id === serviceInst.id) || {}
          this.$set(serviceInst, 'counting', false)
          this.$set(serviceInst, 'process_count', processCount.count)
        })
      },
      processListRequest(reqParams, reqConfig) {
        return this.$store.dispatch('processInstance/getServiceInstanceProcesses', {
          params: reqParams,
          config: reqConfig
        })
      },
      getRowClassName({ row }) {
        const className = ['instance-table-row']
        if (!row.process_count) {
          className.push('disabled')
        }
        return className.join(' ')
      },
      handlePageChange(page) {
        this.pagination.current = page
        RouterQuery.set({
          node: this.selectedNode.id,
          page,
          _t: Date.now()
        })
      },
      handlePageLimitChange(limit) {
        this.pagination.limit = limit
        this.pagination.page = 1
        RouterQuery.set({
          limit,
          page: 1,
          _t: Date.now()
        })
      },
      beforeExpandChange({ row }) {
        return !!row.process_count
      },
      handleExpandAllChange(expanded) {
        this.list.forEach((row) => {
          row.process_count && this.$refs.instanceTable.toggleRowExpansion(row, expanded)
        })
      },
      async handleExpandChange(row, expandedRows) {
        row.pending = expandedRows.includes(row)
      },
      handleSelectionChange(selection) {
        this.selection = selection
        Bus.$emit('instance-selection-change', selection)
      },
      handleRowClick(row) {
        this.$refs.instanceTable.toggleRowExpansion(row)
      },
      handleExpandResolved(row, list) {
        row.pending = false
        this.handleRefreshCount(row, list.length)
        if (!row.process_count) {
          this.$refs.instanceTable.toggleRowExpansion(row, false)
        }
      },
      handleRefreshCount(row, newCount) {
        row.process_count = newCount
      },
      handleUpdateLabels(row, labels) {
        row.labels = labels
        Bus.$emit('update-labels')
      },
      handleBatchEditLabels() {
        LabelBatchDialog.show({
          serviceInstances: this.selection,
          updateCallback: (removedKeys, newLabels) => {
            this.selection.forEach((instance) => {
              instance.labels && removedKeys.forEach(key => delete instance.labels[key])
              instance.labels = Object.assign({}, instance.labels, newLabels)
            })
          }
        })
      },
      handleEditName(row) {
        this.list.forEach(row => (row.editing.name = false))
        row.editing.name = true
      },
      handleEditNameSuccess(row, value) {
        row.name = value
        row.editing.name = false
      },
      handleCancelEditName(row) {
        row.editing.name = false
      },
      handleClickAddHost() {
        RouterQuery.set({
          tab: 'hostList',
          _t: Date.now(),
          page: '',
          limit: ''
        })
      },
      async handleClearFilter() {
        this.filters = []
        await this.getList()
        this.table.stuff.type = 'default'
        Bus.$emit('filter-clear')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .instance-table {
        /deep/ {
            .instance-table-row {
                &:hover,
                &.expanded {
                    td {
                        background-color: #f0f1f5;
                    }
                    .tag-item {
                        background-color: #dcdee5;
                    }
                }
                &:hover {
                    .tag-edit {
                        visibility: visible;
                    }
                    .tag-empty {
                        display: none;
                    }
                    .instance-name:not(.disabled) {
                        .name-edit {
                            visibility: visible;
                        }
                    }
                }
                &.disabled {
                    .bk-table-expand-icon {
                        display: none;
                        cursor: not-allowed;
                    }
                }
            }
            .bk-table-expanded-cell {
                padding-left: 80px !important;
            }
        }
    }
    .empty-content {
      .text-btn {
        font-size: 14px;
        margin-left: 2px;
        height: auto;
      }
    }
</style>
