package net

import (
	"net"
)

type LAN struct {
	ipNets []*net.IPNet
}

func NewLAN() LAN {
	var lan LAN
	lanCIDRs := []string{"10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"}
	for _, cidr := range lanCIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		lan.ipNets = append(lan.ipNets, ipNet)
	}
	return lan
}

func (l LAN) In(ip string) bool {
	for _, ipNet := range l.ipNets {
		in := ipNet.Contains(net.ParseIP(ip))
		if in {
			return true
		}
	}
	return false
}

// return LAN ip, mac
func (l LAN) NetInfo() (string, string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, i := range ifaces {
		mac := i.HardwareAddr.String()
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address
			_ip := ip.String()
			if l.In(_ip) {
				return _ip, mac
			}
		}
	}
	panic("get net info fail")
}

