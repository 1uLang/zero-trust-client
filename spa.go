package zero_trust_client

import (
	"errors"
	"github.com/1uLang/libspa"
	"github.com/1uLang/zero-trust-client/utils"
	"net"
	"time"

	libspaclt "github.com/1uLang/libspa/client"
)

const (
	method = ""
	key    = ""
	iv     = ""
)

var publicIp net.IP
var deviceId string

// SPA 发送spa 认证包
/*
	参数说明：
		ip 		： 控制器 IP
		port 	： 控制器 spa 服务端口
	功能描述：
		向控制器/网关发送spa报文，开启控制器认证端口/网关业务端口
*/
func SPA(ip string, port int) error {
	var err error

	controlIP = ip
	clt := libspaclt.New()
	clt.Port = port
	clt.Addr = ip
	clt.Protocol = "udp"
	clt.Method = method
	clt.KEY = key
	clt.IV = iv
	publicIp, err = utils.GetExternalIP()
	if err != nil {
		return errors.New("获取本地IP失败：" + err.Error())
	}
	deviceId, err = utils.GetDeviceId()
	if err != nil {
		return errors.New("获取设备ID失败：" + err.Error())
	}

	return clt.Send(&libspa.Body{
		ClientDeviceId: deviceId,
		ClientPublicIP: publicIp,
		ServerPublicIP: net.ParseIP(clt.Addr),
	})
}

// CheckResult 检测与控制器的连通性
func CheckResult() error {
	timeout := time.NewTimer(connectTimeout)
	isConnect := make(chan bool, 1)
	go func() {
		for !timeout.Stop() {
			if clt.clt != nil {
				isConnect <- true
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	select {
	case <-timeout.C:
		return errors.New("认证失败，请重试")
	case <-isConnect:
		return nil
	}
}
