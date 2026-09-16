package readme

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIDefinitionClient_Create(t *testing.T) {
	specContent := `{"openapi":"3.0.0","info":{"title":"Test API","version":"1.0.0"},"paths":{}}`
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/branches/v1/apis", r.URL.Path)
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

		err := r.ParseMultipartForm(32 << 20)
		require.NoError(t, err)

		f, _, err := r.FormFile("schema")
		require.NoError(t, err)
		defer f.Close()
		content, _ := io.ReadAll(f)
		assert.Equal(t, specContent, string(content))

		_, _ = w.Write([]byte(`{"data":{"id":"api-123","version":"v1","title":"Test API"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	err := c.APIDefinitions.Create(context.Background(), "v1", APIDefinitionParams{
		Schema:   specContent,
		FileName: "spec.json",
	})
	require.NoError(t, err)
}

func TestAPIDefinitionClient_Get(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/branches/v1/apis/petstore.json", r.URL.Path)
		_, _ = w.Write([]byte(`{"data":{"filename":"petstore.json","version":"v1","title":"Petstore"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	apiDef, err := c.APIDefinitions.Get(context.Background(), "v1", "petstore.json")
	require.NoError(t, err)
	assert.Equal(t, "petstore.json", apiDef.Filename)
	assert.Equal(t, "Petstore", apiDef.Title)
}

func TestAPIDefinitionClient_Update(t *testing.T) {
	specContent := `{"openapi":"3.0.0","info":{"title":"Updated API","version":"1.0.0"},"paths":{}}`
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/branches/v1/apis/petstore.json", r.URL.Path)
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

		err := r.ParseMultipartForm(32 << 20)
		require.NoError(t, err)

		f, _, err := r.FormFile("schema")
		require.NoError(t, err)
		defer f.Close()
		content, _ := io.ReadAll(f)
		assert.Equal(t, specContent, string(content))

		_, _ = w.Write([]byte(`{"data":{"id":"petstore.json","version":"v1","title":"Updated API"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	err := c.APIDefinitions.Update(context.Background(), "v1", APIDefinitionParams{
		Schema:   specContent,
		FileName: "petstore.json",
	})
	require.NoError(t, err)
}

func TestAPIDefinitionClient_Delete(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/branches/v1/apis/petstore.json", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	err := c.APIDefinitions.Delete(context.Background(), "v1", "petstore.json")
	require.NoError(t, err)
}

func TestAPIDefinitionClient_Validate(t *testing.T) {
	specContent := `{"openapi":"3.0.0","info":{"title":"Test API","version":"1.0.0"},"paths":{}}`
	srv := httptest.NewServer(jsonHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/validate/api", r.URL.Path)
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

		err := r.ParseMultipartForm(32 << 20)
		require.NoError(t, err)

		f, _, err := r.FormFile("schema")
		require.NoError(t, err)
		defer f.Close()
		content, _ := io.ReadAll(f)
		assert.Equal(t, specContent, string(content))

		_, _ = w.Write([]byte(`{"data":{"valid":true,"warnings":["Missing description"]}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	val, err := c.APIDefinitions.Validate(context.Background(), "v1", APIDefinitionParams{
		Schema:   specContent,
		FileName: "spec.json",
	})
	require.NoError(t, err)
	assert.Contains(t, val.Data, "Missing description")
}

func TestAPIDefinitionClient_Validation(t *testing.T) {
	c, _ := New("key")
	cases := []struct {
		name   string
		branch string
		params APIDefinitionParams
		want   string
	}{
		{
			name:   "empty branch",
			branch: "",
			params: APIDefinitionParams{Schema: "{}", FileName: "spec.json"},
			want:   "branch is required",
		},
		{
			name:   "missing schema and url",
			branch: "v1",
			params: APIDefinitionParams{},
			want:   "schema must satisfy required_without_all=UploadSource Url",
		},
		{
			name:   "schema without filename",
			branch: "v1",
			params: APIDefinitionParams{Schema: "{}"},
			want:   "filename must satisfy required_with=Schema",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := c.APIDefinitions.Create(context.Background(), tc.branch, tc.params)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}
