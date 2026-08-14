package http

import "net/http"

// ResponseOption keeps different options for responses used in [HTTPServer] objects.
type ResponseOption struct {
	// StatusCode sets response status code, by default its 200 HTTP OK.
	StatusCode int
	// Headers keeps a "key, value" format for headers.
	// Unbalanced slice is valid, in this case a value will be blank.
	Headers []string
	// Reload rewinds contents back after it has been served. Otherwise blank
	// content will be used.
	//
	// TODO: on false there could be [NotFound] use instead of blank contents.
	//       we need to reconsider.
	Reload *bool
}

func mergeResponseOptions(opts ...*ResponseOption) *ResponseOption {
	def := &ResponseOption{
		StatusCode: http.StatusOK,
	}

	for _, opt := range opts {
		if opt.StatusCode != 0 {
			def.StatusCode = opt.StatusCode
		}

		if opt.Headers != nil {
			def.Headers = opt.Headers
		}

		if opt.Reload != nil {
			def.Reload = opt.Reload
		}
	}

	return def
}
