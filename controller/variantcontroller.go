package controller

import (
	"fmt"
	"net/http"
	"pb-dev-be/models"
	"time"

	"github.com/labstack/echo/v4"
)

func FetchAllVariantData(c echo.Context) error {
	fmt.Println("GET Variant END POINT HIT!")

	result, err := models.FetchAllVariantData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})

	}

	return c.JSON(http.StatusOK, result)
}

func StoreVariant(c echo.Context) error {
	fmt.Println("POST Variant END POINT HIT!")

	var variant models.Variant

	// variant.Id = c.FormValue("variant_id")
	variant.Name = c.FormValue("name")
	variant.UserId = c.FormValue("user_id")
	variant.Created_at = time.Now().String()
	variant.Modified_at = time.Now().String()

	result, err := models.StoreVariantData(variant)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)

}

func UpdateVariant(c echo.Context) error {
	fmt.Println("PUT Variant END POINT HIT!")

	var variant models.Variant

	variant.Id = c.Param("id")
	variant.Name = c.FormValue("name")
	variant.UserId = c.FormValue("user_id")
	variant.Created_at = time.Now().String()
	variant.Modified_at = time.Now().String()

	result, err := models.UpdateVariant(variant)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteVariant(c echo.Context) error {
	fmt.Println("DELETE Variant END POINT HIT!")

	id := c.Param("id")

	result, err := models.DeleteVariant(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
