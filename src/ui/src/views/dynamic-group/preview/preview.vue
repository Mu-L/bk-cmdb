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
  <bk-dialog
    v-model="isShow"
    header-position="left"
    :width="880"
    :title="title"
    :show-footer="false"
    :draggable="false"
    :mask-close="true"
    @after-leave="handleHidden">
    <template slot="header">
      <div :title="title" class="title-text">{{ title }} </div>
    </template>
    <bk-table class="preview-table"
      ref="table"
      v-bkloading="{ isLoading: $loading() }"
      :pagination="table.pagination"
      :data="table.list"
      height="536"
      max-height="536"
      @page-change="handlePageChange"
      @page-limit-change="handlePageLimitChange"
      @sort-change="handleSortChange">
      <bk-table-column v-for="property in previewProperties"
        show-overflow-tooltip
        sortable="custom"
        :key="property.bk_property_id"
        :render-header="h => renderHeader(h, property)"
        :prop="property.bk_property_id">
        <template slot-scope="{ row }">
          <cmdb-property-value
            :value="row[property.bk_property_id]"
            :show-unit="false"
            :show-title="false"
            :property="property">
          </cmdb-property-value>
        </template>
      </bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :stuff="table.stuff">
      </cmdb-table-empty>
    </bk-table>
  </bk-dialog>
</template>

<script>
  import { mapGetters } from 'vuex'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants'

  export default {
    props: {
      id: {
        type: String
      }
    },
    data() {
      return {
        details: null,
        properties: [],
        previewProperties: [],
        table: {
          pagination: this.$tools.getDefaultPaginationConfig({
            limit: 10, 'limit-list': [10, 50, 100, 500]
          }, false),
          sort: '-create_time',
          list: [],
          stuff: {
            type: 'default',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            }
          }
        },
        isShow: false
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapGetters('objectBiz', ['bizId']),
      previewFields() {
        return this.previewProperties.map(property => property.bk_property_id)
      },
      title() {
        const title = [this.$t('预览分组')]
        this.details && title.push(this.details.name)
        return title.join(' - ')
      }
    },
    async created() {
      await this.getDetails()
      await this.getProperties()
      await this.setPreviewProperties()
      this.preview()
    },
    methods: {
      async getDetails() {
        try {
          const details = await this.$store.dispatch('dynamicGroup/details', {
            bizId: this.bizId,
            id: this.id,
            config: {}
          })
          this.details = Object.freeze(details)
        } catch (error) {
          console.error(error)
        }
      },
      async getProperties() {
        try {
          const properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
            params: {
              bk_supplier_account: this.supplierAccount,
              bk_biz_id: this.bizId,
              bk_obj_id: this.details.bk_obj_id
            },
            config: {
              requestId: `dynamic-group-preivew-${this.details.bk_obj_id}`,
              fromCache: true
            }
          })
          this.properties = Object.freeze(properties)
        } catch (error) {
          console.error(error)
        }
      },
      async setPreviewProperties() {
        try {
          const previewProperties = await this.$tools.getDefaultHeaderProperties(this.properties)

          if (this.details.bk_obj_id === BUILTIN_MODELS.HOST) {
            const innerIPv6 = this.properties.find(item => item.bk_property_id === 'bk_host_innerip_v6')
            if (innerIPv6) {
              previewProperties.splice(1, 0, innerIPv6)
            }
          }

          this.previewProperties = Object.freeze(previewProperties)
        } catch (error) {
          console.error(error)
        }
      },
      async preview() {
        try {
          const { count, info } = await this.$store.dispatch('dynamicGroup/preview', {
            bizId: this.bizId,
            id: this.id,
            params: {
              fields: this.previewFields,
              page: {
                ...this.$tools.getPageParams(this.table.pagination),
                sort: this.table.sort
              }
            }
          })
          this.table.pagination.count = count
          this.table.list = Object.freeze(info)
        } catch (error) {
          this.table.pagination.count = 0
          this.table.list = []
          console.error(error)
        }
      },
      handlePageChange(page) {
        this.table.pagination.current = page
        this.preview()
      },
      handlePageLimitChange(limit) {
        this.table.pagination.current = 1
        this.table.pagination.limit = limit
        this.preview()
      },
      handleSortChange(sort) {
        this.table.sort = this.$tools.getSort(sort, '-create_time')
        this.preview()
      },
      renderHeader(h, property) {
        if (!this.table.pagination.count || property.bk_property_id !== 'bk_host_innerip') {
          return <div v-bk-overflow-tips class="table-label">{ this.$tools.getHeaderPropertyName(property) }</div>
        }
        const attrs = {
          // eslint-disable-next-line no-underscore-dangle
          [this.$options._scopeId]: true
        }
        return (
                    <span class="custom-header" {...{ attrs }}>
                        <span>{ this.$tools.getHeaderPropertyName(property) }</span>
                        <i class="icon-cc-copy" v-bk-tooltips="复制IP" {...{ attrs }} on-click={ this.handleCopyIP }></i>
                    </span>
                )
      },
      async handleCopyIP(event) {
        event.stopPropagation()
        const IP = this.table.list.map(item => item.bk_host_innerip)
        try {
          await this.$copyText(IP.join('\n'))
          this.$success(this.$t('复制成功'))
        } catch (error) {
          console.error(error)
          this.$error(this.$t('复制失败'))
        }
      },
      show() {
        this.isShow = true
        setTimeout(this.$refs.table.doLayout, 50)
      },
      handleHidden() {
        this.$emit('close')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .preview-table {
        margin-top: -20px;
    }
    .custom-header {
        display: inline-flex;
        align-items: center;
        .icon-cc-copy {
            cursor: pointer;
            margin: 1px 0 0 4px;
            color: #c0c4cc;
            &:hover {
                color: #63656e;
            }
        }
    }
    .title-text{
        display: inline-block;
        width: 100%;
        font-size: 20px;
        color: #313238;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        margin: 0;
    }
   :deep(.table-label){
        overflow: hidden;
        white-space: nowrap;
        word-wrap: normal;
        text-overflow: ellipsis;
    }
</style>
