package goreadme

import (
	"fmt"
	"strings"
)

// APIError models an error payload returned by the ReadMe API v2.
// Example:
//
//	{
//	  "type": "https://docs.readme.com/main/",
//	  "title": "We encountered validation errors while processing your input.",
//	  "status": 422,
//	  "detail": "The JSON you sent was the right format, but had data our endpoint couldn't process. Check the rest of the `errors` object for more details on what went wrong.",
//	  "poem": ["Data astray,", "Schema's dismay,", "Validation failed,", "Errors unveiled."],
//	  "errors": [{ "key": "category.uri", "message": "The supplied category URI must be of a category within the guide section of your docs." }]
//	}
type APIError struct {
	Type   string          `json:"type,omitempty"`
	Title  string          `json:"title,omitempty"`
	Status int             `json:"status,omitempty"`
	Detail string          `json:"detail,omitempty"`
	Errors []APIFieldError `json:"errors,omitempty"`
}

// APIFieldError represents a single field validation error entry.
type APIFieldError struct {
	Key     string `json:"key,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e APIError) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("%d %s: %s", e.Status, e.Title, e.Detail)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %s: %s", e.Status, e.Title, e.Detail)
	for _, fe := range e.Errors {
		if fe.Key != "" || fe.Message != "" {
			fmt.Fprintf(&b, " [%s: %s]", fe.Key, fe.Message)
		}
	}
	return b.String()
}
