package sql

import (
	"database/sql"

	protoPrice "github.com/asciiu/gomo/analytics-service/proto/price"
	"github.com/lib/pq"
)

/*
This file has these functions:
	InsertPrices
*/

func InsertPrices(db *sql.DB, exchangeRates []*protoPrice.MarketPrice) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("exchange_rates",
		"exchange_name",
		"market_name",
		"closed_at_price",
		"closed_at_time",
	))
	if err != nil {
		return err
	}

	for _, rate := range exchangeRates {
		_, err = stmt.Exec(
			rate.Exchange,
			rate.MarketName,
			rate.ClosedAtPrice,
			rate.ClosedAtTime,
		)

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

	err = txn.Commit()

	return err
}
