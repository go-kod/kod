package circuitbreaker

import (
	"context"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCircuitBreaker(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		cb := NewCircuitBreaker(context.Background(), "case1")
		done, err := cb.Allow()
		assert.Nil(t, err)
		done(context.Canceled)
		done, err = cb.Allow()
		assert.Nil(t, err)
		done(context.Canceled)
		done, err = cb.Allow()
		assert.Nil(t, err)
		done(context.Canceled)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)

		time.Sleep(time.Second)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
	})
}

func TestIsSuccessful(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"case 1", args{nil}, true},
		{"case 2", args{context.DeadlineExceeded}, false},
		{"case 3", args{context.Canceled}, false},
		{"case 4", args{status.Error(codes.DeadlineExceeded, "")}, false},
		{"case 5", args{status.Error(codes.ResourceExhausted, "")}, false},
		{"case 6", args{status.Error(codes.Canceled, "")}, false},
		{"case 7", args{status.Error(codes.Aborted, "")}, false},
		{"case 8", args{status.Error(codes.Internal, "")}, false},
		{"case 9", args{status.Error(codes.Unavailable, "")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSuccessful(tt.args.err); got != tt.want {
				t.Errorf("IsSuccessful() = %v, want %v", got, tt.want)
			}
		})
	}
}
