package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"time"
)

const ObjectId = "PRODUCT"

type Product struct {
	Id                    string           `json:"id"`
	Name                  string           `json:"name"`
	CategoryId            string           `json:"categoryid"`
	Description           string           `json:"description"`
	MinOrder              int              `json:"minimal_order"`
	IsVariant             bool             `json:"using_variant"`
	IsDiscount            bool             `json:"using_discount"`
	IsReadyStock          bool             `json:"ready_stock"`
	IsWholesalePrice      bool             `json:"using_wholesale_price"`
	PublishStatus         bool             `json:"publish_status"`
	ProductDetails        []Details        `json:"details"`
	ProductMedia          []Media          `json:"media"`
	ProductDiscount       Discount         `json:"discount"`
	ProductWholesalePrice []WholesalePrice `json:"wholesale_price"`
	UserId                string           `json:"user_id"`
	Created_at            string           `json:"created_at"`
	Modified_at           string           `json:"modified_at"`
}

type Details struct {
	ItemNumber  int     `json:"index"`
	Weight      int     `json:"weight"`
	BuyPrice    float64 `json:"buy_price"`
	SellPrice   float64 `json:"sell_price"`
	VKey        int     `json:"v_key"`
	VariantType string  `json:"variant_type"`
	Content     string  `json:"content"`
	Stock       float64 `json:"stock"`
}

type Discount struct {
	Value int `json:"value"`
}

type WholesalePrice struct {
	ItemNumber    int     `json:"index"`
	MinimumAmount int     `json:"minimum_amount"`
	Price         float64 `json:"price"`
}

type Media struct {
	ItemNumber int    `json:"index"`
	FilePath   string `json:"file_path"`
}

func FetchAllProductData(is_newest_product bool, start string, limit string, isp bool, user_id string) (Response, error) {
	var res Response
	var arrobj []Product
	var product Product

	con := db.CreateCon()
	qry := ""
	if is_newest_product {
		qry = `SELECT A.* FROM smc_product A
		LEFT JOIN smc_approved B on B.s_id = A.s_sku_id and B.s_object_id = 'PRODUCT'
		WHERE A.s_created_at > now() - interval 7 day 
		and B.s_status = 'VERIFIED' and A.s_publish_status = 1
		ORDER BY A.s_created_at DESC LIMIT 12`
	} else if start != "" {
		qry = `SELECT A.* FROM smc_product A
		LEFT JOIN smc_approved B on B.s_id = A.s_sku_id and B.s_object_id = 'PRODUCT'
		WHERE B.s_status = 'VERIFIED' and A.s_publish_status = 1 LIMIT ` + start + `,` + limit
	} else if isp {
		qry = `SELECT A.* FROM smc_product A 
		LEFT JOIN smc_approved B on B.s_id = A.s_sku_id and B.s_object_id = 'PRODUCT'
		WHERE
		B.s_status = 'VERIFIED' and A.s_publish_status = 1`
	} else if user_id != "" {
		qry = "SELECT * FROM smc_product WHERE s_user_id = '" + user_id + "'"
	} else {
		qry = "SELECT * FROM smc_product"
	}
	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = product
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&product.Id, &product.Name, &product.CategoryId, &product.Description,
			&product.MinOrder, &product.IsVariant, &product.IsDiscount, &product.IsReadyStock,
			&product.IsWholesalePrice, &product.PublishStatus, &product.UserId,
			&product.Created_at, &product.Modified_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = product
			return res, err
		}

		// fmt.Println(product.Id)
		res, err, product := GetProductDetails(con, product)

		if err != nil {
			return res, err
		}

		res2, err, product := GetProductMedia(con, product)
		if err != nil {
			return res2, err
		}

		res3, err, product := GetProductDiscount(con, product)

		if err != nil {
			return res3, err
		}

		res5, err, product := GetProductWholesalePrice(con, product)

		if err != nil {
			return res5, err
		}

		arrobj = append(arrobj, product)
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrobj

	return res, nil
}

func ShowProductById(param_id string) (Response, error) {
	var res Response
	var product Product

	con := db.CreateCon()

	qry := "SELECT * FROM smc_product WHERE s_sku_id = ?"

	rows, err := con.Query(qry, param_id)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Product{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&product.Id, &product.Name, &product.CategoryId, &product.Description,
			&product.MinOrder, &product.IsVariant, &product.IsDiscount, &product.IsReadyStock,
			&product.IsWholesalePrice, &product.PublishStatus, &product.UserId,
			&product.Created_at, &product.Modified_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Product{}
			return res, err
		}

		// fmt.Println(product.Id)
		res, err, prod := GetProductDetails(con, product)
		product.ProductDetails = prod.ProductDetails
		if err != nil {
			return res, err
		}

		res3, err, prod := GetProductMedia(con, product)
		product.ProductMedia = prod.ProductMedia
		if err != nil {
			return res3, err
		}

		res2, err, prod := GetProductDiscount(con, product)
		product.ProductDiscount = prod.ProductDiscount
		if err != nil {
			return res2, err
		}

		res4, err, prod := GetProductWholesalePrice(con, product)
		product.ProductWholesalePrice = prod.ProductWholesalePrice
		if err != nil {
			return res4, err
		}

	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = product

	return res, nil
}

func GetProductDetails(con *sql.DB, product Product) (Response, error, Product) {
	var res Response
	var detail Details

	qry_details := `SELECT s_item_number, s_weight, s_buy_price, s_sell_price, s_v_key,
	s_variant_type, s_content, s_stock
	FROM smc_productdetails WHERE s_sku_id = ?`

	rows_details, err := con.Query(qry_details, product.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetProductDetails - qry - " + product.Id + " - " + err.Error()
		res.Data = Product{}
		return res, err, product
	}

	for rows_details.Next() {
		err := rows_details.Scan(&detail.ItemNumber, &detail.Weight, &detail.BuyPrice, &detail.SellPrice, &detail.VKey,
			&detail.VariantType, &detail.Content, &detail.Stock)

		if err != nil {
			fmt.Println(err.Error() + " - " + product.Id)
			res.Status = http.StatusInternalServerError
			res.Message = "GetProductDetails - scn - " + product.Id + " - " + err.Error()
			res.Data = Product{}
			return res, err, product
		}

		product.ProductDetails = append(product.ProductDetails, detail)
	}
	defer rows_details.Close()

	return res, nil, product
}

func GetProductMedia(con *sql.DB, product Product) (Response, error, Product) {
	var res Response
	var media Media

	qry := `SELECT s_item_number, s_file_path
	FROM smc_productmedia WHERE s_sku_id = ?`

	rows_details, err := con.Query(qry, product.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetProductMedia - qry - " + product.Id + " - " + err.Error()
		res.Data = Product{}
		return res, err, product
	}

	for rows_details.Next() {
		err := rows_details.Scan(&media.ItemNumber, &media.FilePath)

		if err != nil {
			fmt.Println(err.Error() + " - " + product.Id)
			res.Status = http.StatusInternalServerError
			res.Message = "GetProductMedia - scn - " + product.Id + " - " + err.Error()
			res.Data = Product{}
			return res, err, product
		}

		product.ProductMedia = append(product.ProductMedia, media)
	}
	defer rows_details.Close()

	return res, nil, product
}

func GetProductDiscount(con *sql.DB, product Product) (Response, error, Product) {
	var res Response
	var discount Discount

	qry := `SELECT s_discount_value
	FROM smc_productdiscount WHERE s_sku_id = ?`

	rows_details, err := con.Query(qry, product.Id)
	if err != nil {
		fmt.Println(err.Error() + " - " + product.Id)
		res.Status = http.StatusInternalServerError
		res.Message = "GetProductDiscount - qry - " + product.Id + " - " + err.Error()
		res.Data = Product{}
		return res, err, product
	}

	for rows_details.Next() {
		err := rows_details.Scan(&discount.Value)

		if err != nil {
			fmt.Println(err.Error() + " - " + product.Id)
			res.Status = http.StatusInternalServerError
			res.Message = "GetProductDiscount - scn - " + product.Id + " - " + err.Error()
			res.Data = Product{}
			return res, err, product
		}

		product.ProductDiscount.Value = discount.Value
	}
	defer rows_details.Close()

	return res, nil, product
}

func GetProductWholesalePrice(con *sql.DB, product Product) (Response, error, Product) {
	var res Response
	var wholesale WholesalePrice

	qry_details := `SELECT s_item_number, s_minimum_amount, s_wholesale_price
	FROM smc_productwholesaleprice WHERE s_sku_id = ?`

	rows_details, err := con.Query(qry_details, product.Id)
	if err != nil {
		fmt.Println(err.Error() + " - " + product.Id)
		res.Status = http.StatusInternalServerError
		res.Message = "GetProductWholesalePrice - qry - " + product.Id + " - " + err.Error()
		res.Data = Product{}
		return res, err, product
	}

	for rows_details.Next() {
		err := rows_details.Scan(&wholesale.ItemNumber, &wholesale.MinimumAmount, &wholesale.Price)

		if err != nil {
			fmt.Println(err.Error() + " - " + product.Id)
			res.Status = http.StatusInternalServerError
			res.Message = "GetProductWholesalePrice - scn - " + product.Id + " - " + err.Error()
			res.Data = Product{}
			return res, err, product
		}

		product.ProductWholesalePrice = append(product.ProductWholesalePrice, wholesale)
	}
	defer rows_details.Close()

	return res, nil, product
}

func CheckProductExist(id string, con *sql.DB) (bool, error) {
	var obj Product

	qry := "SELECT s_sku_id FROM smc_product WHERE s_sku_id = ?"

	err := con.QueryRow(qry, id).Scan(&obj.Id)

	if err == sql.ErrNoRows {
		cerr := errors.New("Product Id '" + id + "' Not Found")
		return false, cerr
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}

func GetTotalProduts() (Response, error) {
	var res Response
	var count int
	con := db.CreateCon()

	qry := "select COUNT(*) as count FROM smc_product"

	err := con.QueryRow(qry).Scan(&count)

	if err != nil {
		// fmt.Println(err.Error() + " - " + product.Id)
		res.Status = http.StatusInternalServerError
		res.Message = "GetTotalProducts - scn  - " + err.Error()
		res.Data = Product{}
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = count
	return res, nil

}

func StoreProductData(product Product, is_mitra bool) (Response, error) {
	var res Response
	con := db.CreateCon()
	// fmt.Println(product)

	exists, err := CheckProductExist(product.Id, con)

	if exists {
		cerr := "Product Id '" + product.Id + "' Already Exist"
		fmt.Println(cerr)
		res.Status = http.StatusInternalServerError
		res.Message = cerr
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = product
		return res, err
	}

	qry := "INSERT INTO smc_product VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	qry_details := "INSERT INTO smc_productdetails VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	qry_media := "INSERT INTO smc_productmedia VALUES(?, ?, ?)"
	qry_discount := "INSERT INTO smc_productdiscount VALUES(?, ?)"
	qry_wholesale := "INSERT INTO smc_productwholesaleprice VALUES(?, ?, ?, ?)"

	//Product Header
	product.Created_at = time.Now().String()
	product.Modified_at = time.Now().String()

	_, err = tx.ExecContext(ctx, qry, product.Id, product.Name, product.CategoryId, product.Description, product.MinOrder,
		product.IsVariant, product.IsDiscount, product.IsReadyStock, product.IsWholesalePrice, product.PublishStatus,
		product.UserId, product.Created_at, product.Modified_at)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Header"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = product
		return res, errors.New(er)
	}
	//End

	//Product Details
	for idx := range product.ProductDetails {
		detail := product.ProductDetails[idx]
		detail.ItemNumber = idx

		_, err := tx.ExecContext(ctx, qry_details, product.Id, detail.ItemNumber, detail.Weight,
			detail.BuyPrice, detail.SellPrice, detail.VKey, detail.VariantType, detail.Content, detail.Stock)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Details"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = product
			return res, errors.New(er)
		}
	}
	//End

	//Product Media
	for idx := range product.ProductMedia {
		media := product.ProductMedia[idx]
		media.ItemNumber = idx

		_, err := tx.ExecContext(ctx, qry_media, product.Id, media.ItemNumber, media.FilePath)
		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Media"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = product
			return res, errors.New(er)
		}
	}
	//End

	//Product Discount
	_, err = tx.ExecContext(ctx, qry_discount, product.Id, product.ProductDiscount.Value)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Discount"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = product
		return res, errors.New(er)
	}
	//End

	//Product Wholesaleprice
	for idx := range product.ProductWholesalePrice {
		wholesale_price := product.ProductWholesalePrice[idx]
		wholesale_price.ItemNumber = idx

		_, err := tx.ExecContext(ctx, qry_wholesale, product.Id, wholesale_price.ItemNumber, wholesale_price.MinimumAmount,
			wholesale_price.Price)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - WholesalePrice"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = product
			return res, errors.New(er)
		}
	}
	//End

	var app = new(Approved)
	app.ObjectId = ObjectId
	app.Id = product.Id
	app.Approved_at = time.Now().String()
	app.UserId = product.UserId
	if is_mitra {
		app.Status = "NOT_VERIFIED"
	} else {
		app.Status = "VERIFIED"
	}
	err = CreateApproved(*app)
	if err != nil {
		tx.Rollback()
		// fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = "Error While Creating Approved" + " - " + err.Error()
		res.Data = app
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		er := err.Error() + " - Commit"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = product
		return res, errors.New(er)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]string{
		"sku_id": product.Id,
	}

	return res, nil
}

func UpdateVerified(product Product) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = product
		return res, err
	}

	qry := `UPDATE smc_product SET s_publish_status = ?, s_modified_at = ? WHERE s_sku_id = ?`
	product.Modified_at = time.Now().String()

	result, err := tx.ExecContext(ctx, qry, product.PublishStatus, product.Modified_at, product.Id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	var app = new(Approved)
	app.ObjectId = ObjectId
	app.Id = product.Id
	app.Approved_at = time.Now().String()
	app.Status = "VERIFIED"
	app.UserId = product.UserId

	err = CreateApproved(*app)
	if err != nil {
		tx.Rollback()
		// fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = "Error While Creating Approved" + " - " + err.Error()
		res.Data = app
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = product
		return res, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int64{
		"rows_affected": affectedRows,
	}

	return res, nil
}

func UpdateProductData(product Product, param_id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	// res.Status = http.StatusOK
	// res.Message = "Test"
	// res.Data = product
	// return res, nil

	exists, err := CheckProductExist(param_id, con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	if !exists {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = product
		return res, err
	}

	qry := `UPDATE smc_product SET s_sku_id = ?, s_name = ?, s_category_id = ?, s_description = ?, s_min_order = ?, 
	s_is_variant = ?, s_is_discount = ?, s_is_ready_stock = ?, s_is_wholesale_price = ?, s_publish_status = ?, 
	s_user_id = ?, s_modified_at = ? WHERE s_sku_id = ?`

	qry_details := "INSERT INTO smc_productdetails VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	qry_media := "INSERT INTO smc_productmedia VALUES(?, ?, ?)"
	qry_discount := "INSERT INTO smc_productdiscount VALUES(?, ?)"
	qry_wholesale := "INSERT INTO smc_productwholesaleprice VALUES(?, ?, ?, ?)"

	product.Modified_at = time.Now().String()

	result, err := tx.ExecContext(ctx, qry, product.Id, product.Name, product.CategoryId, product.Description, product.MinOrder,
		product.IsVariant, product.IsDiscount, product.IsReadyStock, product.IsWholesalePrice, product.PublishStatus,
		product.UserId, product.Modified_at, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = product
		return res, err
	}
	//End

	//Delete The Item Data First
	tx, err = DeleteProductDetailsData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = product
		return res, err
	}

	tx, err = DeleteProductMediaData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	tx, err = DeleteProductDiscountData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	tx, err = DeleteProductWholesalePriceData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}
	//End

	//Product Details
	for idx := range product.ProductDetails {
		detail := product.ProductDetails[idx]
		detail.ItemNumber = idx

		_, err := tx.ExecContext(ctx, qry_details, product.Id, detail.ItemNumber, detail.Weight,
			detail.BuyPrice, detail.SellPrice, detail.VKey, detail.VariantType, detail.Content, detail.Stock)

		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error() + " - Details"
			res.Data = product.ProductDetails

			return res, err
		}
	}
	//End

	//Product Media
	for idx := range product.ProductMedia {
		media := product.ProductMedia[idx]
		media.ItemNumber = idx

		_, err := tx.ExecContext(ctx, qry_media, product.Id, media.ItemNumber, media.FilePath)
		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error() + " - Media"
			res.Data = product
			return res, err
		}
	}
	//End

	//Product Discount
	_, err = tx.ExecContext(ctx, qry_discount, product.Id, product.ProductDiscount.Value)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error() + " - Discount"
		res.Data = product
		return res, err
	}
	//End

	//Product Wholesaleprice
	for idx := range product.ProductWholesalePrice {
		wholesale_price := product.ProductWholesalePrice[idx]
		wholesale_price.ItemNumber = idx

		_, err := tx.ExecContext(ctx, qry_wholesale, product.Id, wholesale_price.ItemNumber, wholesale_price.MinimumAmount,
			wholesale_price.Price)

		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error() + " - Wholesale"
			res.Data = product
			return res, err
		}
	}
	//End

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = product
		return res, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int64{
		"rows_affected": affectedRows,
	}
	return res, nil
}

func DeleteProduct(param_id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckProductExist(param_id, con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	if !exists {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}
	tx, err = DeleteProductHeaderData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	tx, err = DeleteProductDetailsData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	tx, err = DeleteProductMediaData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	tx, err = DeleteProductDiscountData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	tx, err = DeleteProductWholesalePriceData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = nil

	return res, nil
}

func DeleteProductHeaderData(ctx context.Context, tx *sql.Tx, param_id string) (*sql.Tx, error) {
	qry := "DELETE FROM smc_product WHERE s_sku_id = ?"
	_, err := tx.ExecContext(ctx, qry, param_id)

	if err != nil {
		return tx, err
	}

	return tx, nil

}

func DeleteProductDetailsData(ctx context.Context, tx *sql.Tx, param_id string) (*sql.Tx, error) {
	qry := "DELETE FROM smc_productdetails WHERE s_sku_id = ?"

	_, err := tx.ExecContext(ctx, qry, param_id)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func DeleteProductMediaData(ctx context.Context, tx *sql.Tx, param_id string) (*sql.Tx, error) {
	qry := "DELETE FROM smc_productmedia WHERE s_sku_id = ?"

	_, err := tx.ExecContext(ctx, qry, param_id)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func DeleteProductDiscountData(ctx context.Context, tx *sql.Tx, param_id string) (*sql.Tx, error) {
	qry := "DELETE FROM smc_productdiscount WHERE s_sku_id = ?"

	_, err := tx.ExecContext(ctx, qry, param_id)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func DeleteProductWholesalePriceData(ctx context.Context, tx *sql.Tx, param_id string) (*sql.Tx, error) {
	qry := "DELETE FROM smc_productwholesaleprice WHERE s_sku_id = ?"

	_, err := tx.ExecContext(ctx, qry, param_id)
	if err != nil {
		return tx, err
	}

	return tx, nil
}
