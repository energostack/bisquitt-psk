package clientmap_test

import (
	"testing"

	"bisquitt-psk/pkg/clientmap"

	"github.com/stretchr/testify/assert"
)

func TestKeyExistsInMapReturnsTrue(t *testing.T) {
	m := clientmap.New()
	m.Store("key", []byte("value"))
	_, ok := m.Load("key")
	assert.True(t, ok)
}

func TestKeyDoesNotExistInMapReturnsFalse(t *testing.T) {
	m := clientmap.New()
	_, ok := m.Load("nonexistent")
	assert.False(t, ok)
}

func TestStoreAddsKeyToMap(t *testing.T) {
	m := clientmap.New()
	m.Store("key", []byte("value"))
	value, _ := m.Load("key")
	assert.Equal(t, []byte("value"), value)
}

func TestDeleteRemovesKeyFromMap(t *testing.T) {
	m := clientmap.New()
	m.Store("key", []byte("value"))
	m.Delete("key")
	_, ok := m.Load("key")
	assert.False(t, ok)
}

func TestLoadOrStoreReturnsExistingValue(t *testing.T) {
	m := clientmap.New()
	m.Store("key", []byte("value"))
	value, _ := m.LoadOrStore("key", []byte("newvalue"))
	assert.Equal(t, []byte("value"), value)
}

func TestLoadOrStoreStoresNewValue(t *testing.T) {
	m := clientmap.New()
	value, _ := m.LoadOrStore("key", []byte("value"))
	assert.Equal(t, []byte("value"), value)
}

func TestRangeIteratesOverAllKeys(t *testing.T) {
	m := clientmap.New()
	m.Store("key1", []byte("value1"))
	m.Store("key2", []byte("value2"))
	keys := make([]string, 0)
	m.Range(func(key string, value []byte) bool {
		keys = append(keys, key)
		return true
	})
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
}

func TestGetReturnsAllKeys(t *testing.T) {
	m := clientmap.New()
	m.Store("key1", []byte("value1"))
	m.Store("key2", []byte("value2"))
	mapCopy := m.Get()
	assert.Equal(t, []byte("value1"), mapCopy["key1"])
	assert.Equal(t, []byte("value2"), mapCopy["key2"])
}

func TestSetReplacesAllKeys(t *testing.T) {
	m := clientmap.New()
	m.Store("key1", []byte("value1"))
	m.Set(map[string][]byte{
		"key2": []byte("value2"),
	})
	_, ok := m.Load("key1")
	assert.False(t, ok)
	value, _ := m.Load("key2")
	assert.Equal(t, []byte("value2"), value)
}

func TestConcurrentAccessDoesNotCauseRaceCondition(t *testing.T) {
	m := clientmap.New()
	go func() {
		for i := 0; i < 1000; i++ {
			m.Store("key", []byte("value"))
		}
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			m.Load("key")
		}
	}()
}
