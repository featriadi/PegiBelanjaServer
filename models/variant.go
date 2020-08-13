package models

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"pb-dev-be/db"
)

type Variant struct {
	Id          string `json:"variant_id"`
	Name        string `json:"name"`
	UserId      string `json:"user_id"`
	Created_at  string `json:"created_at"`
	Modified_at string `json:"modified_at"`
	// error_res   Response `json:"err"`
}

func FetchAllVariantData() (Response, error) {
	var obj Variant
	var arrobj []Variant
	var resp Response

	con := db.CreateCon()

	qry := "SELECT * FROM smc_variant"

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()
		resp.Data = obj
		return resp, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.UserId, &obj.Created_at, &obj.Modified_at)
		if err != nil {
			fmt.Println(err.Error())
			resp.Status = http.StatusInternalServerError
			resp.Message = err.Error()
			resp.Data = obj
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

func StoreVariantData(variant Variant) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = variant
		return res, err
	}

	qry := `INSERT INTO smc_variant(s_name, s_user_id, s_created_at, s_modified_at) VALUES(?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, qry, variant.Name, variant.UserId, variant.Created_at, variant.Modified_at)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = variant

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = variant
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = variant

	return res, nil
}

func UpdateVariant(variant Variant) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckVariantExist(variant.Id, con)

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
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	qry := `UPDATE smc_variant SET s_name = ?, s_user_id = ?, s_modified_at = ? WHERE s_variant_id = ?`

	_, err = tx.ExecContext(ctx, qry, variant.Name, variant.UserId, variant.Modified_at, variant.Id)

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
	res.Data = variant

	return res, nil
}

func DeleteVariant(id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckVariantExist(id, con)

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
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	qry := "DELETE FROM smc_variant WHERE s_variant_id = ?"

	_, err = tx.ExecContext(ctx, qry, id)

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

func CheckVariantExist(id string, con *sql.DB) (bool, error) {
	var obj Variant

	qry := "SELECT * FROM smc_variant WHERE s_variant_id = ?"

	err := con.QueryRow(qry, id).Scan(&obj.Id, &obj.Name, &obj.UserId, &obj.Created_at, &obj.Modified_at)

	if err == sql.ErrNoRows {
		fmt.Println("Variant Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}
