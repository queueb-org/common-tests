package http

import "testing"

func TestDeref(t *testing.T) {
	for _, test := range []struct {
		name     string
		p        *string
		defaults string
		expected string
	}{
		{"ok", new("ok"), "defaults", "ok"},
		{"nil", nil, "defaults", "defaults"},
	} {
		t.Run(test.name, func(in *testing.T) {
			if result := deref(test.p, test.defaults); result != test.expected {
				in.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
