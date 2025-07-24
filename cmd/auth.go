package main

import (
	"argentum/db"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	refsecret = "ref-secret"
)

type jwtClaims struct {
	UID   int `json:"uid"`
	Level int `json:"level"`
	jwt.RegisteredClaims
}

func authTime() time.Time {
	return time.Now().Add(time.Minute * 1)

}

func refTime() time.Time {
	return time.Now().Add(time.Minute * 5)
}

func Login(c echo.Context) error {

	tel := c.FormValue("tel")
	pass := c.FormValue("pass")
	pwd := []byte(pass)

	u, err := db.GetUser(tel)
	if err != nil {
		fmt.Println("get user err")
		return echo.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Hash), pwd); err != nil {
		fmt.Println("compare hash")
		return echo.ErrUnauthorized
	}

	claims := &jwtClaims{
		u.Id,
		u.Level,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(authTime()),
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
		Expires:  authTime(),
	}

	c.SetCookie(cookie)

	refClaims := &jwtClaims{
		u.Id,
		u.Level,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refTime()),
		},
	}

	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims)

	s, err := refToken.SignedString([]byte(refsecret))
	if err != nil {
		return err
	}

	if err := db.UpdateToken(s, u.Id); err != nil {
		fmt.Println(err)
		return err
	}

	refCookie := &http.Cookie{
		Name:     "ref",
		Value:    s,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  refTime(),
	}

	c.SetCookie(refCookie)

	switch claims.Level {
	case 10, 50, 100:
		return c.Redirect(http.StatusSeeOther, "/menu")
	default:
		fmt.Println("default err")
		return echo.ErrUnauthorized
	}

}

func Logout(c echo.Context) error {

	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/")
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

func GetClaims(c echo.Context) *jwtClaims {

	user := c.Get("user").(*jwt.Token)
	return user.Claims.(*jwtClaims)
}

func AutoRefreshJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		path := c.Request().URL.Path
		if public[path] {
			return next(c)
		}

		cookie, err := c.Cookie("token")
		if err == nil {
			acc := cookie.Value
			_, err := jwt.ParseWithClaims(acc, &jwtClaims{}, func(t *jwt.Token) (any, error) {
				return []byte("secret"), nil
			})

			if err == nil {
				return next(c)
			}
		}

		fmt.Println(err)

		refCookie, err := c.Cookie("ref")
		if err != nil {
			fmt.Println("no ref cookie")
			fmt.Println(err)
			return echo.ErrUnauthorized
		}

		refToken, err := jwt.ParseWithClaims(refCookie.Value, &jwtClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(refsecret), nil
		})

		if err != nil || !refToken.Valid {
			return echo.ErrUnauthorized
		}

		claims, ok := refToken.Claims.(*jwtClaims)
		if !ok {
			return echo.ErrUnauthorized
		}

		newClaims := &jwtClaims{
			UID:   claims.UID,
			Level: claims.Level,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(authTime()),
			},
		}

		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
		s, err := newToken.SignedString([]byte("secret"))
		if err != nil {
			return echo.ErrUnauthorized
		}

		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    s,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			Expires:  newClaims.ExpiresAt.Time,
		})

		c.Request().Header.Set("Authorization", "Bearer "+s)

		return next(c)
	}
}
