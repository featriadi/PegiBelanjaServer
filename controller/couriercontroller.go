package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pb-dev-be/models"
	"strings"

	"github.com/labstack/echo/v4"
)

type RequestCostParam struct {
	Origin          string
	OriginType      string
	Destination     string
	DestinationType string
	Weight          string
	Courier         string
}

func FetchAllCourier(c echo.Context) error {
	result, err := models.FetchAllCourier()

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func GetThirdPartyCourier(c echo.Context) error {
	url := "https://pro.rajaongkir.com/api/cost"

	var param = new(RequestCostParam)

	err := c.Bind(param)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	payload := strings.NewReader("origin=" + param.Origin + "&originType=" + param.OriginType + "&destination=" + param.Destination + "&destinationType=" + param.DestinationType + "&weight=" + param.Weight + "&courier=" + param.Courier)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("key", "355d769dde694b80a08280d6331ef125")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var gdat map[string]interface{}
	if err := json.Unmarshal(body, &gdat); err != nil {
		panic(err)
	}

	return c.JSONPretty(http.StatusOK, gdat, "  ")
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
