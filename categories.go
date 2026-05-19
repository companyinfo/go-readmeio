package readme

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

// categoryEnvelope wraps responses from the Categories API.
// It supports both list responses (data: [ ... ]) and single-item responses (data: { ... })
// by keeping the raw JSON for post-processing.
type categoryEnvelope struct {
	Data json.RawMessage `json:"data"`
}

// CategoryService defines Category CRUD operations against the ReadMe API v2.
//
// All endpoints are scoped to a branch (a.k.a. version) supplied via the
// `branch` path parameter, and to a `section` ("guides" or "reference").
//
// See: https://docs.readme.com/main/reference/getcategories-1
type CategoryService interface {
	// Create creates a new category on the given branch. The category's section
	// ("guides" or "reference") is provided in params and sent in the request body.
	Create(ctx context.Context, branch string, params CategoryParams) (*Category, error)
	// List retrieves all categories on the given branch for the given section
	// ("guides" or "reference").
	List(ctx context.Context, branch, section string) ([]Category, error)
	// GetByTitle retrieves a single category by its title on the given branch
	// and section.
	GetByTitle(ctx context.Context, branch, section, title string) (*Category, error)
	// Update updates an existing category identified by its title on the given
	// branch and section. Only the fields set in params are sent.
	Update(ctx context.Context, branch, section, title string, params CategoryParams) (*Category, error)
	// Delete removes a category identified by its title on the given branch
	// and section.
	Delete(ctx context.Context, branch, section, title string) error
}

// CategoryClient implements CategoryService.
type CategoryClient struct {
	client *Client
}

// NewCategoryClient returns a new CategoryClient bound to the provided root client.
func NewCategoryClient(c *Client) *CategoryClient {
	return &CategoryClient{client: c}
}

func (c *CategoryClient) newRequest(ctx context.Context) *resty.Request {
	return c.client.HTTPClient.R().SetContext(ctx)
}

// Create creates a new category on the given branch.
// ReadMe API v2: POST /branches/{branch}/categories
func (c *CategoryClient) Create(ctx context.Context, branch string, params CategoryParams) (*Category, error) {
	if err := check().
		Branch(branch).
		Title(params.Title).
		Section(string(params.Section), nil).
		Err(); err != nil {
		return nil, err
	}

	var env categoryEnvelope
	resp, err := c.newRequest(ctx).
		SetPathParams(map[string]string{
			"branch": branch,
		}).
		SetBody(params).
		SetResult(&env).
		SetError(&APIError{}).
		Post("/branches/{branch}/categories")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, resp.Error().(*APIError)
	}

	var cat Category
	if err := json.Unmarshal(env.Data, &cat); err != nil {
		return nil, err
	}
	return &cat, nil
}

// List retrieves all categories on the given branch for the given section.
// ReadMe API v2: GET /branches/{branch}/categories/{section}
func (c *CategoryClient) List(ctx context.Context, branch, section string) ([]Category, error) {
	if err := check().
		Branch(branch).
		Section(section, &section).
		Err(); err != nil {
		return nil, err
	}

	var env categoryEnvelope
	resp, err := c.newRequest(ctx).
		SetPathParams(map[string]string{
			"branch":  branch,
			"section": section,
		}).
		SetResult(&env).
		SetError(&APIError{}).
		Get("/branches/{branch}/categories/{section}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, resp.Error().(*APIError)
	}

	var list []Category
	if err := json.Unmarshal(env.Data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// GetByTitle retrieves a single category by its title.
// ReadMe API v2: GET /branches/{branch}/categories/{section}/{title}
func (c *CategoryClient) GetByTitle(ctx context.Context, branch, section, title string) (*Category, error) {
	if err := check().
		Branch(branch).
		Section(section, &section).
		Title(title).
		Err(); err != nil {
		return nil, err
	}

	var env categoryEnvelope
	resp, err := c.newRequest(ctx).
		SetPathParams(map[string]string{
			"branch":  branch,
			"section": section,
			"title":   title,
		}).
		SetResult(&env).
		SetError(&APIError{}).
		Get("/branches/{branch}/categories/{section}/{title}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, resp.Error().(*APIError)
	}

	var cat Category
	if err := json.Unmarshal(env.Data, &cat); err != nil {
		return nil, err
	}
	return &cat, nil
}

// Update updates an existing category identified by its title.
// ReadMe API v2: PATCH /branches/{branch}/categories/{section}/{title}
func (c *CategoryClient) Update(ctx context.Context, branch, section, title string, params CategoryParams) (*Category, error) {
	if err := check().
		Branch(branch).
		Section(section, &section).
		Title(title).
		Err(); err != nil {
		return nil, err
	}

	var env categoryEnvelope
	resp, err := c.newRequest(ctx).
		SetPathParams(map[string]string{
			"branch":  branch,
			"section": section,
			"title":   title,
		}).
		SetBody(params).
		SetResult(&env).
		SetError(&APIError{}).
		Patch("/branches/{branch}/categories/{section}/{title}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, resp.Error().(*APIError)
	}

	var cat Category
	if err := json.Unmarshal(env.Data, &cat); err != nil {
		return nil, err
	}
	return &cat, nil
}

// Delete removes a category identified by its title.
// ReadMe API v2: DELETE /branches/{branch}/categories/{section}/{title}
func (c *CategoryClient) Delete(ctx context.Context, branch, section, title string) error {
	if err := check().
		Branch(branch).
		Section(section, &section).
		Title(title).
		Err(); err != nil {
		return err
	}

	resp, err := c.newRequest(ctx).
		SetPathParams(map[string]string{
			"branch":  branch,
			"section": section,
			"title":   title,
		}).
		SetError(&APIError{}).
		Delete("/branches/{branch}/categories/{section}/{title}")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return resp.Error().(*APIError)
	}
	return nil
}
