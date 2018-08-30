package sql

import (
	"database/sql"

	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

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
