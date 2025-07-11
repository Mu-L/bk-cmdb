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

// Package redis TODO
package redis

import (
	"context"
	"strconv"
	"strings"

	"configcenter/src/common/ssl"

	"github.com/go-redis/redis/v7"
)

// Config define redis config
type Config struct {
	Address          string
	Password         string
	Database         string
	MasterName       string
	SentinelPassword string
	// for datacollection, notify if the snapshot redis is in use
	Enable       string
	MaxOpenConns int
	// tls config
	TLSConfig *ssl.TLSClientConfig
}

// NewFromConfig returns new redis client from config
func NewFromConfig(cfg Config) (Client, error) {
	dbNum, err := strconv.Atoi(cfg.Database)
	if err != nil {
		return nil, err
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 3000
	}

	tlsConf, useTLS, err := ssl.NewTLSConfigFromConf(cfg.TLSConfig)
	if err != nil {
		return nil, err
	}

	var client Client
	if cfg.MasterName == "" {
		option := &redis.Options{
			Addr:     cfg.Address,
			Password: cfg.Password,
			DB:       dbNum,
			PoolSize: cfg.MaxOpenConns,
		}
		if useTLS {
			option.TLSConfig = tlsConf
		}
		client = NewClient(option)
	} else {
		hosts := strings.Split(cfg.Address, ",")
		option := &redis.FailoverOptions{
			MasterName:       cfg.MasterName,
			SentinelAddrs:    hosts,
			Password:         cfg.Password,
			DB:               dbNum,
			PoolSize:         cfg.MaxOpenConns,
			SentinelPassword: cfg.SentinelPassword,
		}
		if useTLS {
			option.TLSConfig = tlsConf
		}
		client = NewFailoverClient(option)
	}

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, err
}

// IsNilErr returns whether err is nil error
func IsNilErr(err error) bool {
	return redis.Nil == err
}
