package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"pb-dev-be/db"
)

type Bank struct {
	Id            string `json:"bank_id"`
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	AccountNumber string `json:"account_number"`
	UserId        string `json:"user_id"`
	Created_at    string `json:"created_at"`
	Modified_at   string `json:"modified_at"`
}

func FetchAllBankData() (Response, error) {
	var obj Bank
	var arrobj []Bank
	var res Response

	con := db.CreateCon()

	qry := "SELECT * FROM smc_bankaccount"

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Category{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Owner, &obj.AccountNumber, &obj.UserId, &obj.Created_at, &obj.Modified_at)
		if err != nil {
			fmt.Println(err.Error())
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Bank{}

			return res, err
		}

		arrobj = append(arrobj, obj)
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrobj

	return res, nil
}

func StoreBank(bank Bank) (Response, error) {
	var res Response
	con := db.CreateCon()

	qry := "INSERT INTO smc_bankaccount VALUES (?, ?, ?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(qry)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = bank
		return res, err
	}

	result, err := stmt.Exec(bank.Id, bank.Name, bank.Owner, bank.AccountNumber, bank.UserId, bank.Created_at, bank.Modified_at)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = bank
		return res, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = bank
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]string{
		"bank_id": bank.Id,
	}

	return res, nil
}

func UpdateBank(bank Bank, param_id string) (Response, error) {
	var res Response

	con := db.CreateCon()

	exists, err := CheckBankExist(param_id, con)

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

	qry_update := `UPDATE smc_bankaccount SET s_bank_id = ?, s_bank_name = ?, s_owner = ?,
				s_account_number = ?, s_user_id = ?, s_modified_at = ? WHERE s_bank_id = ?`

	stmt, err := con.Prepare(qry_update)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	result, err := stmt.Exec(bank.Id, bank.Name, bank.Owner, bank.AccountNumber, bank.UserId, bank.Modified_at, param_id)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil

	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int64{
		"rows_affected": affectedRows,
	}

	return res, nil
}

func CheckBankExist(id string, con *sql.DB) (bool, error) {
	var bank Bank

	qry := "SELECT * FROM smc_bankaccount WHERE s_bank_id = ?"

	err := con.QueryRow(qry, id).Scan(&bank.Id, &bank.Name, &bank.Owner, &bank.AccountNumber, &bank.UserId,
		&bank.Created_at, &bank.Modified_at,
	)

	if err == sql.ErrNoRows {
		fmt.Println("Bank Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}

func DeleteBank(id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckBankExist(id, con)

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

	qry := "DELETE FROM smc_bankaccount WHERE s_bank_id = ?"

	stmt, err := con.Prepare(qry)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	result, err := stmt.Exec(id)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int64{
		"rows_affected": affectedRows,
	}

	return res, nil

}
