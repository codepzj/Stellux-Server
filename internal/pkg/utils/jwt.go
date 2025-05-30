package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var jwtKey = []byte(viper.GetString("JWT_SECRET"))

type JwtCustomClaims struct {
	ID string
	jwt.RegisteredClaims
}

func GenerateAccessToken(id string) (string, error) {
	now := time.Now().Local()
	claims := JwtCustomClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "stellux",
			Subject:   "stellux",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateRefreshToken(id string) (string, error) {
	now := time.Now().Local()
	claims := JwtCustomClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "stellux",
			Subject:   "stellux",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(jwtStr string) (*JwtCustomClaims, error) {
	claims := new(JwtCustomClaims)
	token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("无效的token")
	}
	return claims, nil
}
