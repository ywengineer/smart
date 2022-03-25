package mr_smart

import (
	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/bytedance/gopkg/util/gopool"
	"sync/atomic"
)

// LoadBalance sets the load balancing method.
type LoadBalance int

const (
	// Random requests that connections are randomly distributed.
	Random LoadBalance = iota
	// Fixed. requests that connections are bind to a fixed pool.
	Fixed
	RoundRobin
)

// loadBalance sets the load balancing method for []*Pool
type loadBalance interface {
	LoadBalance() LoadBalance
	// Choose the most qualified Pool
	Pick(id int) gopool.Pool
}

func newLoadBalance(lb LoadBalance, pools []gopool.Pool) loadBalance {
	switch lb {
	case Random:
		return newRandomLB(pools)
	case Fixed:
		return newFixedLB(pools)
	case RoundRobin:
		return newRoundRobinLB(pools)
	}
	return newFixedLB(pools)
}

// randomLB
func newRandomLB(pools []gopool.Pool) loadBalance {
	return &randomLB{pools: pools, poolSize: len(pools)}
}

type randomLB struct {
	pools    []gopool.Pool
	poolSize int
}

func (b *randomLB) LoadBalance() LoadBalance {
	return Random
}

func (b *randomLB) Pick(id int) gopool.Pool {
	idx := fastrand.Intn(b.poolSize)
	return b.pools[idx]
}

// fixedLB
func newFixedLB(pools []gopool.Pool) loadBalance {
	return &fixedLB{pools: pools, poolSize: len(pools)}
}

type fixedLB struct {
	pools    []gopool.Pool
	poolSize int
}

func (b *fixedLB) LoadBalance() LoadBalance {
	return Fixed
}

func (b *fixedLB) Pick(id int) gopool.Pool {
	idx := id % b.poolSize
	return b.pools[idx]
}

// roundRobinLB
func newRoundRobinLB(pools []gopool.Pool) loadBalance {
	return &roundRobinLB{pools: pools, poolSize: len(pools)}
}

type roundRobinLB struct {
	pools    []gopool.Pool
	accepted uintptr // accept counter
	poolSize int
}

func (b *roundRobinLB) LoadBalance() LoadBalance {
	return Fixed
}

func (b *roundRobinLB) Pick(id int) gopool.Pool {
	idx := int(atomic.AddUintptr(&b.accepted, 1)) % b.poolSize
	return b.pools[idx]
}
