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

func TestCategoryClient_Create(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/branches/v1/categories", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		body, _ := io.ReadAll(r.Body)
		var got CategoryParams
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, "Intro", got.Title)
		assert.Equal(t, CategoryTypeGuides, got.Section)

		_, _ = w.Write([]byte(`{"data":{"title":"Intro","section":"guides","uri":"/categories/guides/intro"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	cat, err := c.Categories.Create(context.Background(), "v1", CategoryParams{Title: "Intro", Section: CategoryTypeGuides})
	require.NoError(t, err)
	assert.Equal(t, "Intro", cat.Title)
	assert.Equal(t, "/categories/guides/intro", cat.URI)
}

func TestCategoryClient_List(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/branches/v1/categories/guides", r.URL.Path)
		_, _ = w.Write([]byte(`{"data":[{"title":"A","section":"guides","uri":"/a"},{"title":"B","section":"guides","uri":"/b"}]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	cats, err := c.Categories.List(context.Background(), "v1", "guides")
	require.NoError(t, err)
	require.Len(t, cats, 2)
	assert.Equal(t, "A", cats[0].Title)
	assert.Equal(t, "B", cats[1].Title)
}

func TestCategoryClient_GetByTitle(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/branches/v1/categories/guides/Intro", r.URL.Path)
		_, _ = w.Write([]byte(`{"data":{"title":"Intro","section":"guides","uri":"/x"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	cat, err := c.Categories.GetByTitle(context.Background(), "v1", "guides", "Intro")
	require.NoError(t, err)
	assert.Equal(t, "Intro", cat.Title)
}

func TestCategoryClient_Update(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/branches/v1/categories/guides/Intro", r.URL.Path)
		_, _ = w.Write([]byte(`{"data":{"title":"Intro2","section":"guides","uri":"/x"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	cat, err := c.Categories.Update(context.Background(), "v1", "guides", "Intro", CategoryParams{Title: "Intro2"})
	require.NoError(t, err)
	assert.Equal(t, "Intro2", cat.Title)
}

func TestCategoryClient_Delete(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/branches/v1/categories/guides/Intro", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	assert.NoError(t, c.Categories.Delete(context.Background(), "v1", "guides", "Intro"))
}

func TestCategoryClient_APIError(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"status":422,"title":"Validation","detail":"bad","errors":[{"key":"title","message":"required"}]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, err := c.Categories.Create(context.Background(), "v1", CategoryParams{Title: "x", Section: CategoryTypeGuides})
	require.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok, "expected *APIError, got %T", err)
	assert.Equal(t, 422, apiErr.Status)
	require.Len(t, apiErr.Errors, 1)
	assert.Equal(t, "title", apiErr.Errors[0].Key)
}

func TestCategoryClient_Validation(t *testing.T) {
	c, _ := New("k")
	cases := []struct {
		name string
		fn   func() error
		want string
	}{
		{"create empty branch", func() error {
			_, e := c.Categories.Create(context.Background(), "", CategoryParams{Title: "t", Section: CategoryTypeGuides})
			return e
		}, "branch"},
		{"create empty title", func() error {
			_, e := c.Categories.Create(context.Background(), "v1", CategoryParams{Section: CategoryTypeGuides})
			return e
		}, "title"},
		{"create bad section", func() error {
			_, e := c.Categories.Create(context.Background(), "v1", CategoryParams{Title: "t", Section: "bogus"})
			return e
		}, "section"},
		{"get bad section", func() error {
			_, e := c.Categories.List(context.Background(), "v1", "bogus")
			return e
		}, "section"},
		{"get-by-title empty title", func() error {
			_, e := c.Categories.GetByTitle(context.Background(), "v1", "guides", "")
			return e
		}, "title"},
		{"delete empty title", func() error {
			return c.Categories.Delete(context.Background(), "v1", "guides", "")
		}, "title"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}
