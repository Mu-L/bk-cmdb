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
  <div class="create-node-layout">
    <h2 class="node-create-title">{{$t('新建集群')}}</h2>
    <div class="node-create-path" :title="topoPath">{{$t('添加节点已选择')}}：{{topoPath}}</div>
    <div class="node-create-form">
      <bk-radio-group class="form-item mb20" v-model="withTemplate">
        <bk-radio :value="true">{{$t('从模板新建')}}</bk-radio>
        <bk-radio :value="false">{{$t('直接新建')}}</bk-radio>
      </bk-radio-group>
      <div class="form-item" v-if="withTemplate">
        <label>{{$t('集群模板')}}</label>
        <bk-select style="width: 100%;"
          :clearable="false"
          :searchable="setTemplateList.length > 7"
          :placeholder="$t('请选择xx', { name: $t('集群模板') })"
          :loading="$loading(request.setTemplate)"
          v-model="setTemplate"
          v-validate.disabled="'required'"
          data-vv-name="setTemplate">
          <bk-option v-for="option in setTemplateList"
            :key="option.id"
            :id="option.id"
            :name="option.name">
          </bk-option>
          <div class="add-template" slot="extension" @click="handleAddTemplate" v-if="!setTemplateList.length">
            <i class="bk-icon icon-plus-circle"></i>
            <span>{{$t('创建集群模板')}}</span>
          </div>
        </bk-select>
        <span class="form-error" v-if="errors.has('setTemplate')">{{errors.first('setTemplate')}}</span>
      </div>
      <div class="form-item">
        <label>
          {{$t('集群名称')}}
          <span class="red-star">*</span>
        </label>
        <bk-input class="form-textarea"
          type="textarea"
          data-vv-name="setName"
          v-validate="`required|longchar|businessTopoInstNames|emptySetName|setNameMap|setNameLen|splitMaxLength:100
          ,${$t('超过限制，一次最多支持创建n个', { n: 100 })}`"
          v-model="setName"
          :rows="rows"
          :placeholder="$t('集群多个创建提示')"
          @keydown="handleKeydown"
          @paste="handlePaste">
        </bk-input>
        <span class="form-error" v-if="errors.has('setName')">{{errors.first('setName')}}</span>
      </div>
    </div>
    <div class="node-create-options">
      <bk-button theme="primary" class="mr10" v-test-id="'createSetSave'"
        :disabled="$loading() || errors.any()"
        @click="handleCreateSet">
        {{$t('提交')}}
      </bk-button>
      <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
    </div>
  </div>
</template>

<script>
  import has from 'has'
  import { MENU_BUSINESS_SET_TEMPLATE } from '@/dictionary/menu-symbol'
  export default {
    props: {
      parentNode: {
        type: Object,
        required: true
      }
    },
    data() {
      return {
        withTemplate: true,
        setTemplate: '',
        setName: '',
        rows: 1,
        setTemplateList: [],
        request: {
          setTemplate: Symbol('setTemplate')
        }
      }
    },
    computed: {
      topoPath() {
        const nodePath = [...this.parentNode.parents, this.parentNode]
        return nodePath.map(node => node.data.bk_inst_name).join('/')
      },
      business() {
        return this.$store.getters['objectBiz/bizId']
      },
      setTemplateMap() {
        return this.$store.state.businessHost.setTemplateMap
      }
    },
    watch: {
      withTemplate(value) {
        if (value) {
          this.setTemplate = this.setTemplateList[0] ? this.setTemplateList[0].id : ''
        } else {
          this.setTemplate = 0
        }
      }
    },
    created() {
      this.getSetTemplates()
    },
    methods: {
      setRows() {
        setTimeout(() => {
          const rows = this.setName.split('\n').length
          this.rows = Math.min(3, Math.max(rows, 1))
        })
      },
      handleKeydown(value, keyEvent) {
        if (['Enter', 'NumpadEnter'].includes(keyEvent.code)) {
          this.rows = Math.min(this.rows + 1, 3)
        } else if (keyEvent.code === 'Backspace') {
          this.setRows()
        }
      },
      handlePaste() {
        this.setRows()
      },
      async getSetTemplates() {
        if (has(this.setTemplateMap, this.business)) {
          this.setTemplateList = this.setTemplateMap[this.business]
        } else {
          try {
            const data = await this.$store.dispatch('setTemplate/getSetTemplates', {
              bizId: this.business,
              params: { page: { sort: '-last_time' } },
              config: {
                requestId: this.request.setTemplate
              }
            })
            const list = (data.info || []).map(template => ({ ...template.set_template }))
            this.setTemplateList = list
            this.$store.commit('businessHost/setSetTemplate', {
              id: this.business,
              templates: list
            })
          } catch (e) {
            console.error(e)
          }
        }
        this.setTemplate = this.setTemplateList[0] ? this.setTemplateList[0].id : ''
      },
      handleCreateSet() {
        this.$validator.validateAll().then((isValid) => {
          if (isValid) {
            const nameList = this.setName.split('\n').filter(name => name.trim().length)
              .map(name => name.trim())
            const sets = nameList.map(name => ({
              set_template_id: this.setTemplate,
              bk_set_name: name
            }))
            this.$emit('submit', {
              set_template_id: this.setTemplate,
              sets
            })
          }
        })
      },
      handleAddTemplate() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SET_TEMPLATE,
          params: {
            bizId: this.business
          }
        })
      },
      handleCancel() {
        this.$emit('cancel')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .node-create-layout {
        position: relative;
    }
    .node-create-title {
        margin-top: -15px;
        padding: 0 26px;
        line-height: 30px;
        font-size: 24px;
        color: #444444;
        font-weight: normal;
    }
    .node-create-path {
        padding: 14px 26px 0;
        margin: 0 0 -5px 0;
        font-size: 12px;
        color: #63656E;
        @include ellipsis;
    }
    .node-create-form {
        padding: 20px 26px 32px;
    }
    .form-item {
        margin: 15px 0 0 0;
        position: relative;
        .bk-form-radio {
            display: inline-block;
            margin-right: 70px;
            /deep/ input[type=radio] {
                margin-top: 2px;
            }
        }
        label {
            display: block;
            padding: 0 0 10px;
            line-height: 19px;
            font-size: 14px;
            color: #63656E;
            > span {
                color: #979BA5;
                font-size: 12px;
            }
            .red-star {
              color: #f00;
              font-size: 14px;
            }
        }
        .form-error {
            position: absolute;
            top: 100%;
            left: 0;
            font-size: 12px;
            color: $cmdbDangerColor;
            &.second-class {
                left: 270px;
            }
        }
        .form-textarea {
            /deep/ textarea {
                min-height: auto !important;
                line-height: 22px;
                @include scrollbar-y(6px);
            }
        }
    }
    .add-template {
        width: 20%;
        line-height: 38px;
        cursor: pointer;
        color: #63656E;
        font-size: 12px;
        .icon-plus-circle {
            margin-top: -2px;
            font-size: 14px;
            color: #979BA5;
        }
    }
    .node-create-options {
        padding: 9px 20px;
        border-top: 1px solid $cmdbBorderColor;
        text-align: right;
        background-color: #FAFBFD;
        font-size: 0;
    }
</style>
