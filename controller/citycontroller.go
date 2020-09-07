package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pb-dev-be/models"

	"github.com/labstack/echo/v4"
)

func FetchAllCity(c echo.Context) error {
	var result models.Response
	var err error
	query_param := c.QueryParam("province_id")

	if query_param == "" {
		result, err = models.FetchAllCityData()
	} else {
		result, err = models.ShowCityByProvinceIdData(query_param)
	}

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")

}

func StoreCity(c echo.Context) error {
	url := "https://pro.rajaongkir.com/api/city"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("key", "355d769dde694b80a08280d6331ef125")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "  ")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var gdat map[string]interface{}
	var obj = new([]models.City)
	if err := json.Unmarshal(body, &gdat); err != nil {
		panic(err)
	}

	dat := gdat["rajaongkir"]
	mapDat := dat.(map[string]interface{})
	pDat := mapDat["results"]

	byteData, _ := json.Marshal(pDat)
	err = json.Unmarshal(byteData, &obj)

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "  ")
	}

	result, err := models.StoreCityBulk(*obj)

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}
