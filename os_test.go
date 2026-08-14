package tests

import (
	"os"
	"reflect"
	"testing"
)

func TestWithArgs(t *testing.T) {
	origin := os.Args

	// test set.
	t.Run("set", func(in *testing.T) {
		expected := []string{"this", "is", "a", "test"}
		WithArgs(in, expected)
		if !reflect.DeepEqual(os.Args, expected) {
			in.Errorf("expected: `%#+v`, got: `%#+v`", expected, os.Args)
		}
	})

	// test rollback
	if !reflect.DeepEqual(os.Args, origin) {
		t.Errorf("expected: `%#+v`, got: `%#+v`", origin, os.Args)
	}
}
