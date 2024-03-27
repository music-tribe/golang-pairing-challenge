package service

import (
	"io"

	"github.com/music-tribe/golang-pairing-challenge/ports"
)

func (s *Service) GetFile(req ports.GetFileRequest) (io.ReadCloser, error) {
	record, err := s.database.Fetch(req.Id)
	if err != nil {
		return nil, err
	}

	file, err := s.store.GetFile(record.Filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
