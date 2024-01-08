package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestBBR(t *testing.T) {
	t.Run("case 1ms", func(t *testing.T) {
		bbr := NewLimiter(context.Background(), "bbr")
		bbr.opts.CPUThreshold = 0

		done, err := bbr.Allow()
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond)
		done()

		_, _ = bbr.Allow()
		_, _ = bbr.Allow()

		_, err = bbr.Allow()
		if err != ErrLimitExceed {
			t.Fatal(err)
		}
	})

	t.Run("case 1s", func(t *testing.T) {
		bbr := NewLimiter(context.Background(), "bbr")

		done, err := bbr.Allow()
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second)
		done()

		_, _ = bbr.Allow()
		_, _ = bbr.Allow()

		_, _ = bbr.Allow()
		if err != nil {
			t.Fatal(err)
		}

		stat := bbr.Stat()
		if stat.InFlight != 3 {
			t.Fatal("inflight not expected")
		}
	})
}
