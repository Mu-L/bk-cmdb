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
  <cmdb-form-table class="cmdb-form-process-table"
    v-bind="$attrs"
    v-model="localValue"
    :options="options"
    :mode="mode">
    <template v-for="column in options" #[column.bk_property_id]="rowProps">
      <div class="process-table-content"
        :key="`row-${rowProps.index}-${column.bk_property_id}`"
        :class="{ 'is-lock': isLocked(rowProps) }">
        <component class="content-value"
          size="small"
          font-size="small"
          v-bind="getEvents(column)"
          v-validate="getRules(rowProps, column)"
          :data-vv-name="column.bk_property_id"
          :data-vv-as="column.bk_property_name"
          :data-vv-scope="column.bk_property_group || 'bind_info'"
          :is="getComponentType(column)"
          :options="getOptions(rowProps, column)"
          :placeholder="getPlaceholder(column)"
          :value="localValue[rowProps.index][column.bk_property_id]"
          :auto-select="false"
          :readonly="readonly"
          :disabled="column.bk_property_type === 'bool' && readonly"
          @input="handleColumnValueChange(rowProps, ...arguments)">
        </component>
        <span :class="['property-lock-state', { disabled: readonly }]"
          v-bk-tooltips="{
            placement: 'top',
            interactive: false,
            content: $t('进程模板加解锁提示语'),
            delay: [100, 0]
          }"
          tabindex="-1"
          @click="setLockState(rowProps)">
          <i class="icon-cc-lock-fill" v-if="isLocked(rowProps)"></i>
          <i class="icon-cc-unlock-fill" v-else></i>
        </span>
      </div>
    </template>
  </cmdb-form-table>
</template>

<script>
  import ProcessFormPropertyIp from './process-form-property-ip'
  import { PROCESS_BIND_IPV4_MAP, PROCESS_BIND_IPV6_MAP } from '@/dictionary/process-bind-ip.js'
  import { PROCESS_BIND_PROTOCOL_V4_LIST, PROCESS_BIND_PROTOCOL_V6_LIST } from '@/dictionary/process-bind-protocol.js'

  export default {
    components: {
      ProcessFormPropertyIp
    },
    props: {
      value: {
        type: Array,
        default: () => ([])
      },
      options: {
        type: Array,
        required: true
      },
      readonly: Boolean
    },
    inject: {
      type: { default: '' } // from ./process-form
    },
    computed: {
      localValue: {
        get() {
          return (this.value || []).map((row) => {
            const rowValue = {}
            Object.keys(row).forEach((key) => {
              if (['process_id', 'row_id'].includes(key)) {
                rowValue[key] = row[key]
              } else {
                rowValue[key] = row[key].value
              }
            })
            return rowValue
          })
        },
        set(values) {
          const newValues = this.transformValueToTemplateValue(values)
          this.$emit('input', newValues)
          this.$emit('change', newValues)
        }
      },
      lockStates: {
        get() {
          return (this.value || []).map((row) => {
            const rowState = {}
            Object.keys(row).forEach((key) => {
              if (['process_id', 'row_id'].includes(key)) {
                rowState[key] = row[key]
              } else {
                rowState[key] = row[key].as_default_value
              }
            })
            return rowState
          })
        },
        set(states) {
          const newValues = this.transformStateToTemplateValue(states)
          this.$emit('input', newValues)
          this.$emit('change', newValues)
        }
      },
      defaultRowValue() {
        return {
          // ip为字符串类型，模板提供内置两种枚举选项，模板锁定时，默认选择127.0.0.1
          locked: this.$tools.getInstFormValues(this.options, { ip: '1' }, true),
          unlocked: this.$tools.getInstFormValues(this.options, {}, false)
        }
      },
      IPV4Keys() {
        return Object.keys(PROCESS_BIND_IPV4_MAP)
      },
      IPV6Keys() {
        return Object.keys(PROCESS_BIND_IPV6_MAP)
      },
      mode() {
        return this.readonly ? 'info' : 'update'
      }
    },
    methods: {
      isLocked({ column, index }) {
        return this.lockStates[index][column.property]
      },
      setLockState(rowProps) {
        if (this.readonly) {
          return
        }
        const { column, index } = rowProps
        const lockState = { ...(this.lockStates[index] || {}) }
        lockState[column.property] = !this.isLocked(rowProps)
        const newStates = [...this.lockStates]
        newStates.splice(index, 1, lockState)
        this.lockStates = newStates
      },
      getComponentType(property) {
        if (property.bk_property_id === 'ip') {
          return 'process-form-property-ip'
        }
        return `cmdb-form-${property.bk_property_type}`
      },
      getPlaceholder(property) {
        const placeholderTxt = ['enum', 'list'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
        return this.$t(placeholderTxt, { name: property.bk_property_name })
      },
      getEvents(property) {
        if (property.bk_property_id === 'ip') {
          return {}
        }
        return this.$tools.getValidateEvents(property)
      },
      getRules(rowProps, property) {
        const rules = this.$tools.getValidateRules(property)
        rules.required = true
        // IP字段在模板上被构造为枚举，无法通过ip的正则，此处忽略IP正则
        if (property.bk_property_id === 'ip') {
          rules.required = false
          delete rules.remoteString
        }
        return rules
      },
      getOptions({ row }, property) {
        // protocol列的选项与ip类型联动，ipv4对应v4的协议选项，ipv6则对应v6的协议选项
        if (property.bk_property_id === 'protocol') {
          if (this.IPV4Keys.includes(row.ip)) {
            return property.option.filter(item => PROCESS_BIND_PROTOCOL_V4_LIST.includes(item.id))
          }
          if (this.IPV6Keys.includes(row.ip)) {
            return property.option.filter(item => PROCESS_BIND_PROTOCOL_V6_LIST.includes(item.id))
          }
        }

        return property.option || []
      },
      handleColumnValueChange({ row, column, index }, value) {
        const rowValue = { ...row }

        // 变更ip值时，如果ip类型与协议类型值不一致则需要重置其值
        if (column.property === 'ip') {
          if (this.IPV4Keys.includes(value) && !PROCESS_BIND_PROTOCOL_V4_LIST.includes(rowValue.protocol)) {
            rowValue.protocol = ''
          }
          if (this.IPV6Keys.includes(value) && !PROCESS_BIND_PROTOCOL_V6_LIST.includes(rowValue.protocol)) {
            rowValue.protocol = ''
          }
        }

        rowValue[column.property] = value
        const newValues = [...this.localValue]
        newValues.splice(index, 1, rowValue)
        this.localValue = newValues
      },
      // 将常规表格数据，转换成服务模板需要的数据格式
      transformValueToTemplateValue(values) {
        const isAddOrDelete = values.length !== this.localValue.length
        return values.map((row, rowIndex) => {
          const templateRowValue = {}
          // 获取新value中每行对应的老数据的index，用于正确的获取checkbox勾选状态
          const index = isAddOrDelete ? this.localValue.indexOf(row) : rowIndex

          // 创建模式并且是新添加的行，使其默认锁定
          const defaultLocked = this.type === 'create' && rowIndex === values.length - 1

          Object.keys(row).forEach((key) => {
            if (['process_id', 'row_id'].includes(key)) {
              templateRowValue[key] = row[key]
            } else {
              templateRowValue[key] = {
                value: row[key],
                as_default_value: !!(this.lockStates[index] || {})[key] || defaultLocked
              }
            }
          })
          return templateRowValue
        })
      },
      // 将常规表格行锁定状态，转换成服务模板需要的数据格式
      transformStateToTemplateValue(states) {
        return states.map((row, rowIndex) => {
          const templateRowValue = {}
          Object.keys(row).forEach((key) => {
            if (['process_id', 'row_id'].includes(key)) {
              templateRowValue[key] = row[key]
            } else {
              const value = this.localValue[rowIndex][key]
              templateRowValue[key] = {
                value,
                as_default_value: row[key]
              }
            }
          })
          return templateRowValue
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    @mixin property-lock-state-visible {
        display: inline-flex;
        border: 1px solid #c4c6cc;
        border-left: none;
    }
    @mixin no-right-radius {
        border-top-right-radius: 0;
        border-bottom-right-radius: 0;
    }
    .cmdb-form-process-table {
        .process-table-content {
            display: flex;
            align-items: center;
            justify-content: flex-start;
            &:hover,
            &.is-lock {
                .property-lock-state {
                    @include property-lock-state-visible;
                    &.disabled {
                      color: #dcdee5;
                    }
                }
                .content-value /deep/ {
                    .bk-form-input,
                    .bk-form-textarea,
                    .bk-textarea-wrapper {
                        @include no-right-radius;
                    }
                }
                .content-value.bk-select {
                    @include no-right-radius;
                }
            }
            .content-value {
                &:not(.bk-switcher) {
                  flex: 1;
                  width: calc(100% - 24px);
                }
                &.control-active /deep/ {
                    .bk-form-input,
                    .bk-form-textarea,
                    .bk-textarea-wrapper {
                        @include no-right-radius;
                    }
                }
                &.control-active ~ .property-lock-state {
                    @include property-lock-state-visible;
                }
                &.is-focus {
                    @include no-right-radius;
                }
                &.form-bool {
                    flex: none;
                    & ~ .property-lock-state {
                        border: none;
                        background-color: transparent;
                        transition: none;
                    }
                }
            }
            .property-lock-state {
                display: none;
                width: 24px;
                height: 26px;
                align-items: center;
                justify-content: center;
                background-color: #f2f4f8;
                font-size: 14px;
                overflow: hidden;
                transition: width .1s linear;
                cursor: pointer;
            }
        }
    }
</style>
