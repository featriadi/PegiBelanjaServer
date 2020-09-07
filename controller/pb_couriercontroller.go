package controller

import (
	"net/http"
	"pb-dev-be/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func FetchAllPBCourier(c echo.Context) error {
	result, err := models.FetchAllPBCourier()
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func StorePBCourier(c echo.Context) error {
	var pb_courier = new(models.PBCourier)

	err := c.Bind(pb_courier)

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": "Bind Error - " + err.Error()}, "  ")
	}

	result, err := models.StorePBCourier(*pb_courier)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func UpdatePBCourier(c echo.Context) error {
	var pb_courier = new(models.PBCourier)

	err := c.Bind(pb_courier)

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": "Bind Error - " + err.Error()}, "  ")
	}

	result, err := models.UpdatePBCourier(*pb_courier)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func DeletePBCourier(c echo.Context) error {
	param_id := c.Param("id")

	id, err := strconv.Atoi(param_id)

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": "Error While Converting Parameter To Int - " + err.Error()}, "  ")
	}

	result, err := models.DeletePBCourier(id)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}
