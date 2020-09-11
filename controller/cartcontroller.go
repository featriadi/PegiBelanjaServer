package controller

import (
	"net/http"
	"pb-dev-be/models"

	"github.com/labstack/echo/v4"
)

func GetCartByCustomerId(c echo.Context) error {

	customer_id := c.QueryParam("customer_id")

	result, err := models.GetCartById(customer_id)

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func UpdateCart(c echo.Context) error {
	var cart = new(models.Cart)

	err := c.Bind(cart)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.UpdateCart(*cart)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}

	return c.JSON(http.StatusOK, result)
}

func CreateCart(c echo.Context) error {
	var cart = new(models.Cart)

	err := c.Bind(cart)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.CreateCart(*cart)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteCart(c echo.Context) error {

	param_id := c.Param("customer_id")

	result, err := models.DeleteCartByCustomerId(param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}

	return c.JSON(http.StatusOK, result)
}
