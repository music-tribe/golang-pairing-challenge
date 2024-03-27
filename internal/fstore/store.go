package fstore

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

type Store struct {
	dir string
}

func NewStore(dir string) (*Store, error) {
	return &Store{
		dir: dir,
	}, nil
}

func (s *Store) SaveFile(key string, value io.Reader) error {
	fullpath := path.Join(s.dir, key)

	if err := os.MkdirAll(filepath.Dir(fullpath), 0770); err != nil && err != io.EOF {
		return err
	}

	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, value)

	return err
}
