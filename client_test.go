package ozonCore

import (
	"context"
	"net/http"
	"testing"
	"time"
)

const (
	testTimeout = 5 * time.Second
)

type TestRequestRequest struct {
	FirstField  string `json:"first_field"`
	SecondField int32  `json:"second_field"`
}

// TestRequestWithMapBody 回归测试：请求体是 map（如 /v1/notification/* 的 body 直接传
// map[string]interface{}）时，getDefaultValues 不能对 map 调用 v.Elem() 而 panic
// （线上曾出现 "reflect: call of reflect.Value.Elem on map Value"）。
func TestRequestWithMapBody(t *testing.T) {
	c := NewMockClient(NewMockHttpHandler(http.StatusOK, `{"code": 0}`, nil))

	body := map[string]interface{}{
		"url":   "https://example.com/callback",
		"types": []string{"TYPE_NEW_POSTING"},
	}
	var respData map[string]interface{}
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	resp, err := c.Request(ctx, http.MethodPost, "/v1/notification/set", body, &respData, nil)
	if err != nil {
		t.Fatalf("map body 请求不应报错: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, http.StatusOK)
	}
	if respData["code"] == nil {
		t.Errorf("response not parsed: %v", respData)
	}
}

type TestRequestResponse struct {
	FirstField  string `json:"first_header"`
	SecondField struct {
		Key int32 `json:"key"`
	} `json:"second_header"`
	ThirdField []struct {
		FirstElement string `json:"first_element"`
	} `json:"third_header"`
}

func TestRequest(t *testing.T) {
	tests := []struct {
		statusCode int
		headers    map[string]string
		params     *TestRequestRequest
		response   string
	}{
		{
			http.StatusOK,
			map[string]string{
				"first_header":  "first-value",
				"second_header": "second-value",
			},
			&TestRequestRequest{
				FirstField:  "test",
				SecondField: 123,
			},
			`{
				"first_field": "first-value",
				"second_field": {
					"key": 12
				},
				"third_field": [
					{
						"first_element": "first-element-value"
					}
				]
			}`,
		},
	}

	for _, test := range tests {
		c := NewMockClient(NewMockHttpHandler(test.statusCode, test.response, test.headers))

		respStruct := &TestRequestResponse{}
		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Request(ctx, http.MethodPost, "/", test.params, respStruct, nil)
		cancel()

		if err != nil {
			t.Error(err)
			continue
		}

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}
	}
}
