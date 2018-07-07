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

// FindOrderByID ...
// func FindPlanByID(db *sql.DB, req *protoPlan.GetUserPlanRequest) (*protoPlan.Plan, error) {
// 	var p protoPlan.Plan
// 	var condition sql.NullString
// 	var exchangeOrderID sql.NullString
// 	var exchangeMarketName sql.NullString
// 	var price sql.NullFloat64
// 	var ids []string
// 	err := db.QueryRow(`SELECT id,
// 		user_id,
// 		user_key_id,
// 		exchange_name,
// 		exchange_order_id,
// 		exchange_market_name,
// 		market_name,
// 		order_ids,
// 		base_balance,
// 		currency_balance,
// 		status
// 		FROM plans WHERE id = $1`, req.PlanID).
// 		Scan(&p.PlanID,
// 			&p.UserID,
// 			&p.KeyID,
// 			&p.Exchange,
// 			&exchangeOrderID,
// 			&exchangeMarketName,
// 			&p.MarketName,
// 			&ids,
// 			&p.BaseBalance,
// 			&p.CurrencyBalance,
// 			&p.Status)

// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// if exchangeOrderID.Valid {
// 	// 	p.ExchangeOrderID = exchangeOrderID.String
// 	// }
// 	// if exchangeMarketName.Valid {
// 	// 	p.ExchangeMarketName = exchangeMarketName.String
// 	// }
// 	// if price.Valid {
// 	// 	o.Price = price.Float64
// 	// }

// 	// if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
// 	// 	return nil, err
// 	// }

// 	return &p, nil
// }

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
			&plan.Secret,
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
			order.Price = price.Float64
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
		&plan.Secret,
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
		order.Price = price.Float64
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

// func FindOrdersByUserID(db *sql.DB, req *orderProto.GetUserOrdersRequest) ([]*orderProto.Order, error) {
// 	results := make([]*orderProto.Order, 0)

// 	rows, err := db.Query(`SELECT id,
// 		user_id,
// 		user_key_id,
// 		exchange_name,
// 		exchange_order_id,
// 		exchange_market_name,
// 		market_name,
// 		side,
// 		type,
// 		price,
// 		base_quantity,
// 		base_percent,
// 		currency_quantity,
// 		currency_percent,
// 		status,
// 		chain_status,
// 		conditions,
// 		condition,
// 		parent_order_id
// 		FROM orders WHERE user_id = $1`, req.UserID)

// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var o orderProto.Order
// 		var condition sql.NullString
// 		var exchangeOrderID sql.NullString
// 		var exchangeMarketName sql.NullString
// 		var price sql.NullFloat64
// 		err := rows.Scan(&o.OrderID,
// 			&o.UserID,
// 			&o.KeyID,
// 			&o.Exchange,
// 			&exchangeOrderID,
// 			&exchangeMarketName,
// 			&o.MarketName,
// 			&o.Side,
// 			&o.OrderType,
// 			&price,
// 			&o.BaseQuantity,
// 			&o.BasePercent,
// 			&o.CurrencyQuantity,
// 			&o.CurrencyPercent,
// 			&o.Status,
// 			&o.ChainStatus,
// 			&o.Conditions,
// 			&condition,
// 			&o.ParentOrderID)

// 		if err != nil {
// 			log.Fatal(err)
// 			return nil, err
// 		}
// 		if condition.Valid {
// 			o.Condition = condition.String
// 		}
// 		if exchangeOrderID.Valid {
// 			o.ExchangeOrderID = exchangeOrderID.String
// 		}
// 		if exchangeMarketName.Valid {
// 			o.ExchangeMarketName = exchangeMarketName.String
// 		}
// 		if price.Valid {
// 			o.Price = price.Float64
// 		}

// 		if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
// 			return nil, err
// 		}
// 		results = append(results, &o)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}

// 	return results, nil
// }

// TODO this should be inserted as a transaction
func InsertPlan(db *sql.DB, req *protoPlan.PlanRequest) (*protoPlan.Plan, error) {
	planID := uuid.New()
	orderIDs := make([]uuid.UUID, 0)
	orders := make([]*protoPlan.Order, 0)

	planStatus := plan.Active
	if !req.Active {
		planStatus = plan.Inactive
	}

	// generate order ids
	for range req.Orders {
		orderIDs = append(orderIDs, uuid.New())
	}

	sqlStatement := `insert into plans (id, 
		user_id, 
		user_key_id, 
		exchange_name, 
		market_name, 
		base_balance, 
		currency_balance, 
		plan_order_ids,
		status) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := db.Exec(sqlStatement,
		planID,
		req.UserID,
		req.KeyID,
		req.Exchange,
		req.MarketName,
		req.BaseBalance,
		req.CurrencyBalance,
		pq.Array(orderIDs),
		planStatus)

	if err != nil {
		return nil, err
	}

	// loop through and insert orders
	for i, or := range req.Orders {
		orderStatus := status.Pending
		// first order of active plan is active
		if i == 0 && planStatus == plan.Active {
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
			order_type,
			conditions,
			price,
			next_plan_order_id,
			status) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		_, err = db.Exec(sqlInsertOrder,
			orderID,
			planID,
			or.BasePercent,
			or.CurrencyPercent,
			or.Side,
			or.OrderType,
			jsonCond,
			or.Price,
			nextOrderID,
			orderStatus)

		if err != nil {
			return nil, err
		}
		order := protoPlan.Order{
			OrderID:         orderID.String(),
			OrderType:       or.OrderType,
			Side:            or.Side,
			Price:           or.Price,
			BasePercent:     or.BasePercent,
			CurrencyPercent: or.CurrencyPercent,
			Status:          orderStatus,
			Conditions:      or.Conditions,
			NextOrderID:     nextOrderID.String(),
		}
		orders = append(orders, &order)
	}

	plan := protoPlan.Plan{
		PlanID:          planID.String(),
		UserID:          req.UserID,
		KeyID:           req.KeyID,
		Exchange:        req.Exchange,
		MarketName:      req.MarketName,
		BaseBalance:     req.BaseBalance,
		CurrencyBalance: req.CurrencyBalance,
		Orders:          orders,
		Status:          planStatus,
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
func UpdateOrderStatus(db *sql.DB, orderID, status string) (*protoPlan.Order, error) {
	sqlStatement := `UPDATE plan_orders SET status = $1 WHERE id = $2 
	RETURNING next_plan_order_id`

	var o protoPlan.Order
	err := db.QueryRow(sqlStatement, status, orderID).
		Scan(&o.NextOrderID)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func UpdatePlanStatus(db *sql.DB, planID, status string) (*protoPlan.Plan, error) {
	sqlStatement := `UPDATE plan SET status = $1 WHERE id = $2 
	RETURNING id`

	var p protoPlan.Plan
	err := db.QueryRow(sqlStatement, status, planID).
		Scan(&p.PlanID)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
