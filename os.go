package tests

import (
	"os"
	"testing"
)

// WithArgs sets os.Args with given set, once test is over it returns os.Args
// to the original state.
func WithArgs(t testing.TB, set []string) {
	origin := os.Args
	os.Args = set

	t.Cleanup(func() {
		os.Args = origin
	})
}
