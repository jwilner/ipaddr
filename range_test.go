package ipaddr

import (
	"bytes"
	"context"
	"net"
	"reflect"
	"testing"
)

func TestNext(t *testing.T) {
	for _, tt := range []struct {
		name string
		nt   string
		want string
	}{
		{"incrs", "10.0.0.0/16", "10.1.0.0/16"},
		{"wraps bit", "10.7.0.0/16", "10.8.0.0/16"},
		{"wraps byte", "10.255.0.0/16", "11.0.0.0/16"},
		{"out of bound", "0.0.0.0/0", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			nt, want := parseHelper(t, tt.nt, tt.want)
			got := Next(nt)
			checkWant(t, got, want)
		})
	}
}

func TestPrev(t *testing.T) {
	for _, tt := range []struct {
		name string
		nt   string
		want string
	}{
		{"decrs", "10.1.0.0/16", "10.0.0.0/16"},
		{"wraps bit", "10.8.0.0/16", "10.7.0.0/16"},
		{"wraps byte", "11.0.0.0/16", "10.255.0.0/16"},
		{"out of bound", "0.0.0.0/0", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			nt, want := parseHelper(t, tt.nt, tt.want)
			got := Prev(nt)
			checkWant(t, got, want)
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name       string
		start      string
		step, take int
		want       []string
	}{
		{"finds next", "10.0.0.0/8", 1, 3, []string{"10.0.0.0/8", "11.0.0.0/8", "12.0.0.0/8"}},
		{"steps next", "10.0.0.0/8", 2, 3, []string{"10.0.0.0/8", "12.0.0.0/8", "14.0.0.0/8"}},
		{"no more", "255.0.0.0/8", 1, 3, []string{"255.0.0.0/8"}},
		{"one more", "254.0.0.0/8", 1, 3, []string{"254.0.0.0/8", "255.0.0.0/8"}},

		{"finds next", "2001:db8::/35", 1, 3, []string{"2001:db8::/35", "2001:db8:2000::/35", "2001:db8:4000::/35"}},
		{"steps next", "2001:db8::/35", 2, 3, []string{"2001:db8::/35", "2001:db8:4000::/35", "2001:db8:8000::/35"}},
		{"no more", "ffff::/16", 1, 3, []string{"ffff::/16"}},
		{"one more", "fffe::/16", 1, 3, []string{"fffe::/16", "ffff::/16"}},

		{"finds next", "10.0.0.0/8", -1, 3, []string{"10.0.0.0/8", "9.0.0.0/8", "8.0.0.0/8"}},
		{"steps next", "10.0.0.0/8", -2, 3, []string{"10.0.0.0/8", "8.0.0.0/8", "6.0.0.0/8"}},
		{"no more", "0.0.0.0/8", -1, 3, []string{"0.0.0.0/8"}},
		{"one more", "1.0.0.0/8", -1, 3, []string{"1.0.0.0/8", "0.0.0.0/8"}},

		{"finds next", "f::/16", -1, 3, []string{"f::/16", "e::/16", "d::/16"}},
		{"steps next", "f::/16", -2, 3, []string{"f::/16", "d::/16", "b::/16"}},
		{"no more", "::/16", -1, 3, []string{"::/16"}},
		{"one more", "1::/16", -1, 3, []string{"1::/16", "::/16"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cncl := context.WithCancel(context.Background())
			defer cncl()

			_, cidr, err := net.ParseCIDR(tt.start)
			if err != nil {
				t.Fatalf("expected valid cidr: %v", err)
			}

			var found []string
			for _, f := range Range(ctx, *cidr, tt.step).Take(ctx, tt.take) {
				found = append(found, f.String())
			}

			if !reflect.DeepEqual(found, tt.want) {
				t.Errorf("NextSeq() = %v, want %v", found, tt.want)
			}
		})
	}
}

func eqNet(a, b net.IPNet) bool {
	return a.IP.Equal(b.IP) && bytes.Equal(a.Mask, b.Mask)
}

func checkWant(t *testing.T, got, want *net.IPNet) {
	if (got == nil) != (want == nil) || (got != nil && !eqNet(*got, *want)) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func parseHelper(t *testing.T, ttNt, ttWant string) (net.IPNet, *net.IPNet) {
	_, nt, err := net.ParseCIDR(ttNt)
	if err != nil {
		t.Fatalf("invalid nt: %v", err)
	}

	var want *net.IPNet
	if ttWant != "" {
		if _, want, err = net.ParseCIDR(ttWant); err != nil {
			t.Fatalf("invalid want: %v", err)
		}
	}

	return *nt, want
}
