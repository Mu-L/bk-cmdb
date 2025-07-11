/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package model

import (
	"strings"
)

const (
	mongoFieldNotAllowed string = "$."
	// The maximum length of <database>.<collection>) is 120 bytes, so MongoCollMaxLength is 120-len("cmdb.")
	mongoCollMaxLength int = 115
)

// satisfyCharLimit 字符格式限制检查
func satisfyCharLimit(name string) bool {
	if strings.ContainsAny(name, mongoFieldNotAllowed) {
		return false
	}
	return true
}

// SatisfyMongoCollLimit 考虑到后续按模型拆分collection，需要限制模型名的字符格式，以满足mongo对collection名的限制要求
// mongo collection名限制可见 https://docs.mongodb.com/manual/reference/limits/#Restriction-on-Collection-Names
func SatisfyMongoCollLimit(collName string) bool {
	if len(collName) > mongoCollMaxLength {
		return false
	}
	return satisfyCharLimit(collName)
}

// SatisfyMongoFieldLimit 考虑到后续按模型拆分collection，需要限制模型属性名的字符格式，以满足mongo对字段名的限制要求
// mongo field名限制可见 https://docs.mongodb.com/manual/reference/limits/#Restrictions-on-Field-Names
func SatisfyMongoFieldLimit(fieldName string) bool {
	return satisfyCharLimit(fieldName)
}
