package mr_smart

import (
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/ywengineer/mr.smart/log"
	"go.uber.org/zap"
)

type WorkerManager interface {
	Pick(id int) Worker
	Status() interface{}
}

type Worker interface {
	Run(ctx context.Context, f func())
	Status() interface{}
}

// NewWorkerManager create default worker manager for process client channel
func NewWorkerManager(poolSize int, lb LoadBalance) (WorkerManager, error) {
	if poolSize < 1 {
		log.GetLogger().Error("set invalid poolSize", zap.Int("poolSize", poolSize))
		return nil, fmt.Errorf("set invalid poolSize[%d]", poolSize)
	}
	manager := &defaultWorkerManager{}
	for idx := 0; idx < poolSize; idx++ {
		p := gopool.NewPool(fmt.Sprintf("smart-handlers-%d", idx), 1, gopool.NewConfig())
		p.SetPanicHandler(manager.errorHandler)
		manager.workers = append(manager.workers, &defaultWorker{runner: p})
	}
	manager.balance = newLoadBalance(lb, manager.workers)
	return manager, nil
}

type defaultWorkerManager struct {
	balance loadBalance // load balancing method
	workers []Worker    // all the workers
}

// Pick will select the poller for use each time based on the LoadBalance.
func (m *defaultWorkerManager) Pick(id int) Worker {
	return m.balance.Pick(id)
}

func (m *defaultWorkerManager) errorHandler(ctx context.Context, err interface{}) {
	log.GetLogger().Error("process worker task error", zap.Any("error", err))
}

func (m *defaultWorkerManager) Status() interface{} {
	return nil
}

type defaultWorker struct {
	runner gopool.Pool
}

func (w *defaultWorker) Run(ctx context.Context, f func()) {
	w.runner.CtxGo(ctx, f)
}

func (w *defaultWorker) Status() interface{} {
	return nil
}
