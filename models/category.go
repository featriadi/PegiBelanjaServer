package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"pb-dev-be/db"
)

type Category struct {
	Id          string `json:"category_id"`
	Name        string `json:"name"`
	UserId      string `json:"user_id"`
	Created_at  string `json:"created_at"`
	Modified_at string `json:"modified_at"`
	// error_res   Response `json:"err"`
}

func FetchAllCategoryData() (Response, error) {
	var obj Category
	var arrobj []Category
	var resp Response

	con := db.CreateCon()

	qry := "SELECT * FROM smc_category"

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()
		resp.Data = Category{}
		return resp, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.UserId, &obj.Created_at, &obj.Modified_at)
		if err != nil {
			fmt.Println(err.Error())
			resp.Status = http.StatusInternalServerError
			resp.Message = err.Error()
			resp.Data = Category{}
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

func GetCategoryById(param_id string) (Category, error) {
	var obj Category
	var resp Response

	con := db.CreateCon()

	qry := "SELECT * FROM smc_category WHERE s_category_id = ?"

	rows, err := con.Query(qry, param_id)

	if err != nil {
		fmt.Println(err.Error())
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()
		resp.Data = Category{}
		return obj, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.UserId, &obj.Created_at, &obj.Modified_at)
		if err != nil {
			fmt.Println(err.Error())
			resp.Status = http.StatusInternalServerError
			resp.Message = err.Error()
			resp.Data = Category{}
			return obj, err
		}
	}
	defer rows.Close()

	resp.Status = http.StatusOK
	resp.Message = "Success"
	resp.Data = obj

	return obj, nil
}

func StoreCategory(cat Category) (Response, error) {
	var res Response
	con := db.CreateCon()

	qry := "INSERT INTO smc_category VALUES (?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(qry)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cat
		return res, nil
	}

	// fmt.Println(cat)
	result, err := stmt.Exec(cat.Id, cat.Name, cat.UserId, cat.Created_at, cat.Modified_at)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cat
		return res, nil
	}

	_, err = result.LastInsertId()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cat
		return res, nil
	}

	// kalau misalnya mau bikin variable baru
	// cat2 := Category{}

	// kalau mau gunain variable yang udh ada
	// cat2 := cat

	// res.Status = http.StatusOK
	// res.Message = "Success"
	// res.Data = map[string]string{
	// 	"category_id": cat.Id,
	// }

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = cat

	return res, nil
}

func UpdateCategory(cat Category, param_id string) (Response, error) {
	var res Response

	con := db.CreateCon()

	exists, err := CheckCategoryExist(param_id, con)

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

	qry_update := "UPDATE smc_category SET s_category_id = ?, s_category_name = ?, s_user_id = ?, s_modified_at = ? WHERE s_category_id = ?"

	stmt, err := con.Prepare(qry_update)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil
	}

	result, err := stmt.Exec(cat.Id, cat.Name, cat.UserId, cat.Modified_at, param_id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, nil

	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
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

func CheckCategoryExist(id string, con *sql.DB) (bool, error) {
	var cat Category

	qry := "SELECT * FROM smc_category WHERE s_category_id = ?"

	err := con.QueryRow(qry, id).Scan(&cat.Id, &cat.Name, &cat.UserId, &cat.Created_at, &cat.Modified_at)

	if err == sql.ErrNoRows {
		fmt.Println("Category Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}

func DeleteCategory(id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckCategoryExist(id, con)

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

	qry := "DELETE FROM smc_category WHERE s_category_id = ?"

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
