package getfile

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/music-tribe/golang-pairing-challenge/service/mocksvc"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("when the bind fails we should receive a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		router := mux.NewRouter()
		router.HandleFunc("/getfile/{id}/{user_id}", Handler(mocksvc.NewMockService(ctrl))).Methods(http.MethodGet)

		req := httptest.NewRequest(http.MethodGet, "/getfile/invalid/invalid", nil)
		res := httptest.NewRecorder()

		svc := mocksvc.NewMockService(ctrl)
		Handler(svc)(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "invalid UUID length")
	})

	t.Run("when the business logic fails we should receive a 500 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		svc := mocksvc.NewMockService(ctrl)
		svc.EXPECT().GetFile(gomock.Any()).Return(nil, errors.New("some svc error"))
		handler := Handler(svc)

		router := mux.NewRouter()
		router.HandleFunc("/getfile/{id}/{user_id}", handler).Methods(http.MethodGet)

		req := httptest.NewRequest(http.MethodGet, "/getfile/f965681c-0c52-47ce-802e-6aa112c08065/8dd7e4ac-00c8-40e2-9a3e-40fdebc3b8fa", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Contains(t, res.Body.String(), "some svc error")
	})

	t.Run("when the request is good we should get a 200 and the file", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		fileData := []byte("Hello, World!")

		someFile := &closer{bytes.NewBuffer(fileData)}

		svc := mocksvc.NewMockService(ctrl)
		svc.EXPECT().GetFile(gomock.Any()).Return(someFile, nil)
		handler := Handler(svc)

		router := mux.NewRouter()
		router.HandleFunc("/getfile/{id}/{user_id}", handler).Methods(http.MethodGet)

		req := httptest.NewRequest(http.MethodGet, "/getfile/f965681c-0c52-47ce-802e-6aa112c08065/8dd7e4ac-00c8-40e2-9a3e-40fdebc3b8fa", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, fileData, res.Body.Bytes())
	})
}

type closer struct {
	*bytes.Buffer
}

func (c *closer) Close() error {
	return nil
}

func Test_bindRequest(t *testing.T) {
	id := uuid.New()
	userId := uuid.New()

	type args struct {
		vars map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return error when id is invalid",
			args: args{
				vars: map[string]string{
					"id":      "invalid",
					"user_id": userId.String(),
				},
			},
			wantErr: true,
		},
		{
			name: "should return error when userId is invalid",
			args: args{
				vars: map[string]string{
					"id":      id.String(),
					"user_id": "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "should return a completed GetFileRequest",
			args: args{
				vars: map[string]string{
					"id":      id.String(),
					"user_id": userId.String(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := Request{}
			err := bindRequest(tt.args.vars, &params)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			expect := Request{
				Id:     id,
				UserId: userId,
			}

			assert.NoError(t, err)
			assert.Equal(t, expect, params)
		})
	}
}

func Test_sendResponse(t *testing.T) {
	fileData := []byte("Hello, World!")
	preRead := bytes.NewBuffer(fileData)
	r := make([]byte, len(fileData))
	preRead.Read(r)

	type args struct {
		fileData []byte
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "should return 200 and the correct file data",
			args: args{
				fileData: []byte("Hello, World!"),
			},
			wantErr:    false,
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			err := sendResponse(res, bytes.NewReader(tt.args.fileData))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.Code)
			assert.Equal(t, tt.args.fileData, res.Body.Bytes())
		})
	}
}
