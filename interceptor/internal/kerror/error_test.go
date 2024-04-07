package kerror

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
			if got := !IsCritical(tt.args.err); got != tt.want {
				t.Errorf("IsSuccessful() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsError(t *testing.T) {
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
		{"case 10", args{status.Error(codes.Unknown, "")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := !IsSystem(tt.args.err); got != tt.want {
				t.Errorf("IsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBusinessError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"case 1", args{nil}, false},
		{"case 2", args{context.DeadlineExceeded}, false},
		{"case 3", args{context.Canceled}, false},
		{"case 4", args{status.Error(codes.DeadlineExceeded, "")}, false},
		{"case 5", args{status.Error(codes.ResourceExhausted, "")}, false},
		{"case 6", args{status.Error(codes.Canceled, "")}, false},
		{"case 7", args{status.Error(codes.Aborted, "")}, false},
		{"case 8", args{status.Error(codes.Internal, "")}, false},
		{"case 9", args{status.Error(codes.Unavailable, "")}, false},
		{"case 10", args{status.Error(codes.Unknown, "")}, false},
		{"case 11", args{status.Error(codes.Unauthenticated+1, "")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBusiness(tt.args.err); got != tt.want {
				t.Errorf("IsBusinessError() = %v, want %v", got, tt.want)
			}
		})
	}
}
