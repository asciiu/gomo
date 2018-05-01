package sql

import (
	"database/sql"
	"encoding/json"
	"log"

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
	err := db.QueryRow(`SELECT id, user_id, user_key_id, exchange_name, exchange_order_id, exchange_market_name,
		 market_name, side, type, base_quantity, base_percent, currency_quantity, currency_percent, status, 
		 conditions, condition, parent_order_id FROM orders WHERE id = $1`, req.OrderID).
		Scan(&o.OrderID, &o.UserID, &o.KeyID, &o.Exchange, &o.ExchangeOrderID, &o.ExchangeMarketName,
			&o.MarketName, &o.Side, &o.OrderType, &o.BaseQuantity, &o.BasePercent, &o.CurrencyQuantity,
			&o.CurrencyPercent, &o.Status, &o.Conditions, &o.Condition, &o.ParentOrderID)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
		return nil, err
	}

	return &o, nil
}

// FindOrderWithParentID will be invoked by the order listeners after an order has been filled
// and when the filled order has a child order
func FindOrderWithParentID(db *sql.DB, parentOrderID string) (*orderProto.Order, error) {
	var o orderProto.Order
	err := db.QueryRow(`SELECT id, user_id, user_key_id,
		 market_name, side, type, base_quantity, base_percent, currency_quantity, currency_percent, status, 
		 conditions, parent_order_id FROM orders WHERE parent_order_id = $1`, parentOrderID).
		Scan(&o.OrderID, &o.UserID, &o.KeyID, &o.MarketName, &o.Side, &o.OrderType, &o.BaseQuantity,
			&o.BasePercent, &o.CurrencyQuantity, &o.CurrencyPercent, &o.Status, &o.Conditions, &o.ParentOrderID)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
		return nil, err
	}

	return &o, nil
}

func FindOrdersByUserID(db *sql.DB, req *orderProto.GetUserOrdersRequest) ([]*orderProto.Order, error) {
	results := make([]*orderProto.Order, 0)

	rows, err := db.Query(`SELECT id, user_id, user_key_id, exchange_name, exchange_order_id, exchange_market_name,
		market_name, side, type, base_quantity, base_percent, currency_quantity, currency_percent, status, 
		conditions, condition, parent_order_id FROM orders WHERE user_id = $1`, req.UserID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o orderProto.Order
		err := rows.Scan(&o.OrderID, &o.UserID, &o.KeyID, &o.Exchange, &o.ExchangeOrderID, &o.ExchangeMarketName,
			&o.MarketName, &o.Side, &o.OrderType, &o.BaseQuantity, &o.BasePercent, &o.CurrencyQuantity,
			&o.CurrencyPercent, &o.Status, &o.Conditions, &o.Condition, &o.ParentOrderID)

		if err != nil {
			log.Fatal(err)
			return nil, err
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

func InsertOrder(db *sql.DB, req *orderProto.OrderRequest) (*orderProto.Order, error) {
	newID := uuid.New()

	jsonCond, err := json.Marshal(req.Conditions)
	if err != nil {
		return nil, err
	}

	// the exchange_order_id and exchange_market_name must be "" and not null
	// when scanning in order data null cannot be set on a type string. Therefore,
	// just default those cols to "".
	sqlStatement := `insert into orders (id, user_id, user_key_id, exchange_name, exchange_order_id, 
		exchange_market_name, market_name, side, type, base_quantity, base_percent, 
		currency_quantity, currency_percent, status, conditions, parent_order_id) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	_, err = db.Exec(sqlStatement, newID, req.UserID, req.KeyID, req.Exchange, "", "", req.MarketName,
		req.Side, req.OrderType, req.BaseQuantity, req.BasePercent, req.CurrencyQuantity, req.CurrencyPercent,
		"pending", jsonCond, req.ParentOrderID)

	if err != nil {
		return nil, err
	}
	order := &orderProto.Order{
		OrderID:          newID.String(),
		KeyID:            req.KeyID,
		UserID:           req.UserID,
		Exchange:         req.Exchange,
		MarketName:       req.MarketName,
		Side:             req.Side,
		OrderType:        req.OrderType,
		BaseQuantity:     req.BaseQuantity,
		BasePercent:      req.BasePercent,
		CurrencyQuantity: req.CurrencyQuantity,
		CurrencyPercent:  req.CurrencyPercent,
		Status:           "pending",
		Conditions:       req.Conditions,
		ParentOrderID:    req.ParentOrderID,
	}
	return order, nil
}

// UpdateOrder should be an update head chain buy order
// todo add sql for update other types of orders - UpdateHeadBuyOrder, UpdateHeadSellOrder, UpdateChainBuyOrder, UpdateChainSellOrder
func UpdateOrder(db *sql.DB, req *orderProto.OrderRequest) (*orderProto.Order, error) {
	sqlStatement := `UPDATE orders SET conditions = $1, base_quantity = $2,
	WHERE id = $4 and user_id = $5 RETURNING id, user_id, exchange_name, market_name, user_key, side, 
	base_quantity, base_quantity_remainder, status, conditions, parent_order_id`

	jsonCond, err := json.Marshal(req.Conditions)
	if err != nil {
		return nil, err
	}

	var o orderProto.Order
	err = db.QueryRow(sqlStatement, jsonCond, req.BaseQuantity, req.BaseQuantity, req.OrderID, req.UserID).
		Scan(&o.OrderID, &o.UserID, &o.Exchange, &o.MarketName, &o.KeyID,
			&o.Side, &o.BaseQuantity, &o.Status, &o.Conditions, &o.ParentOrderID)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(o.Conditions), &o.Conditions); err != nil {
		return nil, err
	}

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
