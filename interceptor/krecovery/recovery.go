package krecovery

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/go-kod/kod"
)

type panicError struct {
	Panic any
	Stack []byte
}

func (e *panicError) Error() string {
	return fmt.Sprintf("panic caught: %v\n\n%s", e.Panic, e.Stack)
}

func recoverFrom(p any) error {
	stack := make([]byte, 64<<10)
	stack = stack[:runtime.Stack(stack, false)]

	return &panicError{Panic: p, Stack: stack}
}

// Interceptor returns an interceptor that recovers from panics.
func Interceptor() kod.Interceptor {
	return func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) (err error) {
		normalReturn := false
		defer func() {
			if !normalReturn {
				if r := recover(); r != nil {
					err = recoverFrom(r)
					os.Stderr.Write([]byte(err.Error()))
				}
			}
		}()

		err = invoker(ctx, info, req, reply)
		normalReturn = true

		return err
	}
}
