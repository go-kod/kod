package kod

import (
	"context"
)

type deferFunc struct {
	Name string
	Fn   func(context.Context) error
}

// addDefer adds a defer function to the Kod instance.
func (k *Kod) addDefer(d deferFunc) {
	k.deferMux.Lock()
	defer k.deferMux.Unlock()
	k.defers = append(k.defers, d)
}

// runDefer runs the defer functions in reverse order.
func (k *Kod) runDefer(ctx context.Context) {
	ctx, timeoutCancel := context.WithTimeout(context.WithoutCancel(ctx), k.config.ShutdownTimeout)
	defer timeoutCancel()

	k.deferMux.Lock()
	defer k.deferMux.Unlock()
	for i := len(k.defers) - 1; i >= 0; i-- {
		err := k.defers[i].Fn(ctx)
		if err != nil {
			k.log.Error("component %q stop failed: %w", k.defers[i].Name, err)
		}
	}
}
