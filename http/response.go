package http

import (
	"io"
	"net/http"
	"strings"
)

// Response reflect a simple object which is used in [Router] responses.
type Response struct {
	// StatusCode keeps http return status code, e.g. HTTP 200 OK.
	StatusCode int
	// Headers in a plain form, a slice can be unbalanced, in this case blank value is used.
	Headers []string
	// Body represents response contents. For testing purposes there's only [io.Reader]
	// requirement. A consumer should maintain Close operation on their own.
	Body io.Reader
	// Reload keeps a switch if [Body] has to be reloaded (i.e. seek back to the beginning).
	Reload bool
}

// NewResponse creates Response with any given options, see ROpt... functions
// as source of available options.
func NewResponse(reader io.Reader, options ...*ResponseOption) *Response {
	opt := mergeResponseOptions(options...)

	return &Response{
		Body:       reader,
		StatusCode: opt.StatusCode,
		Headers:    opt.Headers,
		Reload:     deref(opt.Reload, true),
	}
}

// NewResponseString is the same as NewResponse, src string will be used.
func NewResponseString(src string, options ...*ResponseOption) *Response {
	return NewResponse(strings.NewReader(src), options...)
}

// rewind rewinds reader to its beginning if it's supported.
func rewind(r *Response) {
	if rewinder, ok := r.Body.(io.Seeker); ok && r.Reload {
		// seek file to its beginning for its further use.
		_, _ = rewinder.Seek(0, io.SeekStart)
	}
}

// Option builders

// ROptStatusOK : HTTP 200 OK
func ROptStatusOK() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusOK}
}

// ROptWithOK alias for [ROptStatusOK].
func ROptWithOK() *ResponseOption {
	return ROptStatusOK()
}

// ROptCreated : HTTP 201 CREATED
func ROptCreated() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusCreated}
}

// ROptAccepted : HTTP 202 ACCEPTED
func ROptAccepted() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusAccepted}
}

// ROptNoContent : HTTP 204 NO CONTENT
func ROptNoContent() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusNoContent}
}

// ROptBadRequest : 400 BAD REQUEST
func ROptBadRequest() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusBadRequest}
}

// ROptUnauthorized : 401 UNAUTHORIZED
func ROptUnauthorized() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusUnauthorized}
}

// ROptForbidden : HTTP 403 FORBIDDEN
func ROptForbidden() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusForbidden}
}

// ROptNotFound : 404 NOT FOUND
func ROptNotFound() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusNotFound}
}

// ROptServerError : 500 INTERNAL SERVER ERROR
func ROptServerError() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusInternalServerError}
}

// ROptNotImplemented : 501 NOT IMPLEMENTED
func ROptNotImplemented() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusNotImplemented}
}

// ROptBadGateway : 502 BAD GATEWAY
func ROptBadGateway() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusBadGateway}
}

// ROptServiceUnavailable : 503 SERVICE UNAVAILABLE
func ROptServiceUnavailable() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusServiceUnavailable}
}

// ROptGatewayTimeout : 503 GATEWAY TIMEOUT
func ROptGatewayTimeout() *ResponseOption {
	return &ResponseOption{StatusCode: http.StatusGatewayTimeout}
}

// ROptUnknownError : 520 Web Server Returned an Unknown Error.
func ROptUnknownError() *ResponseOption {
	return &ResponseOption{StatusCode: HTTPUnknownError}
}

// ROptWithHeaders sets http headers using key, value repeating arguments.
//
// Example:
//
//	ROptWithHeaders("X-Key", "value", "X-Ref")  // the same as below
//	ROptWithHeaders("X-Key", "value", "X-Ref", "")  // the same as above
func ROptWithHeaders(keyValues ...string) *ResponseOption {
	return &ResponseOption{Headers: keyValues}
}
