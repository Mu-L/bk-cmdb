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
      <bk-date-picker class="options-date"
        type="datetimerange"
        transfer
        :placeholder="$t('请选择xx', { name: $t('时间范围') })"
        v-model="timeRange"
        @change="getHistory">
      </bk-date-picker>
    </div>
    <div class="history-table" v-bkloading="{ isLoading: $loading(requestId) }">
      <bk-table ref="table"
        :data="histories"
        :max-height="$APP.height - 180"
        :row-style="{ cursor: 'pointer' }"
        @cell-click="handleView">
        <bk-table-column type="expand" width="30" align="center">
          <task-details-history-content
            slot-scope="{ row }"
            :details="row.bk_detail">
          </task-details-history-content>
        </bk-table-column>
        <bk-table-column :label="$t('操作概要')" prop="bk_summary" width="200" show-overflow-tooltip>
          <i18n slot-scope="{ row }" path="新增N台主机，更新M台主机">
            <template #new><span class="summary-count">{{getCount(row, 'new_add')}}</span></template>
            <template #update><span class="summary-count">{{getCount(row, 'update')}}</span></template>
          </i18n>
        </bk-table-column>
        <bk-table-column :label="$t('状态')" prop="bk_sync_status">
          <div class="row-status" slot-scope="{ row }">
            <i :class="['status', { 'is-error': row.bk_sync_status !== 'cloud_sync_success' }]"></i>
            {{row.bk_sync_status === 'cloud_sync_success' ? $t('成功') : $t('失败')}}
          </div>
        </bk-table-column>
        <bk-table-column :label="$t('时间')" prop="create_time">
          <template slot-scope="{ row }">{{row.create_time | formatter('time')}}</template>
        </bk-table-column>
      </bk-table>
    </div>
  </div>
</template>

<script>
  import TaskDetailsHistoryContent from './task-details-history-content.vue'
  export default {
    name: 'task-details-history',
    components: {
      [TaskDetailsHistoryContent.name]: TaskDetailsHistoryContent
    },
    props: {
      id: Number
    },
    data() {
      return {
        timeRange: [],
        histories: [],
        pagination: this.$tools.getDefaultPaginationConfig(),
        requestId: Symbol('getHistory')
      }
    },
    created() {
      this.getHistory()
    },
    methods: {
      async getHistory() {
        try {
          const params = {
            bk_task_id: this.id,
            page: this.$tools.getPageParams(this.pagination)
          }
          if (this.timeRange.length) {
            params.start_time = this.$tools.formatTime(this.timeRange[0])
            params.end_time = this.$tools.formatTime(this.timeRange[1])
          }
          const data = await this.$store.dispatch('cloud/resource/findHistory', {
            params,
            config: {
              requestId: this.requestId
            }
          })
          this.pagination.count = data.count
          this.histories = data.info
        } catch (e) {
          console.error(e)
          this.histories = []
        }
      },
      getCount(row, type) {
        return (row.bk_detail[type] || {}).count || 0
      },
      handleView(row, column) {
        column.type !== 'expand' && this.$refs.table.toggleRowExpansion(row)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .history-layout {
        .history-options {
            padding: 12px 24px;
            .options-date {
                width: 320px;
            }
        }
        .history-table {
            padding: 0 24px;
            /deep/ {
                .bk-table-expanded-cell {
                    padding: 0;
                }
            }
        }
        .summary-count {
            font-weight: bold;
            padding: 0 2px;
        }
        .row-status {
            display: inline-block;
            .status {
                display: inline-block;
                margin-right: 4px;
                width: 7px;
                height: 7px;
                border-radius: 50%;
                background-color: $successColor;
            }
        }
    }
</style>
