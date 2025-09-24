package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	refSecret = "ref-secret"
	accSecret = "secret"
)

type jwtClaims struct {
	UID   int `json:"uid"`
	Level int `json:"level"`
	jwt.RegisteredClaims
}

// func JWTCooked(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		cookie, err := c.Cookie("token")
// 		if err == nil {
// 			c.Request().Header.Set("Authorization", "Bearer "+cookie.Value)
// 		}
// 		return next(c)
// 	}
// }

// func JWTRoles(minLevel int) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			user := c.Get("user").(*jwt.Token)
// 			claims := user.Claims.(*jwtClaims)
// 			if claims.Level >= minLevel {
// 				return next(c)
// 			}
// 			return echo.ErrForbidden
// 		}
// 	}
// }

// func autoRefreshJWT(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		path := c.Request().URL.Path
// 		if public[path] {
// 			return next(c)
// 		}
// 		if strings.HasPrefix(path, "/css/") {
// 			return next(c)
// 		}
// 		cookie, err := c.Cookie("token")
// 		if err == nil {
// 			acc := cookie.Value
// 			_, err := jwt.ParseWithClaims(acc, &jwtClaims{}, func(t *jwt.Token) (any, error) {
// 				return []byte("secret"), nil
// 			})
// 			if err == nil {
// 				return next(c)
// 			}
// 		}
// 		fmt.Println(err)
// 		refCookie, err := c.Cookie("ref")
// 		if err != nil {
// 			return echo.ErrUnauthorized
// 		}
// 		refToken, err := jwt.ParseWithClaims(refCookie.Value, &jwtClaims{}, func(t *jwt.Token) (any, error) {
// 			return []byte(refSecret), nil
// 		})
// 		if err != nil || !refToken.Valid {
// 			return echo.ErrUnauthorized
// 		}
// 		claims, ok := refToken.Claims.(*jwtClaims)
// 		if !ok {
// 			return echo.ErrUnauthorized
// 		}
// 		newClaims := &jwtClaims{
// 			UID:   claims.UID,
// 			Level: claims.Level,
// 			RegisteredClaims: jwt.RegisteredClaims{
// 				ExpiresAt: jwt.NewNumericDate(accTime()),
// 			},
// 		}
// 		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
// 		s, err := newToken.SignedString([]byte("secret"))
// 		if err != nil {
// 			return echo.ErrUnauthorized
// 		}
// 		c.SetCookie(&http.Cookie{
// 			Name:     "token",
// 			Value:    s,
// 			Path:     "/",
// 			HttpOnly: true,
// 			Secure:   true,
// 			SameSite: http.SameSiteLaxMode,
// 			Expires:  newClaims.ExpiresAt.Time,
// 		})
// 		c.Request().Header.Set("Authorization", "Bearer "+s)
// 		return next(c)
// 	}
// }

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
				fmt.Println(err)
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

			signed, _, err := signToken(refClaims.UID, refClaims.Level, aTime, refSecret)
			if err != nil {
				fmt.Println(err)
				return echo.ErrUnauthorized
			}

			setCookie(c, "token", signed, aTime)
			c.Request().Header.Set("Authorization", "Bearer "+signed)
			return echojwt.WithConfig(config)(next)(c)
		}
	}
}
