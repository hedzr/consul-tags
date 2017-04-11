package util

import (
	"fmt"
	"net"
	"os"
)

func ParseHostOrIp(s string, ipv4 bool) (net.IP, error) {
	addrs, err := net.LookupHost(s)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return net.IP{}, err
	}

	var ipx = net.IPv4zero
	for _, a := range addrs {
		if ip := net.ParseIP(a); ip != nil && !ip.IsLoopback() {
			if ipv4 {
				if ip.To4() != nil {
					//os.Stdout.WriteString(ip.String() + "\n")
					if ipx.IsUnspecified() {
						ipx = ip //.String()
					}
				}
			} else {
				if ip.To16() != nil {
					//os.Stdout.WriteString(ip.String() + "\n")
					if ipx.IsUnspecified() {
						ipx = ip //.String()
					}
				}
			}
		}
	}
	return ipx, nil
}

func GetLanIp(ipv4 bool) (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		return "", err //os.Exit(1)
	}

	ip := ""
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipv4 {
				if ipnet.IP.To4() != nil {
					os.Stdout.WriteString(ipnet.IP.String() + "\n")
					if ip == "" {
						ip = ipnet.IP.String()
					}
				}
			} else {
				if ipnet.IP.To16() != nil {
					os.Stdout.WriteString(ipnet.IP.String() + "\n")
					if ip == "" {
						ip = ipnet.IP.String()
					}
				}
			}
		}
	}
	return ip, nil
}

func GetLanIp_v2(ipv4 bool) (string, error) {
	ipx := ""

	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return "", err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return "", err
	}

	for _, a := range addrs {
		if ip := net.ParseIP(a); ip != nil && !ip.IsLoopback() {
			if ipv4 {
				if ip.To4() != nil {
					os.Stdout.WriteString(ip.String() + "\n")
					if ipx == "" {
						ipx = ip.String()
					}
				}
			} else {
				if ip.To16() != nil {
					os.Stdout.WriteString(ip.String() + "\n")
					if ipx == "" {
						ipx = ip.String()
					}
				}
			}
		}
	}
	return ipx, nil
}
