package controller

import (
	"net/http"
	"pb-dev-be/helpers"

	"github.com/labstack/echo/v4"
)

func SendEmail(c echo.Context) error {
	err := helpers.SendMailTest()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Email Sent"})
}
