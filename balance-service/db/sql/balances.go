package sql

import (
	"database/sql"
	"log"

	bp "github.com/asciiu/gomo/balance-service/proto/balance"
	"github.com/google/uuid"
)

func FindBalance(db *sql.DB, req *bp.GetUserBalanceRequest) (*bp.Balance, error) {
	b := bp.Balance{}

	err := db.QueryRow(`SELECT id, exchange_name, available, locked, exchange_total, exchange_available,
		exchange_locked FROM user_balances WHERE user_id = $1 and user_key_id = $2 and currency_name = $3`,
		req.UserID, req.KeyID, req.Currency).
		Scan(&b.ID, &b.ExchangeName, &b.Available, &b.Locked, &b.ExchangeTotal, &b.ExchangeAvailable,
			&b.ExchangeLocked)

	b.UserID = req.UserID
	b.KeyID = req.KeyID
	b.CurrencyName = req.Currency

	if err != nil {
		return nil, err
	}
	return &b, nil
}

func FindSymbolBalancesByUserID(db *sql.DB, req *bp.GetUserBalancesRequest) (*bp.AccountBalances, error) {
	balances := make([]*bp.Balance, 0)

	rows, err := db.Query(`SELECT id, user_id, user_key_id, exchange_name, currency_name, available, 
		locked, exchange_total, exchange_available, exchange_locked FROM user_balances 
		WHERE user_id = $1 and lower($2) = lower(currency_name)`,
		req.UserID, req.Symbol)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		b := bp.Balance{}

		err := rows.Scan(&b.ID, &b.UserID, &b.KeyID, &b.ExchangeName, &b.CurrencyName, &b.Available,
			&b.Locked, &b.ExchangeTotal, &b.ExchangeAvailable, &b.ExchangeLocked)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		balances = append(balances, &b)
	}
	accBalances := bp.AccountBalances{
		Balances: balances,
	}

	return &accBalances, nil
}

func FindAllBalancesByUserID(db *sql.DB, req *bp.GetUserBalancesRequest) (*bp.AccountBalances, error) {
	balances := make([]*bp.Balance, 0)

	rows, err := db.Query(`SELECT id, user_id, user_key_id, exchange_name, currency_name, available, 
		locked, exchange_total, exchange_available, exchange_locked FROM user_balances 
		WHERE user_id = $1`,
		req.UserID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		b := bp.Balance{}

		err := rows.Scan(&b.ID, &b.UserID, &b.KeyID, &b.ExchangeName, &b.CurrencyName, &b.Available,
			&b.Locked, &b.ExchangeTotal, &b.ExchangeAvailable, &b.ExchangeLocked)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		balances = append(balances, &b)
	}
	accBalances := bp.AccountBalances{
		Balances: balances,
	}

	return &accBalances, nil
}

// return count of inserts
func UpsertBalances(db *sql.DB, req *bp.AccountBalances) (int, error) {
	sqlStatement := `INSERT INTO user_balances (id, user_id, user_key_id, exchange_name, currency_name, 
		available, locked, exchange_total, exchange_available, exchange_locked) VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (user_key_id, currency_name) 
		DO UPDATE SET available = $6, locked = $7, exchange_total = $8, exchange_available = $9,
		exchange_locked = $10;`

	count := 0

	for _, balance := range req.Balances {
		newID := uuid.New()
		available := balance.Available
		locked := balance.Locked
		total := available + locked

		_, err := db.Exec(sqlStatement, newID, balance.UserID, balance.KeyID, balance.ExchangeName,
			balance.CurrencyName, available, locked, total, balance.ExchangeAvailable, balance.ExchangeLocked)

		if err != nil {
			return count, err
		}
		count += 1
	}

	return count, nil
}
