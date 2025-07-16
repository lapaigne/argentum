package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthHandler(c echo.Context) error {

	tel := c.FormValue("tel")
	pass := c.FormValue("pass")

	if tel != "1234" || pass != "abcd" {
		return echo.ErrUnauthorized
	}

	claims := &jwtClaims{
		"admin",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
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

	return c.Redirect(http.StatusSeeOther, "/index")
}

func AuthTel(c echo.Context) error {
	tel := c.FormValue("tel")
	return c.Render(200, "tel-err", tel)
}

type jwtClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
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
