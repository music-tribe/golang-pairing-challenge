package main

import (
	"log"

	"github.com/music-tribe/golang-pairing-challenge/httprest"
	"github.com/music-tribe/golang-pairing-challenge/internal/boltdb"
	"github.com/music-tribe/golang-pairing-challenge/internal/fstore"
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
	api := httprest.NewAPI(svc)

	log.Fatal(api.Run("8080"))
}
