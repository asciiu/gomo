package sql

import (
	"context"
	"database/sql"

	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
)

func DeleteOrders(txn *sql.Tx, ctx context.Context, orderIDs []string) error {
	_, err := txn.ExecContext(ctx, `
		DELETE FROM orders 
		WHERE
			id = Any($1)`, pq.Array(orderIDs))

	if err != nil {
		return err
	}

	return nil
}

func FindAccountOrders(db *sql.DB, accountID string) ([]*protoOrder.Order, error) {
	rows, err := db.Query(`SELECT 
		id as order_id,
		account_id,
		parent_order_id,
		plan_id,
		plan_depth,
		exchange_name,
		market_name,
		initial_currency_symbol,
		initial_currency_balance,
		initial_currency_traded,
		order_template_id,
		order_type,
		side,
		limit_price, 
		status,
		grupo,
		created_on,
		updated_on
		FROM orders 
		WHERE account_id = $1`, accountID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := make([]*protoOrder.Order, 0)

	for rows.Next() {
		var price sql.NullFloat64
		var order protoOrder.Order

		err := rows.Scan(
			&order.OrderID,
			&order.AccountID,
			&order.ParentOrderID,
			&order.PlanID,
			&order.PlanDepth,
			&order.Exchange,
			&order.MarketName,
			&order.InitialCurrencySymbol,
			&order.InitialCurrencyBalance,
			&order.InitialCurrencyTraded,
			&order.OrderTemplateID,
			&order.OrderType,
			&order.Side,
			&price,
			&order.Status,
			&order.Grupo,
			&order.CreatedOn,
			&order.UpdatedOn)

		if err != nil {
			return nil, err
		}
		if price.Valid {
			order.LimitPrice = price.Float64
		}

		orders = append(orders, &order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func InsertOrders(txn *sql.Tx, orders []*protoOrder.Order) error {
	stmt, err := txn.Prepare(pq.CopyIn("orders",
		"id",
		"account_id",
		"parent_order_id",
		"plan_id",
		"plan_depth",
		"exchange_name",
		"market_name",
		"initial_currency_symbol",
		"initial_currency_balance",
		"grupo",
		"order_priority",
		"order_template_id",
		"order_type",
		"side",
		"limit_price",
		"status",
		"created_on",
		"updated_on"))
	if err != nil {
		return err
	}

	for _, order := range orders {
		_, err = stmt.Exec(
			uuid.FromStringOrNil(order.OrderID),
			uuid.FromStringOrNil(order.AccountID),
			uuid.FromStringOrNil(order.ParentOrderID),
			uuid.FromStringOrNil(order.PlanID),
			order.PlanDepth,
			order.Exchange,
			order.MarketName,
			order.InitialCurrencySymbol,
			order.InitialCurrencyBalance,
			order.Grupo,
			order.OrderPriority,
			order.OrderTemplateID,
			order.OrderType,
			order.Side,
			order.LimitPrice,
			order.Status,
			order.CreatedOn,
			order.UpdatedOn)

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

// returns (planID, plan_depth, error)
func UpdateOrderStatus(db *sql.DB, orderID, status string) (string, int32, error) {
	updatePlanOrderSql := `UPDATE orders SET status = $1 WHERE id = $2 
	RETURNING plan_id, plan_depth`

	var planID string
	var depth int32
	err := db.QueryRow(updatePlanOrderSql, status, orderID).
		Scan(&planID, &depth)

	if err != nil {
		return "", 0, err
	}

	return planID, depth, nil
}

// returns (planID, depth, err)
func UpdateOrderStatusAndBalance(db *sql.DB, orderID, status string, initBalance float64) (string, int32, error) {
	updatePlanOrderSql := `
	UPDATE orders 
	SET status = $1, initial_currency_balance = $2
	WHERE id = $3 
	RETURNING plan_id, plan_depth`

	var planID string
	var depth int32
	err := db.QueryRow(updatePlanOrderSql, status, initBalance, orderID).
		Scan(&planID, &depth)

	if err != nil {
		return "", 0, err
	}

	return planID, depth, nil
}

func UpdateOrderResults(db *sql.DB, orderID string, traded, remainder, finalBalance float64, finalSymbol string) error {
	updatePlanOrderSql := `
	UPDATE orders 
	SET 
		initial_currency_traded = $1,
		initial_currency_remainder = $2,
		final_currency_symbol = $3,
		final_currency_balance = $4
	WHERE id = $5`

	_, err := db.Exec(updatePlanOrderSql,
		traded,
		remainder,
		finalSymbol,
		finalBalance,
		orderID)

	return err
}
