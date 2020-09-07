package controller

import (
	"net/http"
	"pb-dev-be/models"

	"github.com/labstack/echo/v4"
)

func FetchAllCourier(c echo.Context) error {
	result, err := models.FetchAllCourier()

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func StoreCourier(c echo.Context) error {
	var courier = new(models.Courier)

	err := c.Bind(courier)

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": "Bind Error - " + err.Error()}, "  ")
	}

	result, err := models.StoreCourier(*courier)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func UpdateCourier(c echo.Context) error {
	var courier = new(models.Courier)

	err := c.Bind(courier)
	param_id := c.Param("id")

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": "Bind Error - " + err.Error()}, "  ")
	}

	result, err := models.UpdateCourier(*courier, param_id)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func DeleteCourier(c echo.Context) error {
	param_id := c.Param("id")

	id := param_id

	result, err := models.DeleteCourier(id)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}
