/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

/* eslint-disable no-unused-vars */

import $http from '@/api'

const state = {}

const getters = {}

const actions = {
  /**
     * 获取服务分类
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  searchServiceCategory({ commit, state, dispatch, rootGetters }, { params, config }) {
    return $http.post('findmany/proc/service_category/with_statistics', params, config)
  },
  searchServiceCategoryWithoutAmout(context, { params, config }) {
    return $http.post('findmany/proc/service_category', params, config)
  },
  /**
     * 创建服务分类
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  createServiceCategory({ commit, state, dispatch, rootGetters }, { params, config }) {
    return $http.post('create/proc/service_category', params, config)
  },
  /**
     * 更新服务分类
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  updateServiceCategory({ commit, state, dispatch, rootGetters }, { params, config }) {
    return $http.put('update/proc/service_category', params, config)
  },
  /**
     * 创建服务分类
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  deleteServiceCategory({ commit, state, dispatch, rootGetters }, { params, config }) {
    return $http.delete('delete/proc/service_category', params, config)
  }
}

const mutations = {}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
