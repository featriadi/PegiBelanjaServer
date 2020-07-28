package models

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"time"
)

type DiscountType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	UserCreated string `json:"user_id"`
	Created_at  string `json:"created_at"`
	Modified_at string `json:"modified_at"`
}

func FetchAllDiscountType() (Response, error) {
	var res Response
	var arrObj []DiscountType
	var discType DiscountType

	con := db.CreateCon()

	qry := "SELECT * FROM smc_discounttype"

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = DiscountType{}
		return res, err
	}

	for rows.Next() {
		err := rows.Scan(&discType.Id, &discType.Name, &discType.UserCreated, &discType.Created_at, &discType.Modified_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = DiscountType{}
			return res, err
		}
		arrObj = append(arrObj, discType)
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"discount_type": arrObj,
	}

	return res, nil
}

func StoreDiscountType(discType DiscountType) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = discType
		return res, err
	}

	discType.Created_at = time.Now().String()
	discType.Modified_at = time.Now().String()

	qry := "INSERT INTO smc_discounttype VALUES(?, ?, ?, ?, ?)"

	_, err = tx.ExecContext(ctx, qry, discType.Id, discType.Name, discType.UserCreated, discType.Created_at, discType.Modified_at)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = discType

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = discType
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"id":         discType.Id,
		"created_at": discType.Created_at,
	}

	return res, nil
}

func UpdateDiscountTypeData(discType DiscountType) (Response, error) {
	var res Response
	con := db.CreateCon()

	exist, err := CheckDiscountTypeExist(discType.Id, con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	if !exist {
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
		res.Data = discType
		return res, err
	}

	discType.Modified_at = time.Now().String()

	qry := `UPDATE smc_discounttype SET s_discount_type_name = ?, s_user_id = ?, s_modified_at = ?
			WHERE s_discount_type_id = ?`

	result, err := tx.ExecContext(ctx, qry, discType.Name, discType.UserCreated, discType.Modified_at, discType.Id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = discType

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = discType
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

func DeleteDiscountTypeData(param_id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exist, err := CheckDiscountTypeExist(param_id, con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	if !exist {
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

	qry := "DELETE FROM smc_discounttype WHERE s_discount_type_id = ?"

	_, err = tx.ExecContext(ctx, qry, param_id)

	if err != nil {
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

func CheckDiscountTypeExist(param_id string, con *sql.DB) (bool, error) {
	var obj Coupon

	qry := "SELECT s_discount_type_id FROM smc_discounttype WHERE s_discount_type_id = ?"

	err := con.QueryRow(qry, param_id).Scan(&obj.Id)

	if err == sql.ErrNoRows {
		fmt.Println("Discount Type Id '" + param_id + "' Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}
