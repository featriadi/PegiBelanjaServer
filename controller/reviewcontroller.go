package controller

import (
	"net/http"
	"pb-dev-be/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func FetchReviewByProductId(c echo.Context) error {

	product_id := c.Param("product_id")
	result, err := models.GetReviewByProductId(product_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func FetchAllReview(c echo.Context) error {
	result, err := models.GetAllReview()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StoreReview(c echo.Context) error {
	var review models.Review

	rating, err := strconv.ParseFloat(c.FormValue("rating"), 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error While Parsing Float : " + err.Error()})
	}

	review.Content = c.FormValue("content")
	review.Rating = rating
	review.ProductId = c.FormValue("product_id")
	review.CustomerId = c.FormValue("customer_id")

	result, err := models.CreateReview(review)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
