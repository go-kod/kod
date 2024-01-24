package hooks

import (
	"context"
	"sync"
)

type HookFunc struct {
	Name string
	Fn   func(context.Context) error
}

type Hooker struct {
	mux   sync.Mutex
	hooks []HookFunc
}

func New() *Hooker {
	return &Hooker{
		mux:   sync.Mutex{},
		hooks: make([]HookFunc, 0),
	}
}

// Add adds a hook function to the Kod instance.
func (k *Hooker) Add(d HookFunc) {
	k.mux.Lock()
	defer k.mux.Unlock()
	k.hooks = append(k.hooks, d)
}

// Do runs the hook functions in reverse order.
func (k *Hooker) Do(ctx context.Context) {

	k.mux.Lock()
	defer k.mux.Unlock()
	for i := len(k.hooks) - 1; i >= 0; i-- {
		_ = k.hooks[i].Fn(ctx)
	}
}
