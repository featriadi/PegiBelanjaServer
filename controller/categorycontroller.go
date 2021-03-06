package controller

import (
	"fmt"
	"net/http"
	"pb-dev-be/models"

	"github.com/labstack/echo/v4"
)

func FetchAllCategoryData(c echo.Context) error {
	fmt.Println("GET Category END POINT HIT!")

	result, err := models.FetchAllCategoryData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})

	}

	return c.JSON(http.StatusOK, result)
}

func StoreCategory(c echo.Context) error {
	fmt.Println("POST Category END POINT HIT!")

	var cat = new(models.Category)
	err := c.Bind(cat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.StoreCategory(*cat)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateCategory(c echo.Context) error {
	fmt.Println("PUT Category END POINT HIT!")

	var cat = new(models.Category)
	param_id := c.Param("id")

	err := c.Bind(cat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.UpdateCategory(*cat, param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteCategory(c echo.Context) error {
	fmt.Println("DELETE Category END POINT HIT!")

	id := c.FormValue("category_id")

	result, err := models.DeleteCategory(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
