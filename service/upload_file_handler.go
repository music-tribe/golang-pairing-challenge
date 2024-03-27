package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gorilla/mux"
	"github.com/music-tribe/golang-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
)

type Request struct {
	UserId uuid.UUID
	Id     uuid.UUID
}

type Response struct {
	Id       uuid.UUID `json:"id"`
	Filepath string    `json:"filepath"`
}

func (s *Service) UploadFile(w http.ResponseWriter, r *http.Request) {
	// check the request content type for multipart form
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "request content type is not multipart/form-data", http.StatusBadRequest)
		return
	}

	// Parse the multipart form in the request
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// get the user_id from the path params and create new id
	vars := mux.Vars(r)
	userId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := uuid.New()

	// snif the file content type
	headSize := 1024
	head := make([]byte, headSize)
	bytesRead, err := file.Read(head)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bytes.Equal(head, make([]byte, headSize)) { // check if []byte content is empty
		http.Error(w, "cannot decipher file type", http.StatusInternalServerError)
		return
	}

	contentType := mimetype.Detect(head[:bytesRead])

	// TODO: we should remove this as part of the test
	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create the filepath
	filepath := fmt.Sprintf("%s/%s/%s", userId, id, fileHeader.Filename)

	// Save the file
	if err = s.store.SaveFile(filepath, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the file metadata
	if err = s.database.Insert(&domain.ShowFile{
		Id:           id,
		UserId:       userId,
		ContenctType: contentType.String(),
		Filename:     fileHeader.Filename,
		Size:         fileHeader.Size,
		Filepath:     filepath,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send the response
	resp := Response{
		Id:       id,
		Filepath: filepath,
	}

	byt, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(byt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
