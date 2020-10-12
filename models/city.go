package models

import (
	"context"
	"net/http"
	"pb-dev-be/db"
)

type City struct {
	CityId     string `json:"city_id"`
	Name       string `json:"city_name"`
	ProvinceId string `json:"province_id"`
	PostalCode string `json:"postal_code"`
	Type       string `json:"type"`
}

func FetchAllCityData() (Response, error) {
	var obj City
	var arrObj []City
	var res Response

	con := db.CreateCon()

	qry := `SELECT * FROM smc_city ORDER BY s_name ASC`

	rows, err := con.Query(qry)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = City{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.CityId, &obj.Name, &obj.ProvinceId, &obj.PostalCode, &obj.Type)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = City{}
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

func GetCityByIdData(param_id string) (City, error) {
	var obj City
	var res Response

	con := db.CreateCon()

	qry := `SELECT * FROM smc_city WHERE s_city_id = ? ORDER BY s_name ASC`

	rows, err := con.Query(qry, param_id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = City{}
		return obj, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.CityId, &obj.Name, &obj.ProvinceId, &obj.PostalCode, &obj.Type)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = City{}
			return obj, err
		}
	}
	defer rows.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = obj

	return obj, nil
}

func ShowCityByProvinceIdData(param_id string) (Response, error) {
	var obj City
	var arrObj []City
	var res Response

	con := db.CreateCon()

	qry := `SELECT * FROM smc_city WHERE s_province_id = ? ORDER BY s_name ASC`

	rows, err := con.Query(qry, param_id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = City{}
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.CityId, &obj.Name, &obj.ProvinceId, &obj.PostalCode, &obj.Type)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = City{}
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

func StoreCityBulk(city []City) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = city
		return res, err
	}

	for idx := range city {
		obj := city[idx]

		qry := `INSERT INTO smc_city VALUES (?, ?, ?, ?, ?)`

		_, err = tx.ExecContext(ctx, qry, obj.CityId, obj.Name, obj.ProvinceId, obj.PostalCode, obj.Type)

		if err != nil {
			tx.Rollback()
			res.Status = http.StatusInternalServerError
			res.Message = err.Error()
			res.Data = city

			return res, err
		}

	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = city
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = city

	return res, nil
}

func StoreCity(city City) (Response, error) {
	var res Response
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = city
		return res, err
	}

	qry := `INSERT INTO smc_province VALUES (?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, qry, city.CityId, city.Name, city.ProvinceId, city.PostalCode, city.Type)

	if err != nil {
		tx.Rollback()
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = city

		return res, err
	}

	err = tx.Commit()
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		res.Data = city
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = city

	return res, nil
}
