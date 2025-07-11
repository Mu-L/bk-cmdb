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
  <div class="resource-layout clearfix">
    <bk-tab
      :active.sync="activeTab"
      class="scope-tab"
      type="unborder-card"
      @tab-change="handleTabChange">
      <bk-tab-panel v-for="item in scopeList"
        :key="item.id"
        :name="item.id"
        :label="item.label">
      </bk-tab-panel>
    </bk-tab>
    <div class="content">
      <cmdb-resize-layout
        v-if="isResourcePool"
        :class="['resize-layout fl', { 'is-collapse': layout.collapse }]"
        :handler-offset="3"
        :min="200"
        :max="480"
        :disabled="layout.collapse"
        direction="right">
        <resource-directory ref="resourceDirectory"></resource-directory>
        <i class="directory-collapse-icon bk-icon icon-angle-left"
          @click="layout.collapse = !layout.collapse">
        </i>
      </cmdb-resize-layout>
      <resource-hosts class="main" ref="resourceHost" @refresh="handleRefresh"></resource-hosts>
    </div>
    <router-subview></router-subview>
    <cmdb-model-fast-link :obj-id="objId"></cmdb-model-fast-link>
  </div>
</template>

<script>
  import FilterStore from '@/components/filters/store'
  import resourceDirectory from './children/directory.vue'
  import resourceHosts from './children/host-list.vue'
  import Bus from '@/utils/bus.js'
  import RouterQuery from '@/router/query'
  import cmdbModelFastLink from '@/components/model-fast-link'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'

  export default {
    components: {
      resourceDirectory,
      resourceHosts,
      cmdbModelFastLink
    },
    data() {
      return {
        layout: {
          collapse: false
        },
        activeTab: RouterQuery.get('scope', '1'),
        scopeList: [{
          id: '1',
          label: this.$t('未分配')
        }, {
          id: '0',
          label: this.$t('已分配')
        }, {
          id: 'all',
          label: this.$t('全部')
        }]
      }
    },
    computed: {
      isResourcePool() {
        return this.activeTab.toString() === '1'
      },
      objId() {
        return BUILTIN_MODELS.HOST
      }
    },
    watch: {
      activeTab(val) {
        const scope = RouterQuery.get('scope', '1')
        if (scope !== val) RouterQuery.set('scope', val)
      }
    },
    methods: {
      handleRefresh() {
        const { resourceHost, resourceDirectory } = this.$refs
        resourceHost.getHostList()
        resourceDirectory.getDirectoryList()
      },
      handleTabChange(tab) {
        Bus.$emit('toggle-host-filter', false)
        Bus.$emit('reset-host-filter')

        // 设置scope
        FilterStore.setResourceScope(tab)

        // 此时selected为上一个scope的，需要清空，在setupNormalProperty方法中会使用在设置条件时已保存的值
        FilterStore.updateSelected([])
        FilterStore.setupNormalProperty()

        RouterQuery.set({
          scope: isNaN(tab) ? tab : parseInt(tab, 10),
          ip: '',
          bk_asset_id: '',
          page: 1,
          _t: Date.now()
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .resource-layout{
        .scope-tab {
            height: auto;
            margin: 0 20px;
            /deep/ .bk-tab-header {
                padding: 0;
            }
        }
        .content {
            height: calc(100% - 58px);
            overflow: hidden;
        }
        .resize-layout {
            position: relative;
            width: 280px;
            height: 100%;
            border-right: 1px solid $cmdbLayoutBorderColor;
            &.is-collapse {
                width: 0 !important;
                border-right: none;
                .directory-collapse-icon:before {
                    display: inline-block;
                    transform: rotate(180deg);
                }
            }
            .directory-collapse-icon {
                position: absolute;
                left: 100%;
                top: 50%;
                width: 16px;
                height: 100px;
                line-height: 100px;
                background: $cmdbLayoutBorderColor;
                border-radius: 0px 12px 12px 0px;
                transform: translateY(-50%);
                text-align: center;
                text-indent: -2px;
                font-size: 20px;
                color: #fff;
                cursor: pointer;
                &:hover {
                    background: #699DF4;
                }
            }
        }
        .main {
            height: 100%;
            padding: 10px 20px 0 20px;
            overflow: hidden;
        }
    }
    .assign-dialog {
        /deep/ .bk-dialog-body {
            padding: 0 50px 40px;
        }
        .assign-info span {
            color: #3c96ff;
        }
        .assign-footer {
            padding-top: 20px;
            font-size: 0;
            text-align: center;
            .bk-button-normal {
                width: 96px;
            }
        }
    }
</style>
