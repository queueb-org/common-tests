package tests

import (
	"os"
	"testing"
)

const (
	defaultFileMode = os.FileMode(0640)
)

// WithReadFile reads a file and returns its content on success.
func WithReadFile(t testing.TB, filename string) (contents []byte) {
	var err error
	if contents, err = os.ReadFile(filename); err == nil {
		return
	}

	t.Fatalf("WithReadFile(): %v", err)
	return
}

// WithReadFileString reads a file and returns its content on success.
func WithReadFileString(t testing.TB, filename string) (contents string) {
	return string(WithReadFile(t, filename))
}

// WithWriteFile writes contents to filename.
func WithWriteFile(t testing.TB, filename string, contents []byte) {
	if err := os.WriteFile(filename, contents, defaultFileMode); err != nil {
		t.Fatalf("WithWriteFile(): %v", err)
	}
}

// WithAppendFile appends new contents to a file.
func WithAppendFile(t testing.TB, filename string, contents []byte) {
	mode := os.O_RDWR | os.O_CREATE | os.O_APPEND
	var (
		err error
		fd  *os.File
	)

	if fd, err = os.OpenFile(filename, mode, defaultFileMode); err == nil {
		if _, err = fd.Seek(0, 2); err == nil {
			_, err = fd.Write(contents)
		}
	}

	if err != nil {
		t.Fatalf("WithAppendFile(): %v", err)
	}
}

// WithCopyFiles copies files from one location to another
// pack reflects: src1, dst1, src2, dst2, ..., srcN, dstN
func WithCopyFiles(t testing.TB, pack ...string) {
	if len(pack)%2 != 0 {
		t.Fatalf("amount of arguments should be even (not odd)")
		return
	}

	for i := 0; i < len(pack); i += 2 {
		source, dest := pack[i], pack[i+1]
		WithWriteFile(t, dest, WithReadFile(t, source))
	}
}

// WithOpenFile opens a file and returns its descriptor, once test is over the file
// will be automatically closed.
func WithOpenFile(t testing.TB, filename string) *os.File {
	fd, err := os.Open(filename)
	if err != nil {
		t.Fatalf("could not open file (%s): %v", filename, err)
		return nil
	}

	t.Cleanup(func() {
		if fd != nil {
			_ = fd.Close()
		}
	})
	return fd
}
