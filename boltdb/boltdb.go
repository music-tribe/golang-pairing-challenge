package boltdb

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"github.com/music-tribe/golang-pairing-challenge/domain"
)

type Database struct {
	session *bolt.DB
}

func NewDatabase() (*Database, error) {
	sess, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Database{
		session: sess,
	}, nil
}

func (d *Database) Insert(sf *domain.ShowFile) error {
	byt, err := json.Marshal(sf)
	if err != nil {
		return err
	}

	return d.session.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("showFiles"))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(sf.Id.String()), byt)
	})
}
