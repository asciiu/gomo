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
func DeleteOrder(db *sql.DB, orderId string) error {
	_, err := db.Exec("DELETE FROM orders WHERE id = $1", orderId)
	return err
}

// FindOrderById ...
func FindOrderById(db *sql.DB, req *orderProto.GetUserOrderRequest) (*orderProto.Order, error) {
	var o orderProto.Order
	err := db.QueryRow(`SELECT id, user_id, user_key_id, exchange_name, exchange_order_id, exchange_market_name,
		 market_name, side, type, base_quantity, base_quantity_remainder, currency_quantity, status, 
		 conditions, condition, parent_order_id FROM order WHERE id = $1`, req.OrderId).
		Scan(&o.OrderId, &o.UserId, &o.ApiKeyId, &o.Exchange, &o.ExchangeOrderId, &o.ExchangeMarketName,
			&o.MarketName, &o.Side, &o.OrderType, &o.BaseQuantity, &o.BaseQuantityRemainder, &o.CurrencyQuantity,
			&o.Status, &o.Conditions, &o.Condition, &o.ParentOrderId)

	if err != nil {
		return nil, err
	}
	return &o, nil
}

func FindOrdersByUserId(db *sql.DB, req *orderProto.GetUserOrdersRequest) ([]*orderProto.Order, error) {
	results := make([]*orderProto.Order, 0)

	rows, err := db.Query(`SELECT id, user_id, user_key_id, exchange_name, exchange_order_id, exchange_market_name,
		market_name, side, type, base_quantity, base_quantity_remainder, currency_quantity, status, 
		conditions, condition, parent_order_id FROM order WHERE user_id = $1`, req.UserId)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o orderProto.Order
		err := rows.Scan(&o.OrderId, &o.UserId, &o.ApiKeyId, &o.Exchange, &o.ExchangeOrderId, &o.ExchangeMarketName,
			&o.MarketName, &o.Side, &o.OrderType, &o.BaseQuantity, &o.BaseQuantityRemainder, &o.CurrencyQuantity,
			&o.Status, &o.Conditions, &o.Condition, &o.ParentOrderId)

		if err != nil {
			log.Fatal(err)
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
	newId := uuid.New()

	jsonCond, err := json.Marshal(req.Conditions)
	if err != nil {
		return nil, err
	}

	// the exchange_order_id and exchange_market_name must be "" and not null
	// when scanning in order data null cannot be set on a type string. Therefore,
	// just default those cols to "".
	sqlStatement := `insert into orders (id, user_id, user_key_id, exchange_name, exchange_order_id, 
		exchange_market_name, market_name, side, type, base_quantity, base_quantity_remainder, status, 
		conditions, parent_order_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err = db.Exec(sqlStatement, newId, req.UserId, req.ApiKeyId, req.Exchange, "", "", req.MarketName,
		req.Side, req.OrderType, req.BaseQuantity, req.BaseQuantity, "pending", jsonCond, req.ParentOrderId)

	if err != nil {
		return nil, err
	}
	order := &orderProto.Order{
		OrderId:               newId.String(),
		ApiKeyId:              req.ApiKeyId,
		UserId:                req.UserId,
		Exchange:              req.Exchange,
		MarketName:            req.MarketName,
		Side:                  req.Side,
		OrderType:             req.OrderType,
		BaseQuantity:          req.BaseQuantity,
		BaseQuantityRemainder: req.BaseQuantity,
		Status:                "pending",
		Conditions:            req.Conditions,
		ParentOrderId:         req.ParentOrderId,
	}
	return order, nil
}

func UpdateOrder(db *sql.DB, req *orderProto.OrderRequest) (*orderProto.Order, error) {
	sqlStatement := `UPDATE orders SET conditions = $1, base_quantity = $2, base_quantity_remainder = $3 
	WHERE id = $4 and user_id = $5 RETURNING id, user_id, exchange_name, market_name, user_key, side, 
	base_quantity, base_quantity_remainder, status, conditions, parent_order_id`

	jsonCond, err := json.Marshal(req.Conditions)
	if err != nil {
		return nil, err
	}

	var o orderProto.Order
	err = db.QueryRow(sqlStatement, jsonCond, req.BaseQuantity, req.BaseQuantity, req.OrderId, req.UserId).
		Scan(&o.OrderId, &o.UserId, &o.Exchange, &o.MarketName, &o.ApiKeyId,
			&o.Side, &o.BaseQuantity, &o.BaseQuantityRemainder, &o.Status, &o.Conditions, &o.ParentOrderId)

	if err != nil {
		return nil, err
	}
	return &o, nil
}

func UpdateOrderStatus(db *sql.DB, evt *evt.OrderEvent) (*orderProto.Order, error) {
	sqlStatement := `UPDATE orders SET status = $1, condition = $2, exchange_order_id = $3, exchange_market_name = $4 
	WHERE id = $5 RETURNING id, user_id, exchange_name, market_name, user_key_id, side, base_quantity, 
	base_quantity_remainder, status, conditions`

	var o orderProto.Order
	err := db.QueryRow(sqlStatement, evt.Status, evt.Condition, evt.ExchangeOrderId,
		evt.ExchangeMarketName, evt.OrderId).
		Scan(&o.OrderId, &o.UserId, &o.Exchange, &o.MarketName, &o.ApiKeyId,
			&o.Side, &o.BaseQuantity, &o.BaseQuantityRemainder, &o.Status, &o.Conditions)

	if err != nil {
		return nil, err
	}
	return &o, nil
}
