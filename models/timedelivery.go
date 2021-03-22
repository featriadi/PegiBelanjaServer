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

type TimeDelivery struct {
	Id          string `json:"time_delivery_id"`
	Name        string `json:"time_delivery_name"`
	TimeStart   string `json:"time_start"`
	TimeEnd     string `json:"time_end"`
	UserId      string `json:"user_id"`
	Created_at  string `json:"created_at"`
	Modified_at string `json:"modified_at"`
}

func FetchAllTimeDeliveryData() (Response, error) {
	var obj TimeDelivery
	var arrobj []TimeDelivery
	var resp Response

	con := db.CreateCon()

	qry := "SELECT * FROM smc_delivery_time"

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()
		resp.Data = TimeDelivery{}
		return resp, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.TimeStart, &obj.TimeEnd, &obj.UserId, &obj.Created_at, &obj.Modified_at)
		if err != nil {
			fmt.Println(err.Error())
			resp.Status = http.StatusInternalServerError
			resp.Message = err.Error()
			resp.Data = TimeDelivery{}
			return resp, err
		}

		arrobj = append(arrobj, obj)
	}
	defer rows.Close()

	resp.Status = http.StatusOK
	resp.Message = "Success"
	resp.Data = arrobj

	return resp, nil
}

func GetTimeDeliveryById(param_id string) (TimeDelivery, error) {
	var obj TimeDelivery
	var resp Response

	con := db.CreateCon()

	qry := "SELECT * FROM smc_delivery_time WHERE s_time_delivery_id = ?"

	rows, err := con.Query(qry, param_id)

	if err != nil {
		fmt.Println(err.Error())
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()
		resp.Data = TimeDelivery{}
		return obj, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.TimeStart, &obj.TimeEnd, &obj.UserId, &obj.Created_at, &obj.Modified_at)
		if err != nil {
			fmt.Println(err.Error())
			resp.Status = http.StatusInternalServerError
			resp.Message = err.Error()
			resp.Data = TimeDelivery{}
			return obj, err
		}
	}
	defer rows.Close()

	resp.Status = http.StatusOK
	resp.Message = "Success"
	resp.Data = obj

	return obj, nil
}

func StoreTimeDelivery(td TimeDelivery) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckTimeDelivery(td.Id, con)

	if exists {
		cerr := "Time Delivery With Id '" + td.Id + "' Already Exist"
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
		res.Data = td
		return res, err
	}

	qry := "INSERT INTO smc_delivery_time VALUES (?, ?, ?, ?, ?, ?, ?)"

	//TimeDelivery
	td.Created_at = time.Now().Format("2006-01-02 15:04:05")
	td.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	_, err = tx.ExecContext(ctx, qry, td.Id, td.Name, td.TimeStart, td.TimeEnd, td.UserId, td.Created_at, td.Modified_at)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - TimeDelivery"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = td
		return res, errors.New(er)
	}
	//End TimeDelivery

	err = tx.Commit()
	if err != nil {
		er := err.Error() + " - Commit"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = td
		return res, errors.New(er)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]string{
		"time_delivery": td.Id,
	}

	return res, nil

	// kalau misalnya mau bikin variable baru
	// cat2 := TimeDelivery{}

	// kalau mau gunain variable yang udh ada
	// cat2 := td

	// res.Status = http.StatusOK
	// res.Message = "Success"
	// res.Data = map[string]string{
	// 	"category_id": td.Id,
	// }
}

func UpdateTimeDelivery(td TimeDelivery, param_id string) (Response, error) {
	var res Response

	con := db.CreateCon()

	exists, err := CheckTimeDelivery(param_id, con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	if !exists {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = td
		return res, err
	}

	qry_update := `UPDATE smc_delivery_time SET s_time_delivery_id = ?, s_time_delivery_name = ?, s_time_start = ?, 
	s_time_start = ?, s_user_id = ?, s_modified_at = ? WHERE s_time_delivery_id = ?`

	//TimeDelivery
	td.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	result, err := tx.ExecContext(ctx, qry_update, td.Id, td.Name, td.TimeStart, td.TimeEnd, td.UserId, td.Modified_at, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = td
		return res, err
	}
	//End

	//Delete Sub TimeDelivery Data First
	tx, err = DeleteSubCategoryData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = td
		return res, err
	}
	//End

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = td
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

func DeleteTimeDeliveryData(ctx context.Context, tx *sql.Tx, param_id string) (*sql.Tx, error) {
	qry := "DELETE FROM smc_delivery_time WHERE s_time_delivery_id = ?"
	_, err := tx.ExecContext(ctx, qry, param_id)

	if err != nil {
		return tx, err
	}

	return tx, nil

}

func CheckTimeDelivery(id string, con *sql.DB) (bool, error) {
	var td TimeDelivery

	qry := "SELECT s_time_delivery_id FROM smc_category WHERE s_time_delivery_id = ?"

	err := con.QueryRow(qry, id).Scan(&td.Id)

	if err == sql.ErrNoRows {
		fmt.Println("Time Delivery Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}

func DeleteTimeCategory(id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckTimeDelivery(id, con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	if !exists {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
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

	tx, err = DeleteCategoryData(ctx, tx, id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error() + "TimeDelivery"
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
