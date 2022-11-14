package account

import (
	"encoding/json"
	"errors"
	"github.com/1uLang/zero-trust-client"
	log "github.com/sirupsen/logrus"
)

const (
	createPath = ""
	editPath   = ""
	listPath   = ""
	deletePath = ""
)

type accountItem struct {
	ApplicationId   int64  `json:"application_id"`   // 所属应用ID
	ApplicationName string `json:"application_name"` // 应用名称
	Username        string `json:"username"`         // 账户名称
	CreatedAt       string `json:"createdAt"`        // 创建时间
}
type responsePublic struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Create 申请应用账户
func Create(appId int64, username, password string) error {
	respBytes, err := zero_trust_client.Request(createPath, "post", map[string]interface{}{
		"appId":    appId,
		"username": username,
		"password": password,
	})
	if err != nil {
		log.Warnf("申请应用账户失败：%s", err)
		return errors.New("申请应用账户失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return errors.New("申请应用账户失败：服务器异常！")
	}
	if resp.Code != 0 {
		return errors.New("申请应用账户失败：" + resp.Message)
	}
	return nil
}

// Edit 修改应用账户密码
func Edit(id int64, password string) error {
	respBytes, err := zero_trust_client.Request(editPath, "post", map[string]interface{}{
		"id":       id,
		"password": password,
	})
	if err != nil {
		log.Warnf("修改应用账户密码失败：%s", err)
		return errors.New("修改应用账户密码失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return errors.New("修改应用账户密码失败：服务器异常！")
	}
	if resp.Code != 0 {
		return errors.New("修改应用账户密码失败：" + resp.Message)
	}
	return nil
}

// List  获取 应用账户列表
func List(keyword string, state, index, size int) (int64, []accountItem, error) {
	respBytes, err := zero_trust_client.Request(listPath, "get", map[string]interface{}{
		"keyword":   keyword,
		"state":     state,
		"pageIndex": index,
		"pageSize":  size,
	})
	if err != nil {
		log.Warnf("获取应用账户列表失败：%s", err)
		return 0, nil, errors.New("获取应用账户列表失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
		Data struct {
			Total int64         `json:"total"`
			Items []accountItem `json:"items"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return 0, nil, errors.New("获取应用账户列表失败：服务器异常！")
	}
	if resp.Code != 0 {
		return 0, nil, errors.New("获取应用账户列表失败：" + resp.Message)
	}
	return resp.Data.Total, resp.Data.Items, nil
}

// Delete 删除应用账户
func Delete(id int64) error {
	respBytes, err := zero_trust_client.Request(deletePath, "delete", map[string]interface{}{
		"id": id,
	})
	if err != nil {
		log.Warnf("删除应用账户失败：%s", err)
		return errors.New("删除应用账户失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return errors.New("删除应用账户失败：服务器异常！")
	}
	if resp.Code != 0 {
		return errors.New("删除应用账户失败：" + resp.Message)
	}
	return nil
}
