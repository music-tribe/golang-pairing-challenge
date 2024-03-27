package boltdb

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"github.com/music-tribe/golang-pairing-challenge/domain"
	"github.com/music-tribe/golang-pairing-challenge/ports"
	"github.com/music-tribe/uuid"
)

const bucketName = "showFiles"

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
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(sf.Id.String()), byt)
	})
}

func (d *Database) Fetch(id uuid.UUID) (*domain.ShowFile, error) {
	var sf domain.ShowFile
	err := d.session.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ports.ErrNotFound
		}

		byt := bucket.Get([]byte(id.String()))
		if byt == nil {
			return ports.ErrNotFound
		}

		return json.Unmarshal(byt, &sf)
	})

	return &sf, err
}
