package ip

import (
	"fmt"
	"net"
)

type GetIp interface {
	GetPublicIp() (ipv4 string, ipv6 string, err error)
}

type LocalInterface struct {
}

func (l LocalInterface) GetPublicIp() (ipv4 string, ipv6 string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		err = fmt.Errorf("error getting interfaces: %w", err)
		return
	}
	var ipv4s []net.IP
	var ipv6s []net.IP
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue // 跳过未启用或回环接口
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() &&
				!ipnet.IP.IsPrivate() && ipnet.IP.IsGlobalUnicast() {
				if ip := ipnet.IP.To4(); ip != nil {
					ipv4s = append(ipv4s, ip)
				} else if ip := ipnet.IP.To16(); ip != nil {
					//if ip[8]&0x40 == 0 {
					//	continue
					//}
					ipv6s = append(ipv6s, ip)
				}
			}
		}

	}
	if len(ipv4s) == 0 {
		ipv4 = ""
	} else {
		ipv4 = ipv4s[0].String()
	}
	if len(ipv6s) == 0 {
		ipv6 = ""
	} else {
		ipv6 = ipv6s[0].String()
	}
	fmt.Printf("Public IPv4s: %s, Public IPv6s: %s\n", ipv4s, ipv6s)
	return
}
