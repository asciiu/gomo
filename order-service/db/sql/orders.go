package sql

import (
	"database/sql"
	"encoding/json"
	"log"

	orderProto "github.com/asciiu/gomo/order-service/proto/order"
	"github.com/google/uuid"
)

func DeleteOrder(db *sql.DB, orderId string) error {
	_, err := db.Exec("DELETE FROM orders WHERE id = $1", orderId)
	return err
}

func FindOrderById(db *sql.DB, req *orderProto.GetUserOrderRequest) (*orderProto.Order, error) {
	var o orderProto.Order
	err := db.QueryRow("SELECT id, user_id, user_key_id, exchange_name, exchange_order_id, exchange_market_name, market_name, side, type, price, quantity, quantity_remaining, status, conditions  FROM orders WHERE id = $1", req.OrderId).
		Scan(&o.OrderId, &o.UserId, &o.ApiKeyId, &o.Exchange, &o.ExchangeOrderId, &o.ExchangeMarketName, &o.MarketName, &o.Side, &o.OrderType, &o.Price, &o.Qty, &o.QtyRemaining, &o.Status, &o.Conditions)

	if err != nil {
		return nil, err
	}
	return &o, nil
}

func FindOrdersByUserId(db *sql.DB, req *orderProto.GetUserOrdersRequest) ([]*orderProto.Order, error) {
	results := make([]*orderProto.Order, 0)

	rows, err := db.Query("SELECT id, user_id, user_key_id, exchange_name, exchange_order_id, exchange_market_name, market_name, side, type, price, quantity, quantity_remaining, status, conditions  FROM orders WHERE user_id = $1", req.UserId)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var o orderProto.Order
		err := rows.Scan(&o.OrderId, &o.UserId, &o.ApiKeyId, &o.Exchange, &o.ExchangeOrderId, &o.ExchangeMarketName, &o.MarketName, &o.Side, &o.OrderType, &o.Price, &o.Qty, &o.QtyRemaining, &o.Status, &o.Conditions)
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

	sqlStatement := `insert into orders (id, user_id, user_key_id, exchange_name, market_name, side, type, price, quantity, quantity_remaining, status, conditions) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err = db.Exec(sqlStatement, newId, req.UserId, req.ApiKeyId, req.Exchange, req.MarketName, req.Side, req.OrderType, req.Price, req.Qty, req.Qty, "pending", jsonCond)

	if err != nil {
		return nil, err
	}
	order := &orderProto.Order{
		OrderId:      newId.String(),
		ApiKeyId:     req.ApiKeyId,
		UserId:       req.UserId,
		Exchange:     req.Exchange,
		MarketName:   req.MarketName,
		Side:         req.Side,
		OrderType:    req.OrderType,
		Price:        req.Price,
		Qty:          req.Qty,
		QtyRemaining: req.Qty,
		Status:       "pending",
		Conditions:   req.Conditions,
	}
	return order, nil
}

func UpdateOrder(db *sql.DB, req *orderProto.OrderRequest) (*orderProto.Order, error) {
	sqlStatement := `UPDATE orders SET conditions = $1, price = $2, quantity = $3 WHERE id = $4 and user_id = $5 RETURNING exchange_name, api_key, status`

	var o orderProto.Order
	err := db.QueryRow(sqlStatement, req.Conditions, req.Price, req.Qty, req.OrderId, req.UserId).
		Scan(&o.Exchange, &o.ApiKeyId, &o.Status)

	if err != nil {
		return nil, err
	}
	order := &orderProto.Order{
		OrderId:    req.OrderId,
		UserId:     req.UserId,
		ApiKeyId:   o.ApiKeyId,
		Exchange:   o.Exchange,
		Status:     o.Status,
		Conditions: req.Conditions,
		Price:      req.Price,
		Qty:        req.Qty,
	}
	return order, nil
}

func UpdateApiKeyStatus(db *sql.DB, req *orderProto.OrderStatusRequest) (*orderProto.Order, error) {
	sqlStatement := `UPDATE orders SET status = $1 WHERE id = $2 RETURNING user_id, exchange_name, market_name, user_key, side, price, quantity, status, conditions`

	var o orderProto.Order
	err := db.QueryRow(sqlStatement, req.Status, req.OrderId).
		Scan(&o.UserId, &o.Exchange, &o.MarketName, &o.ApiKeyId, &o.Side, &o.Price, &o.Qty, &o.Status, &o.Conditions)

	if err != nil {
		return nil, err
	}
	order := &orderProto.Order{
		OrderId:    req.OrderId,
		ApiKeyId:   o.ApiKeyId,
		UserId:     o.UserId,
		Exchange:   o.Exchange,
		MarketName: o.MarketName,
		Side:       o.Side,
		Price:      o.Price,
		Qty:        o.Qty,
		Status:     o.Status,
		Conditions: o.Conditions,
	}
	return order, nil
}
