package getfile

import (
	"io"

	"github.com/music-tribe/golang-pairing-challenge/ports"
)

func Action(db ports.Database, st ports.Store, req ports.GetFileRequest) (io.ReadCloser, error) {
	record, err := db.Fetch(req.Id)
	if err != nil {
		return nil, err
	}

	file, err := st.GetFile(record.Filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
