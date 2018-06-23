package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	asql "github.com/asciiu/gomo/api/db/sql"
	"github.com/asciiu/gomo/common/constants/key"
	orderValidator "github.com/asciiu/gomo/common/constants/order"
	"github.com/asciiu/gomo/common/constants/response"
	sideValidator "github.com/asciiu/gomo/common/constants/side"
	keys "github.com/asciiu/gomo/key-service/proto/key"
	orders "github.com/asciiu/gomo/order-service/proto/order"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
	"golang.org/x/net/context"
)

type TradeController struct {
	DB     *sql.DB
	Orders orders.OrderServiceClient
	Keys   keys.KeyServiceClient
	// map of ticker symbol to full name
	currencies map[string]string
}

type UserStrategyData struct {
	Strategy *Strategy `json:"strategy"`
}

type UserStrategiesData struct {
	Strategies []*Strategy `json:"strategies"`
}

// This is the response struct for order
type Strategy struct {
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
type StrategyRequest struct {
	// Required this is our api key ID (string uuid) assigned to the user's exchange key and secret.
	// in: body
	KeyID string `json:"keyID"`
	// Required e.g. ADA-BTC. Base pair should be the suffix.
	// in: body
	MarketName string `json:"marketName"`
	// Required When first order is buy. Base should be in suffix of market name.
	// in: body
	BaseBalance float64 `json:"baseBalance"`
	// Required When first order is buy. Base should be in suffix of market name.
	// in: body
	CurrencyBalance float64 `json:"currencyBalance"`
	// Optional make this strategy live immeadiately - live (true) or make inactive (false)
	// in: body
	Live bool `json:"live"`

	// Required array of orders. The sequence of orders is assumed to be the sequence of execution.
	// in: body
	Orders []*OrderReq `json:"orders"`
}

type OrderReq struct {
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
}

// swagger:parameters updateOrder
type UpdateStrategyRequest struct {
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
type ResponseStrategySuccess struct {
	Status string         `json:"status"`
	Data   *UserOrderData `json:"data"`
}

// A ResponseKeysSuccess will always contain a status of "successful".
// swagger:model responseOrdersSuccess
type ResponseStategiesSuccess struct {
	Status string          `json:"status"`
	Data   *UserOrdersData `json:"data"`
}

func NewTradeController(db *sql.DB) *TradeController {
	// Create a new service. Optionally include some options here.
	service := k8s.NewService(micro.Name("apikey.client"))
	service.Init()

	controller := TradeController{
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
func (controller *TradeController) HandleGetTrade(c echo.Context) error {
	return nil
	// token := c.Get("user").(*jwt.Token)
	// claims := token.Claims.(jwt.MapClaims)
	// userID := claims["jti"].(string)
	// orderID := c.Param("orderID")

	// getRequest := orders.GetUserOrderRequest{
	// 	OrderID: orderID,
	// 	UserID:  userID,
	// }

	// r, _ := controller.Orders.GetUserOrder(context.Background(), &getRequest)
	// if r.Status != response.Success {
	// 	res := &ResponseError{
	// 		Status:  r.Status,
	// 		Message: r.Message,
	// 	}

	// 	if r.Status == response.Fail {
	// 		return c.JSON(http.StatusBadRequest, res)
	// 	}
	// 	if r.Status == response.Error {
	// 		return c.JSON(http.StatusInternalServerError, res)
	// 	}
	// }

	// names := strings.Split(r.Data.Order.MarketName, "-")
	// baseCurrencySymbol := names[1]
	// baseCurrencyName := controller.currencies[baseCurrencySymbol]
	// currencySymbol := names[0]
	// currencyName := controller.currencies[currencySymbol]

	// res := &ResponseOrderSuccess{
	// 	Status: response.Success,
	// 	Data: &UserOrderData{
	// 		Order: &Order{
	// 			OrderID:            r.Data.Order.OrderID,
	// 			KeyID:              r.Data.Order.KeyID,
	// 			Exchange:           r.Data.Order.Exchange,
	// 			ExchangeOrderID:    r.Data.Order.ExchangeOrderID,
	// 			ExchangeMarketName: r.Data.Order.ExchangeMarketName,
	// 			MarketName:         r.Data.Order.MarketName,
	// 			Side:               r.Data.Order.Side,
	// 			OrderType:          r.Data.Order.OrderType,
	// 			BaseCurrencySymbol: baseCurrencySymbol,
	// 			BaseCurrencyName:   baseCurrencyName,
	// 			BaseQuantity:       r.Data.Order.BaseQuantity,
	// 			BasePercent:        r.Data.Order.BasePercent,
	// 			CurrencySymbol:     currencySymbol,
	// 			CurrencyName:       currencyName,
	// 			CurrencyQuantity:   r.Data.Order.CurrencyQuantity,
	// 			CurrencyPercent:    r.Data.Order.CurrencyPercent,
	// 			Status:             r.Data.Order.Status,
	// 			ChainStatus:        r.Data.Order.ChainStatus,
	// 			Conditions:         r.Data.Order.Conditions,
	// 			Condition:          r.Data.Order.Condition,
	// 			ParentOrderID:      r.Data.Order.ParentOrderID,
	// 			Price:              r.Data.Order.Price,
	// 		},
	// 	},
	// }

	// return c.JSON(http.StatusOK, res)
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
func (controller *TradeController) HandleListTrades(c echo.Context) error {
	return nil
	// token := c.Get("user").(*jwt.Token)
	// claims := token.Claims.(jwt.MapClaims)
	// userID := claims["jti"].(string)

	// getRequest := orders.GetUserOrdersRequest{
	// 	UserID: userID,
	// }

	// r, _ := controller.Orders.GetUserOrders(context.Background(), &getRequest)
	// if r.Status != response.Success {
	// 	res := &ResponseError{
	// 		Status:  r.Status,
	// 		Message: r.Message,
	// 	}

	// 	if r.Status == response.Fail {
	// 		return c.JSON(http.StatusBadRequest, res)
	// 	}
	// 	if r.Status == response.Error {
	// 		return c.JSON(http.StatusInternalServerError, res)
	// 	}
	// }

	// data := make([]*Order, len(r.Data.Orders))
	// for i, o := range r.Data.Orders {

	// 	names := strings.Split(o.MarketName, "-")
	// 	baseCurrencySymbol := names[1]
	// 	baseCurrencyName := controller.currencies[baseCurrencySymbol]
	// 	currencySymbol := names[0]
	// 	currencyName := controller.currencies[currencySymbol]

	// 	data[i] = &Order{
	// 		OrderID:            o.OrderID,
	// 		KeyID:              o.KeyID,
	// 		Exchange:           o.Exchange,
	// 		ExchangeOrderID:    o.ExchangeOrderID,
	// 		ExchangeMarketName: o.ExchangeMarketName,
	// 		MarketName:         o.MarketName,
	// 		Side:               o.Side,
	// 		OrderType:          o.OrderType,
	// 		BaseCurrencySymbol: baseCurrencySymbol,
	// 		BaseCurrencyName:   baseCurrencyName,
	// 		BaseQuantity:       o.BaseQuantity,
	// 		BasePercent:        o.BasePercent,
	// 		CurrencySymbol:     currencySymbol,
	// 		CurrencyName:       currencyName,
	// 		CurrencyQuantity:   o.CurrencyQuantity,
	// 		CurrencyPercent:    o.CurrencyPercent,
	// 		Status:             o.Status,
	// 		ChainStatus:        o.ChainStatus,
	// 		Conditions:         o.Conditions,
	// 		Condition:          o.Condition,
	// 		ParentOrderID:      o.ParentOrderID,
	// 		Price:              o.Price,
	// 	}
	// }

	// res := &ResponseOrdersSuccess{
	// 	Status: response.Success,
	// 	Data: &UserOrdersData{
	// 		Orders: data,
	// 	},
	// }

	// return c.JSON(http.StatusOK, res)
}

// func fail(c echo.Context, msg string) error {
// 	res := &ResponseError{
// 		Status:  response.Fail,
// 		Message: msg,
// 	}

// 	return c.JSON(http.StatusBadRequest, res)
// }

// swagger:route POST /strategies orders addOrder
//
// create a new strategy (protected)
//
// This will create a new chain of orders for the user.
// Example:
//{
//    "keyID": "680d6bbf-1feb-4122-bd10-0e7ce080676a",
//    "marketName": "ADA-BTC",
//    "baseBalance": 0.1,
//    "currencyBalance": 0.0,
//    "live": true,
//    "orders": [
//        {
//            "side": "buy",
//            "basePercent": 0.1,
//            "orderType": "market",
//            "conditions": "price <= 0.00002800"
//        },
//        {
//            "side": "sell",
//            "currencyPercent": 1.0,
//            "orderType": "market",
//            "conditions": "price <= 0.00002200"
//        }
//    ]
//}
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
func (controller *TradeController) HandlePostTrade(c echo.Context) error {
	//defer c.Request().Body.Close()
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	// read strategy from post body
	strategy := new(StrategyRequest)
	err := json.NewDecoder(c.Request().Body).Decode(&strategy)
	if err != nil {
		return fail(c, err.Error())
	}

	// market name and api key are required
	if strategy.MarketName == "" || strategy.KeyID == "" {
		return fail(c, "marketName and keyID required!")
	}
	if !strings.Contains(strategy.MarketName, "-") {
		return fail(c, "marketName must be currency-base: e.g. ADA-BTC")
	}

	// error check all orders
	for i, order := range strategy.Orders {
		log.Printf("order %d: %+v\n", i, order)

		if !orderValidator.ValidateOrderType(order.OrderType) {
			return fail(c, "market, limit, or virtual orders only!")
		}

		if !sideValidator.ValidateSide(order.Side) {
			return fail(c, "buy or sell required for side!")
		}
	}

	getRequest := keys.GetUserKeyRequest{
		KeyID:  strategy.KeyID,
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
	// key must be verified status
	if r.Data.Key.Status != key.Verified {
		return fail(c, "key must be verified")
	}

	//exchangeName = r.Data.Key.Exchange
	//request := orders.StrategyRequest{

	//		UserID:           userID,
	//		KeyID:            order.KeyID,
	//		Exchange:         exchangeName,
	//		MarketName:       order.MarketName,
	//      Live:             stratgey.Live,
	//      Orders:           &orders,
	//}

	//		Side:             order.Side,
	//		OrderType:        order.OrderType,
	//		Price:            order.Price,
	//		BaseQuantity:     order.BaseQuantity,
	//		BasePercent:      order.BasePercent,
	//		CurrencyQuantity: order.CurrencyQuantity,
	//		CurrencyPercent:  order.CurrencyPercent,
	//		Conditions:       order.Conditions,
	//	}
	//	requests = append(requests, &request)
	//}

	//orderRequests := orders.OrdersRequest{
	//	Orders: requests,
	//}

	//// add order returns nil for error
	//r, _ := controller.Orders.AddOrders(context.Background(), &orderRequests)
	//if r.Status != response.Success {
	//	res := &ResponseError{
	//		Status:  r.Status,
	//		Message: r.Message,
	//	}

	//	if r.Status == response.Fail {
	//		return c.JSON(http.StatusBadRequest, res)
	//	}
	//	if r.Status == response.Error {
	//		return c.JSON(http.StatusInternalServerError, res)
	//	}
	//}

	//data := make([]*Order, len(r.Data.Orders))
	//for i, o := range r.Data.Orders {
	//	names := strings.Split(o.MarketName, "-")
	//	baseCurrencySymbol := names[1]
	//	baseCurrencyName := controller.currencies[baseCurrencySymbol]
	//	currencySymbol := names[0]
	//	currencyName := controller.currencies[currencySymbol]

	//	data[i] = &Order{
	//		OrderID:            o.OrderID,
	//		KeyID:              o.KeyID,
	//		Exchange:           o.Exchange,
	//		ExchangeOrderID:    o.ExchangeOrderID,
	//		ExchangeMarketName: o.ExchangeMarketName,
	//		MarketName:         o.MarketName,
	//		Side:               o.Side,
	//		OrderType:          o.OrderType,
	//		Price:              o.Price,
	//		BaseCurrencySymbol: baseCurrencySymbol,
	//		BaseCurrencyName:   baseCurrencyName,
	//		BaseQuantity:       o.BaseQuantity,
	//		BasePercent:        o.BasePercent,
	//		CurrencySymbol:     currencySymbol,
	//		CurrencyName:       currencyName,
	//		CurrencyQuantity:   o.CurrencyQuantity,
	//		CurrencyPercent:    o.CurrencyPercent,
	//		Status:             o.Status,
	//		Conditions:         o.Conditions,
	//		Condition:          o.Condition,
	//		ParentOrderID:      o.ParentOrderID,
	//		ChainStatus:        o.ChainStatus,
	//	}
	//}

	//res := &ResponseOrdersSuccess{
	//	Status: response.Success,
	//	Data: &UserOrdersData{
	//		Orders: data,
	//	},
	//}

	return c.JSON(http.StatusOK, "ok")
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
// func (controller *TradeController) HandleUpdateOrder(c echo.Context) error {
// 	token := c.Get("user").(*jwt.Token)
// 	claims := token.Claims.(jwt.MapClaims)
// 	userID := claims["jti"].(string)
// 	orderID := c.Param("orderID")

// 	orderRequest := UpdateOrderRequest{}

// 	err := json.NewDecoder(c.Request().Body).Decode(&orderRequest)
// 	if err != nil {
// 		res := &ResponseError{
// 			Status:  response.Fail,
// 			Message: err.Error(),
// 		}

// 		return c.JSON(http.StatusBadRequest, res)
// 	}

// 	// client can only update description
// 	updateRequest := orders.OrderRequest{
// 		OrderID:      orderID,
// 		UserID:       userID,
// 		Conditions:   orderRequest.Conditions,
// 		BaseQuantity: orderRequest.BaseQuantity,
// 	}

// 	r, _ := controller.Orders.UpdateOrder(context.Background(), &updateRequest)
// 	if r.Status != response.Success {
// 		res := &ResponseError{
// 			Status:  r.Status,
// 			Message: r.Message,
// 		}

// 		if r.Status == response.Fail {
// 			return c.JSON(http.StatusBadRequest, res)
// 		}
// 		if r.Status == response.Error {
// 			return c.JSON(http.StatusInternalServerError, res)
// 		}
// 	}

// 	names := strings.Split(r.Data.Order.MarketName, "-")
// 	baseCurrencySymbol := names[1]
// 	baseCurrencyName := controller.currencies[baseCurrencySymbol]
// 	currencySymbol := names[0]
// 	currencyName := controller.currencies[currencySymbol]

// 	res := &ResponseOrderSuccess{
// 		Status: response.Success,
// 		Data: &UserOrderData{
// 			Order: &Order{
// 				OrderID:            r.Data.Order.OrderID,
// 				KeyID:              r.Data.Order.KeyID,
// 				Exchange:           r.Data.Order.Exchange,
// 				ExchangeOrderID:    r.Data.Order.ExchangeOrderID,
// 				ExchangeMarketName: r.Data.Order.ExchangeMarketName,
// 				MarketName:         r.Data.Order.MarketName,
// 				Side:               r.Data.Order.Side,
// 				Price:              r.Data.Order.Price,
// 				OrderType:          r.Data.Order.OrderType,
// 				BaseCurrencySymbol: baseCurrencySymbol,
// 				BaseCurrencyName:   baseCurrencyName,
// 				BaseQuantity:       r.Data.Order.BaseQuantity,
// 				BasePercent:        r.Data.Order.BasePercent,
// 				CurrencySymbol:     currencySymbol,
// 				CurrencyName:       currencyName,
// 				CurrencyQuantity:   r.Data.Order.CurrencyQuantity,
// 				CurrencyPercent:    r.Data.Order.CurrencyPercent,
// 				Status:             r.Data.Order.Status,
// 				ChainStatus:        r.Data.Order.ChainStatus,
// 				Conditions:         r.Data.Order.Conditions,
// 				Condition:          r.Data.Order.Condition,
// 				ParentOrderID:      r.Data.Order.ParentOrderID,
// 			},
// 		},
// 	}

// 	return c.JSON(http.StatusOK, res)
// }

// swagger:route DELETE /orders/:orderID orders deleteOrder
//
// Remove and order (protected)
//
// Cannot remove orders that have already executed.
//
// responses:
//  200: responseOrderSuccess data will be null with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
// func (controller *TradeController) HandleDeleteOrder(c echo.Context) error {
// 	token := c.Get("user").(*jwt.Token)
// 	claims := token.Claims.(jwt.MapClaims)
// 	userID := claims["jti"].(string)
// 	orderID := c.Param("orderID")

// 	removeRequest := orders.RemoveOrderRequest{
// 		OrderID: orderID,
// 		UserID:  userID,
// 	}

// 	r, _ := controller.Orders.RemoveOrder(context.Background(), &removeRequest)
// 	if r.Status != response.Success {
// 		res := &ResponseError{
// 			Status:  r.Status,
// 			Message: r.Message,
// 		}

// 		if r.Status == response.Fail {
// 			return c.JSON(http.StatusBadRequest, res)
// 		}
// 		if r.Status == response.Error {
// 			return c.JSON(http.StatusInternalServerError, res)
// 		}
// 	}

// 	res := &ResponseOrderSuccess{
// 		Status: response.Success,
// 	}

// 	return c.JSON(http.StatusOK, res)
// }
