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
  <div class="property">
    <div class="group"
      v-for="(group, index) in $sortedGroups"
      :key="index">
      <h2 class="group-name">{{group.bk_group_name}}</h2>
      <ul class="property-list">
        <li :class="['property-item', property.bk_property_type]"
          v-for="property in $groupedProperties[index]"
          :key="property.id"
          :id="`property-item-${property.id}`">
          <div class="property-name" v-bk-overflow-tips>
            <span v-if="isContainerObjects(property.bk_obj_id)" class="no-show-more">
              {{property.bk_property_name}}
            </span>
            <span v-else
              @mouseenter="(event) => handlePropertyNameMouseenter(event, property)"
            >
              {{property.bk_property_name}}
            </span>
            <i class="property-name-tooltips icon-cc-tips"
              v-if="property.placeholder && $tools.isIconTipProperty(property.bk_property_type)"
              v-bk-tooltips.top="{
                theme: 'light',
                trigger: 'mouseenter',
                content: property.placeholder
              }">
            </i>
          </div>
          <template v-if="property.bk_property_type !== PROPERTY_TYPES.INNER_TABLE">
            <div :class="['property-value', { 'is-loading': loadingState.includes(property) }]"
              v-if="property !== editState.property">
              <cmdb-property-value
                tag="div"
                :is-show-overflow-tips="isShowOverflowTips(property)"
                :ref="`property-value-${property.bk_property_id}`"
                :value="instState[property.bk_property_id]"
                :property="property"
                :instance="instState">
              </cmdb-property-value>
            </div>
            <template v-if="!loadingState.includes(property)">
              <template v-if="!readonly && !isPropertyEditable(property)">
                <i class="is-related property-edit icon-cc-edit-shape"
                  v-bk-tooltips="{
                    content: $t('系统限定不可修改'),
                    placement: 'top',
                    onShow: () => {
                      setFocus(`#property-item-${property.id}`, true)
                    },
                    onHide: () => {
                      setFocus(`#property-item-${property.id}`, false)
                    }
                  }">
                </i>
              </template>
              <template v-else-if="!readonly">
                <cmdb-auth
                  style="margin: 8px 0 0 8px; font-size: 0;"
                  :auth="authData"
                  v-show="property !== editState.property">
                  <bk-button slot-scope="{ disabled }"
                    text
                    theme="primary"
                    class="property-edit-btn"
                    :disabled="disabled"
                    @click="setEditState(property)">
                    <i class="property-edit icon-cc-edit-shape"></i>
                  </bk-button>
                </cmdb-auth>
                <div class="property-form" v-if="property === editState.property"
                  @keyup="(event) => keyupCallMethodFn(event)">
                  <div :class="['form-component', property.bk_property_type]">
                    <component
                      :is="`cmdb-form-${property.bk_property_type}`"
                      :class="[property.bk_property_type, { error: errors.has(property.bk_property_id) }]"
                      :unit="property.unit"
                      :options="property.option || []"
                      :data-vv-name="property.bk_property_id"
                      :data-vv-as="property.bk_property_name"
                      :placeholder="$tools.getPropertyPlaceholder(property)"
                      :auto-check="false"
                      :multiple="property.ismultiple"
                      v-bind="$tools.getValidateEvents(property)"
                      v-validate="$tools.getValidateRules(property)"
                      v-model.trim="editState.value"
                      @focus="handleFocus"
                      @blur="handleBlur"
                      v-bk-tooltips.top="{
                        disabled: !property.placeholder || $tools.isIconTipProperty(property.bk_property_type),
                        theme: 'light',
                        showOnInit: true,
                        trigger: 'click',
                        content: property.placeholder
                      }"
                      :ref="`component-${property.bk_property_id}`">
                    </component>
                  </div>
                  <i class="form-confirm bk-icon icon-check-1" @click="confirm"></i>
                  <i class="form-cancel bk-icon icon-close" @click="exitForm"></i>
                  <cmdb-default-picker
                    v-if="showDefault(property.bk_property_id)"
                    :value="propertyDefaults[property.bk_property_id]"
                    :property="property"
                    :instance="instState"
                    @pick-default="handlePickDefault">
                  </cmdb-default-picker>
                  <span class="form-error"
                    v-else-if="errors.has(property.bk_property_id)">
                    {{errors.first(property.bk_property_id)}}
                  </span>
                </div>
              </template>
              <template v-if="!$tools.isEmptyPropertyValue(instState[property.bk_property_id])
                && property !== editState.property">
                <div class="copy-box">
                  <i
                    class="property-copy icon-copy bk-cmdb-icon"
                    @click="handleCopy(property.bk_property_id)">
                  </i>
                  <transition name="fade">
                    <span class="copy-tips"
                      :style="{ width: $i18n.locale === 'en' ? '100px' : '70px' }"
                      v-if="showCopyTips === property.bk_property_id">
                      {{$t('复制成功')}}
                    </span>
                  </transition>
                </div>
              </template>
            </template>
          </template>
          <template v-else>
            <cmdb-form-innertable
              class="form-component-innertable"
              :mode="'update'"
              :disabled="!isPropertyEditable(property)"
              :disabled-tips="$t('系统限定不可修改')"
              :property="property"
              :obj-id="objId"
              :instance-id="instanceId"
              :auth="authData"
              :ref="`component-${property.bk_property_id}`" />
          </template>
        </li>
      </ul>
    </div>
    <slot name="append"></slot>

    <!-- hover字段摘要 -->
    <div ref="propertySummary" v-show="hoverPropertyPopover.show">
      <dl class="property-summary-content">
        <dt class="content-head">
          <span class="name">{{hoverPropertyPopover.data.bk_property_name}}</span>
          <span class="more-link" @click="handleViewProperty(hoverPropertyPopover.data.bk_property_id)">
            {{$t('更多信息')}}<i class="link-icon icon-cc-share"></i>
          </span>
        </dt>
        <div class="content-body">
          <dd class="row-item">
            <span class="item-name">{{$t('唯一标识')}}</span>
            <span class="item-value">{{hoverPropertyPopover.data.bk_property_id}}</span>
            <i class="copy-icon icon-cc-details-copy"
              @click="handleCopyPropertyId(hoverPropertyPopover.data.bk_property_id)"></i>
          </dd>
          <dd class="row-item">
            <span class="item-name">{{$t('字段名称')}}</span>
            <span class="item-value">{{hoverPropertyPopover.data.bk_property_name}}</span>
          </dd>
          <dd class="row-item">
            <span class="item-name">{{$t('字段类型')}}</span>
            <span class="item-value">
              {{fieldTypeMap[hoverPropertyPopover.data.bk_property_type]}}
              ({{hoverPropertyPopover.data.bk_property_type}})
            </span>
          </dd>
          <dd class="row-item" v-show="hoverPropertyPopover.data.placeholder">
            <span class="item-name">{{$t('用户提示')}}</span>
            <span class="item-value">{{hoverPropertyPopover.data.placeholder}}</span>
          </dd>
        </div>
      </dl>
    </div>
  </div>
</template>

<script>
  import { mapGetters, mapActions } from 'vuex'
  import isEqual from 'lodash/isEqual'
  import qs from 'qs'
  import formMixins from '@/mixins/form'
  import {
    BUILTIN_MODELS,
    BUILTIN_MODEL_PROPERTY_KEYS,
    BUILTIN_MODEL_RESOURCE_TYPES
  } from '@/dictionary/model-constants.js'
  import { MENU_MODEL_DETAILS } from '@/dictionary/menu-symbol'
  import businessSetService from '@/service/business-set/index.js'
  import projectService from '@/service/project/index.js'
  import authMixin from './mixin-auth'
  import { PROPERTY_TYPES, PROPERTY_TYPE_NAMES } from '@/dictionary/property-constants'
  import { keyupCallMethod } from '@/utils/util'
  import cmdbDefaultPicker from '@/components/ui/other/default-value-picker'
  import { isContainerObjects } from '@/utils/tools'

  export default {
    filters: {
      filterShowText(value, unit) {
        return value === '--' ? '--' : value + unit
      }
    },
    components: {
      cmdbDefaultPicker
    },
    mixins: [formMixins, authMixin],
    props: {
      inst: {
        type: Object,
        required: true
      },
      resourceType: {
        type: String,
        default: ''
      },
      readonly: {
        type: Boolean,
        default: false
      },
      objId: {
        type: String
      },
      showDefaultValue: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        PROPERTY_TYPES,
        instState: this.inst,
        propertyDefaults: this.$tools.getInstFormDefaults(this.properties),
        editState: {
          property: null,
          value: null,
          focus: false
        },
        loadingState: [],
        showCopyTips: false,
        fieldTypeMap: PROPERTY_TYPE_NAMES,
        hoverPropertyPopover: {
          show: false,
          data: {},
          instance: null
        }
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['models', 'getModelById']),
      authData() {
        const auths = {
          [BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS]]: this.INST_AUTH.U_BUSINESS,
          [BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS_SET]]: this.INST_AUTH.U_BUSINESS_SET,
          [BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.PROJECT]]: this.INST_AUTH.U_PROJECT
        }
        return auths[this.resourceType] || this.INST_AUTH.U_INST
      },
      instanceId() {
        return this.instState?.[BUILTIN_MODEL_PROPERTY_KEYS[this.objId]?.ID || 'bk_inst_id']
      }
    },
    watch: {
      inst(val) {
        this.instState = this.$tools.getInstFormValues(this.properties, val, false)
        this.propertyDefaults = this.$tools.getInstFormDefaults(this.properties)
      }
    },
    methods: {
      ...mapActions('objectCommonInst', ['updateInst']),
      ...mapActions('objectBiz', ['updateBusiness']),
      setFocus(id, focus) {
        const item = this.$el.querySelector(id)
        focus ? item.classList.add('focus') : item.classList.remove('focus')
      },
      keyupCallMethodFn(event) {
        keyupCallMethod(event, this.confirm)
      },
      getPlaceholder(property) {
        const placeholderTxt = ['enum', 'list', 'organization'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
        return this.$t(placeholderTxt, { name: property.bk_property_name })
      },
      isContainerObjects(objId) {
        return isContainerObjects(objId)
      },
      isPropertyEditable(property) {
        return property.editable && !property.bk_isapi
      },
      isShowOverflowTips(property) {
        const complexTypes = [PROPERTY_TYPES.MAP, PROPERTY_TYPES.ENUMQUOTE]
        return !complexTypes.includes(property.bk_property_type)
      },
      setEditState(property) {
        const value = this.instState[property.bk_property_id]
        this.editState.value = (value === null || value === undefined) ? '' : value
        this.editState.property = property
        setTimeout(() => {
          const component = this.$refs[`component-${property.bk_property_id}`]
          component?.[0]?.focus()
        }, 100)
      },
      async confirm() {
        const { property, value } = this.editState
        try {
          const isValid = await this.$validator.validateAll()
          if (!isValid) {
            return false
          }
          this.exitForm()
          const oldValue = this.instState[property.bk_property_id]
          if (isEqual(oldValue, value)) return

          this.loadingState.push(property)
          const values = { [property.bk_property_id]: this.$tools.formatValue(value, property) }

          if (this.resourceType === BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS]) {
            await this.updateBusiness({
              bizId: this.instState.bk_biz_id,
              params: values
            })
          } else if (this.resourceType === BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS_SET]) {
            const MODEL_ID_KEY = BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID
            await businessSetService.update({
              bk_biz_set_ids: [this.instState[MODEL_ID_KEY]],
              data: {
                bk_biz_set_attr: { ...values },
              }
            })
          } else if (this.resourceType === BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.PROJECT]) {
            const params = {
              ids: [this.instState.id],
              data: values
            }
            await projectService.update(params)
          } else {
            await this.updateInst({
              objId: this.instState.bk_obj_id,
              instId: this.instState.bk_inst_id,
              params: values
            })
          }

          this.$success(this.$t('修改成功'))

          this.instState = { ...this.instState, ...values }

          this.loadingState = this.loadingState.filter(exist => exist !== property)

          this.$emit('after-update')
        } catch (e) {
          console.error(e)
          this.loadingState = this.loadingState.filter(exist => exist !== property)
        }
      },
      exitForm() {
        this.editState.property = null
        this.editState.value = null
      },
      showPopover() {
        this.hoverPropertyPopover.instance.show()
      },
      showDefault(propertyId) {
        const { value, focus } = this.editState
        return this.propertyDefaults[propertyId]
          && this.showDefaultValue
          && !value
          && focus
      },
      handleFocus() {
        this.editState.focus = true
      },
      handleBlur() {
        this.editState.focus = false
      },
      handlePickDefault(val) {
        this.editState.value = val
      },
      handleCopy(propertyId) {
        const [component] = this.$refs[`property-value-${propertyId}`]
        const copyText = component?.getCopyValue() ?? ''
        this.$copyText(copyText).then(() => {
          this.showCopyTips = propertyId
          const timer = setTimeout(() => {
            this.showCopyTips = false
            clearTimeout(timer)
          }, 200)
        }, () => {
          this.$error(this.$t('复制失败'))
        })
      },
      handleCopyPropertyId(propertyId) {
        this.$copyText(propertyId).then(() => {
          this.$success(this.$t('复制成功'))
        }, () => {
          this.$error(this.$t('复制失败'))
        })
      },
      handlePropertyNameMouseenter(event, property) {
        const { bk_property_id: propertyId } = property
        const { bk_property_id: lastPropertyId } = this.hoverPropertyPopover.data
        if (propertyId === lastPropertyId) {
          return this.showPopover()
        }

        this.hoverPropertyPopover.instance?.destroy?.()
        this.hoverPropertyPopover.data = property
        this.hoverPropertyPopover.instance = this.$bkPopover(event.target, {
          content: this.$refs.propertySummary,
          delay: [300, 0],
          hideOnClick: true,
          interactive: true,
          placement: 'top',
          animateFill: false,
          sticky: true,
          theme: 'light property-summary-popover',
          boundary: 'window',
          trigger: 'mouseenter', // 'manual mouseenter',
          arrow: true,
          onShow: () => {
            this.hoverPropertyPopover.show = true
          },
          onHidden: () => {
            this.hoverPropertyPopover.show = false
          }
        })

        this.showPopover()
      },
      handleViewProperty(propertyId) {
        this.$routerActions.open({
          name: MENU_MODEL_DETAILS,
          params: {
            modelId: this.objId,
          },
          query: {
            action: 'view-field',
            payload: qs.stringify({
              id: propertyId
            })
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .property-name-tooltips {
      margin-right: 2px;
    }
    .property {
        height: 100%;
        overflow: auto;
        @include scrollbar-y;
    }
    .group {
        margin: 22px 0 0 0;
        .group-name {
            line-height: 21px;
            font-size: 16px;
            font-weight: normal;
            color: #333948;
            &:before {
                content: "";
                display: inline-block;
                vertical-align: -2px;
                width: 4px;
                height: 14px;
                margin-right: 9px;
                background-color: $cmdbBorderColor;
            }
        }
    }
    .property-list {
        width: 1208px;
        margin: 25px 0 0 0;
        color: #63656e;
        display: flex;
        flex-wrap: wrap;
        .property-item {
            flex: 0 0 50%;
            max-width: 50%;
            padding-bottom: 8px;
            display: flex;
            &:hover,
            &.focus {
                .property-edit {
                    opacity: 1;
                }
                .property-copy {
                    display: inline-block;
                }
            }

            &.innertable {
              flex: 0 0 100%;
              max-width: 100%;
              .property-name {
                flex: none;
              }
              .form-component-innertable {
                flex: none;
                width: calc(100% - 260px);
                margin-top: 6px;
              }
            }

            .property-name {
                position: relative;
                width: 260px;
                line-height: 32px;
                padding: 0 16px 0 36px;
                font-size: 14px;
                color: #63656E;
                text-align: right;
                @include ellipsis;
                &:after {
                    position: absolute;
                    right: 2px;
                    content: "：";
                }
                :first-child:not(.no-show-more) {
                  &:hover {
                    color: $primaryColor;
                  }
                }
            }
            .property-value {
                margin: 6px 0 0 4px;
                max-width: 286px;
                font-size: 14px;
                color: #313237;
                overflow:hidden;
                text-overflow:ellipsis;
                word-break: break-all;
                display: -webkit-box;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
                &.is-loading {
                    font-size: 0;
                    &:before {
                        content: "";
                        display: inline-block;
                        width: 16px;
                        height: 16px;
                        margin: 2px 0;
                        background-image: url("../../assets/images/icon/loading.svg");
                    }
                }
            }
            .property-edit-btn {
                height: auto;
                font-size: 0;
            }
            .property-edit {
                font-size: 16px;
                opacity: 0;
                &.is-related {
                    display: inline-block;
                    vertical-align: middle;
                    width: 16px;
                    height: 16px;
                    margin: 8px 0 0 8px;
                    line-height: 1;
                    &:hover {
                        color: inherit !important;
                    }
                }
                &:hover {
                  color: #1768EF;
                }
            }
            .property-copy {
                margin: 8px 0 0 8px;
                color: #3c96ff;
                cursor: pointer;
                display: none;
                font-size: 16px;
            }
            .copy-box {
                position: relative;
                font-size: 0;
                .copy-tips {
                    position: absolute;
                    top: -22px;
                    left: -18px;
                    min-width: 70px;
                    height: 26px;
                    line-height: 26px;
                    font-size: 12px;
                    color: #ffffff;
                    text-align: center;
                    background-color: #9f9f9f;
                    border-radius: 2px;
                }
                .fade-enter-active, .fade-leave-active {
                    transition: all 0.5s;
                }
                .fade-enter {
                    top: -14px;
                    opacity: 0;
                }
                .fade-leave-to {
                    top: -28px;
                    opacity: 0;
                }
            }
        }
    }
    .property-form {
        font-size: 0;
        position: relative;
        .bk-icon {
            display: inline-block;
            vertical-align: middle;
            width: 32px;
            height: 32px;
            margin: 0 0 0 6px;
            border-radius: 2px;
            border: 1px solid #c4c6cc;
            line-height: 30px;
            font-size: 20px;
            text-align: center;
            cursor: pointer;
            &.form-confirm {
                color: #0082ff;
                font-size: 20px;
                &:before {
                    display: inline-block;
                }
            }
            &.form-cancel {
                color: #979ba5;
                font-size: 20px;
                &:before {
                    display: inline-block;
                }
            }
            &:hover {
                font-weight: bold;
            }
        }
        .form-error {
            position: absolute;
            top: 100%;
            left: 0;
            font-size: 12px;
            line-height: 1;
            color: $cmdbDangerColor;
        }
        .form-component {
            display: inline-block;
            vertical-align: middle;
            height: 32px;
            width: 260px;
            margin: 0 4px 0 0;
            position: relative;
            &.bool {
                width: 42px;
                height: 24px;
            }
            &.longchar {
                height: auto;
                vertical-align: top;
            }
        }
    }

    .property-summary-content {
      font-size: 14px;
      padding: 8px 12px;
      min-width: 260px;
      .content-head {
        display: flex;
        justify-content: space-between;
        align-items: center;
        .name {
          font-weight: 700;
        }
        .more-link {
          display: flex;
          align-items: center;
          gap: 2px;
          color: $primaryColor;
          cursor: pointer;
        }
      }
      .content-body {
        display: flex;
        flex-direction: column;
        gap: 6px;
        margin-top: 8px;
        .copy-icon {
          color: $primaryColor;
          cursor: pointer;
          margin-left: 8px;
        }
        .row-item {
          display: flex;
          align-items: center;
          .item-name {
            position: relative;
            padding-right: 14px;
            &::after {
              position: absolute;
              right: 0;
              content: "：";
            }
          }
        }
      }
    }
</style>
<style>
  .tippy-tooltip.property-summary-popover-theme {
    box-shadow: 0 0 6px 0.5px #dcdee5;
  }
</style>
