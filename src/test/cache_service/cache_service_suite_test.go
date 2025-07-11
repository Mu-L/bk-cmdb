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

package cache_service_test

import (
	"context"
	"testing"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/src/common/mapstr"
	"configcenter/src/test"
	"configcenter/src/test/reporter"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	header        = test.GetHeader()
	db            = test.GetDB()
	clientSet     = test.GetClientSet()
	cacheCli      = clientSet.CacheService().Cache()
	generalResCli = cacheCli.GeneralRes()
	hostCli       = cacheCli.Host()
	hostSvrCli    = clientSet.HostServer()
	topoSvrCli    = clientSet.TopoServer()
)

func TestCacheService(t *testing.T) {
	RegisterFailHandler(util.Fail)
	reporters := []Reporter{
		reporter.NewHtmlReporter(test.GetReportDir()+"cacheservice.html", test.GetReportUrl(), true),
	}
	RunSpecsWithDefaultAndCustomReporters(t, "CacheService Suite", reporters)
}

func deleteAllFullSyncCond() {
	err := db.Table(fullsynccond.BKTableNameFullSyncCond).Delete(context.Background(), make(mapstr.MapStr))
	Expect(err).NotTo(HaveOccurred())
}
