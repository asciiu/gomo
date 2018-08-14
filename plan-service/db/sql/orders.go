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

func InsertOrders(txn *sql.Tx, orders []*protoOrder.Order) error {
	stmt, err := txn.Prepare(pq.CopyIn("orders",
		"id",
		"user_key_id",
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
			uuid.FromStringOrNil(order.KeyID),
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

// returns (planID, depth, err)
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
