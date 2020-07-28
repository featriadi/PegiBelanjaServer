package models

import (
	"database/sql"
	"fmt"
	"pb-dev-be/db"
	"pb-dev-be/helpers"
)

func CheckLogin(user User) (bool, error, User) {
	var obj User
	var pwd string

	con := db.CreateCon()

	qry := "SELECT * FROM smc_user WHERE s_email = ?"

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
