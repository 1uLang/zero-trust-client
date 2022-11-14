package application

import (
	"encoding/json"
	"errors"
	"github.com/1uLang/zero-trust-client"
	log "github.com/sirupsen/logrus"
)

const (
	listPath    = ""
	optionsPath = ""
)

type application struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
type applicationItem struct {
	Logo string `json:"logo"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type responsePublic struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 应用模块

// List 获取 用户应用列表
func List(keyword string, index, size int) (int64, []applicationItem, error) {
	respBytes, err := zero_trust_client.Request(listPath, "get", map[string]interface{}{
		"keyword":   keyword,
		"pageIndex": index,
		"pageSize":  size,
	})
	if err != nil {
		log.Warnf("获取应用列表失败：%s", err)
		return 0, nil, errors.New("获取应用列表失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
		Data struct {
			Total int64             `json:"total"`
			Items []applicationItem `json:"items"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return 0, nil, errors.New("获取应用列表失败：服务器异常！")
	}
	if resp.Code != 0 {
		return 0, nil, errors.New("获取应用列表失败：" + resp.Message)
	}
	return resp.Data.Total, resp.Data.Items, nil
}

// Options 应用下拉框
func Options() ([]application, error) {
	respBytes, err := zero_trust_client.Request(optionsPath, "get", map[string]interface{}{})
	if err != nil {
		log.Warnf("获取应用失败：%s", err)
		return nil, errors.New("获取应用失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
		Data struct {
			Options []application `json:"options"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return nil, errors.New("获取应用失败：服务器异常！")
	}
	if resp.Code != 0 {
		return nil, errors.New("获取应用失败：" + resp.Message)
	}
	return resp.Data.Options, nil
}
