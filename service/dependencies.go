package service

import (
	"io"

	"github.com/music-tribe/golang-pairing-challenge/domain"
)

type Store interface {
	SaveFile(filepath string, file io.Reader) error
}

type Database interface {
	Insert(sf *domain.ShowFile) error
}
