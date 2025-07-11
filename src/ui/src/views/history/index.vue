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
  <div class="history-layout">
    <div class="history-options">
      <cmdb-form-date-range class="history-date-range"
        v-model="condition.operation_time"
        :clearable="false"
        @change="handlePageChange(1)">
      </cmdb-form-date-range>
      <bk-input class="history-host-filter ml10"
        v-if="isHost"
        right-icon="icon-search"
        clearable
        v-model="condition.resource_name"
        :placeholder="$t('请输入xx', { name: 'IP' })"
        @enter="handlePageChange(1)"
        @clear="handlePageChange(1)"
        @right-icon-click="handlePageChange(1)">
      </bk-input>
    </div>
    <bk-table class="history-table"
      v-bkloading="{ isLoading: $loading() }"
      :pagination="pagination"
      :data="history"
      :max-height="$APP.height - 190"
      :row-style="{ cursor: 'pointer' }"
      @page-change="handlePageChange"
      @page-limit-change="handleSizeChange"
      @row-click="handleRowClick">
      <bk-table-column prop="resource_id" label="ID"></bk-table-column>
      <bk-table-column prop="resource_name" :label="isHost ? 'IP' : $t('资源')" show-overflow-tooltip></bk-table-column>
      <bk-table-column prop="operation_time" :label="$t('更新时间')">
        <template slot-scope="{ row }">{{$tools.formatTime(row.operation_time)}}</template>
      </bk-table-column>
      <bk-table-column prop="user" :label="$t('操作账号')">
        <template slot-scope="{ row }">
          <cmdb-form-objuser :value="row.user" type="info"></cmdb-form-objuser>
        </template>
      </bk-table-column>
      <cmdb-table-empty slot="empty" :stuff="table.stuff" @clear="handleClearFilter"></cmdb-table-empty>
    </bk-table>
  </div>
</template>

<script>
  import AuditDetails from '@/components/audit-history/details.js'
  import tools from '@/utils/tools'

  const today = tools.formatTime(new Date(), 'YYYY-MM-DD')
  const formatValue = () => ({
    operation_time: [today, today],
    resource_name: ''
  })

  export default {
    data() {
      return {
        dictionary: [],
        history: [],
        pagination: this.$tools.getDefaultPaginationConfig(),
        condition: {
          ...formatValue(),
          action: ['delete']
        },
        table: {
          stuff: {
            type: 'default',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            }
          }
        },
        requestId: Symbol('getHistory')
      }
    },
    computed: {
      objId() {
        if (this.$route.name === 'hostHistory') {
          return 'host'
        }
        return this.$route.params.objId
      },
      isHost() {
        return this.objId === 'host'
      },
      resourceType() {
        return this.isHost ? 'host' : 'model_instance'
      },
      bizId() {
        return this.isHost ? 1 : 0
      }
    },
    watch: {
      objId: {
        immediate: true,
        handler(objId) {
          const model = this.$store.getters['objectModelClassify/getModelById'](objId) || {}
          this.$store.commit('setTitle', `${model.bk_obj_name} ${this.$t('删除历史')}`)
        }
      }
    },
    created() {
      this.getAuditDictionary()
      this.getHistory()
    },
    methods: {
      async getAuditDictionary() {
        try {
          this.dictionary = await this.$store.dispatch('audit/getDictionary', {
            fromCache: true,
            globalPermission: false
          })
        } catch (error) {
          this.dictionary = []
        }
      },
      async getHistory() {
        try {
          const { info, count } = await this.$store.dispatch('audit/getInstList', {
            params: {
              condition: this.getUsefulConditon(),
              page: {
                ...this.$tools.getPageParams(this.pagination),
                sort: '-operation_time'
              }
            },
            config: {
              requestId: this.requestId,
              globalPermission: false
            }
          })
          this.table.stuff.type = this.condition.resource_name ? 'search' : 'default'
          this.pagination.count = count
          this.history = info
        } catch ({ permission }) {
          if (permission) {
            this.table.stuff = {
              type: 'permission',
              payload: { permission }
            }
          }
          this.history = []
        }
      },
      getUsefulConditon() {
        const usefuleCondition = {
          bk_obj_id: this.objId,
          bk_biz_id: this.bizId,
          resource_type: this.resourceType
        }
        Object.keys(this.condition).forEach((key) => {
          const value = this.condition[key]
          if (String(value).length) {
            usefuleCondition[key] = value
          }
        })
        if (usefuleCondition.operation_time) {
          const [start, end] = usefuleCondition.operation_time
          usefuleCondition.operation_time = {
            start: `${start} 00:00:00`,
            end: `${end} 23:59:59`
          }
        }
        return usefuleCondition
      },
      handleSizeChange(size) {
        this.pagination.limit = size
        this.handlePageChange(1)
      },
      handlePageChange(current) {
        this.pagination.current = current
        this.getHistory()
      },
      handleRowClick(row) {
        AuditDetails.show({
          aduitTarget: 'instance',
          id: row.id,
          resourceType: this.resourceType,
          bizId: this.bizId,
          objId: this.objId
        })
      },
      handleClearFilter() {
        this.condition = {
          ...formatValue(),
          action: ['delete']
        }
        this.getHistory()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .history-layout{
        padding: 15px 20px 0;
    }
    .history-options{
        font-size: 0px;
        .history-host-filter,
        .history-date-range {
            width: 260px !important;
            display: inline-block;
            vertical-align: top;
        }
    }
    .history-table{
        margin-top: 15px;
    }

</style>
