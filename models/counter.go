package models

import (
	"context"
	"pb-dev-be/db"
)

type Counter struct {
	CounterId string `json:"counter_id"`
	Count     int    `json:"count"`
}

func GetCounterById(param_id string) (Counter, error) {
	var counter Counter

	con := db.CreateCon()
	qry := `SELECT * FROM smc_counter WHERE s_counter_id = ?`

	rows, err := con.Query(qry, param_id)
	if err != nil {
		return counter, err
	}

	for rows.Next() {
		err = rows.Scan(&counter.CounterId, &counter.Count)
		if err != nil {
			return counter, err
		}
	}
	defer rows.Close()

	return counter, nil
}

func UpdateCounterCount(param_id string) error {
	count, err := GetCounterById(param_id)
	con := db.CreateCon()

	if err != nil {
		return err
	}

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	count.Count += 1
	qry := `UPDATE smc_counter SET s_count = ? WHERE s_counter_id = ?`

	_, err = tx.ExecContext(ctx, qry, count.Count, count.CounterId)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
