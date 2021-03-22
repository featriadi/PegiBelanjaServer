package controller

import (
	"net/http"
	"pb-dev-be/models"

	"github.com/labstack/echo/v4"
)

func CreateVerification(c echo.Context) error {

	var ver = new(models.VerificationCode)

	err := c.Bind(ver)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	res, err := models.StoreVerificationData(*ver)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}
