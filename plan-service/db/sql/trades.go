package sql

import (
	"database/sql"

	protoTrade "github.com/asciiu/gomo/plan-service/proto/trade"
)

/*
This file has these functions:
	InsertTradeResult
*/

func InsertTradeResult(db *sql.DB, trade *protoTrade.Trade) error {
	stmt := `insert into trades (
		id, 
		order_id, 
		initial_currency_symbol,
		initial_currency_balance,
		initial_currency_traded,
		initial_currency_remainder,
		initial_currency_price,
		final_currency_symbol,
		final_currency_balance,
		fee_currency_symbol,
		fee_currency_amount,
		exchange_time,
		side,
		created_on,
		updated_on) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err := db.Exec(stmt,
		trade.TradeID,
		trade.OrderID,
		trade.InitialCurrencySymbol,
		trade.InitialCurrencyBalance,
		trade.InitialCurrencyTraded,
		trade.InitialCurrencyRemainder,
		trade.InitialCurrencyPrice,
		trade.FinalCurrencySymbol,
		trade.FinalCurrencyBalance,
		trade.FeeCurrencySymbol,
		trade.FeeCurrencyAmount,
		trade.ExchangeTime,
		trade.Side,
		trade.CreatedOn,
		trade.UpdatedOn,
	)

	return err
}
