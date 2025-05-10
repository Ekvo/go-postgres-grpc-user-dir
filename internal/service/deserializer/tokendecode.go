package deserializer

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/Ekvo/go-postgres-grpc-user-dir/pkg/utils"
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
		return utils.ErrUtilsAutoHeaderMissing
	}
	token, err := utils.ParseAuthorization(authHeader)
	if err != nil {
		return err
	}
	td.token = token
	return nil
}
