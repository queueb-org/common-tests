package tests

import (
	"bytes"
	"os"
	"path"
	"testing"
)

const testFile = "testdata/go.mod.txt"

func TestWithReadFile(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		test, sub := F(in).WithTestFatalF()

		result := WithReadFile(test, testFile)
		if sub.Failed() {
			in.Errorf("WithReadFile() should not fail")
			in.FailNow()
		}
		expected, err := os.ReadFile(testFile)
		if !(err == nil && bytes.Equal(result, expected)) {
			in.Errorf("err: `%v`\nexp: `%s`\ngot: `%s`", err, expected, result)
		}
	})

	t.Run("ok/string", func(in *testing.T) {
		test, sub := F(in).WithTestFatalF()
		result := WithReadFileString(test, testFile)

		if sub.Failed() {
			in.Errorf("WithReadFile() should not fail")
			in.FailNow()
		}
		expected, err := os.ReadFile(testFile)
		if !(err == nil && result == string(expected)) {
			in.Errorf("err: `%v`\nexp: `%s`\ngot: `%s`", err, expected, result)
		}
	})

	t.Run("failed", func(in *testing.T) {
		test, sub := F(in).WithTestFatalF()
		WithReadFile(test, "non existent file")
		if !sub.Failed() {
			in.Errorf("%T expected to fail", t.Name())
		}
	})
}

func TestWithWriteFile(t *testing.T) {
	dir := t.TempDir()
	t.Run("ok", func(in *testing.T) {
		filename := path.Join(dir, "a-test-file")
		test, sub := F(in).WithTestFatalF()
		WithWriteFile(test, filename, []byte("this is a test"))
		if sub.Failed() {
			in.Errorf("WithWriteFile() should not fail")
		}
	})

	t.Run("ok/with-sub-dir", func(in *testing.T) {
		filename := path.Join(dir, "sub-dir/a-test-file")
		test, sub := F(in).WithTestFatalF()
		WithWriteFile(test, filename, []byte("this is a test"))
		if sub.Failed() {
			in.Errorf("WithWriteFile() should not fail")
		}
	})

	t.Run("cant-write", func(in *testing.T) {
		filename := path.Join(dir, string([]byte{0x0, 0x1}))
		test, sub := F(in).WithTestFatalF()
		WithWriteFile(test, filename, []byte("this is a test"))

		if !sub.Failed() {
			in.Errorf("WithWriteFile() should fail")
		}
	})

	t.Run("cant-create-dir", func(in *testing.T) {
		filename := path.Join(dir, string([]byte{'/', 0x0, 0x1, '/', 'f', 'i', 'l', 'e'}))
		test, sub := F(in).WithTestFatalF()
		WithWriteFile(test, filename, []byte("this is a test"))

		if !sub.Failed() {
			in.Errorf("WithWriteFile() should fail")
		}
	})
}

func TestTestWithWriteFileString(t *testing.T) {
	dir := t.TempDir()

	t.Run("ok", func(in *testing.T) {
		filename := path.Join(dir, "a-test-file")
		test, sub := F(in).WithTestFatalF()
		WithWriteFileString(test, filename, "this is a test")
		if sub.Failed() {
			in.Errorf("WithWriteFile() should not fail")
		}
	})
}

func TestWithAppendFile(t *testing.T) {
	dir := t.TempDir()
	t.Run("ok", func(in *testing.T) {
		filename := path.Join(dir, "a-test-file")
		test, sub := F(in).WithTestFatalF()
		WithAppendFile(test, filename, []byte("this is a test"))

		if sub.Failed() {
			in.Errorf("WithWriteFile() should not fail")
		}
	})

	t.Run("cant-write", func(in *testing.T) {
		filename := path.Join(dir, string([]byte{0x0, 0x1}))
		test, sub := F(in).WithTestFatalF()
		WithAppendFile(test, filename, []byte("this is a test"))

		if !sub.Failed() {
			in.Errorf("WithWriteFile() should fail")
		}
	})
}

func TestWithCopyFiles(t *testing.T) {
	t.Run("odd-arguments", func(in *testing.T) {
		test, sub := F(in).WithTestFatalF()
		WithCopyFiles(test, "one", "two", "three")
		if !sub.Failed() {
			in.Errorf("test expected to be failed")
		}
	})

	t.Run("ok", func(in *testing.T) {
		srcTempDir := t.TempDir()
		dstTempDir := t.TempDir()
		srcFile := path.Join(srcTempDir, "test.txt")
		dstFile := path.Join(dstTempDir, "a-copied-test.txt")
		expected := []byte("this is\na test file")

		WithWriteFile(in, srcFile, expected)

		test, sub := F(t).WithTestFatalF()
		WithCopyFiles(test, srcFile, dstFile)
		if sub.Failed() {
			in.Errorf("test should not be failed")
		}
		if result := WithReadFile(in, dstFile); !bytes.Equal(result, expected) {
			in.Errorf("expected: %v, got: %v\n", expected, result)
		}
	})
}

func TestWithOpenFile(t *testing.T) {
	var fd *os.File

	t.Run("ok", func(in *testing.T) {
		fd = WithOpenFile(in, testFile)
		if fd == nil {
			in.Errorf("expected *os.File, got nil instead")
		}
	})

	t.Run("ok/cant-open-file", func(in *testing.T) {
		test, sub := F(t).WithTestFatalF()
		fd = WithOpenFile(test, "non existent")

		if fd != nil {
			in.Errorf("expected nil got: %v", fd)
		}
		if !sub.Failed() {
			in.Errorf("test expected to be failed")
		}
	})
}
