// Scry Info.  All rights reserved.
// license that can be found in the license file.

package lb

import (
	"google.golang.org/grpc/balancer/roundrobin"
	"strings"

	"google.golang.org/grpc"
)

const (
	Round = "round"
	First = "first"
)

//If not found, then return Round Balance
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

func BalancerRound() grpc.DialOption { //todo grpc realization balance management use global variables, and do not consider multi thread condition, this is a temporary point
	return grpc.WithBalancerName(roundrobin.Name)
}

func BalancerFirst() grpc.DialOption { //todo grpc realization balance management use global variables, and do not consider multi thread condition, this is a temporary point
	return grpc.WithBalancerName(grpc.PickFirstBalancerName)
}
