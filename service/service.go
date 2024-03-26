package service

import (
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

func (s *Service) Run() error {
	return http.ListenAndServe(":8080", s.router)
}
