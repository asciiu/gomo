package sql

import (
	"context"
	"database/sql"
	"errors"

	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	"github.com/lib/pq"
)

/*
All account specific queries shall live here.

This file has these functions:
	FindAccount
	FindAccounts
	InsertAccount
	UpdateAccountStatus
	UpdateAccountTxn
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

func FindAccount(db *sql.DB, accountID, userID string) (*protoAccount.Account, error) {
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
		WHERE a.id = $1 AND a.user_id = $2`

	rows, err := db.Query(query, accountID, userID)
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

	if account.AccountID == "" {
		return nil, sql.ErrNoRows
	}

	return account, nil
}

func FindAccounts(db *sql.DB, userID string) ([]*protoAccount.Account, error) {
	accounts := make([]*protoAccount.Account, 0)
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
		WHERE a.user_id = $1 ORDER BY a.created_on`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		account := new(protoAccount.Account)
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

		var found = false
		for _, a := range accounts {
			if a.AccountID == account.AccountID {
				// append balance to this account
				a.Balances = append(a.Balances, balance)
				found = true
				break
			}
		}
		if !found {
			// append new account
			account.Balances = append(account.Balances, balance)
			accounts = append(accounts, account)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func FindAccountKeys(db *sql.DB, accountIDs []string) ([]*protoAccount.AccountKey, error) {
	results := make([]*protoAccount.AccountKey, 0)

	rows, err := db.Query(`SELECT 
		id, 
		user_id, 
		exchange_name, 
		key_public, 
		key_secret FROM accounts WHERE id = Any($1)`, pq.Array(accountIDs))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var k protoAccount.AccountKey
		err := rows.Scan(&k.AccountID,
			&k.UserID,
			&k.Exchange,
			&k.KeyPublic,
			&k.KeySecret)

		if err != nil {
			return nil, err
		}
		results = append(results, &k)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func UpdateAccountStatus(db *sql.DB, accountID, userID, status string) (*protoAccount.Account, error) {
	stmt := `
		UPDATE accounts 
		SET 
			status = $1
		WHERE
			id = $2 AND user_id = $3
		RETURNING 
			id, 
			user_id, 
			exchange_name, 
			key_public, 
			key_secret, 
			description, 
			status,
			created_on,
			updated_on`

	account := new(protoAccount.Account)
	err := db.QueryRow(stmt, status, accountID, userID).
		Scan(
			&account.AccountID,
			&account.UserID,
			&account.Exchange,
			&account.KeyPublic,
			&account.KeySecret,
			&account.Description,
			&account.Status,
			&account.CreatedOn,
			&account.UpdatedOn)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func UpdateAccountTxn(txn *sql.Tx, ctx context.Context, accountID, userID, public, secret, description string) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE accounts 
		SET 
			key_public = $1, 
			key_secret = $2, 
			description = $3
		WHERE
			id = $4 AND user_id = $5`,
		public, secret, description, accountID, userID)

	return err
}
