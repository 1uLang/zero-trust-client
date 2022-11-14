package user

import (
	"encoding/json"
	"errors"
	"github.com/1uLang/zero-trust-client"
	log "github.com/sirupsen/logrus"
)

const (
	editPasswordPath = ""
	enableOTPPath    = ""
	infoPath         = ""
	optPath          = ""
	editPath         = ""
)

type responsePublic struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type info struct {
	Header         string `json:"header"`           // 头像
	Name           string `json:"name"`             // 姓名
	Account        string `json:"account"`          // 账户
	Mobile         string `json:"mobile"`           // 手机号
	Email          string `json:"email"`            // 邮箱地址
	Device         string `json:"device"`           // 设备名称
	LastOnlineTime string `json:"last_online_time"` // 上次登录时间
}

// EditPassword 修改密码
func EditPassword(old, new string) error {
	respBytes, err := zero_trust_client.Request(editPasswordPath, "post", map[string]interface{}{
		"old": old,
		"new": new,
	})
	if err != nil {
		log.Warnf("修改密码失败：%s", err)
		return errors.New("修改密码失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return errors.New("修改密码失败：服务器异常！")
	}
	if resp.Code != 0 {
		return errors.New("修改密码失败：" + resp.Message)
	}
	return nil
}

// EnableOTP 开启双因素认证
func EnableOTP() error {
	respBytes, err := zero_trust_client.Request(enableOTPPath, "post", nil)
	if err != nil {
		log.Warnf("开启双因素认证失败：%s", err)
		return errors.New("开启双因素认证失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return errors.New("开启双因素认证失败：服务器异常！")
	}
	if resp.Code != 0 {
		return errors.New("开启双因素认证失败：" + resp.Message)
	}
	return nil
}

// OPT 获取opt二维码
func OPT() (string, error) {
	respBytes, err := zero_trust_client.Request(optPath, "get", map[string]interface{}{})
	if err != nil {
		log.Warnf("获取opt二维码失败：%s", err)
		return "", errors.New("获取opt二维码失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
		Data struct {
			OPT string `json:"opt"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return "", errors.New("获取opt二维码失败：服务器异常！")
	}
	if resp.Code != 0 {
		return "", errors.New("获取opt二维码失败：" + resp.Message)
	}
	return resp.Data.OPT, nil
}

// Info 获取个人信息
func Info() (*info, error) {
	respBytes, err := zero_trust_client.Request(infoPath, "get", map[string]interface{}{})
	if err != nil {
		log.Warnf("获取个人信息失败：%s", err)
		return nil, errors.New("获取个人信息失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
		Data struct {
			Info *info `json:"info"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return nil, errors.New("获取个人信息失败：服务器异常！")
	}
	if resp.Code != 0 {
		return nil, errors.New("获取个人信息失败：" + resp.Message)
	}
	return resp.Data.Info, nil
}

func Edit(header, name, mobile, email string) error {
	respBytes, err := zero_trust_client.Request(editPath, "post", map[string]interface{}{
		"header": header,
		"name":   name,
		"mobile": mobile,
		"email":  email,
	})
	if err != nil {
		log.Warnf("修改用户信息失败：%s", err)
		return errors.New("修改用户信息失败：服务器错误！")
	}
	resp := &struct {
		responsePublic
	}{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		log.Warnf("响应数据解析错误：%s", err)
		return errors.New("修改用户信息失败：服务器异常！")
	}
	if resp.Code != 0 {
		return errors.New("修改用户信息失败：" + resp.Message)
	}
	return nil
}
