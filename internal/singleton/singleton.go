package singleton

import (
	"sync"
)

// Singleton 是一个泛型结构体，用于存储多个单一实例
type Singleton[T any] struct {
	instances map[string]*T
	mu        sync.Mutex
}

// NewSingleton 创建一个新的Singleton实例。
func NewSingleton[T any]() *Singleton[T] {
	return &Singleton[T]{
		instances: make(map[string]*T),
	}
}

// Get 返回与给定名称对应的单例实例。如果该实例还不存在，则使用提供的初始化函数initFn来创建它。
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
