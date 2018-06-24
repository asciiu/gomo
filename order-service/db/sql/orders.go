package sql

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/asciiu/gomo/common/constants/key"
	"github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
	orderProto "github.com/asciiu/gomo/order-service/proto/order"
	"github.com/google/uuid"
)

// DeleteOrder will completely remove the order from the DB
func DeleteOrder(db *sql.DB, orderID string) error {
	_, err := db.Exec("DELETE FROM orders WHERE id = $1", orderID)
	return err
}

// FindOrderByID ...
func FindOrderByID(db *sql.DB, req *orderProto.GetUserOrderRequest) (*orderProto.Order, error) {
	var o orderProto.Order
	var condition sql.NullString
	var exchangeOrderID sql.NullString
	var exchangeMarketName sql.NullString
	var price sql.NullFloat64
	err := db.QueryRow(`SELECT id, 
		user_id, 
		user_key_id, 
		exchange_name, 
		exchange_order_id, 
		exchange_market_name,
		market_name, 
		side, 
		type, 
		price,
		base_quantity, 
		base_percent, 
		currency_quantity, 
		currency_percent, 
		status, 
		chain_status, 
		conditions, 
		condition, 
		parent_order_id 
		FROM orders WHERE id = $1`, req.OrderID).
		Scan(&o.OrderID,
			&o.UserID,
			&o.KeyID,
			&o.Exchange,
			&exchangeOrderID,
			&exchangeMarketName,
			&o.MarketName,
			&o.Side,
			&o.OrderType,
			&price,
			&o.BaseQuantity,
			&o.BasePercent,
			&o.CurrencyQuantity,
			&o.CurrencyPercent,
			&o.Status,
			&o.ChainStatus,
			&o.Conditions,
			&condition,
			&o.ParentOrderID)

	if err != nil {
		return nil, err
	}
	if condition.Valid {
		o.Condition = condition.String
	}
	if exchangeOrderID.Valid {
		o.ExchangeOrderID = exchangeOrderID.String
	}
	if exchangeMarketName.Valid {
		o.ExchangeMarketName = exchangeMarketName.String
	}
	if price.Valid {
		o.Price = price.Float64
	}

	if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
		return nil, err
	}

	return &o, nil
}

// FindOrderWithParentID will be invoked by the order listeners after an order has been filled
// and when the filled order has a child order. Chain must be active.
func FindOrderWithParentID(db *sql.DB, parentOrderID string) (*orderProto.Order, error) {
	var o orderProto.Order
	var price sql.NullFloat64

	err := db.QueryRow(`SELECT id, 
		user_id, 
		user_key_id,
		market_name, 
		side, 
		type, 
		price,
		base_quantity, 
		base_percent, 
		currency_quantity, 
		currency_percent, 
		status, 
		conditions, 
		parent_order_id 
		FROM orders 
		WHERE parent_order_id = $1 and chain_status = $2`, parentOrderID, plan.Active).
		Scan(&o.OrderID,
			&o.UserID,
			&o.KeyID,
			&o.MarketName,
			&o.Side,
			&o.OrderType,
			&price,
			&o.BaseQuantity,
			&o.BasePercent,
			&o.CurrencyQuantity,
			&o.CurrencyPercent,
			&o.Status,
			&o.Conditions,
			&o.ParentOrderID)

	if err != nil {
		return nil, err
	}

	if price.Valid {
		o.Price = price.Float64
	}

	if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
		return nil, err
	}

	return &o, nil
}

// Find all active orders in the DB. This wil load the keys for each order.
func FindActiveOrders(db *sql.DB) ([]*orderProto.Order, error) {
	results := make([]*orderProto.Order, 0)

	// chain must be active
	rows, err := db.Query(`SELECT o.id, 
		o.user_id, 
		o.user_key_id, 
		o.exchange_name,
		o.exchange_order_id, 
		o.exchange_market_name,
		o.market_name, 
		o.side, 
		o.type, 
		o.price,
		o.base_quantity, 
		o.base_percent, 
		o.currency_quantity, 
		o.currency_percent, 
		o.status, 
		o.conditions, 
		o.condition, 
		o.parent_order_id, 
		u.api_key, 
		u.secret 
		FROM orders o 
		JOIN user_keys u on u.id = o.user_key_id 
		WHERE o.status = $1 AND u.status = $2 AND o.chain_status = $3`, status.Active, key.Verified, plan.Active)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o orderProto.Order
		var condition sql.NullString
		var exchangeOrderID sql.NullString
		var exchangeMarketName sql.NullString
		var price sql.NullFloat64
		err := rows.Scan(&o.OrderID,
			&o.UserID,
			&o.KeyID,
			&o.Exchange,
			&exchangeOrderID,
			&exchangeMarketName,
			&o.MarketName,
			&o.Side,
			&o.OrderType,
			&price,
			&o.BaseQuantity,
			&o.BasePercent,
			&o.CurrencyQuantity,
			&o.CurrencyPercent,
			&o.Status,
			&o.Conditions,
			&condition,
			&o.ParentOrderID,
			&o.Key,
			&o.Secret)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		if condition.Valid {
			o.Condition = condition.String
		}
		if exchangeOrderID.Valid {
			o.ExchangeOrderID = exchangeOrderID.String
		}
		if exchangeMarketName.Valid {
			o.ExchangeMarketName = exchangeMarketName.String
		}
		if price.Valid {
			o.Price = price.Float64
		}

		if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return results, nil
}

func FindOrdersByUserID(db *sql.DB, req *orderProto.GetUserOrdersRequest) ([]*orderProto.Order, error) {
	results := make([]*orderProto.Order, 0)

	rows, err := db.Query(`SELECT id, 
		user_id, 
		user_key_id, 
		exchange_name, 
		exchange_order_id, 
		exchange_market_name,
		market_name, 
		side, 
		type, 
		price,
		base_quantity, 
		base_percent, 
		currency_quantity, 
		currency_percent, 
		status, 
		chain_status,
		conditions, 
		condition, 
		parent_order_id 
		FROM orders WHERE user_id = $1`, req.UserID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o orderProto.Order
		var condition sql.NullString
		var exchangeOrderID sql.NullString
		var exchangeMarketName sql.NullString
		var price sql.NullFloat64
		err := rows.Scan(&o.OrderID,
			&o.UserID,
			&o.KeyID,
			&o.Exchange,
			&exchangeOrderID,
			&exchangeMarketName,
			&o.MarketName,
			&o.Side,
			&o.OrderType,
			&price,
			&o.BaseQuantity,
			&o.BasePercent,
			&o.CurrencyQuantity,
			&o.CurrencyPercent,
			&o.Status,
			&o.ChainStatus,
			&o.Conditions,
			&condition,
			&o.ParentOrderID)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		if condition.Valid {
			o.Condition = condition.String
		}
		if exchangeOrderID.Valid {
			o.ExchangeOrderID = exchangeOrderID.String
		}
		if exchangeMarketName.Valid {
			o.ExchangeMarketName = exchangeMarketName.String
		}
		if price.Valid {
			o.Price = price.Float64
		}

		if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return results, nil
}

func InsertOrder(db *sql.DB, req *orderProto.OrderRequest, status string) (*orderProto.Order, error) {
	newID := uuid.New()

	jsonCond, err := json.Marshal(req.Conditions)
	if err != nil {
		return nil, err
	}
	chainStatus := plan.Active
	if !req.Active {
		chainStatus = plan.Inactive
	}

	// the exchange_order_id and exchange_market_name must be "" and not null
	// when scanning in order data null cannot be set on a type string. Therefore,
	// just default those cols to "".
	sqlStatement := `insert into orders (id, 
		user_id, 
		user_key_id, 
		exchange_name, 
		market_name, 
		side, 
		type, 
		price,
		base_quantity, 
		base_percent, 
		currency_quantity, 
		currency_percent, 
		status, 
		chain_status,
		conditions, 
		parent_order_id) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	_, err = db.Exec(sqlStatement,
		newID,
		req.UserID,
		req.KeyID,
		req.Exchange,
		req.MarketName,
		req.Side,
		req.OrderType,
		req.Price,
		req.BaseQuantity,
		req.BasePercent,
		req.CurrencyQuantity,
		req.CurrencyPercent,
		status,
		chainStatus,
		jsonCond,
		req.ParentOrderID)

	if err != nil {
		return nil, err
	}
	order := &orderProto.Order{
		OrderID:          newID.String(),
		UserID:           req.UserID,
		KeyID:            req.KeyID,
		Exchange:         req.Exchange,
		MarketName:       req.MarketName,
		Side:             req.Side,
		OrderType:        req.OrderType,
		Price:            req.Price,
		BaseQuantity:     req.BaseQuantity,
		BasePercent:      req.BasePercent,
		CurrencyQuantity: req.CurrencyQuantity,
		CurrencyPercent:  req.CurrencyPercent,
		Status:           status,
		ChainStatus:      chainStatus,
		Conditions:       req.Conditions,
		ParentOrderID:    req.ParentOrderID,
	}
	return order, nil
}

// UpdateOrder should be an update head chain buy order
// todo add sql for update other types of orders - UpdateHeadBuyOrder, UpdateHeadSellOrder, UpdateChainBuyOrder, UpdateChainSellOrder
// TODO needs work!!
func UpdateOrder(db *sql.DB, req *orderProto.OrderRequest) (*orderProto.Order, error) {
	// sqlStatement := `UPDATE orders SET
	// conditions = $1,
	// base_quantity = $2,
	// WHERE id = $4 and user_id = $5
	// RETURNING id, user_id, exchange_name, market_name, user_key, side,
	// base_quantity, base_quantity_remainder, status, conditions, parent_order_id`

	// jsonCond, err := json.Marshal(req.Conditions)
	// if err != nil {
	// 	return nil, err
	// }

	var o orderProto.Order
	// err = db.QueryRow(sqlStatement, jsonCond, req.BaseQuantity, req.BaseQuantity, req.OrderID, req.UserID).
	// 	Scan(&o.OrderID, &o.UserID, &o.Exchange, &o.MarketName, &o.KeyID,
	// 		&o.Side, &o.BaseQuantity, &o.Status, &o.Conditions, &o.ParentOrderID)

	// if err != nil {
	// 	return nil, err
	// }

	// if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
	// 	return nil, err
	// }

	return &o, nil
}

func UpdateOrderStatus(db *sql.DB, evt *evt.OrderEvent) (*orderProto.Order, error) {
	sqlStatement := `UPDATE orders SET status = $1, condition = $2, exchange_order_id = $3, exchange_market_name = $4 
	WHERE id = $5 RETURNING id, user_id, exchange_name, market_name, user_key_id, side, base_quantity, 
	base_percent, currency_quantity, currency_percent, status, conditions`

	var o orderProto.Order
	err := db.QueryRow(sqlStatement, evt.Status, evt.Condition, evt.ExchangeOrderID, evt.ExchangeMarketName, evt.OrderID).
		Scan(&o.OrderID, &o.UserID, &o.Exchange, &o.MarketName, &o.KeyID, &o.Side, &o.BaseQuantity,
			&o.BasePercent, &o.CurrencyQuantity, &o.CurrencyPercent, &o.Status, &o.Conditions)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
		return nil, err
	}

	return &o, nil
}
