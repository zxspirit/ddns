package local

import (
	"net"
)

// getPublicIp 获取本地的公网IPv6地址
func getPublicIp() (ipv4s []string, ipv6s []string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ipv4s, ipv6s
	}
	// https://en.wikipedia.org/wiki/IPv6_address#General_allocation
	//_, ipv6Unicast, _ := net.ParseCIDR("2000::/3")
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue // 跳过未启用或回环接口
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue // 如果获取地址失败，跳过该接口
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok &&
				ipnet.IP.IsGlobalUnicast() &&
				!ipnet.IP.IsLoopback() &&
				!ipnet.IP.IsPrivate() {
				_, bits := ipnet.Mask.Size()
				if bits == 128 {
					ipv6s = append(ipv6s, ipnet.IP.String())
				}
				if bits == 32 {
					ipv4s = append(ipv4s, ipnet.IP.String())
				}
			}
		}
	}
	return ipv4s, ipv6s
}
