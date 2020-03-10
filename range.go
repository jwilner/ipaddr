package ipaddr

import (
	"context"
	"math/big"
	"net"
)

// Next returns the greater and adjacent IPNet or nil if none exists.
func Next(nt net.IPNet) *net.IPNet {
	var ip net.IP
	if ntIP := nt.IP.To4(); ntIP != nil {
		ip = ntIP
	} else {
		ip = make(net.IP, len(nt.IP))
	}
	copy(ip, nt.IP)

	ones, bits := nt.Mask.Size()
	mask := net.CIDRMask(ones, bits)

	res := net.IPNet{IP: ip, Mask: mask}
	for i := ones - 1; i >= 0; i-- {
		a, b := i/8, 7-(i%8)
		if res.IP[a]&(1<<b) == 0 {
			res.IP[a] |= 1 << b
			return &res
		}
		res.IP[a] &^= 1 << b
	}
	return nil
}

// Prev returns the lesser and adjacent IPNet or nil if none exists.
func Prev(nt net.IPNet) *net.IPNet {
	ip := make(net.IP, len(nt.IP))
	copy(ip, nt.IP)

	ones, bits := nt.Mask.Size()
	mask := net.CIDRMask(ones, bits)

	res := net.IPNet{IP: ip, Mask: mask}
	for i := ones - 1; i >= 0; i-- {
		a, b := i/8, 7-(i%8)
		if res.IP[a]&(1<<b) != 0 {
			res.IP[a] &^= 1 << b
			return &res
		}
		res.IP[a] |= 1 << b
	}
	return nil
}

// Range returns the sequence of nets with a given mask starting with `start` and stepping by `step`. If arguments are
// invalid, no networks are provided.
func Range(ctx context.Context, start net.IPNet, step int) NetSeq {
	ch := make(chan net.IPNet)

	ones, bits := start.Mask.Size()
	if ip := start.IP.To4(); ip != nil && bits == ipv4MaskLen && step != 0 {
		if step > 0 {
			go nextSeq4(ctx, ip, ones, step, ch)
		} else {
			go prevSeq4(ctx, ip, ones, step, ch)
		}
	} else if ip = start.IP.To16(); ip != nil && bits == ipv6MaskLen && step != 0 {
		if step > 0 {
			go nextSeq6(ctx, ip, ones, step, ch)
		} else {
			go prevSeq6(ctx, ip, ones, step, ch)
		}
	} else {
		close(ch)
	}
	return ch
}

func prevSeq6(ctx context.Context, ip net.IP, ones, step int, ch chan net.IPNet) {
	defer close(ch)

	cur := new(big.Int).SetBytes(ip)
	cur.Rsh(cur, uint(ipv6MaskLen-ones))

	s := big.NewInt(int64(step))
	s.Abs(s)

	mask := net.CIDRMask(ones, ipv6MaskLen)
	for ; cur.Sign() != -1; cur.Sub(cur, s) {
		dist := new(big.Int).Lsh(cur, uint(ipv6MaskLen-ones)).Bytes()
		newIP := make(net.IP, net.IPv6len)

		copy(newIP[net.IPv6len-len(dist):], dist)

		select {
		case <-ctx.Done():
			return
		case ch <- net.IPNet{IP: newIP, Mask: mask}:
		}
	}
}

func prevSeq4(ctx context.Context, ip net.IP, ones, step int, ch chan<- net.IPNet) {
	defer close(ch)

	mask := net.CIDRMask(ones, ipv4MaskLen)

	cur := int64(uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3]))
	cur >>= ipv4MaskLen - ones

	s := int64(step)

	for ; cur >= 0; cur += s {
		v := cur << (32 - ones)
		select {
		case <-ctx.Done():
			return
		case ch <- net.IPNet{
			IP: net.IP{
				uint8(v >> 24),
				uint8(v >> 16),
				uint8(v >> 8),
				uint8(v),
			},
			Mask: mask,
		}:
		}
	}
}

func nextSeq6(ctx context.Context, ip net.IP, ones, step int, ch chan net.IPNet) {
	defer close(ch)

	cur := new(big.Int).SetBytes(ip)
	cur.Rsh(cur, uint(ipv6MaskLen-ones))

	max := big.NewInt(1)
	max.Lsh(max, uint(ones))

	s := big.NewInt(int64(step))

	mask := net.CIDRMask(ones, ipv6MaskLen)
	for ; cur.Cmp(max) < 0; cur.Add(cur, s) {
		dist := new(big.Int).Lsh(cur, uint(ipv6MaskLen-ones)).Bytes()
		newIP := make(net.IP, net.IPv6len)

		copy(newIP[net.IPv6len-len(dist):], dist)

		select {
		case <-ctx.Done():
			return
		case ch <- net.IPNet{IP: newIP, Mask: mask}:
		}
	}
}

func nextSeq4(ctx context.Context, ip net.IP, ones, step int, ch chan<- net.IPNet) {
	defer close(ch)

	mask := net.CIDRMask(ones, ipv4MaskLen)

	cur := uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
	cur >>= ipv4MaskLen - ones

	s := uint32(step)

	for max := uint32(1 << ones); cur < max; cur += s {
		v := cur << (32 - ones)
		select {
		case <-ctx.Done():
			return
		case ch <- net.IPNet{
			IP: net.IP{
				uint8(v >> 24),
				uint8(v >> 16),
				uint8(v >> 8),
				uint8(v),
			},
			Mask: mask,
		}:
		}
	}
}
