package readme

import (
	"context"
	"encoding/json"
	"errors"
)

// categoryResponse wraps responses from the Categories API.
// It supports both list responses (data: [ ... ]) and single-item responses (data: { ... })
// by keeping the raw JSON for post-processing.
type categoryResponse struct {
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
	List(ctx context.Context, branch string, section CategoryType) ([]Category, error)
	// GetByTitle retrieves a single category by its title on the given branch
	// and section. The title is a unique identifier within a (branch, section) pair.
	GetByTitle(ctx context.Context, branch string, section CategoryType, title string) (*Category, error)
	// Update updates an existing category identified by its title on the given
	// branch and section. Only the fields set in params are sent.
	Update(ctx context.Context, branch string, section CategoryType, title string, params CategoryParams) (*Category, error)
	// Delete removes a category identified by its title on the given branch
	// and section.
	Delete(ctx context.Context, branch string, section CategoryType, title string) error
}

// CategoryClient implements CategoryService.
type CategoryClient struct {
	client *Client
}

// NewCategoryClient returns a new CategoryClient bound to the provided root client.
func NewCategoryClient(c *Client) *CategoryClient {
	return &CategoryClient{client: c}
}

// validateBranch ensures the branch path parameter is supplied.
func validateBranch(branch string) error {
	if branch == "" {
		return errors.New("branch is required")
	}
	return nil
}

// validateTitle ensures the category title path parameter is supplied.
func validateTitle(title string) error {
	if title == "" {
		return errors.New("title is required")
	}
	return nil
}

// Create creates a new category on the given branch.
// ReadMe API v2: POST /branches/{branch}/categories
func (c *CategoryClient) Create(ctx context.Context, branch string, params CategoryParams) (*Category, error) {
	if err := validateBranch(branch); err != nil {
		return nil, err
	}
	if err := validateParams(params); err != nil {
		return nil, err
	}

	var env categoryResponse
	resp, err := c.client.NewRequest(ctx).
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
		return nil, apiErrorFromResponse(resp)
	}

	var cat Category
	if err := json.Unmarshal(env.Data, &cat); err != nil {
		return nil, err
	}
	return &cat, nil
}

// List retrieves all categories on the given branch for the given section.
// ReadMe API v2: GET /branches/{branch}/categories/{section}
func (c *CategoryClient) List(ctx context.Context, branch string, section CategoryType) ([]Category, error) {
	if err := validateBranch(branch); err != nil {
		return nil, err
	}
	if err := validateSection(section); err != nil {
		return nil, err
	}
	sec := canonSection(section)

	var env categoryResponse
	resp, err := c.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch":  branch,
			"section": sec,
		}).
		SetResult(&env).
		SetError(&APIError{}).
		Get("/branches/{branch}/categories/{section}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, apiErrorFromResponse(resp)
	}

	var list []Category
	if err := json.Unmarshal(env.Data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// GetByTitle retrieves a single category by its title, which is a unique
// identifier within a (branch, section) pair.
// ReadMe API v2: GET /branches/{branch}/categories/{section}/{title}
func (c *CategoryClient) GetByTitle(ctx context.Context, branch string, section CategoryType, title string) (*Category, error) {
	if err := validateBranch(branch); err != nil {
		return nil, err
	}
	if err := validateSection(section); err != nil {
		return nil, err
	}
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	sec := canonSection(section)

	var env categoryResponse
	resp, err := c.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch":  branch,
			"section": sec,
			"title":   title,
		}).
		SetResult(&env).
		SetError(&APIError{}).
		Get("/branches/{branch}/categories/{section}/{title}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, apiErrorFromResponse(resp)
	}

	var cat Category
	if err := json.Unmarshal(env.Data, &cat); err != nil {
		return nil, err
	}
	return &cat, nil
}

// Update updates an existing category identified by its title.
// ReadMe API v2: PATCH /branches/{branch}/categories/{section}/{title}
func (c *CategoryClient) Update(ctx context.Context, branch string, section CategoryType, title string, params CategoryParams) (*Category, error) {
	if err := validateBranch(branch); err != nil {
		return nil, err
	}
	if err := validateSection(section); err != nil {
		return nil, err
	}
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	sec := canonSection(section)

	var env categoryResponse
	resp, err := c.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch":  branch,
			"section": sec,
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
		return nil, apiErrorFromResponse(resp)
	}

	var cat Category
	if err := json.Unmarshal(env.Data, &cat); err != nil {
		return nil, err
	}
	return &cat, nil
}

// Delete removes a category identified by its title.
// ReadMe API v2: DELETE /branches/{branch}/categories/{section}/{title}
func (c *CategoryClient) Delete(ctx context.Context, branch string, section CategoryType, title string) error {
	if err := validateBranch(branch); err != nil {
		return err
	}
	if err := validateSection(section); err != nil {
		return err
	}
	if err := validateTitle(title); err != nil {
		return err
	}
	sec := canonSection(section)

	resp, err := c.client.NewRequest(ctx).
		SetPathParams(map[string]string{
			"branch":  branch,
			"section": sec,
			"title":   title,
		}).
		SetError(&APIError{}).
		Delete("/branches/{branch}/categories/{section}/{title}")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return apiErrorFromResponse(resp)
	}
	return nil
}
