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

// Package parser TODO
package parser

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"configcenter/src/ac/meta"
	"configcenter/src/common/backbone"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful/v3"
)

// ParseAttribute TODO
func ParseAttribute(req *restful.Request, engine *backbone.Engine) (*meta.AuthAttribute, error) {
	elements, err := urlParse(req.Request.URL.Path)
	if err != nil {
		return nil, err
	}

	requestContext := &RequestContext{
		Rid:      httpheader.GetRid(req.Request.Header),
		Header:   req.Request.Header,
		Method:   req.Request.Method,
		URI:      req.Request.URL.Path,
		Elements: elements,
		getBody: func() (body []byte, err error) {
			body, err = util.PeekRequest(req.Request)
			if err != nil {
				return nil, err
			}
			return
		},
	}

	stream, err := newParseStream(requestContext, engine)
	if err != nil {
		return nil, err
	}

	return stream.Parse()
}

// ParseCommonInfo get common info from req, aims at avoiding too much repeat code
func ParseCommonInfo(requestHeader http.Header) (*meta.CommonInfo, error) {
	commonInfo := new(meta.CommonInfo)

	userInfo, err := ParseUserInfo(requestHeader)
	if err != nil {
		return nil, err
	}
	commonInfo.User = *userInfo

	return commonInfo, nil
}

// ParseUserInfo TODO
func ParseUserInfo(requestHeader http.Header) (*meta.UserInfo, error) {
	userInfo := new(meta.UserInfo)
	user := httpheader.GetUser(requestHeader)
	if len(user) == 0 {
		return nil, errors.New("parse user info failed, miss user header in your request header")
	}
	userInfo.UserName = user
	supplierID := httpheader.GetSupplierAccount(requestHeader)
	if len(supplierID) == 0 {
		return nil, errors.New("parse user info failed, miss bk_supplier_id in your request header")
	}
	userInfo.SupplierAccount = supplierID
	return userInfo, nil
}

// url example: /api/v3/create/model
var urlRegex = regexp.MustCompile(`^/api/([^/]+)(/[^/]+)+/?$`)

func urlParse(url string) (elements []string, err error) {
	if !urlRegex.MatchString(url) {
		return nil, fmt.Errorf("invalid url format, url=%s", url)
	}

	return strings.Split(url, "/")[1:], nil
}
