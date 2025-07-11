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

package service_test

import (
	"net/http"
	"testing"

	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/language"
	"configcenter/src/source_controller/coreservice/app/options"
	"configcenter/src/source_controller/coreservice/service"
	"configcenter/src/storage/dal/mongo"

	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/require"
)

var defaultHeader = func() http.Header {

	header := headerutil.GenCommonHeader("test_user", "test_owner", "test_req_id")
	httpheader.SetLanguage(header, "en")
	return header
}()

func startCoreService(t *testing.T, ip string, port uint) {

	// create a logics service
	coreService := service.New()

	// register the server hander
	bonServer := backbone.Server{
		ListenAddr: ip,
		ListenPort: port,
		Handler:    restful.NewContainer().Add(coreService.WebService()),
		TLS:        backbone.TLSConfig{},
	}

	// set backbone config
	bonC := &backbone.Config{
		Server: bonServer,
	}

	// new server instance
	engine, err := backbone.NewMockBackbone(bonC)
	require.NoError(t, err)

	errE, err := errors.New("../../../../resources/errors/")
	require.NoError(t, err)
	lan, err := language.New("../../../../resources/language/")
	require.NoError(t, err)
	engine.CCErr = errE
	engine.Language = lan
	// set the config
	coreService.SetConfig(options.Config{
		Mongo: mongo.Config{
			Address:  "cc:cc@localhost:27010,localhost:27011,localhost:27012,localhost:27013",
			User:     "cc",
			Password: "cc",
			Database: "cmdb",
		},
	}, engine, engine.CCErr, engine.Language)

	return
}
