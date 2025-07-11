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

package types

import (
	"fmt"
	"os"

	"configcenter/src/common/core/cc/config"
	"configcenter/src/common/version"

	"github.com/rs/xid"
)

// zk path
const (
	CC_SERV_BASEPATH        = "/cc/services/endpoints"
	CC_SERVCONF_BASEPATH    = "/cc/services/config"
	CC_SERVERROR_BASEPATH   = "/cc/services/errors"
	CC_SERVLANG_BASEPATH    = "/cc/services/language"
	CC_SERVNOTICE_BASEPATH  = "/cc/services/notice"
	CC_SERVLIMITER_BASEPATH = "/cc/services/limiter"

	CC_DISCOVERY_PREFIX = "cc_"
)

// cc modules
const (
	CC_MODULE_DATACOLLECTION = "datacollection"
	CC_MODULE_HOST           = "host"
	CC_MODULE_MIGRATE        = "migrate"
	CC_MODULE_PROC           = "proc"
	CC_MODULE_TOPO           = "topo"
	CC_MODULE_APISERVER      = "apiserver"
	CC_MODULE_WEBSERVER      = "webserver"
	CC_MODULE_EVENTSERVER    = "eventserver"
	CC_MODULE_CORESERVICE    = "coreservice"
	GSE_MODULE_PROCSERVER    = "gseprocserver"
	// CC_MODULE_SYNCHRONZESERVER multiple cmdb synchronize data server
	CC_MODULE_SYNCHRONZESERVER = "sync"
	CC_MODULE_OPERATION        = "operation"
	CC_MODULE_TASK             = "task"
	CC_MODULE_CLOUD            = "cloud"
	CC_MODULE_AUTH             = "auth"
	// CC_MODULE_CACHE 缓存服务
	CC_MODULE_CACHESERVICE = "cacheservice"
	// CC_MODULE_TRANSFERSERVICE is a service to transfer data from one cmdb to another
	CC_MODULE_TRANSFERSERVICE = "transfer_service"
)

// AllModule all cc module
var AllModule = map[string]bool{
	CC_MODULE_DATACOLLECTION: true,
	CC_MODULE_HOST:           true,
	CC_MODULE_MIGRATE:        true,
	CC_MODULE_PROC:           true,
	CC_MODULE_TOPO:           true,
	CC_MODULE_APISERVER:      true,
	CC_MODULE_WEBSERVER:      true,
	CC_MODULE_EVENTSERVER:    true,
	CC_MODULE_CORESERVICE:    true,
	// CC_MODULE_SYNCHRONZESERVER: true,
	CC_MODULE_OPERATION:       true,
	CC_MODULE_TASK:            true,
	CC_MODULE_CLOUD:           true,
	CC_MODULE_AUTH:            true,
	CC_MODULE_CACHESERVICE:    true,
	CC_MODULE_TRANSFERSERVICE: true,
}

// cc functionality define
const (
	CCFunctionalityServicediscover = "servicediscover"
	CCFunctionalityMongo           = "mongo"
	CCFunctionalityRedis           = "redis"
)

const (
	// CCConfigureRedis TODO
	CCConfigureRedis = "redis"
	// CCConfigureMongo TODO
	CCConfigureMongo = "mongodb"
	// CCConfigureCommon TODO
	CCConfigureCommon = "common"
	// CCConfigureExtra TODO
	CCConfigureExtra = "extra"
)

// ServerInfo define base server information
type ServerInfo struct {
	IP         string `json:"ip"`
	Port       uint   `json:"port"`
	RegisterIP string `json:"registerip"`
	HostName   string `json:"hostname"`
	Scheme     string `json:"scheme"`
	Version    string `json:"version"`
	Pid        int    `json:"pid"`
	// UUID is used to distinguish which service is master in zookeeper
	UUID string `json:"uuid"`
	// Environment is the server's environment, servers can only discover other servers in the same environment
	Environment string `json:"env"`
}

// NewServerInfo new a ServerInfo object
func NewServerInfo(conf *config.CCAPIConfig) (*ServerInfo, error) {
	ip, err := conf.GetAddress()
	if err != nil {
		return nil, err
	}

	port, err := conf.GetPort()
	if err != nil {
		return nil, err
	}

	registerIP := conf.RegisterIP
	// if no registerIP is set, default to be the ip
	if registerIP == "" {
		registerIP = ip
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	info := &ServerInfo{
		IP:          ip,
		Port:        port,
		RegisterIP:  registerIP,
		HostName:    hostname,
		Scheme:      "http",
		Version:     version.GetVersion(),
		Pid:         os.Getpid(),
		UUID:        xid.New().String(),
		Environment: conf.Environment,
	}
	return info, nil
}

// APIServerServInfo apiserver informaiton
type APIServerServInfo struct {
	ServerInfo
}

// WebServerInfo web server information
type WebServerInfo struct {
	ServerInfo
}

// AuditControllerServInfo audit-controller server information
type AuditControllerServInfo struct {
	ServerInfo
}

// HostControllerServInfo host-controller server information
type HostControllerServInfo struct {
	ServerInfo
}

// MigrateControllerServInfo migrate-controller server information
type MigrateControllerServInfo struct {
	ServerInfo
}

// ObjectControllerServInfo object-controller server information
type ObjectControllerServInfo struct {
	ServerInfo
}

// ProcControllerServInfo proc-controller server information
type ProcControllerServInfo struct {
	ServerInfo
}

// DataCollectionServInfo data-conllection server information
type DataCollectionServInfo struct {
	ServerInfo
}

// HostServerInfo host server information
type HostServerInfo struct {
	ServerInfo
}

// MigrateServInfo migrate server information
type MigrateServInfo struct {
	ServerInfo
}

// ProcServInfo proc server information
type ProcServInfo struct {
	ServerInfo
}

// TopoServInfo topo server information
type TopoServInfo struct {
	ServerInfo
}

// EventServInfo topo server information
type EventServInfo struct {
	ServerInfo
}

// RegisterAddress TODO
// Address convert struct to host address
func (s *ServerInfo) RegisterAddress() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme, s.RegisterIP, s.Port)
}

// Instance TODO
func (s *ServerInfo) Instance() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
