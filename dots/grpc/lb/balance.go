// Scry Info.  All rights reserved.
// license that can be found in the license file.

package lb

import (
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

const (
	Round = "round"
	First = "first"
)

//如果没有找到，那么返回 Round Balance
func Balance(bname string) grpc.DialOption {
	var do grpc.DialOption = nil
	switch strings.ToLower(bname) {
	case Round:
		do = BalancerRound()
	case First:
		do = BalancerFirst()
	default:
		do = BalancerRound()
	}
	return do
}

func BalancerRound() grpc.DialOption { //todo grpc实现 balance管理里使用了 全局变量，且没有考虑多线程的情况， 这里是一个改时点
	return grpc.WithBalancerName(roundrobin.Name)
}

func BalancerFirst() grpc.DialOption { //todo grpc实现 balance管理里使用了 全局变量，且没有考虑多线程的情况， 这里是一个改时点
	return grpc.WithBalancerName(grpc.PickFirstBalancerName)
}
