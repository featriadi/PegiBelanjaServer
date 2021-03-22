package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pb-dev-be/models"
	"time"

	"github.com/labstack/echo/v4"
)

func FetchAllMitraData(c echo.Context) error {
	fmt.Println("GET Mitra END POINT Hit !")

	result, err := models.FetchAllMitraData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StoreMitra(c echo.Context) error {
	fmt.Println("POST Mitra END POINT HIT")

	var mitra = new(models.Mitra)

	err := c.Bind(mitra)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.StoreMitraData(*mitra)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Process Error : " + err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StoreMitraWithUser(c echo.Context) error {
	fmt.Println("POST Mitra END POINT HIT")

	var mitra = new(models.Mitra)
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
	err := json.Unmarshal(bodyBytes, &mitra)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	mitra.Created_at = time.Now().Format("2006-01-02 15:04:05")
	mitra.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	//Bind To User Data
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error User : " + err.Error()})
	}

	result, err := models.StoreMitraAndRegisterUser(*mitra, user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Process Error : " + err.Error()})
	}

	return c.JSONPretty(http.StatusOK, result, "")
}

func UpdateMitra(c echo.Context) error {
	fmt.Println("PUT Mitra End Point Hit")

	param_id := c.Param("id")

	var mitra = new(models.Mitra)
	err := c.Bind(mitra)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.UpdateMitra(*mitra, param_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Process Error : " + err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func DeleteMitra(c echo.Context) error {
	fmt.Println("DELETE Mitra End Point Hit")

	param_id := c.Param("id")

	result, err := models.DeleteMitra(param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func ShowMitraDataById(c echo.Context) error {
	fmt.Println("GET Mitra By ID END POINT HIT!")

	param_id := c.Param("id")
	result, err := models.GetMitraById(param_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})

	}

	return c.JSON(http.StatusOK, result)
}
