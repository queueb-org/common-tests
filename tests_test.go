package tests

import (
	"testing"
)

func TestPanicExpected(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		//: no reason to check test status due to panic will automatically fail
		//: the test
		defer PanicExpected(in)()
		panic("welcome to the fun palace")
	})

	t.Run("failure", func(in *testing.T) {
		test := &testing.T{}
		defer func() { // invokes secondary
			if !test.Failed() {
				in.Errorf("issue expected but didn't happen")
			}
		}()
		defer PanicExpected(test)() // invokes first
	})
}

func TestPanicNotExpected(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		test := &testing.T{}
		defer func() { // invokes secondary
			if test.Failed() {
				in.Errorf("issue expected but didn't happen")
			}
		}()
		defer PanicNotExpected(in)()
	})

	t.Run("failure", func(in *testing.T) {
		test := &testing.T{}
		defer func() { // invokes secondary
			if !test.Failed() {
				in.Errorf("issue expected but didn't happen")
			}
		}()
		defer PanicNotExpected(test)() // invokes first
		panic("panic!")
	})
}

func TestF(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		test, sub := F(in).WithTestFatalF()
		WithOpenFile(test, "README.rst")
		if sub.Failed() {
			in.Errorf("test expected to succeed")
		}
	})

	t.Run("failed", func(in *testing.T) {
		test, sub := F(in).WithTestFatalF()
		WithOpenFile(test, "non existent")
		if !sub.Failed() {
			in.Errorf("test expected to fail")
		}
	})

}
