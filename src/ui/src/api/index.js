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

import Axios from 'axios'
import md5 from 'md5'
import has from 'has'
import { $error } from '@/magicbox'
import i18n, { language } from '@/i18n'
import { showLoginModal } from '@/utils/login-helper'
import CachedPromise from './_cached-promise'
import RequestQueue from './_request-queue'
import customHeaders from './custom-header'

// axios实例
const axiosInstance = Axios.create({
  baseURL: window.API_PREFIX,
  xsrfCookieName: 'data_csrftoken',
  xsrfHeaderName: 'X-CSRFToken',
  withCredentials: true
})

// axios实例拦截器
axiosInstance.interceptors.request.use(
  (config) => {
    for (const [key, value] of Object.entries(customHeaders())) {
      config.headers[key] = value
    }
    return config
  },
  error => Promise.reject(error)
)

axiosInstance.interceptors.response.use(
  response => response,
  error => Promise.reject(error)
)

const $http = {
  queue: new RequestQueue(),
  cache: new CachedPromise(),
  cancelRequest: requestId => $http.queue.cancel(requestId),
  cancelCache: requestId => $http.cache.delete(requestId),
  cancel: requestId => Promise.all([$http.cancelRequest(requestId), $http.cancelCache(requestId)]),
  setHeader: (key, value) => {
    axiosInstance.defaults.headers[key] = value
  },
  deleteHeader: (key) => {
    delete axiosInstance.defaults.headers[key]
  },
  download
}

const methodsWithoutData = ['delete', 'get', 'head', 'options']
const methodsWithData = ['post', 'put', 'patch']
const allMethods = [...methodsWithoutData, ...methodsWithData]

// 在自定义对象$http上添加各请求方法
allMethods.forEach((method) => {
  Object.defineProperty($http, method, {
    get() {
      return getRequest(method)
    }
  })
})

/**
 * 获取http不同请求方式对应的函数
 * @param {method} http method 与 axios实例中的method保持一致
 * @return {Function} 实际调用的请求函数
 */
function getRequest(method) {
  if (methodsWithData.includes(method)) {
    return (url, data, config) => getPromise(method, url, data, config)
  }
  return (url, config) => getPromise(method, url, null, config)
}

/**
 * 实际发起http请求的函数，根据配置调用缓存的promise或者发起新的请求
 * @param {method} http method 与 axios实例中的method保持一致
 * @param {url} 请求地址
 * @param {data} 需要传递的数据, 仅 post/put/patch 三种请求方式可用
 * @param {userConfig} 用户配置，包含axios的配置与本系统自定义配置
 * @return {Promise} 本次http请求的Promise
 */
async function getPromise(method, url, data, userConfig = {}) {
  const config = initConfig(method, url, userConfig)
  let promise
  if (config.cancelPrevious) {
    await $http.cancel(config.requestId)
  }
  if (config.clearCache) {
    $http.cache.delete(config.requestId)
  } else {
    promise = $http.cache.get(config.requestId)
  }
  if (config.fromCache && promise) {
    return promise
  }
  promise = new Promise((resolve, reject) => {
    const axiosRequest = methodsWithData.includes(method)
      ? axiosInstance[method](url, data, config)
      : axiosInstance[method](url, config)
    axiosRequest.then((response) => {
      Object.assign(config, response.config)
      handleResponse({ config, response, resolve, reject })
    }).catch((error) => {
      Object.assign(config, error.config)
      reject(error)
    })
  }).catch(error => handleReject(error, config))
    .finally(() => {
      $http.queue.delete(config.requestId, config.requestSymbol)
    })
  // 添加请求队列
  $http.queue.set(config)
  // 添加请求缓存
  $http.cache.set(config.requestId, promise, config)
  return promise
}

/**
 * 处理http请求成功结果
 * @param {config} 请求配置
 * @param {response} cgi原始返回数据
 * @param {resolve} promise完成函数
 * @param {reject} promise拒绝函数
 * @return
 */
const PermissionCode = 9900403
const TokenInvalidCode = 1306000
function handleResponse({ config, response, resolve, reject }) {
  const transformedResponse = response.data
  const { bk_error_msg: message, bk_error_code: code, permission } = transformedResponse
  if (code === PermissionCode) {
    config.globalPermission && popupPermissionModal(permission)
    return reject({ message, permission, code })
  }
  if (code === TokenInvalidCode) {
    return reject({ message, permission, code })
  }
  if (!transformedResponse.result && config.globalError) {
    reject({ message })
    return
  }
  if (config.originalResponse) {
    resolve(response)
    return
  }
  if (config.transformData) {
    resolve(transformedResponse.data)
    return
  }
  resolve(transformedResponse)
}

/**
 * 处理http请求失败结果
 * @param {error} Error 对象
 * @param {config} 请求配置
 * @return Promise.reject
 */
function handleReject(error, config) {
  if (error.code && error.code === PermissionCode) {
    return Promise.reject(error)
  }

  if (error.code && error.code === TokenInvalidCode) {
    showLoginModal()
  }

  if (Axios.isCancel(error)) {
    return Promise.reject(error)
  }
  if (error.response) {
    const { status, data } = error.response
    const nextError = { message: error.message, status }
    if (status === 401) {
      showLoginModal()
    } else if (data && data.bk_error_msg) {
      nextError.message = data.bk_error_msg
    } else if (status === 403) {
      nextError.message = language === 'en' ? 'You don\'t have permission.' : '无权限操作'
    } else if (status === 500) {
      nextError.message = language === 'en' ? 'System error, please contact developers.' : '系统出现异常, 请记录下错误场景并与开发人员联系, 谢谢!'
    }
    config.globalError && status !== 401 && $error(nextError.message)
    return Promise.reject(nextError)
  } if (error.message === 'Network Error') {
    $error(i18n.t('资源请求失败提示'))
  } else {
    config.globalError && $error(error.message)
  }
  return Promise.reject(error)
}

function popupPermissionModal(permission = []) {
  window.permissionModal && window.permissionModal.show(permission)
}

/**
 * 初始化本系统http请求的各项配置
 * @param {method} http method 与 axios实例中的method保持一致
 * @param {url} 请求地址, 结合method 生成md5 requestId
 * @param {userConfig} 用户配置，包含axios的配置与本系统自定义配置
 * @return {Promise} 本次http请求的Promise
 */

function initConfig(method, url, userConfig) {
  if (has(userConfig, 'requestGroup')) {
    userConfig.requestGroup = userConfig.requestGroup instanceof Array
      ? userConfig.requestGroup
      : [userConfig.requestGroup]
  }
  const defaultConfig = {
    ...getCancelToken(),
    // http请求默认id
    requestId: md5(method + url),
    requestGroup: [],
    requestSymbol: Symbol('requestSymbol'),
    // 是否全局捕获异常
    globalError: true,
    // 是否直接复用缓存的请求
    fromCache: false,
    // 是否在请求发起前清楚缓存
    clearCache: false,
    // 响应结果是否返回原始数据
    originalResponse: false,
    // 转换返回数据，仅返回data对象
    transformData: true,
    // 当路由变更时取消请求
    cancelWhenRouteChange: false,
    // 取消上次请求
    cancelPrevious: false,
    // 是否全局捕获权限异常
    globalPermission: true
  }
  return Object.assign(defaultConfig, userConfig)
}

/**
 * 生成http请求的cancelToken，用于取消尚未完成的请求
 * @return {Object}
 *      cancelToken: axios实例使用的cancelToken
 *      cancelExcutor: 取消http请求的可执行函数
 */
function getCancelToken() {
  let cancelExcutor
  const cancelToken = new Axios.CancelToken((excutor) => {
    cancelExcutor = excutor
  })
  return {
    cancelToken,
    cancelExcutor
  }
}

async function download(options = {}) {
  const { url, method = 'post', data, name } = options
  const config = Object.assign({
    globalError: false,
    originalResponse: true,
    responseType: 'blob'
  }, options.config)
  if (!url) {
    const error = new Error('Empty download url')
    $error(error.message)
    return Promise.reject(error)
  }
  let promise
  if (methodsWithData.includes(method)) {
    promise = $http[method](url, data, config)
  } else {
    promise = $http[method](url, config)
  }
  try {
    const response = await promise
    if (response.data.type.indexOf('application') === -1) {
      const text = await new Response(response.data).text()
      throw new Error(JSON.parse(text).bk_error_msg)
    }
    const disposition = response.headers['content-disposition']
    const fileName = name || disposition.substring(disposition.indexOf('filename') + 9)
    const downloadUrl = window.URL.createObjectURL(new Blob([response.data], {
      type: response.headers['content-type']
    }))
    const link = document.createElement('a')
    link.style.display = 'none'
    link.href = downloadUrl
    link.setAttribute('download', fileName)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    return Promise.resolve(response)
  } catch (error) {
    if (Axios.isCancel(error)) {
      return Promise.reject(error)
    }
    $error(error.message)
    return Promise.reject(error)
  }
}

export * from './jsonp'

export const useHttp = () => $http

export default $http
