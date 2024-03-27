package ports

import (
	"io"
	"net/http"

	"github.com/music-tribe/uuid"
)

//go:generate mockgen -destination=./../service/mocksvc/mock.go -package=mocksvc -source=service.go
type Service interface {
	// new actions
	GetFile(GetFileRequest) (io.ReadCloser, error)

	// leagacy handlers
	UploadFile(w http.ResponseWriter, r *http.Request) // TODO: refactor me!!!
}

type GetFileRequest struct {
	Id     uuid.UUID `validate:"required"`
	UserId uuid.UUID `validate:"required"`
}

type UploadFileRequest struct {
	// TODO: refactor me!!!
}
