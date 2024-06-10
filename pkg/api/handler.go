package api

import (
	"encoding/json"
	"net/http"

	"bisquitt-psk/pkg/mapstore"
)

// Data is an object of client IDs (string) and  PSK (byte slice).
type Data struct {
	Client string `json:"client"`
	Psk    []byte `json:"psk"`
}

type ResponseError struct {
	Message string `json:"message"`
}

// CustomHandler is a handler for the API.
type CustomHandler struct {
	mapStoreController *mapstore.Controller
}

// NewCustomHandler creates a new handler.
func NewCustomHandler(controller *mapstore.Controller) *CustomHandler {
	return &CustomHandler{
		mapStoreController: controller,
	}
}

// GetClient is a handler for getting a client by ID.
//
//	@Summary		Get client
//	@Description	Get a client by ID
//	@Tags			clients
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Data
//	@Failure		500	{object}	ResponseError
//	@Param			id	path		string	true	"Client ID"
//	@Router			/clients/{id} [get]
func (ch *CustomHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	psk, ok := ch.mapStoreController.GetPSK(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		data, _ := json.Marshal(ResponseError{Message: "Client not found"})
		w.Write(data)
		return
	}
	data, err := json.Marshal(Data{Client: id, Psk: psk})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data, _ = json.Marshal(ResponseError{Message: "Failed to marshal data"})
		w.Write(data)
		return
	}
	w.Write(data)
}
