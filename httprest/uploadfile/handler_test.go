package uploadfile

import (
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("TestHandler", func(t *testing.T) {
		Handler(nil)
	})
}
