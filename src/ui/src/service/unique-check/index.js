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

import http from '@/api'
import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'

export const find = async (modelId, config) => {
  try {
    const result = await http.post(`find/objectunique/object/${modelId}`, [], config)
    return result
  } catch (error) {
    console.error(error)
    return []
  }
}

export const findMany = async (models, config) => Promise.all(models.map(modelId => find(modelId, config)))

export const findBizSet = () => find(BUILTIN_MODELS.BUSINESS_SET)

export default {
  find,
  findMany,
  findBizSet
}
