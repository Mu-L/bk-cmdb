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

package operator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// IamIDKey TODO
	IamIDKey = "id"
	// IamPathKey TODO
	IamPathKey = "_bk_iam_path_"
)

// Policy TODO
type Policy struct {
	Operator OperType `json:"op"`
	// Element is a pointer interface point to the implements struct,
	// which should be one of Content or FieldValue.
	Element
}

// UnmarshalJSON TODO
func (p *Policy) UnmarshalJSON(i []byte) error {
	if string(i) == "{}" {
		return nil
	}

	broker := new(policyBroker)
	err := json.Unmarshal(i, broker)
	if err != nil {
		return err
	}

	p.Operator = broker.Operator

	if broker.Operator == And || broker.Operator == Or {
		content := new(Content)
		if err := json.Unmarshal(broker.Content, &content.Content); err != nil {
			return err
		}
		p.Element = content
		return nil
	}

	if broker.Operator == In || broker.Operator == Nin {
		to := make([]interface{}, 0)
		if err := json.Unmarshal(broker.Value, &to); err != nil {
			return err
		}

		p.Element = &FieldValue{
			Field: broker.Field,
			Value: to,
		}

	} else {
		to := new(interface{})
		if err := json.Unmarshal(broker.Value, &to); err != nil {
			return err
		}

		p.Element = &FieldValue{
			Field: broker.Field,
			Value: *to,
		}
	}

	return nil
}

type policyBroker struct {
	Operator OperType        `json:"op"`
	Content  json.RawMessage `json:"content"`
	Field    Field           `json:"field"`
	Value    json.RawMessage `json:"value"`
}

// MarshalJSON is used to marshal the policy to the standard
// iam policy protocol, which is not correspond to the struct
// we defined here.
// Note: when you marshal the policy, the policy must be a pointer,
// otherwise, the marshaled json struct is wrong.
func (p *Policy) MarshalJSON() ([]byte, error) {
	js, err := json.Marshal(p.Element)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	buf.WriteString(`{"op":"`)
	buf.WriteString(string(p.Operator))
	buf.WriteString(`",`)
	buf.Write(js[1 : len(js)-1])
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// Element TODO
type Element interface {
	EleName() string
}

// Content TODO
type Content struct {
	// Content is only exist when OperType is "And" or "OR"
	Content []*Policy `json:"content"`
}

// EleName TODO
func (e *Content) EleName() string {
	return "content"
}

// FieldValue TODO
type FieldValue struct {
	// Field and Value is only exist when OperType is not
	// one of "And" or "OR"
	Field Field       `json:"field"`
	Value interface{} `json:"value"`
}

// EleName TODO
func (f *FieldValue) EleName() string {
	return "field_value"
}

// Field TODO
type Field struct {
	Resource  string
	Attribute string
}

// UnmarshalJSON TODO
func (f *Field) UnmarshalJSON(i []byte) error {
	if string(i) == "\"\"" {
		f.Attribute = ""
		f.Resource = ""
		return nil
	}
	index := bytes.IndexByte(i, '.')
	if index < 0 {
		return errors.New("invalid \"field\"")
	}

	f.Resource = string(bytes.TrimLeft(i[:index], "\""))
	f.Attribute = string(bytes.TrimRight(i[index+1:], "\""))

	if f.Resource == "" || f.Attribute == "" {
		return errors.New("invalid \"field\"")
	}

	return nil
}

// MarshalJSON TODO
func (f *Field) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s.%s\"", f.Resource, f.Attribute)), nil
}
