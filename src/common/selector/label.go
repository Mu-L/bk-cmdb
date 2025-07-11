// Package selector TODO
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
package selector

import (
	"fmt"
	"regexp"
)

// Labels TODO
type Labels map[string]string

var (
	// LabelNGKeyRule TODO
	LabelNGKeyRule = regexp.MustCompile(`^[a-zA-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`)
	// LabelNGValueRule TODO
	LabelNGValueRule = regexp.MustCompile(`^[a-z0-9A-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`)
)

// Validate TODO
func (lng Labels) Validate() (string, error) {
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set
	// https://www.replex.io/blog/9-best-practices-and-examples-for-working-with-kubernetes-labels
	for key, value := range lng {
		// validate key
		if LabelNGKeyRule.MatchString(key) == false {
			return key, fmt.Errorf("key: %s format error", key)
		}
		if len(key) >= 64 {
			return key, fmt.Errorf("key: %s exceed max length 63", key)
		}

		// validate value
		field := fmt.Sprintf("%s:%s", key, value)
		if LabelNGValueRule.MatchString(value) == false {
			return field, fmt.Errorf("value: %s format error", field)
		}
		if len(value) >= 64 {
			return field, fmt.Errorf("value: %s exceed max length 63", field)
		}
	}
	return "", nil
}

// AddLabel TODO
func (lng Labels) AddLabel(l Labels) {
	for key, value := range l {
		lng[key] = value
	}
}

// RemoveLabel TODO
func (lng Labels) RemoveLabel(keys []string) {
	for _, key := range keys {
		delete(lng, key)
	}
}

// LabelInstance TODO
type LabelInstance struct {
	Labels Labels `bson:"labels" json:"labels"`
}
