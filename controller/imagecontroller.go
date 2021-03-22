package controller

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func GetEncodedImage(c echo.Context) error {
	fmt.Println("GET Encoded Image END POINT HIT!")

	_type := c.Param("type")
	file_name := c.Param("file_name")

	file_dir := ""
	if _type == "product" {
		file_dir = "static/image/product/" + file_name
	} else {
		file_dir = "static/image/banner/" + file_name
	}
	imgFile, err := os.Open(file_dir) // a QR code image

	if err != nil {
		os.Exit(1)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Open File Error : " + err.Error()})

	}
	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	_, err = fReader.Read(buf)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "File Read Error : " + err.Error()})
	}

	// if you create a new image instead of loading from file, encode the image to buffer instead with png.Encode()

	// png.Encode(&buf, image)

	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)

	return c.JSON(http.StatusOK, map[string]string{"data": imgBase64Str})

	// Embed into an html without PNG file
	// img2html := "<html><body><img src=\"data:image/png;base64," + imgBase64Str + "\" /></body></html>"

	// w.Write([]byte(fmt.Sprintf(img2html)))
}
