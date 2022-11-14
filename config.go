package zero_trust_client

import "time"

const (
	sdpControlPort   = 20000
	connectTimeout   = 3 * time.Second
	keepaliveTimeout = 30 * time.Second
)

// 设置的控制器IP
var controlIP string
