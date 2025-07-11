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

package util

import (
	"reflect"
	"testing"
)

func TestInArray(t *testing.T) {
	type args struct {
		obj    interface{}
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				target: []string{"a", "b", "c"},
				obj:    "a",
			},
			want: true,
		},
		{
			args: args{
				target: []string{"a", "b", "c"},
				obj:    "d",
			},
			want: false,
		},
		{
			args: args{
				target: []interface{}{"a", "b", "c", 1},
				obj:    1,
			},
			want: true,
		},
		{
			args: args{
				target: []interface{}{"a", "b", "c", 1},
				obj:    int64(1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InArray(tt.args.obj, tt.args.target); got != tt.want {
				t.Errorf("InArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayUnique(t *testing.T) {
	type args struct {
		a interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet []interface{}
	}{
		{
			args: args{
				[]interface{}{"a", "b", "c", 1},
			},
			wantRet: []interface{}{"a", "b", "c", 1},
		},
		{
			args: args{
				[]interface{}{"a", "b", "c", 1, 1, "a", ""},
			},
			wantRet: []interface{}{"a", "b", "c", 1, ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := ArrayUnique(tt.args.a); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("ArrayUnique() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestRemoveDuplicatesAndEmpty(t *testing.T) {
	type args struct {
		a []string
	}
	tests := []struct {
		name    string
		args    args
		wantRet []string
	}{
		{
			args: args{
				[]string{"a", "b", "c", "a", ""},
			},
			wantRet: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := RemoveDuplicatesAndEmpty(tt.args.a); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("RemoveDuplicatesAndEmpty() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestStrArrDiff(t *testing.T) {
	type args struct {
		slice1 []string
		slice2 []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{
				[]string{"a", "b", "c", "a", ""},
				[]string{"a", "b"},
			},
			want: []string{"c", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrArrDiff(tt.args.slice1, tt.args.slice2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrArrDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInStrArr(t *testing.T) {
	type args struct {
		arr []string
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{[]string{"key"}, "key"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InStrArr(tt.args.arr, tt.args.key); got != tt.want {
				t.Errorf("InStrArr() = %v, want %v", got, tt.want)
			}
		})
	}
}
