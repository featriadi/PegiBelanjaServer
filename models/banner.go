package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"pb-dev-be/db"
)

type Banner struct {
	Id          string `json:"banner_id"`
	Name        string `json:"name"`
	Link        string `json:"link"`
	ImageData   string `json:"image_data"`
	UserId      string `json:"user_id"`
	Created_at  string `json:"created_at"`
	Modified_at string `json:"modified_at"`
}

func FetchAllBannerData() (Response, error) {
	var obj Banner
	var arrobj []Banner
	var res Response

	con := db.CreateCon()

	qry := "SELECT * FROM smc_banner"

	rows, err := con.Query(qry)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Category{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Link, &obj.ImageData, &obj.UserId, &obj.Created_at, &obj.Modified_at)
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

func CheckBannerExist(id string, con *sql.DB) (bool, error) {
	var obj Banner

	qry := "SELECT * FROM smc_banner WHERE s_banner_id = ?"

	err := con.QueryRow(qry, id).Scan(&obj.Id, &obj.Name, &obj.Link, &obj.ImageData,
		&obj.UserId, &obj.Created_at, &obj.Modified_at,
	)

	if err == sql.ErrNoRows {
		fmt.Println("Banner Not Found")
		return false, err
	}

	if err != nil {
		fmt.Println("Query Error")
		return false, err
	}

	return true, nil
}

func StoreBanner(ban Banner) (Response, error) {
	var res Response
	con := db.CreateCon()

	qry := "INSERT INTO smc_banner VALUES(?, ?, ?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(qry)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = ban
		return res, err
	}

	result, err := stmt.Exec(ban.Id, ban.Name, ban.Link, ban.ImageData, ban.UserId, ban.Created_at, ban.Modified_at)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = ban
		return res, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = ban
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]string{
		"banner_id": ban.Id,
	}

	return res, nil
}

func DeleteBanner(id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckBannerExist(id, con)

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

	qry := "DELETE FROM smc_banner WHERE s_banner_id = ?"

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

func UpdateBanner(ban Banner, param_id string) (Response, error) {
	var res Response
	con := db.CreateCon()

	exists, err := CheckBannerExist(param_id, con)

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

	qry := `UPDATE smc_banner SET s_banner_id = ?, s_banner_name = ?, s_link = ?, s_image_data = ?, s_user_id = ?
	, s_modified_at = ? WHERE s_banner_id = ?`

	stmt, err := con.Prepare(qry)
	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = ban
		return res, err
	}

	result, err := stmt.Exec(ban.Id, ban.Name, ban.Link, ban.ImageData, ban.UserId, ban.Modified_at, param_id)

	if err != nil {
		fmt.Println(err.Error())
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = ban
		return res, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = ban
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int64{
		"rows_affected": affected,
	}

	return res, nil
}
