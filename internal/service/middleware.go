// contains middleware for request method checking and authorization if necessary
package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"google.golang.org/grpc"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/lib/jwtsign"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/deserializer"
)

var ErrServiceMethodInvalid = errors.New("invalid method")

// Authorization - middleware function
// check method
// 1. without auth -> next(ctx, req)
// 2. otherwise check the bearer token -> next(ctx, req)
func Authorization(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	next grpc.UnaryHandler) (resp any, err error) {
	log.Printf("service: request received for method - {%s};", info.FullMethod)
	method, err := methodSuffix(info.FullMethod)
	if err != nil {
		log.Printf("service: Authorization method error - {%v};", err)
		return nil, ErrServiceInternal
	}
	if !isAuth(method) {
		return next(ctx, req)
	}

	deserialize := deserializer.NewTokenDecode()
	if err := deserialize.Decode(ctx); err != nil {
		log.Printf("service: Authorization token error - {%v};", err)
		return nil, ErrServiceAuthorizationInvalid
	}

	content, err := jwtsign.GetContentFromToken(deserialize.Token())
	if err != nil {
		log.Printf("service: parse token error - {%v};", err)
		return nil, ErrServiceAuthorizationInvalid
	}
	ctx = context.WithValue(ctx, "content", content)

	return next(ctx, req)
}

// isAuth - return true if method with Authorization
func isAuth(method string) bool {
	if method == "UserData" || method == "UserUpdate" || method == "UserDelete" {
		return true
	}
	return false
}

// methodSuffix - parse Suffix from method
func methodSuffix(fullMethod string) (string, error) {
	parts := strings.Split(fullMethod, "/")
	if n := len(parts); n > 1 {
		return parts[n-1], nil
	}
	return "", ErrServiceMethodInvalid
}
