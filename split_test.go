package ipaddr

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name string
		orig string
		mask net.IPMask
		want []string
	}{
		{"target cannot be less", "10.0.0.0/16", net.CIDRMask(15, 32), []string{}},
		{"target wrong length", "10.0.0.0/16", net.CIDRMask(15, 128), []string{}},

		{"no change 4", "10.0.0.0/16", net.CIDRMask(16, 32), []string{"10.0.0.0/16"}},
		{"splits 4", "10.1.0.0/16", net.CIDRMask(18, 32), []string{"10.1.0.0/18", "10.1.64.0/18", "10.1.128.0/18", "10.1.192.0/18"}},

		{"no change 6", "::/0", net.CIDRMask(0, 128), []string{"::/0"}},
		{"splits 6", "::/0", net.CIDRMask(2, 128), []string{"::/2", "4000::/2", "8000::/2", "c000::/2"}},
		{
			"splits 6",
			"2001:db8::/32",
			net.CIDRMask(35, 128),
			[]string{
				"2001:db8::/35",
				"2001:db8:2000::/35",
				"2001:db8:4000::/35",
				"2001:db8:6000::/35",
				"2001:db8:8000::/35",
				"2001:db8:a000::/35",
				"2001:db8:c000::/35",
				"2001:db8:e000::/35",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, orig, err := net.ParseCIDR(tt.orig)
			if err != nil {
				t.Fatalf("orig: %v", err)
			}

			ctx, cncl := context.WithCancel(context.Background())
			defer cncl()

			want := parseNetSlice(t, tt.want)

			i := 0
			for n := range Split(ctx, *orig, tt.mask) {
				if i >= len(want) {
					t.Fatalf("Wanted %d nets, but had more starting with %v", len(want), n.String())
				}
				if w := want[i]; !n.IP.Equal(w.IP) || !bytes.Equal(n.Mask, w.Mask) {
					t.Fatalf("nets %d unequal -- wanted %v, got %v", i, w.String(), n.String())
				}
				i++
			}
			if i < len(want) {
				t.Fatalf("Wanted %d nets, but only got %d", len(want), i)
			}
		})
	}

	t.Run("zero", func(t *testing.T) {
		if r, ok := <-Split(context.Background(), net.IPNet{}, nil); ok {
			t.Fatalf("wanted nil but got %v", r)
		}
	})
}

func parseNetSlice(t *testing.T, nets []string) []net.IPNet {
	want := make([]net.IPNet, 0, len(nets))
	for i, w := range nets {
		_, n, err := net.ParseCIDR(w)
		if err != nil {
			t.Fatalf("nets[%d]: %v", i, err)
		}
		want = append(want, *n)
	}
	return want
}
