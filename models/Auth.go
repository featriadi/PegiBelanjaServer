package models

import (
	"database/sql"
	"fmt"
	"pb-dev-be/db"
	"pb-dev-be/helpers"
)

func CheckLogin(user User, is_customer bool, is_mitra bool, is_admin bool) (bool, error, User) {
	var obj User
	var pwd string

	con := db.CreateCon()
	qry := ""

	if is_customer {
		qry = `SELECT B.s_customer_id,A.s_name,A.s_email,A.s_password,A.s_userrole_id,A.is_verified,
		A.s_remember_me,A.s_user_created, A.s_created_at,A.s_modified_at,A.s_lastlogin 
		FROM smc_user A LEFT JOIN smc_customer B on B.s_email = A.s_email WHERE A.s_email = ?`
	} else if is_mitra {
		qry = `SELECT B.s_mitra_id,A.s_name,A.s_email,A.s_password,A.s_userrole_id,A.is_verified,
		A.s_remember_me,A.s_user_created, A.s_created_at,A.s_modified_at,A.s_lastlogin 
		FROM smc_user A LEFT JOIN smc_mitra B on B.s_email = A.s_email WHERE A.s_email = ?`
	} else if is_admin {
		qry = `SELECT * FROM smc_user where s_email = ? and s_userrole_id = 'SADM'`
	} else {
		qry = "SELECT * FROM smc_user WHERE s_email = ?"
	}

	err := con.QueryRow(qry, user.Email).Scan(&obj.Id, &obj.Name, &obj.Email, &pwd,
		&obj.UserRole, &obj.IsVerified, &obj.RememberMe, &obj.UserCreated, &obj.Created_at, &obj.Modified_at, &obj.LastLogin,
	)

	if err == sql.ErrNoRows {
		fmt.Println("User Not Found")
		return false, err, obj
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err, obj
	}

	match, err := helpers.CheckPasswordHash(user.Password, pwd)

	if !match {
		fmt.Println("Password not match")
		return false, err, obj
	}

	return true, nil, obj
}
