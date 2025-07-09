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

package cmd

import (
	"fmt"
)

// WithRedColor TODO
func WithRedColor(str string) string {
	return fmt.Sprintf("%c[1;40;31m>> %s %c[0m\n", 0x1B, str, 0x1B)
}

// WithGreenColor TODO
func WithGreenColor(str string) string {
	return fmt.Sprintf("%c[1;40;32m>> %s %c[0m\n", 0x1B, str, 0x1B)
}

// WithBlueColor TODO
func WithBlueColor(str string) string {
	return fmt.Sprintf("%c[1;40;34m>> %s %c[0m\n", 0x1B, str, 0x1B)
}
