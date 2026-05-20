package readme

import (
	"context"
)

// guideEnvelope wraps single-item responses from the Guides API.
type guideEnvelope struct {
	Data Guide `json:"data"`
}

// GuideService defines Guide CRUD operations against the ReadMe API v2.
//
// All endpoints are scoped to a branch (a.k.a. version) supplied via the
// `branch` path parameter. A single guide page is identified by its `slug`.
//
// See:
//   - https://docs.readme.com/main/reference/createguide
//   - https://docs.readme.com/main/reference/getguide
//   - https://docs.readme.com/main/reference/updateguide
//   - https://docs.readme.com/main/reference/deleteguide
type GuideService interface {
	// Create creates a new guide on the given branch.
	Create(ctx context.Context, branch string, params GuideParams) (*Guide, error)
	// Get retrieves a single guide by its slug on the given branch.
	Get(ctx context.Context, branch, slug string) (*Guide, error)
	// Update updates an existing guide identified by its slug. Only the fields
	// set in params are sent.
	Update(ctx context.Context, branch, slug string, params GuideParams) (*Guide, error)
	// Delete removes a guide identified by its slug.
	Delete(ctx context.Context, branch, slug string) error
}

// GuideClient implements GuideService.
type GuideClient struct {
	client *Client
}

// NewGuideClient returns a new GuideClient bound to the provided root client.
func NewGuideClient(c *Client) *GuideClient {
	return &GuideClient{client: c}
}

// Create creates a new guide on the given branch.
// ReadMe API v2: POST /branches/{branch}/guides
func (g *GuideClient) Create(ctx context.Context, branch string, params GuideParams) (*Guide, error) {
	if err := check().
		Branch(branch).
		Title(params.Title).
		Category(params.Category).
		Err(); err != nil {
		return nil, err
	}

	var out guideEnvelope
	resp, err := g.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
		}).
		SetBody(params).
		SetResult(&out).
		SetError(&APIError{}).
		Post("/branches/{branch}/guides")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, apiErrorFromResponse(resp)
	}
	return &out.Data, nil
}

// Get retrieves a single guide by its slug.
// ReadMe API v2: GET /branches/{branch}/guides/{slug}
func (g *GuideClient) Get(ctx context.Context, branch, slug string) (*Guide, error) {
	if err := check().Branch(branch).Slug(slug).Err(); err != nil {
		return nil, err
	}

	var out guideEnvelope
	resp, err := g.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
			"slug":   slug,
		}).
		SetResult(&out).
		SetError(&APIError{}).
		Get("/branches/{branch}/guides/{slug}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, apiErrorFromResponse(resp)
	}
	return &out.Data, nil
}

// Update updates an existing guide identified by its slug.
// ReadMe API v2: PATCH /branches/{branch}/guides/{slug}
func (g *GuideClient) Update(ctx context.Context, branch, slug string, params GuideParams) (*Guide, error) {
	if err := check().Branch(branch).Slug(slug).Err(); err != nil {
		return nil, err
	}

	var out guideEnvelope
	resp, err := g.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
			"slug":   slug,
		}).
		SetBody(params).
		SetResult(&out).
		SetError(&APIError{}).
		Patch("/branches/{branch}/guides/{slug}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, apiErrorFromResponse(resp)
	}
	return &out.Data, nil
}

// Delete removes a guide identified by its slug.
// ReadMe API v2: DELETE /branches/{branch}/guides/{slug}
func (g *GuideClient) Delete(ctx context.Context, branch, slug string) error {
	if err := check().Branch(branch).Slug(slug).Err(); err != nil {
		return err
	}

	resp, err := g.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
			"slug":   slug,
		}).
		SetError(&APIError{}).
		Delete("/branches/{branch}/guides/{slug}")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return apiErrorFromResponse(resp)
	}
	return nil
}
