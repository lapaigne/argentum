package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
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
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	if !c.Response().Committed {
		c.JSON(code, map[string]any{
			"error": msg,
		})
	}
}
