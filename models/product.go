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
	Id               string `json:"id"`
	Name             string `json:"name"`
	CategoryId       string `json:"categoryid"`
	CategoryName     string `json:"category_name"`
	Description      string `json:"description"`
	MinOrder         int    `json:"minimal_order"`
	IsVariant        bool   `json:"using_variant"`
	IsDiscount       bool   `json:"using_discount"`
	IsReadyStock     bool   `json:"ready_stock"`
	IsWholesalePrice bool   `json:"using_wholesale_price"`
	// IsPromo               bool             `json:"using_promo"`
	PublishStatus         bool             `json:"publish_status"`
	ProductDetails        []Details        `json:"details"`
	ProductSubCategory    []SubCategory    `json:"sub_category"`
	ProductMedia          []Media          `json:"media"`
	ProductDiscount       Discount         `json:"discount"`
	ProductWholesalePrice []WholesalePrice `json:"wholesale_price"`
	UserId                string           `json:"user_id"`
	Created_at            string           `json:"created_at"`
	Modified_at           string           `json:"modified_at"`
	AppStatus             string           `json:"app_status"`
}

type Details struct {
	ItemNumber  int     `json:"index"`
	Weight      int     `json:"weight"`
	BuyPrice    float64 `json:"buy_price"`
	SellPrice   float64 `json:"sell_price"`
	VKey        int     `json:"v_key"`
	VariantType string  `json:"variant_type"`
	VariantName string  `json:"variant_name"`
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

func FetchAllProductData(is_newest_product bool, category_id string, subCategory_id string, isp bool, user_id string, param_search string) (Response, error) {
	var res Response
	var arrobj []Product
	var product Product
	con := db.CreateCon()
	qry := ""
	if is_newest_product {
		qry = `SELECT A.*, IFNULL(B.s_status, ''), IFNULL(C.s_category_name, '') as 's_category_name' FROM smc_product A
		LEFT JOIN smc_approved B on B.s_id = A.s_sku_id and B.s_object_id = 'PRODUCT'
        LEFT JOIN smc_category C on C.s_category_id = A.s_category_id
		WHERE A.s_created_at > now() - interval 2 month 
		and B.s_status = 'VERIFIED' and A.s_publish_status = 1
		ORDER BY A.s_created_at DESC LIMIT 12`
	} else if isp {
		qry = `SELECT A.*, IFNULL(B.s_status, ''), IFNULL(C.s_category_name, '') as 's_category_name' FROM smc_product A 
		LEFT JOIN smc_approved B on B.s_id = A.s_sku_id and B.s_object_id = 'PRODUCT'
		LEFT JOIN smc_category C on C.s_category_id = A.s_category_id
		WHERE
		B.s_status = 'VERIFIED' and A.s_publish_status = 1
		ORDER BY A.s_created_at DESC`
	} else if user_id != "" {
		qry = `SELECT A.*, IFNULL(B.s_status, ''), IFNULL(C.s_category_name, '') as 's_category_name' FROM smc_product A  
		LEFT JOIN smc_approved B on B.s_id = A.s_sku_id and B.s_object_id = 'PRODUCT'
		LEFT JOIN smc_category C on C.s_category_id = A.s_category_id
		WHERE A.s_user_id = '` + user_id + `'`
	} else if param_search != "" {
		qry = `SELECT A.*, IFNULL(B.s_status, ''), IFNULL(C.s_category_name, '') as 's_category_name' FROM smc_product A 
		LEFT JOIN smc_approved B on B.s_id = A.s_sku_id and B.s_object_id = 'PRODUCT'
		LEFT JOIN smc_category C on C.s_category_id = A.s_category_id
		WHERE
		B.s_status = 'VERIFIED' and A.s_publish_status = 1 and A.s_name like '%` + param_search + `%'`

	} else if category_id != "" {
		if subCategory_id != "" {
			qry = `SELECT A.*, IFNULL(B.s_status, ''), IFNULL(D.s_category_name, '') as 's_category_name' FROM smc_product A
			LEFT JOIN smc_approved B on B.s_id = A.s_sku_id and B.s_object_id = 'PRODUCT'
			LEFT JOIN smc_productsubcategory C on C.s_sku_id = A.s_sku_id
			LEFT JOIN smc_category D on D.s_category_id = A.s_category_id
			WHERE B.s_status = 'VERIFIED' and A.s_publish_status = 1 and A.s_category_id = '` + category_id + `' 
			AND C.s_sub_category_id='` + subCategory_id + `'
			ORDER BY A.s_created_at DESC`
		} else {
			qry = `SELECT A.*, IFNULL(B.s_status, ''), IFNULL(C.s_category_name, '') as 's_category_name' FROM smc_product A
			LEFT JOIN smc_approved B on B.s_id = A.s_sku_id and B.s_object_id = 'PRODUCT'
			LEFT JOIN smc_category C on C.s_category_id = A.s_category_id
			WHERE B.s_status = 'VERIFIED' and A.s_publish_status = 1 and A.s_category_id = '` + category_id + `'
			ORDER BY A.s_created_at DESC`
		}
	} else {
		qry = `SELECT A.*, IFNULL(B.s_status, ''), IFNULL(C.s_category_name, '') as 's_category_name' FROM smc_product A 
		LEFT JOIN smc_approved B on B.s_id = A.s_sku_id and B.s_object_id = 'PRODUCT'
		LEFT JOIN smc_category C on C.s_category_id = A.s_category_id`

	}

	// fmt.Println(is_newest_product)
	// fmt.Println(qry)
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
			&product.Created_at, &product.Modified_at, &product.AppStatus, &product.CategoryName)

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

		res4, err, product := GetProductSubCategory(con, product)
		if err != nil {
			return res4, err
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

	// if category_id != "" {
	// 	qry_params := `s_category_id = '` + category_id + `'`
	// 	ttResponse, tErr := GetTotalProduts(qry_params)

	// 	if tErr != nil {
	// 		return ttResponse, tErr
	// 	}

	// 	res.Status = http.StatusOK
	// 	res.Message = "Success"
	// 	res.Data = map[string]interface{}{
	// 		"data":       arrobj,
	// 		"rows_total": ttResponse.Data,
	// 	}
	// } else {
	// 	res.Status = http.StatusOK
	// 	res.Message = "Success"
	// 	res.Data = arrobj
	// }

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrobj

	return res, nil
}

func GetProductSubCategory(con *sql.DB, product Product) (Response, error, Product) {
	var res Response
	var subCat SubCategory

	qry_details := `SELECT A.s_item_number, A.s_sub_category_id, B.s_sub_category_name FROM smc_productsubcategory A
	LEFT JOIN smc_sub_category B on B.s_sub_category_id = A.s_sub_category_id
	WHERE A.s_sku_id = ?`

	rows_details, err := con.Query(qry_details, product.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetProductSubCategory - qry - " + product.Id + " - " + err.Error()
		res.Data = Product{}
		return res, err, product
	}

	for rows_details.Next() {
		err := rows_details.Scan(&subCat.Index, &subCat.Id, &subCat.Name)

		if err != nil {
			fmt.Println(err.Error() + " - " + product.Id)
			res.Status = http.StatusInternalServerError
			res.Message = "GetProductSubCategory - scn - " + product.Id + " - " + err.Error()
			res.Data = Product{}
			return res, err, product
		}

		product.ProductSubCategory = append(product.ProductSubCategory, subCat)
	}
	defer rows_details.Close()

	return res, nil, product
}

func GetProductById(param_id string) (Product, error) {
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
		return product, err
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
			return product, err
		}

		// fmt.Println(product.Id)
		_, err, prod := GetProductDetails(con, product)
		product.ProductDetails = prod.ProductDetails
		if err != nil {
			return product, err
		}

		_, err, product := GetProductSubCategory(con, product)
		if err != nil {
			return product, err
		}

		_, err, prod = GetProductMedia(con, product)
		product.ProductMedia = prod.ProductMedia
		if err != nil {
			return product, err
		}

		_, err, prod = GetProductDiscount(con, product)
		product.ProductDiscount = prod.ProductDiscount
		if err != nil {
			return product, err
		}

		_, err, prod = GetProductWholesalePrice(con, product)
		product.ProductWholesalePrice = prod.ProductWholesalePrice
		if err != nil {
			return product, err
		}

	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = product

	return product, nil
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

		res4, err, prod := GetProductSubCategory(con, product)
		product.ProductSubCategory = prod.ProductSubCategory
		if err != nil {
			return res4, err
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

		res5, err, prod := GetProductWholesalePrice(con, product)
		product.ProductWholesalePrice = prod.ProductWholesalePrice
		if err != nil {
			return res5, err
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

	qry_details := `SELECT A.s_item_number, A.s_weight, A.s_buy_price, A.s_sell_price, A.s_v_key,
	A.s_variant_type, A.s_content, A.s_stock, IFNULL(B.s_name,'') as 'variant_name'
	FROM smc_productdetails A 
	LEFT JOIN smc_variant B on B.s_variant_id = A.s_variant_type    
	WHERE A.s_sku_id = ?`

	rows_details, err := con.Query(qry_details, product.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetProductDetails - qry - " + product.Id + " - " + err.Error()
		res.Data = Product{}
		return res, err, product
	}

	for rows_details.Next() {
		err := rows_details.Scan(&detail.ItemNumber, &detail.Weight, &detail.BuyPrice, &detail.SellPrice, &detail.VKey,
			&detail.VariantType, &detail.Content, &detail.Stock, &detail.VariantName)

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

func GetTotalProduts(query_params string) (Response, error) {
	var res Response
	var count int
	con := db.CreateCon()
	qry := ""

	if query_params != "" {
		qry = "select COUNT(*) as count FROM smc_product WHERE " + query_params
	} else {
		qry = "select COUNT(*) as count FROM smc_product"
	}

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
	qry_sub_cat := "INSERT INTO smc_productsubcategory VALUES(?, ?, ?)"
	qry_media := "INSERT INTO smc_productmedia VALUES(?, ?, ?)"
	qry_discount := "INSERT INTO smc_productdiscount VALUES(?, ?)"
	qry_wholesale := "INSERT INTO smc_productwholesaleprice VALUES(?, ?, ?, ?)"

	//Product Header
	product.Created_at = time.Now().Format("2006-01-02 15:04:05")
	product.Modified_at = time.Now().Format("2006-01-02 15:04:05")

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

	//Product Sub Category
	for idx := range product.ProductSubCategory {
		subCat := product.ProductSubCategory[idx]
		subCat.Index = idx

		_, err := tx.ExecContext(ctx, qry_sub_cat, product.Id, subCat.Index, subCat.Id)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Sub Category"
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
	app.Approved_at = time.Now().Format("2006-01-02 15:04:05")
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
	product.Modified_at = time.Now().Format("2006-01-02 15:04:05")

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
	app.Approved_at = time.Now().Format("2006-01-02 15:04:05")
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
	qry_sub_cat := "INSERT INTO smc_productsubcategory VALUES(?, ?, ?)"
	// qry_media := "INSERT INTO smc_productmedia VALUES(?, ?, ?)"
	qry_discount := "INSERT INTO smc_productdiscount VALUES(?, ?)"
	qry_wholesale := "INSERT INTO smc_productwholesaleprice VALUES(?, ?, ?, ?)"

	//Product Header
	product.Modified_at = time.Now().Format("2006-01-02 15:04:05")

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

	tx, err = DeleteProductSubCategoryData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = product
		return res, err
	}

	//Off Temporarily
	// tx, err = DeleteProductMediaData(ctx, tx, param_id)

	// if err != nil {
	// 	tx.Rollback()
	// 	fmt.Println(err.Error())
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = err.Error()
	// 	res.Data = map[string]int64{
	// 		"rows_affected": 0,
	// 	}
	// 	return res, err
	// }

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

	//Product Sub Category
	for idx := range product.ProductSubCategory {
		subCat := product.ProductSubCategory[idx]
		subCat.Index = idx

		_, err := tx.ExecContext(ctx, qry_sub_cat, product.Id, subCat.Index, subCat.Id)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Sub Category"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = product
			return res, errors.New(er)
		}
	}
	//End

	//Off Temporarily
	//Product Media
	// for idx := range product.ProductMedia {
	// 	media := product.ProductMedia[idx]
	// 	media.ItemNumber = idx

	// 	_, err := tx.ExecContext(ctx, qry_media, product.Id, media.ItemNumber, media.FilePath)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		fmt.Println(err.Error())
	// 		res.Status = http.StatusInternalServerError
	// 		res.Message = err.Error() + " - Media"
	// 		res.Data = product
	// 		return res, err
	// 	}
	// }
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

	tx, err = DeleteProductSubCategoryData(ctx, tx, param_id)

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

func DeleteProductSubCategoryData(ctx context.Context, tx *sql.Tx, param_id string) (*sql.Tx, error) {
	qry := "DELETE FROM smc_productsubcategory WHERE s_sku_id = ?"

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
