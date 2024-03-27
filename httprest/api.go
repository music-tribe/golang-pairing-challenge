package httprest

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/music-tribe/golang-pairing-challenge/core/ports"
	"github.com/music-tribe/golang-pairing-challenge/httprest/getfile"
)

type API struct {
	svc    ports.Service
	router *mux.Router
}

func NewAPI(svc ports.Service) *API {
	api := &API{
		svc:    svc,
		router: mux.NewRouter(),
	}

	api.RegisterRoutes()

	return api
}

func (a *API) RegisterRoutes() {
	// new routes
	a.router.HandleFunc("/api/files/{user_id}/{id}", getfile.Handler(a.svc)).Methods("GET")

	// lagacy routes
	a.router.HandleFunc("/api/files/{user_id}", a.svc.UploadFile).Methods("POST")
}

func (a *API) Run(port string) error {
	a.RegisterRoutes()
	errChan := make(chan error)

	go func() {
		defer close(errChan)
		errChan <- http.ListenAndServe(":"+port, a.router)
	}()

	log.Printf("API listening on port %s\n", port)

	return <-errChan
}
