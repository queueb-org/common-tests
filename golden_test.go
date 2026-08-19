package tests

import (
	"os"
	"path"
	"testing"
	"time"
)

func WithLocationPrefix(t testing.TB, location string) {
	orig := locationPrefix
	locationPrefix = location

	t.Cleanup(func() {
		locationPrefix = orig
	})
}

func WithUpdate(t testing.TB, v bool) {
	origin := update
	update = &v

	t.Cleanup(func() {
		update = origin
	})
}

func TestGolden(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		tempDir := in.TempDir()
		WithLocationPrefix(in, tempDir+"/")

		sub := &testing.T{}
		location := goldenName(sub)

		expected := "this is a test"
		if err := os.MkdirAll(path.Dir(location), defaultDirPerm); err != nil {
			in.Errorf("could not create directory: %v", err)
			in.FailNow()
		}

		// put expected result
		WithWriteFile(in, location, []byte(expected))
		// Read file
		if result := Golden(sub, nil); result != expected {
			in.Errorf("expected: `%v`, got: `%v`", expected, result)
		}
	})

	t.Run("on-err", func(in *testing.T) {
		test, sub := F(in).WithTestFatalF()
		Golden(test, nil)

		if !sub.Failed() {
			in.Errorf("Golden() expected to fail")
		}
	})

	t.Run("on-update-ok", func(in *testing.T) {
		WithUpdate(in, true)
		tempDir := in.TempDir()
		WithLocationPrefix(in, tempDir+"/")

		test, sub := F(in).WithTestFatalF()
		location := goldenName(sub)

		if err := os.MkdirAll(path.Dir(location), defaultDirPerm); err != nil {
			in.Errorf("could not create directory: %v", err)
			in.FailNow()
		}

		Golden(test, []byte("contents to write"))
		time.Sleep(time.Millisecond * 15)
		if !sub.Skipped() {
			in.Errorf("test expected to be skipped")
		}
	})

	t.Run("on-update-err", func(in *testing.T) {
		WithUpdate(in, true)
		tempDir := in.TempDir()
		WithLocationPrefix(in, tempDir+"/")

		test, sub := F(in).WithTestFatalF()
		Golden(test, []byte("contents to write"))

		if !sub.Failed() {
			in.Errorf("test expected to be failed")
		}
	})
}
