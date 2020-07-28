package controller

import (
	"fmt"
	"net/http"
	"pb-dev-be/models"
	"time"

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
	fmt.Println("POST Bank END POINT HIT!")

	var cat models.Category

	cat.Id = c.FormValue("category_id")
	cat.Name = c.FormValue("name")
	cat.UserId = c.FormValue("user_id")
	cat.Created_at = time.Now().String()
	cat.Modified_at = time.Now().String()

	result, err := models.StoreCategory(cat)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateCategory(c echo.Context) error {
	fmt.Println("PUT Bank END POINT HIT!")

	var cat models.Category
	param_id := c.Param("id")
	cat.Id = c.FormValue("category_id")
	cat.Name = c.FormValue("name")
	cat.UserId = c.FormValue("user_id")
	cat.Created_at = time.Now().String()
	cat.Modified_at = time.Now().String()

	result, err := models.UpdateCategory(cat, param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteCategory(c echo.Context) error {
	fmt.Println("DELETE Bank END POINT HIT!")

	id := c.FormValue("category_id")

	result, err := models.DeleteCategory(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
