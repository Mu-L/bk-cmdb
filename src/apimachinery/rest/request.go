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

package rest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	commonUtil "configcenter/src/common/util"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

// map[url]responseDataString
var mockResponseMap map[string]string
var once = sync.Once{}

func init() {
	once.Do(func() {
		mockResponseMap = make(map[string]string)
	})
}

// VerbType TODO
// http request verb type
type VerbType string

const (
	// PUT TODO
	PUT VerbType = http.MethodPut
	// POST TODO
	POST VerbType = http.MethodPost
	// GET TODO
	GET VerbType = http.MethodGet
	// DELETE TODO
	DELETE VerbType = http.MethodDelete
	// PATCH TODO
	PATCH VerbType = http.MethodPatch
)

// Request TODO
type Request struct {
	parent *RESTClient

	capability *util.Capability

	verb    VerbType
	params  url.Values
	headers http.Header
	body    []byte
	ctx     context.Context

	// prefixed url
	baseURL string
	// sub path of the url, will be append to baseURL
	subPath string
	// sub path format args
	subPathArgs []interface{}

	// metric additional labels
	metricDimension string

	// request timeout value
	timeout time.Duration

	peek bool
	err  error

	wrappers RequestWrapperChain
}

// WithMetricDimension TODO
// add this request a addition dimension value, which helps us to separate
// request metrics with a dimension label.
func (r *Request) WithMetricDimension(value string) *Request {
	r.metricDimension = value
	return r
}

// WithParams TODO
func (r *Request) WithParams(params map[string]string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	for paramName, value := range params {
		r.params[paramName] = append(r.params[paramName], value)
	}
	return r
}

// WithParamsFromURL TODO
func (r *Request) WithParamsFromURL(u *url.URL) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	params := u.Query()
	for paramName, value := range params {
		r.params[paramName] = append(r.params[paramName], value...)
	}
	return r
}

// WithParam TODO
func (r *Request) WithParam(paramName, value string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[paramName] = append(r.params[paramName], value)
	return r
}

// WithHeaders TODO
func (r *Request) WithHeaders(header http.Header) *Request {
	if r.headers == nil {
		r.headers = header
		return r
	}

	for key, values := range header {
		for _, v := range values {
			r.headers.Add(key, v)
		}
	}
	return r
}

// Peek TODO
func (r *Request) Peek() *Request {
	r.peek = true
	return r
}

// WithContext TODO
func (r *Request) WithContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// WithTimeout TODO
func (r *Request) WithTimeout(d time.Duration) *Request {
	r.timeout = d
	return r
}

// WithBaseURL set request base url
func (r *Request) WithBaseURL(baseURL string) *Request {
	r.baseURL = baseURL
	return r
}

// SubResourcef TODO
func (r *Request) SubResourcef(subPath string, args ...interface{}) *Request {
	r.subPathArgs = args
	return r.subResource(subPath)
}

func (r *Request) subResource(subPath string) *Request {
	subPath = strings.TrimLeft(subPath, "/")
	r.subPath = subPath
	return r
}

// Body TODO
func (r *Request) Body(body interface{}) *Request {
	if nil == body {
		r.body = []byte("")
		return r
	}

	bodyBytes, ok := body.([]byte)
	if ok {
		r.body = bodyBytes
		return r
	}

	valueOf := reflect.ValueOf(body)
	switch valueOf.Kind() {
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.Slice:
		if valueOf.IsNil() {
			r.body = []byte("")
			return r
		}
		break

	case reflect.Struct:
		break

	default:
		r.err = errors.New("body should be one of interface, map, pointer or slice value")
		r.body = []byte("")
		return r
	}

	data, err := json.Marshal(body)
	if err != nil {
		r.err = err
		r.body = []byte("")
		return r
	}

	r.body = data
	return r
}

// WrapURL TODO
func (r *Request) WrapURL() *url.URL {
	finalUrl := &url.URL{}
	if len(r.baseURL) != 0 {
		u, err := url.Parse(r.baseURL)
		if err != nil {
			r.err = err
			return new(url.URL)
		}
		*finalUrl = *u
	}

	if len(r.subPathArgs) > 0 {
		finalUrl.Path = finalUrl.Path + fmt.Sprintf(r.subPath, r.subPathArgs...)
	} else {
		finalUrl.Path = finalUrl.Path + r.subPath
	}

	query := url.Values{}
	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	if r.timeout != 0 {
		query.Set("timeout", r.timeout.String())
	}

	finalUrl.RawQuery = query.Encode()
	return finalUrl
}

func (r *Request) checkToleranceLatency(start *time.Time, url, subPath, rid string) {
	if time.Since(*start) < r.capability.ToleranceLatencyTime {
		return
	}

	if strings.Contains(url, "/watch/resource/") || strings.Contains(url, "/watch/cache/event") ||
		strings.Contains(url, "/cache/event/node/with_start_from") {
		// except resource watch api.
		return
	}

	// request time larger than the maxToleranceLatencyTime time, then log the request
	blog.InfofDepthf(3, "[apimachinery] request exceeded max latency time. cost: %d ms, code: %s, user: %s, %s, "+
		"url: %s, body: %s, rid: %s", time.Since(*start)/time.Millisecond, httpheader.GetAppCode(r.headers),
		httpheader.GetUser(r.headers), r.verb, url, commonUtil.FormatHttpBody(subPath, r.body), rid)
}

// Do http request
func (r *Request) Do() *Result {
	for _, wrapper := range r.wrappers {
		r = wrapper(r)
	}

	result := new(Result)

	rid := commonUtil.ExtractRequestIDFromContext(r.ctx)
	if rid == "" {
		rid = httpheader.GetRid(r.headers)
	}
	result.Rid = rid

	if r.err != nil {
		result.Err = r.err
		return result
	}

	if r.capability.Mock.Mocked {
		return r.handleMockResult()
	}

	client := r.capability.Client
	if client == nil {
		client = http.DefaultClient
	}

	hosts, err := r.capability.Discover.GetServers()
	if err != nil {
		result.Err = err
		return result
	}

	maxRetryCycle := 3
	var retries int
	for try := 0; try < maxRetryCycle; try++ {
		for index, host := range hosts {
			retries = try + index
			url := host + r.WrapURL().String()
			req, err := http.NewRequest(string(r.verb), url, bytes.NewReader(r.body))
			if err != nil {
				result.Err = err
				result.Rid = rid
				return result
			}

			if r.ctx != nil {
				req = req.WithContext(r.ctx)
			}

			req.Header = commonUtil.CloneHeader(r.headers)
			if len(req.Header) == 0 {
				req.Header = make(http.Header)
			}
			// 删除 Accept-Encoding 避免返回值被压缩
			req.Header.Del("Accept-Encoding")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			if retries > 0 {
				r.tryThrottle(url)
			}

			var needRetry bool
			result, needRetry = r.do(client, req, result, rid)
			if needRetry {
				continue
			}
			return result
		}

	}

	result.Err = errors.New("unexpected error")
	return result
}

func (r *Request) do(client util.HttpClient, req *http.Request, result *Result, rid string) (*Result, bool) {
	url := req.URL.String()
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		// "Connection reset by peer" is a special err which in most scenario is a a transient error.
		// Which means that we can retry it. And so does the GET operation.
		// While the other "write" operation can not simply retry it again, because they are not idempotent.

		blog.Errorf("[apimachinery] %s %s with body %s, but %v, rid: %s", string(r.verb), url,
			commonUtil.FormatHttpBody(r.subPath, r.body), err, rid)
		r.checkToleranceLatency(&start, url, r.subPath, rid)
		if !isConnectionReset(err) || r.verb != GET {
			result.Err = err
			return result, false
		}

		// retry now
		time.Sleep(20 * time.Millisecond)
		return result, true
	}

	// collect request metrics
	if r.parent.requestDuration != nil {
		labels := prometheus.Labels{
			"handler":     r.subPath,
			"status_code": strconv.Itoa(result.StatusCode),
			"dimension":   r.metricDimension,
		}

		r.parent.requestDuration.With(labels).Observe(float64(time.Since(start) / time.Millisecond))
	}

	// record latency if needed
	r.checkToleranceLatency(&start, url, r.subPath, rid)

	var body []byte
	if resp.Body != nil {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				// retry now
				time.Sleep(20 * time.Millisecond)
				return result, true
			}
			result.Err = err
			blog.Errorf("[apimachinery] %s %s with body %s, err: %v, rid: %s", string(r.verb), url,
				commonUtil.FormatHttpBody(r.subPath, r.body), err, rid)
			return result, false
		}
		body = data
	}

	if blog.V(4) {
		blog.V(4).InfoDepthf(2, "[apimachinery] cost: %dms, %s %s with body %s, response status: %s, "+
			"response body: %s, rid: %s", time.Since(start)/time.Millisecond, string(r.verb), url,
			commonUtil.FormatHttpBody(r.subPath, r.body), resp.Status,
			commonUtil.FormatHttpBody(r.subPath, body), rid)
	}

	result.Body = body
	result.StatusCode = resp.StatusCode
	result.Status = resp.Status
	result.Header = resp.Header

	return result, false
}

const maxLatency = 100 * time.Millisecond

func (r *Request) tryThrottle(url string) {
	now := time.Now()
	if r.capability.Throttle != nil {
		r.capability.Throttle.Accept()
	}

	if latency := time.Since(now); latency > maxLatency {
		blog.V(3).Infof("Throttling request took %d ms, verb: %s, request: %s", latency, r.verb, url)
	}
}

// Result TODO
type Result struct {
	Rid        string
	Body       []byte
	Err        error
	StatusCode int
	Status     string
	Header     http.Header
}

// IntoCmdbResp decode result into cmdb inner response
func (r *Result) IntoCmdbResp(obj interface{}) error {
	if r.Err != nil {
		return r.Err
	}

	if len(r.Body) != 0 {
		// parse api gateway response format to cmdb inner response format
		if gjson.GetBytes(r.Body, common.BkAPIErrorCode).Exists() {
			keyMap := map[string]string{
				common.BkAPIErrorCode:    common.HTTPBKAPIErrorCode,
				common.BkAPIErrorMessage: common.HTTPBKAPIErrorMessage,
			}
			var err error
			r.Body, err = json.ReplaceJsonKey(r.Body, keyMap)
			if err != nil {
				return fmt.Errorf("replace resp(%s) key failed, err: %v", string(r.Body), err)
			}
		}

		err := json.Unmarshal(r.Body, obj)
		if err != nil {
			if r.StatusCode >= 300 {
				return fmt.Errorf("http request err: %s", string(r.Body))
			}
			blog.Errorf("invalid response body, unmarshal json failed, reply:%s, error:%s", r.Body, err.Error())
			return fmt.Errorf("http response err: %v, raw data: %s", err, r.Body)
		}
		return nil
	}

	if r.StatusCode >= 300 {
		return fmt.Errorf("http request failed: %s", r.Status)
	}
	return nil
}

// Into TODO
func (r *Result) Into(obj interface{}) error {
	if nil != r.Err {
		return r.Err
	}

	if 0 != len(r.Body) {
		err := json.Unmarshal(r.Body, obj)
		if err != nil {
			if r.StatusCode >= 300 {
				return fmt.Errorf("http request err: %s", string(r.Body))
			}
			blog.Errorf("invalid response body, unmarshal json failed, reply:%s, error:%s", r.Body, err.Error())
			return fmt.Errorf("http response err: %v, raw data: %s", err, r.Body)
		}
	} else if r.StatusCode >= 300 {
		return fmt.Errorf("http request failed: %s", r.Status)
	}
	return nil
}

// IntoJsonString TODO
func (r *Result) IntoJsonString() (*metadata.JsonStringResp, error) {
	if nil != r.Err {
		return nil, r.Err
	}

	if 0 == len(r.Body) {
		return nil, fmt.Errorf("http request failed: %s", r.Status)
	}
	elements := gjson.GetManyBytes(r.Body, "result", "bk_error_code", "bk_error_msg", "permission", "data")

	// check result
	if !elements[0].Exists() {
		blog.Errorf("invalid http response, no result field, body: %s, rid: %s", r.Body, r.Rid)
		return nil, fmt.Errorf("invalid http response, body: %s", r.Body)
	}

	// check error code
	if !elements[1].Exists() {
		blog.Errorf("invalid http response, no bk_error_code field, body: %s, rid: %s", r.Body, r.Rid)
		return nil, fmt.Errorf("invalid http response, body: %s", r.Body)
	}

	// check error message
	if !elements[2].Exists() {
		blog.Errorf("invalid http response, no bk_error_msg field, body: %s, rid: %s", r.Body, r.Rid)
		return nil, fmt.Errorf("invalid http response, body: %s", r.Body)
	}

	// check data
	if !elements[4].Exists() {
		blog.Errorf("invalid http response, no data field, body: %s, rid: %s", r.Body, r.Rid)
		return nil, fmt.Errorf("invalid http response, body: %s", r.Body)
	}

	resp := new(metadata.JsonStringResp)
	resp.Result = elements[0].Bool()
	resp.Code = int(elements[1].Int())
	resp.ErrMsg = elements[2].String()
	// parse permission field
	if elements[3].Exists() {
		raw := elements[3].Raw
		if len(raw) != 0 {
			perm := new(metadata.IamPermission)
			if err := json.Unmarshal([]byte(raw), &perm); err != nil {
				blog.Errorf("invalid http response, invalid permission field, body: %s, rid: %s", r.Body, r.Rid)
				return nil, fmt.Errorf("http response with invalid permission field, body: %s", r.Body)
			}
			resp.Permissions = perm
		}
	}
	resp.Data = elements[4].Raw

	return resp, nil
}

// IntoJsonCntInfoString TODO
func (r *Result) IntoJsonCntInfoString() (*metadata.JsonCntInfoResp, error) {
	if nil != r.Err {
		return nil, r.Err
	}

	if 0 == len(r.Body) {
		return nil, fmt.Errorf("http request failed: %s", r.Status)
	}
	elements := gjson.GetManyBytes(r.Body, "result", "bk_error_code", "bk_error_msg", "permission", "data")

	// check result
	if !elements[0].Exists() {
		blog.Errorf("invalid http response, no result field, body: %s, rid: %s", r.Body, r.Rid)
		return nil, fmt.Errorf("invalid http response, body: %s", r.Body)
	}

	// check error code
	if !elements[1].Exists() {
		blog.Errorf("invalid http response, no bk_error_code field, body: %s, rid: %s", r.Body, r.Rid)
		return nil, fmt.Errorf("invalid http response, body: %s", r.Body)
	}

	// check error message
	if !elements[2].Exists() {
		blog.Errorf("invalid http response, no bk_error_msg field, body: %s, rid: %s", r.Body, r.Rid)
		return nil, fmt.Errorf("invalid http response, body: %s", r.Body)
	}

	// check data
	if !elements[4].Exists() {
		blog.Errorf("invalid http response, no data field, body: %s, rid: %s", r.Body, r.Rid)
		return nil, fmt.Errorf("invalid http response, body: %s", r.Body)
	}

	// check data.count
	if !elements[4].Get("count").Exists() {
		blog.Errorf("invalid http response, no data.count field, body: %s, rid: %s", r.Body, r.Rid)
		return nil, fmt.Errorf("invalid http response, body: %s", r.Body)
	}

	// check data.info
	if !elements[4].Get("info").Exists() {
		blog.Errorf("invalid http response, no data.info field, body: %s, rid: %s", r.Body, r.Rid)
		return nil, fmt.Errorf("invalid http response, body: %s", r.Body)
	}

	resp := new(metadata.JsonCntInfoResp)
	resp.Result = elements[0].Bool()
	resp.Code = int(elements[1].Int())
	resp.ErrMsg = elements[2].String()

	// parse permission field
	if elements[3].Exists() {
		raw := elements[3].Raw
		if len(raw) != 0 {
			perm := new(metadata.IamPermission)
			if err := json.Unmarshal([]byte(raw), perm); err != nil {
				blog.Errorf("invalid http response, invalid permission field, body: %s, rid: %s", r.Body, r.Rid)
				return nil, fmt.Errorf("http response with invalid permission field, body: %s", r.Body)
			}
			resp.Permissions = perm
		}
	}

	// set count field
	resp.Data.Count = elements[4].Get("count").Int()

	// set info field
	resp.Data.Info = elements[4].Get("info").Raw

	return resp, nil
}

func (r *Request) handleMockResult() *Result {
	if r.capability.Mock.SetMockData {
		if r.capability.Mock.MockData == nil {
			mockResponseMap[r.WrapURL().String()] = ""
			return &Result{
				Body:       []byte(""),
				Err:        nil,
				StatusCode: http.StatusOK,
			}
		}

		switch reflect.ValueOf(r.capability.Mock.MockData).Kind() {
		case reflect.String:
			body := r.capability.Mock.MockData.(string)
			mockResponseMap[r.WrapURL().String()] = body
			return &Result{
				Body:       []byte(body),
				Err:        nil,
				StatusCode: http.StatusOK,
			}
		case reflect.Interface:
			fallthrough
		case reflect.Map:
			fallthrough
		case reflect.Ptr:
			fallthrough
		case reflect.Struct:
			js, err := json.Marshal(r.capability.Mock.MockData)
			if err != nil {
				return &Result{
					Body:       nil,
					Err:        err,
					StatusCode: http.StatusOK,
				}
			}
			mockResponseMap[r.WrapURL().String()] = string(js)
			return &Result{
				Body:       js,
				Err:        nil,
				StatusCode: http.StatusOK,
			}
		default:
			panic("unsupported mock data")
		}
	}

	body, exist := mockResponseMap[r.WrapURL().String()]
	if exist {
		return &Result{
			Body:       []byte(body),
			Err:        nil,
			StatusCode: http.StatusOK,
		}
	}

	panic("got empty mock response")
}

// isConnectionReset TODO
// Returns if the given err is "connection reset by peer" error.
func isConnectionReset(err error) bool {
	if urlErr, ok := err.(*url.Error); ok {
		err = urlErr.Err
	}
	if opErr, ok := err.(*net.OpError); ok {
		err = opErr.Err
	}
	if osErr, ok := err.(*os.SyscallError); ok {
		err = osErr.Err
	}
	if errno, ok := err.(syscall.Errno); ok && errno == syscall.ECONNRESET {
		return true
	}
	return false
}
