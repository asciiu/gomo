package sql

import (
	"database/sql"
	"fmt"
	"strings"

	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
)

func InsertOrders(db *sql.DB, planID string, orders []*protoOrder.Order) error {
	totalColumns := 12
	valueStrings := make([]string, 0, len(orders))
	valueArgs := make([]interface{}, 0, len(orders)*totalColumns)

	for _, order := range orders {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, order.OrderID)
		valueArgs = append(valueArgs, planID)
		valueArgs = append(valueArgs, order.BalancePercent)
		valueArgs = append(valueArgs, order.Side)
		valueArgs = append(valueArgs, order.OrderTemplateID)
		valueArgs = append(valueArgs, order.OrderNumber)
		valueArgs = append(valueArgs, order.OrderType)
		valueArgs = append(valueArgs, order.LimitPrice)
		valueArgs = append(valueArgs, order.ParentOrderNumber)
		valueArgs = append(valueArgs, order.Status)
		valueArgs = append(valueArgs, order.CreatedOn)
		valueArgs = append(valueArgs, order.UpdatedOn)
	}

	stmt := fmt.Sprintf(`INSERT INTO orders (
		id, 
		plan_id, 
		balance_percent, 
		side,
		order_template_id,
		order_number,
		order_type,
		limit_price,
		next_order_id,
		status,
		created_on,
		updated_on) VALUES %s`, strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	return err
}
