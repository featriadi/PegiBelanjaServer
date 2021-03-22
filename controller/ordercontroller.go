package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"pb-dev-be/models"
	"strconv"

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

func PostOrder(c echo.Context) error {
	var order = new(models.Order)
	formData := url.Values{
		"name": {"rey"},
	}

	err := c.Bind(order)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := http.PostForm("https://app.sandbox.midtrans.com/snap/v1/transactions", formData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}

	return c.JSON(http.StatusOK, result)
}

func CreateOrderTracking(c echo.Context) error {
	var tracking = new(models.OrderTracking)

	err := c.Bind(tracking)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}
	waybill := c.QueryParam("waybill")

	result, err := models.CreateOrderTracking(*tracking, waybill)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}

	return c.JSON(http.StatusOK, result)
}

func GetOrderStats(c echo.Context) error {

	result, err := models.GetOrderStats()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}

	return c.JSON(http.StatusOK, result)
}

func GetOrderTracking(c echo.Context) error {
	order_id := c.Param("order_id")

	result, err := models.GetOrderTracking(order_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateWaybillOrder(c echo.Context) error {
	order_id := c.Param("order_id")
	waybill := c.FormValue("waybill")

	_Id, err := strconv.Atoi(order_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error While Converting String"})
	}

	result, err := models.UpdateWaybillOrder(_Id, waybill)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}

	return c.JSON(http.StatusOK, result)
}

func GetOrderData(c echo.Context) error {
	fmt.Println("GET Order END POINT HIT!")

	_isAdmin := c.QueryParam("isAdmin")
	_isMitra := c.QueryParam("isMitra")
	_userid := c.QueryParam("user_id")
	_paramS := c.QueryParam("paramS")

	_startMY := c.QueryParam("startMY")
	_endMY := c.QueryParam("endMY")

	if _isAdmin == "" {
		_isAdmin = "false"
	}

	if _isMitra == "" {
		_isMitra = "false"
	}

	isAdmin, err := strconv.ParseBool(_isAdmin)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Parse Boolean : " + err.Error()})
	}

	isMitra, err := strconv.ParseBool(_isMitra)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Parse Boolean 202 : " + err.Error()})
	}

	result, err := models.GetOrder(isAdmin, isMitra, _userid, _paramS, _startMY, _endMY)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func GetTransactionStatus(c echo.Context) error {
	param := c.Param("order_id")

	result, err := models.GetTransactionStatus(param)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}
