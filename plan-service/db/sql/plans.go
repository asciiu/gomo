package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	constAccount "github.com/asciiu/gomo/account-service/constants"
	constPlan "github.com/asciiu/gomo/plan-service/constants"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
)

/*
This file has these functions:
	DeletePlan
	InsertPlan
	FindActivePlans
	FindChildOrders
	FindParentAndChildren
	FindPlanOrders
	FindPlanWithUnexecutedOrders
	FindUserExchangePlansWithStatus
	FindUserMarketPlansWithStatus
	FindUserPlans
	FindUserPlansWithStatus
	FindUserPlanCount
	UpdatePlanBaseCurrencyTxn
	UpdatePlanCloseOnCompleteTxn
	UpdatePlanContext
	UpdatePlanContextTxn
	UpdatePlanExecutedOrder
	UpdatePlanStatus
	UpdatePlanStatusTxn
	UpdatePlanTemplateTxn
	UpdatePlanTimestampTxn
	UpdatePlanTitleTxn
	UpdatePlanTotalDepthTxn
	UpdatePlanTxn
*/

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
		p.title,
		p.total_depth,
		p.exchange_name,
		p.user_currency_symbol,
		p.active_currency_symbol,
		p.active_currency_balance,
		p.initial_currency_symbol,
		p.initial_currency_balance,
		p.initial_timestamp,
		p.last_executed_plan_depth,
		p.last_executed_order_id,
		p.user_plan_number,
		p.status,
		p.close_on_complete,
		p.created_on,
		p.updated_on,
		a.id as account_id,
		a.key_public,
		a.key_secret,
		a.description,
		o.id as order_id,
		o.parent_order_id,
		o.plan_id,
		o.plan_depth,
		o.exchange_name,
		o.exchange_order_id,
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
		JOIN accounts a on a.id = o.account_id
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
		var initialTimestamp sql.NullString
		var exchangeOrderID sql.NullString
		var actionsStr string
		var order protoOrder.Order

		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.Title,
			&plan.TotalDepth,
			&plan.Exchange,
			&plan.UserCurrencySymbol,
			&plan.CommittedCurrencySymbol,
			&plan.CommittedCurrencyAmount,
			&plan.InitialCurrencySymbol,
			&plan.InitialCurrencyBalance,
			&initialTimestamp,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&plan.UserPlanNumber,
			&plan.Status,
			&plan.CloseOnComplete,
			&plan.CreatedOn,
			&plan.UpdatedOn,
			&order.AccountID,
			&order.KeyPublic,
			&order.KeySecret,
			&order.KeyDescription,
			&order.OrderID,
			&order.ParentOrderID,
			&order.PlanID,
			&order.PlanDepth,
			&order.Exchange,
			&exchangeOrderID,
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
		if exchangeOrderID.Valid {
			order.ExchangeOrderID = exchangeOrderID.String
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
		if initialTimestamp.Valid {
			plan.InitialTimestamp = initialTimestamp.String
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
		p.title,
		p.total_depth,
		p.exchange_name,
		p.user_currency_symbol,
		p.active_currency_symbol,
		p.active_currency_balance,
		p.initial_currency_symbol,
		p.initial_currency_balance,
		p.initial_timestamp,
		p.last_executed_plan_depth,
		p.last_executed_order_id,
		p.reference_price,
		p.close_on_complete,
		p.user_plan_number,
		p.status,
		a.id as account_id,
		a.account_type,
		a.key_public,
		a.key_secret,
		a.description,
		o.id as order_id,
		o.parent_order_id,
		o.plan_id,
		o.plan_depth,
		o.exchange_name,
		o.exchange_order_id,
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
		JOIN accounts a on o.account_id = a.id
		WHERE p.status = $1 AND a.status = $2 ORDER BY p.id, o.id, t.index`, constPlan.Active, constAccount.AccountValid)

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
		var initialTimestamp sql.NullString
		var exchangeOrderID sql.NullString
		var plan protoPlan.Plan
		var order protoOrder.Order
		var referencePrice sql.NullFloat64

		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.Title,
			&plan.TotalDepth,
			&plan.Exchange,
			&plan.UserCurrencySymbol,
			&plan.CommittedCurrencySymbol,
			&plan.CommittedCurrencyAmount,
			&plan.InitialCurrencySymbol,
			&plan.InitialCurrencyBalance,
			&initialTimestamp,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&referencePrice,
			&plan.CloseOnComplete,
			&plan.UserPlanNumber,
			&plan.Status,
			&order.AccountID,
			&order.AccountType,
			&order.KeyPublic,
			&order.KeySecret,
			&order.KeyDescription,
			&order.OrderID,
			&order.ParentOrderID,
			&order.PlanID,
			&order.PlanDepth,
			&order.Exchange,
			&exchangeOrderID,
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
		if exchangeOrderID.Valid {
			order.ExchangeOrderID = exchangeOrderID.String
		}
		if initialTimestamp.Valid {
			plan.InitialTimestamp = initialTimestamp.String
		}
		if price.Valid {
			order.LimitPrice = price.Float64
		}
		if err := json.Unmarshal([]byte(actionsStr), &trigger.Actions); err != nil {
			return nil, err
		}
		if referencePrice.Valid {
			plan.ReferencePrice = referencePrice.Float64
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
		p.title,
		p.total_depth,
		p.exchange_name,
		p.user_currency_symbol,
		p.user_currency_balance_at_init,
		p.active_currency_symbol,
		p.active_currency_balance,
		p.initial_currency_symbol,
		p.initial_currency_balance,
		p.initial_timestamp,
		p.last_executed_plan_depth,
		p.last_executed_order_id,
		p.plan_template_id,
		p.user_plan_number,
		p.status,
		p.close_on_complete,
		p.created_on,
		p.updated_on,
		a.id as account_id,
		a.key_public,
		a.key_secret,
		a.description,
		o.id as order_id,
		o.parent_order_id,
		o.plan_id,
		o.plan_depth,
		o.exchange_name,
		o.exchange_order_id,
		o.exchange_price,
		o.exchange_time,
		o.market_name,
		o.fee_currency_symbol,
		o.fee_currency_amount,
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
		o.errors,
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
		JOIN accounts a on o.account_id = a.id
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
		var initialTimestamp sql.NullString
		var orderErr sql.NullString
		var exchangePrice sql.NullFloat64
		var exchangeTime sql.NullString
		var feeSymbol sql.NullString
		var feeAmount sql.NullFloat64
		var userBalanceInit sql.NullFloat64
		var finalCurrencySymbol sql.NullString
		var exchangeOrderID sql.NullString
		var actionsStr string
		var order protoOrder.Order

		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.Title,
			&plan.TotalDepth,
			&plan.Exchange,
			&plan.UserCurrencySymbol,
			&userBalanceInit,
			&plan.CommittedCurrencySymbol,
			&plan.CommittedCurrencyAmount,
			&plan.InitialCurrencySymbol,
			&plan.InitialCurrencyBalance,
			&initialTimestamp,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&plan.PlanTemplateID,
			&plan.UserPlanNumber,
			&plan.Status,
			&plan.CloseOnComplete,
			&plan.CreatedOn,
			&plan.UpdatedOn,
			&order.AccountID,
			&order.KeyPublic,
			&order.KeySecret,
			&order.KeyDescription,
			&order.OrderID,
			&order.ParentOrderID,
			&order.PlanID,
			&order.PlanDepth,
			&order.Exchange,
			&exchangeOrderID,
			&exchangePrice,
			&exchangeTime,
			&order.MarketName,
			&feeSymbol,
			&feeAmount,
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
			&orderErr,
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
			&trigger.UpdatedOn,
		)

		if err != nil {
			return nil, err
		}
		if exchangeOrderID.Valid {
			order.ExchangeOrderID = exchangeOrderID.String
		}
		if initialTimestamp.Valid {
			plan.InitialTimestamp = initialTimestamp.String
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
		if feeSymbol.Valid {
			order.FeeCurrencySymbol = feeSymbol.String
		}
		if feeAmount.Valid {
			order.FeeCurrencyAmount = feeAmount.Float64
		}
		if exchangePrice.Valid {
			order.ExchangePrice = exchangePrice.Float64
		}
		if exchangeTime.Valid {
			order.ExchangeTime = exchangeTime.String
		}
		if orderErr.Valid {
			order.Errors = orderErr.String
		}
		if userBalanceInit.Valid {
			plan.UserCurrencyBalanceAtInit = userBalanceInit.Float64
		}
		if err := json.Unmarshal([]byte(actionsStr), &trigger.Actions); err != nil {
			return nil, err
		}
		trigger.OrderID = order.OrderID

		// loop through the orders and append trigger to existing
		var foundOrder = false
		for _, o := range plan.Orders {
			if o.OrderID == order.OrderID {
				// add trigger to orders
				var foundTrigger = false
				for _, t := range o.Triggers {
					if t.TriggerID == trigger.TriggerID {
						foundTrigger = true
						break
					}
				}

				// new trigger for this order append it!
				if !foundTrigger {
					o.Triggers = append(o.Triggers, &trigger)
				}

				foundOrder = true
				break
			}
		}

		if !foundOrder {
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
		p.title,
		p.total_depth,
		p.exchange_name,
		p.user_currency_symbol,
		p.active_currency_symbol,
		p.active_currency_balance,
		p.initial_currency_symbol,
		p.initial_currency_balance,
		p.initial_timestamp,
		p.last_executed_plan_depth,
		p.last_executed_order_id,
		p.reference_price,
		p.close_on_complete,
		p.user_plan_number,
		p.status,
		p.created_on,
		p.updated_on,
		a.id as account_id,
		a.account_type,
		a.key_public,
		a.key_secret,
		a.description,
		o.id as order_id,
		o.parent_order_id,
		o.plan_id,
		o.plan_depth,
		o.exchange_name,
		o.exchange_order_id,
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
		JOIN accounts a on o.account_id = a.id
		WHERE p.id = $1 AND o.parent_order_id = $2 AND a.status = $3 
		ORDER BY o.id, t.index`, planID, parentOrderID, constAccount.AccountValid)

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
		var initialTimestamp sql.NullString
		var exchangeOrderID sql.NullString
		var actionsStr string
		var order protoOrder.Order
		var referencePrice sql.NullFloat64

		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.Title,
			&plan.TotalDepth,
			&plan.Exchange,
			&plan.UserCurrencySymbol,
			&plan.CommittedCurrencySymbol,
			&plan.CommittedCurrencyAmount,
			&plan.InitialCurrencySymbol,
			&plan.InitialCurrencyBalance,
			&initialTimestamp,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&referencePrice,
			&plan.CloseOnComplete,
			&plan.UserPlanNumber,
			&plan.Status,
			&plan.CreatedOn,
			&plan.UpdatedOn,
			&order.AccountID,
			&order.AccountType,
			&order.KeyPublic,
			&order.KeySecret,
			&order.KeyDescription,
			&order.OrderID,
			&order.ParentOrderID,
			&order.PlanID,
			&order.PlanDepth,
			&order.Exchange,
			&exchangeOrderID,
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
		if exchangeOrderID.Valid {
			order.ExchangeOrderID = exchangeOrderID.String
		}
		if initialTimestamp.Valid {
			plan.InitialTimestamp = initialTimestamp.String
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
		if referencePrice.Valid {
			plan.ReferencePrice = referencePrice.Float64
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

// Includes last executed and next order of executed orders.
func FindParentAndChildren(db *sql.DB, planID, parentOrderID string) (*protoPlan.Plan, error) {
	rows, err := db.Query(`SELECT 
		p.id as plan_id,
		p.user_id,
		p.title,
		p.total_depth,
		p.exchange_name,
		p.user_currency_symbol,
		p.active_currency_symbol,
		p.active_currency_balance,
		p.initial_currency_symbol,
		p.initial_currency_balance,
		p.initial_timestamp,
		p.last_executed_plan_depth,
		p.last_executed_order_id,
		p.close_on_complete,
		p.user_plan_number,
		p.status,
		p.created_on,
		p.updated_on,
		a.id as account_id,
		a.key_public,
		a.key_secret,
		a.description,
		o.id as order_id,
		o.parent_order_id,
		o.plan_id,
		o.plan_depth,
		o.exchange_name,
		o.exchange_order_id,
		o.exchange_price,
		o.exchange_time,
		o.market_name,
		o.fee_currency_symbol,
		o.fee_currency_amount,
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
		o.errors,
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
		JOIN accounts a on o.account_id = a.id
		WHERE p.id = $1 AND (o.parent_order_id = $2 or o.id = $2)
		ORDER BY o.plan_depth, o.order_priority, o.id, t.index`, planID, parentOrderID)

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
		var initialTimestamp sql.NullString
		var orderErr sql.NullString
		var exchangePrice sql.NullFloat64
		var exchangeTime sql.NullString
		var feeSymbol sql.NullString
		var feeAmount sql.NullFloat64
		var finalCurrencySymbol sql.NullString
		var exchangeOrderID sql.NullString
		var actionsStr string
		var order protoOrder.Order

		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.Title,
			&plan.TotalDepth,
			&plan.Exchange,
			&plan.UserCurrencySymbol,
			&plan.CommittedCurrencySymbol,
			&plan.CommittedCurrencyAmount,
			&plan.InitialCurrencySymbol,
			&plan.InitialCurrencyBalance,
			&initialTimestamp,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&plan.CloseOnComplete,
			&plan.UserPlanNumber,
			&plan.Status,
			&plan.CreatedOn,
			&plan.UpdatedOn,
			&order.AccountID,
			&order.KeyPublic,
			&order.KeySecret,
			&order.KeyDescription,
			&order.OrderID,
			&order.ParentOrderID,
			&order.PlanID,
			&order.PlanDepth,
			&order.Exchange,
			&exchangeOrderID,
			&exchangePrice,
			&exchangeTime,
			&order.MarketName,
			&feeSymbol,
			&feeAmount,
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
			&orderErr,
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
		if exchangeOrderID.Valid {
			order.ExchangeOrderID = exchangeOrderID.String
		}
		if initialTimestamp.Valid {
			plan.InitialTimestamp = initialTimestamp.String
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
		if feeSymbol.Valid {
			order.FeeCurrencySymbol = feeSymbol.String
		}
		if feeAmount.Valid {
			order.FeeCurrencyAmount = feeAmount.Float64
		}
		if exchangePrice.Valid {
			order.ExchangePrice = exchangePrice.Float64
		}
		if exchangeTime.Valid {
			order.ExchangeTime = exchangeTime.String
		}
		if orderErr.Valid {
			order.Errors = orderErr.String
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
 			title,
 			total_depth,
			exchange_name,
			user_currency_symbol,
 			active_currency_symbol,
 			active_currency_balance,
 			initial_currency_symbol,
			initial_currency_balance,
			initial_timestamp,
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
		var initialTimestamp sql.NullString

		err := rows.Scan(
			&plan.PlanID,
			&plan.Title,
			&plan.TotalDepth,
			&plan.Exchange,
			&plan.UserCurrencySymbol,
			&plan.CommittedCurrencySymbol,
			&plan.CommittedCurrencyAmount,
			&plan.InitialCurrencySymbol,
			&plan.InitialCurrencyBalance,
			&initialTimestamp,
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
		if initialTimestamp.Valid {
			plan.InitialTimestamp = initialTimestamp.String
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

// This was needed to pull the plans that have orders with an account ID.
// The status of these plans need to be set to closed if active, or inactive
func FindAccountPlans(db *sql.DB, accountID string) ([]*protoPlan.Plan, error) {
	plans := make([]*protoPlan.Plan, 0)
	query := `SELECT
			distinct
			p.id,
			p.title,
			p.total_depth,
			p.exchange_name,
			p.user_currency_symbol,
			p.active_currency_symbol,
			p.active_currency_balance,
			p.initial_currency_symbol,
			p.initial_currency_balance,
			p.initial_timestamp,
			p.last_executed_plan_depth,
			p.last_executed_order_id,
			p.plan_template_id,
			p.close_on_complete,
			p.user_plan_number,
			p.status,
			p.created_on,
			p.updated_on
		FROM plans p 
		JOIN orders o on o.plan_id = p.id
 		WHERE o.account_id = $1`

	rows, err := db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var plan protoPlan.Plan
		var planTemplateID sql.NullString
		var initialTimestamp sql.NullString

		err := rows.Scan(
			&plan.PlanID,
			&plan.Title,
			&plan.TotalDepth,
			&plan.Exchange,
			&plan.UserCurrencySymbol,
			&plan.CommittedCurrencySymbol,
			&plan.CommittedCurrencyAmount,
			&plan.InitialCurrencySymbol,
			&plan.InitialCurrencyBalance,
			&initialTimestamp,
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
		if initialTimestamp.Valid {
			plan.InitialTimestamp = initialTimestamp.String
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

	return plans, nil
}

// func FindUserPlansWithStatus(db *sql.DB, userID, status string, page, pageSize uint32) (*protoPlan.PlansPage, error) {

// 	var count uint32
// 	queryCount := `SELECT count(*) FROM plans WHERE user_id = $1 AND status = $2`
// 	err := db.QueryRow(queryCount, userID, status).Scan(&count)

// 	plans := make([]*protoPlan.Plan, 0)
// 	query := `SELECT
// 			id,
// 			title,
// 			total_depth,
// 			exchange_name,
// 			user_currency_symbol,
// 			active_currency_symbol,
// 			active_currency_balance,
// 			initial_currency_symbol,
// 			initial_currency_balance,
// 			initial_timestamp,
// 			last_executed_plan_depth,
// 			last_executed_order_id,
// 			plan_template_id,
// 			close_on_complete,
// 			user_plan_number,
// 			status,
// 			created_on,
// 			updated_on
// 		FROM plans
// 		WHERE user_id = $1 AND status = $2 ORDER BY created_on OFFSET $3 LIMIT $4`

// 	rows, err := db.Query(query, userID, status, page, pageSize)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var plan protoPlan.Plan
// 		var planTemplateID sql.NullString
// 		var initialTimestamp sql.NullString

// 		err := rows.Scan(
// 			&plan.PlanID,
// 			&plan.Title,
// 			&plan.TotalDepth,
// 			&plan.Exchange,
// 			&plan.UserCurrencySymbol,
// 			&plan.CommittedCurrencySymbol,
// 			&plan.CommittedCurrencyAmount,
// 			&plan.InitialCurrencySymbol,
// 			&plan.InitialCurrencyBalance,
// 			&initialTimestamp,
// 			&plan.LastExecutedPlanDepth,
// 			&plan.LastExecutedOrderID,
// 			&planTemplateID,
// 			&plan.CloseOnComplete,
// 			&plan.UserPlanNumber,
// 			&plan.Status,
// 			&plan.CreatedOn,
// 			&plan.UpdatedOn)

// 		if err != nil {
// 			return nil, err
// 		}
// 		if initialTimestamp.Valid {
// 			plan.InitialTimestamp = initialTimestamp.String
// 		}
// 		if planTemplateID.Valid {
// 			plan.PlanTemplateID = planTemplateID.String
// 		}

// 		plans = append(plans, &plan)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}

// 	result := protoPlan.PlansPage{
// 		Page:     page,
// 		PageSize: pageSize,
// 		Total:    count,
// 		Plans:    plans,
// 	}

// 	return &result, nil
// }

// Find all user plans for exchange with status
//		WHERE user_id = $1 AND status = $2 AND exchange_name = $3 ORDER BY created_on OFFSET $4 LIMIT $5`
func FindUserPlansWithStatus(db *sql.DB, userID, status string, page, pageSize uint32) (*protoPlan.PlansPage, error) {

	var count uint32
	queryCount := `SELECT count(*) FROM plans WHERE user_id = $1 AND status like '%' || $2 || '%'`
	err := db.QueryRow(queryCount, userID, status).Scan(&count)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`SELECT 
	p.id as plan_id,
	p.user_id,
	p.title,
	p.total_depth,
	p.exchange_name,
	p.user_currency_symbol,
	p.user_currency_balance_at_init,
	p.active_currency_symbol,
	p.active_currency_balance,
	p.initial_currency_symbol,
	p.initial_currency_balance,
	p.initial_timestamp,
	p.last_executed_plan_depth,
	p.last_executed_order_id,
	p.close_on_complete,
	p.user_plan_number,
	p.status,
	p.created_on,
	p.updated_on,
	a.id as account_id,
	a.key_public,
	a.key_secret,
	a.description,
	o.id as order_id,
	o.parent_order_id,
	o.plan_id,
	o.plan_depth,
	o.exchange_name,
	o.exchange_order_id,
	o.exchange_price,
	o.exchange_time,
	o.market_name,
	o.fee_currency_symbol,
	o.fee_currency_amount,
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
	o.errors,
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
	JOIN accounts a on o.account_id = a.id
	WHERE p.user_id = $1 AND p.status like '%' || $2 || '%' AND (o.parent_order_id = p.last_executed_order_id or o.plan_depth <= (p.last_executed_plan_depth +1))
	ORDER BY p.id, o.plan_depth, o.order_priority, o.id, t.index
	OFFSET $3 LIMIT $4`, userID, status, page, pageSize)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	plans := make([]*protoPlan.Plan, 0)
	for rows.Next() {
		var plan protoPlan.Plan
		var initialTimestamp sql.NullString
		var trigger protoOrder.Trigger
		var price sql.NullFloat64
		var triggerTemplateID sql.NullString
		var triggeredPrice sql.NullFloat64
		var triggeredCondition sql.NullString
		var triggeredTimestamp sql.NullString
		var orderErr sql.NullString
		var exchangePrice sql.NullFloat64
		var exchangeTime sql.NullString
		var feeSymbol sql.NullString
		var feeAmount sql.NullFloat64
		var finalCurrencySymbol sql.NullString
		var exchangeOrderID sql.NullString
		var actionsStr string
		var order protoOrder.Order
		var balanceAtInit sql.NullFloat64

		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.Title,
			&plan.TotalDepth,
			&plan.Exchange,
			&plan.UserCurrencySymbol,
			&balanceAtInit,
			&plan.CommittedCurrencySymbol,
			&plan.CommittedCurrencyAmount,
			&plan.InitialCurrencySymbol,
			&plan.InitialCurrencyBalance,
			&initialTimestamp,
			&plan.LastExecutedPlanDepth,
			&plan.LastExecutedOrderID,
			&plan.CloseOnComplete,
			&plan.UserPlanNumber,
			&plan.Status,
			&plan.CreatedOn,
			&plan.UpdatedOn,
			&order.AccountID,
			&order.KeyPublic,
			&order.KeySecret,
			&order.KeyDescription,
			&order.OrderID,
			&order.ParentOrderID,
			&order.PlanID,
			&order.PlanDepth,
			&order.Exchange,
			&exchangeOrderID,
			&exchangePrice,
			&exchangeTime,
			&order.MarketName,
			&feeSymbol,
			&feeAmount,
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
			&orderErr,
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
		if exchangeOrderID.Valid {
			order.ExchangeOrderID = exchangeOrderID.String
		}
		if initialTimestamp.Valid {
			plan.InitialTimestamp = initialTimestamp.String
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
		if feeSymbol.Valid {
			order.FeeCurrencySymbol = feeSymbol.String
		}
		if feeAmount.Valid {
			order.FeeCurrencyAmount = feeAmount.Float64
		}
		if exchangePrice.Valid {
			order.ExchangePrice = exchangePrice.Float64
		}
		if exchangeTime.Valid {
			order.ExchangeTime = exchangeTime.String
		}
		if orderErr.Valid {
			order.Errors = orderErr.String
		}
		if balanceAtInit.Valid {
			plan.UserCurrencyBalanceAtInit = balanceAtInit.Float64
		}
		trigger.OrderID = order.OrderID

		// does the plan structure already exist in the plans
		isNewPlan := true
		for _, p := range plans {
			// append structures to existing plan
			if p.PlanID == plan.PlanID {
				// is this a new order?
				isNewOrder := true
				for _, o := range p.Orders {
					// if not new order append trigger
					if o.OrderID == order.OrderID {
						isNewTrigger := true
						for _, t := range o.Triggers {
							if t.TriggerID == trigger.TriggerID {
								isNewTrigger = false
								break
							}
						}

						if isNewTrigger {
							o.Triggers = append(o.Triggers, &trigger)
						}

						isNewOrder = false
						break
					}
				}

				if isNewOrder {
					// append the trigger to the new order
					order.Triggers = append(order.Triggers, &trigger)
					p.Orders = append(p.Orders, &order)
				}

				isNewPlan = false
				break
			}
		}

		// new plan append to plans
		if isNewPlan {
			order.Triggers = append(order.Triggers, &trigger)
			plan.Orders = append(plan.Orders, &order)
			// TODO merge all under the appropriate structure
			plans = append(plans, &plan)
		}
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

func FindPlanUserCurrencySymbol(db *sql.DB, planID string) (string, error) {

	var userCurrencySymbol sql.NullString
	query := `SELECT user_currency_symbol FROM plans WHERE id = $1`
	err := db.QueryRow(query, planID).Scan(&userCurrencySymbol)

	if err != nil {
		return "", err
	}

	var symbol string
	if userCurrencySymbol.Valid {
		symbol = userCurrencySymbol.String
	}

	return symbol, nil
}

// Find all user plans for exchange with status
func FindUserMarketPlansWithStatus(db *sql.DB, userID, status, exchange, marketName string, page, pageSize uint32) (*protoPlan.PlansPage, error) {

	var count uint32
	queryCount := `SELECT count(*) FROM plans WHERE user_id = $1 AND status = $2 AND exchange_name = $3 AND market_name = $4`
	err := db.QueryRow(queryCount, userID, status, exchange, marketName).Scan(&count)

	plans := make([]*protoPlan.Plan, 0)
	query := `SELECT 
			id,
			title,
			total_depth,
			exchange_name,
			user_currency_symbol,
			active_currency_symbol,
			active_currency_balance,
			initial_currency_symbol,
			initial_currency_balance,
			initial_timestamp,
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
		var initialTimestamp sql.NullString

		err := rows.Scan(
			&plan.PlanID,
			&plan.Title,
			&plan.TotalDepth,
			&plan.Exchange,
			&plan.UserCurrencySymbol,
			&plan.CommittedCurrencySymbol,
			&plan.CommittedCurrencyAmount,
			&plan.InitialCurrencySymbol,
			&plan.InitialCurrencyBalance,
			&initialTimestamp,
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
		if initialTimestamp.Valid {
			plan.InitialTimestamp = initialTimestamp.String
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
		title,
		total_depth,
		exchange_name,
		user_currency_symbol,
		active_currency_symbol,
		active_currency_balance,
		initial_currency_symbol,
		initial_currency_balance,
		last_executed_plan_depth,
		last_executed_order_id,
		plan_template_id,
		close_on_complete,
		user_plan_number,
		status, 
		created_on,
		updated_on) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`)

	if err != nil {
		txn.Rollback()
		return errors.New("prepare insert plan failed: " + err.Error())
	}
	_, err = txn.Stmt(sqlStatement).Exec(
		newPlan.PlanID,
		newPlan.UserID,
		newPlan.Title,
		newPlan.TotalDepth,
		newPlan.Exchange,
		newPlan.UserCurrencySymbol,
		newPlan.CommittedCurrencySymbol,
		newPlan.CommittedCurrencyAmount,
		newPlan.InitialCurrencySymbol,
		newPlan.InitialCurrencyBalance,
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

func UpdatePlanInitTimestamp(db *sql.DB, planID, timestamp string, balance float64) error {
	stmt := `
		UPDATE plans 
		SET 
			initial_timestamp = $1,
			user_currency_balance_at_init = $2
		WHERE
			id = $3`

	_, err := db.Exec(stmt, timestamp, balance, planID)

	return err
}

// returns (initial_time, error)
func UpdatePlanContext(db *sql.DB, planID, executedOrderID, exchange, symbol string, activeBalance, referencePrice float64, planDepth int32) (string, error) {
	stmt := `
		UPDATE plans 
		SET 
			last_executed_order_id = $1,
			last_executed_plan_depth = $2,
			active_currency_symbol = $3,
			active_currency_balance = $4,
			exchange_name = $5, 
			reference_price = $6
		WHERE
			id = $7
		RETURNING initial_timestamp`

	var initTime sql.NullString
	err := db.QueryRow(stmt, executedOrderID, planDepth, symbol, activeBalance, exchange, referencePrice, planID).Scan(&initTime)
	if err != nil {
		return "", err
	}

	return initTime.String, nil
}

func UpdatePlanContextTxn(txn *sql.Tx, ctx context.Context, planID, symbol, exchange string, activeBalance float64) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE plans 
		SET 
			active_currency_symbol = $1,
			active_currency_balance = $2,
			initial_currency_symbol = $3,
			initial_currency_balance = $4,
			exchange_name = $5
		WHERE
			id = $6`,
		symbol, activeBalance, symbol, activeBalance, exchange, planID)

	return err
}

func UpdatePlanBaseCurrencyTxn(txn *sql.Tx, ctx context.Context, planID, baseCurrencySym string) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE plans 
		SET 
			user_currency_symbol = $1 
		WHERE
			id = $2`,
		baseCurrencySym, planID)

	return err
}

func UpdatePlanInitTimeTxn(txn *sql.Tx, ctx context.Context, planID, timestamp string) error {
	stmt := `
		UPDATE plans 
		SET 
			initial_timestamp = $1
		WHERE
			id = $2`

	_, err := txn.ExecContext(ctx, stmt, timestamp, planID)

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

func UpdatePlanTotalDepthTxn(txn *sql.Tx, ctx context.Context, planID string, totalDepth uint32) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE plans 
		SET 
			total_depth = $1 
		WHERE
			id = $2`,
		totalDepth, planID)

	return err
}

func UpdatePlanTitleTxn(txn *sql.Tx, ctx context.Context, planID, title string) error {
	_, err := txn.ExecContext(ctx, `
		UPDATE plans 
		SET 
			title = $1 
		WHERE
			id = $2`,
		title, planID)

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
			active_currency_symbol = $2,
			active_currency_balance = $3,
			plan_template_id = $4,
			close_on_complete = $5,
			status = $6 
		WHERE
			id = $7`,
		plan.Exchange,
		plan.CommittedCurrencySymbol,
		plan.CommittedCurrencyAmount,
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
