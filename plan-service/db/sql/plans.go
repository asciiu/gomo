package sql

import (
	"context"
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

// Retrieve all unexecuted orders for the plan including the last executed order.
func FindPlanWithUnexecutedOrders(db *sql.DB, planID string) (*protoPlan.Plan, error) {
	rows, err := db.Query(`SELECT 
		p.id as plan_id,
		p.user_id,
		p.exchange_name,
		p.market_name,
		p.active_currency_symbol,
		p.active_currency_balance,
		p.last_executed_plan_depth,
		p.last_executed_order_id,
		p.user_plan_number,
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
		o.initial_currency_symbol,
		o.initial_currency_balance,
		o.initial_currency_traded,
		o.initial_currency_remainder,
		o.final_currency_symbol,
		o.final_currency_balance,
		o.order_priority,
		o.order_template_id,
		o.order_type,
		o.side,
		o.limit_price, 
		o.status,
		o.grupo,
		o.created_on,
		o.updated_on,
		t.id as trigger_id,
		t.name,
		t.code,
		array_to_json(t.actions),
		t.index,
		t.title,
		t.triggered,
		t.triggered_price,
		t.triggered_condition,
		t.triggered_timestamp,
		t.trigger_template_id,
		t.created_on,
		t.updated_on
		FROM plans p 
		JOIN orders o on p.id = o.plan_id AND p.last_executed_plan_depth <= o.plan_depth
		JOIN triggers t on o.id = t.order_id
		JOIN user_keys k on o.user_key_id = k.id
		WHERE p.id = $1
		ORDER BY o.plan_depth, o.order_priority, o.id, t.index`, planID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var plan protoPlan.Plan
	var price sql.NullFloat64

	for rows.Next() {
		var trigger protoOrder.Trigger
		var triggerTemplateID sql.NullString
		var finalCurrencySymbol sql.NullString
		var triggeredPrice sql.NullFloat64
		var triggeredCondition sql.NullString
		var triggeredTimestamp sql.NullString
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
			&plan.UserPlanNumber,
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
			&order.InitialCurrencySymbol,
			&order.InitialCurrencyBalance,
			&order.InitialCurrencyTraded,
			&order.InitialCurrencyRemainder,
			&finalCurrencySymbol,
			&order.FinalCurrencyBalance,
			&order.OrderPriority,
			&order.OrderTemplateID,
			&order.OrderType,
			&order.Side,
			&price,
			&order.Status,
			&order.Grupo,
			&order.CreatedOn,
			&order.UpdatedOn,
			&trigger.TriggerID,
			&trigger.Name,
			&trigger.Code,
			&actionsStr,
			&trigger.Index,
			&trigger.Title,
			&trigger.Triggered,
			&triggeredPrice,
			&triggeredCondition,
			&triggeredTimestamp,
			&triggerTemplateID,
			&trigger.CreatedOn,
			&trigger.UpdatedOn)

		if err != nil {
			return nil, err
		}
		if price.Valid {
			order.LimitPrice = price.Float64
		}
		if finalCurrencySymbol.Valid {
			order.FinalCurrencySymbol = finalCurrencySymbol.String
		}
		if triggerTemplateID.Valid {
			trigger.TriggerTemplateID = triggerTemplateID.String
		}
		if triggeredPrice.Valid {
			trigger.TriggeredPrice = triggeredPrice.Float64
		}
		if triggeredCondition.Valid {
			trigger.TriggeredCondition = triggeredCondition.String
		}
		if triggeredTimestamp.Valid {
			trigger.TriggeredTimestamp = triggeredTimestamp.String
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
		p.close_on_complete,
		p.user_plan_number,
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
		o.initial_currency_symbol,
		o.initial_currency_balance,
		o.initial_currency_traded,
		o.order_template_id,
		o.order_type,
		o.side,
		o.limit_price, 
		o.status,
		o.grupo,
		o.created_on,
		o.updated_on,
		t.id as trigger_id,
		t.order_id as tio,
		t.index,
		t.title,
		t.name,
		t.code,
		array_to_json(t.actions),
		t.triggered,
		t.trigger_template_id,
		t.created_on,
		t.updated_on
		FROM plans p 
		JOIN orders o on p.id = o.plan_id AND p.last_executed_order_id = o.parent_order_id
		JOIN triggers t on o.id = t.order_id
		JOIN user_keys k on o.user_key_id = k.id
		WHERE p.status = $1 AND k.status = $2 ORDER BY p.id, o.id, t.index`, plan.Active, key.Verified)

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
			&plan.CloseOnComplete,
			&plan.UserPlanNumber,
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
			&order.UpdatedOn,
			&trigger.TriggerID,
			&trigger.OrderID,
			&trigger.Index,
			&trigger.Title,
			&trigger.Name,
			&trigger.Code,
			&actionsStr,
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
		p.plan_template_id,
		p.user_plan_number,
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
		o.initial_currency_symbol,
		o.initial_currency_balance,
		o.initial_currency_traded,
		o.initial_currency_remainder,
		o.final_currency_symbol,
		o.final_currency_balance,
		o.order_priority,
		o.order_template_id,
		o.order_type,
		o.side,
		o.limit_price, 
		o.status,
		o.grupo,
		o.created_on,
		o.updated_on,
		t.id as trigger_id,
		t.index,
		t.title,
		t.name,
		t.code,
		array_to_json(t.actions),
		t.triggered,
		t.triggered_price,
		t.triggered_condition,
		t.triggered_timestamp,
		t.trigger_template_id,
		t.created_on,
		t.updated_on
		FROM plans p 
		JOIN orders o on p.id = o.plan_id
		JOIN triggers t on o.id = t.order_id
		JOIN user_keys k on o.user_key_id = k.id
		WHERE p.id = $1 AND o.plan_depth BETWEEN $2 AND $3 
		ORDER BY o.plan_depth, o.order_priority, o.id, t.index`, req.PlanID, req.PlanDepth, req.PlanDepth+req.PlanLength)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var plan protoPlan.Plan
	var price sql.NullFloat64

	for rows.Next() {
		var trigger protoOrder.Trigger
		var triggerTemplateID sql.NullString
		var triggeredPrice sql.NullFloat64
		var triggeredCondition sql.NullString
		var triggeredTimestamp sql.NullString
		var finalCurrencySymbol sql.NullString
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
			&plan.PlanTemplateID,
			&plan.UserPlanNumber,
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
			&order.InitialCurrencySymbol,
			&order.InitialCurrencyBalance,
			&order.InitialCurrencyTraded,
			&order.InitialCurrencyRemainder,
			&finalCurrencySymbol,
			&order.FinalCurrencyBalance,
			&order.OrderPriority,
			&order.OrderTemplateID,
			&order.OrderType,
			&order.Side,
			&price,
			&order.Status,
			&order.Grupo,
			&order.CreatedOn,
			&order.UpdatedOn,
			&trigger.TriggerID,
			&trigger.Index,
			&trigger.Title,
			&trigger.Name,
			&trigger.Code,
			&actionsStr,
			&trigger.Triggered,
			&triggeredPrice,
			&triggeredCondition,
			&triggeredTimestamp,
			&triggerTemplateID,
			&trigger.CreatedOn,
			&trigger.UpdatedOn)

		if err != nil {
			return nil, err
		}
		if price.Valid {
			order.LimitPrice = price.Float64
		}
		if finalCurrencySymbol.Valid {
			order.FinalCurrencySymbol = finalCurrencySymbol.String
		}
		if triggeredPrice.Valid {
			trigger.TriggeredPrice = triggeredPrice.Float64
		}
		if triggeredCondition.Valid {
			trigger.TriggeredCondition = triggeredCondition.String
		}
		if triggeredTimestamp.Valid {
			trigger.TriggeredTimestamp = triggeredTimestamp.String
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

// Returns child orders from parentOrderID with plan deets.
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
		p.close_on_complete,
		p.user_plan_number,
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
		o.initial_currency_symbol,
		o.initial_currency_balance,
		o.initial_currency_traded,
		o.initial_currency_remainder,
		o.final_currency_symbol,
		o.final_currency_balance,
		o.order_template_id,
		o.order_type,
		o.side,
		o.limit_price, 
		o.status,
		o.grupo,
		o.created_on,
		o.updated_on,
		t.id as trigger_id,
		t.index,
		t.title,
		t.name,
		t.code,
		array_to_json(t.actions),
		t.triggered,
		t.triggered_price,
		t.triggered_condition,
		t.triggered_timestamp,
		t.trigger_template_id,
		t.created_on,
		t.updated_on
		FROM plans p 
		JOIN orders o on p.id = o.plan_id
		JOIN triggers t on o.id = t.order_id
		JOIN user_keys k on o.user_key_id = k.id
		WHERE p.id = $1 AND o.parent_order_id = $2 AND k.status = $3 
		ORDER BY o.id, t.index`, planID, parentOrderID, key.Verified)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var plan protoPlan.Plan
	var price sql.NullFloat64

	for rows.Next() {
		var trigger protoOrder.Trigger
		var triggerTemplateID sql.NullString
		var triggeredPrice sql.NullFloat64
		var triggeredCondition sql.NullString
		var triggeredTimestamp sql.NullString
		var finalCurrencySymbol sql.NullString
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
			&plan.CloseOnComplete,
			&plan.UserPlanNumber,
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
			&order.InitialCurrencySymbol,
			&order.InitialCurrencyBalance,
			&order.InitialCurrencyTraded,
			&order.InitialCurrencyRemainder,
			&finalCurrencySymbol,
			&order.FinalCurrencyBalance,
			&order.OrderTemplateID,
			&order.OrderType,
			&order.Side,
			&price,
			&order.Status,
			&order.Grupo,
			&order.CreatedOn,
			&order.UpdatedOn,
			&trigger.TriggerID,
			&trigger.Index,
			&trigger.Title,
			&trigger.Name,
			&trigger.Code,
			&actionsStr,
			&trigger.Triggered,
			&triggeredPrice,
			&triggeredCondition,
			&triggeredTimestamp,
			&triggerTemplateID,
			&trigger.CreatedOn,
			&trigger.UpdatedOn)

		if err != nil {
			return nil, err
		}
		if price.Valid {
			order.LimitPrice = price.Float64
		}
		if finalCurrencySymbol.Valid {
			order.FinalCurrencySymbol = finalCurrencySymbol.String
		}
		if triggeredPrice.Valid {
			trigger.TriggeredPrice = triggeredPrice.Float64
		}
		if triggeredCondition.Valid {
			trigger.TriggeredCondition = triggeredCondition.String
		}
		if triggeredTimestamp.Valid {
			trigger.TriggeredTimestamp = triggeredTimestamp.String
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

	// no child orders
	if plan.PlanID == "" {
		return nil, sql.ErrNoRows
	}

	return &plan, nil
}

func FindUserPlans(db *sql.DB, userID string, page, pageSize uint32) (*protoPlan.PlansPage, error) {

	var count uint32
	queryCount := `SELECT count(*) FROM plans WHERE user_id = $1`
	err := db.QueryRow(queryCount, userID).Scan(&count)

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
			user_plan_number,
			status,
			created_on,
			updated_on
		FROM plans 
		WHERE user_id = $1 ORDER BY created_on OFFSET $2 LIMIT $3`

	rows, err := db.Query(query, userID, page, pageSize)
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
			&plan.UserPlanNumber,
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
			user_plan_number,
			status,
			created_on,
			updated_on
		FROM plans 
		WHERE user_id = $1 AND status = $2 ORDER BY created_on OFFSET $3 LIMIT $4`

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
			&plan.UserPlanNumber,
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
			user_plan_number,
			status,
			created_on,
			updated_on
		FROM plans 
		WHERE user_id = $1 AND status = $2 AND exchange_name = $3 ORDER BY created_on OFFSET $4 LIMIT $5`

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
			&plan.UserPlanNumber,
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
			user_plan_number,
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
			&plan.UserPlanNumber,
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
		user_plan_number,
		status, 
		created_on,
		updated_on) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`)

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
		newPlan.UserPlanNumber,
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

func UpdatePlanStatus(db *sql.DB, planID, status string) error {
	stmt := `
	UPDATE plans 
	SET 
		status = $1 
	WHERE
		id = $2`

	_, err := db.Exec(stmt, status, planID)

	return err
}

// func UpdatePlanBalances(db *sql.DB, planID string, base, currency float64) error {
// 	sqlStatement := `UPDATE plans SET base_balance = $1, currency_balance = $2 WHERE id = $3`

// 	_, error := db.Exec(sqlStatement, base, currency, planID)
// 	return error
// }

func UpdatePlanContext(db *sql.DB, planID, executedOrderID, exchange, marketName, symbol string, activeBalance float64, planDepth int32) error {
	stmt := `
		UPDATE plans 
		SET 
			last_executed_order_id = $1,
			last_executed_plan_depth = $2,
			active_currency_symbol = $3,
			active_currency_balance = $4,
			market_name = $5,
			exchange_name = $6 
		WHERE
			id = $7`

	_, err := db.Exec(stmt, executedOrderID, planDepth, symbol, activeBalance, marketName, exchange, planID)

	return err
}

func UpdatePlanContextTxn(txn *sql.Tx, ctx context.Context, planID, symbol, exchange, marketName string, activeBalance float64) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE plans 
		SET 
			active_currency_symbol = $1,
			active_currency_balance = $2,
			market_name = $3,
			exchange_name = $4 
		WHERE
			id = $5`,
		symbol, activeBalance, marketName, exchange, planID)

	return err
}

func UpdatePlanStatusTxn(txn *sql.Tx, ctx context.Context, planID, status string) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE plans 
		SET 
			status = $1 
		WHERE
			id = $2`,
		status, planID)

	return err
}

func UpdatePlanCloseOnCompleteTxn(txn *sql.Tx, ctx context.Context, planID string, closeOnComplete bool) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE plans 
		SET 
			close_on_complete = $1 
		WHERE
			id = $2`,
		closeOnComplete, planID)

	return err
}

func UpdatePlanTemplateTxn(txn *sql.Tx, ctx context.Context, planID, template string) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE plans 
		SET 
			plan_template_id = $1 
		WHERE
			id = $2`,
		template, planID)

	return err
}

func UpdatePlanTimestampTxn(txn *sql.Tx, ctx context.Context, planID, timestamp string) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE plans 
		SET 
			updated_on = $1 
		WHERE
			id = $2`,
		timestamp, planID)

	return err
}

// this should only be called when updating the an executed plan since the state of plan
// is dictated by the first root order
func UpdatePlanTxn(txn *sql.Tx, ctx context.Context, plan *protoPlan.Plan) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE plans 
		SET 
			exchange_name = $1,
			market_name = $2,
			active_currency_symbol = $3,
			active_currency_balance = $4,
			plan_template_id = $5,
			close_on_complete = $6,
			status = $7 
		WHERE
			id = $8`,
		plan.Exchange,
		plan.MarketName,
		plan.ActiveCurrencySymbol,
		plan.ActiveCurrencyBalance,
		plan.PlanTemplateID,
		plan.CloseOnComplete,
		plan.Status,
		plan.PlanID)

	return err
}

func FindUserPlanCount(db *sql.DB, userID string) uint64 {
	var count uint64
	queryCount := `SELECT count(*) FROM plans WHERE user_id = $1`
	err := db.QueryRow(queryCount, userID).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
