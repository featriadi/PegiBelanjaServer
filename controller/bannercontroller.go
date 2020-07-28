package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pb-dev-be/models"
	"time"

	"github.com/labstack/echo/v4"
)

func FetchAllBannerData(c echo.Context) error {
	fmt.Println("GET Banner END POINT HIT!")

	result, err := models.FetchAllBannerData()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StoreBanner(c echo.Context) error {
	fmt.Println("POST Banner END POINT HIT!")

	var banner models.Banner
	filedir, err := UploadFile(c)

	if err != nil {
		if err != http.ErrMissingFile {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error(), "status": http.StatusInternalServerError})
		}
	}

	banner.Id = c.FormValue("banner_id")
	banner.Name = c.FormValue("name")
	banner.Link = c.FormValue("link")
	banner.ImageData = filedir
	banner.UserId = c.FormValue("user_id")
	banner.Created_at = time.Now().String()
	banner.Modified_at = time.Now().String()

	result, err := models.StoreBanner(banner)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func UploadFile(c echo.Context) (string, error) {

	//this function returns the filename(to save in database) of the saved file or an error if it occurs
	err := c.Request().ParseMultipartForm(10 << 20)
	if err != nil {
		return "", err
	}

	//ParseMultipartForm parses a request body as multipart/form-data
	file, err := c.FormFile("image_data")

	if err == http.ErrMissingFile {
		return "static/image/banner/default.jpg", nil
	}

	ext := filepath.Ext(file.Filename)

	// fmt.Println(ext)

	if err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("static/image/banner/" + time.Now().Format("20060102150405") + "_" + c.FormValue("banner_id") + ext)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return dst.Name(), nil

	// filedir, err := filepath.Abs(filepath.Dir())

}

func UpdateBanner(c echo.Context) error {
	fmt.Println("PUT Banner END POINT HIT!")

	var banner models.Banner
	filedir, err := UploadFile(c)
	param_id := c.Param("id")
	if err != nil {
		if err != http.ErrMissingFile {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error(), "status": http.StatusInternalServerError})
		}
	}

	banner.Id = c.FormValue("banner_id")
	banner.Name = c.FormValue("name")
	banner.Link = c.FormValue("link")
	banner.ImageData = filedir
	banner.UserId = c.FormValue("user_id")
	banner.Created_at = time.Now().String()
	banner.Modified_at = time.Now().String()

	result, err := models.UpdateBanner(banner, param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteBanner(c echo.Context) error {
	fmt.Println("DELETE Banner END POINT HIT!")

	// filedir, err := UploadFile(c)
	param_id := c.Param("id")

	result, err := models.DeleteBanner(param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
