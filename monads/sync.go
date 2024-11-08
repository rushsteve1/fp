package monads

import "sync"

// This file implements some minor synchronization types

type Mutex[T any] struct {
	Cell[T]
	mu *sync.Mutex
}

func NewMutex[T any](v T) Mutex[T] {
	return Mutex[T]{
		Cell: NewCell(v),
		mu:   &sync.Mutex{},
	}
}

func (m Mutex[T]) Get() (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Cell.Get()
}

func (m *Mutex[T]) Set(v T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Cell.Set(v)
}

type RWLock[T any] struct {
	Cell[T]
	mu *sync.Mutex
}

func NewRWLock[T any](v T) RWLock[T] {
	return RWLock[T]{
		Cell: NewCell(v),
		mu:   &sync.Mutex{},
	}
}

func (m RWLock[T]) Get() (T, error) {
	return m.Cell.Get()
}

func (m *RWLock[T]) Set(v T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Cell.Set(v)
}
