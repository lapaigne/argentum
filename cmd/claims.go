package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func accTime(now time.Time) time.Time {
	return now.Add(time.Minute * 15)
}

func refTime(now time.Time) time.Time {
	return now.Add(time.Hour * 24)
}

func getClaims(c echo.Context) *jwtClaims {

	user := c.Get("user")
	if user == nil {
		return nil
	}

	t, ok := user.(*jwt.Token)
	if !ok {
		return nil
	}

	claims, ok := t.Claims.(*jwtClaims)
	if !ok {
		return nil
	}
	return claims
}

func signToken(uid, level int, t time.Time, secret string) (string, *jwtClaims, error) {

	claims := &jwtClaims{
		UID:   uid,
		Level: level,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))

	return signed, claims, err
}
