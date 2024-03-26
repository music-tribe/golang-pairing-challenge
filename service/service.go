package service

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Service struct {
	router   *mux.Router
	database Database
	store    Store
}

func NewService(database Database, store Store) *Service {
	svc := &Service{
		router:   mux.NewRouter(),
		database: database,
		store:    store,
	}

	svc.router.HandleFunc("/api/file/{user_id}", svc.UploadFile).Methods("POST")

	return svc
}

func (s *Service) Run(port string) error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		errChan <- http.ListenAndServe(":"+port, s.router)
	}()

	log.Println("Server is running on port:", port)

	return <-errChan
}
