package solarpanelapi

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/WilliamDeKeyzer/bdd-workshop-sdk/sdk/solarpanel"
)

type handler struct {
	mutex sync.Mutex
	panel *solarpanel.SolarPanel
}

type reading struct {
	Watts int `json:"watts"`
}

func NewHandler() http.Handler {
	api := &handler{panel: solarpanel.NewSolarPanel()}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /generation", api.getGeneration)
	mux.HandleFunc("PUT /generation", api.setGeneration)
	return mux
}

func (api *handler) getGeneration(writer http.ResponseWriter, request *http.Request) {
	api.mutex.Lock()
	defer api.mutex.Unlock()
	respond(writer, http.StatusOK, reading{Watts: api.panel.GetInjection()})
}

func (api *handler) setGeneration(writer http.ResponseWriter, request *http.Request) {
	var input reading
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		respond(writer, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	api.mutex.Lock()
	defer api.mutex.Unlock()
	if err := api.panel.SetNewInjection(input.Watts); err != nil {
		respond(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	respond(writer, http.StatusOK, reading{Watts: api.panel.GetInjection()})
}

func respond(writer http.ResponseWriter, status int, body any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(body)
}
