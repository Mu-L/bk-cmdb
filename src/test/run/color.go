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

package run

import "fmt"

const (
	escape      = "\x1b"
	reset       = "0"
	lightRed    = "91"
	lightGreen  = "92"
	lightYellow = "93"
	lightBlue   = "34"
)

func setBoldColor(color string, s interface{}) string {
	var ret string
	// format
	ret = fmt.Sprintf("%s[1;%sm%v", escape, color, s)
	// unformat
	ret += fmt.Sprintf("%s[%sm", escape, reset)
	return ret

}

// SetGreen TODO
func SetGreen(v interface{}) string {
	return setBoldColor(lightGreen, v)
}

// SetRed TODO
func SetRed(v interface{}) string {
	return setBoldColor(lightRed, v)
}

// SetYellow TODO
func SetYellow(v interface{}) string {
	return setBoldColor(lightYellow, v)
}

// SetBlue TODO
func SetBlue(v interface{}) string {
	return setBoldColor(lightBlue, v)
}
