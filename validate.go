package readme

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// validate is the package-wide go-playground/validator instance used to run
// struct-tag-based validation on request DTOs and on the small path-parameter
// DTOs declared in the `*_dto.go` files.
//
// It is configured with a tag-name function that reports field names using
// their `json` tag (falling back to the Go field name when no json tag is
// set), so validation error messages refer to fields by their wire name
// (e.g. "title", "category.uri", "branch", "section") rather than their
// Go identifier.
var validate = func() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return fld.Name
		}
		return name
	})
	return v
}()

// validateParams runs struct validation on a single value and returns the
// first violation encountered, formatted as a concise message.
func validateParams(v any) error {
	if v == nil {
		return nil
	}
	if err := validate.Struct(v); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// formatValidationError converts the first validator.FieldError into a
// concise message of the form "<json field path> is required" (or, for
// non-required tags, a tag-specific message such as
// "<json field path> must be one of [guides reference]").
//
// We deliberately surface only the first error to keep messages short and
// stable for callers that match on substrings.
func formatValidationError(err error) error {
	var fes validator.ValidationErrors
	if !errors.As(err, &fes) || len(fes) == 0 {
		return err
	}
	fe := fes[0]
	// Drop the top-level struct name from the namespace (e.g.
	// "GuideParams.category.uri" -> "category.uri") so messages stay close
	// to the JSON wire format.
	ns := fe.Namespace()
	if i := strings.IndexByte(ns, '.'); i >= 0 {
		ns = ns[i+1:]
	}
	switch fe.Tag() {
	case "required":
		return fmt.Errorf("%s is required", ns)
	case "oneof":
		return fmt.Errorf("%s must be one of [%s]", ns, fe.Param())
	default:
		if fe.Param() != "" {
			return fmt.Errorf("%s must satisfy %s=%s", ns, fe.Tag(), fe.Param())
		}
		return fmt.Errorf("%s must satisfy %s", ns, fe.Tag())
	}
}
