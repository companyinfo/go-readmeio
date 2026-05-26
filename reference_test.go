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

func TestReferenceClient_Create(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/branches/v1/reference", r.URL.Path)

		body, _ := io.ReadAll(r.Body)
		var got ReferenceParams
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, "Pets", got.Title)
		require.NotNil(t, got.Category)
		assert.Equal(t, "/cats/ref", got.Category.URI)
		require.NotNil(t, got.API)
		assert.Equal(t, "GET", got.API.Method)
		assert.Equal(t, "/pets", got.API.Path)

		_, _ = w.Write([]byte(`{"data":{"slug":"pets","title":"Pets","api":{"method":"GET","path":"/pets"}}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	ref, err := c.Reference.Create(context.Background(), "v1", ReferenceParams{
		Title:    "Pets",
		Category: &ResourceRef{URI: "/cats/ref"},
		API:      &ReferenceAPI{Method: "GET", Path: "/pets"},
	})
	require.NoError(t, err)
	assert.Equal(t, "pets", ref.Slug)
	require.NotNil(t, ref.API)
	assert.Equal(t, "GET", ref.API.Method)
	assert.Equal(t, "/pets", ref.API.Path)
}

func TestReferenceClient_List(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/branches/v1/reference/pets", r.URL.Path)
		_, _ = w.Write([]byte(`{"data":{"slug":"pets","title":"Pets"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	ref, err := c.Reference.Get(context.Background(), "v1", "pets")
	require.NoError(t, err)
	assert.Equal(t, "pets", ref.Slug)
}

func TestReferenceClient_Update(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/branches/v1/reference/pets", r.URL.Path)
		_, _ = w.Write([]byte(`{"data":{"slug":"pets","title":"Pets2"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	ref, err := c.Reference.Update(context.Background(), "v1", "pets", ReferenceParams{Title: "Pets2"})
	require.NoError(t, err)
	assert.Equal(t, "Pets2", ref.Title)
}

func TestReferenceClient_Delete(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/branches/v1/reference/pets", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	assert.NoError(t, c.Reference.Delete(context.Background(), "v1", "pets"))
}

func TestReferenceClient_APIError(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"status":422,"title":"Validation","detail":"bad"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, err := c.Reference.Create(context.Background(), "v1", ReferenceParams{
		Title:    "P",
		Category: &ResourceRef{URI: "/c"},
	})
	require.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok, "expected *APIError, got %T (%v)", err, err)
	assert.Equal(t, 422, apiErr.Status)
}

func TestReferenceClient_Validation(t *testing.T) {
	c, _ := New("k")
	valid := ReferenceParams{Title: "T", Category: &ResourceRef{URI: "/c"}}
	cases := []struct {
		name string
		fn   func() error
		want string
	}{
		{"create empty branch", func() error {
			_, e := c.Reference.Create(context.Background(), "", valid)
			return e
		}, "branch"},
		{"create missing title", func() error {
			_, e := c.Reference.Create(context.Background(), "v1", ReferenceParams{Category: &ResourceRef{URI: "/c"}})
			return e
		}, "title"},
		{"create missing category", func() error {
			_, e := c.Reference.Create(context.Background(), "v1", ReferenceParams{Title: "T"})
			return e
		}, "category"},
		{"get empty slug", func() error {
			_, e := c.Reference.Get(context.Background(), "v1", "")
			return e
		}, "slug"},
		{"update empty branch", func() error {
			_, e := c.Reference.Update(context.Background(), "", "s", ReferenceParams{})
			return e
		}, "branch"},
		{"delete empty slug", func() error {
			return c.Reference.Delete(context.Background(), "v1", "")
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
