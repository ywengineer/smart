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
	// Hash. requests that connections are bind to a fixed pool.
	Hash
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
	case Hash:
		return newHashLB(pools)
	case RoundRobin:
		return newRoundRobinLB(pools)
	}
	return newRoundRobinLB(pools)
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

// hashLB
func newHashLB(pools []gopool.Pool) loadBalance {
	return &hashLB{pools: pools, poolSize: len(pools)}
}

type hashLB struct {
	pools    []gopool.Pool
	poolSize int
}

func (b *hashLB) LoadBalance() LoadBalance {
	return Hash
}

func (b *hashLB) Pick(id int) gopool.Pool {
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
	return Hash
}

func (b *roundRobinLB) Pick(id int) gopool.Pool {
	idx := int(atomic.AddUintptr(&b.accepted, 1)) % b.poolSize
	return b.pools[idx]
}
