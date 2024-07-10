package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/energostack/bisquitt-psk/pkg/api"
	"github.com/energostack/bisquitt-psk/pkg/clientmap"
	"github.com/energostack/bisquitt-psk/pkg/mapstore"

	"github.com/stretchr/testify/assert"
)

type ReaderMock struct {
	clientMap *clientmap.Map
}

func NewReaderMock() *ReaderMock {
	mock := &ReaderMock{
		clientMap: clientmap.New(),
	}
	mock.clientMap.Store("1", []byte("psk"))
	return mock
}

func (r *ReaderMock) Read() (*clientmap.Map, error) {
	return r.clientMap, nil
}

func (r *ReaderMock) Updates() <-chan bool {
	return make(chan bool)
}

func TestGetClientReturnsSuccess(t *testing.T) {
	reader := NewReaderMock()
	controller := mapstore.NewController(reader)
	handler := api.NewCustomHandler(controller)

	req, _ := http.NewRequest("GET", "/clients/1", nil)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()

	handlerFunc := http.HandlerFunc(handler.GetClient)

	handlerFunc.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetClientReturnsNotFoundWhenClientDoesNotExist(t *testing.T) {
	reader := NewReaderMock()
	controller := mapstore.NewController(reader)
	handler := api.NewCustomHandler(controller)

	req, _ := http.NewRequest("GET", "/clients/nonexistent", nil)
	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.GetClient)

	handlerFunc.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
