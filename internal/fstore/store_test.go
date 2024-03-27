package fstore

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	t.Run("test store lifecycle", func(t *testing.T) {
		dir := os.TempDir()
		store, err := NewStore(dir)
		if err != nil {
			t.Fatal(err)
		}

		filepath := "some/path/to/file.txt"

		fileData := []byte("some file data")
		err = store.SaveFile(filepath, bytes.NewBuffer(fileData))
		require.NoError(t, err)

		file, err := store.GetFile(filepath)
		require.NoError(t, err)

		byt, err := io.ReadAll(file)
		require.NoError(t, err)

		err = file.Close()
		require.NoError(t, err)

		assert.Equal(t, fileData, byt)
	})
}
