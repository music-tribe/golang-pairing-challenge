package fstore

import (
	"bytes"
	"io"
	"os"
	"path"
)

type Store struct {
	dir string
}

func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &Store{
		dir: dir,
	}, nil
}

func (s *Store) SaveFile(filepath string, file io.Reader) error {
	buf := bytes.NewBuffer([]byte{})

	if _, err := io.Copy(buf, file); err != nil {
		return err
	}

	return os.WriteFile(path.Join(s.dir, filepath), buf.Bytes(), 0644)
}
