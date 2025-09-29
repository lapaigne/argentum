package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (e *Endpoints) me(c echo.Context) error {
	data_old.Fill()

	claims := getClaims(c)
	uid := claims.UID

	ctx := struct {
		Data   Data
		Helper Helper
	}{
		Data:   data,
		Helper: helper,
	}

	fmt.Println(ctx)

	ctx0 := struct {
		Raw    RawData0
		Mapped MappedData0
		Helper Helper0
		UID    int
	}{
		Raw:    data_old.Raw,
		Mapped: data_old.Mapped,
		Helper: data_old.Helper,
		UID:    uid,
	}

	return c.Render(200, "me", ctx0)
}

func (e *Endpoints) me_s(c echo.Context) error {

	url := c.Request().URL.Path
	switch url {
	default:
	}

	fmt.Println(url)

	return echo.ErrTeapot
}
