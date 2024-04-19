package mapstore_test

import (
	"testing"

	"bisquitt-psk/pkg/clientmap"
	"bisquitt-psk/pkg/mapstore"

	"github.com/stretchr/testify/assert"
)

type ReaderMock struct {
	clientMap *clientmap.Map
	updateCh  chan bool
}

func NewReaderMock() *ReaderMock {
	mock := &ReaderMock{
		clientMap: clientmap.New(),
		updateCh:  make(chan bool),
	}
	mock.clientMap.Store("1", []byte("psk"))
	return mock
}

func (r *ReaderMock) Read() (*clientmap.Map, error) {
	return r.clientMap, nil
}

func (r *ReaderMock) Updates() <-chan bool {
	return r.updateCh
}

func (r *ReaderMock) TriggerUpdate() {
	r.clientMap.Store("2", []byte("psk"))
	r.updateCh <- true
}

func TestControllerSyncsClients(t *testing.T) {
	reader := NewReaderMock()
	controller := mapstore.NewController(reader)
	_, ok := controller.GetPSK("2")
	assert.False(t, ok, "Expected client 2 to not be in the map")
	reader.TriggerUpdate()

	_, ok = controller.GetPSK("2")
	assert.True(t, ok, "Expected client 2 to be in the map")
}
