package ipaddr

import (
	"context"
	"net"
)

// NetSeq is a lazy sequence of net.IPNet values
type NetSeq <-chan net.IPNet

// Next returns the next IPNet if there is one otherwise nil
func (n NetSeq) Next(ctx context.Context) *net.IPNet {
	select {
	case <-ctx.Done():
	case nt, ok := <-n:
		if ok {
			return &nt
		}
	}
	return nil
}

// All returns all of the nets in the sequence or stops on context end.
func (n NetSeq) All(ctx context.Context) (res []net.IPNet) {
	for {
		select {
		case <-ctx.Done():
			return
		case nt, ok := <-n:
			if !ok {
				return
			}
			res = append(res, nt)
		}
	}
}

// Take returns as many as `num` of the nets in the sequence or stops on context end.
func (n NetSeq) Take(ctx context.Context, num int) (res []net.IPNet) {
	res = make([]net.IPNet, 0, num)
	for i := 0; i < num; i++ {
		select {
		case <-ctx.Done():
			return
		case nt, ok := <-n:
			if !ok {
				return
			}
			res = append(res, nt)
		}
	}
	return
}
