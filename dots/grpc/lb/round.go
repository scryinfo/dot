// Scry Info.  All rights reserved.
// license that can be found in the license file.

package lb

import (
	"context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"math/rand"
	"sync"
	"time"
)

// Name is the name of round_robin balancer.
const Name = Round

// newBuilder creates a new roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilderWithConfig(Name, &rrPickerBuilder{r: rand.New(rand.NewSource(time.Now().UnixNano()))}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type rrPickerBuilder struct {
	r         *rand.Rand
	randMutex sync.Mutex
}

// Intn implements rand.Intn on the rrPickerBuilder source.
func (c *rrPickerBuilder) Intn(n int) int {
	c.randMutex.Lock()
	res := c.r.Intn(n)
	c.randMutex.Unlock()
	return res
}

func (c *rrPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	grpclog.Infof("roundrobinPicker: newPicker called with readySCs: %v", readySCs)
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for _, sc := range readySCs {
		scs = append(scs, sc)
	}
	return &rrPicker{
		subConns: scs,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		next: c.Intn(len(scs)),
	}
}

type rrPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn

	mu   sync.Mutex
	next int
}

func (c *rrPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	c.mu.Lock()
	sc := c.subConns[c.next]
	c.next = (c.next + 1) % len(c.subConns)
	c.mu.Unlock()
	return sc, nil, nil
}
