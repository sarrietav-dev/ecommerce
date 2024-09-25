package main

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.Email,
		Issuer:    "auth.service",
		Audience:  jwt.ClaimStrings{"user"},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		ID:        strconv.FormatUint(uint64(user.ID), 10),
	})
	return token.SignedString([]byte("secret"))
}
