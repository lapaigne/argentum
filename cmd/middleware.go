package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type jwtClaims struct {
	UID   int `json:"uid"`
	Level int `json:"level"`
	jwt.RegisteredClaims
}

func jwtMiddleware(config echojwt.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}

			var tokenStr string
			if cookie, err := c.Cookie("token"); err == nil {
				tokenStr = cookie.Value
			}

			if tokenStr != "" {
				tok, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (any, error) {
					return []byte(accSecret), nil
				})

				if err == nil && tok.Valid {
					c.Request().Header.Set("Authorization", "Bearer "+tokenStr)
					return echojwt.WithConfig(config)(next)(c)
				}

			}

			refCookie, err := c.Cookie("ref")
			if err != nil {
				fmt.Println(err)
				return echo.ErrUnauthorized
			}

			refTok, err := jwt.ParseWithClaims(refCookie.Value, &jwtClaims{}, func(t *jwt.Token) (any, error) {
				return []byte(refSecret), nil
			})

			if err != nil || !refTok.Valid {
				fmt.Println(err)
				return echo.ErrUnauthorized
			}

			refClaims, ok := refTok.Claims.(*jwtClaims)
			if !ok {
				return echo.ErrUnauthorized
			}

			aTime := accTime(time.Now())

			signed, claims, err := signToken(refClaims.UID, refClaims.Level, aTime, accSecret)
			if err != nil {
				fmt.Println(err)
				return echo.ErrUnauthorized
			}

			setCookie(c, "token", signed, aTime)
			c.Request().Header.Set("Authorization", "Bearer "+signed)

			c.Set("user", &jwt.Token{
				Claims: claims,
				Valid:  true,
			})

			return next(c)
		}
	}
}
