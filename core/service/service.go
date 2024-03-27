package service

import (
	"io"

	"github.com/music-tribe/golang-pairing-challenge/core/ports"
	"github.com/music-tribe/golang-pairing-challenge/core/service/getfile"
)

type Service struct {
	database ports.Database
	store    ports.Store
}

func NewService(database ports.Database, store ports.Store) *Service {
	svc := &Service{
		database: database,
		store:    store,
	}

	return svc
}

func (s *Service) GetFile(req ports.GetFileRequest) (io.ReadCloser, error) {
	return getfile.Action(s.database, s.store, req)
}
