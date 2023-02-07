package network

import "net"

func CheckDomainIsAlive(domain string) bool {
	_, dnsResolvErr := net.LookupIP(domain)
	if dnsResolvErr == nil {
		return true
	}
	return false
}
