package main

import (
	"argentum/db"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var AUTHTIME = time.Now().Add(time.Minute * 5)

type jwtClaims struct {
	Level int `json:"level"`
	jwt.RegisteredClaims
}

func AuthHandler(c echo.Context) error {

	tel := c.FormValue("tel")
	pass := c.FormValue("pass")

	pwd := []byte(pass)

	u, err := db.GetUser(tel)

	if err != nil {
		return echo.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Hash), pwd); err != nil {
		return echo.ErrUnauthorized
	}

	claims := &jwtClaims{
		u.Level,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(AUTHTIME),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    t,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/menu")
}

func AuthTel(c echo.Context) error {
	tel := c.FormValue("tel")
	return c.Render(200, "tel-err", tel)
}

func JWTCooked() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err == nil {
				c.Request().Header.Set("Authorization", "Bearer "+cookie.Value)
			}
			return next(c)
		}
	}
}

func JWTRoles(minLevel int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*jwtClaims)
			if claims.Level >= minLevel {
				return next(c)
			}
			return echo.ErrForbidden
		}
	}
}
