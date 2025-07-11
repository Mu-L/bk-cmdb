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
  <div class="property-filter clearfix">
    <cmdb-selector class="property-selector fl" style="width: 135px;"
      :list="filteredProperties"
      setting-key="bk_property_id"
      display-key="bk_property_name"
      v-model="localSelected.id"
      @on-selected="handlePropertySelected">
    </cmdb-selector>
    <cmdb-selector class="operator-selector fl" style="width: 135px;"
      :list="operatorOptions"
      setting-key="value"
      display-key="label"
      v-model="localSelected.operator"
      @on-selected="handleOperatorSelected">
    </cmdb-selector>
    <div class="property-value fl" style="width: 315px;"
      v-if="Object.keys(selectedProperty).length">
      <component
        class="search-form-el"
        :is="`cmdb-search-${selectedProperty['bk_property_type']}`"
        :options="selectedProperty.option || []"
        :placeholder="$t('请输入关键字')"
        :property="selectedProperty"
        v-model.trim="localSelected.value">
      </component>
    </div>
  </div>
</template>
<script>
  import { mapGetters, mapActions } from 'vuex'
  import has from 'has'
  import { QUERY_OPERATOR } from '@/utils/query-builder-operator'

  export default {
    props: {
      objId: {
        type: String,
        required: true
      },
      excludeType: {
        type: Array,
        default() {
          return []
        }
      },
      excludeId: {
        type: Array,
        default() {
          return []
        }
      }
    },
    data() {
      return {
        localSelected: {
          id: '',
          operator: '',
          value: ''
        },
        filteredProperties: [],
        propertyOperator: {
          default: ['$eq', '$ne'],
          singlechar: ['$regex', '$eq', '$ne'],
          longchar: ['$regex', '$eq', '$ne'],
          objuser: ['$in', '$nin'],
          singleasst: ['$regex', '$eq', '$ne'],
          multiasst: ['$regex', '$eq', '$ne'],
          list: ['$in', '$nin'],
          timezone: ['$in', '$nin'],
          enummulti: [QUERY_OPERATOR.IN, QUERY_OPERATOR.NIN],
          enumquote: [QUERY_OPERATOR.IN, QUERY_OPERATOR.NIN],
          organization: [QUERY_OPERATOR.IN, QUERY_OPERATOR.NIN]
        },
        operatorLabel: {
          $nin: this.$t('不包含'),
          $in: this.$t('包含'),
          $regex: this.$t('包含'),
          $eq: this.$t('等于'),
          $ne: this.$t('不等于')
        }
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      selectedProperty() {
        // eslint-disable-next-line max-len
        return this.filteredProperties.find(({ bk_property_id: bkPropertyId }) => bkPropertyId === this.localSelected.id) || {}
      },
      operatorOptions() {
        if (this.selectedProperty) {
          if (['bk_host_innerip', 'bk_host_outerip'].includes(this.selectedProperty.bk_property_id) || this.objId === 'biz') {
            return [{ label: this.operatorLabel.$regex, value: '$regex' }]
          }
          const propertyType = this.selectedProperty.bk_property_type
          // eslint-disable-next-line max-len
          const propertyOperator = has(this.propertyOperator, propertyType) ? this.propertyOperator[propertyType] : this.propertyOperator.default
          return propertyOperator.map(operator => ({
            label: this.operatorLabel[operator],
            value: operator
          }))
        }
        return []
      }
    },
    watch: {
      filteredProperties(properties) {
        if (properties.length) {
          this.localSelected.id = properties[0].bk_property_id
          this.$emit('on-property-selected', properties[0].bk_property_id, properties[0])
        } else {
          this.localSelected.id = ''
          this.$emit('on-property-selected', '', null)
        }
      },
      operatorOptions(operatorOptions) {
        this.localSelected.operator = operatorOptions.length ? operatorOptions[0].value : ''
        this.$emit('handleOperatorSelected', this.localSelected.operator)
      },
      'localSelected.id'() {
        this.localSelected.value = ''
      },
      'localSelected.value'(value) {
        this.$emit('on-value-change', value)
      },
      async objId(objId) {
        try {
          const properties = await this.searchObjectAttribute({
            params: {
              bk_obj_id: objId,
              bk_supplier_account: this.supplierAccount
            },
            config: {
              requestId: `post_searchObjectAttribute_${objId}`,
              fromCache: true
            }
          })
          this.filteredProperties = properties.filter((property) => {
            const {
              bk_isapi: bkIsapi,
              bk_property_type: bkPropertyType,
              bk_property_id: bkPropertyId
            } = property
            return !bkIsapi && !this.excludeType.includes(bkPropertyType) && !this.excludeId.includes(bkPropertyId)
          })
        } catch (err) {
          console.error(err)
        }
      }
    },
    methods: {
      ...mapActions('objectModelProperty', ['searchObjectAttribute']),
      handlePropertySelected(value, data) {
        this.$emit('on-property-selected', value, data)
      },
      handleOperatorSelected(value, data) {
        this.$emit('on-operator-selected', value, data)
      },
      clearFilter() {
        this.localSelected.value = ''
      }
    }
  }
</script>
<style lang="scss" scoped>
    .property-selector{
        width: 135px;
    }
    .operator-selector{
        width: 135px;
        margin: 0 10px;
    }
    .property-value{
        width: 245px;
        position: relative;

        .search-form-el {
          width: 100%;
        }
    }
</style>
