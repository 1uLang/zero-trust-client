package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func HTTPRequest(addr, method string,
	params map[string]interface{},
	headers map[string]string) (respBody []byte, err error) {
	var body io.Reader
	method = strings.ToUpper(method)
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	//非get 参数设置在body中 以json形式传输
	if method != "GET" && len(params) > 0 {
		buf, _ := json.Marshal(params)
		body = bytes.NewReader(buf)
	}
	req, err := http.NewRequest(method, addr, body)
	if err != nil {
		return nil, err
	}
	// 设置URL参数

	if method == "GET" {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()
	}
	// 设置header
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
