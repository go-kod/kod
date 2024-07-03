package snowflake

import (
	"context"
	"time"

	"github.com/sony/sonyflake"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/examples/infra/redis"
)

type impl struct {
	kod.Implements[Component]
	redis kod.Ref[redis.Component]

	snowflake *sonyflake.Sonyflake
}

func (ins *impl) Init(ctx context.Context) error {
	ins.snowflake = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: func() (uint16, error) {
			machineId, err := ins.redis.Get().Client().Incr(ctx, "snowflake:id").Uint64()
			return uint16(machineId), err
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
