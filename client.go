package ozonCore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	baseUrl string
	Options map[string]string

	client HttpClient
}

func NewClient(client HttpClient, baseUrl string, opts map[string]string) *Client {
	return &Client{
		Options: opts,
		client:  client,
		baseUrl: baseUrl,
	}
}

func NewMockClient(handler http.HandlerFunc) *Client {
	return &Client{
		client: NewMockHttpClient(handler),
	}
}

func (c Client) newRequest(ctx context.Context, method string, uri string, body interface{}) (*http.Request, error) {
	var err error
	var bodyJson []byte

	// Set default values for empty fields if `default` tag is present
	// And body is not nil
	if body != nil {
		if err := getDefaultValues(reflect.ValueOf(body)); err != nil {
			return nil, err
		}

		bodyJson, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	uri, err = url.JoinPath(c.baseUrl, uri)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}

	for k, v := range c.Options {
		req.Header.Add(k, v)
	}

	return req, nil
}

const (
	// 429 限流时最多重试次数
	max429Retries = 3
)

// Request 发送请求并解析响应。
//
// 对 429 Too Many Requests 会读取 Retry-After 头做退避重试（最多 3 次）：
// Ozon 的限流错误不会处理请求，因此对所有方法重试都是安全的。
func (c Client) Request(
	ctx context.Context,
	method string,
	path string,
	req,
	resp interface{},
	options map[string]string,
) (*Response, error) {
	delay := 1 * time.Second

	for attempt := 0; ; attempt++ {
		response, err := c.doRequest(ctx, method, path, req, resp, options)
		if err != nil {
			return nil, err
		}

		if response.StatusCode == http.StatusTooManyRequests && attempt < max429Retries {
			if retryAfter := getRetryAfter(response); retryAfter > 0 {
				delay = retryAfter
			}

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return response, ctx.Err()
			}
			delay *= 2
			continue
		}

		return response, nil
	}
}

// getRetryAfter 从响应头解析重试等待时间。
// 优先读取 Item-Retry-After（Ozon 商品操作接口专用，单位分钟），
// 其次 Retry-After（秒或 HTTP 日期）。
func getRetryAfter(response *Response) time.Duration {
	if itemRetryAfter := response.Header.Get("Item-Retry-After"); itemRetryAfter != "" {
		if minutes, err := strconv.Atoi(itemRetryAfter); err == nil && minutes > 0 {
			return time.Duration(minutes) * time.Minute
		}
	}

	retryAfter := response.Header.Get("Retry-After")
	if retryAfter == "" {
		return 0
	}

	if seconds, err := strconv.Atoi(retryAfter); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}

	if t, err := http.ParseTime(retryAfter); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}

	return 0
}

func (c Client) doRequest(
	ctx context.Context,
	method string,
	path string,
	req,
	resp interface{},
	options map[string]string,
) (*Response, error) {
	httpReq, err := c.newRequest(ctx, method, path, req)
	if err != nil {
		fmt.Println("c Request:", err)
		return nil, err
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		fmt.Println("c Request 2:", err)
		return nil, err
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Println("c Request 3:", err)
		return nil, err
	}

	response := &Response{}
	response.Data = resp
	response.StatusCode = httpResp.StatusCode
	response.Header = httpResp.Header
	if httpResp.StatusCode == http.StatusOK {
		if options["Content-Type"] == "" || options["Content-Type"] == "application/json" {
			err = json.Unmarshal(body, &response.Data)
		} else {
			response.Data = body
		}

	} else {
		if options["Content-Type"] == "" || options["Content-Type"] == "application/json" {
			err = json.Unmarshal(body, &response)
		} else {
			response.Data = body
		}
	}
	if err != nil {
		bodyText := string(body)
		if len(body) > 1010 {
			bodyText = string(body[:1000])
		}
		fmt.Println("c Request 4:", err, bodyText)
		//fmt.Println(fmt.Sprintf("request: %s", string(body)))
		return nil, err
	}

	return response, nil
}

type MockHttpClient struct {
	handler http.HandlerFunc
}

func NewMockHttpClient(handler http.HandlerFunc) *MockHttpClient {
	return &MockHttpClient{
		handler: handler,
	}
}

func (c MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	rr := httptest.NewRecorder()
	c.handler.ServeHTTP(rr, req)

	return rr.Result(), nil
}

func NewMockHttpHandler(statusCode int, json string, headers map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(headers) > 0 {
			for key, value := range headers {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(statusCode)
		w.Write([]byte(json))
	}
}
