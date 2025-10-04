package main

import (
	"argentum/db"
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func isAdminErr(c echo.Context) error {

	acc := getClaims(c).Level
	if acc != db.AccessLevels["admin"] {
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
	case db.AccessLevels["worker"]:
	case db.AccessLevels["dispatcher"]:
	case db.AccessLevels["admin"]:
	default:
		return echo.ErrUnauthorized
	}

	return nil
}

func register(c echo.Context) error {

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
		Token:  sql.NullString{},
	}

	db.AddUser(u)

	return nil
}

func login(c echo.Context) error {

	tel := c.FormValue("tel")
	pass := c.FormValue("pass")

	u, err := db.GetUser(tel)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(pass)) != nil {
		return echo.ErrUnauthorized
	}

	now := time.Now()

	aTime := accTime(now)
	rTime := refTime(now)

	accSigned, accClaims, err := signToken(u.Worker, u.Level, aTime, accSecret)
	if err != nil {
		return echo.ErrUnauthorized
	}

	refSigned, _, err := signToken(u.Worker, u.Level, rTime, refSecret)
	if err != nil {
		return echo.ErrUnauthorized
	}

	setCookie(c, "token", accSigned, aTime)
	setCookie(c, "ref", refSigned, rTime)

	if err := db.UpdateToken(refSigned, u.Worker); err != nil {
		return echo.ErrUnauthorized
	}

	switch accClaims.Level {
	case 10, 50, 100:
		return c.Redirect(http.StatusSeeOther, "/menu/")
	default:
		return echo.ErrUnauthorized
	}
}

func logout(c echo.Context) error {

	clearCookie(c, "token")
	clearCookie(c, "ref")

	return c.Redirect(http.StatusSeeOther, "/")
}
