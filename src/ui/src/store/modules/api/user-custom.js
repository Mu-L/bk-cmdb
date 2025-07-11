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
import Vue from 'vue'
import has from 'has'
import { MENU_RESOURCE_COLLECTION } from '@/dictionary/menu-symbol'
import { BUILTIN_MODEL_COLLECTION_KEYS } from '@/dictionary/model-constants.js'

const state = {
  usercustom: {},
  globalUsercustom: {}
}

const getters = {
  classifyNavigationKey: () => 'custom_classify_navigation',
  firstEntryKey: () => 'custom_first_entry',
  recentlyKey: () => 'custom_recently',
  usercustom: state => state.usercustom,
  getCustomData: state => (key, defaultData = null) => {
    if (has(state.usercustom, key)) {
      return state.usercustom[key]
    }
    return defaultData
  },
  resourceCollection: (state, getters, rootState, rootGetters) => {
    const collection = [...(state.usercustom[MENU_RESOURCE_COLLECTION] || [])]

    Object.keys(BUILTIN_MODEL_COLLECTION_KEYS).forEach((modelId) => {
      const collected = state.usercustom[BUILTIN_MODEL_COLLECTION_KEYS[modelId]] ?? true
      if (collected) {
        collection.unshift(modelId)
      }
    })

    const models = rootGetters['objectModelClassify/models']

    return collection.filter(modelId => models.some(model => model.bk_obj_id === modelId))
  }
}

const actions = {
  /**
     * 保存用户字段配置
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  saveUsercustom({ commit, state, dispatch }, usercustom = {}, config = {}) {
    return $http.post('usercustom', usercustom, { cancelWhenRouteChange: false, ...config }).then(() => {
      $http.cancelCache('searchUserCustom')
      commit('setUsercustom', usercustom)
      return state.usercustom
    })
  },

  /**
     * 获取用户字段配置
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
  searchUsercustom({ commit, state, dispatch }, { config }) {
    const mergedConfig = Object.assign({
      requestId: 'searchUserCustom'
    }, config)
    return $http.post('usercustom/user/search', {}, mergedConfig).then((usercustom) => {
      commit('setUsercustom', usercustom)
      return usercustom
    })
  },

  /**
     * 获取默认字段配置
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
  getUserDefaultCustom({ commit, state, dispatch }) {
    return $http.post('usercustom/default/search')
  },

  setRencentlyData({ commit, state, dispatch }, { id }) {
    const usercustomData = state.usercustom.recently_models || []
    const isExist = usercustomData.some(target => target === id)
    let newUsercustomData = [...usercustomData]
    if (isExist) {
      newUsercustomData = newUsercustomData.filter(target => target !== id)
    }
    newUsercustomData.unshift(id)
    dispatch('saveUsercustom', {
      recently_models: newUsercustomData
    })
  },

  saveGlobalUsercustom({ commit }, { objId, params, config }) {
    return $http.post(`usercustom/default/model/${objId}`, params, config).then((data) => {
      commit('setGlobalUsercustom', {
        [`${objId}_global_custom_table_columns`]: params.global_custom_table_columns
      })
      return data
    })
  },

  getGlobalUsercustom({ commit }, { config }) {
    const mergedConfig = Object.assign({
      requestId: 'getGlobalUsercustom'
    }, config)
    return $http.post('usercustom/default/model', {}, mergedConfig).then((usercustom) => {
      commit('setGlobalUsercustom', usercustom)
      return usercustom
    })
  }
}

const mutations = {
  setUsercustom(state, usercustom = {}) {
    // eslint-disable-next-line no-restricted-syntax
    for (const key in usercustom) {
      Vue.set(state.usercustom, key, usercustom[key])
    }
  },
  setGlobalUsercustom(state, globalUsercustom = {}) {
    // eslint-disable-next-line no-restricted-syntax
    for (const key in globalUsercustom) {
      Vue.set(state.globalUsercustom, key, globalUsercustom[key])
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
