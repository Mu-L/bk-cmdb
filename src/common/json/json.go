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

// Package json TODO
package json

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var iteratorJson = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	UseNumber:              true,
}.Froze()

// MarshalToString TODO
func MarshalToString(v interface{}) (string, error) {
	return iteratorJson.MarshalToString(v)
}

// Marshal TODO
func Marshal(v interface{}) ([]byte, error) {
	return iteratorJson.Marshal(v)
}

// MarshalIndent TODO
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return iteratorJson.MarshalIndent(v, prefix, indent)
}

// UnmarshalFromString TODO
func UnmarshalFromString(str string, v interface{}) error {
	return iteratorJson.UnmarshalFromString(str, v)
}

// Unmarshal TODO
func Unmarshal(data []byte, v interface{}) error {
	return iteratorJson.Unmarshal(data, v)
}

// UnmarshalArray TODO
func UnmarshalArray(items []string, result interface{}) error {
	strArrJSON := "[" + strings.Join(items, ",") + "]"
	return iteratorJson.Unmarshal([]byte(strArrJSON), result)
}
