package kod

import (
	"context"
	"fmt"
)

type deferFn struct {
	Name string
	Fn   func(context.Context) error
}

func (k *Kod) deferFn(d deferFn) {
	k.deferMux.Lock()
	defer k.deferMux.Unlock()
	k.defers = append(k.defers, d)
}

func (k *Kod) runDefer(ctx context.Context) error {
	k.deferMux.Lock()
	defer k.deferMux.Unlock()
	for i := len(k.defers) - 1; i >= 0; i-- {
		err := k.defers[i].Fn(ctx)
		if err != nil {
			return fmt.Errorf("component %q stop failed: %w", k.defers[i].Name, err)
		}
	}

	return nil
}
