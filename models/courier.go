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

type Courier struct {
	CourierId   string `json:"courier_id"`
	Name        string `json:"name"`
	IsActive    bool   `json:"is_active"`
	UserId      string `json:"user_id"`
	Created_at  string `json:"created_at"`
	Modified_at string `json:"modified_at"`
}

func FetchAllCourier() (Response, error) {
	var obj Courier
	var arrObj []Courier
	var res Response
	con := db.CreateCon()

	qry := `SELECT * FROM smc_courier`

	rows, err := con.Query(qry)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Courier{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.CourierId, &obj.Name, &obj.IsActive, &obj.UserId, &obj.Created_at, &obj.Modified_at)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Courier{}
			return res, err
		}

		arrObj = append(arrObj, obj)
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrObj

	return res, nil
}

func GetCourierById(param_id string) (Courier, error) {
	var obj Courier
	var res Response
	con := db.CreateCon()

	qry := `SELECT * FROM smc_courier WHERE s_courier_id = ?`

	rows, err := con.Query(qry, param_id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Courier{}
		return obj, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.CourierId, &obj.Name, &obj.IsActive, &obj.UserId, &obj.Created_at, &obj.Modified_at)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Courier{}
			return obj, err
		}
	}
	defer rows.Close()

	return obj, nil
}

func StoreCourier(courier Courier) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = courier
		return res, err
	}

	courier.Created_at = time.Now().Format("2006-01-02 15:04:05")
	courier.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	qry := `INSERT INTO smc_courier VALUES (?, ?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, qry, courier.CourierId, courier.Name, courier.IsActive, courier.UserId, courier.Created_at, courier.Modified_at)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = courier

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = courier
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"courier_id": courier.CourierId,
		"created_at": courier.Created_at,
	}

	return res, nil
}

func UpdateCourier(courier Courier, param_id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckCourierExists(param_id, con)

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
		res.Data = courier
		return res, err
	}

	courier.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	qry := `UPDATE smc_courier SET s_courier_id = ?, s_name = ?, s_is_active = ?, s_user_id = ?, s_modified_at = ? WHERE s_courier_id = ?`

	_, err = tx.ExecContext(ctx, qry, courier.CourierId, courier.Name, courier.IsActive, courier.UserId, courier.Created_at, courier.Modified_at, param_id)

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
		"courier_id": courier.CourierId,
		"created_at": courier.Created_at,
	}

	return res, nil
}

func DeleteCourier(Id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckCourierExists(Id, con)

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

	qry := "DELETE FROM smc_courier WHERE s_courier_id = ?"

	_, err = tx.ExecContext(ctx, qry, Id)

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

func CheckCourierExists(id string, con *sql.DB) (bool, error) {
	var obj Customer

	qry := "SELECT s_courier_id FROM smc_courier WHERE s_courier_id = ?"

	err := con.QueryRow(qry, id).Scan(&obj.Id)

	if err == sql.ErrNoRows {
		// fmt.Println()
		newErr := errors.New("Courier Id '" + id + "' Not Found")
		return false, newErr
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}
