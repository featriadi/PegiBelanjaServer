package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pb-dev-be/models"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func FetchAllCustomerData(c echo.Context) error {
	fmt.Println("GET Customer END POINT Hit !")

	result, err := models.FetchAllCustomerData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func ShowCustomerDataById(c echo.Context) error {
	fmt.Println("GET Customer By ID END POINT Hit !")

	param_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	result, err := models.ShowCustomerById(param_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StoreCustomerData(c echo.Context) error {
	fmt.Println("POST Customer END POINT HIT")

	var cust = new(models.Customer)

	err := c.Bind(cust)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.StoreCustomerData(*cust)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Process Error : " + err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StoreCustomerWithUser(c echo.Context) error {
	fmt.Println("POST Customer END POINT HIT")

	var cust = new(models.Customer)
	var user = new(models.User)

	//Read The Content (Request)
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	}

	//Restore the io.ReadCloser to its original state
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// bodyString := string(bodyBytes)

	//Bind To Mitra Data
	err := json.Unmarshal(bodyBytes, &cust)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	cust.Created_at = time.Now().Format("2006-01-02 15:04:05")
	cust.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	//Bind To User Data
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error User : " + err.Error()})
	}

	result, err := models.StoreCustomerAndUserData(*cust, user.Password)
	if err != nil {
		return c.JSON(http.StatusOK, result)
	}

	return c.JSONPretty(http.StatusOK, result, "")
}

func StoreCustomer(c echo.Context) error {
	fmt.Println("POST Customer END POINT HIT")

	var cust = new(models.Customer)

	err := c.Bind(cust)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.StoreCustomerData(*cust)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Process Error : " + err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateCustomer(c echo.Context) error {
	fmt.Println("PUT Customer End Point Hit")

	param_id := c.Param("id")

	var cust = new(models.Customer)
	err := c.Bind(cust)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.UpdateCustomer(*cust, param_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Process Error : " + err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func DeleteCustomer(c echo.Context) error {
	fmt.Println("DELETE Customer End Point Hit")

	param_id := c.Param("id")

	result, err := models.DeleteCustomer(param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func CheckCustomerByEmail(c echo.Context) error {
	fmt.Println("Check Customer End Point Hit")

	email := c.Param("email")

	result, err := models.CheckCustomer(email)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
