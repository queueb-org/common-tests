package http

import (
	"net/http"
	"reflect"
	"testing"
)

func TestMergeResponseOptions(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		opt := mergeResponseOptions(
			ROptAccepted(),
			ROptBadGateway(),
			ROptWithHeaders("test", "none"),
			&ResponseOption{
				Reload:  new(false),
				Headers: []string{"test", "value"}, // overrides previously set
			},
		)

		expected := &ResponseOption{
			StatusCode: http.StatusBadGateway,
			Headers:    []string{"test", "value"},
			Reload:     new(false),
		}

		if !reflect.DeepEqual(opt, expected) {
			in.Errorf("expected: %#v, got: %#v", expected, opt)
		}
	})
}
