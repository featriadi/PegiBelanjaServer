package controller

import (
	"fmt"
	"net/http"
	"pb-dev-be/models"

	"github.com/labstack/echo/v4"
)

func FetchAllTimeDeliveryData(c echo.Context) error {
	fmt.Println("GET Time Delivery END POINT HIT!")

	result, err := models.FetchAllTimeDeliveryData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})

	}

	return c.JSON(http.StatusOK, result)
}

func StoreTimeDelivery(c echo.Context) error {
	fmt.Println("POST Time Delivery END POINT HIT!")

	var td = new(models.TimeDelivery)
	err := c.Bind(td)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.StoreTimeDelivery(*td)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateTimeDelivery(c echo.Context) error {
	fmt.Println("PUT Time Delivery END POINT HIT!")

	var td = new(models.TimeDelivery)
	param_id := c.Param("id")

	err := c.Bind(td)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.UpdateTimeDelivery(*td, param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// func DeleteTimeDelivery(c echo.Context) error {
// 	fmt.Println("DELETE Time Delivery END POINT HIT!")

// 	id := c.FormValue("time_delivery_id")

// 	result, err := models.DeleteTimeDelivery(id)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusOK, result)
// }
