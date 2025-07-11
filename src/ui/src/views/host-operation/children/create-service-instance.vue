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
  <section class="create-layout">
    <cmdb-tips
      :tips-style="{
        background: 'none',
        border: 'none',
        fontSize: '12px',
        lineHeight: '30px',
        padding: 0
      }"
      :icon-style="{
        color: '#63656E',
        fontSize: '14px',
        lineHeight: '30px'
      }">
      {{$t('新增服务实例提示')}}
    </cmdb-tips>
    <bk-select class="sort-options"
      :style="{ width: $i18n.locale === 'en' ? '110px' : '90px' }"
      :clearable="false"
      prefix-icon="icon-cc-order"
      v-model="sort"
      @selected="sortInfo">
      <bk-option id="ip" name="IP"></bk-option>
      <bk-option id="module" :name="$t('模块')"></bk-option>
    </bk-select>
    <service-instance-table class="service-instance-table"
      v-for="(instance, index) in instances"
      ref="serviceInstance"
      :key="`${instance.bk_module_id}-${instance.bk_host_id}`"
      :index="index"
      :id="instance.bk_host_id"
      :name="getName(instance)"
      :deletable="false"
      :editable="false"
      :editing="getEditState(instance)"
      :topology="$parent.getModulePath(instance.bk_module_id)"
      :templates="getServiceTemplates(instance)"
      :addible="!instance.service_template"
      :source-processes="getSourceProcesses(instance)"
      :class="{ 'is-first': index === 0 }"
      :instance="instance"
      :biz-id="bizId"
      @edit-process="handleEditProcess(instance, ...arguments)"
      @edit-name="handleEditName(instance)"
      @confirm-edit-name="handleConfirmEditName(instance, ...arguments)"
      @cancel-edit-name="handleCancelEditName(instance)">
    </service-instance-table>
  </section>
</template>

<script>
  import ServiceInstanceTable from '@/components/service/instance-table'
  import { mapGetters } from 'vuex'
  import has from 'has'
  import { PROCESS_BIND_IP_ALL_LOOPBACK_AND_NON_MAP } from '@/dictionary/process-bind-ip.js'

  export default {
    name: 'create-service-instance',
    components: {
      ServiceInstanceTable
    },
    props: {
      info: {
        type: Array,
        required: true
      }
    },
    data() {
      return {
        sort: 'module',
        instances: [],
        processChangeState: {}
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId'])
    },
    watch: {
      info: {
        immediate: true,
        handler() {
          this.sortInfo()
        }
      }
    },
    methods: {
      sortInfo() {
        let instances = []
        if (this.sort === 'module') {
          const order = this.$parent.targetModules
          instances = [...this.info].sort((A, B) => order.indexOf(A.bk_module_id) - order.indexOf(B.bk_module_id))
        } else {
          instances = [...this.info].sort((A, B) => this.getName(A).localeCompare(this.getName(B)))
        }
        this.instances = instances.map(instance => ({ ...instance, name: '', editing: { name: false } }))
      },
      getName(instance) {
        if (instance.name) {
          return instance.name
        }
        const data = this.$parent.hostInfo.find(data => data.host.bk_host_id === instance.bk_host_id)
        if (data) {
          return data.host.bk_host_innerip ? data.host.bk_host_innerip : data.host.bk_host_innerip_v6
        }
        return '--'
      },
      getEditState(instance) {
        return instance.editing
      },
      getServiceTemplates(instance) {
        if (instance.service_template) {
          return instance.service_template.process_templates
        }
        return []
      },
      getSourceProcesses(instance) {
        const templates = this.getServiceTemplates(instance)
        return templates.map((template) => {
          const value = {}
          Object.keys(template.property).forEach((key) => {
            const templateValue = template.property[key]
            if (key === 'bind_info') {
              value[key] = (templateValue.value || []).map((info) => {
                const infoValue = {}
                Object.keys(info).forEach((infoKey) => {
                  if (infoKey === 'ip') {
                    infoValue[infoKey] = this.getBindIp(instance, info)
                  } else if (infoKey === 'row_id') {
                    infoValue.template_row_id = info.row_id
                  } else if (typeof info[infoKey] === 'object') {
                    infoValue[infoKey] = info[infoKey].value
                  }
                })
                return infoValue
              })
            } else {
              value[key] = templateValue.value
            }
          })
          return value
        })
      },
      getBindIp(instance, info) {
        const ipValue = info.ip.value
        const mapping = PROCESS_BIND_IP_ALL_LOOPBACK_AND_NON_MAP

        if (has(mapping, ipValue)) {
          return mapping[ipValue]
        }

        const { host } = this.$parent.hostInfo.find(data => data.host.bk_host_id === instance.bk_host_id)
        // 第一内网IP
        if (ipValue === '3') {
          const [innerIP] = host.bk_host_innerip.split(',')
          // 可能不存在，允许返回为空
          return innerIP
        }
        // 第一外网IP
        if (ipValue === '4') {
          const [outerIP] = host.bk_host_outerip.split(',')
          return outerIP
        }
        // 第一内网IPv6
        if (ipValue === '7') {
          const [innerIP] = host?.bk_host_innerip_v6?.split(',') || []
          return innerIP || mapping?.[3]
        }
        // 第一外网IPv6
        if (ipValue === '8') {
          const [outerIP] = host?.bk_host_outerip_v6?.split(',') || []
          return outerIP || mapping?.[3]
        }
        return ''
      },
      getServiceInstanceOptions() {
        const instanceOptions = {
          created: [],
          updated: []
        }

        this.instances.forEach((instance, index) => {
          const component = this.$refs.serviceInstance.find(component => component.index === index)
          const changedProcesses = this.getChangedProcessList(instance, component)
          const addedProcesses = this.getAddedProcessList(instance, component)

          // 有模板
          if (instance.service_template) {
            // 空进程的不能作为添加项
            if (addedProcesses.length) {
              instanceOptions.created.push({
                bk_module_id: instance.bk_module_id,
                bk_host_id: instance.bk_host_id
              })
            }

            // 有修改的项放入到updated
            if (changedProcesses.length) {
              instanceOptions.updated.push({
                bk_module_id: instance.bk_module_id,
                bk_host_id: instance.bk_host_id,
                processes: changedProcesses
              })
            }
          } else {
            // 无模板，在created中放入有添加进程的
            if (addedProcesses.length) {
              instanceOptions.created.push({
                bk_module_id: instance.bk_module_id,
                bk_host_id: instance.bk_host_id,
                processes: addedProcesses
              })
            }
          }
        })

        return instanceOptions
      },
      /**
       * 解决后端性能问题: 用服务模板生成的实例仅传递有被用户主动触发过编辑的进程信息
       */
      getChangedProcessList(instance, component) {
        if (instance.service_template) {
          const processes = []
          const stateKey = `${instance.bk_module_id}-${instance.bk_host_id}`
          const changedState = this.processChangeState[stateKey] || new Set()
          component.processList.forEach((process, listIndex) => {
            if (!changedState.has(listIndex)) return
            processes.push({
              process_template_id: component.templates[listIndex] ? component.templates[listIndex].id : 0,
              process_info: process
            })
          })
          return processes
        }
        return component.processList.map((process, listIndex) => ({
          process_template_id: component.templates[listIndex] ? component.templates[listIndex].id : 0,
          process_info: process
        }))
      },
      /**
       * 解决后端性能问题: 记录用服务模板生成的实例是否触发编辑动作
       */
      handleEditProcess(instance, processIndex) {
        if (!instance.service_template) return
        const key = `${instance.bk_module_id}-${instance.bk_host_id}`
        const state = this.processChangeState[key] || new Set()
        state.add(processIndex)
        this.processChangeState[key] = state
      },
      getAddedProcessList(instance, component) {
        const processes = []
        component.processList.forEach((process, index) => {
          const item = { process_info: process }
          if (instance.service_template) {
            item.process_template_id = component.templates[index] ? component.templates[index].id : 0
          }
          processes.push(item)
        })
        return processes
      },
      handleEditName(instance) {
        this.instances.forEach(instance => (instance.editing.name = false))
        instance.editing.name = true
      },
      handleConfirmEditName(instance, name) {
        instance.name = name
        instance.editing.name = false
      },
      handleCancelEditName(instance) {
        instance.editing.name = false
      }
    }
  }
</script>

<style lang="scss" scoped>
    .create-layout {
        position: relative;
        .sort-options {
            position: absolute;
            top: -4px;
            right: 0;
            /deep/ .icon-cc-order {
                font-size: 14px;
                color: #979BA5;
            }
        }
    }
    .service-instance-table {
        &.is-first {
            margin-top: 8px;
        }
        & + .service-instance-table {
            margin: 15px 0 0 0;
        }
    }
</style>
