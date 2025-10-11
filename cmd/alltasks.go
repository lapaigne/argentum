package main

import "github.com/labstack/echo/v4"

func (e Endpoints) alltasks_GET(c echo.Context) error {
	d := struct {
		Data   *Data
		Helper *Helper
	}{
		Data:   &data,
		Helper: &helper,
	}

	return c.Render(200, "alltasks", d)
}

func (e Endpoints) alltasks_worker(c echo.Context) error {
	return nil
}

func (e Endpoints) alltasks_created(c echo.Context) error {
	return nil
}

func (e Endpoints) alltasks_addr(c echo.Context) error {
	return nil
}
