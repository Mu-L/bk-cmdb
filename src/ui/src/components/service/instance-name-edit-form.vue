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
  <div class="instance-name-edit-form">
    <bk-input size="small" font-size="normal" :style="{ width: width ? `${width}px` : '' }" ref="input"
      :placeholder="localPlaceholder"
      v-model.trim="localValue"
      @enter="handleConfirm">
    </bk-input>
    <bk-button class="btn" theme="primary" :disabled="disabled" text @click.stop="handleConfirm">
      {{$t('确定')}}
    </bk-button>
    <span class="divider">|</span>
    <bk-button class="btn" theme="primary" text @click.stop="handleCancel">{{$t('取消')}}</bk-button>
  </div>
</template>

<script>
  export default {
    props: {
      value: {
        type: String,
        default: ''
      },
      placeholder: {
        type: String,
        default: ''
      },
      width: {
        type: Number
      }
    },
    data() {
      return {
        localValue: '',
        localPlaceholder: this.placeholder || this.$t('请输入实例名称')
      }
    },
    computed: {
      disabled() {
        return this.value === this.localValue || !this.localValue.length
      }
    },
    watch: {
      value: {
        handler(value) {
          this.localValue = value || ''
        },
        immediate: true
      }
    },
    methods: {
      focus() {
        this.$refs.input && this.$refs.input.focus()
      },
      handleConfirm() {
        this.$emit('confirm', this.localValue)
      },
      handleCancel() {
        this.$emit('cancel')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .instance-name-edit-form {
        display: flex;
        align-items: center;
        .btn {
            flex: none;
            line-height: normal;
            font-size: 12px;
            margin: 0 6px;
            &.is-disabled {
                &.bk-button-text {
                    color: $textDisabledColor;
                }
            }
        }
        .divider {
            color: $textDisabledColor;
        }
    }
</style>
