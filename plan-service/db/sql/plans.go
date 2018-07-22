package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/asciiu/gomo/common/constants/key"
	"github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/status"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	"github.com/lib/pq"

	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	"github.com/google/uuid"
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
	var planTemp sql.NullString
	var orderTemp sql.NullString
	var price sql.NullFloat64

	// note: always assume that a plan has at least one order.
	// a plan with no order should be removed from the system
	rows, err := db.Query(`SELECT 
		p.id,
		p.plan_template_id,
		p.user_id,
		p.user_key_id,
		k.description,
		p.last_executed_order_id,
		p.exchange_name,
		p.market_name,
		p.base_balance,
		p.currency_balance,
		p.status,
		o.id,
		o.balance_percent,
		o.side,
		o.parent_order_id,
		o.order_type,
		o.order_template_id,
		o.limit_price,
		o.status
		FROM plans p
		JOIN orders o on p.id = o.plan_id and p.last_executed_order_id = o.parent_order_id
		JOIN user_keys k on p.user_key_id = k.id
		WHERE p.id = $1`, planID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var order protoOrder.Order
		err := rows.Scan(
			&plan.PlanID,
			&planTemp,
			&plan.UserID,
			&plan.KeyID,
			&plan.KeyDescription,
			&plan.LastExecutedOrderID,
			&plan.Exchange,
			&plan.MarketName,
			&plan.BaseBalance,
			&plan.CurrencyBalance,
			&plan.Status,
			&order.OrderID,
			&order.BalancePercent,
			&order.Side,
			&order.ParentOrderID,
			&order.OrderType,
			&orderTemp,
			&price,
			&order.Status)

		if err != nil {
			return nil, err
		}
		if price.Valid {
			order.LimitPrice = price.Float64
		}
		if orderTemp.Valid {
			order.OrderTemplateID = orderTemp.String
		}
		if planTemp.Valid {
			plan.PlanTemplateID = planTemp.String
		}
		plan.Orders = append(plan.Orders, &order)

	}
	return &plan, nil
}

// FindPlanWithPagedOrders ...
func FindPlanWithPagedOrders(db *sql.DB, req *protoPlan.GetUserPlanRequest) (*protoPlan.PlanWithPagedOrders, error) {

	var planPaged protoPlan.PlanWithPagedOrders
	var page protoOrder.OrdersPage

	var planTemp sql.NullString
	var orderTemp sql.NullString
	var price sql.NullFloat64

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
		o.order_type,
		o.order_template_id,
		o.limit_price,
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
		WHERE p.id = $1 ORDER BY o.plan_depth, t.trigger_number 
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
			&order.OrderType,
			&orderTemp,
			&price,
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
func FindChildOrders(db *sql.DB, planID, orderID string) (*protoPlan.Plan, error) {
	var plan protoPlan.Plan
	var order protoOrder.Order

	var price sql.NullFloat64

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
		o.parent_order_id,
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
		WHERE p.id = $1 AND o.parent_order_id = $2 AND k.status = $3 ORDER BY t.trigger_number`, planID, orderID, key.Verified)

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
			&order.ParentOrderID,
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
		last_executed_plan_depth,
		last_executed_order_id,
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
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
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
		last_executed_plan_depth,
		last_executed_order_id,
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
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
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
		last_executed_plan_depth,
		last_executed_order_id,
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
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
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
func InsertPlan(db *sql.DB, req *protoPlan.NewPlanRequest) (*protoPlan.Plan, error) {
	none := uuid.Nil.String()

	txn, err := db.Begin()
	if err != nil {
		return nil, err
	}

	sqlStatement, err := txn.Prepare(`insert into plans (
		id, 
		user_id, 
		user_key_id, 
		plan_template_id,
		last_executed_plan_depth,
		last_executed_order_id,
		exchange_name, 
		market_name, 
		base_balance, 
		currency_balance, 
		status, 
		created_on,
		updated_on) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`)

	if err != nil {
		txn.Rollback()
		return nil, errors.New("prepare insert plan failed: " + err.Error())
	}

	planID := uuid.New()
	now := string(pq.FormatTimestamp(time.Now().UTC()))

	_, err = txn.Stmt(sqlStatement).Exec(
		planID,
		req.UserID,
		req.KeyID,
		req.PlanTemplateID,
		0,
		none,
		req.Exchange,
		req.MarketName,
		req.BaseBalance,
		req.CurrencyBalance,
		req.Status,
		now,
		now)

	if err != nil {
		txn.Rollback()
		return nil, errors.New("plan insert failed: " + err.Error())
	}

	newTriggers := make([]*protoOrder.Trigger, 0, len(req.Orders))
	newOrders := make([]*protoOrder.Order, 0, len(req.Orders))

	for _, or := range req.Orders {
		//orderID := uuid.New()
		orderStatus := status.Inactive
		depth := uint32(1)

		if or.ParentOrderID != none {
			for _, o := range newOrders {
				if o.OrderID == or.ParentOrderID {
					depth = o.PlanDepth + 1
					break
				}
			}
		}

		if or.ParentOrderID == none && req.Status == plan.Active {
			orderStatus = status.Active
		}

		// collect triggers for this order
		triggers := make([]*protoOrder.Trigger, 0, len(or.Triggers))
		for j, cond := range or.Triggers {
			triggerID := uuid.New()
			trigger := protoOrder.Trigger{
				TriggerID:         triggerID.String(),
				TriggerNumber:     uint32(j),
				TriggerTemplateID: cond.TriggerTemplateID,
				OrderID:           or.OrderID,
				Name:              cond.Name,
				Code:              cond.Code,
				Actions:           cond.Actions,
				Triggered:         false,
				CreatedOn:         now,
				UpdatedOn:         now,
			}
			triggers = append(triggers, &trigger)
			// append to all new triggers so we can bulk insert
			newTriggers = append(newTriggers, &trigger)
		}

		order := protoOrder.Order{
			OrderID:         or.OrderID,
			OrderType:       or.OrderType,
			OrderTemplateID: or.OrderTemplateID,
			ParentOrderID:   or.ParentOrderID,
			PlanDepth:       depth,
			Side:            or.Side,
			LimitPrice:      or.LimitPrice,
			BalancePercent:  or.BalancePercent,
			Status:          orderStatus,
			Triggers:        triggers,
			CreatedOn:       now,
			UpdatedOn:       now,
		}
		newOrders = append(newOrders, &order)
	}

	if err := InsertOrders(txn, planID.String(), newOrders); err != nil {
		txn.Rollback()
		return nil, errors.New("bulk insert InsertOrders failed: " + err.Error())
	}
	if err := InsertTriggers(txn, newTriggers); err != nil {
		txn.Rollback()
		return nil, errors.New("bulk insert InsertTriggers failed: " + err.Error())
	}

	err = txn.Commit()
	if err != nil {
		return nil, err
	}

	plan := protoPlan.Plan{
		PlanID:         planID.String(),
		PlanTemplateID: req.PlanTemplateID,
		UserID:         req.UserID,
		KeyID:          req.KeyID,
		LastExecutedPlanDepth: 0,
		LastExecutedOrderID:   none,
		Exchange:              req.Exchange,
		MarketName:            req.MarketName,
		BaseBalance:           req.BaseBalance,
		CurrencyBalance:       req.CurrencyBalance,
		Orders:                newOrders,
		Status:                req.Status,
		CreatedOn:             now,
		UpdatedOn:             now,
	}

	return &plan, nil
}

// Maybe this should return more of the updated order but all I need from this
// as of now is the order number so I can load the next order orders with parent_order_number.
func UpdatePlanExecutedOrder(db *sql.DB, planID, orderID string) error {
	updatePlanSql := `UPDATE plans SET last_executed_order_id = $1 WHERE id = $2`
	_, err := db.Exec(updatePlanSql, orderID, planID)

	if err != nil {
		return err
	}

	return nil
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
