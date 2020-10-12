package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"strconv"
	"time"
)

type Customer struct {
	Id             int               `json:"id"`
	Name           string            `json:"name"`
	PhoneNumber    string            `json:"phone_number"`
	Email          string            `json:"email"`
	Point          float64           `json:"point"`
	WalletBallance float64           `json:"wallet_ballance"`
	Status         bool              `json:"status"`
	CustAddress    []CustomerAddress `json:"address"`
	UserCreated    string            `json:"user_id"`
	Created_at     string            `json:"created_at"`
	Modified_at    string            `json:"modified_at"`
}

type CustomerAddress struct {
	ItemNumber      int    `json:"index"`
	AddressName     string `json:"name"`
	Recipient       string `json:"recipient"`
	PhoneNumber     string `json:"phone_number"`
	Province        string `json:"province"`
	ProvinceName    string `json:"province_name"`
	City            string `json:"city"`
	CityName        string `json:"city_name"`
	SubDistrict     string `json:"sub_district"`
	SubDistrictName string `json:"sub_district_name"`
	PostalCode      string `json:"postal_code"`
	Address         string `json:"address"`
	IsMain          bool   `json:"main_address"`
}

func FetchAllCustomerData() (Response, error) {
	var res Response
	var arrObj []Customer
	var cust Customer

	con := db.CreateCon()

	qry := "SELECT * FROM smc_customer"

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Customer{}
		return res, err
	}

	for rows.Next() {
		err := rows.Scan(&cust.Id, &cust.Name, &cust.PhoneNumber, &cust.Email, &cust.Point,
			&cust.WalletBallance, &cust.Status, &cust.UserCreated, &cust.Created_at, &cust.Modified_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Customer{}
			return res, err
		}

		res, err, cust := GetCustomerAddress(con, cust)
		if err != nil {
			return res, err
		}

		arrObj = append(arrObj, cust)
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrObj

	return res, nil

}

func GetCustomerById(param_id int) (Customer, error) {
	var res Response
	var cust Customer

	con := db.CreateCon()

	qry := "SELECT * FROM smc_customer WHERE s_customer_id = ?"

	rows, err := con.Query(qry, param_id)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Customer{}
		return cust, err
	}

	for rows.Next() {
		err := rows.Scan(&cust.Id, &cust.Name, &cust.PhoneNumber, &cust.Email, &cust.Point,
			&cust.WalletBallance, &cust.Status, &cust.UserCreated, &cust.Created_at, &cust.Modified_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Customer{}
			return cust, err
		}

		_, err, cust_a := GetCustomerAddress(con, cust)
		cust.CustAddress = cust_a.CustAddress
		// fmt.Println("Nyampe kemari")
		if err != nil {
			return cust, err
		}
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = cust

	return cust, nil
}

func ShowCustomerById(param_id int) (Response, error) {
	var res Response
	var cust Customer

	con := db.CreateCon()

	qry := "SELECT * FROM smc_customer WHERE s_customer_id = ?"

	rows, err := con.Query(qry, param_id)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Customer{}
		return res, err
	}

	for rows.Next() {
		err := rows.Scan(&cust.Id, &cust.Name, &cust.PhoneNumber, &cust.Email, &cust.Point,
			&cust.WalletBallance, &cust.Status, &cust.UserCreated, &cust.Created_at, &cust.Modified_at)

		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Customer{}
			return res, err
		}

		res, err, cust_a := GetCustomerAddress(con, cust)
		cust.CustAddress = cust_a.CustAddress
		// fmt.Println("Nyampe kemari")
		if err != nil {
			return res, err
		}
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = cust

	return res, nil
}

func GetCustomerAddress(con *sql.DB, cust Customer) (Response, error, Customer) {
	var res Response
	var c_address CustomerAddress

	qry := `SELECT A.s_item_number, A.s_address_name, A.s_recipient, A.s_phone_number, A.s_province, B.s_name as 's_province_name', 
	A.s_city, C.s_name as 's_city_name', A.s_sub_district, D.s_name as 's_sub_district_name',
	A.s_postal_code, A.s_address, A.s_is_main FROM smc_customeraddress A
	LEFT JOIN smc_province B on B.s_province_id = A.s_province
	LEFT JOIN smc_city C on C.s_city_id = A.s_city
	LEFT JOIN smc_subdistrict D on D.s_subdistrict_id = A.s_sub_district
	WHERE A.s_customer_id = ?`

	rows, err := con.Query(qry, cust.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetCustomerAddress - qry - " + strconv.Itoa(cust.Id) + " - " + err.Error()
		res.Data = Customer{}
		return res, err, cust
	}

	for rows.Next() {
		err := rows.Scan(&c_address.ItemNumber, &c_address.AddressName, &c_address.Recipient, &c_address.PhoneNumber, &c_address.Province, &c_address.ProvinceName,
			&c_address.City, &c_address.CityName, &c_address.SubDistrict, &c_address.SubDistrictName, &c_address.PostalCode, &c_address.Address, &c_address.IsMain)

		if err != nil {
			fmt.Println(err.Error() + " - " + strconv.Itoa(cust.Id))
			res.Status = http.StatusInternalServerError
			res.Message = "GetCustomerAddress - scn - " + strconv.Itoa(cust.Id) + " - " + err.Error()
			res.Data = Customer{}
			return res, err, cust
		}

		cust.CustAddress = append(cust.CustAddress, c_address)
	}
	defer rows.Close()

	return res, nil, cust
}

func StoreCustomerData(cust Customer) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cust
		return res, err
	}

	cust.Created_at = time.Now().Format("2006-01-02 15:04:05")
	cust.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	qry := `INSERT INTO smc_customer VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	qry_address := `INSERT INTO smc_customeraddress VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	//Customer Header
	gen_id, err := GenerateCustomerId(con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = "Failed To Generate Id - " + err.Error()
		res.Data = cust
		return res, err
	}

	cust.Id = gen_id

	_, err = tx.ExecContext(ctx, qry, cust.Id, cust.Name, cust.PhoneNumber, cust.Email, cust.Point,
		cust.WalletBallance, cust.Status, cust.UserCreated, cust.Created_at, cust.Modified_at)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cust

		return res, err
	}
	//End Customer

	//Mitra Address
	for idx := range cust.CustAddress {
		address := cust.CustAddress[idx]
		if address.AddressName != "" {
			address.ItemNumber = idx
			_, err := tx.ExecContext(ctx, qry_address, cust.Id, address.ItemNumber, address.AddressName, address.Recipient, address.PhoneNumber,
				address.Province, address.City, address.SubDistrict, address.PostalCode, address.Address, address.IsMain)

			if err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				res.Status = http.StatusInternalServerError
				res.Message = err.Error()
				res.Data = cust

				return res, err
			}
		}
	}
	//End

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cust
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"cust_id": cust.Id,
		"email":   cust.Email,
	}

	return res, nil
}

func StoreCustomerAndUserData(cust Customer, password string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckCustomerExistByEmail(cust.Email, con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	if exists {
		er := errors.New("Email Already Used!")
		res.Status = http.StatusInternalServerError
		res.Message = er.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, er
	}

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cust
		return res, err
	}

	cust.Created_at = time.Now().Format("2006-01-02 15:04:05")
	cust.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	qry := `INSERT INTO smc_customer VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	qry_address := `INSERT INTO smc_customeraddress VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	//Customer Header
	gen_id, err := GenerateCustomerId(con)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = "Failed To Generate Id - " + err.Error()
		res.Data = cust
		return res, err
	}

	cust.Id = gen_id

	_, err = tx.ExecContext(ctx, qry, cust.Id, cust.Name, cust.PhoneNumber, cust.Email, cust.Point,
		cust.WalletBallance, cust.Status, cust.UserCreated, cust.Created_at, cust.Modified_at)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cust

		return res, err
	}
	//End Customer

	//Customer Address
	for idx := range cust.CustAddress {
		address := cust.CustAddress[idx]
		address.ItemNumber = idx
		_, err := tx.ExecContext(ctx, qry_address, cust.Id, address.ItemNumber, address.AddressName, address.Recipient, address.PhoneNumber,
			address.Province, address.City, address.SubDistrict, address.PostalCode, address.Address, address.IsMain)

		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = cust

			return res, err
		}
	}
	//End

	var user User
	user.Name = cust.Name
	user.Email = cust.Email
	user.Password = password
	user.UserRole = "CUST"
	user.IsVerified = true
	user.RememberMe = "0"
	user.Created_at = time.Now().Format("2006-01-02 15:04:05")
	user.Modified_at = time.Now().Format("2006-01-02 15:04:05")
	user.LastLogin = time.Now().Format("2006-01-02 15:04:05")

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
		res.Data = cust
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]interface{}{
		"cust_id": cust.Id,
		"email":   cust.Email,
	}

	return res, nil
}

func UpdateCustomer(cust Customer, id string) (Response, error) {
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

	exists, err := CheckCustomerExist(param_id, con)

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

	qry := `UPDATE smc_customer SET s_customer_name = ?, s_phone_number = ?, s_email = ?, s_status = ?, s_user_id = ?, s_modified_at = ?
	WHERE s_customer_id = ?`

	qry_address := `INSERT INTO smc_customeraddress VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	cust.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	//Customer Header
	result, err := tx.ExecContext(ctx, qry, cust.Name, cust.PhoneNumber, cust.Email, cust.Status, cust.UserCreated, cust.Modified_at, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cust
		return res, err
	}
	//End Customer Header

	//Delete Customer Address data and then insert the data from request
	tx, err = DeleteCustomerAddressData(ctx, tx, param_id)

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

	//Customer Address
	for idx := range cust.CustAddress {
		address := cust.CustAddress[idx]
		address.ItemNumber = idx
		_, err := tx.ExecContext(ctx, qry_address, cust.Id, address.ItemNumber, address.AddressName, address.Recipient, address.PhoneNumber,
			address.Province, address.City, address.SubDistrict, address.PostalCode, address.Address, address.IsMain)

		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = cust

			return res, err
		}
	}
	//End
	//END

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cust
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

func DeleteCustomer(id string) (Response, error) {
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

	exist, err := CheckCustomerExist(param_id, con)

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

	tx, err = DeleteCustomerHeaderData(ctx, tx, param_id)

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

	tx, err = DeleteCustomerAddressData(ctx, tx, param_id)

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

func CheckCustomerExist(id int, con *sql.DB) (bool, error) {
	var obj Customer

	qry := "SELECT s_customer_id FROM smc_customer WHERE s_customer_id = ?"

	err := con.QueryRow(qry, id).Scan(&obj.Id)

	if err == sql.ErrNoRows {
		fmt.Println("Customer Id '" + strconv.Itoa(id) + "' Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}

func CheckCustomerExistByEmail(email string, con *sql.DB) (bool, error) {
	var obj Customer

	qry := "SELECT s_customer_id FROM smc_customer WHERE s_email = ?"

	err := con.QueryRow(qry, email).Scan(&obj.Id)

	if obj.Id != 0 {
		// fmt.Println("Customer Id '" + strconv.Itoa(id) + "' Not Found")
		return true, err
	}

	if err == sql.ErrNoRows {
		// fmt.Println("Customer Id '" + strconv.Itoa(id) + "' Not Found")
		return false, nil
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}

func DeleteCustomerHeaderData(ctx context.Context, tx *sql.Tx, param_id int) (*sql.Tx, error) {
	qry := "DELETE FROM smc_customer WHERE s_customer_id = ?"

	_, err := tx.ExecContext(ctx, qry, param_id)

	if err != nil {
		return tx, err
	}

	return tx, nil

}

func DeleteCustomerAddressData(ctx context.Context, tx *sql.Tx, param_id int) (*sql.Tx, error) {
	qry := "DELETE FROM smc_customeraddress WHERE s_customer_id = ?"

	_, err := tx.ExecContext(ctx, qry, param_id)

	if err != nil {
		return tx, err
	}

	return tx, nil

}

func GenerateCustomerId(con *sql.DB) (int, error) {
	var mitra_id int
	var gen_id int

	qry := `SELECT IFNULL(max(s_customer_id),0) as 's_customer_id' FROM smc_customer`

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
