package singleton

import (
	"sync"
)

// Singleton[T any] provides a common structure for components,
type Singleton[T any] struct {
	instances map[string]*T
	mu        sync.RWMutex
}

// NewSingleton creates a new Singleton[T].
func NewSingleton[T any]() *Singleton[T] {
	return &Singleton[T]{
		instances: make(map[string]*T),
	}
}

// Get returns the instance of the component with the given name.
func (s *Singleton[T]) Get(name string, initFn func() *T) *T {
	s.mu.RLock()

	if instance, exists := s.instances[name]; exists {
		s.mu.RUnlock()
		return instance
	}

	s.mu.RUnlock()

	instance := initFn()

	s.mu.Lock()
	s.instances[name] = instance
	s.mu.Unlock()

	return instance
}
