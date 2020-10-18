package controller

import (
	"fmt"
	"net/http"
	"pb-dev-be/models"
	"time"

	"github.com/labstack/echo/v4"
)

func FetchAllBankData(c echo.Context) error {
	fmt.Println("GET Bank END POINT HIT!")

	result, err := models.FetchAllBankData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StoreBank(c echo.Context) error {
	fmt.Println("POST Bank END POINT HIT!")

	var bank models.Bank

	bank.Id = c.FormValue("bank_id")
	bank.Name = c.FormValue("name")
	bank.Owner = c.FormValue("owner")
	bank.AccountNumber = c.FormValue("account_number")
	bank.UserId = c.FormValue("user_id")
	bank.Created_at = time.Now().Format("2006-01-02 15:04:05")
	bank.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	result, err := models.StoreBank(bank)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateBank(c echo.Context) error {
	fmt.Println("PUT Bank END POINT HIT!")

	var bank models.Bank

	param_id := c.Param("id")
	bank.Id = c.FormValue("bank_id")
	bank.Name = c.FormValue("name")
	bank.Owner = c.FormValue("owner")
	bank.AccountNumber = c.FormValue("account_number")
	bank.UserId = c.FormValue("user_id")
	bank.Created_at = time.Now().Format("2006-01-02 15:04:05")
	bank.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	result, err := models.UpdateBank(bank, param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteBank(c echo.Context) error {
	fmt.Println("DELETE Bank END POINT HIT!")

	id := c.FormValue("bank_id")

	result, err := models.DeleteBank(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
