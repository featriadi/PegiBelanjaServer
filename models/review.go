package models

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"strconv"
	"time"
)

type Review struct {
	Id         string  `json:"review_id"`
	Content    string  `json:"content"`
	Rating     float64 `json:"rating"`
	ProductId  string  `json:"product_id"`
	CustomerId string  `json:"customer_id"`
	Created_at string  `json:"created_at"`
}

func GetAllReview() (Response, error) {
	var res Response
	var review Review

	con := db.CreateCon()

	qry := "SELECT * FROM smc_review"

	rows, err := con.Query(qry)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Review{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&review.Id, &review.Content, &review.Rating, &review.ProductId, &review.CustomerId, &review.Created_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Review{}
			return res, err
		}
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = review

	return res, nil
}

func GetReviewByProductId(product_id string) (Response, error) {
	var res Response
	var review Review

	con := db.CreateCon()

	_, err := CheckProductReview(product_id, con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = review
		return res, err
	}

	qry := "SELECT * FROM smc_review WHERE s_sku_id = ?"

	rows, err := con.Query(qry, product_id)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = review
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&review.Id, &review.Content, &review.Rating, &review.ProductId, &review.CustomerId, &review.Created_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = review
			return res, err
		}
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = review

	return res, nil
}

func CreateReview(review Review) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = review
		return res, err
	}

	qry := "INSERT INTO smc_review VALUES(?, ?, ?, ?, ?, ?)"

	gen_id, err := GenerateId(con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = "Failed To Generate Id - " + err.Error()
		res.Data = review
		return res, err
	}

	review.Id = strconv.Itoa(gen_id)
	review.Created_at = time.Now().Format("2006-01-02 15:04:05")
	_, err = tx.ExecContext(ctx, qry, review.Id, review.Content, review.Rating, review.ProductId, review.CustomerId, review.Created_at)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = review

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = review
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = review

	return res, nil
}

func GenerateId(con *sql.DB) (int, error) {
	var review_id int
	var gen_id int

	qry := `SELECT IFNULL(max(s_review_id),0) as 's_review_id' FROM smc_review`

	rows, err := con.Query(qry)

	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err := rows.Scan(&gen_id)

		if err != nil {
			return 0, err
		}

		review_id = gen_id + 1
	}

	return review_id, err
}

func CheckProductReview(param_id string, con *sql.DB) (bool, error) {
	var obj Review

	qry := "SELECT s_review_id FROM smc_review WHERE s_sku_id = ?"

	err := con.QueryRow(qry, param_id).Scan(&obj.Id)

	if err == sql.ErrNoRows {
		fmt.Println("Product Review '" + param_id + "' Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}
