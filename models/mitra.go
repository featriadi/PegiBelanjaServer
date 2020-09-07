package models

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"strconv"
	"time"
)

type Mitra struct {
	Id           int       `json:"mitra_id"`
	Name         string    `json:"name"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	NoKTP        string    `json:"no_ktp"`
	StoreName    string    `json:"store_name"`
	MitraAddress []Address `json:"address"`
	UserCreated  string    `json:"user_id"`
	Created_at   string    `json:"created_at"`
	Modified_at  string    `json:"modified_at"`
}

type Address struct {
	ItemNumber  int    `json:"index"`
	Province    string `json:"province"`
	City        string `json:"city"`
	SubDistrict string `json:"sub_district"`
	PostalCode  string `json:"postal_code"`
	Address     string `json:"address"`
}

func FetchAllMitraData() (Response, error) {
	var res Response
	var arrObj []Mitra
	var mitra Mitra

	con := db.CreateCon()

	qry := "SELECT * FROM smc_mitra"

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = mitra
		return res, err
	}

	for rows.Next() {
		err := rows.Scan(&mitra.Id, &mitra.Name, &mitra.PhoneNumber, &mitra.Email, &mitra.StoreName, &mitra.NoKTP,
			&mitra.UserCreated, &mitra.Created_at, &mitra.Modified_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = mitra
			return res, err
		}

		res, err, mitra := GetMitraAddress(con, mitra)

		if err != nil {
			return res, err
		}

		arrObj = append(arrObj, mitra)
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrObj

	return res, nil
}

func GetMitraAddress(con *sql.DB, mitra Mitra) (Response, error, Mitra) {
	var res Response
	var address Address

	qry := `SELECT s_item_number, s_province, s_city, s_sub_district, s_postal_code, s_address FROM smc_mitraaddress
	WHERE s_mitra_id = ?`

	rows, err := con.Query(qry, mitra.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetMitraAddress - qry - " + strconv.Itoa(mitra.Id) + " - " + err.Error()
		res.Data = mitra
		return res, err, mitra
	}

	for rows.Next() {
		err := rows.Scan(&address.ItemNumber, &address.Province, &address.City, &address.SubDistrict,
			&address.PostalCode, &address.Address)

		if err != nil {
			fmt.Println(err.Error() + " - " + strconv.Itoa(mitra.Id))
			res.Status = http.StatusInternalServerError
			res.Message = "GetMitraAddress - scn - " + strconv.Itoa(mitra.Id) + " - " + err.Error()
			res.Data = mitra
			return res, err, mitra
		}

		mitra.MitraAddress = append(mitra.MitraAddress, address)
	}
	defer rows.Close()

	return res, nil, mitra
}

func StoreMitraData(mitra Mitra) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = mitra
		return res, err
	}

	qry := `INSERT INTO smc_mitra VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	qry_address := `INSERT INTO smc_mitraaddress VALUES(?, ?, ?, ?, ?, ?, ?)`

	//Mitra
	gen_id, err := GenerateMitraId(con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = "Failed To Generate Id - " + err.Error()
		res.Data = mitra
		return res, err
	}

	mitra.Id = gen_id
	_, err = tx.ExecContext(ctx, qry, mitra.Id, mitra.Name, mitra.PhoneNumber,
		mitra.Email, mitra.NoKTP, mitra.StoreName, mitra.UserCreated, mitra.Created_at, mitra.Modified_at)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = mitra

		return res, err
	}
	//End

	//Mitra Address
	for idx := range mitra.MitraAddress {
		address := mitra.MitraAddress[idx]
		address.ItemNumber = idx
		_, err := tx.ExecContext(ctx, qry_address, mitra.Id, address.ItemNumber, address.Province, address.City, address.SubDistrict,
			address.PostalCode, address.Address)

		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = mitra

			return res, err
		}
	}
	//End

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = mitra
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int{
		"mitra_id": mitra.Id,
	}

	return res, nil
}

func StoreMitraAndRegisterUser(mitra Mitra, password string) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = mitra
		return res, err
	}

	qry := `INSERT INTO smc_mitra VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	qry_address := `INSERT INTO smc_mitraaddress VALUES(?, ?, ?, ?, ?, ?, ?)`

	//Mitra
	gen_id, err := GenerateMitraId(con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = "Failed To Generate Id - " + err.Error()
		res.Data = mitra
		return res, err
	}

	mitra.Id = gen_id
	_, err = tx.ExecContext(ctx, qry, mitra.Id, mitra.Name, mitra.PhoneNumber,
		mitra.Email, mitra.NoKTP, mitra.StoreName, mitra.UserCreated, mitra.Created_at, mitra.Modified_at)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = mitra

		return res, err
	}
	//End

	//Mitra Address
	for idx := range mitra.MitraAddress {
		address := mitra.MitraAddress[idx]
		address.ItemNumber = idx
		_, err := tx.ExecContext(ctx, qry_address, mitra.Id, address.ItemNumber, address.Province, address.City, address.SubDistrict,
			address.PostalCode, address.Address)

		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = mitra

			return res, err
		}
	}
	//End

	var user User
	user.Name = mitra.Name
	user.Email = mitra.Email
	user.Password = password
	user.UserRole = "MIT"
	user.IsVerified = false
	user.RememberMe = "0"
	user.Created_at = time.Now().String()
	user.Modified_at = time.Now().String()
	user.LastLogin = time.Now().String()

	resUser, errUser := StoreUserData(user)

	if errUser != nil {
		// fmt.Println(err.Error())
		resUser.Status = http.StatusInternalServerError
		resUser.Message = "Error While Creating User :" + errUser.Error()
		return resUser, errUser
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = mitra
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"mitra_id": mitra.Id,
		"email":    mitra.Email,
	}

	return res, nil
}

func UpdateMitra(mitra Mitra, id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	param_id, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = "Error While Converting string parameter to int - " + err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	exists, err := CheckMitraExist(param_id, con)

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

	qry := `UPDATE smc_mitra SET s_name = ?, s_phone_number = ?, s_email = ?, s_no_ktp = ?, s_store_name = ?,
	s_user_id = ?, s_modified_at = ? WHERE s_mitra_id = ?`

	qry_address := `INSERT INTO smc_mitraaddress VALUES(?, ?, ?, ?, ?, ?, ?)`

	mitra.Modified_at = time.Now().String()

	//Mitra Header
	result, err := tx.ExecContext(ctx, qry, mitra.Name, mitra.PhoneNumber, mitra.Email,
		mitra.NoKTP, mitra.StoreName, mitra.UserCreated, mitra.Modified_at, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = mitra
		return res, err
	}
	//End Mitra Header

	//Delete Mitra Header and then insert the data from request
	tx, err = DeleteMitraAddressData(ctx, tx, param_id)

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

	//Mitra Address
	for idx := range mitra.MitraAddress {
		address := mitra.MitraAddress[idx]
		address.ItemNumber = idx
		_, err := tx.ExecContext(ctx, qry_address, mitra.Id, address.ItemNumber, address.Province, address.City, address.SubDistrict,
			address.PostalCode, address.Address)

		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = mitra

			return res, err
		}
	}
	//End

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = mitra
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

func DeleteMitra(id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	param_id, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = "Error While Converting string parameter to int - " + err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	exists, err := CheckMitraExist(param_id, con)

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

	tx, err = DeleteMitraHeaderData(ctx, tx, param_id)

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

	tx, err = DeleteMitraAddressData(ctx, tx, param_id)

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

func DeleteMitraHeaderData(ctx context.Context, tx *sql.Tx, param_id int) (*sql.Tx, error) {
	qry := "DELETE FROM smc_mitra WHERE s_mitra_id = ?"

	_, err := tx.ExecContext(ctx, qry, param_id)

	if err != nil {
		return tx, err
	}

	return tx, nil

}

func DeleteMitraAddressData(ctx context.Context, tx *sql.Tx, param_id int) (*sql.Tx, error) {
	qry := "DELETE FROM smc_mitraaddress WHERE s_mitra_id = ?"

	_, err := tx.ExecContext(ctx, qry, param_id)

	if err != nil {
		return tx, err
	}

	return tx, nil

}

func CheckMitraExist(id int, con *sql.DB) (bool, error) {
	var obj Mitra

	qry := "SELECT s_mitra_id FROM smc_mitra WHERE s_mitra_id = ?"

	err := con.QueryRow(qry, id).Scan(&obj.Id)

	if err == sql.ErrNoRows {
		fmt.Println("Mitra Id '" + strconv.Itoa(id) + "' Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}

func GenerateMitraId(con *sql.DB) (int, error) {
	var mitra_id int
	var gen_id int

	qry := `SELECT IFNULL(max(s_mitra_id),0) as 's_mitra_id' FROM smc_mitra`

	rows, err := con.Query(qry)

	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err := rows.Scan(&gen_id)

		if err != nil {
			return 0, err
		}

		mitra_id = gen_id + 1
	}

	return mitra_id, err
}
