package zero_trust_client

import "github.com/1uLang/zero-trust-client/utils"

var managerSvrAddr string
var header map[string]string

// Request 向管理后台发送请求
func Request(path, method string, params map[string]interface{}) (response []byte, err error) {
	return utils.HTTPRequest(managerSvrAddr+path, method, params, header)
}
