package deserializer

import "errors"

var (
	ErrDeserializerEmpty = errors.New("empty")

	ErrDeserializerInvalid = errors.New("invalid")

	ErrDeserializerMetadataEmpty = errors.New("missing metadata")
)
