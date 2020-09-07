package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pb-dev-be/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func FetchAllSubDistrict(c echo.Context) error {
	var result models.Response
	var err error
	query_param := c.QueryParam("city_id")

	if query_param == "" {
		result, err = models.FetchAllSubDistrictData()
	} else {
		result, err = models.ShowSubDistrictByCityId(query_param)
	}

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, result, "  ")
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func StoreSubdistrict(c echo.Context) error {
	var result models.Response
	for i := 1; i < 502; i++ {
		url := "https://pro.rajaongkir.com/api/subdistrict?city=" + strconv.Itoa(i)

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("key", "355d769dde694b80a08280d6331ef125")

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "  ")
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		var gdat map[string]interface{}
		var obj = new([]models.SubDistrict)
		if err = json.Unmarshal(body, &gdat); err != nil {
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

		result, err = models.StoreSubDistrictBulk(*obj)
		if err != nil {
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "  ")
		}
	}
	return c.JSONPretty(http.StatusOK, result, "  ")
}
