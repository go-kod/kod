package kerror

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsCritical returns true if the error is critical.
func IsCritical(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}

	switch status.Code(err) {
	case codes.Canceled,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.Internal,
		codes.Unavailable:
		return true
	}

	return false
}

// IsSystem returns true if the error is a system error.
func IsSystem(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}

	switch status.Code(err) {
	case codes.Canceled,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.Internal,
		codes.Unavailable,
		codes.Unknown:
		return true
	}

	return false
}

// IsBusiness returns true if the error is a business error.
func IsBusiness(err error) bool {
	return err != nil && !IsSystem(err)
}
