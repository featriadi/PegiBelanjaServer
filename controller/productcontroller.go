package controller

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pb-dev-be/models"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func GetTotalProducts(c echo.Context) error {
	result, err := models.GetTotalProduts()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func FetchAllProductData(c echo.Context) error {
	fmt.Println("GET Product END POINT HIT!")

	query_param := c.QueryParam("newest_product")
	isp := c.QueryParam("isp")
	_start := c.QueryParam("_start")
	_limit := c.QueryParam("_limit")
	_userid := c.FormValue("user_id")

	if query_param == "" {
		query_param = "false"
	}

	if isp == "" {
		isp = "false"
	}

	newest_product, err := strconv.ParseBool(query_param)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Parse Boolean : " + err.Error()})
	}

	_isp, err := strconv.ParseBool(isp)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Parse Boolean 202 : " + err.Error()})
	}

	result, err := models.FetchAllProductData(newest_product, _start, _limit, _isp, _userid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func ShowProductDataById(c echo.Context) error {
	fmt.Println("GET Product By ID END POINT HIT!")

	param_id := c.Param("id")
	result, err := models.ShowProductById(param_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})

	}

	return c.JSON(http.StatusOK, result)
}

func StoreProduct(c echo.Context) error {
	// fmt.Println("POST Product End Point Hit")

	var product = new(models.Product)

	err := c.Bind(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	query_param := c.QueryParam("is_mitra")
	if query_param == "" {
		query_param = "false"
	}

	is_mitra, err := strconv.ParseBool(query_param)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Parse Boolean : " + err.Error()})
	}

	// filedir, err := UploadMultipleFile(c, product.Id)

	// if err != nil {
	// 	if err != http.ErrMissingFile {
	// 		fmt.Println(err.Error())
	// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error(), "status": http.StatusInternalServerError})
	// 	}
	// }

	// for idx := range filedir {
	// 	filepath := filedir[idx]
	// 	product_media := new(models.Media)

	// 	product_media.FilePath = filepath
	// 	product.ProductMedia = append(product.ProductMedia, *product_media)
	// }

	result, err := models.StoreProductData(*product, is_mitra)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}
	return c.JSON(http.StatusOK, result)
}

func Test(c echo.Context) error {

	a := c.FormValue("content")

	return c.JSON(http.StatusOK, a)
	// filedir, err := UploadMultipleFile(c, "aaaa")
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Process Error : " + err.Error()})
	// }
	// return c.JSON(http.StatusOK, filedir)

}

func UploadMultipleFile(c echo.Context, product_id string) ([]string, error) {
	//this function returns the filename(to save in database) of the saved file or an error if it occurs
	err := c.Request().ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}

	//ParseMultipartForm parses a request body as multipart/form-data
	fhs := c.Request().MultipartForm.File["image_data"]
	if fhs == nil {
		return nil, errors.New("There's No File On Request")

	}

	var arr []string
	for idx, file := range fhs {
		if err != nil {
			return nil, err
		}

		ext := filepath.Ext(file.Filename)

		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create("static/image/product/" + time.Now().Format("20060102150405") + "_" + product_id + "_" + strconv.Itoa(idx) + ext)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return nil, err
		}
		arr = append(arr, dst.Name())
	}

	return arr, nil
}

func UpdateProduct(c echo.Context) error {
	fmt.Println("PUT Product End Point Hit")

	var product = new(models.Product)
	param_id := c.Param("id")

	err := c.Bind(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bind Error : " + err.Error()})
	}

	result, err := models.UpdateProductData(*product, param_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}
	return c.JSON(http.StatusOK, result)
}

func DeleteProduct(c echo.Context) error {
	fmt.Println("DELETE Product End Point Hit")

	param_id := c.Param("id")

	result, err := models.DeleteProduct(param_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, result)
	}

	return c.JSON(http.StatusOK, result)
}

func UploadFileProduct(c echo.Context) error {
	product_id := c.FormValue("product_id")

	filedir, err := UploadMultipleFile(c, product_id)
	var file_dir []interface{}

	if err != nil {
		if err != http.ErrMissingFile {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error(), "status": http.StatusInternalServerError})
		} else {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error(), "status": http.StatusInternalServerError})
		}
	}

	for idx := range filedir {
		filepath := filedir[idx]
		product_media := new(models.Media)

		product_media.ItemNumber = idx
		product_media.FilePath = filepath
		file_dir = append(file_dir, product_media)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": file_dir})
}
