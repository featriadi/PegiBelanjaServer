package models

import (
	"context"
	"database/sql"
	"net/http"
	"pb-dev-be/db"
	"strconv"
	"time"
)

type Order struct {
	Id                    string              `json:"order_id"`
	Invoice               string              `json:"invoice_number"`
	OrderAt               string              `json:"order_at"`
	CustomerId            string              `json:"customer_id"`
	Total                 float64             `json:"total_order"`
	Status                string              `json:"status"`
	IsTakeFromStore       bool                `json:"is_take_from_store"`
	IsDropshipper         bool                `json:"is_dropshipper"`
	CourierId             string              `json:"courier_id"`
	BankId                string              `json:"bank_id"`
	Note                  string              `json:"note"`
	IsUsingCoupon         bool                `json:"is_using_coupon"`
	IsUsingPartialBalance bool                `json:"is_using_partial_balance"`
	UserId                string              `json:"user_id"`
	Created_at            string              `json:"created_at"`
	Coupon                OrderCoupon         `json:"coupon"`
	Dropshipper           OrderDropshipper    `json:"dropshipper"`
	PartialBalance        OrderPartialBalance `json:"partial_balance"`
	Item                  OrderProduct        `json:"item"`
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

	gen_id, err := GenerateOrderId(con)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}

	order.Id = strconv.Itoa(gen_id)
	order.OrderAt = time.Now().String()
	order.Created_at = time.Now().String()

	gen_inv, err := GenerateInvoiceNumber(con)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = nil
		return res, err
	}
	order.Invoice = gen_inv

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = order
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = order

	return res, nil
}

func GenerateOrderId(con *sql.DB) (int, error) {
	var order_id int
	var gen_id int

	qry := `SELECT IFNULL(max(order_id),0) as 'order_id' FROM smc_torder`

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
