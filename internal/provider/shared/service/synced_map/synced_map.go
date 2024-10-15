// Package synced_map implements a thread safe map with generics support.
package synced_map

import (
	"sync"
)

// SyncedMap is a thread safe map with generics support.
type SyncedMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// Get retrieves a value from the map. The second returned value denotes whether
// the key exists or not.
func (m *SyncedMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.m[key]
	return v, ok
}

func (m *SyncedMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	m.m[key] = value
	m.mu.Unlock()
}

func NewSyncedMap[K comparable, V any]() SyncedMap[K, V] {
	m := make(map[K]V)

	return SyncedMap[K, V]{m: m}
}
