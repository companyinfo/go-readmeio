package readme

import (
	"context"
)

// referenceResponse wraps single-item responses from the Reference API.
type referenceResponse struct {
	Data Reference `json:"data"`
}

// ReferenceService defines Reference CRUD operations against the ReadMe API v2.
//
// All endpoints are scoped to a branch (a.k.a. version) supplied via the
// `branch` path parameter. A single reference page is identified by its `slug`.
//
// See:
//   - https://docs.readme.com/main/reference/createreference
//   - https://docs.readme.com/main/reference/getreference
//   - https://docs.readme.com/main/reference/updatereference
//   - https://docs.readme.com/main/reference/deletereference
type ReferenceService interface {
	// Create creates a new reference page on the given branch.
	Create(ctx context.Context, branch string, params ReferenceParams) (*Reference, error)
	// Get retrieves a single reference page by its slug on the given branch.
	Get(ctx context.Context, branch, slug string) (*Reference, error)
	// Update updates an existing reference page identified by its slug. Only
	// the fields set in params are sent.
	Update(ctx context.Context, branch, slug string, params ReferenceParams) (*Reference, error)
	// Delete removes a reference page identified by its slug.
	Delete(ctx context.Context, branch, slug string) error
}

// ReferenceClient implements ReferenceService.
type ReferenceClient struct {
	client *Client
}

// NewReferenceClient returns a new ReferenceClient bound to the provided root client.
func NewReferenceClient(c *Client) *ReferenceClient {
	return &ReferenceClient{client: c}
}

// Create creates a new reference page on the given branch.
// ReadMe API v2: POST /branches/{branch}/reference
func (r *ReferenceClient) Create(ctx context.Context, branch string, params ReferenceParams) (*Reference, error) {
	if err := validateBranch(branch); err != nil {
		return nil, err
	}
	if err := validateParams(params); err != nil {
		return nil, err
	}

	var out referenceResponse
	resp, err := r.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
		}).
		SetBody(params).
		SetResult(&out).
		SetError(&APIError{}).
		Post("/branches/{branch}/reference")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, apiErrorFromResponse(resp)
	}
	return &out.Data, nil
}

// Get retrieves a single reference page by its slug.
// ReadMe API v2: GET /branches/{branch}/reference/{slug}
func (r *ReferenceClient) Get(ctx context.Context, branch, slug string) (*Reference, error) {
	if err := validateBranch(branch); err != nil {
		return nil, err
	}
	if err := validateSlug(slug); err != nil {
		return nil, err
	}

	var out referenceResponse
	resp, err := r.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
			"slug":   slug,
		}).
		SetResult(&out).
		SetError(&APIError{}).
		Get("/branches/{branch}/reference/{slug}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, apiErrorFromResponse(resp)
	}
	return &out.Data, nil
}

// Update updates an existing reference page identified by its slug.
// ReadMe API v2: PATCH /branches/{branch}/reference/{slug}
func (r *ReferenceClient) Update(ctx context.Context, branch, slug string, params ReferenceParams) (*Reference, error) {
	if err := validateBranch(branch); err != nil {
		return nil, err
	}
	if err := validateSlug(slug); err != nil {
		return nil, err
	}

	var out referenceResponse
	resp, err := r.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
			"slug":   slug,
		}).
		SetBody(params).
		SetResult(&out).
		SetError(&APIError{}).
		Patch("/branches/{branch}/reference/{slug}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, apiErrorFromResponse(resp)
	}
	return &out.Data, nil
}

// Delete removes a reference page identified by its slug.
// ReadMe API v2: DELETE /branches/{branch}/reference/{slug}
func (r *ReferenceClient) Delete(ctx context.Context, branch, slug string) error {
	if err := validateBranch(branch); err != nil {
		return err
	}
	if err := validateSlug(slug); err != nil {
		return err
	}

	resp, err := r.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
			"slug":   slug,
		}).
		SetError(&APIError{}).
		Delete("/branches/{branch}/reference/{slug}")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return apiErrorFromResponse(resp)
	}
	return nil
}
