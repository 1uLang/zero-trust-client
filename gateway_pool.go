package zero_trust_client

import (
	"errors"
	"github.com/1uLang/zero-trust-client/utils/rands"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// 网关池

// 全局可用的网关
var gl_all_gateways []*gatewayCfg
var gl_all_gateway_locker sync.Locker

// 定时循环体 - 定时检测每个网关 发送spa 确保切换网关更加快捷
func loop() {
	gl_all_gateway_locker.Lock()
	defer gl_all_gateway_locker.Unlock()

	for _, gateway := range gl_all_gateways {
		err := SPA(gateway.IP, gateway.SPAPort)
		if err != nil {
			log.Warnf("[SPA Loop] send to %s:%d spa error:%s", gateway.IP, gateway.Port, err)
		}
	}
}

// 方式1 所有网关都发spa 选择一个认证通过最快的
func get1() (*gatewayCfg, error) {
	var node *gatewayCfg
	quitChan := make(chan bool)
	var locker sync.Locker
	gl_all_gateway_locker.Lock()
	for _, gateway := range gl_all_gateways {
		go func(gateway *gatewayCfg) {
			// 先发spa  打开端口 后检测 端口是否打开。 返回最快生效的网关节点
			err := SPA(gateway.IP, gateway.SPAPort)
			if err != nil {
				log.Warn("send gateway spa failed:", err)
				return
			}
			// 验证端口
			ok, err := checkPort(gateway.IP, gateway.Port)
			if err != nil {
				log.Warnf(" gateway %s:%d connect error: %s", gateway.IP, gateway.Port, err)
				return
			}
			if !ok {
				log.Warnf(" gateway %s:%d not be connect : %s", gateway.IP, gateway.Port, err)
				return
			}
			locker.Lock()
			if node == nil {
				node = gateway
				quitChan <- true
			}
			locker.Unlock()
		}(gateway)
	}
	gl_all_gateway_locker.Unlock()
	timeout := time.NewTimer(3 * time.Second)
	select {
	case <-quitChan:
		return node, nil
	case <-timeout.C:
		return nil, errors.New("not available connect gateway timeout")
	}
}

// 方式2 ： 随机找一个网关节点 并发送spa报
func get2() (*gatewayCfg, error) {
	gl_all_gateway_locker.Lock()
	node := gl_all_gateways[rands.Int(0, len(gl_all_gateways)-1)]
	gl_all_gateway_locker.Unlock()
	return node, SPA(node.IP, node.SPAPort)
}

// 设置全局可用的所有网关
func setGatewayPool(gateways []*gatewayCfg) {
	gl_all_gateway_locker.Lock()
	defer gl_all_gateway_locker.Unlock()
	gl_all_gateways = gateways
}
