package controller

import (
	"net/http"
	"pb-dev-be/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CreateStockHInAndOut(c echo.Context) error {
	var stockh models.Stock

	parse_qty, err := strconv.ParseFloat(c.FormValue("qty"), 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error While Parsing Float : " + err.Error()})
	}

	is_in, err := strconv.ParseBool(c.FormValue("type"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error While Parsing Boolean : " + err.Error()})
	}

	stockh.ProductId = c.FormValue("product_id")
	stockh.Qty = parse_qty
	stockh.UserId = c.FormValue("user_id")

	result, err := models.CreateStockInAndOut(stockh, is_in)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
