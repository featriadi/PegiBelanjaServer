package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"strconv"
	"time"
)

type PBCourier struct {
	Id             int     `json:"id"`
	Province       string  `json:"province"`
	City           string  `json:"city"`
	SubDistrict    string  `json:"sub_district"`
	Price          float64 `json:"price"`
	IsFreeDelivery bool    `json:"is_free_delivery"`
	UserId         string  `json:"user_id"`
	Created_at     string  `json:"created_at"`
	Modified_at    string  `json:"modified_at"`
}

func FetchAllPBCourier() (Response, error) {
	var obj PBCourier
	var arrobj []PBCourier
	var res Response
	con := db.CreateCon()

	qry := `SELECT * FROM smc_pb_courier`

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = PBCourier{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Province, &obj.City, &obj.SubDistrict, &obj.Price, &obj.IsFreeDelivery, &obj.UserId,
			&obj.Created_at, &obj.Modified_at)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = PBCourier{}
			return res, err
		}

		arrobj = append(arrobj, obj)
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrobj

	return res, nil
}

func StorePBCourier(pb_courier PBCourier) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = pb_courier
		return res, err
	}

	pb_courier.Created_at = time.Now().String()
	pb_courier.Modified_at = time.Now().String()

	qry := `INSERT INTO smc_pb_courier VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, qry, pb_courier.Id, pb_courier.Province, pb_courier.City, pb_courier.SubDistrict, pb_courier.IsFreeDelivery, pb_courier.Price,
		pb_courier.UserId, pb_courier.Created_at, pb_courier.Modified_at)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = pb_courier

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = pb_courier
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"province":     pb_courier.Province,
		"city":         pb_courier.City,
		"sub_district": pb_courier.SubDistrict,
		"price":        pb_courier.Price,
	}

	return res, nil
}

func UpdatePBCourier(pb_courier PBCourier) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckPBCourierExists(pb_courier.Id, con)

	if err != nil {
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
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	pb_courier.Modified_at = time.Now().String()

	qry := `UPDATE smc_pb_courier SET s_province = ?, s_city = ?, s_sub_district = ?, s_price = ?, s_is_free_delivery = ?, 
			s_user_id = ?, s_modified_at = ? WHERE s_id = ?`

	_, err = tx.ExecContext(ctx, qry, pb_courier.Province, pb_courier.City, pb_courier.SubDistrict, pb_courier.Price, pb_courier.IsFreeDelivery,
		pb_courier.UserId, pb_courier.Modified_at, pb_courier.Id)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"province":     pb_courier.Province,
		"city":         pb_courier.City,
		"sub_district": pb_courier.SubDistrict,
		"price":        pb_courier.Price,
	}

	return res, nil
}

func DeletePBCourier(id int) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckPBCourierExists(id, con)

	if err != nil {
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
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	qry := "DELETE FROM smc_pb_courier WHERE s_id = ?"

	_, err = tx.ExecContext(ctx, qry, id)

	if err != nil {
		tx.Rollback()
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

func CheckPBCourierExists(id int, con *sql.DB) (bool, error) {
	var obj Customer

	qry := "SELECT s_id FROM smc_pb_courier WHERE s_id = ?"

	err := con.QueryRow(qry, id).Scan(&obj.Id)

	if err == sql.ErrNoRows {
		// fmt.Println()
		newErr := errors.New("PB Courier Id '" + strconv.Itoa(id) + "' Not Found")
		return false, newErr
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}
