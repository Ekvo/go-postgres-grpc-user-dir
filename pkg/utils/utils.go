package utils

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Message map[string]any

func (msg Message) String() string {
	lineMsg := make([]string, 0, len(msg))
	for k, v := range msg {
		lineMsg = append(lineMsg, fmt.Sprintf(`{%s:%v}`, k, v))
	}
	sort.Strings(lineMsg)
	return strings.Join(lineMsg, ",")
}

func GenerateJWT(secret, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().UTC().Add(time.Hour * 24 * 7).Unix(),
	})
	return token.SignedString([]byte(secret))
}

var ErrUtilsJWTTokenInvalid = errors.New("invalid token")

func VerifyJWT(secret, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrHashUnavailable
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", ErrUtilsJWTTokenInvalid
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}
	exploration, ok := claims["exp"].(float64)
	if !ok {
		return "", jwt.ErrInvalidKey
	}
	if int64(exploration) < time.Now().UTC().Unix() {
		return "", jwt.ErrTokenExpired
	}
	return claims["user_id"].(string), nil
}

var (
	ErrUtilsAutoHeaderMissing = errors.New("missing authorization header")

	ErrUtilsTokenMissing = errors.New("missing token")

	ErrUtilsTokenInvalid = errors.New("invalid token")
)

func ParseAuthorization(authHeader []string) (string, error) {
	if len(authHeader) == 0 {
		return "", ErrUtilsAutoHeaderMissing
	}
	token := strings.TrimSpace(authHeader[0])
	if token == "" {
		return "", ErrUtilsTokenMissing
	}
	tokenSplit := strings.Split(token, " ")
	if len(tokenSplit) != 2 {
		return "", ErrUtilsTokenInvalid
	}
	if strings.ToLower(tokenSplit[0]) != "bearer" {
		return "", ErrUtilsTokenInvalid
	}
	return tokenSplit[1], nil
}
