package models

import (
	"context"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"pb-dev-be/helpers"
	"time"
)

type VerificationCode struct {
	Id         int    `json:"id"`
	Code       int    `json:"code"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Created_at string `json:"created_at"`
}

func StoreVerificationData(ver VerificationCode) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = ver
		return res, err
	}

	exist, err := CheckCustomerExistByEmail(ver.Email, con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Customer{}
		return res, err
	}

	ver.Code = helpers.GenerateVerCode()
	ver.Created_at = time.Now().Format("2006-01-02 15:04:05")

	qry := `INSERT INTO smc_verification (s_code,s_email,s_created_at) VALUES (?, ?, ?)`

	_, err = tx.ExecContext(ctx, qry, ver.Code, ver.Email, ver.Created_at)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = ver

		return res, err
	}

	errE := helpers.SendMailVerification(ver.Code, ver.Name, ver.Email)

	if errE != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = errE.Error()
		res.Data = ver

		return res, errE
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = ver
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"email":   ver.Email,
		"code":    ver.Code,
		"isExist": exist,
	}

	return res, nil
}
