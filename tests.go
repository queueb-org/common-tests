package tests

import (
	"sync"
	"testing"
)

// PanicExpected returns a callback that checks if function has been failed with
// a panic call, otherwise it fails the test.
func PanicExpected(t testing.TB) func() {
	return func() {
		r := recover()
		if r == nil {
			t.Errorf("didn't panic")
			t.Fail()
		}
	}
}

// PanicNotExpected returns a callback that checks if function hasn't been failed with
// a panic call, otherwise it fails the test.
func PanicNotExpected(t testing.TB) func() {
	return func() {
		if r := recover(); r != nil {
			t.Errorf("got panic: %v", r)
			t.Fail()
		}
	}
}

// ErrorFormatFunc keeps [testing.T] error formatting functions signature.
type ErrorFormatFunc = func(format string, args ...interface{})

// F creates testing function helper.
// Example:
//
//	func TestOpenFile(t *testing.T){
//		test, sub := F(t).WithTestFatalF()
//
//		WithOpenFile(test, "non existent file")
//		if !sub.Failed() {
//			in.Errorf("test expected to fail")
//		}
//	}
func F(t testing.TB) *f {
	return &f{
		TB:     t,
		fatalf: t.Fatalf,
	}
}

// f represents a helper to deal with Fatalf calls upon testing objects.
type f struct {
	testing.TB

	fatalf ErrorFormatFunc
	mu     sync.RWMutex
}

// Fatalf represents [testing.TB.Fatalf] adapter applicable for internal testing.
func (f *f) Fatalf(format string, args ...interface{}) {
	f.fatalf(format, args...)
}

// WithTestFatalF overrides Fatalf calls with safe ones and returns testing object
// that will be used (and record met errors).
func (f *f) WithTestFatalF() (*f, *testing.T) {
	f.mu.Lock()
	orig := f.fatalf

	f.TB.Cleanup(func() {
		f.fatalf = orig
	})
	f.mu.Unlock()

	sub := &testing.T{}
	f.fatalf = sub.Errorf
	return f, sub
}
