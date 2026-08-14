package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime"
	"testing"
)

func FuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type OptionFunc func() *ResponseOption

func TestWithStatuses_Suite(t *testing.T) {
	for _, entry := range []struct {
		in       OptionFunc
		expected *ResponseOption
	}{
		{ROptWithOK, &ResponseOption{StatusCode: http.StatusOK}},
		{ROptStatusOK, &ResponseOption{StatusCode: http.StatusOK}},
		{ROptCreated, &ResponseOption{StatusCode: http.StatusCreated}},
		{ROptAccepted, &ResponseOption{StatusCode: http.StatusAccepted}},
		{ROptNoContent, &ResponseOption{StatusCode: http.StatusNoContent}},
		{ROptBadRequest, &ResponseOption{StatusCode: http.StatusBadRequest}},
		{ROptUnauthorized, &ResponseOption{StatusCode: http.StatusUnauthorized}},
		{ROptForbidden, &ResponseOption{StatusCode: http.StatusForbidden}},
		{ROptNotFound, &ResponseOption{StatusCode: http.StatusNotFound}},
		{ROptServerError, &ResponseOption{StatusCode: http.StatusInternalServerError}},
		{ROptNotImplemented, &ResponseOption{StatusCode: http.StatusNotImplemented}},
		{ROptServiceUnavailable, &ResponseOption{StatusCode: http.StatusServiceUnavailable}},
		{ROptGatewayTimeout, &ResponseOption{StatusCode: http.StatusGatewayTimeout}},
		{ROptUnknownError, &ResponseOption{StatusCode: HTTPUnknownError}},
	} {
		t.Run(FuncName(entry.in), func(in *testing.T) {
			result := entry.in()
			if !reflect.DeepEqual(result, entry.expected) {
				in.Errorf("expected: %#v, got: %#v", entry.expected, result)
			}
		})
	}
}

func TestWithHeaders(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       []string
		expected *ResponseOption
	}{
		{
			name:     "nil",
			in:       nil,
			expected: &ResponseOption{},
		},
		{
			name:     "ok",
			in:       []string{"key", "value", "key2", "value"},
			expected: &ResponseOption{Headers: []string{"key", "value", "key2", "value"}},
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if result := ROptWithHeaders(entry.in...); !reflect.DeepEqual(result, entry.expected) {
				in.Errorf("\nexp: %#v\ngot: %#v\n", entry.expected, result)
			}
		})
	}
}

func TestNewHTTPServer(t *testing.T) {
	var server *httptest.Server
	// a subtest is used to get coverage for test cleanup
	t.Run("ok", func(in *testing.T) {
		server = NewHTTPServer(in, Router{})
		if server == nil {
			in.Errorf("server should not be nil")
		}
	})
}

func TestStackedRouter_ServeHTTP(t *testing.T) {
	var responses []*Response
	for i := 0; i < 10; i++ {
		responses = append(responses, NewResponseString(fmt.Sprintf("%d", i)))
	}
	router := NewStackedRouter(responses)
	for i := 0; i < 10; i++ {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, nil)
		if result := recorder.Result().StatusCode; result != http.StatusOK {
			t.Errorf("expected status: %d, got: %d", http.StatusOK, result)
		}
	}
	// other requests should return 404 (NotFound response)
	for i := 0; i < 10; i++ {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, nil)
		if result := recorder.Result().StatusCode; result != http.StatusNotFound {
			t.Errorf("expected status: %d, got: %d", http.StatusNotFound, result)
		}
	}
}

func TestStackedRouter_WithServerURL(t *testing.T) {
	expected := "https://fake.fqdn"
	router := NewStackedRouter(nil).WithServerURL(expected)
	if result := router.context.ServerURL; result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	}
}

func TestContextRouter_ServeHTTP(t *testing.T) {
	response := NewResponseString("{}", ROptWithHeaders("Location", "/test"))
	router := &ContextRouter{
		context: context{ServerURL: "https://my-shiny-server"},
		Router:  Router{"/test": response},
	}
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(recorder, request)
	resp := recorder.Result()

	expectedLocation := "https://my-shiny-server/test"
	result := resp.Header.Get("Location")
	if result != expectedLocation {
		t.Errorf("expected: %v, got: %v", expectedLocation, result)
	}
}

func TestContextRouter_WithServerURL(t *testing.T) {
	expected := "https://fake.fqdn"
	router := (&ContextRouter{}).WithServerURL(expected)
	if result := router.context.ServerURL; result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	}
}

func TestRouter_ServeHTTP(t *testing.T) {
	expected := "this is the contents"
	response := NewResponseString(expected)
	router := &Router{"/test": response}
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)

	for i := 0; i < 10; i++ {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		resp := recorder.Result()
		result, _ := io.ReadAll(resp.Body)
		if string(result) != expected {
			t.Errorf("[%d] expected: %v, got: %q", i, expected, result)
		}
	}
}

func TestRouter_match(t *testing.T) {
	router := Router{
		`/test`:  NewResponseString("test"),
		"(?P<]}": NewResponseString("broken path"),
	}
	for _, entry := range []struct {
		name     string
		in       string
		expected *Response
	}{
		{
			name:     "found",
			in:       "/test",
			expected: router[`/test`],
		},
		{
			name:     "not-found",
			in:       "/blank",
			expected: NotFound,
		},
		{
			name:     "broken-regex",
			in:       `(?P<])`,
			expected: NotFound,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			u, _ := url.Parse(entry.in)
			if result := router.match(u); !reflect.DeepEqual(result, entry.expected) {
				in.Errorf("\nexp: %#v\ngot: %#v\n", entry.expected, result)
			}
		})
	}
}

func TestContext_Value(t *testing.T) {
	ctx := context{
		ServerURL: "https://fake.fqdn",
	}
	type args struct {
		key, value string
	}
	for _, entry := range []struct {
		name     string
		args     args
		expected string
	}{
		{
			name:     "defaults",
			args:     args{"test", "value"},
			expected: "value",
		},
		{
			name:     "set",
			args:     args{"Location", "/v2/test"},
			expected: "https://fake.fqdn/v2/test",
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if result := ctx.Value(entry.args.key, entry.args.value); !reflect.DeepEqual(result, entry.expected) {
				in.Errorf("\nexp: %#v\ngot: %#v\n", entry.expected, result)
			}
		})
	}
}
