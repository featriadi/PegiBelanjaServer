package models

import (
	"context"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"pb-dev-be/helpers"
	"time"
)

type User struct {
	Id          string `json:"user_id"`
	Name        string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	UserRole    string `json:"user_role"`
	IsVerified  string `json:"verified"`
	RememberMe  string `json:"remember_me"`
	UserCreated string `json:"user_created"`
	Created_at  string `json:"created_at"`
	Modified_at string `json:"modified_at"`
	LastLogin   string `json:"last_login"`
}

func StoreUserData(user User) (Response, error) {
	var res Response
	con := db.CreateCon()

	hash, _ := helpers.HashPassword(user.Password)
	user.Password = hash

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = user
		return res, err
	}

	qry := `INSERT INTO smc_user (s_name,s_email,s_password,s_userrole_id,is_verified,remember_me,s_user_created,s_created_at,
	s_modified_at, s_lastlogin)
	 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, qry, user.Name, user.Email, user.Password, user.UserRole, user.IsVerified, user.RememberMe,
		user.UserCreated, user.Created_at, user.Modified_at, user.LastLogin)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = user

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = user
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]string{
		"email": user.Email,
	}

	return res, nil
}

func FetchAllUserData() (Response, error) {
	var obj User
	var arrobj []User
	var resp Response

	con := db.CreateCon()

	qry := `SELECT s_user_id,s_name,s_email,s_userrole_id,is_verified,s_remember_me,
	s_user_created,s_created_at,s_modified_at,s_lastlogin FROM smc_user`

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()
		resp.Data = Variant{}
		return resp, nil
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Email, &obj.UserRole, &obj.IsVerified, &obj.RememberMe, &obj.UserCreated,
			&obj.Created_at, &obj.Modified_at, &obj.LastLogin)
		if err != nil {
			fmt.Println(err.Error())
			resp.Status = http.StatusInternalServerError
			resp.Message = err.Error()
			resp.Data = Variant{}
			return resp, nil
		}

		arrobj = append(arrobj, obj)
	}
	defer rows.Close()

	resp.Status = http.StatusOK
	resp.Message = "Success"
	resp.Data = arrobj

	return resp, nil
}

func ShowUserDataById(param_id string) (Response, error) {
	var obj User
	var arrobj []User
	var resp Response

	con := db.CreateCon()

	qry := `SELECT s_user_id,s_name,s_email,s_userrole_id,is_verified,s_remember_me,
	s_user_created,s_created_at,s_modified_at,s_lastlogin FROM smc_user WHERE s_user_id = ?`

	rows, err := con.Query(qry, param_id)

	if err != nil {
		fmt.Println(err.Error())
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()
		resp.Data = obj
		return resp, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Email, &obj.UserRole, &obj.IsVerified, &obj.RememberMe, &obj.UserCreated,
			&obj.Created_at, &obj.Modified_at, &obj.LastLogin)
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

func UpdateLastLogin(user User) (Response, error) {
	var res Response
	con := db.CreateCon()

	user.LastLogin = time.Now().String()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = user.LastLogin
		return res, err
	}

	qry := "UPDATE smc_user SET s_lastlogin = ? WHERE s_email = ?"

	_, err = tx.ExecContext(ctx, qry, user.LastLogin, user.Email)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = user

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = user
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]string{
		"email": user.Email,
		"role":  user.UserRole,
	}

	return res, err
}
