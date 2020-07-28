package controller

import (
	"fmt"
	"net/http"
	"pb-dev-be/models"

	"github.com/labstack/echo/v4"
)

func FetchAllCouponData(c echo.Context) error {
	fmt.Println("GET Coupon END POINT HIT!")

	result, err := models.FetchAllCouponData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})

	}

	return c.JSON(http.StatusOK, result)
}

func StoreCouponData(c echo.Context) error {
	fmt.Println("POST Coupon END POINT HIT!")

	var coupon = new(models.Coupon)

	err := c.Bind(coupon)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.StoreCouponData(*coupon)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateCouponData(c echo.Context) error {
	fmt.Println("PUT Coupon End Point Hit")

	param_id := c.Param("id")

	var coupon = new(models.Coupon)
	err := c.Bind(coupon)
	coupon.Id = param_id
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.UpdateCouponData(*coupon)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Process Error : " + err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func DeleteCoupon(c echo.Context) error {
	fmt.Println("DELETE Coupon END POINT HIT!")

	id := c.Param("id")

	result, err := models.DeleteCoupon(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
