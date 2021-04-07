package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"pb-dev-be/config"
	"pb-dev-be/db"
	"strconv"
	"time"
)

type MidtransSuccessResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

type MidtransErrorResponse struct {
	ErrorMessages []string `json:"error_messages"`
}

type MidtransOrder struct {
	TransactionDetails MidtransTransactionDetails `json:"transaction_details"`
	ItemDetails        []MidtransItemDetails      `json:"item_details"`
	CustomerDetails    MidtransCustomerDetails    `json:"customer_details"`
	EnabledPayment     []string                   `json:"enabled_payment"`
	// BCAVA              MidtransBCAVA              `json:"bca_va"`
	BNIVA        MidtransBNIVA     `json:"bni_va"`
	Callbacks    MidtransCallbacks `json:"finish"`
	Expiry       MidtransExpiry    `json:"expiry"`
	CustomField1 string            `json:"custom_field1"`
	CustomField2 string            `json:"custom_field2"`
	CustomField3 string            `json:"custom_field3"`
}

type MidtransCustomerDetails struct {
	FirstName       string                  `json:"first_name"`
	LastName        string                  `json:"last_name"`
	Email           string                  `json:"email"`
	Phone           string                  `json:"phone"`
	BillingAddress  MidtransBillingAddress  `json:"billing_addres"`
	ShippingAddress MidtransShippingAddress `json:"shipping_address"`
}

type MidtransExpiry struct {
	StartTime string  `json:"start_time"`
	Unit      string  `json:"unit"`
	Duration  float64 `json:"duration"`
}

type MidtransCallbacks struct {
	Finish string `json:"finish"`
}

type MidtransBNIVA struct {
	VANumber string `json:"va_number"`
}

type MidtransBCAVA struct {
	VANumber       string                `json:"va_number"`
	SubCompanyCode string                `json:"sub_company_code"`
	FreeText       MidtransBCAVAFreeText `json:"free_text"`
}

type MidtransBCAVAFreeText struct {
	Inquiry []MidtransBCAVAFreeTextInquiry `json:"inquiry"`
	Payment []MidtransBCAVAFreeTextPayment `json:"payment"`
}

type MidtransBCAVAFreeTextInquiry struct {
	EN string `json:"en"`
	ID string `json:"id"`
}

type MidtransBCAVAFreeTextPayment struct {
	EN string `json:"en"`
	ID string `json:"id"`
}

type MidtransShippingAddress struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

type MidtransBillingAddress struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

type MidtransItemDetails struct {
	Id           string  `json:"id"`
	Price        float64 `json:"price"`
	Qty          float64 `json:"quantity"`
	Name         string  `json:"name"`
	Brand        string  `json:"brand"`
	Category     string  `json:"category"`
	MerchantName string  `json:"merchant_name"`
}

type MidtransTransactionDetails struct {
	OrderId     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

type Order struct {
	Id                    int                  `json:"order_id"`
	Invoice               string               `json:"invoice_number"`
	OrderAt               string               `json:"order_at"`
	CustomerId            string               `json:"customer_id"`
	CustomerName          string               `json:"customer_name"`
	Tax                   float64              `json:"tax"`
	Total                 float64              `json:"total"`
	Status                string               `json:"status"`
	IsTakeFromStore       bool                 `json:"is_take_from_store"`
	IsDropshipper         bool                 `json:"is_dropshipper"`
	PaymentMethod         string               `json:"payment_method"`
	Note                  string               `json:"note"`
	IsUsingCoupon         bool                 `json:"is_using_coupon"`
	IsUsingPartialBalance bool                 `json:"is_using_partial_balance"`
	UserId                string               `json:"user_id"`
	Created_at            string               `json:"created_at"`
	Item                  []OrderProduct       `json:"item"`
	Delivery              OrderDelivery        `json:"delivery"`
	ShippingAddress       OrderShippingAddress `json:"shipping_address"`
	Coupon                OrderCoupon          `json:"coupon"`
	Dropshipper           OrderDropshipper     `json:"dropshipper"`
	PartialBalance        OrderPartialBalance  `json:"partial_balance"`
	ServerKey             string
	// ServerKey             string               `json:"ServerKey"`
}

type OrderCore struct {
	Id                    int                  `json:"order_id"`
	Invoice               string               `json:"invoice_number"`
	OrderAt               string               `json:"order_at"`
	CustomerId            string               `json:"customer_id"`
	CustomerName          string               `json:"customer_name"`
	Tax                   float64              `json:"tax"`
	Total                 float64              `json:"total"`
	Status                string               `json:"status"`
	IsTakeFromStore       bool                 `json:"is_take_from_store"`
	IsDropshipper         bool                 `json:"is_dropshipper"`
	PaymentMethod         string               `json:"payment_method"`
	Note                  string               `json:"note"`
	IsUsingCoupon         bool                 `json:"is_using_coupon"`
	IsUsingPartialBalance bool                 `json:"is_using_partial_balance"`
	UserId                string               `json:"user_id"`
	Created_at            string               `json:"created_at"`
	Item                  []OrderProduct       `json:"item"`
	Delivery              OrderDelivery        `json:"delivery"`
	ShippingAddress       OrderShippingAddress `json:"shipping_address"`
	Coupon                OrderCoupon          `json:"coupon"`
	Dropshipper           OrderDropshipper     `json:"dropshipper"`
	PartialBalance        OrderPartialBalance  `json:"partial_balance"`
	ServerKey             string
	// ServerKey             string               `json:"ServerKey"`
}

type OrderDelivery struct {
	TimeDeliveryId   string  `json:"time_delivery_id"`
	TimeDeliveryName string  `json:"time_delivery_name"`
	Fee              float64 `json:"delivery_fee"`
	WayBill          string  `json:"waybill"`
}

type OrderCoupon struct {
	CouponCode string `json:"code"`
}

type OrderDropshipper struct {
	ShipperName        string `json:"shipper_name"`
	ShipperPhoneNumber string `json:"phone_number"`
}

type OrderPartialBalance struct {
	PartialBalance float64 `json:"balance"`
}

type OrderProduct struct {
	ItemNumber     int     `json:"index"`
	ProductId      string  `json:"product_id"`
	ProductName    string  `json:"product_name"`
	Qty            float64 `json:"qty"`
	Price          float64 `json:"price"`
	VariantId      string  `json:"variant_id"`
	VariantName    string  `json:"variant_name"`
	VariantContent string  `json:"variant_content"`
}

type OrderShippingAddress struct {
	AddressName     string `json:"name"`
	Recipient       string `json:"recipient"`
	PhoneNumber     string `json:"phone_number"`
	Province        string `json:"province"`
	ProvinceName    string `json:"province_name"`
	City            string `json:"city"`
	CityName        string `json:"city_name"`
	SubDistrict     string `json:"sub_district"`
	SubDistrictName string `json:"sub_district_name"`
	PostalCode      string `json:"postal_code"`
	Address         string `json:"address"`
}

type OrderTracking struct {
	OrderId        string `json:"order_id"`
	ItemNumber     int    `json:"index"`
	TrackingStatus string `json:"tracking_status"`
	Created_at     string `json:"created_at"`
}

func GetOrderStats() (Response, error) {
	var res Response
	var total int
	var totalMonth int

	con := db.CreateCon()

	res1, err, total := GetOrderStatsTotal(con)
	if err != nil {
		return res1, err
	}

	res2, err, totalMonth := GetOrderStatsByMonth(con)
	if err != nil {
		return res2, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"total":       total,
		"total_month": totalMonth,
	}

	return res, nil
}

func GetOrderStatsTotal(con *sql.DB) (Response, error, int) {
	var res Response
	var total int

	qry := `SELECT (sum(s_total) + sum(s_tax)) as 'total' FROM smc_torder where s_status = 'Finish'`

	rows, err := con.Query(qry)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetOrderStatsTotal - qry - " + err.Error()
		res.Data = nil
		return res, err, 0
	}

	for rows.Next() {
		err := rows.Scan(&total)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = "GetOrderStatsByMonth - scn - " + err.Error()
			res.Data = nil
			return res, err, 0
		}
	}
	defer rows.Close()

	return res, nil, total
}

func GetOrderStatsByMonth(con *sql.DB) (Response, error, int) {
	var res Response
	var totalMonth int

	qry := `SELECT IFNULL((sum(s_total) + sum(s_tax)), 0) as 'total' FROM smc_torder 
	WHERE s_status = 'Finish' and month(s_order_at) = ? and year(s_order_at) = ?`

	rows, err := con.Query(qry, int(time.Now().Month()), time.Now().Year())

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetOrderStatsByMonth - qry - " + err.Error()
		res.Data = nil
		return res, err, 0
	}

	for rows.Next() {
		err := rows.Scan(&totalMonth)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = "GetOrderStatsByMonth - scn - " + err.Error()
			res.Data = nil
			return res, err, 0
		}
	}
	defer rows.Close()

	return res, nil, totalMonth
}

func GetTransactionStatus(order_id string) (Response, error) {
	conf := config.GetConfig()
	var res Response

	// url := "https://api.sandbox.midtrans.com/v2/" + order_id + "/status"
	url := "https://api.midtrans.com/v2/" + order_id + "/status"
	// str := base64.StdEncoding.EncodeToString([]byte(conf.MIDTRANS_SERVER_KEY_SANDBOX))
	str := base64.StdEncoding.EncodeToString([]byte(conf.MIDTRANS_SERVER_KEY))

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+str)

	if err != nil {
		// fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = "Request Error " + err.Error()
		res.Data = nil
		return res, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Response Error " + err.Error()
		res.Data = nil
		return res, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Body Read Error " + err.Error()
		res.Data = nil
		return res, err
	}

	var gdat map[string]interface{}
	if err := json.Unmarshal(body, &gdat); err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "UnMarshal Error " + err.Error()
		res.Data = nil
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Request Success"
	res.Data = gdat
	return res, nil
}

func GetOrderTracking(orderId string) (Response, error) {
	var res Response
	var arrobj []OrderTracking
	var track OrderTracking
	con := db.CreateCon()

	qry := `SELECT * FROM smc_ordertracking WHERE s_order_id = ?`

	rows, err := con.Query(qry, orderId)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = track
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&track.OrderId, &track.ItemNumber, &track.TrackingStatus, &track.Created_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = track
			return res, err
		}

		arrobj = append(arrobj, track)

	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrobj

	return res, nil
}

func CreateOrderTracking(track OrderTracking, waybill string) (Response, error) {
	var res Response
	con := db.CreateCon()
	var order Order
	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}
	qry := `INSERT INTO smc_ordertracking VALUES(?, ?, ?, ?)`

	index, err := CheckOrderTrackingIndex(con, track.OrderId)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Check Tracking Index"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = track
		return res, errors.New(er)
	}

	track.ItemNumber = index
	track.Created_at = time.Now().Format("2006-01-02 15:04:05")

	if waybill != "" {
		order.Id, _ = strconv.Atoi(track.OrderId)
		res2, err, order := GetOrderDelivery(con, order)

		if err != nil {
			tx.Rollback()
			er := err.Error()
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = track
			return res2, err
		}
		track.TrackingStatus = track.TrackingStatus + " (" + order.Delivery.TimeDeliveryId + " - " + waybill + ")"
	}

	_, err = tx.ExecContext(ctx, qry, track.OrderId, track.ItemNumber, track.TrackingStatus, track.Created_at)

	if err != nil {
		tx.Rollback()
		er := err.Error()
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = track
		return res, errors.New(er)
	}

	if track.TrackingStatus == "Pesanan Selesai" {
		_, err1 := UpdateOrderStatus(track.OrderId, "Finish")
		if err != nil {
			tx.Rollback()
			er := err1.Error()
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = track
			return res, errors.New(er)
		}
	} else {
		_, err1 := UpdateOrderStatus(track.OrderId, "InProgress")
		if err != nil {
			tx.Rollback()
			er := err1.Error()
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = track
			return res, errors.New(er)
		}
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = track

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = track
		return res, err
	}

	return res, nil
}

func CreateOrder(order Order) (Response, error) {
	var res Response
	conf := config.GetConfig()

	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}

	qry := `INSERT INTO smc_torder VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	qry_item := `INSERT INTO smc_torderproduct VALUES(?, ?, ?, ?, ?, ?, ?)`
	qry_delivery := `INSERT INTO smc_torderdelivery VALUES(?, ?, ?, ?)`
	qry_address := `INSERT INTO smc_tordershippingaddress VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	qry_coupon := `INSERT INTO smc_tordercoupon VALUES(?, ?)`
	qry_dropshipper := `INSERT INTO smc_torderdropshipper VALUES(?, ?, ?)`
	qry_partialbalance := `INSERT INTO smc_torderpartialbalance VALUES(?, ?)`

	gen_id, err := GenerateOrderId(con)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}

	//Order Header
	date := time.Now()
	order.Id = gen_id
	order.OrderAt = date.Format("2006-01-02 15:04:05")
	order.Created_at = date.Format("2006-01-02 15:04:05")
	order.Status = "Accepted"
	gen_inv, err := GenerateInvoiceNumber(con)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}
	order.Invoice = gen_inv

	_, err = tx.ExecContext(ctx, qry, order.Id, order.Invoice, order.OrderAt, order.CustomerId, order.Tax, order.Total,
		order.Status, order.IsTakeFromStore, order.IsDropshipper, order.PaymentMethod, order.Note, order.IsUsingCoupon,
		order.IsUsingPartialBalance, order.UserId, order.Created_at)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Header"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}
	//End Order Header

	//Order Item
	for idx := range order.Item {
		item := order.Item[idx]
		item.ItemNumber = idx

		_, err := tx.ExecContext(ctx, qry_item, order.Id, item.ItemNumber, item.ProductId, item.Qty, item.Price, item.VariantId, item.VariantContent)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Item"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}
	}
	//End Order Item

	//Order Delivery
	_, err = tx.ExecContext(ctx, qry_delivery, order.Id, order.Delivery.TimeDeliveryId, order.Delivery.Fee, order.Delivery.WayBill)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Delivery"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}
	//End Order Delivery

	//Order Address
	_, err = tx.ExecContext(ctx, qry_address, order.Id, order.ShippingAddress.AddressName, order.ShippingAddress.Recipient, order.ShippingAddress.PhoneNumber, order.ShippingAddress.Province,
		order.ShippingAddress.City, order.ShippingAddress.SubDistrict, order.ShippingAddress.PostalCode, order.ShippingAddress.Address)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Address"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}
	//End Order Address

	//Order Coupon
	if order.Coupon.CouponCode != "" {
		_, err = tx.ExecContext(ctx, qry_coupon, order.Id, order.Coupon.CouponCode)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Coupon"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}
	}
	//End Order Coupon

	//Order Dropshipper
	if order.Dropshipper.ShipperName != "" {
		_, err = tx.ExecContext(ctx, qry_dropshipper, order.Id, order.Dropshipper.ShipperName, order.Dropshipper.ShipperPhoneNumber)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Dropshipper"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}
	}
	//End Order Dropshipper

	//Order Partial Balance
	if order.PartialBalance.PartialBalance != 0 {
		_, err = tx.ExecContext(ctx, qry_partialbalance, order.Id, order.PartialBalance.PartialBalance)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Partial Balance"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}
	}
	//End Order Partial Balance

	// if !order.IsTakeFromStore {
	//Prepare Data For Midtrans
	var mdData = new(MidtransOrder)

	//Transaction Details
	mdData.TransactionDetails.OrderId = strconv.Itoa(order.Id)
	mdData.TransactionDetails.GrossAmount = order.Total
	//End Transaction Details

	//Item Details
	//------------------------1
	//Courier
	for idx := range order.Item {
		item := order.Item[idx]
		var mdItemData MidtransItemDetails
		var product Product
		product, err = GetProductById(item.ProductId)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Prepare Midtrans Item Detail (Get Product)"
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}

		category, err := GetCategoryById(product.CategoryId)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Prepare Midtrans Item Detail (Get Category)"
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}

		var productName string
		var variantName = ""
		if product.IsVariant {
			for _, v := range product.ProductDetails {
				if v.VariantType == item.VariantId {
					variantName = v.VariantName
				}
			}

			productName = product.Name + " " + variantName + " " + item.VariantContent
		} else {
			productName = product.Name
		}

		mdItemData.Id = item.ProductId
		mdItemData.Price = item.Price
		mdItemData.Qty = item.Qty
		mdItemData.Name = productName
		mdItemData.Brand = ""
		mdItemData.Category = category.Name
		mdItemData.MerchantName = ""

		mdData.ItemDetails = append(mdData.ItemDetails, mdItemData)
	}
	//End
	//------------------------
	//Delivery Fee
	var mdItemData MidtransItemDetails
	// var courier Courier
	// courier, err = GetCourierById(order.Delivery.CourierId)
	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Prepare Midtrans Item Detail (Get Courier)"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = order
	// 	return res, errors.New(er)
	// }

	mdItemData.Id = order.Delivery.TimeDeliveryId
	mdItemData.Price = order.Delivery.Fee
	mdItemData.Qty = 1
	mdItemData.Name = order.Delivery.TimeDeliveryId + " ~ " + order.Delivery.TimeDeliveryName + " - Delivery"
	mdItemData.Brand = order.Delivery.TimeDeliveryId
	mdItemData.Category = "DELIVERY"
	mdItemData.MerchantName = ""

	mdData.ItemDetails = append(mdData.ItemDetails, mdItemData)

	//End Delivery Fee

	//End Item Details

	//Customer Details
	param_custId, err := strconv.Atoi(order.CustomerId)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Prepare Midtrans Customer Details (Convert Cust Id)"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}

	cust, err := GetCustomerById(param_custId)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Prepare Midtrans Customer Details (Get Customer)"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}

	city, err := GetCityByIdData(order.ShippingAddress.City)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Prepare Midtrans Customer Details (Get City)"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}

	mdData.CustomerDetails.FirstName = cust.Name
	mdData.CustomerDetails.LastName = ""
	mdData.CustomerDetails.Email = cust.Email
	mdData.CustomerDetails.Phone = cust.PhoneNumber

	mdData.CustomerDetails.BillingAddress.FirstName = cust.Name
	mdData.CustomerDetails.BillingAddress.LastName = ""
	mdData.CustomerDetails.BillingAddress.Email = cust.Email
	mdData.CustomerDetails.BillingAddress.Phone = order.ShippingAddress.PhoneNumber
	mdData.CustomerDetails.BillingAddress.Address = order.ShippingAddress.Address
	mdData.CustomerDetails.BillingAddress.City = city.Name
	mdData.CustomerDetails.BillingAddress.PostalCode = order.ShippingAddress.PostalCode
	mdData.CustomerDetails.BillingAddress.CountryCode = "IDN"

	mdData.CustomerDetails.ShippingAddress.FirstName = cust.Name
	mdData.CustomerDetails.ShippingAddress.LastName = ""
	mdData.CustomerDetails.ShippingAddress.Email = cust.Email
	mdData.CustomerDetails.ShippingAddress.Phone = order.ShippingAddress.PhoneNumber
	mdData.CustomerDetails.ShippingAddress.Address = order.ShippingAddress.Address
	mdData.CustomerDetails.ShippingAddress.City = city.Name
	mdData.CustomerDetails.ShippingAddress.PostalCode = order.ShippingAddress.PostalCode
	mdData.CustomerDetails.ShippingAddress.CountryCode = "IDN"
	//End Customer Details

	//Enabled Payment
	mdData.EnabledPayment = []string{"echannel", "bni_va"}
	//End Enabled Payment

	//BCA VA
	var freetextInq MidtransBCAVAFreeTextInquiry
	var freetextPay MidtransBCAVAFreeTextPayment
	// mdData.BCAVA.VANumber = "12345678911"
	// mdData.BCAVA.SubCompanyCode = "00000"

	freetextInq.EN = ""
	freetextInq.ID = ""

	freetextPay.EN = ""
	freetextPay.ID = ""
	// mdData.BCAVA.FreeText.Inquiry = append(mdData.BCAVA.FreeText.Inquiry, freetextInq)
	// mdData.BCAVA.FreeText.Payment = append(mdData.BCAVA.FreeText.Payment, freetextPay)
	//End BCA VA

	//BNI VA
	mdData.BNIVA.VANumber = "12345678"
	//End BNI VA

	//Callbacks
	mdData.Callbacks.Finish = "www.pegiblanja.com"
	//End Callbacks

	//Expiry

	mdData.Expiry.StartTime = date.Format("2006-01-02 15:04:05 +0700")
	mdData.Expiry.Unit = "days"
	mdData.Expiry.Duration = 1
	//End Expiry

	//Custom Fields
	mdData.CustomField1 = order.Invoice
	mdData.CustomField2 = ""
	mdData.CustomField3 = ""
	//End

	//End

	//Server KEY !important
	// str := base64.StdEncoding.EncodeToString([]byte(conf.MIDTRANS_SERVER_KEY_SANDBOX))
	str := base64.StdEncoding.EncodeToString([]byte(conf.MIDTRANS_SERVER_KEY))
	fmt.Println(str)

	jc, err := json.Marshal(mdData)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Prepare Midtrans JSON Convert"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}

	//URL SANDBOX
	// url := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	// URL Production
	// url := "https://app.midtrans.com/snap/v1/transactions"
	url := "https://api.midtrans.com/v2/charge"
	method := "POST"
	payload := bytes.NewBuffer(jc)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Request To Midtrans"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	req.Header.Add("Authorization", "Basic "+str)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(order.ServerKey, "Mid-server-HnkexXzOfhbI6jCPpRDt89Ea")
	// req.Header.Add("Accept", "application/json")
	response, err := client.Do(req)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Response From Midtrans"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Read Response From Midtrans"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	var errRes = new(MidtransErrorResponse)
	var succRes = new(MidtransSuccessResponse)

	err = json.Unmarshal(body, errRes)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Read Error Response From Midtrans"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	err = json.Unmarshal(body, succRes)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Read Success Response From Midtrans"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	if errRes.ErrorMessages != nil {
		tx.Rollback()

		var ErrMessage string
		for idx := range errRes.ErrorMessages {
			obj := errRes.ErrorMessages[idx]
			if idx == 0 {
				ErrMessage = obj
			} else {
				ErrMessage += " - " + obj

			}
		}

		res.Status = http.StatusInternalServerError
		res.Message = "failed"
		res.Data = map[string]interface{}{
			"request_body":  mdData,
			"response_body": errRes,
		}
		return res, errors.New(ErrMessage)
	}

	// url := "https://api.midtrans.com/v2/charge"
	// method := "POST"
	// payload := bytes.NewBuffer(jc)

	// client := &http.Client{}
	// req, err := http.NewRequest(method, url, payload)

	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Request To Midtrans"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	// req.Header.Add("Authorization", "Basic "+str)
	// req.Header.Add("Content-Type", "application/json")
	// response, err := client.Do(req)

	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Response From Midtrans"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	// defer response.Body.Close()
	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Read Response From Midtrans"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	// var errRes = new(MidtransErrorResponse)
	// var succRes = new(MidtransSuccessResponse)

	// err = json.Unmarshal(body, errRes)
	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Read Error Response From Midtrans"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	// err = json.Unmarshal(body, succRes)
	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Read Success Response From Midtrans"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	// if errRes.ErrorMessages != nil {
	// 	tx.Rollback()

	// 	var ErrMessage string
	// 	for idx := range errRes.ErrorMessages {
	// 		obj := errRes.ErrorMessages[idx]
	// 		if idx == 0 {
	// 			ErrMessage = obj
	// 		} else {
	// 			ErrMessage += " - " + obj

	// 		}
	// 	}

	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = "failed"
	// 	res.Data = map[string]interface{}{
	// 		"request_body":  mdData,
	// 		"response_body": errRes,
	// 	}
	// 	return res, errors.New(ErrMessage)
	// }

	err = UpdateCounterCount("inv_order")
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - UpdateCount"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = map[string]interface{}{
		"request_body":  mdData,
		"response_body": succRes,
	}
	// } else {
	// 	err = UpdateCounterCount("inv_order")
	// 	if err != nil {
	// 		tx.Rollback()
	// 		er := err.Error() + " - UpdateCount"
	// 		res.Status = http.StatusInternalServerError
	// 		res.Message = er
	// 		res.Data = order
	// 		return res, errors.New(er)
	// 	}

	// 	res.Status = http.StatusOK
	// 	res.Message = "success"
	// 	res.Data = order
	// }

	// var _orderTacking OrderTracking
	// _orderTacking.OrderId = strconv.Itoa(order.Id)
	// _orderTacking.ItemNumber = 0
	// _orderTacking.Created_at = time.Now().Format("2006-01-02 15:04:05 +0700")
	// _orderTacking.TrackingStatus = "Menunggu Pembayaran"

	// _, err = CreateOrderTracking(_orderTacking, "")

	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - OrderTracking"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = order
		return res, err
	}
	return res, nil
}

func CreateOrderCore(order OrderCore) (Response, error) {
	var res Response
	conf := config.GetConfig()

	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}

	qry := `INSERT INTO smc_torder VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	qry_item := `INSERT INTO smc_torderproduct VALUES(?, ?, ?, ?, ?, ?, ?)`
	qry_delivery := `INSERT INTO smc_torderdelivery VALUES(?, ?, ?, ?)`
	qry_address := `INSERT INTO smc_tordershippingaddress VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	qry_coupon := `INSERT INTO smc_tordercoupon VALUES(?, ?)`
	qry_dropshipper := `INSERT INTO smc_torderdropshipper VALUES(?, ?, ?)`
	qry_partialbalance := `INSERT INTO smc_torderpartialbalance VALUES(?, ?)`

	gen_id, err := GenerateOrderId(con)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}

	//Order Header
	date := time.Now()
	order.Id = gen_id
	order.OrderAt = date.Format("2006-01-02 15:04:05")
	order.Created_at = date.Format("2006-01-02 15:04:05")
	order.Status = "Accepted"
	gen_inv, err := GenerateInvoiceNumber(con)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}
	order.Invoice = gen_inv

	_, err = tx.ExecContext(ctx, qry, order.Id, order.Invoice, order.OrderAt, order.CustomerId, order.Tax, order.Total,
		order.Status, order.IsTakeFromStore, order.IsDropshipper, order.PaymentMethod, order.Note, order.IsUsingCoupon,
		order.IsUsingPartialBalance, order.UserId, order.Created_at)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Header"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}
	//End Order Header

	//Order Item
	for idx := range order.Item {
		item := order.Item[idx]
		item.ItemNumber = idx

		_, err := tx.ExecContext(ctx, qry_item, order.Id, item.ItemNumber, item.ProductId, item.Qty, item.Price, item.VariantId, item.VariantContent)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Item"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}
	}
	//End Order Item

	//Order Delivery
	_, err = tx.ExecContext(ctx, qry_delivery, order.Id, order.Delivery.TimeDeliveryId, order.Delivery.Fee, order.Delivery.WayBill)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Delivery"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}
	//End Order Delivery

	//Order Address
	_, err = tx.ExecContext(ctx, qry_address, order.Id, order.ShippingAddress.AddressName, order.ShippingAddress.Recipient, order.ShippingAddress.PhoneNumber, order.ShippingAddress.Province,
		order.ShippingAddress.City, order.ShippingAddress.SubDistrict, order.ShippingAddress.PostalCode, order.ShippingAddress.Address)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Address"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}
	//End Order Address

	//Order Coupon
	if order.Coupon.CouponCode != "" {
		_, err = tx.ExecContext(ctx, qry_coupon, order.Id, order.Coupon.CouponCode)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Coupon"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}
	}
	//End Order Coupon

	//Order Dropshipper
	if order.Dropshipper.ShipperName != "" {
		_, err = tx.ExecContext(ctx, qry_dropshipper, order.Id, order.Dropshipper.ShipperName, order.Dropshipper.ShipperPhoneNumber)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Dropshipper"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}
	}
	//End Order Dropshipper

	//Order Partial Balance
	if order.PartialBalance.PartialBalance != 0 {
		_, err = tx.ExecContext(ctx, qry_partialbalance, order.Id, order.PartialBalance.PartialBalance)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Partial Balance"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}
	}
	//End Order Partial Balance

	// if !order.IsTakeFromStore {
	//Prepare Data For Midtrans
	var mdData = new(MidtransOrder)

	//Transaction Details
	mdData.TransactionDetails.OrderId = strconv.Itoa(order.Id)
	mdData.TransactionDetails.GrossAmount = order.Total
	//End Transaction Details

	//Item Details
	//------------------------1
	//Courier
	for idx := range order.Item {
		item := order.Item[idx]
		var mdItemData MidtransItemDetails
		var product Product
		product, err = GetProductById(item.ProductId)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Prepare Midtrans Item Detail (Get Product)"
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}

		category, err := GetCategoryById(product.CategoryId)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Prepare Midtrans Item Detail (Get Category)"
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = order
			return res, errors.New(er)
		}

		var productName string
		var variantName = ""
		if product.IsVariant {
			for _, v := range product.ProductDetails {
				if v.VariantType == item.VariantId {
					variantName = v.VariantName
				}
			}

			productName = product.Name + " " + variantName + " " + item.VariantContent
		} else {
			productName = product.Name
		}

		mdItemData.Id = item.ProductId
		mdItemData.Price = item.Price
		mdItemData.Qty = item.Qty
		mdItemData.Name = productName
		mdItemData.Brand = ""
		mdItemData.Category = category.Name
		mdItemData.MerchantName = ""

		mdData.ItemDetails = append(mdData.ItemDetails, mdItemData)
	}
	//End
	//------------------------
	//Delivery Fee
	var mdItemData MidtransItemDetails
	// var courier Courier
	// courier, err = GetCourierById(order.Delivery.CourierId)
	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Prepare Midtrans Item Detail (Get Courier)"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = order
	// 	return res, errors.New(er)
	// }

	mdItemData.Id = order.Delivery.TimeDeliveryId
	mdItemData.Price = order.Delivery.Fee
	mdItemData.Qty = 1
	mdItemData.Name = order.Delivery.TimeDeliveryId + " ~ " + order.Delivery.TimeDeliveryName + " - Delivery"
	mdItemData.Brand = order.Delivery.TimeDeliveryId
	mdItemData.Category = "DELIVERY"
	mdItemData.MerchantName = ""

	mdData.ItemDetails = append(mdData.ItemDetails, mdItemData)

	//End Delivery Fee

	//End Item Details

	//Customer Details
	param_custId, err := strconv.Atoi(order.CustomerId)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Prepare Midtrans Customer Details (Convert Cust Id)"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}

	cust, err := GetCustomerById(param_custId)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Prepare Midtrans Customer Details (Get Customer)"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}

	city, err := GetCityByIdData(order.ShippingAddress.City)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Prepare Midtrans Customer Details (Get City)"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}

	mdData.CustomerDetails.FirstName = cust.Name
	mdData.CustomerDetails.LastName = ""
	mdData.CustomerDetails.Email = cust.Email
	mdData.CustomerDetails.Phone = cust.PhoneNumber

	mdData.CustomerDetails.BillingAddress.FirstName = cust.Name
	mdData.CustomerDetails.BillingAddress.LastName = ""
	mdData.CustomerDetails.BillingAddress.Email = cust.Email
	mdData.CustomerDetails.BillingAddress.Phone = order.ShippingAddress.PhoneNumber
	mdData.CustomerDetails.BillingAddress.Address = order.ShippingAddress.Address
	mdData.CustomerDetails.BillingAddress.City = city.Name
	mdData.CustomerDetails.BillingAddress.PostalCode = order.ShippingAddress.PostalCode
	mdData.CustomerDetails.BillingAddress.CountryCode = "IDN"

	mdData.CustomerDetails.ShippingAddress.FirstName = cust.Name
	mdData.CustomerDetails.ShippingAddress.LastName = ""
	mdData.CustomerDetails.ShippingAddress.Email = cust.Email
	mdData.CustomerDetails.ShippingAddress.Phone = order.ShippingAddress.PhoneNumber
	mdData.CustomerDetails.ShippingAddress.Address = order.ShippingAddress.Address
	mdData.CustomerDetails.ShippingAddress.City = city.Name
	mdData.CustomerDetails.ShippingAddress.PostalCode = order.ShippingAddress.PostalCode
	mdData.CustomerDetails.ShippingAddress.CountryCode = "IDN"
	//End Customer Details

	//Enabled Payment
	mdData.EnabledPayment = []string{"echannel", "bni_va"}
	//End Enabled Payment

	//BCA VA
	var freetextInq MidtransBCAVAFreeTextInquiry
	var freetextPay MidtransBCAVAFreeTextPayment
	// mdData.BCAVA.VANumber = "12345678911"
	// mdData.BCAVA.SubCompanyCode = "00000"

	freetextInq.EN = ""
	freetextInq.ID = ""

	freetextPay.EN = ""
	freetextPay.ID = ""
	// mdData.BCAVA.FreeText.Inquiry = append(mdData.BCAVA.FreeText.Inquiry, freetextInq)
	// mdData.BCAVA.FreeText.Payment = append(mdData.BCAVA.FreeText.Payment, freetextPay)
	//End BCA VA

	//BNI VA
	mdData.BNIVA.VANumber = "12345678"
	//End BNI VA

	//Callbacks
	mdData.Callbacks.Finish = "www.pegiblanja.com"
	//End Callbacks

	//Expiry

	mdData.Expiry.StartTime = date.Format("2006-01-02 15:04:05 +0700")
	mdData.Expiry.Unit = "days"
	mdData.Expiry.Duration = 1
	//End Expiry

	//Custom Fields
	mdData.CustomField1 = order.Invoice
	mdData.CustomField2 = ""
	mdData.CustomField3 = ""
	//End

	//End

	//Server KEY !important
	// str := base64.StdEncoding.EncodeToString([]byte(conf.MIDTRANS_SERVER_KEY_SANDBOX))
	str := base64.StdEncoding.EncodeToString([]byte(conf.MIDTRANS_SERVER_KEY))
	fmt.Println(str)

	jc, err := json.Marshal(mdData)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Prepare Midtrans JSON Convert"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}

	//URL SANDBOX
	// url := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	// URL Production
	// url := "https://app.midtrans.com/snap/v1/transactions"
	url := "https://api.midtrans.com/v2/charge"
	method := "POST"
	payload := bytes.NewBuffer(jc)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Request To Midtrans"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	req.Header.Add("Authorization", "Basic "+str)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(order.ServerKey, "Mid-server-HnkexXzOfhbI6jCPpRDt89Ea")
	response, err := client.Do(req)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Response From Midtrans"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Read Response From Midtrans"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	var errRes = new(MidtransErrorResponse)
	var succRes = new(MidtransSuccessResponse)

	err = json.Unmarshal(body, errRes)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Read Error Response From Midtrans"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	err = json.Unmarshal(body, succRes)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Read Success Response From Midtrans"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	if errRes.ErrorMessages != nil {
		tx.Rollback()

		var ErrMessage string
		for idx := range errRes.ErrorMessages {
			obj := errRes.ErrorMessages[idx]
			if idx == 0 {
				ErrMessage = obj
			} else {
				ErrMessage += " - " + obj

			}
		}

		res.Status = http.StatusInternalServerError
		res.Message = "failed"
		res.Data = map[string]interface{}{
			"request_body":  mdData,
			"response_body": errRes,
		}
		return res, errors.New(ErrMessage)
	}

	// url := "https://api.midtrans.com/v2/charge"
	// method := "POST"
	// payload := bytes.NewBuffer(jc)

	// client := &http.Client{}
	// req, err := http.NewRequest(method, url, payload)

	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Request To Midtrans"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	// req.Header.Add("Authorization", "Basic "+str)
	// req.Header.Add("Content-Type", "application/json")
	// response, err := client.Do(req)

	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Response From Midtrans"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	// defer response.Body.Close()
	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Read Response From Midtrans"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	// var errRes = new(MidtransErrorResponse)
	// var succRes = new(MidtransSuccessResponse)

	// err = json.Unmarshal(body, errRes)
	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Read Error Response From Midtrans"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	// err = json.Unmarshal(body, succRes)
	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - Read Success Response From Midtrans"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	// if errRes.ErrorMessages != nil {
	// 	tx.Rollback()

	// 	var ErrMessage string
	// 	for idx := range errRes.ErrorMessages {
	// 		obj := errRes.ErrorMessages[idx]
	// 		if idx == 0 {
	// 			ErrMessage = obj
	// 		} else {
	// 			ErrMessage += " - " + obj

	// 		}
	// 	}

	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = "failed"
	// 	res.Data = map[string]interface{}{
	// 		"request_body":  mdData,
	// 		"response_body": errRes,
	// 	}
	// 	return res, errors.New(ErrMessage)
	// }

	err = UpdateCounterCount("inv_order")
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - UpdateCount"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = map[string]interface{}{
		"request_body":  mdData,
		"response_body": succRes,
	}
	// } else {
	// 	err = UpdateCounterCount("inv_order")
	// 	if err != nil {
	// 		tx.Rollback()
	// 		er := err.Error() + " - UpdateCount"
	// 		res.Status = http.StatusInternalServerError
	// 		res.Message = er
	// 		res.Data = order
	// 		return res, errors.New(er)
	// 	}

	// 	res.Status = http.StatusOK
	// 	res.Message = "success"
	// 	res.Data = order
	// }

	// var _orderTacking OrderTracking
	// _orderTacking.OrderId = strconv.Itoa(order.Id)
	// _orderTacking.ItemNumber = 0
	// _orderTacking.Created_at = time.Now().Format("2006-01-02 15:04:05 +0700")
	// _orderTacking.TrackingStatus = "Menunggu Pembayaran"

	// _, err = CreateOrderTracking(_orderTacking, "")

	// if err != nil {
	// 	tx.Rollback()
	// 	er := err.Error() + " - OrderTracking"
	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = er
	// 	res.Data = payload
	// 	return res, errors.New(er)
	// }

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = order
		return res, err
	}
	return res, nil
}

func UpdateOrderStatus(order_id string, order_status string) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}

	qry := `UPDATE smc_torder SET s_status = ? WHERE s_order_id = ?`

	_, err = tx.ExecContext(ctx, qry, order_status, order_id)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Header"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = map[string]interface{}{
			"order_id": order_id,
		}
		return res, errors.New(er)
	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]interface{}{
			"order_id": order_id,
		}
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = map[string]interface{}{
		"order_id": order_id,
	}
	return res, nil
}

func UpdateWaybillOrder(order_id int, waybill string) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}

	qry := `UPDATE smc_torderdelivery SET s_waybill = ? WHERE s_order_id = ?`

	_, err = tx.ExecContext(ctx, qry, waybill, order_id)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Header"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = map[string]interface{}{
			"order_id": order_id,
			"waybill":  waybill,
		}
		return res, errors.New(er)
	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]interface{}{
			"order_id": order_id,
			"waybill":  waybill,
		}
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = map[string]interface{}{
		"order_id": order_id,
		"waybill":  waybill,
	}
	return res, nil
}

func GetOrder(isAdmin bool, isMitra bool, userId string, _paramS string, _startMY string, _endMY string) (Response, error) {
	var res Response
	var arrobj []Order
	var order Order
	con := db.CreateCon()
	qry := ""
	if isAdmin {
		qry = `SELECT A.*, IFNULL(CS.s_customer_name,'') as 's_customer_name' FROM smc_torder A
		LEFT JOIN smc_customer CS on CS.s_customer_id = A.s_customer_id
		order by A.s_order_at desc`
	} else if isMitra {
		qry = `SELECT A.*,IFNULL(CS.s_customer_name,'') as 's_customer_name' FROM smc_torder A
		LEFT JOIN smc_torderproduct B on B.s_order_id = A.s_order_id
		LEFT JOIN smc_product C on C.s_sku_id = B.s_sku_id
		LEFT JOIN smc_customer CS on CS.s_customer_id = A.s_customer_id
		WHERE C.s_user_id = '` + userId + `' order by A.s_order_at desc`
	} else if _startMY != "" {
		if _endMY != "" {
			qry = `SELECT A.*, IFNULL(CS.s_customer_name,'') as 's_customer_name' FROM smc_torder A
			LEFT JOIN smc_customer CS on CS.s_customer_id = A.s_customer_id
			WHERE DATE_FORMAT(A.s_order_at, "%m-%Y") BETWEEN '` + _startMY + `' AND '` + _endMY + `'
			order by A.s_order_at desc`
		} else {
			qry = `SELECT A.*, IFNULL(CS.s_customer_name,'') as 's_customer_name' FROM smc_torder A
			LEFT JOIN smc_customer CS on CS.s_customer_id = A.s_customer_id
			WHERE DATE_FORMAT(A.s_order_at, "%m-%Y") = '` + _startMY + `'
			order by A.s_order_at desc`
		}
	} else {
		if userId != "" {
			// fmt.Println(userId)
			if _paramS == "wfp" {
				qry = `SELECT A.*, IFNULL(CS.s_customer_name,'') as 's_customer_name' FROM smc_torder A
				LEFT JOIN smc_customer CS on CS.s_customer_id = A.s_customer_id
                LEFT JOIN smc_ordertracking OT on OT.s_order_id = A.s_order_id
				WHERE A.s_customer_id = '` + userId + `' and OT.s_tracking_status = 'Menunggu Pembayaran'
				and (select MAX(OTT.s_item_number) from smc_ordertracking OTT WHERE OTT.s_order_id = A.s_order_id) < 1
				and CURRENT_TIMESTAMP < DATE_ADD(A.s_order_at, INTERVAL 1 DAY)
                order by A.s_order_at desc`
			} else {
				qry = `SELECT A.*, IFNULL(CS.s_customer_name,'') as 's_customer_name' FROM smc_torder A
				LEFT JOIN smc_customer CS on CS.s_customer_id = A.s_customer_id
				WHERE A.s_customer_id = '` + userId + `' order by A.s_order_at desc`
			}
		}
	}

	rows, err := con.Query(qry)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = order
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&order.Id, &order.Invoice, &order.OrderAt, &order.CustomerId, &order.Tax, &order.Total, &order.Status,
			&order.IsTakeFromStore, &order.IsDropshipper, &order.PaymentMethod, &order.Note, &order.IsUsingCoupon,
			&order.IsUsingPartialBalance, &order.UserId, &order.Created_at, &order.CustomerName)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = order
			return res, err
		}

		res1, err, order := GetOrderCoupon(con, order)
		// order.Coupon = ord.Coupon
		if err != nil {
			return res1, err
		}

		res2, err, order := GetOrderDelivery(con, order)

		if err != nil {
			return res2, err
		}

		// res3, err, order := GetDropShipper(con, order)

		// if err != nil {
		// 	return res3, err
		// }

		// res4, err, order := GetPartialBalance(con, order)

		// if err != nil {
		// 	return res4, err
		// }

		res5, err, order := GetOrderProduct(con, order)

		if err != nil {
			return res5, err
		}

		res6, err, order := GetOrderShippingAddress(con, order)

		if err != nil {
			return res6, err
		}

		arrobj = append(arrobj, order)

	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrobj

	return res, nil
}

func GetOrderShippingAddress(con *sql.DB, order Order) (Response, error, Order) {
	var res Response
	var shipp OrderShippingAddress

	qry := `SELECT A.s_address_name, A.s_recipient, A.s_phone_number, A.s_province, IFNULL(B.s_name,'') as 's_province_name', 
	A.s_city, IFNULL(C.s_name,'') as 's_city_name', A.s_sub_district, IFNULL(D.s_name,'') as 's_sub_district_name',
	A.s_postal_code, A.s_address FROM smc_tordershippingaddress A
	LEFT JOIN smc_province B on B.s_province_id = A.s_province
	LEFT JOIN smc_city C on C.s_city_id = A.s_city
	LEFT JOIN smc_subdistrict D on D.s_subdistrict_id = A.s_sub_district
	WHERE A.s_order_id = ?`

	rows, err := con.Query(qry, order.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetOrderShippingAddress - qry - " + strconv.Itoa(order.Id) + " - " + err.Error()
		res.Data = Order{}
		return res, err, order
	}

	for rows.Next() {
		err := rows.Scan(&shipp.AddressName, &shipp.Recipient, &shipp.PhoneNumber, &shipp.Province, &shipp.ProvinceName,
			&shipp.City, &shipp.CityName, &shipp.SubDistrict, &shipp.SubDistrictName, &shipp.PostalCode, &shipp.Address)

		if err != nil {
			fmt.Println(err.Error() + " - " + strconv.Itoa(order.Id))
			res.Status = http.StatusInternalServerError
			res.Message = "GetOrderShippingAddress - scn - " + strconv.Itoa(order.Id) + " - " + err.Error()
			res.Data = Order{}
			return res, err, order
		}

		order.ShippingAddress.AddressName = shipp.AddressName
		order.ShippingAddress.Recipient = shipp.Recipient
		order.ShippingAddress.PhoneNumber = shipp.PhoneNumber
		order.ShippingAddress.Province = shipp.Province
		order.ShippingAddress.ProvinceName = shipp.ProvinceName
		order.ShippingAddress.City = shipp.City
		order.ShippingAddress.CityName = shipp.CityName
		order.ShippingAddress.SubDistrict = shipp.SubDistrict
		order.ShippingAddress.SubDistrictName = shipp.SubDistrictName
		order.ShippingAddress.PostalCode = shipp.PostalCode
		order.ShippingAddress.Address = shipp.Address
	}
	defer rows.Close()

	return res, nil, order
}

func GetOrderProduct(con *sql.DB, order Order) (Response, error, Order) {
	var res Response
	var product OrderProduct

	qry := `SELECT A.s_item_number,A.s_sku_id,B.s_name,A.s_qty,A.s_price FROM smc_torderproduct A
	INNER JOIN smc_product B on B.s_sku_id = A.s_sku_id
	WHERE A.s_order_id = ?`

	rows, err := con.Query(qry, order.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetOrderProduct - qry - " + strconv.Itoa(order.Id) + " - " + err.Error()
		res.Data = Order{}
		return res, err, order
	}

	for rows.Next() {
		err := rows.Scan(&product.ItemNumber, &product.ProductId, &product.ProductName, &product.Qty, &product.Price)

		if err != nil {
			fmt.Println(err.Error() + " - " + strconv.Itoa(order.Id))
			res.Status = http.StatusInternalServerError
			res.Message = "GetOrderProduct - scn - " + strconv.Itoa(order.Id) + " - " + err.Error()
			res.Data = Order{}
			return res, err, order
		}

		order.Item = append(order.Item, product)
	}
	defer rows.Close()

	return res, nil, order
}

func GetOrderDelivery(con *sql.DB, order Order) (Response, error, Order) {
	var res Response
	var delivery OrderDelivery

	qry := `SELECT s_time_delivery_id,s_delivery_fee,s_waybill FROM smc_torderdelivery WHERE s_order_id = ?`

	rows, err := con.Query(qry, order.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetOrderDelivery - qry - " + strconv.Itoa(order.Id) + " - " + err.Error()
		res.Data = Order{}
		return res, err, order
	}

	for rows.Next() {
		err := rows.Scan(&delivery.TimeDeliveryId, &delivery.Fee, &delivery.WayBill)

		if err != nil {
			fmt.Println(err.Error() + " - " + strconv.Itoa(order.Id))
			res.Status = http.StatusInternalServerError
			res.Message = "GetOrderDelivery - scn - " + strconv.Itoa(order.Id) + " - " + err.Error()
			res.Data = Order{}
			return res, err, order
		}

		order.Delivery.TimeDeliveryId = delivery.TimeDeliveryId
		order.Delivery.Fee = delivery.Fee
		order.Delivery.WayBill = delivery.WayBill
	}
	defer rows.Close()

	return res, nil, order
}

func GetOrderCoupon(con *sql.DB, order Order) (Response, error, Order) {
	var res Response
	var coupon OrderCoupon

	qry := `SELECT s_coupon_code FROM smc_tordercoupon WHERE s_order_id = ?`

	rows, err := con.Query(qry, order.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetOrderCoupon - qry - " + strconv.Itoa(order.Id) + " - " + err.Error()
		res.Data = Order{}
		return res, err, order
	}

	for rows.Next() {
		err := rows.Scan(&coupon.CouponCode)

		if err != nil {
			fmt.Println(err.Error() + " - " + strconv.Itoa(order.Id))
			res.Status = http.StatusInternalServerError
			res.Message = "GetOrderCoupon - scn - " + strconv.Itoa(order.Id) + " - " + err.Error()
			res.Data = Order{}
			return res, err, order
		}

		order.Coupon.CouponCode = coupon.CouponCode
	}
	defer rows.Close()

	return res, nil, order
}

func CheckOrderTrackingIndex(con *sql.DB, order_id string) (int, error) {
	var idx int
	var gen_id int
	var orderId string
	var isExist bool

	qry_chck := "SELECT s_order_id FROM smc_ordertracking WHERE s_order_id = ?"

	err := con.QueryRow(qry_chck, order_id).Scan(&orderId)

	if err == sql.ErrNoRows {
		isExist = false
	} else {
		isExist = true
	}

	qry := `SELECT IFNULL(max(s_item_number),0) as 's_item_number' FROM smc_ordertracking WHERE s_order_id = ?`

	rows, err := con.Query(qry, order_id)

	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err := rows.Scan(&gen_id)

		if err != nil {
			return 0, err
		}

		if isExist {
			idx = gen_id + 1

		}
	}

	return idx, err
}

func GenerateOrderId(con *sql.DB) (int, error) {
	var order_id int
	var gen_id int

	qry := `SELECT IFNULL(max(s_order_id),0) as 's_order_id' FROM smc_torder`

	rows, err := con.Query(qry)

	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err := rows.Scan(&gen_id)

		if err != nil {
			return 0, err
		}

		order_id = gen_id + 1
	}

	return order_id, err
}

func GenerateInvoiceNumber(con *sql.DB) (string, error) {
	counter, err := GetCounterById("inv_order")
	var counter_value int
	if err != nil {
		return "", err
	}

	dateOrder := time.Now().Format("20060102")
	month := intToRoman(int(time.Now().Month()))
	year := time.Now().Format("06")

	if counter.Count == 0 {
		counter_value = 1
	} else {
		counter_value = counter.Count + 1
	}

	inv := "INV/" + dateOrder + "/" + year + "/" + month + "/" + strconv.Itoa(counter_value)

	return inv, nil

}

func intToRoman(num int) string {
	values := []int{
		1000, 900, 500, 400,
		100, 90, 50, 40,
		10, 9, 5, 4, 1,
	}

	symbols := []string{
		"M", "CM", "D", "CD",
		"C", "XC", "L", "XL",
		"X", "IX", "V", "IV",
		"I"}
	roman := ""
	i := 0

	for num > 0 {
		// calculate the number of times this num is completly divisible by values[i]
		// times will only be > 0, when num >= values[i]
		k := num / values[i]
		for j := 0; j < k; j++ {
			//buildup roman numeral
			roman += symbols[i]

			//reduce the value of num.
			num -= values[i]
		}
		i++
	}
	return roman
}
