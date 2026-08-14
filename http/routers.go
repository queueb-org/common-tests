package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

// context keeps different settings that might be used during execution.
// it might help for complex clients, that might require different headers
// based on the current context.
type context struct {
	ServerURL string
}

// Value transforms internal [context] fields into simple header value representation.
func (c *context) Value(key, value string) string {
	switch key {
	case "Location":
		return c.ServerURL + value
	default:
		return value
	}
}

// Router represents an http router based on simple endpoints (i.e. /) and given responses.
// Note, [Router] is not performance optimized. It might work slow for intense operations.
// see [Router.match] for implementation details.
//
// TODO: decide if match should be optimized.
type Router map[string]*Response

// ServeHTTP implements http.Handler.
func (r Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	response := r.match(request.URL)
	response.SetHeaders(context{}, writer.Header())
	writer.WriteHeader(response.StatusCode)
	_, _ = io.Copy(writer, response.Body)
	rewind(response)
}

// match iterates over all [Router] keys, tries to compile a regular expression
// and match it against given URL. On success, first matched response would be used.
func (r Router) match(url *url.URL) *Response {
	for key := range r {
		if regex, err := regexp.Compile(key); err == nil {
			uri := buildURI(url)
			if regex.MatchString(uri) {
				return r[key]
			}
		}
	}

	return NotFound
}

// ContextRouter extends Router with a context.
type ContextRouter struct {
	Router
	context context
}

// ServeHTTP implements http.Handler, which implements a simple http server
// that emulates remote server behavior.
func (router *ContextRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	response := router.Router.match(request.URL)
	response.SetHeaders(router.context, writer.Header())
	writer.WriteHeader(response.StatusCode)
	_, _ = io.Copy(writer, response.Body)
	// re-use response.Body if it supports Seek operations.
	rewind(response)
}

// WithServerURL sets context server url
func (router *ContextRouter) WithServerURL(url string) *ContextRouter {
	router.context.ServerURL = url
	return router
}

// StackedRouter is a Router that returns responses as a stack.
// For example, you have a single endpoint, but it's required to
// return different responses. In this case StackedRouter might help.
// Note, that once stack (Responses) are blank NotFound will be used.
type StackedRouter struct {
	Responses []*Response
	context   context
}

// NewStackedRouter initializes StackedRouter with the given responses.
// Might be nil, but in this case NotFound response will be used.
func NewStackedRouter(responses []*Response) *StackedRouter {
	return &StackedRouter{Responses: responses}
}

// pop provides a classic stack item extracting.
func (router *StackedRouter) pop() *Response {
	if router == nil || len(router.Responses) == 0 {
		return NotFound
	}
	response := router.Responses[0]
	if len(router.Responses) == 1 {
		router.Responses = nil
	} else {
		router.Responses = router.Responses[1:]
	}
	return response
}

// WithServerURL sets Server URL to context.
func (router *StackedRouter) WithServerURL(url string) *StackedRouter {
	router.context.ServerURL = url
	return router
}

// ServeHTTP implements http.Handler, which implements a simple http server
// that emulates remote server behavior.
func (router *StackedRouter) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	response := router.pop()
	response.SetHeaders(router.context, writer.Header())
	writer.WriteHeader(response.StatusCode)
	_, _ = io.Copy(writer, response.Body)
}

// NewHTTPServer creates a testing http server, testing.T object is used
// to properly close all resources automatically.
func NewHTTPServer(t testing.TB, router http.Handler) *httptest.Server {
	server := httptest.NewServer(router)

	t.Cleanup(func() {
		server.Close()
	})

	return server
}

// All routers must implement simple http server.
var (
	_ http.Handler = &Router{}
	_ http.Handler = &ContextRouter{}
	_ http.Handler = &StackedRouter{}
)
