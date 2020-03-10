package ipaddr_test

import (
	"context"
	"fmt"
	"github.com/jwilner/ipaddr"
	"net"
)

func ExampleSplit() {
	ctx, cncl := context.WithCancel(context.Background())
	defer cncl()

	nt := net.IPNet{IP: net.IP{10, 0, 0, 0}, Mask: net.CIDRMask(16, 32)}

	for sub := range ipaddr.Split(ctx, nt, net.CIDRMask(18, 32)) {
		fmt.Println(sub.String())
	}

	// Output:
	// 10.0.0.0/18
	// 10.0.64.0/18
	// 10.0.128.0/18
	// 10.0.192.0/18
}


func ExampleRange() {
	ctx, cncl := context.WithCancel(context.Background())
	defer cncl()

	nt := net.IPNet{IP: net.IP{10, 0, 0, 0}, Mask: net.CIDRMask(16, 32)}

	for _, n := range ipaddr.Range(ctx, nt, 2).Take(ctx, 5) {
		fmt.Println(n.String())
	}

	// Output:
	// 10.0.0.0/16
	// 10.2.0.0/16
	// 10.4.0.0/16
	// 10.6.0.0/16
	// 10.8.0.0/16
}

func ExampleAdd() {
	fmt.Println(ipaddr.Add(net.IPv4(10, 0, 0, 1), -2))
	fmt.Println(ipaddr.Add(net.IPv4(10, 0, 0, 1), 3))

	// Output:
	// 9.255.255.255
	// 10.0.0.4
}

func ExampleNext() {
	res := ipaddr.Next(net.IPNet{IP: net.IP{10, 0, 0, 0}, Mask: net.CIDRMask(16, 32)})
	fmt.Println(res)

	// Output:
	// 10.1.0.0/16
}

func ExamplePrev() {
	res := ipaddr.Prev(net.IPNet{IP: net.IP{10, 0, 0, 0}, Mask: net.CIDRMask(16, 32)})
	fmt.Println(res)

	// Output:
	// 9.255.0.0/16
}
