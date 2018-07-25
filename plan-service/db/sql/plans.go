package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/asciiu/gomo/common/constants/key"
	"github.com/asciiu/gomo/common/constants/plan"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
)

// DeletePlan is a hard delete
func DeletePlan(db *sql.DB, planID string) error {
	_, err := db.Exec("DELETE FROM plans WHERE id = $1", planID)
	return err
}

// FindPlanSummary will load the order with the active order number as its only order
// this is so we can examine the current state of the plan
// func FindPlanSummary(db *sql.DB, planID string) (*protoPlan.Plan, error) {
// 	var plan protoPlan.Plan
// 	var planTemp sql.NullString
// 	var orderTemp sql.NullString
// 	var price sql.NullFloat64

// 	// note: always assume that a plan has at least one order.
// 	// a plan with no order should be removed from the system
// 	rows, err := db.Query(`SELECT
// 		p.id,
// 		p.plan_template_id,
// 		p.user_id,
// 		p.user_key_id,
// 		k.description,
// 		p.last_executed_order_id,
// 		p.exchange_name,
// 		p.market_name,
// 		p.base_balance,
// 		p.currency_balance,
// 		p.status,
// 		o.id,
// 		o.balance_percent,
// 		o.side,
// 		o.parent_order_id,
// 		o.order_type,
// 		o.order_template_id,
// 		o.limit_price,
// 		o.status
// 		FROM plans p
// 		JOIN orders o on p.id = o.plan_id and p.last_executed_order_id = o.parent_order_id
// 		JOIN user_keys k on p.user_key_id = k.id
// 		WHERE p.id = $1`, planID)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var order protoOrder.Order
// 		err := rows.Scan(
// 			&plan.PlanID,
// 			&planTemp,
// 			&plan.UserID,
// 			&plan.KeyID,
// 			&plan.KeyDescription,
// 			&plan.LastExecutedOrderID,
// 			//&plan.Exchange,
// 			//&plan.MarketName,
// 			//&plan.BaseBalance,
// 			//&plan.CurrencyBalance,
// 			&plan.Status,
// 			&order.OrderID,
// 			//&order.BalancePercent,
// 			&order.Side,
// 			&order.ParentOrderID,
// 			&order.OrderType,
// 			&orderTemp,
// 			&price,
// 			&order.Status)

// 		if err != nil {
// 			return nil, err
// 		}
// 		if price.Valid {
// 			order.LimitPrice = price.Float64
// 		}
// 		if orderTemp.Valid {
// 			order.OrderTemplateID = orderTemp.String
// 		}
// 		if planTemp.Valid {
// 			plan.PlanTemplateID = planTemp.String
// 		}
// 		plan.Orders = append(plan.Orders, &order)

// 	}
// 	return &plan, nil
// }

// FindPlanWithPagedOrders ...
// func FindPlanWithPagedOrders(db *sql.DB, req *protoPlan.GetUserPlanRequest) (*protoPlan.PlanWithPagedOrders, error) {

// 	var planPaged protoPlan.PlanWithPagedOrders
// 	var page protoOrder.OrdersPage

// 	var planTemp sql.NullString
// 	var orderTemp sql.NullString
// 	var price sql.NullFloat64

// 	var count uint32
// 	queryCount := `SELECT count(*) FROM orders WHERE plan_id = $1`
// 	err := db.QueryRow(queryCount, req.PlanID).Scan(&count)

// 	if count == 0 {
// 		return nil, sql.ErrNoRows
// 	}

// 	rows, err := db.Query(`SELECT
// 		p.id,
// 		p.plan_template_id,
// 		p.user_id,
// 		p.user_key_id,
// 		k.api_key,
// 		k.description,
// 		p.exchange_name,
// 		p.market_name,
// 		p.base_balance,
// 		p.currency_balance,
// 		p.status,
// 		p.created_on,
// 		p.updated_on,
// 		o.id,
// 		o.balance_percent,
// 		o.side,
// 		o.order_type,
// 		o.order_template_id,
// 		o.limit_price,
// 		o.status,
// 		o.created_on,
// 		o.updated_on,
// 		t.id as trigger_id,
// 		t.name,
// 		t.code,
// 		array_to_json(t.actions),
// 		t.trigger_number,
// 		t.triggered,
// 		t.trigger_template_id,
// 		t.created_on,
// 		t.updated_on
// 		FROM plans p
// 		JOIN orders o on p.id = o.plan_id
// 		JOIN triggers t on o.id = t.order_id
// 		JOIN user_keys k on p.user_key_id = k.id
// 		WHERE p.id = $1 ORDER BY o.plan_depth, t.trigger_number
// 		OFFSET $2 LIMIT $3`, req.PlanID, req.Page, req.PageSize)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var order protoOrder.Order
// 		var trigger protoOrder.Trigger
// 		var triggerTempID sql.NullString
// 		var actionsStr string

// 		err := rows.Scan(
// 			&planPaged.PlanID,
// 			&planTemp,
// 			&planPaged.UserID,
// 			&planPaged.KeyID,
// 			&planPaged.Key,
// 			&planPaged.KeyDescription,
// 			&planPaged.Exchange,
// 			&planPaged.MarketName,
// 			&planPaged.BaseBalance,
// 			&planPaged.CurrencyBalance,
// 			&planPaged.Status,
// 			&planPaged.CreatedOn,
// 			&planPaged.UpdatedOn,
// 			&order.OrderID,
// 			//&order.BalancePercent,
// 			&order.Side,
// 			&order.OrderType,
// 			&orderTemp,
// 			&price,
// 			&order.Status,
// 			&order.CreatedOn,
// 			&order.UpdatedOn,
// 			&trigger.TriggerID,
// 			&trigger.Name,
// 			&trigger.Code,
// 			&actionsStr,
// 			&trigger.TriggerNumber,
// 			&trigger.Triggered,
// 			&triggerTempID,
// 			&trigger.CreatedOn,
// 			&trigger.UpdatedOn,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}
// 		if price.Valid {
// 			order.LimitPrice = price.Float64
// 		}
// 		if orderTemp.Valid {
// 			order.OrderTemplateID = orderTemp.String
// 		}
// 		if planTemp.Valid {
// 			planPaged.PlanTemplateID = planTemp.String
// 		}
// 		if triggerTempID.Valid {
// 			trigger.TriggerTemplateID = triggerTempID.String
// 		}
// 		if err := json.Unmarshal([]byte(actionsStr), &trigger.Actions); err != nil {
// 			return nil, err
// 		}
// 		if err := json.Unmarshal([]byte(trigger.Code), &trigger.Code); err != nil {
// 			return nil, err
// 		}

// 		trigger.OrderID = order.OrderID

// 		var found = false
// 		for _, or := range page.Orders {
// 			if or.OrderID == order.OrderID {
// 				or.Triggers = append(or.Triggers, &trigger)
// 				found = true
// 				break
// 			}
// 		}

// 		if !found {
// 			order.Triggers = append(order.Triggers, &trigger)
// 			page.Orders = append(page.Orders, &order)
// 		}
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}
// 	page.Page = req.Page
// 	page.PageSize = req.PageSize
// 	page.Total = count
// 	planPaged.OrdersPage = &page

// 	return &planPaged, nil
// }

// Find all active orders in the DB. This wil load the keys for each order.
// Returns active plans that have a verified key only.
func FindActivePlans(db *sql.DB) ([]*protoPlan.Plan, error) {
	rows, err := db.Query(`SELECT 
		p.id as plan_id,
		p.user_id,
		p.exchange_name,
		p.market_name,
		p.active_currency_symbol,
		p.active_currency_balance,
		p.last_executed_plan_depth,
		p.last_executed_order_id,
		p.status,
		k.api_key,
		k.secret,
		k.description,
		o.id as order_id,
		o.parent_order_id,
		o.plan_id,
		o.plan_depth,
		o.exchange_name,
		o.market_name,
		o.active_currency_symbol,
		o.active_currency_balance,
		o.active_currency_traded,
		o.order_template_id,
		o.order_type,
		o.side,
		o.limit_price, 
		o.status,
		o.created_on,
		o.updated_on,
		t.id as trigger_id,
		t.order_id as tio,
		t.name,
		t.code,
		array_to_json(t.actions),
		t.trigger_number,
		t.triggered,
		t.trigger_template_id,
		t.created_on,
		t.updated_on
		FROM plans p 
		JOIN orders o on p.id = o.plan_id AND p.last_executed_order_id = o.parent_order_id
		JOIN triggers t on o.id = t.order_id
		JOIN user_keys k on o.user_key_id = k.id
		WHERE p.status = $1 AND k.status = $2 ORDER BY p.id, o.id, t.trigger_number`, plan.Active, key.Verified)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	activePlans := make([]*protoPlan.Plan, 0)

	for rows.Next() {
		var price sql.NullFloat64
		var trigger protoOrder.Trigger
		var triggerTemplateID sql.NullString
		var actionsStr string
		var plan protoPlan.Plan
		var order protoOrder.Order

		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.Exchange,
			&plan.MarketName,
			&plan.ActiveCurrencySymbol,
			&plan.ActiveCurrencyBalance,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&plan.Status,
			&order.KeyPublic,
			&order.KeySecret,
			&order.KeyDescription,
			&order.OrderID,
			&order.ParentOrderID,
			&order.PlanID,
			&order.PlanDepth,
			&order.Exchange,
			&order.MarketName,
			&order.ActiveCurrencySymbol,
			&order.ActiveCurrencyBalance,
			&order.ActiveCurrencyTraded,
			&order.OrderTemplateID,
			&order.OrderType,
			&order.Side,
			&price,
			&order.Status,
			&order.CreatedOn,
			&order.UpdatedOn,
			&trigger.TriggerID,
			&trigger.OrderID,
			&trigger.Name,
			&trigger.Code,
			&actionsStr,
			&trigger.TriggerNumber,
			&trigger.Triggered,
			&triggerTemplateID,
			&trigger.CreatedOn,
			&trigger.UpdatedOn)

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

		// add order to plan if planID matches
		var foundPlan = false
		for _, p := range activePlans {
			// if plan match activePlans
			if p.PlanID == plan.PlanID {

				// add trigger to existing order
				var foundOrder = false
				for _, o := range p.Orders {
					if o.OrderID == order.OrderID {
						// append to the existing order's triggers
						o.Triggers = append(o.Triggers, &trigger)
						foundOrder = true
						break
					}
				}

				if !foundOrder {
					order.Triggers = append(order.Triggers, &trigger)
					p.Orders = append(p.Orders, &order)
				}

				foundPlan = true
				break
			}
		}

		if !foundPlan {
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

func FindPlanOrders(db *sql.DB, req *protoPlan.GetUserPlanRequest) (*protoPlan.Plan, error) {
	rows, err := db.Query(`SELECT 
		p.id as plan_id,
		p.user_id,
		p.exchange_name,
		p.market_name,
		p.active_currency_symbol,
		p.active_currency_balance,
		p.last_executed_plan_depth,
		p.last_executed_order_id,
		p.status,
		p.created_on,
		p.updated_on,
		k.api_key,
		k.secret,
		k.description,
		o.id as order_id,
		o.parent_order_id,
		o.plan_id,
		o.plan_depth,
		o.exchange_name,
		o.market_name,
		o.active_currency_symbol,
		o.active_currency_balance,
		o.active_currency_traded,
		o.order_priority,
		o.order_template_id,
		o.order_type,
		o.side,
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
		JOIN user_keys k on o.user_key_id = k.id
		WHERE p.id = $1 AND o.plan_depth BETWEEN $2 AND $3 
		ORDER BY o.plan_depth, o.order_priority, o.id, t.trigger_number`, req.PlanID, req.PlanDepth, req.PlanDepth+req.PlanLength)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var plan protoPlan.Plan
	var price sql.NullFloat64

	for rows.Next() {
		var trigger protoOrder.Trigger
		var triggerTemplateID sql.NullString
		var actionsStr string
		var order protoOrder.Order

		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.Exchange,
			&plan.MarketName,
			&plan.ActiveCurrencySymbol,
			&plan.ActiveCurrencyBalance,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&plan.Status,
			&plan.CreatedOn,
			&plan.UpdatedOn,
			&order.KeyPublic,
			&order.KeySecret,
			&order.KeyDescription,
			&order.OrderID,
			&order.ParentOrderID,
			&order.PlanID,
			&order.PlanDepth,
			&order.Exchange,
			&order.MarketName,
			&order.ActiveCurrencySymbol,
			&order.ActiveCurrencyBalance,
			&order.ActiveCurrencyTraded,
			&order.OrderPriority,
			&order.OrderTemplateID,
			&order.OrderType,
			&order.Side,
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
			&triggerTemplateID,
			&trigger.CreatedOn,
			&trigger.UpdatedOn)

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
		trigger.OrderID = order.OrderID

		// loop through the orders and append trigger to existing
		var found = false
		for _, o := range plan.Orders {
			if o.OrderID == order.OrderID {
				o.Triggers = append(o.Triggers, &trigger)
				found = true
				break
			}
		}
		if !found {
			order.Triggers = append(order.Triggers, &trigger)
			plan.Orders = append(plan.Orders, &order)
		}
	}

	if plan.PlanID == "" {
		return nil, sql.ErrNoRows
	}

	return &plan, nil
}

// Returns plan with specific order that has a verified key.
func FindChildOrders(db *sql.DB, planID, parentOrderID string) (*protoPlan.Plan, error) {
	rows, err := db.Query(`SELECT 
		p.id as plan_id,
		p.user_id,
		p.exchange_name,
		p.market_name,
		p.active_currency_symbol,
		p.active_currency_balance,
		p.last_executed_plan_depth,
		p.last_executed_order_id,
		p.status,
		k.api_key,
		k.secret,
		k.description,
		o.id as order_id,
		o.parent_order_id,
		o.plan_id,
		o.plan_depth,
		o.exchange_name,
		o.market_name,
		o.active_currency_symbol,
		o.active_currency_balance,
		o.active_currency_traded,
		o.order_template_id,
		o.order_type,
		o.side,
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
		JOIN user_keys k on o.user_key_id = k.id
		WHERE p.id = $1 AND o.parent_order_id = $2 AND k.status = $3 
		ORDER BY o.id, t.trigger_number`, planID, parentOrderID, key.Verified)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var plan protoPlan.Plan
	var price sql.NullFloat64

	for rows.Next() {
		var trigger protoOrder.Trigger
		var triggerTemplateID sql.NullString
		var actionsStr string
		var order protoOrder.Order

		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.Exchange,
			&plan.MarketName,
			&plan.ActiveCurrencySymbol,
			&plan.ActiveCurrencyBalance,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&plan.Status,
			&order.KeyPublic,
			&order.KeySecret,
			&order.KeyDescription,
			&order.OrderID,
			&order.ParentOrderID,
			&order.PlanID,
			&order.PlanDepth,
			&order.Exchange,
			&order.MarketName,
			&order.ActiveCurrencySymbol,
			&order.ActiveCurrencyBalance,
			&order.ActiveCurrencyTraded,
			&order.OrderTemplateID,
			&order.OrderType,
			&order.Side,
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
			&triggerTemplateID,
			&trigger.CreatedOn,
			&trigger.UpdatedOn)

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
		trigger.OrderID = order.OrderID

		// loop through the orders and append trigger to existing
		var found = false
		for _, o := range plan.Orders {
			if o.OrderID == order.OrderID {
				o.Triggers = append(o.Triggers, &trigger)
				found = true
				break
			}
		}
		if !found {
			order.Triggers = append(order.Triggers, &trigger)
			plan.Orders = append(plan.Orders, &order)
		}
	}

	return &plan, nil
}

func FindUserPlansWithStatus(db *sql.DB, userID, status string, page, pageSize uint32) (*protoPlan.PlansPage, error) {

	var count uint32
	queryCount := `SELECT count(*) FROM plans WHERE user_id = $1 AND status = $2`
	err := db.QueryRow(queryCount, userID, status).Scan(&count)

	plans := make([]*protoPlan.Plan, 0)
	query := `SELECT 
			id,
			exchange_name,
			market_name,
			active_currency_symbol,
			active_currency_balance,
			last_executed_plan_depth,
			last_executed_order_id,
			plan_template_id,
			close_on_complete,
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
			&plan.Exchange,
			&plan.MarketName,
			&plan.ActiveCurrencySymbol,
			&plan.ActiveCurrencyBalance,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&planTemplateID,
			&plan.CloseOnComplete,
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
func FindUserExchangePlansWithStatus(db *sql.DB, userID, status, exchange string, page, pageSize uint32) (*protoPlan.PlansPage, error) {

	var count uint32
	queryCount := `SELECT count(*) FROM plans WHERE user_id = $1 AND status = $2 AND exchange_name = $3`
	err := db.QueryRow(queryCount, userID, status, exchange).Scan(&count)

	plans := make([]*protoPlan.Plan, 0)
	query := `SELECT 
			id,
			exchange_name,
			market_name,
			active_currency_symbol,
			active_currency_balance,
			last_executed_plan_depth,
			last_executed_order_id,
			plan_template_id,
			close_on_complete,
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
			&plan.Exchange,
			&plan.MarketName,
			&plan.ActiveCurrencySymbol,
			&plan.ActiveCurrencyBalance,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&planTemplateID,
			&plan.CloseOnComplete,
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
			exchange_name,
			market_name,
			active_currency_symbol,
			active_currency_balance,
			last_executed_plan_depth,
			last_executed_order_id,
			plan_template_id,
			close_on_complete,
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
			&plan.Exchange,
			&plan.MarketName,
			&plan.ActiveCurrencySymbol,
			&plan.ActiveCurrencyBalance,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&planTemplateID,
			&plan.CloseOnComplete,
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
// TODO the orders in the plan request must be pre sorted by depth. Root node is assumed to be at index 0.
func InsertPlan(db *sql.DB, newPlan *protoPlan.Plan) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}

	sqlStatement, err := txn.Prepare(`insert into plans (
		id, 
		user_id, 
		exchange_name,
		market_name,
		active_currency_symbol,
		active_currency_balance,
		last_executed_plan_depth,
		last_executed_order_id,
		plan_template_id,
		close_on_complete,
		status, 
		created_on,
		updated_on) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`)

	if err != nil {
		txn.Rollback()
		return errors.New("prepare insert plan failed: " + err.Error())
	}
	_, err = txn.Stmt(sqlStatement).Exec(
		newPlan.PlanID,
		newPlan.UserID,
		newPlan.Exchange,
		newPlan.MarketName,
		newPlan.ActiveCurrencySymbol,
		newPlan.ActiveCurrencyBalance,
		newPlan.LastExecutedPlanDepth,
		newPlan.LastExecutedOrderID,
		newPlan.PlanTemplateID,
		newPlan.CloseOnComplete,
		newPlan.Status,
		newPlan.CreatedOn,
		newPlan.UpdatedOn)

	if err != nil {
		txn.Rollback()
		return errors.New("plan insert failed: " + err.Error())
	}

	if err := InsertOrders(txn, newPlan.Orders); err != nil {
		txn.Rollback()
		return errors.New("insert orders failed: " + err.Error())
	}

	newTriggers := make([]*protoOrder.Trigger, 0, len(newPlan.Orders))
	for _, o := range newPlan.Orders {
		for _, t := range o.Triggers {
			newTriggers = append(newTriggers, t)
		}
	}

	if err := InsertTriggers(txn, newTriggers); err != nil {
		txn.Rollback()
		return errors.New("bulk triggers failed: " + err.Error())
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
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
	RETURNING id, plan_template_id, user_id, status`

	var p protoPlan.Plan
	err := db.QueryRow(sqlStatement, status, planID).
		Scan(&p.PlanID,
			&p.PlanTemplateID,
			&p.UserID,
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
	status`

	var p protoPlan.Plan
	err := db.QueryRow(sqlStatement, baseBalance, planID).
		Scan(&p.PlanID,
			&p.PlanTemplateID,
			&p.UserID,
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
	status`

	var p protoPlan.Plan
	err := db.QueryRow(sqlStatement, currencyBalance, planID).
		Scan(&p.PlanID,
			&p.PlanTemplateID,
			&p.UserID,
			&p.Status,
		)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
