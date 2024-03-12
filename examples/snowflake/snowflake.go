package snowflake

import (
	"context"
	"time"

	"github.com/go-kod/kod"
	"github.com/sony/sonyflake"
)

type impl struct {
	kod.Implements[Component]
	snowflake *sonyflake.Sonyflake
}

func (ins *impl) Init(ctx context.Context) error {
	ins.snowflake = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: func() (uint16, error) {
			// TODO use redis to get machine id
			return 1, nil
		},
	})

	return nil
}

type NextIDRequest struct{}

type NextIDResponse struct {
	ID uint64
}

// NextID returns a unique ID
func (ins *impl) NextID(ctx context.Context, req *NextIDRequest) (*NextIDResponse, error) {
	id, err := ins.snowflake.NextID()
	if err != nil {
		return nil, err
	}

	return &NextIDResponse{ID: id}, nil
}
