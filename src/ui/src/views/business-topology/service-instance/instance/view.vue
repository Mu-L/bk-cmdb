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
  <section class="view-instance" v-test-id="'viewInst'">
    <instance-options class="instance-options"></instance-options>
    <instance-list class="instance-list"></instance-list>
  </section>
</template>

<script>
  import RouterQuery from '@/router/query'
  import Bus from '../common/bus'
  import RootBus from '@/utils/bus'
  import InstanceOptions from './options'
  import InstanceList from './list'
  export default {
    name: 'view-instance',
    components: {
      InstanceOptions,
      InstanceList
    },
    data() {
      return {
      }
    },
    created() {
      Bus.$on('delete-complete', this.refreshView)
    },
    beforeDestroy() {
      Bus.$off('delete-complete', this.refreshView)
    },
    methods: {
      refreshView() {
        // 通知刷新左侧树节点中的服务实例数
        RootBus.$emit('refresh-count-by-node')

        RouterQuery.set({
          _t: Date.now()
        })
      }
    }
  }
</script>
