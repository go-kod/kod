package knet

import (
	"net"
	"reflect"
	"testing"

	"github.com/samber/lo"
)

func TestGetRegisterAddr(t *testing.T) {
	type args struct {
		hostPort string
		lis      net.Listener
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{

		{
			"",
			args{
				hostPort: "localhost:0",
				lis:      lo.Must(net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 8124})),
			},
			"localhost:8124",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractAddress(tt.args.hostPort, tt.args.lis)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRegisterAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRegisterAddr() = %v, want %v", got, tt.want)
			}
		})
	}

	lis, err := net.Listen("tcp", ":12345")
	if err != nil {
		t.Errorf("expected: %v got %v", nil, err)
	}
	res, err := ExtractAddress("", lis)
	if err != nil {
		t.Errorf("expected: %v got %v", nil, err)
	}
	expect, err := ExtractAddress(lis.Addr().String(), nil)
	if err != nil {
		t.Errorf("expected: %v got %v", nil, err)
	}
	if !reflect.DeepEqual(res, expect) {
		t.Errorf("expected %s got %s", expect, res)
	}
}
