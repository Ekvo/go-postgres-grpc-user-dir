package deserializer

import (
	"context"
	"strconv"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/lib/jwtsign"
)

type IDDecode struct {
	id uint64
}

func NewIDDecode() *IDDecode {
	return &IDDecode{}
}

func (id *IDDecode) UserID() uint {
	return uint(id.id)
}

func (id *IDDecode) Decode(ctx context.Context) (err error) {
	content, ok := ctx.Value("content").(jwtsign.Content)
	if !ok {
		return ErrDeserializerInvalid
	}
	id.id, err = strconv.ParseUint(content["user_id"], 10, 64)
	if err != nil {
		return err
	}
	return nil
}
