package sql

import (
	"database/sql"
	"errors"

	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
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

func FindAccount(db *sql.DB, accountID string) (*protoAccount.Account, error) {
	account := new(protoAccount.Account)
	balances := make([]*protoBalance.Balance, 0)
	query := `SELECT 
			a.id, 
			a.user_id, 
			a.exchange_name, 
			a.key_public, 
			a.key_secret, 
			a.description, 
			a.status,
			a.created_on,
			a.updated_on,
			b.id,
			b.user_id,
			b.account_id,
			b.currency_symbol,
			b.available,
			b.locked,
			b.exchange_total,
			b.exchange_available,
			b.exchange_locked,
			b.created_on,
			b.updated_on 
		FROM accounts a 
		JOIN balances b on a.id = b.account_id 
		WHERE account_id = $1`

	rows, err := db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		balance := new(protoBalance.Balance)
		err := rows.Scan(
			&account.AccountID,
			&account.UserID,
			&account.Exchange,
			&account.KeyPublic,
			&account.KeySecret,
			&account.Description,
			&account.Status,
			&account.CreatedOn,
			&account.UpdatedOn,
			&balance.BalanceID,
			&balance.UserID,
			&balance.AccountID,
			&balance.CurrencySymbol,
			&balance.Available,
			&balance.Locked,
			&balance.ExchangeTotal,
			&balance.ExchangeAvailable,
			&balance.ExchangeLocked,
			&balance.CreatedOn,
			&balance.UpdatedOn,
		)

		if err != nil {
			return nil, err
		}

		balances = append(balances, balance)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	account.Balances = balances

	return account, nil
}
