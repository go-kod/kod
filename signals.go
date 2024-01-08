package kod

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

var (
	ErrorShutdown         = errors.New("shutdown")
	ErrorGracefulShutdown = errors.New("graceful shutdown")
)

var shutdownSignals = []os.Signal{syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM}

// shutdown support twice signal must exit
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
