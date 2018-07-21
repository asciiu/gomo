package sql

import (
	"database/sql"

	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
)

// func InsertOrders(db *sql.DB, planID string, orders []*protoOrder.Order) error {
// 	totalColumns := 11
// 	valueStrings := make([]string, 0, len(orders))
// 	valueArgs := make([]interface{}, 0, len(orders)*totalColumns)

// 	for _, order := range orders {
// 		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
// 		valueArgs = append(valueArgs, order.OrderID)
// 		valueArgs = append(valueArgs, planID)
// 		valueArgs = append(valueArgs, order.BalancePercent)
// 		valueArgs = append(valueArgs, order.Side)
// 		valueArgs = append(valueArgs, order.OrderTemplateID)
// 		valueArgs = append(valueArgs, order.OrderType)
// 		valueArgs = append(valueArgs, order.LimitPrice)
// 		valueArgs = append(valueArgs, order.ParentOrderID)
// 		valueArgs = append(valueArgs, order.Status)
// 		valueArgs = append(valueArgs, order.CreatedOn)
// 		valueArgs = append(valueArgs, order.UpdatedOn)
// 	}

// 	sqlStr := fmt.Sprintf(`INSERT INTO orders (
// 		id,
// 		plan_id,
// 		balance_percent,
// 		side,
// 		order_template_id,
// 		order_type,
// 		limit_price,
// 		parent_order_id,
// 		status,
// 		created_on,
// 		updated_on) VALUES %s`, strings.Join(valueStrings, ", "))

// 	_, err := db.Exec(sqlStr, valueArgs...)

// 	if err != nil {
// 		return errors.New("bulk insert InsertOrders failed: " + err.Error())
// 	}

// 	return nil
// }
func InsertOrders(txn *sql.Tx, planID string, orders []*protoOrder.Order) error {
	stmt, err := txn.Prepare(pq.CopyIn("orders",
		"id",
		"plan_id",
		"plan_depth",
		"balance_percent",
		"side",
		"order_template_id",
		"order_type",
		"limit_price",
		"parent_order_id",
		"status",
		"created_on",
		"updated_on"))
	if err != nil {
		return err
	}

	for _, order := range orders {
		_, err = stmt.Exec(
			uuid.FromStringOrNil(order.OrderID),
			uuid.FromStringOrNil(planID),
			order.PlanDepth,
			order.BalancePercent,
			order.Side,
			order.OrderTemplateID,
			order.OrderType,
			order.LimitPrice,
			uuid.FromStringOrNil(order.ParentOrderID),
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
