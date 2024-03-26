package main

import (
	"log"

	"github.com/music-tribe/golang-pairing-challenge/boltdb"
	"github.com/music-tribe/golang-pairing-challenge/fstore"
	"github.com/music-tribe/golang-pairing-challenge/service"
)

func main() {
	db, err := boltdb.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	store, err := fstore.NewStore("./tmp")
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewService(db, store)
	log.Fatal(svc.Run("8080"))
}
