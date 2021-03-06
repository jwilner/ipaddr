# ipaddr

Because I'm tired of reimplementing the same IP operations everywhere in Go.

Expressive IP utilities implemented with efficient, idiomatic lazy operations.

```go
import "github.com/jwilner/ipaddr"

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
```
