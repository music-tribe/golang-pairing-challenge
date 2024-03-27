package getfile

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/music-tribe/golang-pairing-challenge/core/ports"
	"github.com/music-tribe/uuid"
)

type Request struct {
	Id     uuid.UUID `validate:"required"`
	UserId uuid.UUID `validate:"required"`
}

func Handler(svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		vars := mux.Vars(r)
		if err := bindRequest(vars, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, err := svc.GetFile(ports.GetFileRequest{
			Id:     req.Id,
			UserId: req.UserId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		if err := sendResponse(w, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func bindRequest(vars map[string]string, req *Request) (err error) {
	req.Id, err = uuid.Parse(vars["id"])
	if err != nil {
		return err
	}

	req.UserId, err = uuid.Parse(vars["user_id"])
	if err != nil {
		return err
	}

	return nil
}

func sendResponse(w http.ResponseWriter, file io.Reader) error {
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, file); err != nil {
		return err
	}

	return nil
}
