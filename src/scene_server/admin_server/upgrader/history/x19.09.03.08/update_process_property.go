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

package x19_09_03_08

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func updateProcessIntProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	for _, property := range []string{"priority", "auto_time_gap", "timeout"} {
		filter := map[string]interface{}{
			common.BKObjIDField:      "process",
			common.BKPropertyIDField: property,
		}
		now := time.Now()
		doc := map[string]interface{}{
			"option": map[string]int{
				"min": 0,
				"max": 10000,
			},
			common.CreateTimeField: &now,
			common.LastTimeField:   &now,
		}
		err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc)
		if nil != err {
			blog.Errorf("[upgrade x19_09_03_08] update property %s failed, err: %+v", property, err)
			return err
		}

	}
	return nil
}
