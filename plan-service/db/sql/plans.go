package sql

import (
	"database/sql"
	"encoding/json"

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
// func FindActiveOrders(db *sql.DB) ([]*orderProto.Order, error) {
// 	results := make([]*orderProto.Order, 0)

// 	// chain must be active
// 	rows, err := db.Query(`SELECT o.id,
// 		o.user_id,
// 		o.user_key_id,
// 		o.exchange_name,
// 		o.exchange_order_id,
// 		o.exchange_market_name,
// 		o.market_name,
// 		o.side,
// 		o.type,
// 		o.price,
// 		o.base_quantity,
// 		o.base_percent,
// 		o.currency_quantity,
// 		o.currency_percent,
// 		o.status,
// 		o.conditions,
// 		o.condition,
// 		o.parent_order_id,
// 		u.api_key,
// 		u.secret
// 		FROM orders o
// 		JOIN user_keys u on u.id = o.user_key_id
// 		WHERE o.status = $1 AND u.status = $2 AND o.chain_status = $3`, status.Active, key.Verified, chain.Active)

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
// 			&o.Conditions,
// 			&condition,
// 			&o.ParentOrderID,
// 			&o.Key,
// 			&o.Secret)

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
		sqlInsertOrder := `insert into plan_orders (
			id, 
			plan_id, 
			base_percent, 
			currency_percent, 
			side,
			conditions,
			price,
			status) 
			values ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err = db.Exec(sqlInsertOrder,
			orderID,
			planID,
			or.BasePercent,
			or.CurrencyPercent,
			or.Side,
			jsonCond,
			or.Price,
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

// func UpdateOrderStatus(db *sql.DB, evt *evt.OrderEvent) (*orderProto.Order, error) {
// 	sqlStatement := `UPDATE orders SET status = $1, condition = $2, exchange_order_id = $3, exchange_market_name = $4
// 	WHERE id = $5 RETURNING id, user_id, exchange_name, market_name, user_key_id, side, base_quantity,
// 	base_percent, currency_quantity, currency_percent, status, conditions`

// 	var o orderProto.Order
// 	err := db.QueryRow(sqlStatement, evt.Status, evt.Condition, evt.ExchangeOrderID, evt.ExchangeMarketName, evt.OrderID).
// 		Scan(&o.OrderID, &o.UserID, &o.Exchange, &o.MarketName, &o.KeyID, &o.Side, &o.BaseQuantity,
// 			&o.BasePercent, &o.CurrencyQuantity, &o.CurrencyPercent, &o.Status, &o.Conditions)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
// 		return nil, err
// 	}

// 	return &o, nil
// }
