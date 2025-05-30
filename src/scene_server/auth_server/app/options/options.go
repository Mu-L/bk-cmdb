/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package options TODO
package options

import (
	"configcenter/src/ac/iam"
	"configcenter/src/common/auth"
	"configcenter/src/common/core/cc/config"
	"configcenter/src/common/ssl"

	"github.com/spf13/pflag"
)

// ServerOption define option of server in flags
type ServerOption struct {
	ServConf *config.CCAPIConfig
}

// NewServerOption create a ServerOption object
func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf: config.NewCCAPIConfig(),
	}

	return &s
}

// AddFlags add flags
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	s.ServConf.AddFlags(fs, "127.0.0.1:60014")
	fs.Var(auth.EnableAuthFlag, "enable-auth", "The auth center enable status, true for enabled, false for disabled")
}

// Config TODO
type Config struct {
	Auth iam.AuthConfig
	TLS  ssl.TLSClientConfig
}
