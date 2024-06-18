package signals

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestShutdown(t *testing.T) {
	sig := make(chan os.Signal, 1)

	time.AfterFunc(time.Millisecond, func() {
		sig <- os.Interrupt
	})

	Shutdown(context.Background(), sig, func(grace bool) {
	})

	time.Sleep(time.Second)
}
