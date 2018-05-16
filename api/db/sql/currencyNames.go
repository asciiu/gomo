package sql

import (
	"database/sql"
	"log"

	"github.com/asciiu/gomo/api/models"
)

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
