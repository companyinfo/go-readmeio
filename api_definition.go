package readme

import (
	"context"
	"errors"
	"strings"

	"github.com/go-resty/resty/v2"
)

// apiDefinitionResponse wraps single-item responses from the API Definition API.
type apiDefinitionResponse struct {
	Data APIDefinition `json:"data"`
}

// APIDefinitionService defines API definition operations against the ReadMe API v2.
//
// See:
//   - https://docs.readme.com/main/reference/createapi
//   - https://docs.readme.com/main/reference/getapi
//   - https://docs.readme.com/main/reference/updateapi
//   - https://docs.readme.com/main/reference/deleteapi
//   - https://docs.readme.com/main/reference/validateapi
type APIDefinitionService interface {
	// Create uploads a new API definition (OpenAPI/Swagger) to the given branch.
	Create(ctx context.Context, branch string, params APIDefinitionParams) error
	// Get retrieves a single API definition by its filename.
	Get(ctx context.Context, branch, filename string) (*APIDefinition, error)
	// Update updates an existing API definition identified by its filename.
	Update(ctx context.Context, branch string, params APIDefinitionParams) error
	// Delete removes an API definition identified by its filename.
	Delete(ctx context.Context, branch, filename string) error
	// Validate validates an API definition without uploading it.
	Validate(ctx context.Context, branch string, params APIDefinitionParams) (ApiDefinitionValidationResponse, error)
}

// ApiDefinitionValidationResponse wraps responses from the API Definition Validation API.
type ApiDefinitionValidationResponse struct {
	Data string `json:"data"`
}

// APIDefinitionClient implements APIDefinitionService.
type APIDefinitionClient struct {
	client *Client
}

// NewAPIDefinitionClient returns a new APIDefinitionClient bound to the provided root client.
func NewAPIDefinitionClient(c *Client) *APIDefinitionClient {
	return &APIDefinitionClient{client: c}
}

// Create uploads a new API definition to the given branch.
// ReadMe API v2: POST /branches/{branch}/api-specification
func (a *APIDefinitionClient) Create(ctx context.Context, branch string, params APIDefinitionParams) error {
	if err := validateBranch(branch); err != nil {
		return err
	}
	if err := validateParams(params); err != nil {
		return err
	}

	var out apiDefinitionResponse
	req := a.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
		}).
		SetResult(&out).
		SetError(&APIError{})

	a.setMultipartFields(req, params)

	resp, err := req.Post("/branches/{branch}/apis")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return apiErrorFromResponse(resp)
	}
	return nil
}

func (a *APIDefinitionClient) setMultipartFields(req *resty.Request, params APIDefinitionParams) {
	if params.Schema != "" {
		req.SetMultipartField("schema", params.FileName, "application/json",
			strings.NewReader(params.Schema))
	}
	if params.Url != "" {
		req.SetMultipartFormData(map[string]string{"url": params.Url})
	}
	if params.UploadSource != "" {
		req.SetMultipartFormData(map[string]string{"upload_source": params.UploadSource})
	}
}

// Get retrieves a single API definition by its filename.
// ReadMe API v2: GET /branches/{branch}/apis/{filename}
func (a *APIDefinitionClient) Get(ctx context.Context, branch, filename string) (*APIDefinition, error) {
	if err := validateBranch(branch); err != nil {
		return nil, err
	}
	// filename is treated as a slug/identifier in the path.
	if filename == "" {
		return nil, errors.New("filename is required")
	}

	var out apiDefinitionResponse
	resp, err := a.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch":   branch,
			"filename": filename,
		}).
		SetResult(&out).
		SetError(&APIError{}).
		Get("/branches/{branch}/apis/{filename}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, apiErrorFromResponse(resp)
	}
	return &out.Data, nil
}

// Update updates an existing API definition identified by its filename.
// ReadMe API v2: PUT /branches/{branch}/apis/{filename}
func (a *APIDefinitionClient) Update(ctx context.Context, branch string, params APIDefinitionParams) error {
	if err := validateBranch(branch); err != nil {
		return err
	}
	if params.FileName == "" {
		return errors.New("filename is required")
	}
	if err := validateParams(params); err != nil {
		return err
	}

	var out apiDefinitionResponse
	req := a.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch":   branch,
			"filename": params.FileName,
		}).
		SetResult(&out).
		SetError(&APIError{})

	a.setMultipartFields(req, params)

	resp, err := req.Put("/branches/{branch}/apis/{filename}")

	if err != nil {
		return err
	}
	if resp.IsError() {
		return apiErrorFromResponse(resp)
	}
	return nil
}

// Delete removes an API definition identified by its filename.
// ReadMe API v2: DELETE /branches/{branch}/apis/{filename}
func (a *APIDefinitionClient) Delete(ctx context.Context, branch, filename string) error {
	if err := validateBranch(branch); err != nil {
		return err
	}
	if filename == "" {
		return errors.New("filename is required")
	}

	resp, err := a.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch":   branch,
			"filename": filename,
		}).
		SetError(&APIError{}).
		Delete("/branches/{branch}/apis/{filename}")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return apiErrorFromResponse(resp)
	}
	return nil
}

// Validate validates an API definition without uploading it.
// ReadMe API v2: POST /validate/api
func (a *APIDefinitionClient) Validate(ctx context.Context, branch string, params APIDefinitionParams) (ApiDefinitionValidationResponse, error) {
	if err := validateBranch(branch); err != nil {
		return ApiDefinitionValidationResponse{Data: ""}, err
	}
	if err := validateParams(params); err != nil {
		return ApiDefinitionValidationResponse{Data: ""}, err
	}

	req := a.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
		}).
		SetError(&APIError{})

	a.setMultipartFields(req, params)

	resp, err := req.Post("/validate/api")
	if err != nil {
		return ApiDefinitionValidationResponse{Data: ""}, err
	}
	if resp.IsError() {
		return ApiDefinitionValidationResponse{Data: ""}, apiErrorFromResponse(resp)
	}
	return ApiDefinitionValidationResponse{Data: string(resp.Body())}, nil
}
