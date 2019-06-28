package util

import (
	"net"
)

func GetLocalIP() ([]string, error) {
	ips := make([]string, 0)

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}
