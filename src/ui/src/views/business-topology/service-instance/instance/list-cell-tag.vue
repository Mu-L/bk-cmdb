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
  <div class="tag-wrapper">
    <ul class="tag-list" ref="list" v-if="tags.length">
      <li class="tag-item"
        v-for="(tag, index) in tags"
        :key="index"
        :title="tag">
        {{tag}}
      </li>
      <li class="tag-item ellipsis" ref="ellipsis" v-show="ellipsisCount" @click.stop>+{{ellipsisCount}}</li>
    </ul>
    <span class="tag-empty" v-else>--</span>
    <cmdb-auth
      v-if="!readonly"
      tag="i"
      class="tag-edit icon-cc-edit"
      ref="editTrigger"
      :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
      @click.native.stop
      @click="handleEditLabel">
    </cmdb-auth>
  </div>
</template>

<script>
  import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
  import { mapGetters } from 'vuex'
  import throttle from 'lodash.throttle'
  import LabelDialog from './dialog/label-dialog.js'
  export default {
    name: 'list-cell-tag',
    props: {
      row: Object,
      readonly: {
        type: Boolean,
        default: false,
      }
    },
    data() {
      return {
        ellipsisCount: 0
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      tags() {
        const { labels } = this.row
        if (!labels) {
          return []
        }
        return Object.keys(labels).map(key => `${key} : ${labels[key]}`)
      }
    },
    watch: {
      tags() {
        this.handleResize()
      }
    },
    created() {
      this.scheduleResize = throttle(this.handleResize, 300)
    },
    mounted() {
      addResizeListener(this.$el, this.scheduleResize)
    },
    beforeDestroy() {
      removeResizeListener(this.$el, this.scheduleResize)
    },
    methods: {
      handleResize() {
        this.removeEllipsisTag()
        if (!this.tags.length) {
          this.updateEditPosition()
          return
        }
        this.$nextTick(() => {
          const items = Array.from(this.$refs.list.querySelectorAll('.tag-item'))
          const referenceItemIndex = items.findIndex((item, index) => {
            if (index === 0) {
              return false
            }
            const previousItem = items[index - 1]
            return previousItem.offsetTop !== item.offsetTop
          })
          if (referenceItemIndex > -1) {
            this.ellipsisCount = this.tags.length - referenceItemIndex
            this.insertEllipsisTag(items[referenceItemIndex], referenceItemIndex)
            this.$nextTick(this.doubleCheckEllipsisPosition)
          } else {
            this.removeEllipsisTag()
          }
          this.$nextTick(this.updateEditPosition)
        })
      },
      insertEllipsisTag(reference) {
        const { ellipsis } = this.$refs
        this.$refs.list.insertBefore(ellipsis, reference)
      },
      doubleCheckEllipsisPosition() {
        const { ellipsis } = this.$refs
        const previous = ellipsis.previousElementSibling
        if (previous && ellipsis.offsetTop !== previous.offsetTop) {
          this.$refs.list.insertBefore(ellipsis, previous)
          this.ellipsisCount += 1
        }
        this.setEllipsisTips()
      },
      updateEditPosition() {
        if (this.readonly) return
        const { ellipsis } = this.$refs
        let lastItem = null
        if (ellipsis && ellipsis.previousElementSibling) {
          lastItem = ellipsis
        } else if (this.tags.length) {
          const tagItems = Array.from(this.$refs.list.querySelectorAll('.tag-item')).filter(el => el.clientHeight > 0)
          lastItem = tagItems[tagItems.length - 1]
        }
        this.$refs.editTrigger.$el.style.left = lastItem ? `${lastItem.offsetLeft + lastItem.offsetWidth + 10}px` : 0
      },
      setEllipsisTips() {
        const { ellipsis } = this.$refs
        const tips = this.getTipsInstance()
        const tipsNode = this.$refs.list.cloneNode(false)
        let loopItem = ellipsis
        while (loopItem) {
          const nextItem = loopItem.nextElementSibling
          if (nextItem && nextItem.classList.contains('tag-item')) {
            tipsNode.appendChild(nextItem.cloneNode(true))
            loopItem = nextItem
          } else {
            loopItem = null
          }
        }
        tips.setContent(tipsNode)
      },
      getTipsInstance() {
        if (!this.tips) {
          this.tips = this.$bkPopover(this.$refs.ellipsis, {
            allowHTML: true,
            placement: 'top',
            arrow: true,
            theme: 'light',
            interactive: true
          })
        }
        return this.tips
      },
      removeEllipsisTag() {
        try {
          this.$refs.list.removeChild(this.$refs.ellipsis)
        } catch (e) {}
      },
      handleEditLabel() {
        LabelDialog.show({
          serviceInstance: this.row,
          updateCallback: (labels) => {
            const newLabels = {}
            labels.forEach((label) => {
              newLabels[label.key] = label.value
            })
            this.$emit('update-labels', newLabels)
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .tag-wrapper {
        position: relative;
        padding-right: 10px;
        height: 22px;
        .tag-edit {
            position: absolute;
            left: 0;
            top: 0;
            flex: 22px 0 0;
            height: 22px;
            font-size: 12px;
            line-height: 22px;
            color: $primaryColor;
            cursor: pointer;
            visibility: hidden;
            &:hover {
                opacity: .8;
            }
            &.disabled {
                color: $textDisabledColor;
            }
        }
    }
    .tag-list {
        flex: 1;
        height: 22px;
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        overflow: hidden;
        font-size: 12px;
        .tag-item {
            display: inline-block;
            max-width: 80px;
            padding: 0 6px;
            border-radius: 2px;
            line-height: 22px;
            color: $textColor;
            background-color: #f0f1f5;
            cursor: default;
            @include ellipsis;
            & ~ .tag-item {
                margin-left: 6px;
            }
            &.ellipsis {
                height: 22px;
                text-align: center;
                & ~ .tag-item {
                    display: none;
                }
            }
        }
    }
</style>
