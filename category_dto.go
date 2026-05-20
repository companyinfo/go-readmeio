package readme

// CategoryType identifies the section a category belongs to in ReadMe API v2.
//
// In v2 the section is supplied as a path parameter (or in the request body
// when creating a category). The accepted values are "guides" and "reference".
type CategoryType string

const (
	// CategoryTypeReference represents the API Reference section.
	CategoryTypeReference CategoryType = "reference"
	// CategoryTypeGuides represents the Guides (knowledge base) section.
	CategoryTypeGuides CategoryType = "guides"
)

// Category models a category returned by the ReadMe API v2.
//
// See: https://docs.readme.com/main/reference/getcategories-1
type Category struct {
	// Title is the display title of the category and also serves as its unique
	// resource identifier within a (branch, section) pair in API v2.
	Title string `json:"title"`
	// Section is "guides" or "reference".
	Section CategoryType `json:"section,omitempty"`
	// Links contain the `project` URI for the category.
	Links CategoryLinks `json:"links"`
	// URI to this category.
	URI string `json:"uri"`
}

// CategoryLinks models the `links` block on a category.
type CategoryLinks struct {
	Project string `json:"project"`
}

// CategoryParams contains the fields used when creating or updating a category
// via the ReadMe API v2.
//
// Note: in v2 the version/branch is supplied as a path parameter and is no
// longer part of the request body.
type CategoryParams struct {
	// Title is the title of the category.
	Title string `json:"title"`
	// Section is "guides" or "reference". Required on Create; optional on Update.
	Section CategoryType `json:"section,omitempty"`
}
