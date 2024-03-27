package ports

import (
	"io"

	"github.com/music-tribe/golang-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
)

//go:generate mockgen -destination=./../internal/mockdeps/mock.go -package=mockdeps -source=dependencies.go
type Store interface {
	SaveFile(filepath string, file io.Reader) error
	GetFile(filepath string) (io.ReadCloser, error)
}

type Database interface {
	Insert(sf *domain.ShowFile) error
	Fetch(id uuid.UUID) (domain.ShowFile, error)
}
