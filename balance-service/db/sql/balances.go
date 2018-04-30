package sql

import (
	"database/sql"
	"log"

	bp "github.com/asciiu/gomo/balance-service/proto/balance"
	"github.com/google/uuid"
)

//func DeleteOrder(db *sql.DB, orderId string) error {
//	_, err := db.Exec("DELETE FROM orders WHERE id = $1", orderId)
//	return err
//}
func NewBalance() *bp.Balance {
	balance := bp.Balance{
		Available:         0.0,
		Locked:            0.0,
		ExchangeTotal:     0.0,
		ExchangeAvailable: 0.0,
	}

	return &balance
}

func FindBalance(db *sql.DB, req *bp.GetUserBalanceRequest) (*bp.Balance, error) {
	b := NewBalance()

	err := db.QueryRow(`SELECT id, exchange_name, available, locked, exchange_total, exchange_available,
		exchange_locked FROM user_balances WHERE user_id = $1 and user_key_id = $2 and currency_name = $3`,
		req.UserId, req.ApiKeyId, req.Currency).
		Scan(&b.Id, &b.ExchangeName, &b.Available, &b.Locked, &b.ExchangeTotal, &b.ExchangeAvailable,
			&b.ExchangedLocked)

	b.UserId = req.UserId
	b.UserKeyId = req.ApiKeyId
	b.CurrencyName = req.Currency

	if err != nil {
		return nil, err
	}
	return b, nil
}

func FindBalancesByUserId(db *sql.DB, req *bp.GetUserBalancesRequest) (*bp.AccountBalances, error) {
	balances := make([]*bp.Balance, 0)

	rows, err := db.Query(`SELECT id, user_id, user_key_id, exchange_name, currency_name, available, 
		locked, exchange_total, exchange_available, exchange_locked FROM user_balances 
		WHERE user_id = $1`,
		req.UserId)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		b := NewBalance()

		err := rows.Scan(&b.Id, &b.UserId, &b.UserKeyId, &b.ExchangeName, &b.CurrencyName, &b.Available,
			&b.Locked, &b.ExchangeTotal, &b.ExchangeAvailable, &b.ExchangedLocked)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		balances = append(balances, b)
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
		newId := uuid.New()
		available := balance.Available
		locked := balance.Locked
		total := available + locked

		_, err := db.Exec(sqlStatement, newId, balance.UserId, balance.UserKeyId, balance.ExchangeName,
			balance.CurrencyName, available, locked, total, balance.ExchangeAvailable, balance.ExchangedLocked)

		if err != nil {
			return count, err
		}
		count += 1
	}

	return count, nil
}

//
//func UpdateOrder(db *sql.DB, req *orderProto.OrderRequest) (*orderProto.Order, error) {
//	sqlStatement := `UPDATE orders SET conditions = $1, price = $2, quantity = $3 WHERE id = $4 and user_id = $5 RETURNING exchange_name, user_key_id, status`
//
//	jsonCond, err := json.Marshal(req.Conditions)
//	if err != nil {
//		return nil, err
//	}
//
//	var o orderProto.Order
//	err = db.QueryRow(sqlStatement, jsonCond, req.Price, req.Qty, req.OrderId, req.UserId).
//		Scan(&o.Exchange, &o.ApiKeyId, &o.Status)
//
//	if err != nil {
//		return nil, err
//	}
//	order := &orderProto.Order{
//		OrderId:    req.OrderId,
//		UserId:     req.UserId,
//		ApiKeyId:   o.ApiKeyId,
//		Exchange:   o.Exchange,
//		Status:     o.Status,
//		Conditions: req.Conditions,
//		Price:      req.Price,
//		Qty:        req.Qty,
//	}
//	return order, nil
//}
//
//func UpdateOrderStatus(db *sql.DB, req *orderProto.OrderStatusRequest) (*orderProto.Order, error) {
//	sqlStatement := `UPDATE orders SET status = $1 WHERE id = $2 RETURNING user_id, exchange_name, market_name, user_key, side, price, quantity, status, conditions`
//
//	var o orderProto.Order
//	err := db.QueryRow(sqlStatement, req.Status, req.OrderId).
//		Scan(&o.UserId, &o.Exchange, &o.MarketName, &o.ApiKeyId, &o.Side, &o.Price, &o.Qty, &o.Status, &o.Conditions)
//
//	if err != nil {
//		return nil, err
//	}
//	order := &orderProto.Order{
//		OrderId:    req.OrderId,
//		ApiKeyId:   o.ApiKeyId,
//		UserId:     o.UserId,
//		Exchange:   o.Exchange,
//		MarketName: o.MarketName,
//		Side:       o.Side,
//		Price:      o.Price,
//		Qty:        o.Qty,
//		Status:     o.Status,
//		Conditions: o.Conditions,
//	}
//	return order, nil
//}
