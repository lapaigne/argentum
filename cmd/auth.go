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

func isAdminErr(c echo.Context) error {

	acc := getClaims(c).Level
	if acc != ACC_ADMIN {
		return echo.ErrUnauthorized
	}

	return nil
}

func checkLevel(c echo.Context, min int) error {

	acc := getClaims(c).Level

	if acc < min {
		return echo.ErrUnauthorized
	}

	switch acc {
	case ACC_ADMIN, ACC_DISPATCHER, ACC_WORKER:
	default:
		return echo.ErrUnauthorized
	}

	return nil
}

func (e Endpoints) register(c echo.Context) error {
	tel := c.FormValue("tel")
	pass := c.FormValue("pass")
	pwd := []byte(pass)

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := db.User{
		Worker: 0,
		Login:  tel,
		Hash:   string(hash),
		Level:  0,
	}

	if _, err := db.AddUser(u); err != nil {
		return err
	}

	return nil
}

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
	case 10, 50, 100:
		return c.Redirect(303, "/menu/")
	default:
		return c.Redirect(303, "/signin")
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
