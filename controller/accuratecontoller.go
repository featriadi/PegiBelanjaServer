package controller

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func Accurate(c echo.Context) error {
	contentType := c.Response().Header()
	fmt.Println(contentType, "get accurate")

	// req := c.Request()
	// res := c.Response()

	return c.JSON(fmt.Println(contentType, "get accurate"))
}
