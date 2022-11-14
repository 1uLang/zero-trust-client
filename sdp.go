package zero_trust_client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/1uLang/libnet"
	message2 "github.com/1uLang/libnet/message"
	"github.com/1uLang/libnet/options"
	"github.com/1uLang/libnet/utils/maps"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

const (
	connection_type = iota
	connection_client
	connection_gateway

	clientLoginTimeout = 3 * time.Second
)

// 与控制器建立 mtl 长连接

type client struct {
	clt    *libnet.Client
	c      *conn
	locker sync.Locker
}

var clt = client{}

func (this client) OnConnect(c *libnet.Connection) {
	log.Info("[SDP Client] connect control success")
	// 注册 - 客户端上线
	//todo：加入注册失败机制
	this.locker.Lock()
	if this.c != nil {
		this.c.c.Close("")
	}
	this.c = &conn{c: c}
	this.locker.Unlock()
	if err := this.c.login(maps.Map{"type": connection_client}); err != nil {
		log.Fatal("[SDP Client] login failed : ", err)
		c.Close(err.Error())
		return
	}
	// setup buffer
	clientBuffer := message2.NewBuffer(checkHeader)
	clientBuffer.OptValidateId = true
	clientBuffer.OnMessage(func(msg message2.MessageI) {
		this.c.onMessage((msg).(*message))
	})
	if err := c.SetBuffer(clientBuffer); err != nil {
		log.Fatal("[SDP Client] set message buffer failed : ", err)
		c.Close(err.Error())
		return
	}
}

func (this client) OnMessage(c *libnet.Connection, bytes []byte) {
}

func (this client) OnClose(c *libnet.Connection, reason string) {
	log.Error("[SDP Control] connection close : ", reason)
}

// ConnectControl 连接控制器
func ConnectControl() error {
	if controlIP == "" {
		return errors.New("please send to control spa")
	}
	return runClient()
}

func runClient() error {
	addr := fmt.Sprintf("%s:%d", controlIP, sdpControlPort)
	cert, err := os.ReadFile("sdp.ca")
	if err != nil {
		log.Fatalf("could not open certificate file: %v", err)
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certificate, err := tls.LoadX509KeyPair("sdp.cert", "sdp.key")
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
		return err
	}

	// Create a tls client and supply the created CA pool and certificate
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		ClientCAs:          caCertPool,
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}
	clt.clt, err = libnet.NewClient(addr, clt,
		options.WithTimeout(connectTimeout))

	return clt.clt.DialTLS(tlsConfig)
}

// Login 登录 - 用户身份检测
func Login(body maps.Map) error {
	// 封包 向控制器发送客户端上线包
	if clt.c != nil {
		err := clt.c.login(body)
		if err != nil {
			return err
		}
		clt.c.loginChan = make(chan error, 1)
		timeout := time.NewTimer(clientLoginTimeout)
		select {
		case <-timeout.C:
			return errors.New("服务器响应超时，请重试")
		case err = <-clt.c.loginChan:
			return err
		}
	}
	return errors.New("wait fot connect sdp control")
}

// Logout 退出
func Logout() error {
	// todo 关闭wireguard
	if clt.c != nil {
		return clt.c.logout()
	}
	return nil
}
