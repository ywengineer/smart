package mr_smart

import (
	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/ywengineer/mr.smart/utility"
	"go.uber.org/zap"
	"sync/atomic"
)

// LoadBalance sets the load balancing method.
type LoadBalance int

const (
	// Random requests that connections are randomly distributed.
	Random LoadBalance = iota
	// Hash. requests that connections are bind to a fixed pool.
	Hash
	RoundRobin
)

// loadBalance sets the load balancing method for []*Pool
type loadBalance interface {
	LoadBalance() LoadBalance
	// Choose the most qualified Pool
	Pick(id int) Worker
}

func parseLoadBalance(lb string) LoadBalance {
	switch lb {
	case "random":
		return Random
	case "hash":
		return Hash
	case "rr":
		return RoundRobin
	}
	utility.DefaultLogger().Warn("unknown load balance, default to RoundRobin", zap.String("lb", lb))
	return RoundRobin
}

func newLoadBalance(lb LoadBalance, pools []Worker) loadBalance {
	switch lb {
	case Random:
		return newRandomLB(pools)
	case Hash:
		return newHashLB(pools)
	case RoundRobin:
		return newRoundRobinLB(pools)
	}
	return newRoundRobinLB(pools)
}

// randomLB
func newRandomLB(pools []Worker) loadBalance {
	return &randomLB{pools: pools, poolSize: len(pools)}
}

type randomLB struct {
	pools    []Worker
	poolSize int
}

func (b *randomLB) LoadBalance() LoadBalance {
	return Random
}

func (b *randomLB) Pick(id int) Worker {
	idx := fastrand.Intn(b.poolSize)
	return b.pools[idx]
}

// hashLB
func newHashLB(pools []Worker) loadBalance {
	return &hashLB{pools: pools, poolSize: len(pools)}
}

type hashLB struct {
	pools    []Worker
	poolSize int
}

func (b *hashLB) LoadBalance() LoadBalance {
	return Hash
}

func (b *hashLB) Pick(id int) Worker {
	idx := id % b.poolSize
	return b.pools[idx]
}

// roundRobinLB
func newRoundRobinLB(pools []Worker) loadBalance {
	return &roundRobinLB{pools: pools, poolSize: len(pools)}
}

type roundRobinLB struct {
	pools    []Worker
	accepted uintptr // accept counter
	poolSize int
}

func (b *roundRobinLB) LoadBalance() LoadBalance {
	return Hash
}

func (b *roundRobinLB) Pick(id int) Worker {
	idx := int(atomic.AddUintptr(&b.accepted, 1)) % b.poolSize
	return b.pools[idx]
}
