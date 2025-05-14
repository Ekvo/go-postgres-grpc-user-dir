package deserializer

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

var (
	ErrDeserializerAutoHeaderMissing = errors.New("missing authorization header")

	ErrDeserializerTokenMissing = errors.New("missing token")

	ErrDeserializerTokenInvalid = errors.New("invalid token")
)

type TokenDecode struct {
	TokenHeader []string

	token string
}

func NewTokenDecode() *TokenDecode {
	return &TokenDecode{}
}

func (td *TokenDecode) Token() string {
	return td.token
}

func (td *TokenDecode) Decode(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ErrDeserializerMetadataEmpty
	}
	authHeader, ex := md["authorization"]
	if !ex {
		return ErrDeserializerAutoHeaderMissing
	}
	token, err := parseAuthorization(authHeader)
	if err != nil {
		return err
	}
	td.token = token
	return nil
}

func parseAuthorization(authHeader []string) (string, error) {
	if len(authHeader) == 0 {
		return "", ErrDeserializerAutoHeaderMissing
	}
	token := strings.TrimSpace(authHeader[0])
	if token == "" {
		return "", ErrDeserializerTokenMissing
	}
	tokenSplit := strings.Split(token, " ")
	if len(tokenSplit) != 2 {
		return "", ErrDeserializerTokenInvalid
	}
	if strings.ToLower(tokenSplit[0]) != "bearer" {
		return "", ErrDeserializerTokenInvalid
	}
	return tokenSplit[1], nil
}
