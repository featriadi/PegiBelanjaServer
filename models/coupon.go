package models

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"time"
)

type Coupon struct {
	Id                    string         `json:"id"`
	Code                  string         `json:"code"`
	CouponDiscount        CouponDiscount `json:"discount"`
	ExpiredDate           string         `json:"expired_date"`
	MinimumAmount         float64        `json:"minimum_amount"`
	MaximalOrder          float64        `json:"maximal_order"`
	Description           string         `json:"description"`
	UsageLimit            float64        `json:"usage_limit"`
	UsageLimitPerCustomer float64        `json:"ulpm"`
	UserCreated           string         `json:"user_id"`
	Created_at            string         `json:"created_at"`
	Modified_at           string         `json:"modified_at"`
}

type CouponDiscount struct {
	DiscountType  string  `json:"type"`
	DiscountValue float64 `json:"value"`
}

func FetchAllCouponData() (Response, error) {
	var res Response
	var arrObj []Coupon
	var coupon Coupon

	con := db.CreateCon()

	qry := "SELECT * FROM smc_coupon"

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Coupon{}
		return res, err
	}

	for rows.Next() {
		err := rows.Scan(&coupon.Id, &coupon.Code, &coupon.CouponDiscount.DiscountType, &coupon.CouponDiscount.DiscountValue,
			&coupon.ExpiredDate, &coupon.MinimumAmount, &coupon.MaximalOrder, &coupon.Description, &coupon.UsageLimit,
			&coupon.UsageLimitPerCustomer, &coupon.UserCreated, &coupon.Created_at, &coupon.Modified_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Coupon{}
			return res, err
		}
		arrObj = append(arrObj, coupon)
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrObj

	return res, nil
}

func StoreCouponData(coupon Coupon) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = coupon
		return res, err
	}

	coupon.Created_at = time.Now().String()
	coupon.Modified_at = time.Now().String()

	qry := "INSERT INTO smc_coupon VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err = tx.ExecContext(ctx, qry, coupon.Id, coupon.Code, coupon.CouponDiscount.DiscountType,
		coupon.CouponDiscount.DiscountValue, coupon.ExpiredDate, coupon.MinimumAmount, coupon.MaximalOrder, coupon.Description,
		coupon.UsageLimit, coupon.UsageLimitPerCustomer, coupon.UserCreated, coupon.Created_at, coupon.Modified_at)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = coupon

		return res, nil
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = coupon
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"id":           coupon.Id,
		"code":         coupon.Code,
		"expired_date": coupon.ExpiredDate,
	}

	return res, nil
}

func UpdateCouponData(coupon Coupon) (Response, error) {
	var res Response
	con := db.CreateCon()

	exist, err := CheckCouponExist(coupon.Id, con)

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
		res.Data = coupon
		return res, err
	}

	coupon.Modified_at = time.Now().String()

	qry := `UPDATE smc_coupon SET s_coupon_code = ?, s_discount_type = ?, s_discount_value = ?, s_expired_date = ?,
	s_minimum_amount = ?, s_maximal_order = ?, s_description = ?, s_usage_limit = ?, s_usage_limit_per_customer = ?,
	s_user_id = ?, s_modified_at = ? WHERE s_coupon_id = ?`

	result, err := tx.ExecContext(ctx, qry, coupon.Code, coupon.CouponDiscount.DiscountType,
		coupon.CouponDiscount.DiscountValue, coupon.ExpiredDate, coupon.MinimumAmount, coupon.MaximalOrder, coupon.Description,
		coupon.UsageLimit, coupon.UsageLimitPerCustomer, coupon.UserCreated, coupon.Modified_at, coupon.Id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = coupon

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = coupon
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

func DeleteCoupon(param_id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exist, err := CheckCouponExist(param_id, con)

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

	qry := "DELETE FROM smc_coupon WHERE s_coupon_id = ?"

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

func CheckCouponExist(param_id string, con *sql.DB) (bool, error) {
	var obj Coupon

	qry := "SELECT s_coupon_id FROM smc_coupon WHERE s_coupon_id = ?"

	err := con.QueryRow(qry, param_id).Scan(&obj.Id)

	if err == sql.ErrNoRows {
		fmt.Println("Coupon Id '" + param_id + "' Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}
