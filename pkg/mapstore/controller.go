package mapstore

import (
	"github.com/energostack/bisquitt-psk/pkg/clientmap"
	"github.com/energostack/bisquitt-psk/pkg/reader"

	"github.com/rs/zerolog/log"
)

// Controller is a controller that manages the client map.
type Controller struct {
	reader    reader.Reader
	updateCh  <-chan bool
	clientMap *clientmap.Map
}

// NewController creates a new controller with the specified reader.
func NewController(reader reader.Reader) *Controller {
	controller := Controller{
		reader:   reader,
		updateCh: reader.Updates(),
	}

	clientMap, err := reader.Read()
	if err != nil {
		log.Err(err).Msg("Failed to read clients from reader")
		controller.clientMap = clientmap.New()
	} else {
		controller.clientMap = clientMap
	}

	controller.syncClients()

	return &controller
}

// GetPSK returns the client map.
func (c *Controller) GetPSK(id string) ([]byte, bool) {
	return c.clientMap.Load(id)
}

func (c *Controller) syncClients() {
	go func() {
		for {
			select {
			case <-c.updateCh:
				newClientMap, err := c.reader.Read()
				if err != nil {
					log.Err(err).Msg("Failed to read clients from reader")
					continue
				}
				c.clientMap.Set(newClientMap.Get())
			}
		}
	}()
}
