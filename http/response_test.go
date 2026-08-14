package http

import (
	"io"
	"testing"
)

func TestRewind(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		expected := "this is a test"
		response := NewResponseString(expected)

		for i := range 3 {
			contents, err := io.ReadAll(response.Body)
			if err != nil {
				in.Fatalf("could not read from buffer: %v", err)
			}

			if string(contents) != expected {
				in.Errorf("[%d] %s != %s", i, expected, contents)
			}
			rewind(response)
		}

	})
}
