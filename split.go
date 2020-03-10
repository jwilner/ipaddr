package ipaddr

import (
	"context"
	"math/big"
	"net"
)

const (
	ipv4MaskLen = net.IPv4len * 8
	ipv6MaskLen = net.IPv6len * 8
)

// Split splits the provided network according to the provided mask and returns a lazy sequence of the resulting nets.
// Argument mismatches or a mask less than that of the original network will result in an empty sequence.
func Split(ctx context.Context, orig net.IPNet, mask net.IPMask) NetSeq {
	ch := make(chan net.IPNet)

	prefix, maskLen := orig.Mask.Size()
	targetPrefix, targetLen := mask.Size()

	switch {
	case targetPrefix < prefix, targetLen != maskLen: // skip invalid
		close(ch)

	case maskLen == ipv4MaskLen:
		go split4(ctx, orig.IP.To4(), prefix, targetPrefix, ch)

	case maskLen == ipv6MaskLen:
		go split6(ctx, orig.IP.To16(), prefix, targetPrefix, ch)

	default:
		close(ch)
	}

	return ch
}

func split4(ctx context.Context, ip net.IP, prefix, targetPrefix int, ch chan<- net.IPNet) {
	defer close(ch)

	cidrMask := net.CIDRMask(targetPrefix, ipv4MaskLen)
	count := uint32(1 << (targetPrefix - prefix))
	offset := ipv4MaskLen - targetPrefix
	for i := uint32(0); i < count; i++ {
		v := i << offset
		select {
		case <-ctx.Done():
			return
		case ch <- net.IPNet{
			IP: net.IP{
				ip[0] | uint8(v>>24),
				ip[1] | uint8(v>>16),
				ip[2] | uint8(v>>8),
				ip[3] | uint8(v),
			},
			Mask: cidrMask,
		}:
		}
	}
}

func split6(ctx context.Context, ip net.IP, prefix, targetPrefix int, ch chan<- net.IPNet) {
	defer close(ch)

	one := big.NewInt(1)
	rng := new(big.Int).Lsh(one, uint(targetPrefix-prefix))
	offset := uint(ipv6MaskLen - targetPrefix)
	cidrMask := net.CIDRMask(targetPrefix, ipv6MaskLen)

	for i := big.NewInt(0); i.Cmp(rng) < 0; i.Add(i, one) {
		dist := new(big.Int).Lsh(i, offset).Bytes()

		newIP := make(net.IP, net.IPv6len)
		diff := net.IPv6len - len(dist)
		copy(newIP, ip[:diff])
		for j := diff; j < net.IPv6len; j++ {
			newIP[j] = ip[j] | dist[j-diff]
		}

		select {
		case <-ctx.Done():
			return
		case ch <- net.IPNet{IP: newIP, Mask: cidrMask}:
		}
	}
}
