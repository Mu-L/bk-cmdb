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
  <section class="filter-wrapper" v-if="selected.length || showIPTag">
    <label class="filter-label">
      <i class="label-icon icon-cc-funnel"></i>
      <span class="label-text">{{$t('检索项')}}</span>
      <span class="label-colon">:</span>
    </label>
    <div class="filter-list" ref="filterList">
      <filter-tag-ip v-if="showIPTag"></filter-tag-ip>
      <filter-tag-item
        v-for="property in selected"
        :key="property.id"
        :property="property"
        v-bind="condition[property.id]">
      </filter-tag-item>
      <bk-button class="filter-clear" text
        v-if="showClear"
        @click="handleResetAll">
        {{$t('清空条件')}}
      </bk-button>
    </div>
  </section>
</template>

<script>
  import FilterTagIp from './filter-tag-ip'
  import FilterTagItem from './filter-tag-item'
  import FilterStore from './store'
  import Utils from './utils'
  export default {
    components: {
      FilterTagIp,
      FilterTagItem
    },
    computed: {
      condition() {
        return FilterStore.condition
      },
      showIPTag() {
        const list = Utils.splitIP(FilterStore.IP.text)
        return !!list.length
      },
      selected() {
        return FilterStore.selected.filter((property) => {
          const { value } = this.condition[property.id]
          return value !== null && value !== undefined && !!value.toString().length
        })
      },
      showClear() {
        const count = this.selected.length + (this.showIPTag ? 1 : 0)
        return count > 1
      }
    },
    watch: {
      selected() {
        if (!(this.selected.length || this.showIPTag)
          && FilterStore.activeCollection) FilterStore.setActiveCollection(null)
      }
    },
    methods: {
      handleResetAll() {
        FilterStore.resetAll()
        FilterStore.setActiveCollection(null)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .filter-wrapper {
        display: flex;
        margin: 10px 0 0 0;
        .filter-label {
            display: flex;
            font-size: 12px;
            align-items: center;
            align-self: flex-start;
            line-height: 22px;
            .label-icon {
                color: #979BA5;
            }
            .label-text {
                margin-left: 4px;
            }
            .label-colon {
                margin: 0 5px;
            }
        }
        .filter-list {
            display: flex;
            flex-wrap: wrap;
            flex: 1;
        }
        .filter-clear {
            line-height: initial;
            margin: 0 0 10px 10px;
            font-size: 12px;
        }
    }
</style>
