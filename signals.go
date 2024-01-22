package kod

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM}

// shutdown support twice signal must exit
// first signal: graceful shutdown
// second signal: exit directly
func shutdown(ctx context.Context, stop func(grace bool)) {
	sig := make(chan os.Signal, 2)
	signal.Notify(
		sig,
		shutdownSignals...,
	)
	go func() {
		select {
		case <-ctx.Done():
			return
		case s := <-sig:
			go stop(s != syscall.SIGQUIT)
			<-sig
			os.Exit(128 + int(s.(syscall.Signal))) // second signal. Exit directly.
		}
	}()
}
