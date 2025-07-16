package main

import "github.com/labstack/echo/v4"

type Endpoints struct{}

var endpoints Endpoints

func (e Endpoints) AddTask(c echo.Context) error {
	return nil
}
