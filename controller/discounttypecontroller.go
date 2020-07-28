package controller

import (
	"fmt"
	"net/http"
	"pb-dev-be/models"

	"github.com/labstack/echo/v4"
)

func FetchAllDiscountData(c echo.Context) error {
	fmt.Println("GET Discount Data END POINT Hit !")

	result, err := models.FetchAllDiscountType()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StoreDiscountData(c echo.Context) error {
	fmt.Println("POST Discount Data END POINT HIT")

	var discType models.DiscountType

	discType.Id = c.FormValue("id")
	discType.Name = c.FormValue("name")
	discType.UserCreated = c.FormValue("user_id")

	result, err := models.StoreDiscountType(discType)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateDiscountData(c echo.Context) error {
	fmt.Println("PUT Discount Data END POINT HIT")

	var discType models.DiscountType

	discType.Id = c.Param("id")
	discType.Name = c.FormValue("name")
	discType.UserCreated = c.FormValue("user_id")

	result, err := models.UpdateDiscountTypeData(discType)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})

	}

	return c.JSON(http.StatusOK, result)
}

func DeleteDiscountData(c echo.Context) error {
	fmt.Println("DELETE Discount Data END POINT HIT")

	param_id := c.Param("id")

	result, err := models.DeleteDiscountTypeData(param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
