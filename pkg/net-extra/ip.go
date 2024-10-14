package net_extra

import "net"

// GetOutboundIP 获取本机外网IP
func GetOutboundIP() string {
	// 114.114.114.114
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
