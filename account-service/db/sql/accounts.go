package sql

import (
	"database/sql"
	"errors"

	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
)

/*
All account specific queries shall live here.

This file has this functions:
	InsertAccount
*/

func InsertAccount(db *sql.DB, newAccount *protoAccount.Account) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}

	sqlStatement, err := txn.Prepare(`insert into accounts (
		id, 
		user_id, 
		exchange_name, 
		key_public, 
		key_secret, 
		description, 
		status,
		created_on,
		updated_on) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`)

	if err != nil {
		txn.Rollback()
		return errors.New("prepare insert plan failed: " + err.Error())
	}
	_, err = txn.Stmt(sqlStatement).Exec(
		newAccount.AccountID,
		newAccount.UserID,
		newAccount.Exchange,
		newAccount.KeyPublic,
		newAccount.KeySecret,
		newAccount.Description,
		newAccount.Status,
		newAccount.CreatedOn,
		newAccount.UpdatedOn,
	)

	if err != nil {
		txn.Rollback()
		return errors.New("account insert failed: " + err.Error())
	}
	if err := InsertBalances(txn, newAccount.Balances); err != nil {
		txn.Rollback()
		return errors.New("insert balances failed: " + err.Error())
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return err
}
