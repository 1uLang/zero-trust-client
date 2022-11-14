package zero_trust_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// 当前正在使用的网关节点
var use_gateway *gatewayCfg
var use_gateway_locker sync.Locker

type gatewayCfg struct {
	IP           string `json:"ip"`
	SPAPort      int    `json:"spa_port"`
	Port         int    `json:"port"`
	WireguardCfg string `json:"wireguard_cfg"`
}

// 返回多个网关 如何选择最优的网关节点 以及如何检测错误的网关节点
func getGateway(data []byte) (*gatewayCfg, error) {
	gateways := []*gatewayCfg{}
	if err := json.Unmarshal(data, &gateways); err != nil {
		return nil, err
	}
	if len(gateways) == 0 {
		return nil, errors.New("无可用的网关节点")
	}
	setGatewayPool(gateways)
	return get1()
}

// 检测指定IP，port tcp端口是否能访问
func checkPort(ip string, port int) (bool, error) {
	c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 3*time.Second)
	if err != nil {
		return false, err
	}
	if c != nil {
		defer c.Close()
	}
	return c != nil, nil
}

// 设置当前生效的gateway
func setGateway(node *gatewayCfg) {
	use_gateway_locker.Lock()
	defer use_gateway_locker.Unlock()
	use_gateway = node
}
