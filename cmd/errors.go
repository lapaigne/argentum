package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	code := 500
	msg := "internal server error"

	if herr, ok := err.(*echo.HTTPError); ok {
		code = herr.Code
		if herr.Message != nil {
			msg = fmt.Sprintf("%v", herr.Message)
		}
	} else if err != nil {
		msg = err.Error()
	}

	if code == 401 {
		if public[c.Request().URL.Path] {
			return
		}

		if c.Request().Header.Get("HX-Request") == "true" {
			_ = c.NoContent(401)
		} else {
			_ = c.Redirect(303, "/signin")
		}
		return
	}

	if !c.Response().Committed {
		_ = c.JSON(code, map[string]any{
			"error": msg,
		})
	}
}
