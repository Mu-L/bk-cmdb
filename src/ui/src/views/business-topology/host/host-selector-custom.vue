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
  <div class="custom-layout">
    <textarea class="ip-input"
      :placeholder="$t('请输入IP，换行分隔')"
      :class="{ 'has-error': invalidList.length }"
      v-model="value"
      @focus="handleFocus">
        </textarea>
    <p class="ip-error" v-show="invalidList.length">{{$t('IP不正确或主机不存在')}}</p>
    <bk-button class="ip-confirm" outline theme="primary" @click="handleConfirm">{{$t('添加至列表')}}</bk-button>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import hostSearchService from '@/service/host/search'

  export default {
    data() {
      return {
        value: '',
        validList: [],
        invalidList: []
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('businessHost', ['getDefaultSearchCondition'])
    },
    activated() {
      this.value = ''
      this.validList = []
      this.invalidList = []
    },
    methods: {
      handleFocus() {
        this.invalidList = []
      },
      async handleConfirm() {
        try {
          await this.validateList()
          if (this.validList.length) {
            const result = await hostSearchService.getBizHosts({
              params: this.getParams(),
              config: {
                requestId: this.$parent.request.host
              }
            })
            this.validList.forEach((ip) => {
              const targetHost = result.info.find(target => target.host.bk_host_innerip === ip)
              if (targetHost) {
                this.$parent.handleSelect(targetHost)
              } else {
                this.invalidList.push(ip)
              }
            })
          }
          this.value = this.invalidList.join('\n')
        } catch (e) {
          console.error(e)
        }
      },
      async validateList() {
        const list = [...new Set(this.value.split('\n').map(ip => ip.trim())
          .filter(ip => ip.length))]
        const validateQueue = []
        list.forEach((ip) => {
          validateQueue.push(this.$validator.verify(ip, 'ip'))
        })
        const results = await Promise.all(validateQueue)
        const validList = []
        const invalidList = []
        results.forEach(({ valid }, index) => {
          if (valid) {
            validList.push(list[index])
          } else {
            invalidList.push(list[index])
          }
        })
        this.validList = validList
        this.invalidList = invalidList
      },
      getParams() {
        return {
          bk_biz_id: this.bizId,
          condition: this.getDefaultSearchCondition(),
          ip: { data: this.validList, exact: 1, flag: 'bk_host_innerip' }
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .custom-layout {
        padding: 0 20px;
        position: relative;
    }
    .ip-input {
        display: block;
        width: 100%;
        height: calc(100% - 60px);
        padding: 5px 10px;
        font-size: 12px;
        line-height: 20px;
        background-color: #FFF;
        border-radius:2px;
        border:1px solid #C4C6CC;
        cursor: text;
        outline: 0;
        resize: none;
        @include scrollbar;
        &.has-error {
            color: $dangerColor;
            text-decoration: underline;
        }
    }
    .ip-error {
        position: absolute;
        bottom: 45px;
        left: 20px;
        line-height: 16px;
        font-size: 12px;
        color: $dangerColor;
    }
    .ip-confirm {
        display: block;
        width: 100%;
        margin: 15px 0;
    }
</style>
