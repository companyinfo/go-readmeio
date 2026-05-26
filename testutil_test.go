package readme

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient builds a readme.Client pointing at the provided httptest server.
func newTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	c, err := New("test-token", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	return c
}

// jsonHandler wraps an http.HandlerFunc to set a default application/json
// Content-Type on the response (required by resty to auto-decode results).
func jsonHandler(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	})
}
