package service

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/music-tribe/golang-pairing-challenge/internal/boltdb"
	"github.com/music-tribe/golang-pairing-challenge/internal/fstore"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_UploadFile(t *testing.T) {
	t.Run("Approval test for UploadFile - covers only the happy path", func(t *testing.T) {
		// prepare dependencies
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, err := boltdb.NewDatabase()
		require.NoError(t, err)

		dir := t.TempDir()
		store, err := fstore.NewStore(dir)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := os.RemoveAll("./my.db")
			require.NoError(t, err)
		})

		// init service
		svc := NewService(db, store)

		// init router
		router := mux.NewRouter()
		router.HandleFunc("/api/files/{user_id}", svc.UploadFile).Methods("POST")

		// init mulitpart form data
		body := bytes.NewBuffer(nil)
		mpw := multipart.NewWriter(body)
		testData := []byte("test data")
		fw, err := mpw.CreateFormFile("file", "test.txt")
		require.NoError(t, err)
		_, err = fw.Write(testData)
		require.NoError(t, err)
		err = mpw.Close()
		require.NoError(t, err)

		// init request
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/files/2b652684-f786-404b-a4b1-624e19230fb1", body)
		req.Header.Set("Content-Type", mpw.FormDataContentType())
		router.ServeHTTP(res, req)

		// prep for assert
		resp := Response{}
		err = json.Unmarshal(res.Body.Bytes(), &resp)
		require.NoError(t, err)

		// assert
		assert.Equal(t, http.StatusOK, res.Code)
		assert.NotEqual(t, uuid.Nil, resp.Id)
		assert.NotEmpty(t, resp.Filepath)
		assert.FileExists(t, filepath.Join(dir, resp.Filepath))

		// check file content
		byt, err := os.ReadFile(filepath.Join(dir, resp.Filepath))
		require.NoError(t, err)
		assert.Equal(t, testData, byt)
	})
}
