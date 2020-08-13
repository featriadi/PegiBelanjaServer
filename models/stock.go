package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"time"
)

const In = "In"
const Out = "Out"

type Stock struct {
	ProductId  string  `json:"product_id"`
	Status     string  `json:"status"`
	Qty        float64 `json:"qty"`
	UserId     string  `json:"user_id"`
	Created_at string  `json:"created_at"`
}

func CreateStockInAndOut(stock Stock, isIn bool) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = stock
		return res, err
	}

	qry := "INSERT INTO smc_hstock VALUES(?, ?, ?, ?, ?)"

	if isIn {
		stock.Status = In
	} else {
		stock.Status = Out
	}

	stock.Created_at = time.Now().String()

	_, err = tx.ExecContext(ctx, qry, stock.ProductId, stock.Status, stock.Qty, stock.UserId, stock.Created_at)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = stock
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		er := err.Error() + " - Commit"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = stock
		return res, errors.New(er)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = stock

	return res, nil
}
