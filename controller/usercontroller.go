package controller

import (
	"fmt"
	"net/http"
	"pb-dev-be/models"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func StoreUser(c echo.Context) error {
	fmt.Println("POST User END POINT HIT!")

	var user models.User

	user.Name = c.FormValue("name")
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	user.UserRole = c.FormValue("user_role")
	user.IsVerified, _ = strconv.ParseBool(c.FormValue("verified"))
	user.UserCreated = c.FormValue("user_created")
	user.Created_at = time.Now().String()
	user.Modified_at = time.Now().String()
	user.LastLogin = time.Now().String()

	result, err := models.StoreUserData(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func FetchAllUserData(c echo.Context) error {
	fmt.Println("GET User END POINT HIT!")

	result, err := models.FetchAllUserData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})

	}

	return c.JSON(http.StatusOK, result)
}
