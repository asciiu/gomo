package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	asql "github.com/asciiu/gomo/api/db/sql"
	orderValidator "github.com/asciiu/gomo/common/constants/order"
	"github.com/asciiu/gomo/common/constants/response"
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
	// map of ticker symbol to full name
	currencies map[string]string
}

type UserOrderData struct {
	Order *Order `json:"order"`
}

type UserOrdersData struct {
	Orders []*Order `json:"orders"`
}

type Order struct {
	OrderID            string  `json:"orderID"`
	KeyID              string  `json:"keyID"`
	Exchange           string  `json:"exchange"`
	ExchangeOrderID    string  `json:"exchangeOrderID"`
	ExchangeMarketName string  `json:"exchangeMarketName"`
	MarketName         string  `json:"marketName"`
	MarketCurrency     string  `json:"marketCurrency"`
	MarketCurrencyLong string  `json:"marketCurrencyLong"`
	Side               string  `json:"side"`
	OrderType          string  `json:"orderType"`
	BaseCurrency       string  `json:"baseCurrency"`
	BaseCurrencyLong   string  `json:"baseCurrencyLong"`
	BaseQuantity       float64 `json:"baseQuantity"`
	BasePercent        float64 `json:"basePercent"`
	CurrencyQuantity   float64 `json:"currencyQuantity"`
	CurrencyPercent    float64 `json:"currencyPercent"`
	Status             string  `json:"status"`
	Conditions         string  `json:"conditions"`
	Condition          string  `json:"condition"`
	ParentOrderID      string  `json:"parentOrderID"`
}

// swagger:parameters addOrder
type OrderRequest struct {
	// Required internal api key ID
	// in: body
	KeyID string `json:"keyID"`
	// Required e.g. ADA-BTC
	// in: body
	MarketName string `json:"marketName"`
	// Required "buy" or "sell"
	// in: body
	Side string `json:"side"`
	// Required Valid order types are "market", "limit", "virtual". Orders not within these types will be ignored.
	// in: body
	OrderType string `json:"orderType"`
	// Required for buy side when order is first in chain
	// in: body
	BaseQuantity float64 `json:"baseQuantity"`
	// Required for buy side on chained orders
	// in: body
	BasePercent float64 `json:"basePercent"`
	// Required for sell side when an order is first in a chain
	// in: body
	CurrencyQuantity float64 `json:"currencyQuantity"`
	// Required for sell side for all orders that are chained
	// in: body
	CurrencyPercent float64 `json:"currencyPercent"`
	// Required
	// in: body
	Conditions string `json:"conditions"`

	// Optional parent order ID to add this chain of orders to
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
	baseCurrency := names[1]
	baseCurrencyLong := controller.currencies[baseCurrency]
	marketCurrency := names[0]
	marketCurrencyLong := controller.currencies[marketCurrency]

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
				MarketCurrency:     marketCurrency,
				MarketCurrencyLong: marketCurrencyLong,
				Side:               r.Data.Order.Side,
				OrderType:          r.Data.Order.OrderType,
				BaseCurrency:       baseCurrency,
				BaseCurrencyLong:   baseCurrencyLong,
				BaseQuantity:       r.Data.Order.BaseQuantity,
				BasePercent:        r.Data.Order.BasePercent,
				CurrencyQuantity:   r.Data.Order.CurrencyQuantity,
				CurrencyPercent:    r.Data.Order.CurrencyPercent,
				Status:             r.Data.Order.Status,
				Conditions:         r.Data.Order.Conditions,
				Condition:          r.Data.Order.Condition,
				ParentOrderID:      r.Data.Order.ParentOrderID,
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
		baseCurrency := names[1]
		baseCurrencyLong := controller.currencies[baseCurrency]
		marketCurrency := names[0]
		marketCurrencyLong := controller.currencies[marketCurrency]

		data[i] = &Order{
			OrderID:            o.OrderID,
			KeyID:              o.KeyID,
			Exchange:           o.Exchange,
			ExchangeOrderID:    o.ExchangeOrderID,
			ExchangeMarketName: o.ExchangeMarketName,
			MarketName:         o.MarketName,
			MarketCurrency:     marketCurrency,
			MarketCurrencyLong: marketCurrencyLong,
			Side:               o.Side,
			OrderType:          o.OrderType,
			BaseCurrency:       baseCurrency,
			BaseCurrencyLong:   baseCurrencyLong,
			BaseQuantity:       o.BaseQuantity,
			BasePercent:        o.BasePercent,
			CurrencyQuantity:   o.CurrencyQuantity,
			CurrencyPercent:    o.CurrencyPercent,
			Status:             o.Status,
			Conditions:         o.Conditions,
			Condition:          o.Condition,
			ParentOrderID:      o.ParentOrderID,
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
// create a new order  (protected)
//
// This will create a new order in the system.
// Example request:
// [
//	{
//		"keyID": "680d6bbf-1feb-4122-bd10-0e7ce080676a",
//		"marketName": "ADA-BTC",
//		"side": "buy",
//		"basePercent": 0.50,
//		"orderType": "market",
//		"conditions": "price <= 0.00002800"
//	},
//	{
//		"keyID": "680d6bbf-1feb-4122-bd10-0e7ce080676a",
//		"marketName": "ADA-BTC",
//		"side": "buy",
//		"basePercent": 1.0,
//		"orderType": "market",
//		"conditions": "price <= 0.00002200"
//	}
// ]
//
// responses:
//  200: responseOrdersSuccess "data" will contain list of orders with "status": "success"
//  400: responseError missing params with "status": "fail"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *OrderController) HandlePostOrder(c echo.Context) error {
	defer c.Request().Body.Close()
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

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
			return fail(c, "what kind of order type is this?")
		}

		// side, market name, and api key are required
		if order.Side == "" || order.MarketName == "" || order.KeyID == "" {
			return fail(c, "side, marketName, and keyID required!")
		}

		// assume the first order is head of a chain if the ParentOrderID is empty
		// this means that a new chain of orders has been submitted because the
		// ParentOrderID has not been assigned yet.
		if i == 0 && order.ParentOrderID == "" && order.Side == "buy" && order.BasePercent == 0.0 {
			return fail(c, "head buy in chain requires a basePercent")
		}

		// if the head order side is sell we need a currency quantity
		if i == 0 && order.ParentOrderID == "" && order.Side == "sell" && order.CurrencyPercent == 0.0 {
			return fail(c, "head sell in chain requires a currencyPercent")
		}

		// need to use basePercent for chained buys
		//if i != 0 && order.Side == "buy" && order.BasePercent == 0.0 {
		//	return fail(c, "chained buys require a basePercent")
		//}

		//// need to use currencyPercent for chained buys
		//if i != 0 && order.Side == "sell" && order.CurrencyQuantity == 0.0 {
		//	return fail(c, "chained sells require a currencyPercent")
		//}

		// market name should be formatted as
		// currency-base (e.g. ADA-BTC)
		if !strings.Contains(order.MarketName, "-") {
			return fail(c, "marketName must be currency-base: e.g. ADA-BTC")
		}

		if order.ParentOrderID == "" {
			order.ParentOrderID = "00000000-0000-0000-0000-000000000000"
		}

		request := orders.OrderRequest{
			UserID:           userID,
			KeyID:            order.KeyID,
			MarketName:       order.MarketName,
			Side:             order.Side,
			Conditions:       order.Conditions,
			OrderType:        order.OrderType,
			BaseQuantity:     order.BaseQuantity,
			BasePercent:      order.BasePercent,
			CurrencyQuantity: order.CurrencyQuantity,
			CurrencyPercent:  order.CurrencyPercent,
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
		baseCurrency := names[1]
		baseCurrencyLong := controller.currencies[baseCurrency]
		marketCurrency := names[0]
		marketCurrencyLong := controller.currencies[marketCurrency]

		data[i] = &Order{
			OrderID:            o.OrderID,
			KeyID:              o.KeyID,
			Exchange:           o.Exchange,
			ExchangeOrderID:    o.ExchangeOrderID,
			ExchangeMarketName: o.ExchangeMarketName,
			MarketName:         o.MarketName,
			MarketCurrency:     marketCurrency,
			MarketCurrencyLong: marketCurrencyLong,
			Side:               o.Side,
			OrderType:          o.OrderType,
			BaseQuantity:       o.BaseQuantity,
			BasePercent:        o.BasePercent,
			BaseCurrency:       baseCurrency,
			BaseCurrencyLong:   baseCurrencyLong,
			CurrencyQuantity:   o.CurrencyQuantity,
			CurrencyPercent:    o.CurrencyPercent,
			Status:             o.Status,
			Conditions:         o.Conditions,
			Condition:          o.Condition,
			ParentOrderID:      o.ParentOrderID,
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
	baseCurrency := names[1]
	baseCurrencyLong := controller.currencies[baseCurrency]
	marketCurrency := names[0]
	marketCurrencyLong := controller.currencies[marketCurrency]

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
				MarketCurrency:     marketCurrency,
				MarketCurrencyLong: marketCurrencyLong,
				Side:               r.Data.Order.Side,
				OrderType:          r.Data.Order.OrderType,
				BaseQuantity:       r.Data.Order.BaseQuantity,
				BasePercent:        r.Data.Order.BasePercent,
				BaseCurrency:       baseCurrency,
				BaseCurrencyLong:   baseCurrencyLong,
				CurrencyQuantity:   r.Data.Order.CurrencyQuantity,
				CurrencyPercent:    r.Data.Order.CurrencyPercent,
				Status:             r.Data.Order.Status,
				Conditions:         r.Data.Order.Conditions,
				Condition:          r.Data.Order.Condition,
				ParentOrderID:      r.Data.Order.ParentOrderID,
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
