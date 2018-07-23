package sql

import (
	"database/sql"

	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
)

func InsertOrders(txn *sql.Tx, orders []*protoOrder.Order) error {
	stmt, err := txn.Prepare(pq.CopyIn("orders",
		"id",
		"user_key_id",
		"parent_order_id",
		"plan_id",
		"plan_depth",
		"exchange_name",
		"market_name",
		"currency_symbol",
		"currency_balance",
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
			order.CurrencySymbol,
			order.CurrencyBalance,
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

func UpdateOrderStatus(db *sql.DB, orderID, status string) (*protoOrder.Order, error) {
	updatePlanOrderSql := `UPDATE orders SET status = $1 WHERE id = $2 
	RETURNING plan_id`

	var o protoOrder.Order
	err := db.QueryRow(updatePlanOrderSql, status, orderID).
		Scan(&o.PlanID)

	if err != nil {
		return nil, err
	}

	return &o, nil
}
