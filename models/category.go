package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"pb-dev-be/db"
	"time"
)

type Category struct {
	Id          string        `json:"category_id"`
	Name        string        `json:"name"`
	IconUrl     string        `json:"icon_url"`
	SubCat      []SubCategory `json:"sub_category"`
	UserId      string        `json:"user_id"`
	Created_at  string        `json:"created_at"`
	Modified_at string        `json:"modified_at"`
	// error_res   Response `json:"err"`
}

type SubCategory struct {
	Index   int    `json:"index"`
	Id      string `json:"sub_category_id"`
	Name    string `json:"sub_category_name"`
	IconUrl string `json:"sub_category_icon_url"`
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
		err = rows.Scan(&obj.Id, &obj.Name, &obj.IconUrl, &obj.UserId, &obj.Created_at, &obj.Modified_at)
		if err != nil {
			fmt.Println(err.Error())
			resp.Status = http.StatusInternalServerError
			resp.Message = err.Error()
			resp.Data = Category{}
			return resp, err
		}

		res, err, obj := GetSubCategory(con, obj)
		if err != nil {
			return res, err
		}

		arrobj = append(arrobj, obj)
	}
	defer rows.Close()

	resp.Status = http.StatusOK
	resp.Message = "Success"
	resp.Data = arrobj

	return resp, nil
}

func GetSubCategory(con *sql.DB, category Category) (Response, error, Category) {
	var res Response
	var subCat SubCategory

	qry_details := `SELECT s_item_number, s_sub_category_id, s_sub_category_name, s_icon_url FROM smc_sub_category WHERE s_category_id = ?`

	rows_details, err := con.Query(qry_details, category.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "GetSubCategory - qry - " + category.Id + " - " + err.Error()
		res.Data = Category{}
		return res, err, category
	}

	for rows_details.Next() {
		err := rows_details.Scan(&subCat.Index, &subCat.Id, &subCat.Name, &subCat.IconUrl)

		if err != nil {
			fmt.Println(err.Error() + " - " + category.Id)
			res.Status = http.StatusInternalServerError
			res.Message = "GetSubCategory - scn - " + category.Id + " - " + err.Error()
			res.Data = Category{}
			return res, err, category
		}

		category.SubCat = append(category.SubCat, subCat)
	}
	defer rows_details.Close()

	return res, nil, category
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
		err = rows.Scan(&obj.Id, &obj.Name, &obj.IconUrl, &obj.UserId, &obj.Created_at, &obj.Modified_at)
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

	exists, err := CheckCategoryExist(cat.Id, con)

	if exists {
		cerr := "Cateogry Id '" + cat.Id + "' Already Exist"
		fmt.Println(cerr)
		res.Status = http.StatusInternalServerError
		res.Message = cerr
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
		res.Data = cat
		return res, err
	}

	qry := "INSERT INTO smc_category VALUES (?, ?, ?, ?, ?, ?)"
	qry_item := "INSERT INTO smc_sub_category VALUES (?, ?, ?, ?, ?)"

	//Category
	cat.Created_at = time.Now().Format("2006-01-02 15:04:05")
	cat.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	_, err = tx.ExecContext(ctx, qry, cat.Id, cat.Name, cat.IconUrl, cat.UserId, cat.Created_at, cat.Modified_at)

	if err != nil {
		tx.Rollback()
		er := err.Error() + " - Category"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = cat
		return res, errors.New(er)
	}
	//End Category

	//Sub Category
	for idx := range cat.SubCat {
		subCat := cat.SubCat[idx]
		subCat.Index = idx

		_, err = tx.ExecContext(ctx, qry_item, cat.Id, subCat.Index, subCat.Id, subCat.Name, subCat.IconUrl)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Sub Category"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = cat
			return res, errors.New(er)
		}
	}
	//End

	err = tx.Commit()
	if err != nil {
		er := err.Error() + " - Commit"
		fmt.Println(er)
		res.Status = http.StatusInternalServerError
		res.Message = er
		res.Data = cat
		return res, errors.New(er)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]string{
		"category_id": cat.Id,
	}

	return res, nil

	// kalau misalnya mau bikin variable baru
	// cat2 := Category{}

	// kalau mau gunain variable yang udh ada
	// cat2 := cat

	// res.Status = http.StatusOK
	// res.Message = "Success"
	// res.Data = map[string]string{
	// 	"category_id": cat.Id,
	// }
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

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cat
		return res, err
	}

	qry_update := "UPDATE smc_category SET s_category_id = ?, s_category_name = ?, s_icon_url = ?, s_user_id = ?, s_modified_at = ? WHERE s_category_id = ?"
	qry_item := "INSERT INTO smc_sub_category VALUES (?, ?, ?, ?, ?)"

	//Category
	cat.Modified_at = time.Now().Format("2006-01-02 15:04:05")

	result, err := tx.ExecContext(ctx, qry_update, cat.Id, cat.Name, cat.IconUrl, cat.UserId, cat.Modified_at, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cat
		return res, err
	}
	//End

	//Delete Sub Category Data First
	tx, err = DeleteSubCategoryData(ctx, tx, param_id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cat
		return res, err
	}
	//End

	//Then Insert The New Sub Category Data
	for idx := range cat.SubCat {
		subCat := cat.SubCat[idx]
		subCat.Index = idx

		_, err = tx.ExecContext(ctx, qry_item, cat.Id, subCat.Index, subCat.Id, subCat.Name, subCat.IconUrl)

		if err != nil {
			tx.Rollback()
			er := err.Error() + " - Sub Category"
			fmt.Println(er)
			res.Status = http.StatusInternalServerError
			res.Message = er
			res.Data = cat
			return res, errors.New(er)
		}
	}
	//End

	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = cat
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

func DeleteCategoryData(ctx context.Context, tx *sql.Tx, param_id string) (*sql.Tx, error) {
	qry := "DELETE FROM smc_category WHERE s_category_id = ?"
	_, err := tx.ExecContext(ctx, qry, param_id)

	if err != nil {
		return tx, err
	}

	return tx, nil

}

func DeleteSubCategoryData(ctx context.Context, tx *sql.Tx, param_id string) (*sql.Tx, error) {
	qry := "DELETE FROM smc_sub_category WHERE s_category_id = ?"

	_, err := tx.ExecContext(ctx, qry, param_id)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func CheckCategoryExist(id string, con *sql.DB) (bool, error) {
	var cat Category

	qry := "SELECT s_category_id FROM smc_category WHERE s_category_id = ?"

	err := con.QueryRow(qry, id).Scan(&cat.Id)

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

	tx, err = DeleteCategoryData(ctx, tx, id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error() + "Category"
		res.Data = map[string]int64{
			"rows_affected": 0,
		}
		return res, err
	}

	tx, err = DeleteSubCategoryData(ctx, tx, id)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error() + "Category"
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
