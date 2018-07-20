package sql

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/asciiu/gomo/common/constants/key"
	"github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/status"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"

	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// DeletePlan is a hard delete
func DeletePlan(db *sql.DB, planID string) error {
	_, err := db.Exec("DELETE FROM plans WHERE id = $1", planID)
	return err
}

// FindPlanSummary will load the order with the active order number as its only order
// this is so we can examine the current state of the plan
func FindPlanSummary(db *sql.DB, planID string) (*protoPlan.Plan, error) {
	var plan protoPlan.Plan
	var order protoOrder.Order
	var planTemp sql.NullString
	var orderTemp sql.NullString
	var price sql.NullFloat64
	var nextOrderID sql.NullString

	// note: always assume that a plan has at least one order.
	// a plan with no order should be removed from the system
	err := db.QueryRow(`SELECT 
		p.id,
		p.plan_template_id,
		p.user_id,
		p.user_key_id,
		k.description,
		p.active_order_number,
		p.exchange_name,
		p.market_name,
		p.base_balance,
		p.currency_balance,
		p.status,
		po.id,
		po.balance_percent,
		po.side,
		po.order_number,
		po.order_type,
		po.order_template_id,
		po.limit_price,
		po.next_order_id,
		po.status
		FROM plans p
		JOIN orders o on p.id = o.plan_id and p.active_order_number = o.order_number
		JOIN user_keys k on p.user_key_id = k.id
		WHERE p.id = $1`, planID).Scan(
		&plan.PlanID,
		&planTemp,
		&plan.UserID,
		&plan.KeyID,
		&plan.KeyDescription,
		&plan.ActiveOrderNumber,
		&plan.Exchange,
		&plan.MarketName,
		&plan.BaseBalance,
		&plan.CurrencyBalance,
		&plan.Status,
		&order.OrderID,
		&order.BalancePercent,
		&order.Side,
		&order.OrderNumber,
		&order.OrderType,
		&orderTemp,
		&price,
		&nextOrderID,
		&order.Status)

	if err != nil {
		return nil, err
	}
	if price.Valid {
		order.LimitPrice = price.Float64
	}
	if nextOrderID.Valid {
		order.NextOrderID = nextOrderID.String
	}
	if orderTemp.Valid {
		order.OrderTemplateID = orderTemp.String
	}
	if planTemp.Valid {
		plan.PlanTemplateID = planTemp.String
	}
	plan.Orders = append(plan.Orders, &order)

	return &plan, nil
}

// FindPlanWithPagedOrders ...
func FindPlanWithPagedOrders(db *sql.DB, req *protoPlan.GetUserPlanRequest) (*protoPlan.PlanWithPagedOrders, error) {

	var planPaged protoPlan.PlanWithPagedOrders
	var page protoOrder.OrdersPage

	var planTemp sql.NullString
	var orderTemp sql.NullString
	var price sql.NullFloat64
	var nextOrderID sql.NullString

	var count uint32
	queryCount := `SELECT count(*) FROM orders WHERE plan_id = $1`
	err := db.QueryRow(queryCount, req.PlanID).Scan(&count)

	if count == 0 {
		return nil, sql.ErrNoRows
	}

	rows, err := db.Query(`SELECT 
		p.id,
		p.plan_template_id,
		p.user_id,
		p.user_key_id,
		k.api_key,
		k.description,
		p.exchange_name,
		p.market_name,
		p.base_balance,
		p.currency_balance,
		p.status,
		p.created_on,
		p.updated_on,
		o.id,
		o.balance_percent,
		o.side,
		o.order_number,
		o.order_type,
		o.order_template_id,
		o.limit_price,
		o.next_order_id,
		o.status,
		o.created_on,
		o.updated_on,
		t.id as trigger_id,
		t.name,
		t.code,
		array_to_json(t.actions),
		t.trigger_number,
		t.triggered,
		t.trigger_template_id,
		t.created_on,
		t.updated_on
		FROM plans p
		JOIN orders o on p.id = o.plan_id
		JOIN triggers t on o.id = t.order_id
		JOIN user_keys k on p.user_key_id = k.id
		WHERE p.id = $1 ORDER BY o.order_number, t.trigger_number 
		OFFSET $2 LIMIT $3`, req.PlanID, req.Page, req.PageSize)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var order protoOrder.Order
		var trigger protoOrder.Trigger
		var triggerTempID sql.NullString
		var actionsStr string

		err := rows.Scan(
			&planPaged.PlanID,
			&planTemp,
			&planPaged.UserID,
			&planPaged.KeyID,
			&planPaged.Key,
			&planPaged.KeyDescription,
			&planPaged.Exchange,
			&planPaged.MarketName,
			&planPaged.BaseBalance,
			&planPaged.CurrencyBalance,
			&planPaged.Status,
			&planPaged.CreatedOn,
			&planPaged.UpdatedOn,
			&order.OrderID,
			&order.BalancePercent,
			&order.Side,
			&order.OrderNumber,
			&order.OrderType,
			&orderTemp,
			&price,
			&nextOrderID,
			&order.Status,
			&order.CreatedOn,
			&order.UpdatedOn,
			&trigger.TriggerID,
			&trigger.Name,
			&trigger.Code,
			&actionsStr,
			&trigger.TriggerNumber,
			&trigger.Triggered,
			&triggerTempID,
			&trigger.CreatedOn,
			&trigger.UpdatedOn,
		)

		if err != nil {
			return nil, err
		}
		if price.Valid {
			order.LimitPrice = price.Float64
		}
		if nextOrderID.Valid {
			order.NextOrderID = nextOrderID.String
		}
		if orderTemp.Valid {
			order.OrderTemplateID = orderTemp.String
		}
		if planTemp.Valid {
			planPaged.PlanTemplateID = planTemp.String
		}
		if triggerTempID.Valid {
			trigger.TriggerTemplateID = triggerTempID.String
		}
		if err := json.Unmarshal([]byte(actionsStr), &trigger.Actions); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(trigger.Code), &trigger.Code); err != nil {
			return nil, err
		}

		trigger.OrderID = order.OrderID

		var found = false
		for _, or := range page.Orders {
			if or.OrderID == order.OrderID {
				or.Triggers = append(or.Triggers, &trigger)
				found = true
				break
			}
		}

		if !found {
			order.Triggers = append(order.Triggers, &trigger)
			page.Orders = append(page.Orders, &order)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	page.Page = req.Page
	page.PageSize = req.PageSize
	page.Total = count
	planPaged.OrdersPage = &page

	return &planPaged, nil
}

// Find all active orders in the DB. This wil load the keys for each order.
// Returns active plans that have a verified key only.
func FindActivePlans(db *sql.DB) ([]*protoPlan.Plan, error) {
	activePlans := make([]*protoPlan.Plan, 0)

	rows, err := db.Query(`SELECT p.id as plan_id,
		p.user_id,
		p.user_key_id,
		k.api_key,
		k.secret,
		p.exchange_name,
		p.market_name,
		p.base_balance,
		p.currency_balance,
		p.status,
		o.id as order_id,
		o.balance_percent,
		o.side,
		o.order_type,
		o.limit_price, 
		o.next_order_id,
		o.status,
		t.id as trigger_id,
		t.name,
		t.code,
		array_to_json(t.actions),
		t.triggered
		FROM plans p 
		JOIN orders o on p.id = o.plan_id
		JOIN triggers t on o.id = t.order_id
		JOIN user_keys k on p.user_key_id = k.id
		WHERE p.status = $1 AND p.active_order_number = o.order_number 
		AND k.status = $2 ORDER BY o.id`, plan.Active, key.Verified)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var plan protoPlan.Plan
		var order protoOrder.Order
		var trigger protoOrder.Trigger
		var actionsStr string
		var price sql.NullFloat64
		var nextOrderID sql.NullString
		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.KeyID,
			&plan.Key,
			&plan.KeySecret,
			&plan.Exchange,
			&plan.MarketName,
			&plan.BaseBalance,
			&plan.CurrencyBalance,
			&plan.Status,
			&order.OrderID,
			&order.BalancePercent,
			&order.Side,
			&order.OrderType,
			&price,
			&nextOrderID,
			&order.Status,
			&trigger.TriggerID,
			&trigger.Name,
			&trigger.Code,
			&actionsStr,
			&trigger.Triggered)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		if price.Valid {
			order.LimitPrice = price.Float64
		}
		if nextOrderID.Valid {
			order.NextOrderID = nextOrderID.String
		}
		if err := json.Unmarshal([]byte(actionsStr), &trigger.Actions); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(trigger.Code), &trigger.Code); err != nil {
			return nil, err
		}

		var found = false
		for _, plan := range activePlans {
			if plan.Orders[0].OrderID == order.OrderID {
				plan.Orders[0].Triggers = append(plan.Orders[0].Triggers, &trigger)
				found = true
				break
			}
		}

		if !found {
			order.Triggers = append(order.Triggers, &trigger)
			plan.Orders = append(plan.Orders, &order)
			activePlans = append(activePlans, &plan)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return activePlans, nil
}

// Returns plan with specific order that has a verified key.
func FindPlanWithOrderID(db *sql.DB, orderID string) (*protoPlan.Plan, error) {
	var plan protoPlan.Plan
	var order protoOrder.Order

	var price sql.NullFloat64
	var nextOrderID sql.NullString

	rows, err := db.Query(`SELECT 
		p.id,
		p.user_id,
		p.user_key_id,
		k.api_key,
		k.secret,
		p.exchange_name,
		p.market_name,
		p.base_balance,
		p.currency_balance,
		p.status,
		o.id,
		o.balance_percent,
		o.side,
		o.order_type,
		o.limit_price, 
		o.next_order_id,
		o.status,
		o.created_on,
		o.updated_on,
		t.id as trigger_id,
		t.name,
		t.code,
		array_to_json(t.actions),
		t.trigger_number,
		t.triggered,
		t.trigger_template_id,
		t.created_on,
		t.updated_on
		FROM plans p 
		JOIN orders o on p.id = o.plan_id
		JOIN triggers t on o.id = t.order_id
		JOIN user_keys k on p.user_key_id = k.id
		WHERE o.id = $1 AND k.status = $2 ORDER BY t.trigger_number`, orderID, key.Verified)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var trigger protoOrder.Trigger
		var triggerTemplateID sql.NullString
		var actionsStr string
		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.KeyID,
			&plan.Key,
			&plan.KeySecret,
			&plan.Exchange,
			&plan.MarketName,
			&plan.BaseBalance,
			&plan.CurrencyBalance,
			&plan.Status,
			&order.OrderID,
			&order.BalancePercent,
			&order.Side,
			&order.OrderType,
			&price,
			&nextOrderID,
			&order.Status,
			&order.CreatedOn,
			&order.UpdatedOn,
			&trigger.TriggerID,
			&trigger.Name,
			&trigger.Code,
			&actionsStr,
			&trigger.TriggerNumber,
			&trigger.Triggered,
			&triggerTemplateID,
			&trigger.CreatedOn,
			&trigger.UpdatedOn,
		)

		if err != nil {
			return nil, err
		}
		if price.Valid {
			order.LimitPrice = price.Float64
		}
		if nextOrderID.Valid {
			order.NextOrderID = nextOrderID.String
		}
		if triggerTemplateID.Valid {
			trigger.TriggerTemplateID = triggerTemplateID.String
		}
		if err := json.Unmarshal([]byte(actionsStr), &trigger.Actions); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(trigger.Code), &trigger.Code); err != nil {
			return nil, err
		}
		trigger.OrderID = order.OrderID
		order.Triggers = append(order.Triggers, &trigger)
	}

	plan.Orders = append(plan.Orders, &order)

	return &plan, nil
}

func FindUserPlansWithStatus(db *sql.DB, userID, status string, page, pageSize uint32) (*protoPlan.PlansPage, error) {

	var count uint32
	queryCount := `SELECT count(*) FROM plans WHERE user_id = $1 AND status = $2`
	err := db.QueryRow(queryCount, userID, status).Scan(&count)

	plans := make([]*protoPlan.Plan, 0)
	query := `SELECT 
		id,
		plan_template_id,
		user_key_id,
		exchange_name,
		market_name,
		base_balance,
		currency_balance,
		status,
		created_on,
		updated_on
		FROM plans 
		WHERE user_id = $1 AND status = $2 OFFSET $3 LIMIT $4`

	rows, err := db.Query(query, userID, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var plan protoPlan.Plan
		var planTemplateID sql.NullString

		err := rows.Scan(
			&plan.PlanID,
			&planTemplateID,
			&plan.KeyID,
			&plan.Exchange,
			&plan.MarketName,
			&plan.BaseBalance,
			&plan.CurrencyBalance,
			&plan.Status,
			&plan.CreatedOn,
			&plan.UpdatedOn,
		)

		if err != nil {
			return nil, err
		}
		if planTemplateID.Valid {
			plan.PlanTemplateID = planTemplateID.String
		}

		plans = append(plans, &plan)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result := protoPlan.PlansPage{
		Page:     page,
		PageSize: pageSize,
		Total:    count,
		Plans:    plans,
	}

	return &result, nil
}

// Find all user plans for exchange with status
func FindUserExchangePlansWithStatus(db *sql.DB, userID, status, exchange string, page, pageSize uint32) (*protoPlan.PlansPage, error) {

	var count uint32
	queryCount := `SELECT count(*) FROM plans WHERE user_id = $1 AND status = $2 AND exchange_name = $3`
	err := db.QueryRow(queryCount, userID, status, exchange).Scan(&count)

	plans := make([]*protoPlan.Plan, 0)
	query := `SELECT 
		id,
		plan_template_id,
		user_key_id,
		exchange_name,
		market_name,
		base_balance,
		currency_balance,
		status,
		created_on,
		updated_on
		FROM plans 
		WHERE user_id = $1 AND status = $2 AND exchange_name = $3 OFFSET $4 LIMIT $5`

	rows, err := db.Query(query, userID, status, exchange, page, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var plan protoPlan.Plan
		var planTemplateID sql.NullString

		err := rows.Scan(
			&plan.PlanID,
			&planTemplateID,
			&plan.KeyID,
			&plan.Exchange,
			&plan.MarketName,
			&plan.BaseBalance,
			&plan.CurrencyBalance,
			&plan.Status,
			&plan.CreatedOn,
			&plan.UpdatedOn)

		if err != nil {
			return nil, err
		}
		if planTemplateID.Valid {
			plan.PlanTemplateID = planTemplateID.String
		}

		plans = append(plans, &plan)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result := protoPlan.PlansPage{
		Page:     page,
		PageSize: pageSize,
		Total:    count,
		Plans:    plans,
	}

	return &result, nil
}

// Find all user plans for exchange with status
func FindUserMarketPlansWithStatus(db *sql.DB, userID, status, exchange, marketName string, page, pageSize uint32) (*protoPlan.PlansPage, error) {

	var count uint32
	queryCount := `SELECT count(*) FROM plans WHERE user_id = $1 AND status = $2 AND exchange_name = $3 AND market_name = $4`
	err := db.QueryRow(queryCount, userID, status, exchange, marketName).Scan(&count)

	plans := make([]*protoPlan.Plan, 0)
	query := `SELECT 
		id,
		plan_template_id,
		user_key_id,
		exchange_name,
		market_name,
		base_balance,
		currency_balance,
		status,
		created_on,
		updated_on
		FROM plans 
		WHERE user_id = $1 AND status = $2 AND exchange_name = $3 AND market_name = $4 OFFSET $5 LIMIT $6`

	rows, err := db.Query(query, userID, status, exchange, marketName, page, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var plan protoPlan.Plan
		var planTemplateID sql.NullString

		err := rows.Scan(
			&plan.PlanID,
			&planTemplateID,
			&plan.KeyID,
			&plan.Exchange,
			&plan.MarketName,
			&plan.BaseBalance,
			&plan.CurrencyBalance,
			&plan.Status,
			&plan.CreatedOn,
			&plan.UpdatedOn,
		)

		if err != nil {
			return nil, err
		}
		if planTemplateID.Valid {
			plan.PlanTemplateID = planTemplateID.String
		}

		plans = append(plans, &plan)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result := protoPlan.PlansPage{
		Page:     page,
		PageSize: pageSize,
		Total:    count,
		Plans:    plans,
	}

	return &result, nil
}

// TODO this should be inserted as a transaction
func InsertPlan(db *sql.DB, req *protoPlan.PlanRequest) (*protoPlan.Plan, error) {
	planID := uuid.New()
	orderIDs := make([]uuid.UUID, 0)
	orders := make([]*protoOrder.Order, 0)

	// generate order ids
	for range req.Orders {
		orderIDs = append(orderIDs, uuid.New())
	}

	now := time.Now()

	sqlStatement := `insert into plans (
		id, 
		user_id, 
		user_key_id, 
		plan_template_id,
		active_order_number,
		exchange_name, 
		market_name, 
		base_balance, 
		currency_balance, 
		status, 
		created_on,
		updated_on) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := db.Exec(sqlStatement,
		planID,
		req.UserID,
		req.KeyID,
		req.PlanTemplateID,
		0, // active order number always starts out 0
		req.Exchange,
		req.MarketName,
		req.BaseBalance,
		req.CurrencyBalance,
		req.Status,
		now,
		now)

	if err != nil {
		return nil, err
	}

	// loop through and insert orders
	for i, or := range req.Orders {
		orderStatus := status.Inactive
		// first order of active plan is active
		if i == 0 && req.Status == plan.Active {
			orderStatus = status.Active
		}
		if err != nil {
			return nil, err
		}

		orderID := orderIDs[i]
		var nextOrderID uuid.UUID
		var nextIndex = i + 1
		if nextIndex < len(orderIDs) {
			nextOrderID = orderIDs[nextIndex]
		}
		sqlInsertOrder := `insert into orders (
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
			updated_on) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err = db.Exec(sqlInsertOrder,
			orderID,
			planID,
			or.BalancePercent,
			or.Side,
			or.OrderTemplateID,
			i,
			or.OrderType,
			or.LimitPrice,
			nextOrderID,
			orderStatus,
			now,
			now)

		if err != nil {
			return nil, err
		}

		triggers := make([]*protoOrder.Trigger, 0)
		// next insert all the triggers for this order
		for j, cond := range or.Triggers {
			jsonCode, err := json.Marshal(cond.Code)
			if err != nil {
				return nil, err
			}

			triggerID := uuid.New()
			sqlInsertOrder := `insert into triggers (
				id, 
				order_id, 
				trigger_number,
				trigger_template_id,
				name,
				code,
				actions,
				created_on,
				updated_on) 
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
			_, err = db.Exec(sqlInsertOrder,
				triggerID,
				orderID,
				j,
				cond.TriggerTemplateID,
				cond.Name,
				jsonCode,
				pq.Array(cond.Actions),
				now,
				now)

			if err != nil {
				return nil, err
			}

			trigger := protoOrder.Trigger{
				TriggerID:         triggerID.String(),
				TriggerNumber:     uint32(j),
				TriggerTemplateID: cond.TriggerTemplateID,
				OrderID:           orderID.String(),
				Name:              cond.Name,
				Code:              cond.Code,
				Actions:           cond.Actions,
				Triggered:         false,
				CreatedOn:         now.String(),
				UpdatedOn:         now.String(),
			}
			triggers = append(triggers, &trigger)
		}

		order := protoOrder.Order{
			OrderID:         orderID.String(),
			OrderNumber:     uint32(i),
			OrderType:       or.OrderType,
			OrderTemplateID: or.OrderTemplateID,
			Side:            or.Side,
			LimitPrice:      or.LimitPrice,
			BalancePercent:  or.BalancePercent,
			Status:          orderStatus,
			Triggers:        triggers,
			CreatedOn:       now.String(),
			UpdatedOn:       now.String(),

			NextOrderID: nextOrderID.String(),
		}
		orders = append(orders, &order)
	}

	plan := protoPlan.Plan{
		PlanID:            planID.String(),
		PlanTemplateID:    req.PlanTemplateID,
		UserID:            req.UserID,
		KeyID:             req.KeyID,
		ActiveOrderNumber: 0,
		Exchange:          req.Exchange,
		MarketName:        req.MarketName,
		BaseBalance:       req.BaseBalance,
		CurrencyBalance:   req.CurrencyBalance,
		Orders:            orders,
		Status:            req.Status,
		CreatedOn:         now.String(),
		UpdatedOn:         now.String(),
	}
	return &plan, nil
}

// Maybe this should return more of the updated order but all I need from this
// as of now is the next_order_id so I can load the next order.
func UpdatePlanOrder(db *sql.DB, orderID, stat string) (*protoOrder.Order, error) {
	updatePlanOrderSql := `UPDATE plan_orders SET status = $1 WHERE id = $2 
	RETURNING next_order_id, order_number, plan_id`

	var o protoOrder.Order
	var planID string
	err := db.QueryRow(updatePlanOrderSql, stat, orderID).
		Scan(&o.NextOrderID, &o.OrderNumber, &planID)

	if err != nil {
		return nil, err
	}

	if stat == status.Active {
		updatePlanSql := `UPDATE plans SET active_order_number = $1 WHERE id = $2`
		_, err = db.Exec(updatePlanSql, o.OrderNumber, planID)

		if err != nil {
			return nil, err
		}
	}

	return &o, nil
}

func UpdatePlanBalances(db *sql.DB, planID string, base, currency float64) error {
	sqlStatement := `UPDATE plans SET base_balance = $1, currency_balance = $2 WHERE id = $3`

	_, error := db.Exec(sqlStatement, base, currency, planID)
	return error
}

func UpdatePlanStatus(db *sql.DB, planID, status string) (*protoPlan.Plan, error) {
	sqlStatement := `UPDATE plans SET status = $1 WHERE id = $2 
	RETURNING id, plan_template_id, user_id, user_key_id, exchange_name, market_name, base_balance, currency_balance, status`

	var p protoPlan.Plan
	err := db.QueryRow(sqlStatement, status, planID).
		Scan(&p.PlanID,
			&p.PlanTemplateID,
			&p.UserID,
			&p.KeyID,
			&p.Exchange,
			&p.MarketName,
			&p.BaseBalance,
			&p.CurrencyBalance,
			&p.Status,
		)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func UpdatePlanBaseBalance(db *sql.DB, planID string, baseBalance float64) (*protoPlan.Plan, error) {
	sqlStatement := `UPDATE plans SET base_balance = $1 WHERE id = $2 
	RETURNING 
	id, 
	plan_template_id, 
	user_id, 
	user_key_id, 
	exchange_name, 
	market_name, 
	base_balance, 
	currency_balance, 
	status`

	var p protoPlan.Plan
	err := db.QueryRow(sqlStatement, baseBalance, planID).
		Scan(&p.PlanID,
			&p.PlanTemplateID,
			&p.UserID,
			&p.KeyID,
			&p.Exchange,
			&p.MarketName,
			&p.BaseBalance,
			&p.CurrencyBalance,
			&p.Status,
		)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func UpdatePlanCurrencyBalance(db *sql.DB, planID string, currencyBalance float64) (*protoPlan.Plan, error) {
	sqlStatement := `UPDATE plans SET currency_balance = $1 WHERE id = $2 
	RETURNING 
	id, 
	plan_template_id, 
	user_id, 
	user_key_id, 
	exchange_name, 
	market_name, 
	base_balance, 
	currency_balance, 
	status`

	var p protoPlan.Plan
	err := db.QueryRow(sqlStatement, currencyBalance, planID).
		Scan(&p.PlanID,
			&p.PlanTemplateID,
			&p.UserID,
			&p.KeyID,
			&p.Exchange,
			&p.MarketName,
			&p.BaseBalance,
			&p.CurrencyBalance,
			&p.Status,
		)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
