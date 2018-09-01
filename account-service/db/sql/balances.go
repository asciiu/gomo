package sql

import (
	"context"
	"database/sql"

	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

/*
All balance specific queries shall live here.

This file has these functions:
	InsertBalances
	UpdateBalanceTxn
*/

func InsertBalances(txn *sql.Tx, balances []*protoBalance.Balance) error {
	stmt, err := txn.Prepare(pq.CopyIn("balances",
		"id",
		"user_id",
		"account_id",
		"currency_symbol",
		"available",
		"locked",
		"exchange_total",
		"exchange_available",
		"exchange_locked",
		"created_on",
		"updated_on"))

	if err != nil {
		return err
	}

	for _, balance := range balances {
		balanceID := uuid.Must(uuid.NewV4(), nil)
		_, err = stmt.Exec(
			uuid.FromStringOrNil(balanceID.String()),
			uuid.FromStringOrNil(balance.UserID),
			uuid.FromStringOrNil(balance.AccountID),
			balance.CurrencySymbol,
			balance.Available,
			balance.Locked,
			balance.ExchangeTotal,
			balance.ExchangeAvailable,
			balance.ExchangeLocked,
			balance.CreatedOn,
			balance.UpdatedOn)

		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func UpdateBalanceTxn(txn *sql.Tx, ctx context.Context, accountID, userID, currencySymbol string, total, available, locked float64) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE balances 
		SET 
			exchange_total = $1, 
			exchange_available = $2, 
			exchange_locked = $3
		WHERE
			currency_symbol = $4 AND account_id = $5 AND user_id = $6`,
		total, available, locked, currencySymbol, accountID, userID)

	return err
}
