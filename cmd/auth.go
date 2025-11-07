package main

import (
	"argentum/db"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	ACC_ADMIN      = 100
	ACC_DISPATCHER = 50
	ACC_WORKER     = 10
	ACC_NONE       = 0
)

func (e Endpoints) signin_POST(c echo.Context) error {
	tel := c.FormValue("tel")
	pass := c.FormValue("pass")

	u, err := db.GetUser(tel)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(pass)) != nil {
		return c.Redirect(303, "/signin")
	}

	now := time.Now()

	aTime := accTime(now)
	rTime := refTime(now)

	accSigned, accClaims, err := signToken(u.Worker, u.Level, aTime, accSecret)
	if err != nil {
		return c.Redirect(303, "/signin")
	}

	refSigned, _, err := signToken(u.Worker, u.Level, rTime, refSecret)
	if err != nil {
		return c.Redirect(303, "/signin")
	}

	setCookie(c, "token", accSigned, aTime)
	setCookie(c, "ref", refSigned, rTime)

	switch accClaims.Level {
	case ACC_WORKER:
		return c.Redirect(303, "/me/")
	case ACC_DISPATCHER, ACC_ADMIN:
		return c.Redirect(303, "/menu/")
	default:
		return e.logout(c)
	}
}

func (e Endpoints) logout(c echo.Context) error {
	clearCookie(c, "token")
	clearCookie(c, "ref")

	return c.Redirect(303, "/signin")
}
func (e Endpoints) signin_GET(c echo.Context) error {
	return c.Render(200, "signin", nil)
}

func (e Endpoints) authTel(c echo.Context) error {
	tel := c.FormValue("tel")
	return c.Render(200, "tel-err", tel)
}

func hash(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}
