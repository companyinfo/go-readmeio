package readme

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuideClient_Create(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/branches/v1/guides", r.URL.Path)

		body, _ := io.ReadAll(r.Body)
		var got GuideParams
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, "Hello", got.Title)
		require.NotNil(t, got.Category)
		assert.Equal(t, "/cats/x", got.Category.URI)

		_, _ = w.Write([]byte(`{"data":{"slug":"hello","title":"Hello"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	g, err := c.Guides.Create(context.Background(), "v1", GuideParams{
		Title:    "Hello",
		Category: &ResourceRef{URI: "/cats/x"},
	})
	require.NoError(t, err)
	assert.Equal(t, "hello", g.Slug)
	assert.Equal(t, "Hello", g.Title)
}

func TestGuideClient_List(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/branches/v1/guides/hello", r.URL.Path)
		_, _ = w.Write([]byte(`{"data":{"slug":"hello","title":"Hello"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	g, err := c.Guides.Get(context.Background(), "v1", "hello")
	require.NoError(t, err)
	assert.Equal(t, "hello", g.Slug)
}

func TestGuideClient_Update(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/branches/v1/guides/hello", r.URL.Path)
		_, _ = w.Write([]byte(`{"data":{"slug":"hello","title":"Hello2"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	g, err := c.Guides.Update(context.Background(), "v1", "hello", GuideParams{Title: "Hello2"})
	require.NoError(t, err)
	assert.Equal(t, "Hello2", g.Title)
}

func TestGuideClient_Delete(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/branches/v1/guides/hello", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	assert.NoError(t, c.Guides.Delete(context.Background(), "v1", "hello"))
}

func TestGuideClient_APIError(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status":404,"title":"Not Found","detail":"missing"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, err := c.Guides.Get(context.Background(), "v1", "hello")
	require.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok, "expected *APIError, got %T (%v)", err, err)
	assert.Equal(t, 404, apiErr.Status)
}

func TestGuideClient_Validation(t *testing.T) {
	c, _ := New("k")
	validParams := GuideParams{Title: "T", Category: &ResourceRef{URI: "/c"}}
	cases := []struct {
		name string
		fn   func() error
		want string
	}{
		{"create empty branch", func() error {
			_, e := c.Guides.Create(context.Background(), "", validParams)
			return e
		}, "branch"},
		{"create missing title", func() error {
			_, e := c.Guides.Create(context.Background(), "v1", GuideParams{Category: &ResourceRef{URI: "/c"}})
			return e
		}, "title"},
		{"create missing category", func() error {
			_, e := c.Guides.Create(context.Background(), "v1", GuideParams{Title: "T"})
			return e
		}, "category"},
		{"get empty slug", func() error {
			_, e := c.Guides.Get(context.Background(), "v1", "")
			return e
		}, "slug"},
		{"update empty branch", func() error {
			_, e := c.Guides.Update(context.Background(), "", "s", GuideParams{})
			return e
		}, "branch"},
		{"delete empty slug", func() error {
			return c.Guides.Delete(context.Background(), "v1", "")
		}, "slug"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}
