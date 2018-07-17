package sql

import (
	"database/sql"
	"encoding/json"

	protoOrder "github.com/asciiu/gomo/order-service/proto/order"
	"github.com/lib/pq"
)

// DeleteOrder
func DeleteOrder(db *sql.DB, orderID string) error {
	_, err := db.Exec("DELETE FROM orders WHERE id = $1", orderID)
	return err
}

func FindOrder(db *sql.DB, orderID string) (*protoOrder.Order, error) {
	var order protoOrder.Order
	var planID string
	var orderTemp sql.NullString
	var limitPrice sql.NullFloat64
	var nextOrderID sql.NullString

	// note: always assume that a plan has at least one order.
	// a plan with no order should be removed from the system
	err := db.QueryRow(`SELECT 
		id,
		plan_id,
		order_template_id,
		balance_percent,
		side,
		order_number,
		order_type,
		limit_price,
		next_order_id,
		status
		FROM orders 
		WHERE id = $1`, orderID).Scan(
		&order.OrderID,
		&planID,
		&orderTemp,
		&order.BalancePercent,
		&order.Side,
		&order.OrderNumber,
		&order.OrderType,
		&limitPrice,
		&nextOrderID,
		&order.Status)

	if err != nil {
		return nil, err
	}
	if limitPrice.Valid {
		order.LimitPrice = limitPrice.Float64
	}
	if nextOrderID.Valid {
		order.NextOrderID = nextOrderID.String
	}
	if orderTemp.Valid {
		order.OrderTemplateID = orderTemp.String
	}

	return &order, nil
}

func UpdateBalancePercent(db *sql.DB, orderID string, percent float64) (*protoOrder.Order, error) {

	updatePlanOrderSql := `UPDATE orders SET balance_percent = $1 WHERE id = $2
	RETURNING 
	order_template_id,
	balance_percent,
	side,
	order_number, 
	order_type,
	limit_price,
	next_order_id,
	status`

	var o protoOrder.Order
	var orderTemplateID sql.NullString
	err := db.QueryRow(updatePlanOrderSql, percent, orderID).
		Scan(
			&orderTemplateID,
			&o.BalancePercent,
			&o.Side,
			&o.OrderNumber,
			&o.OrderType,
			&o.LimitPrice,
			&o.NextOrderID,
			&o.Status)

	if err != nil {
		return nil, err
	}
	if orderTemplateID.Valid {
		o.OrderTemplateID = orderTemplateID.String
	}

	return &o, nil
}

func UpdateLimitPrice(db *sql.DB, orderID string, price float64) (*protoOrder.Order, error) {

	updatePlanOrderSql := `UPDATE orders SET limit_price = $1 WHERE id = $2
	RETURNING 
	order_template_id,
	balance_percent,
	side,
	order_number, 
	order_type,
	limit_price,
	next_order_id,
	status`

	var o protoOrder.Order
	var orderTemplateID sql.NullString
	err := db.QueryRow(updatePlanOrderSql, price, orderID).
		Scan(
			&orderTemplateID,
			&o.BalancePercent,
			&o.Side,
			&o.OrderNumber,
			&o.OrderType,
			&o.LimitPrice,
			&o.NextOrderID,
			&o.Status)

	if err != nil {
		return nil, err
	}
	if orderTemplateID.Valid {
		o.OrderTemplateID = orderTemplateID.String
	}

	return &o, nil
}

func UpdateTriggerName(db *sql.DB, triggerID, name string) (*protoOrder.Trigger, error) {

	updatePlanOrderSql := `UPDATE triggers SET name = $1 WHERE id = $2
	RETURNING 
	order_id,
	trigger_number,
	trigger_template_id,
	name,
	code,
	array_to_json(actions),
	triggered`

	var t protoOrder.Trigger
	var triggerTemplateID sql.NullString
	var actionsStr string
	err := db.QueryRow(updatePlanOrderSql, name, triggerID).
		Scan(
			&t.OrderID,
			&t.TriggerNumber,
			&triggerTemplateID,
			&t.Name,
			&t.Code,
			&actionsStr,
			&t.Triggered)

	if err != nil {
		return nil, err
	}
	if triggerTemplateID.Valid {
		t.TriggerTemplateID = triggerTemplateID.String
	}
	if err := json.Unmarshal([]byte(actionsStr), &t.Actions); err != nil {
		return nil, err
	}

	return &t, nil
}

func UpdateTriggerCode(db *sql.DB, triggerID, code string) (*protoOrder.Trigger, error) {

	updatePlanOrderSql := `UPDATE triggers SET code = $1 WHERE id = $2
	RETURNING 
	order_id,
	trigger_number,
	trigger_template_id,
	name,
	code,
	array_to_json(actions),
	triggered`

	var t protoOrder.Trigger
	var triggerTemplateID sql.NullString
	var actionsStr string
	jsonCode, err := json.Marshal(code)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(updatePlanOrderSql, jsonCode, triggerID).
		Scan(
			&t.OrderID,
			&t.TriggerNumber,
			&triggerTemplateID,
			&t.Name,
			&t.Code,
			&actionsStr,
			&t.Triggered)

	if err != nil {
		return nil, err
	}
	if triggerTemplateID.Valid {
		t.TriggerTemplateID = triggerTemplateID.String
	}
	if err := json.Unmarshal([]byte(actionsStr), &t.Actions); err != nil {
		return nil, err
	}

	return &t, nil
}

func UpdateTriggerActions(db *sql.DB, triggerID, actions []string) (*protoOrder.Trigger, error) {

	updatePlanOrderSql := `UPDATE triggers SET actions = $1 WHERE id = $2
	RETURNING 
	order_id,
	trigger_number,
	trigger_template_id,
	name,
	code,
	array_to_json(actions),
	triggered`

	var t protoOrder.Trigger
	var triggerTemplateID sql.NullString
	var actionsStr string
	err := db.QueryRow(updatePlanOrderSql, pq.Array(actions), triggerID).
		Scan(
			&t.OrderID,
			&t.TriggerNumber,
			&triggerTemplateID,
			&t.Name,
			&t.Code,
			&actionsStr,
			&t.Triggered)

	if err != nil {
		return nil, err
	}
	if triggerTemplateID.Valid {
		t.TriggerTemplateID = triggerTemplateID.String
	}
	if err := json.Unmarshal([]byte(actionsStr), &t.Actions); err != nil {
		return nil, err
	}

	return &t, nil
}
