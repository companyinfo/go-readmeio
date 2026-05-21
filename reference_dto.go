package readme

// Reference models a reference (API reference) page returned by ReadMe API v2.
//
// The payload shares all sub-blocks with Guide (see guides_dto.go) and adds
// a Reference-specific `api` block carrying the HTTP method and path of the
// documented endpoint.
//
// See: https://docs.readme.com/main/reference/getreference
type Reference struct {
	Slug          string                `json:"slug"`
	Title         string                `json:"title"`
	Type          string                `json:"type,omitempty"`  // e.g. "basic", "link"
	State         string                `json:"state,omitempty"` // "current" | "deprecated"
	Position      int                   `json:"position,omitempty"`
	APIConfig     string                `json:"api_config,omitempty"`
	API           *ReferenceAPI         `json:"api,omitempty"`
	Content       GuideContent          `json:"content,omitempty"`
	Category      ResourceRef           `json:"category"`
	Parent        *ResourceRef          `json:"parent,omitempty"`
	Privacy       GuidePrivacy          `json:"privacy,omitempty"`
	Appearance    GuideAppearance       `json:"appearance,omitempty"`
	Metadata      ReferenceMetadata     `json:"metadata,omitempty"`
	Connections   *ReferenceConnections `json:"connections,omitempty"`
	AllowCrawlers AllowCrawlers         `json:"allow_crawlers,omitempty"`
	HRef          *GuideHRef            `json:"href,omitempty"`
	Links         *GuideLinks           `json:"links,omitempty"`
	Project       *GuideProject         `json:"project,omitempty"`
	Renderable    *GuideRenderable      `json:"renderable,omitempty"`
	URI           string                `json:"uri,omitempty"`
	CreatedAt     string                `json:"created_at,omitempty"`
	UpdatedAt     string                `json:"updated_at,omitempty"`
}

// ReferenceAPI models the `api` block specific to reference pages.
// In ReadMe API v2 the v1 flat `api.url` field was split into `api.method`
// and `api.path`.
type ReferenceAPI struct {
	Method     string                  `json:"method,omitempty"` // e.g. "GET", "POST"
	Path       string                  `json:"path,omitempty"`   // e.g. "/pets/{petId}"
	Stats      *ReferenceAPIStats      `json:"stats,omitempty"`
	Source     string                  `json:"source,omitempty"`
	URI        string                  `json:"uri,omitempty"`
	Validation *ReferenceAPIValidation `json:"validation,omitempty"`
}

// ReferenceAPIStats captures analysis flags for the underlying OpenAPI.
type ReferenceAPIStats struct {
	AdditionalProperties bool `json:"additional_properties,omitempty"`
	Callbacks            bool `json:"callbacks,omitempty"`
	CircularReferences   bool `json:"circular_references,omitempty"`
	CommonParameters     bool `json:"common_parameters,omitempty"`
	Discriminators       bool `json:"discriminators,omitempty"`
	Links                bool `json:"links,omitempty"`
	Polymorphism         bool `json:"polymorphism,omitempty"`
	References           bool `json:"references,omitempty"`
	ServerVariables      bool `json:"server_variables,omitempty"`
	Style                bool `json:"style,omitempty"`
	Webhooks             bool `json:"webhooks,omitempty"`
	XMLRequests          bool `json:"xml_requests,omitempty"`
	XMLResponses         bool `json:"xml_responses,omitempty"`
	XMLSchemas           bool `json:"xml_schemas,omitempty"`
}

// ReferenceAPIValidation represents validation status for the API definition.
type ReferenceAPIValidation struct {
	Status string `json:"status,omitempty"` // e.g. "pending", "valid", "invalid"
	Reason string `json:"reason,omitempty"`
}

// ReferenceConnections models auxiliary connections for a reference page.
type ReferenceConnections struct {
	Recipes []ReferenceRecipe `json:"recipes,omitempty"`
}

// ReferenceRecipe represents a connected recipe entry.
type ReferenceRecipe struct {
	URI        string                    `json:"uri,omitempty"`
	Title      string                    `json:"title,omitempty"`
	Slug       string                    `json:"slug,omitempty"`
	Appearance ReferenceRecipeAppearance `json:"appearance,omitempty"`
}

// ReferenceRecipeAppearance captures styling for a recipe.
type ReferenceRecipeAppearance struct {
	BackgroundColor string `json:"background_color,omitempty"`
	Emoji           string `json:"emoji,omitempty"`
}

// ReferenceMetadata extends GuideMetadata with additional fields present on Reference.
type ReferenceMetadata struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Image       struct {
		URI string `json:"uri,omitempty"`
		URL string `json:"url,omitempty"`
	} `json:"image,omitempty"`
	Keywords string `json:"keywords,omitempty"`
	XImport  string `json:"x_import,omitempty"`
}

// ReferenceParams is the request body for create/update.
// Update sends only the fields explicitly supplied (omitempty).
type ReferenceParams struct {
	Slug          string             `json:"slug,omitempty"` // optional override
	Title         string             `json:"title,omitempty"`
	Type          string             `json:"type,omitempty"`
	State         string             `json:"state,omitempty"`
	Position      *int               `json:"position,omitempty"`
	APIConfig     *string            `json:"api_config,omitempty"`
	API           *ReferenceAPI      `json:"api,omitempty"`
	Content       *GuideContent      `json:"content,omitempty"`
	Category      *ResourceRef       `json:"category,omitempty"`
	Parent        *ResourceRef       `json:"parent,omitempty"`
	Privacy       *GuidePrivacy      `json:"privacy,omitempty"`
	Appearance    *GuideAppearance   `json:"appearance,omitempty"`
	Metadata      *ReferenceMetadata `json:"metadata,omitempty"`
	AllowCrawlers *AllowCrawlers     `json:"allow_crawlers,omitempty"` // "enabled" | "disabled"
}
