package utils

import (
	"errors"
	"net"
)

var publicIp net.IP

// GetExternalIP 获取ip
func GetExternalIP() (net.IP, error) {
	if publicIp != nil {
		return publicIp, nil
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			publicIp = GetIpFromAddr(addr)
			if publicIp == nil {
				continue
			}
			return publicIp, nil
		}
	}
	return nil, errors.New("connected to the network")
}

// GetIpFromAddr 获取ip
func GetIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
