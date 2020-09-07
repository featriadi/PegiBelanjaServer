package models

import (
	"context"
	"net/http"
	"pb-dev-be/db"
)

type Province struct {
	ProvinceId string `json:"province_id"`
	Name       string `json:"province"`
}

func FetchAllProvinceData() (Response, error) {
	var obj Province
	var arrObj []Province
	var res Response

	con := db.CreateCon()

	qry := `SELECT * FROM smc_province ORDER BY s_name ASC`

	rows, err := con.Query(qry)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = Province{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.ProvinceId, &obj.Name)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = Province{}
			return res, err
		}

		arrObj = append(arrObj, obj)
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrObj

	return res, nil
}

func StoreProvinceBulk(province []Province) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = province
		return res, err
	}

	for idx := range province {
		obj := province[idx]

		qry := `INSERT INTO smc_province VALUES (?, ?)`

		_, err = tx.ExecContext(ctx, qry, obj.ProvinceId, obj.Name)

		if err != nil {
			tx.Rollback()
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = province

			return res, err
		}

	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = province
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = province

	return res, nil
}

func StoreProvince(province Province) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = province
		return res, err
	}

	qry := `INSERT INTO smc_province VALUES (?, ?)`

	_, err = tx.ExecContext(ctx, qry, province.ProvinceId, province.Name)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = province

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = province
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = province

	return res, nil
}
