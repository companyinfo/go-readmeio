package readme

// Guide models a guides page returned by ReadMe API v2.
//
// See: https://docs.readme.com/main/reference/getguide
type Guide struct {
	Slug          string           `json:"slug"`
	Title         string           `json:"title"`
	Type          string           `json:"type,omitempty"`  // e.g. "basic", "link"
	State         string           `json:"state,omitempty"` // "current" | "deprecated"
	Position      int              `json:"position,omitempty"`
	Content       GuideContent     `json:"content,omitempty"`
	Category      ResourceRef      `json:"category,omitempty"` // { uri: "..." }
	Parent        *ResourceRef     `json:"parent,omitempty"`   // { uri: "..." } or null
	Privacy       GuidePrivacy     `json:"privacy,omitempty"`
	Appearance    GuideAppearance  `json:"appearance,omitempty"`
	Metadata      GuideMetadata    `json:"metadata,omitempty"`
	AllowCrawlers AllowCrawlers    `json:"allow_crawlers,omitempty"`
	HRef          *GuideHRef       `json:"href,omitempty"`
	Links         *GuideLinks      `json:"links,omitempty"`
	Project       *GuideProject    `json:"project,omitempty"`
	Renderable    *GuideRenderable `json:"renderable,omitempty"`
	URI           string           `json:"uri,omitempty"`
	CreatedAt     string           `json:"created_at,omitempty"`
	UpdatedAt     string           `json:"updated_at,omitempty"`
}

// GuideContent models the body/excerpt/next block of a guide.
type GuideContent struct {
	Body    string     `json:"body,omitempty"` // Markdown
	Excerpt string     `json:"excerpt,omitempty"`
	Link    *GuideLink `json:"link,omitempty"`
	Next    *GuideNext `json:"next,omitempty"` // "What's next" links
}

// GuideNext models the `content.next` block (the "What's next" footer).
type GuideNext struct {
	Description string          `json:"description,omitempty"`
	Pages       []GuideNextPage `json:"pages,omitempty"`
}

// GuideNextPage models an entry in `content.next.pages`.
type GuideNextPage struct {
	Type  string `json:"type,omitempty"`  // "doc" | "link" | ...
	Slug  string `json:"slug,omitempty"`  // for internal "doc" links
	Title string `json:"title,omitempty"` // display label
	URL   string `json:"url,omitempty"`   // for external "link" entries
}

// ResourceRef is a generic { uri } reference to another v2 resource.
type ResourceRef struct {
	URI string `json:"uri"`
}

// GuideHRef models the set of helpful links for a guide.
type GuideHRef struct {
	Dash      string `json:"dash,omitempty"`
	Hub       string `json:"hub,omitempty"`
	GithubURL string `json:"github_url,omitempty"`
}

// GuideLinks models the `links` block on a guide.
type GuideLinks struct {
	Project string `json:"project,omitempty"`
}

// GuideProject models the `project` block on a guide.
type GuideProject struct {
	Name      string `json:"name,omitempty"`
	Subdomain string `json:"subdomain,omitempty"`
	URI       string `json:"uri,omitempty"`
}

// GuideRenderable models the `renderable` block status.
type GuideRenderable struct {
	Status  bool   `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// GuideLink models the `content.link` block for external links.
type GuideLink struct {
	URL    string `json:"url,omitempty"`
	NewTab bool   `json:"new_tab,omitempty"`
}

// GuidePrivacy models a guide's privacy block.
type GuidePrivacy struct {
	View string `json:"view,omitempty"` // "public" | "anyone_with_link"
}

// GuideAppearance models a guide's appearance block.
type AppearanceType string

const (
	// AppearanceTypeIcon represents a named icon from the theme set.
	AppearanceTypeIcon AppearanceType = "icon"
	// AppearanceTypeEmoji represents a Unicode emoji name/value.
	AppearanceTypeEmoji AppearanceType = "emoji"
)

// GuideAppearance models a guide's appearance block.
type GuideAppearance struct {
	Icon struct {
		Name string         `json:"name,omitempty"`
		Type AppearanceType `json:"type,omitempty"`
	} `json:"icon,omitempty"`
}

// GuideMetadata models the SEO metadata block (replaces v1 flat SEO fields).
type GuideMetadata struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Image       GuideImage `json:"image,omitempty"`
	Keywords    string     `json:"keywords,omitempty"`
}

type GuideImage struct {
	URI string `json:"uri,omitempty"`
}

// AllowCrawlers restricts allowed values for allow_crawlers in GuideParams.
type AllowCrawlers string

const (
	AllowCrawlersEnabled  AllowCrawlers = "enabled"
	AllowCrawlersDisabled AllowCrawlers = "disabled"
)

// GuideParams is the request body for create/update.
// Update sends only the fields explicitly supplied (omitempty).
type GuideParams struct {
	Slug          string           `json:"slug,omitempty"`
	Title         string           `json:"title,omitempty"`
	Type          string           `json:"type,omitempty"`
	State         string           `json:"state,omitempty"`
	Position      *int             `json:"position,omitempty"`
	Content       *GuideContent    `json:"content,omitempty"`
	Category      *ResourceRef     `json:"category,omitempty"`
	Parent        *ResourceRef     `json:"parent,omitempty"`
	Privacy       *GuidePrivacy    `json:"privacy,omitempty"`
	Appearance    *GuideAppearance `json:"appearance,omitempty"`
	Metadata      *GuideMetadata   `json:"metadata,omitempty"`
	AllowCrawlers *AllowCrawlers   `json:"allow_crawlers,omitempty"` // "enabled" | "disabled"
}
