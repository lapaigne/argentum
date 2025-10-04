package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func setCookie(c echo.Context, name, value string, time time.Time) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time,
	}
	c.SetCookie(cookie)
}

func clearCookie(c echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}

	c.SetCookie(cookie)
}
