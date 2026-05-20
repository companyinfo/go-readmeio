package readme

import (
	"fmt"
	"strings"
)

// validator is a tiny fluent helper that accumulates the first validation
// error encountered across a chain of checks. It removes the repetitive
// `if err := validateX(...); err != nil { return ... }` boilerplate from
// the CRUD methods of every client.
//
// Usage:
//
//	if err := check().
//	    Branch(branch).
//	    Slug(slug).
//	    Err(); err != nil {
//	    return nil, err
//	}
//
// Each method short-circuits once the first error has been recorded, so
// later calls are cheap no-ops and the order of checks defines the order
// in which errors surface — typically: path params first, then body.
type validator struct {
	err error
}

// check starts a new validation chain.
func check() *validator { return &validator{} }

// Err returns the first error captured by the chain, or nil if all checks
// passed.
func (v *validator) Err() error { return v.err }

// fail records err if no previous error has been captured yet.
func (v *validator) fail(err error) *validator {
	if v.err == nil {
		v.err = err
	}
	return v
}

// NonEmpty records an error when value is empty or whitespace-only,
// using fieldName in the error message ("<fieldName> is required").
func (v *validator) NonEmpty(fieldName, value string) *validator {
	if v.err != nil {
		return v
	}
	if strings.TrimSpace(value) == "" {
		return v.fail(fmt.Errorf("%s is required", fieldName))
	}
	return v
}

// Branch validates the `branch` path parameter.
func (v *validator) Branch(branch string) *validator {
	return v.NonEmpty("branch", branch)
}

// Slug validates the `slug` path parameter.
func (v *validator) Slug(slug string) *validator {
	return v.NonEmpty("slug", slug)
}

// Title validates the `title` field (used both as a body field and as a
// path parameter for Categories).
func (v *validator) Title(title string) *validator {
	return v.NonEmpty("title", title)
}

// Category validates that a category ResourceRef is non-nil and has a
// non-empty URI. This matches the requirement for creating/updating Guides
// and References.
func (v *validator) Category(ref *ResourceRef) *validator {
	if v.err != nil {
		return v
	}
	if ref == nil || strings.TrimSpace(ref.URI) == "" {
		return v.fail(fmt.Errorf("category is required"))
	}
	return v
}

// Section validates the `section` path parameter (must be "guides" or
// "reference") and writes the canonical lowercase form back through the
// out pointer when provided.
func (v *validator) Section(section CategoryType, out *string) *validator {
	if v.err != nil {
		return v
	}
	s := strings.ToLower(strings.TrimSpace(string(section)))
	if s != "guides" && s != "reference" {
		return v.fail(fmt.Errorf("section must be 'guides' or 'reference'"))
	}
	if out != nil {
		*out = s
	}
	return v
}
