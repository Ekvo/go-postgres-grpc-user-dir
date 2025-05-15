// describes the rule for obtaining data from a request and
// creating objects from the "model" package for use at the business layer
package deserializer

import (
	"errors"
	"regexp"
)

// mark Errors to create a detailed problem object for the error response
var (
	ErrDeserializerEmpty = errors.New("empty")

	ErrDeserializerInvalid = errors.New("invalid")

	ErrDeserializerMetadataEmpty = errors.New("missing metadata")
)

// reEmail - regexp for check email from request
var reEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
