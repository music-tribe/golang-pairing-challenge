package service

import (
	"github.com/music-tribe/golang-pairing-challenge/ports"
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
