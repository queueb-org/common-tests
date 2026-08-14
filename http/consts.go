package http

import (
	"bytes"
	"net/http"
)

// Exportable constants for HTTP:
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes
const (
	// HTTPUnknownError sets http server unknown error.
	HTTPUnknownError = 520
)

// NotFound is a response that returns http non found (404) status with plain text.
// This response is used as a default response, if a test server wasn't configured
// to send any other responses.
var NotFound = &Response{
	Body:       bytes.NewBuffer([]byte("Not found")),
	StatusCode: http.StatusNotFound,
}
