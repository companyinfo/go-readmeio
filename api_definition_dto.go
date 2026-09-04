package readme

// APIDefinition models the response from the ReadMe API v2 for API definitions.
//
// See: https://docs.readme.com/main/reference/getapis
type APIDefinition struct {
	ID      string `json:"id"`
	Version string `json:"version"` // branch/version slug
	Title   string `json:"title"`
}

// APIDefinitionParams is the request body for creating/updating/validating an API definition.
//
// See:
//   - https://docs.readme.com/main/reference/createapi
//   - https://docs.readme.com/main/reference/updateapi
//   - https://docs.readme.com/main/reference/validateapi
type APIDefinitionParams struct {
	// Schema is the OpenAPI or Swagger specification as a string (JSON or YAML).
	FileName     string `json:"filename"`
	Schema       string `json:"schema"`
	UploadSource string `json:"upload_source"`
	Url          string `json:"url"`
}

// APIDefinitionValidation represents the result of validating an API definition.
type APIDefinitionValidation struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}
