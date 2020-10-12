package models

import (
	"context"
	"net/http"
	"pb-dev-be/db"
)

type SubDistrict struct {
	SubDistrictId string `json:"subdistrict_id"`
	Name          string `json:"subdistrict_name"`
	CityId        string `json:"city_id"`
	ProvinceId    string `json:"province_id"`
	Type          string `json:"type"`
}

func FetchAllSubDistrictData() (Response, error) {
	var obj SubDistrict
	var arrObj []SubDistrict
	var res Response

	con := db.CreateCon()

	qry := `SELECT * FROM smc_subsdistrict ORDER BY s_name ASC`

	rows, err := con.Query(qry)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = SubDistrict{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.SubDistrictId, &obj.Name, &obj.CityId, &obj.ProvinceId, &obj.Type)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = SubDistrict{}
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

func ShowSubDistrictByCityId(param_id string) (Response, error) {
	var obj SubDistrict
	var arrObj []SubDistrict
	var res Response

	con := db.CreateCon()

	qry := `SELECT * FROM smc_subdistrict WHERE s_city_id = ? ORDER BY s_name ASC`

	rows, err := con.Query(qry, param_id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = SubDistrict{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.SubDistrictId, &obj.Name, &obj.CityId, &obj.ProvinceId, &obj.Type)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = SubDistrict{}
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

func StoreSubDistrictBulk(subdistrict []SubDistrict) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = subdistrict
		return res, err
	}

	for idx := range subdistrict {
		obj := subdistrict[idx]

		qry := `INSERT INTO smc_subdistrict VALUES (?, ?, ?, ?, ?)`

		_, err = tx.ExecContext(ctx, qry, obj.SubDistrictId, obj.Name, obj.CityId, obj.ProvinceId, obj.Type)

		if err != nil {
			tx.Rollback()
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = subdistrict

			return res, err
		}

	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = subdistrict
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = subdistrict

	return res, nil
}

func StoreSubDistrict(subdistrict SubDistrict) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = subdistrict
		return res, err
	}

	qry := `INSERT INTO smc_subdistrict VALUES (?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, qry, subdistrict.SubDistrictId, subdistrict.Name, subdistrict.CityId, subdistrict.ProvinceId, subdistrict.Type)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = subdistrict

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = subdistrict
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = subdistrict

	return res, nil
}
