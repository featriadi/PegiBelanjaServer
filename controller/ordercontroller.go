package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pb-dev-be/config"
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

func GetTransactionStatus(c echo.Context) error {
	conf := config.GetConfig()
	param := c.Param("order_id")
	// url := "https://api.sandbox.midtrans.com/v2/" + param + "/status"
	url := "https://api.midtrans.com/v2/" + param + "/status"
	// str := base64.StdEncoding.EncodeToString([]byte(conf.MIDTRANS_SERVER_KEY_SANDBOX))
	str := base64.StdEncoding.EncodeToString([]byte(conf.MIDTRANS_SERVER_KEY))

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+str)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	defer resp.Body.Close()
	fmt.Println(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	var gdat map[string]interface{}
	if err := json.Unmarshal(body, &gdat); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	return c.JSONPretty(http.StatusOK, gdat, "  ")
}
