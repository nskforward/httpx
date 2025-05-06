package httpx

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type IPTree struct {
	root  *ipTreeNode
	count int64
}

type ipTreeNode struct {
	left  *ipTreeNode // 0
	right *ipTreeNode // 1
	value any
}

func (tree *IPTree) Search(addr string) any {
	if tree.root == nil {
		return nil
	}
	ip := net.ParseIP(addr)
	if ip == nil {
		return nil
	}
	return tree.root.search(ip)
}

func (tree *IPTree) Count() int64 {
	return tree.count
}

func (tree *IPTree) Dump() {
	tree.root.dump(0)
}

func (tree *IPTree) Add(addr string, val any) error {
	if tree.root == nil {
		tree.root = &ipTreeNode{}
	}
	if strings.Contains(addr, "/") {
		_, cidr, err := net.ParseCIDR(addr)
		if err != nil {
			return err
		}
		err = tree.root.addCIDR(cidr, val)
		if err == nil {
			tree.count++
		}
		return err
	}
	ip := net.ParseIP(addr)
	if ip == nil {
		return fmt.Errorf("cannot parse ip")
	}
	err := tree.root.addIP(ip, val)
	if err == nil {
		tree.count++
	}
	return err
}

func (node *ipTreeNode) addIP(ip net.IP, value any) error {
	ipv4 := ip.To4()
	if ipv4 != nil {
		ip = ipv4
	}
	current := node
	for _, b := range ip {
		for bit := range 8 {
			bitFilled := isBitSet(b, bit)
			current = current.nextNode(bitFilled)
		}
	}
	if current.value != nil {
		return errors.New("ip adress already in use")
	}
	current.value = value
	return nil
}

func (node *ipTreeNode) addCIDR(cidr *net.IPNet, value any) error {
	current := node
loop:
	for i, b := range cidr.Mask {
		for bit := range 8 {
			if !isBitSet(b, bit) {
				if current.value != nil {
					return fmt.Errorf("cidr adress already in use")
				}
				current.value = value
				break loop
			}
			bitFilled := isBitSet(cidr.IP[i], bit)
			current = current.nextNode(bitFilled)
		}
	}
	return nil
}

func (node *ipTreeNode) search(ip net.IP) any {
	ipv4 := ip.To4()
	if ipv4 != nil {
		ip = ipv4
	}
	current := node
	for _, b := range ip {
		for bit := range 8 {
			bitFilled := isBitSet(b, bit)
			if bitFilled && current.right != nil {
				current = current.right
				if current.value != nil {
					return current.value
				}
				continue
			}
			if !bitFilled && current.left != nil {
				current = current.left
				if current.value != nil {
					return current.value
				}
				continue
			}
			return nil
		}
	}
	return nil
}

func (node *ipTreeNode) dump(level int) {
	if node.left != nil {
		fmt.Print(strings.Repeat(".", level))
		fmt.Print("0")
		if node.left.value != nil {
			fmt.Printf(" (%s)", node.left.value)
		}
		fmt.Println()
		node.left.dump(level + 1)
	}
	if node.right != nil {
		fmt.Print(strings.Repeat(".", level))
		fmt.Print("1")
		if node.right.value != nil {
			fmt.Printf(" (%s)", node.right.value)
		}
		fmt.Println()
		node.right.dump(level + 1)
	}
}

func (node *ipTreeNode) nextNode(isSet bool) *ipTreeNode {
	if isSet {
		if node.right == nil {
			node.right = &ipTreeNode{}
		}
		return node.right
	} else {
		if node.left == nil {
			node.left = &ipTreeNode{}
		}
		return node.left
	}
}

func isBitSet(b byte, bit int) bool {
	return b&(1<<(7-bit)) != 0
}
