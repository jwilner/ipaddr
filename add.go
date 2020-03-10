package ipaddr

import "net"

// Add adds `val` to the provided IP address
func Add(ip net.IP, val int) net.IP {
	if i := ip.To4(); i != nil {
		ip = i
	}

	sum := make(net.IP, len(ip))
	copy(sum, ip)

	if val < 0 {
		dist := val * -1
		for i := len(sum) - 1; i >= 0; i-- {
			a, b := dist/256, dist%256

			v := int(sum[i]) - b
			if v < 0 {
				v += 256
				a += 1
			}
			sum[i] = byte(v)

			if a == 0 {
				return sum
			}
			dist = a
		}
		return nil
	}

	for i := len(sum) - 1; i >= 0; i-- {
		added := int(sum[i]) + val
		if added > 255 {
			sum[i] = byte(added % 256)
			val = added / 256
			continue
		}
		sum[i] = uint8(added)
		return sum
	}

	return nil
}
