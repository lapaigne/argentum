package main

import "github.com/labstack/echo/v4"

func (e Endpoints) alltasks(c echo.Context) error {
	d := struct {
		Data   *Data
		Helper *Helper
	}{
		Data:   &data,
		Helper: &helper,
	}

	return c.Render(200, "alltasks", d)
}
