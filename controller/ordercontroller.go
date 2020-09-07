package controller

import (
	"net/http"
	"pb-dev-be/models"

	"github.com/labstack/echo/v4"
)

func CreateOrder(c echo.Context) error {
	var order = new(models.Order)

	err := c.Bind(order)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.CreateOrder(*order)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}

	return c.JSON(http.StatusOK, result)
}
