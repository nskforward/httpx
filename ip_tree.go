package httpx

import (
	"errors"
	"net"
	"strings"
)

type IPTree struct {
	parent *IPTree
	b0     *IPTree
	b1     *IPTree
	value  any
}

var (
	ErrNodeAlreadyExists = errors.New("address already used")
	ErrBadInputFormat    = errors.New("bad address format")
)

func (ipTree *IPTree) Search(ip string) (any, error) {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return nil, ErrBadInputFormat
	}

	ipv4 := netIP.To4()
	if ipv4 != nil {
		netIP = ipv4
	}

	var i int
	bit := byte(0x80)
	node := ipTree

	for node != nil {
		if node.value != nil {
			return node.value, nil
		}

		if netIP[i]&bit != 0 {
			node = node.b1
		} else {
			node = node.b0
		}

		if bit >>= 1; bit == 0 {
			i, bit = i+1, byte(0x80)
			if i >= len(netIP) {
				if node != nil {
					return node.value, nil
				}
				break
			}
		}
	}
	return nil, nil
}

func (ipTree *IPTree) Append(addr string, value any) error {
	if strings.Contains(addr, "/") {
		return ipTree.appendCIDR(addr, value)
	}
	return ipTree.appendIP(addr, value)
}

func (ipTree *IPTree) appendIP(addr string, value any) error {
	ip := net.ParseIP(addr)
	if ip == nil {
		return ErrBadInputFormat
	}

	ipv4 := ip.To4()
	if ipv4 != nil {
		ip = ipv4
	}

	node := ipTree
	next := node
	bit := byte(0x80)
	i := 0

	for {
		if ip[i]&bit != 0 {
			next = node.b1
		} else {
			next = node.b0
		}
		if next == nil {
			break
		}
		node = next
		if bit >>= 1; bit == 0 {
			if i++; i == len(ip) {
				break
			}
			bit = byte(0x80)
		}
	}

	if next != nil {
		if node.value != nil {
			return ErrNodeAlreadyExists
		}
		node.value = value
		return nil
	}

	for {
		next = &IPTree{}
		next.parent = node
		if ip[i]&bit != 0 {
			node.b1 = next
		} else {
			node.b0 = next
		}
		node = next
		if bit >>= 1; bit == 0 {
			if i++; i == len(ip) {
				break
			}
			bit = byte(0x80)
		}
	}
	node.value = value
	return nil
}

func (ipTree *IPTree) appendCIDR(addr string, value any) error {
	_, network, err := net.ParseCIDR(addr)
	if err != nil {
		return err
	}

	node := ipTree
	next := node
	bit := byte(0x80)
	i := 0

	for bit&network.Mask[i] != 0 {
		if network.IP[i]&bit != 0 {
			next = node.b1
		} else {
			next = node.b0
		}
		if next == nil {
			break
		}
		node = next
		if bit >>= 1; bit == 0 {
			if i++; i == len(network.IP) {
				break
			}
			bit = byte(0x80)
		}
	}

	if next != nil {
		if node.value != nil {
			return ErrNodeAlreadyExists
		}
		node.value = value
		return nil
	}

	for bit&network.Mask[i] != 0 {
		next = &IPTree{}
		next.parent = node
		if network.IP[i]&bit != 0 {
			node.b1 = next
		} else {
			node.b0 = next
		}
		node = next
		if bit >>= 1; bit == 0 {
			if i++; i == len(network.IP) {
				break
			}
			bit = byte(0x80)
		}
	}
	node.value = value
	return nil
}
