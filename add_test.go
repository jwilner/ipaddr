package ipaddr

import (
	"net"
	"testing"
)

func TestAdd(t *testing.T) {
	for _, tt := range []struct {
		name string
		ip   string
		val  int
		want string
	}{
		{"add one", "10.0.0.0", 1, "10.0.0.1"},
		{"add one", "9.255.255.255", 1, "10.0.0.0"},
		{"wrap", "10.0.255.255", 1, "10.1.0.0"},
		{"multi", "10.254.255.255", 2, "10.255.0.1"},
		{"none", "255.255.255.255", 1, ""},
		{"big", "0.0.0.0", 1 << 32, ""},

		{"sub one", "10.0.0.1", -1, "10.0.0.0"},
		{"sub one", "10.0.0.0", -1, "9.255.255.255"},
		{"sub one", "10.0.0.0", -3, "9.255.255.253"},
		{"sub one", "10.0.0.0", -259, "9.255.254.253"},
		{"sub one", "255.255.255.255", (-256 * 256 * 256 * 256) + 1, "0.0.0.0"},
		{"sub one", "255.255.255.255", -256 * 256 * 256 * 256, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ip := parseIP(t, tt.ip)
			var want net.IP
			if tt.want != "" {
				want = parseIP(t, tt.want)
			}
			got := Add(ip, tt.val)
			if (got == nil) != (want == nil) || (want != nil && !want.Equal(got)) {
				t.Fatalf("Wanted %v and got %v", want, got)
			}
		})
	}
}

func parseIP(t *testing.T, s string) net.IP {
	ip := net.ParseIP(s)
	if ip == nil {
		t.Fatalf("unable to parse IP: %q", s)
	}
	return ip
}
