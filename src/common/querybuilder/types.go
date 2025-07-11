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

package querybuilder

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/util"
)

const timeLayout = "2006-01-02"

// Rule TODO
type Rule interface {
	GetDeep() int
	Validate(option *RuleOption) (string, error)
	ToMgo() (mgoFilter map[string]interface{}, errKey string, err error)
	Match(matcher Matcher) bool
	// MatchAny if any of the rules matches the matcher, return true
	MatchAny(matcher Matcher) bool
	GetField() []string
}

// Condition TODO
// *************** define condition ************************
type Condition string

// Validate TODO
func (c Condition) Validate() error {
	if c == ConditionAnd || c == ConditionOr {
		return nil
	}
	return fmt.Errorf("unexpected condition: %s", c)
}

// ToMgo TODO
func (c Condition) ToMgo() (mgoOperator string, err error) {
	switch c {
	case ConditionOr:
		return common.BKDBOR, nil
	case ConditionAnd:
		return common.BKDBAND, nil
	default:
		return "", fmt.Errorf("unexpected operator %s", c)
	}
}

var (
	// ConditionAnd TODO
	ConditionAnd = Condition("AND")
	// ConditionOr TODO
	ConditionOr = Condition("OR")
)

// Operator TODO
// *************** define operator ************************
type Operator string

var (
	// OperatorEqual TODO
	OperatorEqual = Operator("equal")
	// OperatorNotEqual TODO
	OperatorNotEqual = Operator("not_equal")

	// OperatorIn TODO
	// set operator
	OperatorIn = Operator("in")
	// OperatorNotIn TODO
	OperatorNotIn = Operator("not_in")

	// OperatorLess TODO
	// numeric compare
	OperatorLess = Operator("less")
	// OperatorLessOrEqual TODO
	OperatorLessOrEqual = Operator("less_or_equal")
	// OperatorGreater TODO
	OperatorGreater = Operator("greater")
	// OperatorGreaterOrEqual TODO
	OperatorGreaterOrEqual = Operator("greater_or_equal")

	// OperatorDatetimeLess TODO
	// datetime operate only use for data type
	OperatorDatetimeLess = Operator("datetime_less")
	// OperatorDatetimeLessOrEqual TODO
	OperatorDatetimeLessOrEqual = Operator("datetime_less_or_equal")
	// OperatorDatetimeGreater TODO
	OperatorDatetimeGreater = Operator("datetime_greater")
	// OperatorDatetimeGreaterOrEqual TODO
	OperatorDatetimeGreaterOrEqual = Operator("datetime_greater_or_equal")

	// OperatorBeginsWith TODO
	// string operator
	OperatorBeginsWith = Operator("begins_with")
	// OperatorNotBeginsWith TODO
	OperatorNotBeginsWith = Operator("not_begins_with")
	// OperatorContains TODO
	OperatorContains = Operator("contains")
	// OperatorNotContains TODO
	OperatorNotContains = Operator("not_contains")
	// OperatorsEndsWith TODO
	OperatorsEndsWith = Operator("ends_with")
	// OperatorNotEndsWith TODO
	OperatorNotEndsWith = Operator("not_ends_with")

	// OperatorIsEmpty TODO
	// array operator
	OperatorIsEmpty = Operator("is_empty")
	// OperatorIsNotEmpty TODO
	OperatorIsNotEmpty = Operator("is_not_empty")

	// OperatorIsNull TODO
	// null check
	OperatorIsNull = Operator("is_null")
	// OperatorIsNotNull TODO
	OperatorIsNotNull = Operator("is_not_null")

	// OperatorExist TODO
	// exist check
	OperatorExist = Operator("exist")
	// OperatorNotExist TODO
	OperatorNotExist = Operator("not_exist")
)

// SupportOperators TODO
var SupportOperators = map[Operator]bool{
	OperatorEqual:    true,
	OperatorNotEqual: true,

	OperatorIn:    true,
	OperatorNotIn: true,

	OperatorLess:           true,
	OperatorLessOrEqual:    true,
	OperatorGreater:        true,
	OperatorGreaterOrEqual: true,

	OperatorDatetimeLess:           true,
	OperatorDatetimeLessOrEqual:    true,
	OperatorDatetimeGreater:        true,
	OperatorDatetimeGreaterOrEqual: true,

	OperatorBeginsWith:    true,
	OperatorNotBeginsWith: true,
	OperatorContains:      true,
	OperatorNotContains:   true,
	OperatorsEndsWith:     true,
	OperatorNotEndsWith:   true,

	OperatorIsEmpty:    true,
	OperatorIsNotEmpty: true,

	OperatorIsNull:    true,
	OperatorIsNotNull: true,

	OperatorExist:    true,
	OperatorNotExist: true,
}

// Validate TODO
func (op Operator) Validate() error {
	if support, ok := SupportOperators[op]; !support || !ok {
		return fmt.Errorf("unsupported operator: %s", op)
	}
	return nil
}

// AtomRule TODO
// *************** define rule ************************
type AtomRule struct {
	Field    string      `json:"field"`
	Operator Operator    `json:"operator"`
	Value    interface{} `json:"value"`
}

// GetDeep TODO
func (r AtomRule) GetDeep() int {
	return int(1)
}

// Validate TODO
func (r AtomRule) Validate(option *RuleOption) (string, error) {
	if err := r.Operator.Validate(); err != nil {
		return "operator", err
	}
	if err := r.validateField(); err != nil {
		return "field", err
	}
	if err := r.validateValue(option); err != nil {
		return "value", err
	}
	return "", nil
}

// Matcher TODO
type Matcher func(r AtomRule) bool

// Match TODO
func (r AtomRule) Match(matcher Matcher) bool {
	return matcher(r)
}

// MatchAny if any of the rules matches the matcher, return true
func (r AtomRule) MatchAny(matcher Matcher) bool {
	return matcher(r)
}

var (
	// ValidFieldPattern TODO
	// TODO: should we support dot field separator here?
	ValidFieldPattern = regexp.MustCompile(`^[a-zA-Z0-9][\d\w\-_.]*$`)
)

func (r AtomRule) validateField() error {
	if !ValidFieldPattern.MatchString(r.Field) {
		return fmt.Errorf("invalid field: %s", r.Field)
	}
	return nil
}

func (r AtomRule) validateValue(option *RuleOption) error {
	switch r.Operator {
	case OperatorEqual, OperatorNotEqual:
		return validateBasicType(r.Value)

	case OperatorIn, OperatorNotIn:
		return validateSliceOfBasicType(r.Value, option.NeedSameSliceElementType, option.MaxSliceElementsCount)

	case OperatorLess, OperatorLessOrEqual, OperatorGreater, OperatorGreaterOrEqual:
		return validateNumericType(r.Value)
	case OperatorDatetimeLess, OperatorDatetimeLessOrEqual, OperatorDatetimeGreater, OperatorDatetimeGreaterOrEqual:
		return validateDatetimeStringType(r.Value)
	case OperatorBeginsWith, OperatorNotBeginsWith, OperatorContains, OperatorNotContains, OperatorsEndsWith, OperatorNotEndsWith:
		return validateNotEmptyStringType(r.Value)
	case OperatorIsEmpty, OperatorIsNotEmpty:
		return nil
	case OperatorIsNull, OperatorIsNotNull:
		return nil
	case OperatorExist, OperatorNotExist:
		return nil
	default:
		return fmt.Errorf("unsupported operator: %s", r.Operator)
	}
}

// ToMgo generate mongo filter from rule
func (r AtomRule) ToMgo() (mgoFiler map[string]interface{}, key string, err error) {
	if key, err := r.Validate(&RuleOption{NeedSameSliceElementType: true}); err != nil {
		return nil, key, fmt.Errorf("validate failed, key: %s, err: %s", key, err)
	}

	filter := make(map[string]interface{})
	switch r.Operator {
	case OperatorEqual:
		value, err := common.ConvertIpv6ToFullWord(r.Field, r.Value)
		if err != nil {
			return nil, "convert ipv6 full word failed", err
		}
		filter[r.Field] = map[string]interface{}{
			common.BKDBEQ: value,
		}
	case OperatorNotEqual:
		value, err := common.ConvertIpv6ToFullWord(r.Field, r.Value)
		if err != nil {
			return nil, "convert ipv6 full word failed", err
		}
		filter[r.Field] = map[string]interface{}{
			common.BKDBNE: value,
		}
	case OperatorIn:
		value, err := common.ConvertIpv6ToFullWord(r.Field, r.Value)
		if err != nil {
			return nil, "convert ipv6 full word failed", err
		}
		filter[r.Field] = map[string]interface{}{
			common.BKDBIN: value,
		}
	case OperatorNotIn:
		value, err := common.ConvertIpv6ToFullWord(r.Field, r.Value)
		if err != nil {
			return nil, "convert ipv6 full word failed", err
		}
		filter[r.Field] = map[string]interface{}{
			common.BKDBNIN: value,
		}
	case OperatorLess:
		filter[r.Field] = map[string]interface{}{
			common.BKDBLT: r.Value,
		}
	case OperatorLessOrEqual:
		filter[r.Field] = map[string]interface{}{
			common.BKDBLTE: r.Value,
		}
	case OperatorGreater:
		filter[r.Field] = map[string]interface{}{
			common.BKDBGT: r.Value,
		}
	case OperatorGreaterOrEqual:
		filter[r.Field] = map[string]interface{}{
			common.BKDBGTE: r.Value,
		}
	case OperatorDatetimeLess:
		value, err := getDatetime(r.Value)
		if err != nil {
			return nil, "value", err
		}
		filter[r.Field] = map[string]interface{}{
			common.BKDBLT: value,
		}
	case OperatorDatetimeLessOrEqual:
		value, err := getDatetime(r.Value)
		if err != nil {
			return nil, "value", err
		}
		filter[r.Field] = map[string]interface{}{
			common.BKDBLTE: value,
		}
	case OperatorDatetimeGreater:
		value, err := getDatetime(r.Value)
		if err != nil {
			return nil, "value", err
		}
		filter[r.Field] = map[string]interface{}{
			common.BKDBGT: value,
		}
	case OperatorDatetimeGreaterOrEqual:
		value, err := getDatetime(r.Value)
		if err != nil {
			return nil, "value", err
		}
		filter[r.Field] = map[string]interface{}{
			common.BKDBGTE: value,
		}
	case OperatorBeginsWith:
		filter[r.Field] = map[string]interface{}{
			common.BKDBLIKE: fmt.Sprintf("^%s", r.Value),
		}
	case OperatorNotBeginsWith:
		filter[r.Field] = map[string]interface{}{
			common.BKDBNot: map[string]interface{}{common.BKDBLIKE: fmt.Sprintf("^%s", r.Value)},
		}
	case OperatorContains:
		filter[r.Field] = map[string]interface{}{
			common.BKDBLIKE:    fmt.Sprintf("%s", r.Value),
			common.BKDBOPTIONS: "i",
		}
	case OperatorNotContains:
		filter[r.Field] = map[string]interface{}{
			common.BKDBNot: map[string]interface{}{common.BKDBLIKE: fmt.Sprintf("%s", r.Value)},
		}
	case OperatorsEndsWith:
		filter[r.Field] = map[string]interface{}{
			common.BKDBLIKE: fmt.Sprintf("%s$", r.Value),
		}
	case OperatorNotEndsWith:
		filter[r.Field] = map[string]interface{}{
			common.BKDBNot: map[string]interface{}{common.BKDBLIKE: fmt.Sprintf("%s$", r.Value)},
		}
	case OperatorIsEmpty:
		// array empty
		filter[r.Field] = map[string]interface{}{
			common.BKDBEQ: make([]interface{}, 0),
		}
	case OperatorIsNotEmpty:
		// array not empty
		filter[r.Field] = map[string]interface{}{
			common.BKDBNE: make([]interface{}, 0),
		}
	case OperatorIsNull:
		filter[r.Field] = map[string]interface{}{
			common.BKDBEQ: nil,
		}
	case OperatorIsNotNull:
		filter[r.Field] = map[string]interface{}{
			common.BKDBNE: nil,
		}
	case OperatorExist:
		filter[r.Field] = map[string]interface{}{
			common.BKDBExists: true,
		}
	case OperatorNotExist:
		filter[r.Field] = map[string]interface{}{
			common.BKDBExists: false,
		}
	default:
		return nil, "operator", fmt.Errorf("unsupported operator: %s", r.Operator)
	}
	return filter, "", nil
}

func getDatetime(value interface{}) (interface{}, error) {
	valString, isStr := value.(string)
	if !isStr {
		return nil, errors.New("get date or time failed, value is not string type")
	}
	// 时间类型
	timeType, isTime := util.IsTime(valString)
	if isTime {
		return util.Str2Time(valString, timeType), nil
	}
	// 日期类型
	_, err := time.Parse(timeLayout, valString)
	if err != nil {
		return nil, err
	}
	return valString, nil
}

// GetField get rule field
func (r AtomRule) GetField() []string {
	return []string{r.Field}
}

// CombinedRule TODO
// *************** define query ************************
type CombinedRule struct {
	Condition Condition `json:"condition"`
	Rules     []Rule    `json:"rules"`
}

var (
	// MaxDeep 嵌套层级的深度按树的高度计算，查询条件最大深度为3即最多嵌套2层
	MaxDeep = 3

	// DefaultMaxSliceElementsCount is max elements count of slice(array) condition value.
	DefaultMaxSliceElementsCount = 500

	// DefaultMaxConditionOrRulesCount is default max rules count of one OR combined condition.
	DefaultMaxConditionOrRulesCount = 20
)

// RuleOption is combined condition rule validator option.
type RuleOption struct {
	// NeedSameSliceElementType whether need same type in one slice(array) value or not.
	NeedSameSliceElementType bool

	// MaxSliceElementsCount max slice(array) value elements count, 0 means no limit.
	MaxSliceElementsCount int

	// MaxConditionOrRulesCount max atom rules count in one OR combined condition, 0 means no limit.
	MaxConditionOrRulesCount int

	// MaxConditionAndRulesCount max atom rules count in one AND combined condition, 0 means no limit.
	MaxConditionAndRulesCount int
}

// GetDeep TODO
func (r CombinedRule) GetDeep() int {
	maxChildDeep := 1
	for _, child := range r.Rules {
		childDeep := child.GetDeep()
		if childDeep > maxChildDeep {
			maxChildDeep = childDeep
		}
	}
	return maxChildDeep + 1
}

// Validate validates combined rules with the options.
func (r CombinedRule) Validate(option *RuleOption) (string, error) {
	if err := r.Condition.Validate(); err != nil {
		return "condition", err
	}

	if len(r.Rules) == 0 {
		return "rules", fmt.Errorf("combined rules shouldn't be empty")
	}

	// validate condition rules count.
	if r.Condition == ConditionOr && option.MaxConditionOrRulesCount > 0 &&
		len(r.Rules) > option.MaxConditionOrRulesCount {

		return "rules", fmt.Errorf("too many rules of OR condition: %d max(%d)",
			len(r.Rules), option.MaxConditionOrRulesCount)
	}

	if r.Condition == ConditionAnd && option.MaxConditionAndRulesCount > 0 &&
		len(r.Rules) > option.MaxConditionAndRulesCount {

		return "rules", fmt.Errorf("too many rules of AND condition: %d max(%d)",
			len(r.Rules), option.MaxConditionAndRulesCount)
	}

	for idx, rule := range r.Rules {
		if key, err := rule.Validate(option); err != nil {
			return fmt.Sprintf("rules[%d].%s", idx, key), err
		}
	}
	return "", nil
}

// ToMgo TODO
func (r CombinedRule) ToMgo() (mgoFilter map[string]interface{}, key string, err error) {
	if err := r.Condition.Validate(); err != nil {
		return nil, "condition", err
	}
	if r.Rules == nil || len(r.Rules) == 0 {
		return nil, "rules", fmt.Errorf("combined rules shouldn't be empty")
	}
	filters := make([]map[string]interface{}, 0)
	for idx, rule := range r.Rules {
		filter, key, err := rule.ToMgo()
		if err != nil {
			return nil, fmt.Sprintf("rules[%d].%s", idx, key), err
		}
		filters = append(filters, filter)
	}
	mgoOperator, err := r.Condition.ToMgo()
	if err != nil {
		return nil, "condition", err
	}
	mgoFilter = map[string]interface{}{
		mgoOperator: filters,
	}
	return mgoFilter, "", nil
}

// Match TODO
func (r CombinedRule) Match(matcher Matcher) bool {
	if len(r.Rules) == 0 {
		return true
	}

	switch r.Condition {
	case ConditionAnd:
		for _, rule := range r.Rules {
			if rule.Match(matcher) == false {
				return false
			}
		}
		return true
	case ConditionOr:
		for _, rule := range r.Rules {
			if rule.Match(matcher) == true {
				return true
			}
		}
		return false
	default:
		panic(fmt.Sprintf("unexpected condition %s", r.Condition))
	}
}

// MatchAny if any of the rules matches the matcher, return true
func (r CombinedRule) MatchAny(matcher Matcher) bool {
	if len(r.Rules) == 0 {
		return true
	}

	for _, rule := range r.Rules {
		if rule.MatchAny(matcher) {
			return true
		}
	}
	return false
}

// GetField get rule field
func (r CombinedRule) GetField() []string {
	if len(r.Rules) == 0 {
		return nil
	}

	result := make([]string, 0)
	for _, rule := range r.Rules {
		fields := rule.GetField()
		if len(fields) != 0 {
			result = append(result, fields...)
		}
	}

	return result
}
