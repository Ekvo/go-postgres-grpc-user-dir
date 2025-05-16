// contains a global non-exported variable 'secretKey' for working with jwt
// sets secretKey during application startup from config.Config
// contains functions for generating and parsing jwt.Token
package jwtsign

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
)

var (
	// ErrJWTSecretKeySetAfterRun - attempt to change non-empty secretkey
	ErrJWTSecretKeySetAfterRun = errors.New("secret key already exists")

	ErrJWTSecretKeyEmpty = errors.New("secret key not found")

	// ErrJWTContentInvalid - len('Content'(map[string]string)==0)
	ErrJWTContentInvalid = errors.New("invalid content")
)

// SecretKey -  key for jwt.Token
var secretKey = ""

// NewSecretKey - call in Run -> during application startup
func NewSecretKey(cfg *config.Config) error {
	if secretKey == "" {
		secretKey = cfg.JWTSecretKey
		return nil
	}
	return ErrJWTSecretKeySetAfterRun
}

// time exploration for jwt.Token see 'TokenGenerator'
const tokenLife = 7 * 24 * time.Hour

type Content map[string]string

// TokenGenerator - create jwt token using specific key
// set time of exploration in claims
func TokenGenerator(content Content) (string, error) {
	if secretKey == "" {
		return "", ErrJWTSecretKeyEmpty
	}
	if len(content) == 0 {
		return "", ErrJWTContentInvalid
	}
	claims := jwt.MapClaims{}
	for key, val := range content {
		claims[key] = val
	}
	claims["exp"] = time.Now().UTC().Add(tokenLife).Unix()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(secretKey))
}

// GetContentFromToken - get all fields without "exploration" from token
func GetContentFromToken(token string) (Content, error) {
	jwtToken, err := tokenRetrive(token)
	if err != nil {
		return nil, err
	}
	return receiveContentFromToken(jwtToken)
}

// tokenRetrive - get jwt.Token from string
func tokenRetrive(value string) (*jwt.Token, error) {
	return jwt.Parse(value, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrHashUnavailable
		}
		if secretKey == "" {
			return nil, ErrJWTSecretKeyEmpty
		}
		return []byte(secretKey), nil
	})
}

// receiveContentFromToken - check token expiration date and get date from jwt.MapClaims
func receiveContentFromToken(token *jwt.Token) (Content, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	exploration, ok := claims["exp"].(float64)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}
	if int64(exploration) < time.Now().UTC().Unix() {
		return nil, jwt.ErrTokenExpired
	}
	delete(claims, "exp")
	contetn := Content{}
	for key, val := range claims {
		line, ok := val.(string)
		if !ok {
			return nil, ErrJWTContentInvalid
		}
		contetn[key] = line
	}
	if len(contetn) == 0 {
		return nil, ErrJWTContentInvalid
	}
	return contetn, nil
}
