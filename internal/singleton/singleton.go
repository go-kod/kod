package singleton

import (
	"sync"
)

// Singleton[T any] provides a common structure for components,
type Singleton[T any] struct {
	instances map[string]*T
	mu        sync.Mutex
}

// NewSingleton creates a new Singleton[T].
func NewSingleton[T any]() *Singleton[T] {
	return &Singleton[T]{
		instances: make(map[string]*T),
	}
}

// Get returns the instance of the component with the given name.
func (s *Singleton[T]) Get(name string, initFn func() *T) *T {
	s.mu.Lock()
	defer s.mu.Unlock()

	if instance, exists := s.instances[name]; exists {
		return instance
	}

	instance := initFn()
	s.instances[name] = instance

	return instance
}
