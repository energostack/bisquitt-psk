package clientmap

import "sync"

// Map is a thread-safe map of clients to pre-shared keys.
type Map struct {
	mutex        sync.RWMutex
	clientPskMap map[string][]byte
}

// New creates a new Map.
func New() *Map {
	return &Map{
		clientPskMap: make(map[string][]byte),
	}
}

// Load retrieves the value for a key.
func (m *Map) Load(key string) (value []byte, ok bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok = m.clientPskMap[key]
	return
}

// Store sets the value for a key.
func (m *Map) Store(key string, value []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.clientPskMap[key] = value
}

// Delete removes the value for a key.
func (m *Map) Delete(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.clientPskMap, key)
}

// LoadOrStore retrieves the existing value for a key or stores a new value.
func (m *Map) LoadOrStore(key string, value []byte) (actual []byte, loaded bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if actual, loaded = m.clientPskMap[key]; loaded {
		return
	}
	m.clientPskMap[key] = value
	return value, false
}

// Range calls f sequentially for each key and value present in the map.
func (m *Map) Range(f func(key string, value []byte) bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for k, v := range m.clientPskMap {
		if !f(k, v) {
			break
		}
	}
}

// Get returns the copy of map.
func (m *Map) Get() map[string][]byte {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	newMap := make(map[string][]byte)
	for k, v := range m.clientPskMap {
		newMap[k] = v
	}
	return newMap
}

// Set insert items to the map.
func (m *Map) Set(newMap map[string][]byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k, _ := range m.clientPskMap {
		delete(m.clientPskMap, k)
	}
	for k, v := range newMap {
		m.clientPskMap[k] = v
	}
}
