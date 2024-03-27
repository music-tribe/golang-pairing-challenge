package boltdb

import (
	"os"
	"testing"

	"github.com/music-tribe/golang-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	t.Run("test database lifecycle", func(t *testing.T) {
		db, err := NewDatabase()
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() {
			err = os.RemoveAll(db.session.Path())
			require.NoError(t, err)

			err := db.session.Close()
			assert.NoError(t, err)
		})

		expect := &domain.ShowFile{
			Id:     uuid.New(),
			UserId: uuid.New(),
		}

		err = db.Insert(expect)
		require.NoError(t, err)

		actual, err := db.Fetch(expect.Id)

		require.NoError(t, err)

		assert.Equal(t, expect, actual)

	})
}
