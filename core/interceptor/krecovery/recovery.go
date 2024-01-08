package krecovery

import (
	"context"
	"os"

	"github.com/go-kod/kod"
)

// Interceptor returns an interceptor that recovers from panics.
func Interceptor() kod.Interceptor {
	return func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) (err error) {
		normalReturn := false
		defer func() {
			if !normalReturn {
				if r := recover(); r != nil {
					err = kod.RecoverFrom(r)
					os.Stderr.Write([]byte(err.Error()))
				}
			}
		}()

		err = invoker(ctx, info, req, reply)
		normalReturn = true

		return err
	}
}
