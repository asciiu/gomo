package sql

import (
	"database/sql"
	"log"

	protoPrice "github.com/asciiu/gomo/analytics-service/proto/analytics"
	"github.com/asciiu/gomo/api/models"
	"github.com/lib/pq"
)

/*
This file has these functions:
	DeleteExchangeRates
	FindExchangeRate
	FindExchangeRates
	InsertPrices
*/

// used in dropping test data
func DeleteExchangeRates(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM exchange_rates")
	return err
}

func FindExchangeRate(db *sql.DB, exchange, marketName, time string) (*protoPrice.MarketPrice, error) {
	p := new(protoPrice.MarketPrice)
	err := db.QueryRow(`SELECT 
			closed_at_price
		FROM exchange_rates 
		WHERE exchange_name = $1 AND market_name = $2 AND closed_at_time = $3`,
		exchange, marketName, time).Scan(&p.ClosedAtPrice)

	if err != nil {
		return nil, err
	}
	p.Exchange = exchange
	p.MarketName = marketName
	p.ClosedAtTime = time

	return p, nil
}

func FindExchangeRates(db *sql.DB, exchange, atTime string) ([]*protoPrice.MarketPrice, error) {
	rates := make([]*protoPrice.MarketPrice, 0)
	query := `SELECT 
		exchange_name,
		market_name,
		closed_at_price,
		closed_at_time
		FROM exchange_rates 
		WHERE exchange_name = $1 AND closed_at_time = $2 ORDER BY market_name`

	rows, err := db.Query(query, exchange, atTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rate protoPrice.MarketPrice

		err := rows.Scan(
			&rate.Exchange,
			&rate.MarketName,
			&rate.ClosedAtPrice,
			&rate.ClosedAtTime)

		if err != nil {
			return nil, err
		}

		rates = append(rates, &rate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return rates, nil
}

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

func GetCurrencyNames(db *sql.DB) ([]*models.Currency, error) {
	results := make([]*models.Currency, 0)

	rows, err := db.Query("SELECT currency_name, currency_symbol FROM currency_names")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.Currency
		err := rows.Scan(&c.CurrencyName, &c.TickerSymbol)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		results = append(results, &c)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return results, nil
}
