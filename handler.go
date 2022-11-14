package zero_trust_client

import (
	"errors"
	"fmt"
	"github.com/1uLang/libnet"
	"github.com/1uLang/libnet/utils/maps"
	log "github.com/sirupsen/logrus"
	"time"
)

type conn struct {
	c         *libnet.Connection
	ticker    *time.Ticker
	loginChan chan error
}

// 定时发送心跳包文
func (this conn) startKeepaliveTicker() {

	this.ticker = time.NewTicker(keepaliveTimeout / 2)

	for range this.ticker.C {
		if !this.c.IsClose() {
			msg := message{Type: KeepaliveRequestCode}
			this.c.Write(msg.Marshal())
		} else {
			return
		}
	}
}

// 处理消息
func (this *conn) onMessage(msg *message) {
	switch msg.Type {
	case LoginResponseCode:
		this.handleLoginAckMessage(msg)
	case AHListRequestCode:
		this.handleAHListMessage(msg)
	case CustomRequestCode:
		this.handleCustomMessage(msg)
	}
}

// 登陆消息
func (this *conn) login(m maps.Map) error {

	reply := &message{
		Type: LoginRequestCode,
		Data: m.AsJSON(),
	}
	_, err := this.c.Write(reply.Marshal())
	fmt.Println("==== send login message failed : ", err)
	return err
}

// AH/IH 处理登录响应
func (this *conn) handleLoginAckMessage(msg *message) {

	m, err := maps.DecodeJSON(msg.Data)
	if err != nil {
		return
	}
	switch m.GetInt8("code") {
	case 0: //登录成功 - 定时发送心跳
		// 开启定时任务 发送心跳报
		this.loginChan <- nil
		go this.startKeepaliveTicker()
		return
	case 1: //记录错误信息：
		log.Error("[SDP Client]authorize failed ", m.GetString("message"))
		this.loginChan <- errors.New("身份认证失败：账户或密码错误。")
		this.c.Close("authority failed : " + m.GetString("message"))
		//无效的认证凭证
	case 2:
		//限制登录
		log.Warn("[SDP Client]login limit ", m.GetString("message"))
		this.loginChan <- errors.New("已限制登录，请联系管理员。")
		this.c.Close("login limit " + m.GetString("message"))
	}
}

// 网关列表
func (this *conn) handleAHListMessage(msg *message) {

	go func() { // 存在堵塞情况 但是不影响 整条连接 故 丢到协程中处理
		gateway, err := getGateway(msg.Data)
		if err != nil {
			log.Warn("[SDP Client]get gateway failed : %s", err)
			return
		}
		setGateway(gateway)
		//todo: 对接wireguard 启动wireguard 客户端
	}()
}

// ah/ih 自定义消息（错误消息）
func (this *conn) handleCustomMessage(msg *message) {
	// todo：处理自定义消息
	log.Info("[SDP Client] custom message : ", msg.Data)
}

func (this *conn) logout() error {
	reply := &message{
		Type: IHLogoutRequestCode,
	}
	if this.c != nil {
		_, err := this.c.Write(reply.Marshal())
		if err != nil {
			log.Warnf("[SDP Client] send logout message failed :%s ", err)
		}
		defer this.c.Close("")
		return err
	}
	return nil
}
