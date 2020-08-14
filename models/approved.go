package models

import (
	"context"
	"errors"
	"pb-dev-be/db"
)

type Approved struct {
	ObjectId    string `json:"object_id"`
	Id          string `json:"id"`
	Approved_at string `json:"approved_at"`
	Status      string `json:"status"`
	UserId      string `json:"user_id"`
}

func CreateApproved(approved Approved) error {
	con := db.CreateCon()

	ctx := context.Background()
	tx, err := con.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	qry := `INSERT INTO smc_approved VALUES(?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, qry, approved.ObjectId, approved.Id, approved.Approved_at, approved.Status, approved.UserId)
	if err != nil {
		tx.Rollback()
		er := "Error While Creating Approved" + err.Error()
		return errors.New(er)
	}

	err = tx.Commit()
	if err != nil {
		er := "Error While Committing Process Approved" + err.Error()
		return errors.New(er)
	}

	return nil
}
