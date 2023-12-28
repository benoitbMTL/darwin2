package utils

import (
	"math/rand"
	"net"
)

func GenerateRandomPublicIP() string {

	for {
		ip := net.IPv4(
			byte(rand.Intn(256)),
			byte(rand.Intn(256)),
			byte(rand.Intn(256)),
			byte(rand.Intn(256)),
		)

		if isPublicIP(ip) {
			return ip.String()
		}
	}
}

func isPublicIP(ip net.IP) bool {
	// Check if the IP is private, loopback, multicast, link-local, or in reserved ranges
	if ip.IsPrivate() || ip.IsLoopback() || ip.IsMulticast() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return false
	}

	// Check for specific reserved ranges
	reservedRanges := []string{
		"0.0.0.0/8",       // Current network (only valid as source address)
		"10.0.0.0/8",      // Private network
		"100.64.0.0/10",   // Shared Address Space
		"127.0.0.0/8",     // Loopback
		"169.254.0.0/16",  // Link Local
		"172.16.0.0/12",   // Private network
		"192.0.0.0/24",    // IETF Protocol Assignments
		"192.0.2.0/24",    // TEST-NET-1, documentation and examples
		"192.88.99.0/24",  // 6to4 Relay Anycast
		"192.168.0.0/16",  // Private network
		"198.18.0.0/15",   // Network benchmark tests
		"198.51.100.0/24", // TEST-NET-2, documentation and examples
		"203.0.113.0/24",  // TEST-NET-3, documentation and examples
		"224.0.0.0/4",     // IP multicast (former Class D network)
		"240.0.0.0/4",     // Reserved (former Class E network)
	}

    for _, cidr := range reservedRanges {
        _, reservedNet, _ := net.ParseCIDR(cidr)
        if reservedNet.Contains(ip) {
            return false
        }
    }

	return true
}
