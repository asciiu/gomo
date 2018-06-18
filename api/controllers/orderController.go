package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/asciiu/gomo/common/constants/key"

	asql "github.com/asciiu/gomo/api/db/sql"
	orderValidator "github.com/asciiu/gomo/common/constants/order"
	"github.com/asciiu/gomo/common/constants/response"
	"github.com/asciiu/gomo/common/constants/side"
	keys "github.com/asciiu/gomo/key-service/proto/key"
	orders "github.com/asciiu/gomo/order-service/proto/order"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
	"golang.org/x/net/context"
)

type OrderController struct {
	DB     *sql.DB
	Orders orders.OrderServiceClient
	Keys   keys.KeyServiceClient
	// map of ticker symbol to full name
	currencies map[string]string
}

type UserOrderData struct {
	Order *Order `json:"order"`
}

type UserOrdersData struct {
	Orders []*Order `json:"orders"`
}

// This is the response struct for order
type Order struct {
	ParentOrderID      string  `json:"parentOrderID"`
	OrderID            string  `json:"orderID"`
	KeyID              string  `json:"keyID"`
	Exchange           string  `json:"exchange"`
	ExchangeOrderID    string  `json:"exchangeOrderID"`
	ExchangeMarketName string  `json:"exchangeMarketName"`
	MarketName         string  `json:"marketName"`
	Side               string  `json:"side"`
	OrderType          string  `json:"orderType"`
	Price              float64 `json:"price"`
	BaseCurrencySymbol string  `json:"baseCurrencySymbol"`
	BaseCurrencyName   string  `json:"baseCurrencyName"`
	BaseQuantity       float64 `json:"baseQuantity"`
	BasePercent        float64 `json:"basePercent"`
	CurrencySymbol     string  `json:"currencySymbol"`
	CurrencyName       string  `json:"currencyName"`
	CurrencyQuantity   float64 `json:"currencyQuantity"`
	CurrencyPercent    float64 `json:"currencyPercent"`
	Status             string  `json:"status"`
	ChainStatus        string  `json:"chainStatus"`
	Conditions         string  `json:"conditions"`
	Condition          string  `json:"condition"`
}

// swagger:parameters addOrder
type OrderRequest struct {
	// Required this is our api key ID (string uuid) assigned to the user's exchange key and secret.
	// in: body
	KeyID string `json:"keyID"`
	// Required e.g. ADA-BTC. Base pair should be the suffix.
	// in: body
	MarketName string `json:"marketName"`
	// Required "buy" or "sell"
	// in: body
	Side string `json:"side"`
	// Required Valid order types are "market", "limit", "virtual". Orders not within these types will be rejected.
	// in: body
	OrderType string `json:"orderType"`
	// Required for buy side when order is first in chain. This will designate the reserve base balance to use during the execution of the chained order stategy.
	// in: body
	BaseQuantity float64 `json:"baseQuantity"`
	// Required for buy side on chained orders. Specifies the precent of the reserve balance to allocate for the order.
	// in: body
	BasePercent float64 `json:"basePercent"`
	// Required for sell side when an order is first in a chain. This is the quantity of market currency to sell.
	// in: body
	CurrencyQuantity float64 `json:"currencyQuantity"`
	// Required for sell side for all orders that are chained. Similar to the basePercent, but for sell orders.
	// in: body
	CurrencyPercent float64 `json:"currencyPercent"`
	// Required these are the conditions that trigger the order to execute. Example: ???
	// in: body
	Conditions string `json:"conditions"`
	// Optional this is required only when the order type is 'limit'. This is the limit order price.
	// in: body
	Price float64 `json:"price"`
	// Optional to set the chain status to active (default true for active)
	Active bool `json:"active"`

	// Optional parent order ID to add this chain of orders to. When you want to add children to an existing order.
	ParentOrderID string `json:"parentOrderID"`
}

// swagger:parameters updateOrder
type UpdateOrderRequest struct {
	// Optional.
	// in: body
	OrderType string `json:"orderType"`
	// Optional.
	// in: body
	Price float64 `json:"price"`
	// Optional.
	// in: body
	BaseQuantity float64 `json:"baseQuantity"`
	// Optional.
	// in: body
	Conditions string `json:"conditions"`
}

// A ResponseKeySuccess will always contain a status of "successful".
// swagger:model responseOrderSuccess
type ResponseOrderSuccess struct {
	Status string         `json:"status"`
	Data   *UserOrderData `json:"data"`
}

// A ResponseKeysSuccess will always contain a status of "successful".
// swagger:model responseOrdersSuccess
type ResponseOrdersSuccess struct {
	Status string          `json:"status"`
	Data   *UserOrdersData `json:"data"`
}

func NewOrderController(db *sql.DB) *OrderController {
	// Create a new service. Optionally include some options here.
	service := k8s.NewService(micro.Name("apikey.client"))
	service.Init()

	controller := OrderController{
		DB:         db,
		Orders:     orders.NewOrderServiceClient("orders", service.Client()),
		Keys:       keys.NewKeyServiceClient("keys", service.Client()),
		currencies: make(map[string]string),
	}

	currencies, err := asql.GetCurrencyNames(db)
	switch {
	case err == sql.ErrNoRows:
		log.Println("Quaid, you need to populate the currency_names table!")
	case err != nil:
	default:
		for _, c := range currencies {
			controller.currencies[c.TickerSymbol] = c.CurrencyName
		}
	}

	return &controller

}

// swagger:route GET /orders/:orderID orders getOrder
//
// show order (protected)
//
// Get info about an order.
//
// responses:
//  200: responseOrderSuccess "data" will contain order stuffs with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *OrderController) HandleGetOrder(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	orderID := c.Param("orderID")

	getRequest := orders.GetUserOrderRequest{
		OrderID: orderID,
		UserID:  userID,
	}

	r, _ := controller.Orders.GetUserOrder(context.Background(), &getRequest)
	if r.Status != response.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == response.Fail {
			return c.JSON(http.StatusBadRequest, res)
		}
		if r.Status == response.Error {
			return c.JSON(http.StatusInternalServerError, res)
		}
	}

	names := strings.Split(r.Data.Order.MarketName, "-")
	baseCurrencySymbol := names[1]
	baseCurrencyName := controller.currencies[baseCurrencySymbol]
	currencySymbol := names[0]
	currencyName := controller.currencies[currencySymbol]

	res := &ResponseOrderSuccess{
		Status: response.Success,
		Data: &UserOrderData{
			Order: &Order{
				OrderID:            r.Data.Order.OrderID,
				KeyID:              r.Data.Order.KeyID,
				Exchange:           r.Data.Order.Exchange,
				ExchangeOrderID:    r.Data.Order.ExchangeOrderID,
				ExchangeMarketName: r.Data.Order.ExchangeMarketName,
				MarketName:         r.Data.Order.MarketName,
				Side:               r.Data.Order.Side,
				OrderType:          r.Data.Order.OrderType,
				BaseCurrencySymbol: baseCurrencySymbol,
				BaseCurrencyName:   baseCurrencyName,
				BaseQuantity:       r.Data.Order.BaseQuantity,
				BasePercent:        r.Data.Order.BasePercent,
				CurrencySymbol:     currencySymbol,
				CurrencyName:       currencyName,
				CurrencyQuantity:   r.Data.Order.CurrencyQuantity,
				CurrencyPercent:    r.Data.Order.CurrencyPercent,
				Status:             r.Data.Order.Status,
				Conditions:         r.Data.Order.Conditions,
				Condition:          r.Data.Order.Condition,
				ParentOrderID:      r.Data.Order.ParentOrderID,
				Price:              r.Data.Order.Price,
			},
		},
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:route GET /orders orders getAllOrders
//
// get all orders (protected)
//
// Currently returns all orders. Eventually going to add params to filter orders.
//
// responses:
//  200: responseOrdersSuccess "data" will contain a list of order info with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *OrderController) HandleListOrders(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	getRequest := orders.GetUserOrdersRequest{
		UserID: userID,
	}

	r, _ := controller.Orders.GetUserOrders(context.Background(), &getRequest)
	if r.Status != response.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == response.Fail {
			return c.JSON(http.StatusBadRequest, res)
		}
		if r.Status == response.Error {
			return c.JSON(http.StatusInternalServerError, res)
		}
	}

	data := make([]*Order, len(r.Data.Orders))
	for i, o := range r.Data.Orders {

		names := strings.Split(o.MarketName, "-")
		baseCurrencySymbol := names[1]
		baseCurrencyName := controller.currencies[baseCurrencySymbol]
		currencySymbol := names[0]
		currencyName := controller.currencies[currencySymbol]

		data[i] = &Order{
			OrderID:            o.OrderID,
			KeyID:              o.KeyID,
			Exchange:           o.Exchange,
			ExchangeOrderID:    o.ExchangeOrderID,
			ExchangeMarketName: o.ExchangeMarketName,
			MarketName:         o.MarketName,
			Side:               o.Side,
			OrderType:          o.OrderType,
			BaseCurrencySymbol: baseCurrencySymbol,
			BaseCurrencyName:   baseCurrencyName,
			BaseQuantity:       o.BaseQuantity,
			BasePercent:        o.BasePercent,
			CurrencySymbol:     currencySymbol,
			CurrencyName:       currencyName,
			CurrencyQuantity:   o.CurrencyQuantity,
			CurrencyPercent:    o.CurrencyPercent,
			Status:             o.Status,
			Conditions:         o.Conditions,
			Condition:          o.Condition,
			ParentOrderID:      o.ParentOrderID,
			Price:              o.Price,
		}
	}

	res := &ResponseOrdersSuccess{
		Status: response.Success,
		Data: &UserOrdersData{
			Orders: data,
		},
	}

	return c.JSON(http.StatusOK, res)
}

func fail(c echo.Context, msg string) error {
	res := &ResponseError{
		Status:  response.Fail,
		Message: msg,
	}

	return c.JSON(http.StatusBadRequest, res)
}

// swagger:route POST /orders orders addOrder
//
// create a new order chain (protected)
//
// This will create a new chain of orders for the user.
// Example:
// [
//	{
//		"keyID": "680d6bbf-1feb-4122-bd10-0e7ce080676a",
//		"marketName": "ADA-BTC",
//		"side": "buy",
//		"baseQuantity": 0.50,
//		"orderType": "market",
//		"conditions": "price <= 0.00002800"
//	},
//	{
//		"keyID": "680d6bbf-1feb-4122-bd10-0e7ce080676a",
//		"marketName": "ADA-BTC",
//		"side": "sell",
//		"currencyPercent": 1.0,
//		"orderType": "market",
//		"conditions": "price <= 0.00002200 or trailingStopPts(0, 0.0)"
//	}
// ]
//
// The order chain starts out with a reserve base balance of 0.5 BTC. This order chain will buy 0.5 BTC at
// market price when the market price reaches 2800 satoshi or less. The following sell order will sell 100% of the order
// strategy's (i.e. chain of orders) cardano balance. Cardano balance should in theory be dictated by the Cardano that this
// chain bought - 100% does not mean 100% of user's cardano balance. The conditions will likely be a json string of conditionLabel: condition.
// example: "stopLoss: price <= 0.00002200, takeProfit: ...."
//
// responses:
//  200: responseOrdersSuccess "data" will contain list of orders with "status": "success"
//  400: responseError missing or incorrect params with "status": "fail"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *OrderController) HandlePostOrder(c echo.Context) error {
	defer c.Request().Body.Close()
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	keyID := ""
	exchangeName := ""

	ordrs := make([]*OrderRequest, 0)
	requests := make([]*orders.OrderRequest, 0)
	dec := json.NewDecoder(c.Request().Body)

	_, err := dec.Token()
	if err != nil {
		return fail(c, err.Error())
	}

	// read all orders from array
	for dec.More() {
		var o OrderRequest

		if err := dec.Decode(&o); err != nil {
			return fail(c, "expected an array")
		}
		ordrs = append(ordrs, &o)
	}

	// error check all orders
	for i, order := range ordrs {
		if !orderValidator.ValidateOrderType(order.OrderType) {
			return fail(c, "market, limit, or virtual orders only!")
		}

		// side, market name, and api key are required
		if order.Side == "" || order.MarketName == "" || order.KeyID == "" {
			return fail(c, "side, marketName, and keyID required!")
		}

		// assume the first order is head of a chain if the ParentOrderID is empty
		// this means that a new chain of orders has been submitted because the
		// ParentOrderID has not been assigned yet.
		if i == 0 && order.ParentOrderID == "" && order.Side == side.Buy && order.BaseQuantity == 0.0 {
			return fail(c, "first buy in order chain requires a baseQuantity")
		}

		// if the head order side is sell we need a currency quantity
		if i == 0 && order.ParentOrderID == "" && order.Side == side.Sell && order.CurrencyQuantity == 0.0 {
			return fail(c, "first sell in order chain requires a currencyQuantity")
		}

		// need to use basePercent for chained buys
		if i != 0 && order.Side == side.Buy && order.BasePercent == 0.0 {
			return fail(c, "child buy orders require a basePercent")
		}

		// need to use currencyPercent for chained sells
		if i != 0 && order.Side == side.Sell && order.CurrencyQuantity == 0.0 {
			return fail(c, "child sell orders require a currencyPercent")
		}

		// market name should be formatted as
		// currency-base (e.g. ADA-BTC)
		if !strings.Contains(order.MarketName, "-") {
			return fail(c, "marketName must be currency-base: e.g. ADA-BTC")
		}

		if order.ParentOrderID == "" {
			order.ParentOrderID = "00000000-0000-0000-0000-000000000000"
		}

		// on first iteration of loop set the keyID
		if keyID == "" {
			keyID = order.KeyID
			getRequest := keys.GetUserKeyRequest{
				KeyID:  keyID,
				UserID: userID,
			}

			// ask key service for key
			r, _ := controller.Keys.GetUserKey(context.Background(), &getRequest)
			if r.Status != response.Success {
				res := &ResponseError{
					Status:  r.Status,
					Message: r.Message,
				}

				if r.Status == response.Fail {
					return c.JSON(http.StatusBadRequest, res)
				}
				if r.Status == response.Error {
					return c.JSON(http.StatusInternalServerError, res)
				}
				if r.Status == response.Nonentity {
					return fail(c, "invalid key")
				}
			}
			// if key found it must be verified status
			if r.Data.Key.Status != key.Verified {
				return fail(c, "invalid key")
			}
			exchangeName = r.Data.Key.Exchange
		}

		// all orders in the chain must use the same key ID
		if keyID != "" && order.KeyID != keyID {
			return fail(c, "all orders must use the same keyID")
		}

		request := orders.OrderRequest{
			UserID:           userID,
			KeyID:            order.KeyID,
			Exchange:         exchangeName,
			MarketName:       order.MarketName,
			Side:             order.Side,
			OrderType:        order.OrderType,
			BaseQuantity:     order.BaseQuantity,
			BasePercent:      order.BasePercent,
			CurrencyQuantity: order.CurrencyQuantity,
			CurrencyPercent:  order.CurrencyPercent,
			Conditions:       order.Conditions,
			Price:            order.Price,
			ParentOrderID:    order.ParentOrderID,
		}
		requests = append(requests, &request)
	}

	orderRequests := orders.OrdersRequest{
		Orders: requests,
	}

	// add order returns nil for error
	r, _ := controller.Orders.AddOrders(context.Background(), &orderRequests)
	if r.Status != response.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == response.Fail {
			return c.JSON(http.StatusBadRequest, res)
		}
		if r.Status == response.Error {
			return c.JSON(http.StatusInternalServerError, res)
		}
	}

	data := make([]*Order, len(r.Data.Orders))
	for i, o := range r.Data.Orders {
		names := strings.Split(o.MarketName, "-")
		baseCurrencySymbol := names[1]
		baseCurrencyName := controller.currencies[baseCurrencySymbol]
		currencySymbol := names[0]
		currencyName := controller.currencies[currencySymbol]

		data[i] = &Order{
			OrderID:            o.OrderID,
			KeyID:              o.KeyID,
			Exchange:           o.Exchange,
			ExchangeOrderID:    o.ExchangeOrderID,
			ExchangeMarketName: o.ExchangeMarketName,
			MarketName:         o.MarketName,
			Side:               o.Side,
			OrderType:          o.OrderType,
			BaseCurrencySymbol: baseCurrencySymbol,
			BaseCurrencyName:   baseCurrencyName,
			BaseQuantity:       o.BaseQuantity,
			BasePercent:        o.BasePercent,
			CurrencySymbol:     currencySymbol,
			CurrencyName:       currencyName,
			CurrencyQuantity:   o.CurrencyQuantity,
			CurrencyPercent:    o.CurrencyPercent,
			Status:             o.Status,
			Conditions:         o.Conditions,
			Condition:          o.Condition,
			ParentOrderID:      o.ParentOrderID,
			Price:              o.Price,
		}
	}

	res := &ResponseOrdersSuccess{
		Status: response.Success,
		Data: &UserOrdersData{
			Orders: data,
		},
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:route PUT /orders/:orderID orders updateOrder
//
// update and order (protected)
//
// You can only update pending orders.
//
// responses:
//  200: responseOrderSuccess "data" will contain order info with "status": "success"
//  400: responseError missing params with "status": "fail"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *OrderController) HandleUpdateOrder(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	orderID := c.Param("orderID")

	orderRequest := UpdateOrderRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&orderRequest)
	if err != nil {
		res := &ResponseError{
			Status:  response.Fail,
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	// client can only update description
	updateRequest := orders.OrderRequest{
		OrderID:      orderID,
		UserID:       userID,
		Conditions:   orderRequest.Conditions,
		BaseQuantity: orderRequest.BaseQuantity,
	}

	r, _ := controller.Orders.UpdateOrder(context.Background(), &updateRequest)
	if r.Status != response.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == response.Fail {
			return c.JSON(http.StatusBadRequest, res)
		}
		if r.Status == response.Error {
			return c.JSON(http.StatusInternalServerError, res)
		}
	}

	names := strings.Split(r.Data.Order.MarketName, "-")
	baseCurrencySymbol := names[1]
	baseCurrencyName := controller.currencies[baseCurrencySymbol]
	currencySymbol := names[0]
	currencyName := controller.currencies[currencySymbol]

	res := &ResponseOrderSuccess{
		Status: response.Success,
		Data: &UserOrderData{
			Order: &Order{
				OrderID:            r.Data.Order.OrderID,
				KeyID:              r.Data.Order.KeyID,
				Exchange:           r.Data.Order.Exchange,
				ExchangeOrderID:    r.Data.Order.ExchangeOrderID,
				ExchangeMarketName: r.Data.Order.ExchangeMarketName,
				MarketName:         r.Data.Order.MarketName,
				Side:               r.Data.Order.Side,
				OrderType:          r.Data.Order.OrderType,
				BaseCurrencySymbol: baseCurrencySymbol,
				BaseCurrencyName:   baseCurrencyName,
				BaseQuantity:       r.Data.Order.BaseQuantity,
				BasePercent:        r.Data.Order.BasePercent,
				CurrencySymbol:     currencySymbol,
				CurrencyName:       currencyName,
				CurrencyQuantity:   r.Data.Order.CurrencyQuantity,
				CurrencyPercent:    r.Data.Order.CurrencyPercent,
				Status:             r.Data.Order.Status,
				Conditions:         r.Data.Order.Conditions,
				Condition:          r.Data.Order.Condition,
				ParentOrderID:      r.Data.Order.ParentOrderID,
				Price:              r.Data.Order.Price,
			},
		},
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:route DELETE /orders/:orderID orders deleteOrder
//
// Remove and order (protected)
//
// Cannot remove orders that have already executed.
//
// responses:
//  200: responseOrderSuccess data will be null with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *OrderController) HandleDeleteOrder(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	orderID := c.Param("orderID")

	removeRequest := orders.RemoveOrderRequest{
		OrderID: orderID,
		UserID:  userID,
	}

	r, _ := controller.Orders.RemoveOrder(context.Background(), &removeRequest)
	if r.Status != response.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == response.Fail {
			return c.JSON(http.StatusBadRequest, res)
		}
		if r.Status == response.Error {
			return c.JSON(http.StatusInternalServerError, res)
		}
	}

	res := &ResponseOrderSuccess{
		Status: response.Success,
	}

	return c.JSON(http.StatusOK, res)
}
