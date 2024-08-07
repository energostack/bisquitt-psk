package reader

import "github.com/energostack/bisquitt-psk/pkg/clientmap"

// Reader is an interface for reading clients.
type Reader interface {
	Read() (*clientmap.Map, error)
	Updates() <-chan bool
}
