package models

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"time"
)

type Cart struct {
	CartId     string  `json:"cart_id"`
	CustomerId string  `json:"customer_id"`
	ProductId  string  `json:"sku_id"`
	VariantId  string  `json:"variant_id"`
	Qty        float64 `json:"qty"`
	Price      float64 `json:"price"`
	Added_at   string  `json:"added_at"`
}

func GetCartById(customer_id string) (Response, error) {
	var obj Cart
	var res Response

	con := db.CreateCon()

	qry := "SELECT  * FROM smc_cart WHERE s_customer_id = ?"

	rows, err := con.Query(qry, customer_id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Cart{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.CartId, &obj.CustomerId, &obj.ProductId, &obj.VariantId, &obj.Qty, &obj.Price, &obj.Added_at)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Cart{}
			return res, err
		}
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = obj

	return res, nil
}

func UpdateCart(cart Cart) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cart
		return res, err
	}

	qry := `UPDATE smc_cart set s_sku_id = ?, s_variant_id = ?, s_qty = ?, s_price = ? WHERE s_customer_id = ?`

	cart.Added_at = time.Now().String()
	_, err = tx.ExecContext(ctx, qry, cart.ProductId, cart.VariantId, cart.Qty, cart.Price, cart.CustomerId)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cart

		return res, err
	}

	err = tx.Commit()

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cart
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = cart

	return res, nil
}

func CreateCart(cart Cart) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cart
		return res, err
	}

	qry := `INSERT INTO smc_cart VALUES(?, ?, ?, ?, ?, ?, ?)`

	cart.Added_at = time.Now().String()
	_, err = tx.ExecContext(ctx, qry, cart.CartId, cart.CustomerId, cart.ProductId, cart.VariantId, cart.Qty, cart.Price, cart.Added_at)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cart

		return res, err
	}

	err = tx.Commit()

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cart
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = cart

	return res, nil
}

func DeleteCartByCustomerId(customer_id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckCartExist(customer_id, con)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	if !exists {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	qry := "DELETE FROM smc_cart WHERE s_customer_id = ?"

	stmt, err := con.Prepare(qry)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	result, err := stmt.Exec(customer_id)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int64{
		"rows_affected": affectedRows,
	}

	return res, nil
}

func CheckCartExist(id string, con *sql.DB) (bool, error) {
	var cart Cart

	qry := "SELECT s_customer_id FROM smc_cart WHERE s_customer_id = ?"

	err := con.QueryRow(qry, id).Scan(&cart.CustomerId)

	if err == sql.ErrNoRows {
		fmt.Println("There's no Cart On This Account Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}
