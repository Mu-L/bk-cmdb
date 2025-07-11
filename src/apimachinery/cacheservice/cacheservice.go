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

// Package cacheservice TODO
package cacheservice

import (
	"fmt"

	"configcenter/src/apimachinery/cacheservice/cache/event"
	"configcenter/src/apimachinery/cacheservice/cache/general"
	"configcenter/src/apimachinery/cacheservice/cache/host"
	"configcenter/src/apimachinery/cacheservice/cache/topology"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
)

// Cache TODO
type Cache interface {
	Host() host.Interface
	Topology() topology.Interface
	Event() event.Interface
	GeneralRes() general.Interface
}

// CacheServiceClientInterface TODO
type CacheServiceClientInterface interface {
	Cache() Cache
}

// NewCacheServiceClient TODO
func NewCacheServiceClient(c *util.Capability, version string) CacheServiceClientInterface {
	base := fmt.Sprintf("/cache/%s", version)
	return &cacheService{
		restCli: rest.NewRESTClient(c, base),
	}
}

type cacheService struct {
	restCli rest.ClientInterface
}

type cache struct {
	restCli rest.ClientInterface
}

// Cache TODO
func (c *cacheService) Cache() Cache {
	return &cache{
		restCli: c.restCli,
	}
}

// Host TODO
func (c *cache) Host() host.Interface {
	return host.NewCacheClient(c.restCli)
}

// Topology TODO
func (c *cache) Topology() topology.Interface {
	return topology.NewCacheClient(c.restCli)
}

// Event TODO
func (c *cache) Event() event.Interface {
	return event.NewCacheClient(c.restCli)
}

// GeneralRes is the general resource cache client
func (c *cache) GeneralRes() general.Interface {
	return general.NewCacheClient(c.restCli)
}
