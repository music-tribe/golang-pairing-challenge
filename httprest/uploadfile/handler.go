package uploadfile

import (
	"net/http"

	"github.com/music-tribe/golang-pairing-challenge/core/ports"
)

func Handler(svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
