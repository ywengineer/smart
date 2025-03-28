package utility

import (
	"context"
	"github.com/ywengineer/smart-kit/pkg/logk"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func WatchQuitSignal() chan os.Signal {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGQUIT)
	//
	return quit
}

func Watch(ctx context.Context, notify chan<- bool) {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
	stop:
		for {
			select {
			case <-ctx.Done():
				logk.Info("terminating: context cancelled")
				notify <- true
				break stop
			case <-ticker.C:
				if ctx.Err() != nil {
					logk.Info("terminating: context cancelled")
					notify <- true
					break stop
				}
			}
		}
		ticker.Stop()
	}()
}
