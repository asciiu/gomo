package sql

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/asciiu/gomo/common/constants/key"
	"github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/status"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// DeletePlan is a hard delete
func DeletePlan(db *sql.DB, planID string) error {
	_, err := db.Exec("DELETE FROM plans WHERE id = $1", planID)
	return err
}

func FindPlanSummary(db *sql.DB, planID string) (*protoPlan.Plan, error) {
	var plan protoPlan.Plan

	err := db.QueryRow(`SELECT 
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
		FROM plans p
		JOIN user_keys k on p.user_key_id = k.id
		WHERE p.id = $1`, planID).Scan(
		&plan.PlanID,
		&plan.PlanTemplateID,
		&plan.UserID,
		&plan.KeyID,
		&plan.KeyDescription,
		&plan.ActiveOrderNumber,
		&plan.Exchange,
		&plan.MarketName,
		&plan.BaseBalance,
		&plan.CurrencyBalance,
		&plan.Status)

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindPlanWithPagedOrders ...
func FindPlanWithPagedOrders(db *sql.DB, req *protoPlan.GetUserPlanRequest) (*protoPlan.PlanWithPagedOrders, error) {

	var planPaged protoPlan.PlanWithPagedOrders
	var page protoPlan.OrdersPage

	var planTemp sql.NullString
	var orderTemp sql.NullString
	var price sql.NullFloat64
	var basePercent sql.NullFloat64
	var currencyPercent sql.NullFloat64
	var nextOrderID sql.NullString

	var count uint32
	queryCount := `SELECT count(*) FROM plan_orders WHERE plan_id = $1`
	err := db.QueryRow(queryCount, req.PlanID).Scan(&count)

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
		po.id,
		po.base_percent,
		po.currency_percent,
		po.side,
		po.order_number,
		po.order_type,
		po.order_template_id,
		po.conditions,
		po.price,
		po.next_plan_order_id,
		po.status
		FROM plans p
		JOIN plan_orders po on p.id = po.plan_id
		JOIN user_keys k on p.user_key_id = k.id
		WHERE p.id = $1 ORDER BY po.order_number OFFSET $2 LIMIT $3`, req.PlanID, req.Page, req.PageSize)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var order protoPlan.Order

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
			&order.OrderID,
			&basePercent,
			&currencyPercent,
			&order.Side,
			&order.OrderNumber,
			&order.OrderType,
			&orderTemp,
			&order.Conditions,
			&price,
			&nextOrderID,
			&order.Status)

		if err != nil {
			return nil, err
		}
		if basePercent.Valid {
			order.BasePercent = basePercent.Float64
		}
		if currencyPercent.Valid {
			order.CurrencyPercent = currencyPercent.Float64
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
		if err := json.Unmarshal([]byte(order.Conditions), &order.Conditions); err != nil {
			return nil, err
		}
		page.Orders = append(page.Orders, &order)
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
	results := make([]*protoPlan.Plan, 0)

	rows, err := db.Query(`SELECT p.id,
		p.user_id,
		p.user_key_id,
		k.api_key,
		k.secret,
		p.exchange_name,
		p.market_name,
		p.plan_order_ids,
		p.base_balance,
		p.currency_balance,
		p.status,
		po.id,
		po.base_percent,
		po.currency_percent,
		po.side,
		po.order_type,
		po.conditions,
		po.price, 
		po.next_plan_order_id,
		po.status
		FROM plans p 
		JOIN plan_orders po on p.id = po.plan_id
		JOIN user_keys k on p.user_key_id = k.id
		WHERE p.status = $1 AND po.status = $2 AND k.status = $3`, plan.Active, status.Active, key.Verified)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var plan protoPlan.Plan
		var order protoPlan.Order
		var planOrderIds pq.StringArray

		var price sql.NullFloat64
		var basePercent sql.NullFloat64
		var currencyPercent sql.NullFloat64
		var nextOrderID sql.NullString
		err := rows.Scan(
			&plan.PlanID,
			&plan.UserID,
			&plan.KeyID,
			&plan.Key,
			&plan.KeySecret,
			&plan.Exchange,
			&plan.MarketName,
			&planOrderIds,
			&plan.BaseBalance,
			&plan.CurrencyBalance,
			&plan.Status,
			&order.OrderID,
			&basePercent,
			&currencyPercent,
			&order.Side,
			&order.OrderType,
			&order.Conditions,
			&price,
			&nextOrderID,
			&order.Status)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		if basePercent.Valid {
			order.BasePercent = basePercent.Float64
		}
		if currencyPercent.Valid {
			order.CurrencyPercent = currencyPercent.Float64
		}
		if price.Valid {
			order.LimitPrice = price.Float64
		}
		if nextOrderID.Valid {
			order.NextOrderID = nextOrderID.String
		}
		if err := json.Unmarshal([]byte(order.Conditions), &order.Conditions); err != nil {
			return nil, err
		}
		plan.Orders = append(plan.Orders, &order)
		results = append(results, &plan)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return results, nil
}

// Returns plan with specific order that has a verified key.
func FindPlanWithOrderID(db *sql.DB, orderID string) (*protoPlan.Plan, error) {
	var plan protoPlan.Plan
	var order protoPlan.Order

	var price sql.NullFloat64
	var basePercent sql.NullFloat64
	var currencyPercent sql.NullFloat64
	var nextOrderID sql.NullString

	err := db.QueryRow(`SELECT 
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
		po.id,
		po.base_percent,
		po.currency_percent,
		po.side,
		po.order_type,
		po.conditions,
		po.price, 
		po.next_plan_order_id,
		po.status
		FROM plans p 
		JOIN plan_orders po on p.id = po.plan_id
		JOIN user_keys k on p.user_key_id = k.id
		WHERE po.id = $1 AND k.status = $2`, orderID, key.Verified).Scan(
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
		&basePercent,
		&currencyPercent,
		&order.Side,
		&order.OrderType,
		&order.Conditions,
		&price,
		&nextOrderID,
		&order.Status)

	if err != nil {
		return nil, err
	}

	if basePercent.Valid {
		order.BasePercent = basePercent.Float64
	}
	if currencyPercent.Valid {
		order.CurrencyPercent = currencyPercent.Float64
	}
	if price.Valid {
		order.LimitPrice = price.Float64
	}
	if nextOrderID.Valid {
		order.NextOrderID = nextOrderID.String
	}
	if err := json.Unmarshal([]byte(order.Conditions), &order.Conditions); err != nil {
		return nil, err
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
		status
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
			&plan.Status)

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
		status
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
			&plan.Status)

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
		status
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
			&plan.Status)

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
	orders := make([]*protoPlan.Order, 0)

	// generate order ids
	for range req.Orders {
		orderIDs = append(orderIDs, uuid.New())
	}

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
		status) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

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
		req.Status)

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
		jsonCond, err := json.Marshal(or.Conditions)
		if err != nil {
			return nil, err
		}

		orderID := orderIDs[i]
		var nextOrderID uuid.UUID
		var nextIndex = i + 1
		if nextIndex < len(orderIDs) {
			nextOrderID = orderIDs[nextIndex]
		}
		sqlInsertOrder := `insert into plan_orders (
			id, 
			plan_id, 
			base_percent, 
			currency_percent, 
			side,
			order_template_id,
			order_number,
			order_type,
			conditions,
			price,
			next_plan_order_id,
			status) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err = db.Exec(sqlInsertOrder,
			orderID,
			planID,
			or.BasePercent,
			or.CurrencyPercent,
			or.Side,
			or.OrderTemplateID,
			i,
			or.OrderType,
			jsonCond,
			or.LimitPrice,
			nextOrderID,
			orderStatus)

		if err != nil {
			return nil, err
		}
		order := protoPlan.Order{
			OrderID:         orderID.String(),
			OrderNumber:     uint32(i),
			OrderType:       or.OrderType,
			OrderTemplateID: or.OrderTemplateID,
			Side:            or.Side,
			LimitPrice:      or.LimitPrice,
			BasePercent:     or.BasePercent,
			CurrencyPercent: or.CurrencyPercent,
			Status:          orderStatus,
			Conditions:      or.Conditions,
			NextOrderID:     nextOrderID.String(),
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
	}
	return &plan, nil
}

// UpdateOrder should be an update head chain buy order
// todo add sql for update other types of orders - UpdateHeadBuyOrder, UpdateHeadSellOrder, UpdateChainBuyOrder, UpdateChainSellOrder
// TODO needs work!!
// func UpdateOrder(db *sql.DB, req *orderProto.OrderRequest) (*orderProto.Order, error) {
// 	// sqlStatement := `UPDATE orders SET
// 	// conditions = $1,
// 	// base_quantity = $2,
// 	// WHERE id = $4 and user_id = $5
// 	// RETURNING id, user_id, exchange_name, market_name, user_key, side,
// 	// base_quantity, base_quantity_remainder, status, conditions, parent_order_id`

// 	// jsonCond, err := json.Marshal(req.Conditions)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	var o orderProto.Order
// 	// err = db.QueryRow(sqlStatement, jsonCond, req.BaseQuantity, req.BaseQuantity, req.OrderID, req.UserID).
// 	// 	Scan(&o.OrderID, &o.UserID, &o.Exchange, &o.MarketName, &o.KeyID,
// 	// 		&o.Side, &o.BaseQuantity, &o.Status, &o.Conditions, &o.ParentOrderID)

// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
// 	// 	return nil, err
// 	// }

// 	return &o, nil
// }

// Maybe this should return more of the updated order but all I need from this
// as of now is the next_plan_order_id so I can load the next order.
func UpdatePlanOrder(db *sql.DB, orderID, stat string) (*protoPlan.Order, error) {
	updatePlanOrderSql := `UPDATE plan_orders SET status = $1 WHERE id = $2 
	RETURNING next_plan_order_id, order_number, plan_id`

	var o protoPlan.Order
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
	RETURNING id`

	var p protoPlan.Plan
	err := db.QueryRow(sqlStatement, status, planID).
		Scan(&p.PlanID)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
