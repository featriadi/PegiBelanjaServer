package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	BCAVA              MidtransBCAVA              `json:"bca_va"`
	BNIVA              MidtransBNIVA              `json:"bni_va"`
	Callbacks          MidtransCallbacks          `json:"finish"`
	Expiry             MidtransExpiry             `json:"expiry"`
	CustomField1       string                     `json:"custom_field1"`
	CustomField2       string                     `json:"custom_field2"`
	CustomField3       string                     `json:"custom_field3"`
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
	Id                    string               `json:"order_id"`
	Invoice               string               `json:"invoice_number"`
	OrderAt               string               `json:"order_at"`
	CustomerId            string               `json:"customer_id"`
	SubTotal              float64              `json:"sub_total"`
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
}

type OrderDelivery struct {
	CourierId string  `json:"courier_id"`
	Fee       float64 `json:"delivery_fee"`
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
	ItemNumber int     `json:"index"`
	ProductId  string  `json:"product_id"`
	Qty        float64 `json:"qty"`
	Price      float64 `json:"price"`
}

type OrderShippingAddress struct {
	AddressName string `json:"name"`
	Recipient   string `json:"recipient"`
	PhoneNumber string `json:"phone_number"`
	Province    string `json:"province"`
	City        string `json:"city"`
	SubDistrict string `json:"sub_district"`
	PostalCode  string `json:"postal_code"`
	Address     string `json:"address"`
}

func CreateOrder(order Order) (Response, error) {
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

	qry := `INSERT INTO smc_torder VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	qry_item := `INSERT INTO smc_torderproduct VALUES(?, ?, ?, ?, ?)`
	qry_delivery := `INSERT INTO smc_torderdelivery VALUES(?, ?, ?)`
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
	order.Id = strconv.Itoa(gen_id)
	order.OrderAt = time.Now().Format("2006-01-02 15:04:05 +0700")
	order.Created_at = time.Now().String()

	gen_inv, err := GenerateInvoiceNumber(con)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}
	order.Invoice = gen_inv

	_, err = tx.ExecContext(ctx, qry, order.Id, order.Invoice, order.OrderAt, order.CustomerId, order.SubTotal, order.Total,
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

		_, err := tx.ExecContext(ctx, qry_item, order.Id, item.ItemNumber, item.ProductId, item.Qty, item.Price)

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
	_, err = tx.ExecContext(ctx, qry_delivery, order.Id, order.Delivery.CourierId, order.Delivery.Fee)

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
	//End Order Coupon

	//Order Dropshipper
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
	//End Order Dropshipper

	//Order Partial Balance
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
	//End Order Partial Balance

	//Prepare Data For Midtrans
	var mdData = new(MidtransOrder)

	//Transaction Details
	mdData.TransactionDetails.OrderId = order.Invoice
	mdData.TransactionDetails.GrossAmount = order.Total
	//End Transaction Details

	//Item Details
	//------------------------
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

		mdItemData.Id = item.ProductId
		mdItemData.Price = item.Price
		mdItemData.Qty = item.Qty
		mdItemData.Name = product.Name
		mdItemData.Brand = ""
		mdItemData.Category = category.Name
		mdItemData.MerchantName = ""

		mdData.ItemDetails = append(mdData.ItemDetails, mdItemData)
	}
	//End
	//------------------------
	//Delivery Fee
	var mdItemData MidtransItemDetails
	var courier Courier
	courier, err = GetCourierById(order.Delivery.CourierId)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Prepare Midtrans Item Detail (Get Courier)"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}

	mdItemData.Id = courier.CourierId
	mdItemData.Price = order.Delivery.Fee
	mdItemData.Qty = 1
	mdItemData.Name = courier.Name + " - Delivery"
	mdItemData.Brand = courier.Name
	mdItemData.Category = "COURIER"
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
	mdData.EnabledPayment = []string{"bca_va", "bni_va"}
	//End Enabled Payment

	//BCA VA
	var freetextInq MidtransBCAVAFreeTextInquiry
	var freetextPay MidtransBCAVAFreeTextPayment
	mdData.BCAVA.VANumber = "12345678911"
	mdData.BCAVA.SubCompanyCode = "00000"

	freetextInq.EN = ""
	freetextInq.ID = ""

	freetextPay.EN = ""
	freetextPay.ID = ""
	mdData.BCAVA.FreeText.Inquiry = append(mdData.BCAVA.FreeText.Inquiry, freetextInq)
	mdData.BCAVA.FreeText.Payment = append(mdData.BCAVA.FreeText.Payment, freetextPay)
	//End BCA VA

	//BNI VA
	mdData.BNIVA.VANumber = "12345678"
	//End BNI VA

	//Callbacks
	mdData.Callbacks.Finish = "www.pegibelanja.com"
	//End Callbacks

	//Expiry

	mdData.Expiry.StartTime = order.OrderAt
	mdData.Expiry.Unit = "days"
	mdData.Expiry.Duration = 1
	//End Expiry

	//Custom Fields
	mdData.CustomField1 = ""
	mdData.CustomField2 = ""
	mdData.CustomField3 = ""
	//End

	//End

	jc, err := json.Marshal(mdData)
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Prepare Midtrans JSON Convert"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = order
		return res, errors.New(er)
	}

	url := "https://app.sandbox.midtrans.com/snap/v1/transactions"
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

	req.Header.Add("Authorization", "Basic U0ItTWlkLXNlcnZlci1lN0JLZjRIakFxV3JqQm4wY0FaTFNWNXY6")
	req.Header.Add("Content-Type", "application/json")

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

	err = UpdateCounterCount("inv_order")
	if err != nil {
		tx.Rollback()
		er := err.Error() + " - UpdateCount"
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = payload
		return res, errors.New(er)
	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = order
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = map[string]interface{}{
		"request_body":  mdData,
		"response_body": succRes,
	}

	return res, nil
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

	if err != nil {
		return "", err
	}

	dateOrder := time.Now().Format("20060102")
	month := intToRoman(int(time.Now().Month()))
	year := time.Now().Format("06")

	inv := "INV/" + dateOrder + "/" + year + "/" + month + "/" + strconv.Itoa(counter.Count)

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
