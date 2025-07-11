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
  <div class="dynamic-group-layout">
    <cmdb-tips class="mb10" tips-key="showCustomQuery"
      :more-link="`${$helpDocUrlPrefix}/UserGuide/Feature/CustomQuery.md`">
      {{$t('动态分组提示')}}
    </cmdb-tips>
    <div class="dynamic-group-options">
      <cmdb-auth class="options-create"
        :auth="{ type: $OPERATION.C_CUSTOM_QUERY, relation: [bizId] }">
        <bk-button slot-scope="{ disabled }"
          theme="primary"
          :disabled="disabled"
          @click="handleCreate">
          {{$t('新建')}}
        </bk-button>
      </cmdb-auth>
      <search-select
        :default-filter="filter"
        @search="handleSearch">
      </search-select>
    </div>
    <div class="dynamic-group-table">
      <bk-table
        class="api-table"
        v-bkloading="{ isLoading: $loading([request.search, request.delete]) }"
        :data="table.list"
        :pagination="table.pagination"
        :max-height="$APP.height - 229"
        @cell-click="handleCellClick"
        @page-change="handlePageChange"
        @page-limit-change="handlePageLimitChange"
        @sort-change="handleSortChange">
        <bk-table-column prop="name" :label="$t('查询名称')" sortable="custom" fixed show-overflow-tooltip>
          <span class="name-text" slot-scope="{ row }">{{row.name}}</span>
        </bk-table-column>
        <bk-table-column prop="id" label="ID" show-overflow-tooltip></bk-table-column>
        <bk-table-column prop="bk_obj_id" :label="$t('查询对象')" sortable="custom">
          <span slot-scope="{ row }">{{getModelName(row)}}</span>
        </bk-table-column>
        <bk-table-column prop="create_user" :label="$t('创建用户')" sortable="custom" show-overflow-tooltip>
          <template slot-scope="{ row }">
            <cmdb-form-objuser :value="row.create_user" type="info"></cmdb-form-objuser>
          </template>
        </bk-table-column>
        <bk-table-column prop="create_time" :label="$t('创建时间')" sortable="custom">
          <template slot-scope="{ row }">
            {{row.create_time | formatter('time')}}
          </template>
        </bk-table-column>
        <bk-table-column prop="modify_user" :label="$t('修改人')" sortable="custom" show-overflow-tooltip>
          <template slot-scope="{ row }">
            <cmdb-form-objuser :value="row.modify_user" type="info"></cmdb-form-objuser>
          </template>
        </bk-table-column>
        <bk-table-column prop="last_time" :label="$t('修改时间')" sortable="custom">
          <template slot-scope="{ row }">
            {{row.last_time | formatter('time')}}
          </template>
        </bk-table-column>
        <bk-table-column prop="operation" :label="$t('操作')" fixed="right" :width="$i18n.locale === 'en' ? 170 : 130">
          <template slot-scope="{ row }">
            <bk-button class="mr10"
              :text="true"
              @click.stop="handlePreview(row)">
              {{$t('预览')}}
            </bk-button>
            <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_CUSTOM_QUERY, relation: [bizId, row.id] }">
              <bk-button slot-scope="{ disabled }"
                :disabled="disabled"
                :text="true"
                @click.stop="handleEdit(row)">
                {{$t('编辑')}}
              </bk-button>
            </cmdb-auth>
            <cmdb-auth :auth="{ type: $OPERATION.D_CUSTOM_QUERY, relation: [bizId, row.id] }">
              <bk-button slot-scope="{ disabled }"
                :disabled="disabled"
                :text="true"
                @click.stop="handleDelete(row)">
                {{$t('删除')}}
              </bk-button>
            </cmdb-auth>
          </template>
        </bk-table-column>
        <cmdb-table-empty
          slot="empty"
          :stuff="table.stuff"
          :auth="{ type: $OPERATION.C_CUSTOM_QUERY, relation: [bizId] }"
          @create="() => handleCreate('create')"
          @clear="handleClearFilter">
        </cmdb-table-empty>
      </bk-table>
    </div>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import RouterQuery from '@/router/query'
  import DynamicGroupForm from './form/form.js'
  import SearchSelect from './children/search-select.vue'

  export default {
    components: {
      SearchSelect
    },
    data() {
      return {
        filter: [],
        table: {
          list: [],
          sort: '-last_time',
          pagination: this.$tools.getDefaultPaginationConfig(),
          stuff: {
            type: 'default',
            payload: {
              resource: this.$t('动态分组')
            }
          }
        },
        request: {
          search: Symbol('search'),
          delete: Symbol('delete')
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('objectModelClassify', ['getModelById']),
      searchParams() {
        const params = {
          condition: {},
          page: {
            ...this.$tools.getPageParams(this.table.pagination),
            sort: this.table.sort
          }
        }
        this.filter.forEach((item) => {
          const itemValue = item.value?.split(',')
          const value = itemValue?.length > 1 ? {
            $in: itemValue
          } : itemValue[0]
          params.condition[item?.id] = value
        })

        return params
      },
    },
    created() {
      this.unwatchQuery = RouterQuery.watch('*', ({ page, limit, sort, action, id, name, bk_obj_id: objId, modify_user: modifyUser  }) => {
        this.table.pagination.current = parseInt(page || this.table.pagination.current, 10)
        this.table.pagination.limit = parseInt(limit || this.table.pagination.limit, 10)
        this.table.sort = sort || this.table.sort
        const queryFilter = [
          { id: 'id', value: id },
          { id: 'name', value: name },
          { id: 'bk_obj_id', value: objId },
          { id: 'modify_user', value: modifyUser }
        ]
        this.filter = queryFilter.filter(item => item.value?.length)
        if (action === 'create') {
          this.handleCreate()
        }
        this.getList()
      }, { immediate: true })
    },
    beforeDestroy() {
      this.unwatchQuery && this.unwatchQuery()
    },
    methods: {
      handleSearch(filter) {
        const query = {
          name: '',
          id: '',
          modify_user: '',
          bk_obj_id: '',
          _t: Date.now(),
          page: 1
        }
        // 处理 如果是bk_obj_id 将路由上的中文name 换成 id
        filter.forEach((item) => {
          const key = item.id === 'bk_obj_id' ? 'id' : 'name'
          query[item.id] = item.values.map(val => val[key]).join(',')
        })
        RouterQuery.set(query)
      },
      async getList() {
        try {
          const params = this.searchParams
          const { info, count } = await this.$store.dispatch('dynamicGroup/search', {
            bizId: this.bizId,
            params,
            config: {
              requestId: this.request.search,
              cancelPrevious: true
            }
          })
          this.table.list = info
          this.table.pagination.count = count
          this.table.stuff.type = this.filter[0] ? 'search' : 'default'
        } catch (error) {
          console.error(error)
          if (error.permission) {
            this.table.stuff.type = {
              type: 'permission',
              payload: { permission: error.permission }
            }
          }
        }
      },
      handleCreate(type = '') {
        if (type) {
          RouterQuery.set({ action: 'create' })
        }
        DynamicGroupForm.show({
          title: this.$t('新建动态分组')
        })
      },
      handleCellClick(row, column, cell, event) {
        if (column.property !== 'name') {
          return false
        }
        const clickTarget = event.target
        if (clickTarget.classList && clickTarget.classList.contains('name-text')) {
          this.handlePreview(row)
        }
      },
      getModelName(row) {
        const model = this.getModelById(row.bk_obj_id)
        return model ? model.bk_obj_name : row.bk_obj_id
      },
      handleEdit(row) {
        DynamicGroupForm.show({
          id: row.id,
          title: this.$t('编辑动态分组', { name: row.name })
        })
      },
      handleDelete(row) {
        this.$bkInfo({
          title: this.$t('确定删除'),
          subTitle: this.$t('确认要删除分组', { name: row.name }),
          extCls: 'bk-dialog-sub-header-center',
          confirmFn: async () => {
            await this.$store.dispatch('dynamicGroup/delete', {
              bizId: this.bizId,
              id: row.id,
              config: {
                requestId: this.request.delete
              }
            })
            this.$success(this.$t('删除成功'))
            const currentPage = this.table.pagination.current
            RouterQuery.set({
              page: this.table.list.length > 1 ? currentPage : ((currentPage - 1) || 1),
              _t: Date.now()
            })
          }
        })
      },
      handlePreview(row) {
        DynamicGroupForm.show({
          id: row.id,
          title: this.$t('动态分组详情', { name: row.name }),
          isPreview: true
        })
      },
      handlePageChange(page) {
        RouterQuery.set({
          page,
          _t: Date.now()
        })
      },
      handlePageLimitChange(limit) {
        RouterQuery.set({
          page: 1,
          limit,
          _t: Date.now()
        })
      },
      handleSortChange(sort) {
        RouterQuery.set({
          sort: this.$tools.getSort(sort, '-last_time'),
          _t: Date.now()
        })
      },
      handleClearFilter() {
        RouterQuery.clear()
      }
    }
  }
</script>

<style lang="scss" scoped>
.dynamic-group-layout {
  padding: 20px;
}
.dynamic-group-options {
  display: flex;
  align-items: center;
  justify-content: space-between;
  .options-filter {
    width: 320px;
  }
}
.dynamic-group-table {
  margin-top: 15px;
  .name-text {
    cursor: pointer;
    color: $primaryColor;
  }
}
</style>
