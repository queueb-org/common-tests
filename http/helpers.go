package http

import (
	"fmt"
	"net/http"
	"net/url"
)

// deref safely de-reference a pointer, if pointer is nil defaults is used.
func deref[T comparable](p *T, defaults T) T {
	if p == nil {
		return defaults
	}

	return *p
}

func buildURI(url *url.URL) string {
	if url.RawQuery != "" {
		return fmt.Sprintf("%s?%s", url.Path, url.RawQuery)
	}

	return url.Path
}

// SetHeaders write headers regarding the context.
// In case of unbalanced Headers slice the blank value will be used.
func (r *Response) SetHeaders(c context, header http.Header) {
	outOfBounds := len(r.Headers)

	for i := 0; i < outOfBounds; i += 2 {
		key := r.Headers[i]
		value := ""
		if i <= outOfBounds {
			value = r.Headers[+1]
		}
		header.Set(key, c.Value(key, value))
	}
}
