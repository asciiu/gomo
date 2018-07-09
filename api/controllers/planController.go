package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	asql "github.com/asciiu/gomo/api/db/sql"
	orderValidator "github.com/asciiu/gomo/common/constants/order"
	"github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/response"
	sideValidator "github.com/asciiu/gomo/common/constants/side"
	keys "github.com/asciiu/gomo/key-service/proto/key"
	plans "github.com/asciiu/gomo/plan-service/proto/plan"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
	"golang.org/x/net/context"
)

type PlanController struct {
	DB    *sql.DB
	Plans plans.PlanServiceClient
	Keys  keys.KeyServiceClient
	// map of ticker symbol to full name
	currencies map[string]string
}

// A ResponsePlansSuccess will always contain a status of "successful".
// swagger:model responsePlansSuccess
type ResponsePlansSuccess struct {
	Status string     `json:"status"`
	Data   *PlansPage `json:"data"`
}

type PlansPage struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Total    uint32  `json:"total"`
	Plans    []*Plan `json:"plans"`
}

type UserPlanData struct {
	Plan *Plan `json:"plan"`
}

type UserPlansData struct {
	PLans []*Plan `json:"plans"`
}

// This is the response struct for order
type Plan struct {
	PlanID             string         `json:"planID"`
	KeyID              string         `json:"keyID"`
	Exchange           string         `json:"exchange"`
	ExchangeMarketName string         `json:"exchangeMarketName"`
	MarketName         string         `json:"marketName"`
	BaseCurrencySymbol string         `json:"baseCurrencySymbol"`
	BaseCurrencyName   string         `json:"baseCurrencyName"`
	BaseBalance        float64        `json:"baseBalance"`
	CurrencySymbol     string         `json:"currencySymbol"`
	CurrencyName       string         `json:"currencyName"`
	CurrencyBalance    float64        `json:"currencyBalance"`
	Status             string         `json:"status"`
	Orders             []*plans.Order `json:"orders"`
}

// swagger:parameters addPlan
type PlanRequest struct {
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
type UpdatePlanRequest struct {
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

// A ResponsePlanSuccess will always contain a status of "successful".
// swagger:model ResponsePlanSuccess
type ResponsePlanSuccess struct {
	Status string        `json:"status"`
	Data   *UserPlanData `json:"data"`
}

// A ResponseKeysSuccess will always contain a status of "successful".
// swagger:model responseOrdersSuccess
type ResponseStategiesSuccess struct {
	Status string          `json:"status"`
	Data   *UserOrdersData `json:"data"`
}

func NewPlanController(db *sql.DB) *PlanController {
	// Create a new service. Optionally include some options here.
	service := k8s.NewService(micro.Name("apikey.client"))
	service.Init()

	controller := PlanController{
		DB:         db,
		Plans:      plans.NewPlanServiceClient("plans", service.Client()),
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
func (controller *PlanController) HandleGetPlan(c echo.Context) error {
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

// swagger:route GET /plans plans getUserPlans
//
// get user plans (protected)
//
// Returns a summary of user plans.
// Query Params: status, marketName, exchange, page, pageSize
//
// The defaults for the params are:
// status - active
// page - 0
// pageSize - 50
//
// example: /plans?status=failed&exchange=binance
//
// responses:
//  200: responsePlansSuccess "data" will contain an array of plan summaries
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *PlanController) HandleListPlans(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	exchange := c.QueryParam("exchange")
	marketName := c.QueryParam("marketName")
	status := c.QueryParam("status")

	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("pageSize")

	// defaults for page and page size here
	// ignore the errors and assume the values are int
	page, _ := strconv.ParseUint(pageStr, 10, 32)
	pageSize, _ := strconv.ParseUint(pageSizeStr, 10, 32)
	if pageSize == 0 {
		pageSize = 50
	}

	// default status should be active plans
	if status == "" {
		status = plan.Active
	}

	getRequest := plans.GetUserPlansRequest{
		UserID:     userID,
		Page:       uint32(page),
		PageSize:   uint32(pageSize),
		Exchange:   exchange,
		MarketName: marketName,
		Status:     status,
	}

	r, _ := controller.Plans.GetUserPlans(context.Background(), &getRequest)
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

	plans := make([]*Plan, 0)
	for _, plan := range r.Data.Plans {
		names := strings.Split(plan.MarketName, "-")
		baseCurrencySymbol := names[1]
		baseCurrencyName := controller.currencies[baseCurrencySymbol]
		currencySymbol := names[0]
		currencyName := controller.currencies[currencySymbol]

		pln := Plan{
			PlanID:             plan.PlanID,
			KeyID:              plan.KeyID,
			Exchange:           plan.Exchange,
			ExchangeMarketName: plan.ExchangeMarketName,
			MarketName:         plan.MarketName,
			BaseCurrencySymbol: baseCurrencySymbol,
			BaseCurrencyName:   baseCurrencyName,
			BaseBalance:        plan.BaseBalance,
			CurrencySymbol:     currencySymbol,
			CurrencyName:       currencyName,
			CurrencyBalance:    plan.CurrencyBalance,
			Status:             plan.Status,
		}
		plans = append(plans, &pln)
	}

	res := &ResponsePlansSuccess{
		Status: response.Success,
		Data: &PlansPage{
			Page:     r.Data.Page,
			PageSize: r.Data.PageSize,
			Total:    r.Data.Total,
			Plans:    plans,
		},
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:route POST /plans plans addPlan
//
// create a new planned strategy (protected)
//
// This will create a new chain of orders for the user. All orders are encapsulated within a plan.
// Example:
//{
//    "keyID": "680d6bbf-1feb-4122-bd10-0e7ce080676a",
//    "marketName": "ADA-BTC",
//    "baseBalance": 1.0,
//    "currencyBalance": 0.0,
//    "live": true,
//    "orders": [
//        {
//            "side": "buy",
//            "basePercent": 0.5,
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
// The order chain starts out with a reserve base balance of 1.0 BTC. This order chain will buy 0.5 BTC at
// market price when the trigger price is less than or equal to 2800 satoshi. The following sell order will sell 100% of the order
// strategy's (i.e. chain of orders) cardano balance. Cardano balance should in theory be dictated by the Cardano that this
// chain bought - 100% does not mean 100% of user's cardano balance. The conditions will likely be a json string of conditionLabel: condition.
// example: "stopLoss: price <= 0.00002200, takeProfit: ...."
//
// responses:
//  200: ResponsePlanSuccess "data" will contain list of orders with "status": "success"
//  400: responseError missing or incorrect params with "status": "fail"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *PlanController) HandlePostPlan(c echo.Context) error {
	//defer c.Request().Body.Close()
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	// read strategy from post body
	newPlan := new(PlanRequest)
	err := json.NewDecoder(c.Request().Body).Decode(&newPlan)
	if err != nil {
		return fail(c, err.Error())
	}

	// market name and api key are required
	if newPlan.MarketName == "" || newPlan.KeyID == "" {
		return fail(c, "marketName and keyID required!")
	}
	if !strings.Contains(newPlan.MarketName, "-") {
		return fail(c, "marketName must be currency-base: e.g. ADA-BTC")
	}
	if len(newPlan.Orders) == 0 {
		return fail(c, "at least one order required for a trade plan")
	}

	// error check all orders
	orders := make([]*plans.OrderRequest, 0)
	for i, order := range newPlan.Orders {
		log.Printf("order %d: %+v\n", i, order)

		if !orderValidator.ValidateOrderType(order.OrderType) {
			return fail(c, "market, limit, or virtual orders only!")
		}

		if !sideValidator.ValidateSide(order.Side) {
			return fail(c, "buy or sell required for side!")
		}
		or := plans.OrderRequest{
			Side:            order.Side,
			OrderType:       order.OrderType,
			BasePercent:     order.BasePercent,
			CurrencyPercent: order.CurrencyPercent,
			Conditions:      order.Conditions,
			Price:           order.Price,
		}
		orders = append(orders, &or)
	}

	newPlanRequest := plans.PlanRequest{
		UserID:          userID,
		KeyID:           newPlan.KeyID,
		MarketName:      newPlan.MarketName,
		BaseBalance:     newPlan.BaseBalance,
		CurrencyBalance: newPlan.CurrencyBalance,
		Active:          newPlan.Live,
		Orders:          orders,
	}

	// add plan returns nil for error
	r, _ := controller.Plans.AddPlan(context.Background(), &newPlanRequest)
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

	names := strings.Split(r.Data.Plan.MarketName, "-")
	baseCurrencySymbol := names[1]
	baseCurrencyName := controller.currencies[baseCurrencySymbol]
	currencySymbol := names[0]
	currencyName := controller.currencies[currencySymbol]

	data := Plan{
		PlanID:             r.Data.Plan.PlanID,
		KeyID:              r.Data.Plan.KeyID,
		Exchange:           r.Data.Plan.Exchange,
		ExchangeMarketName: r.Data.Plan.ExchangeMarketName,
		MarketName:         r.Data.Plan.MarketName,
		BaseCurrencySymbol: baseCurrencySymbol,
		BaseCurrencyName:   baseCurrencyName,
		BaseBalance:        r.Data.Plan.BaseBalance,
		CurrencySymbol:     currencySymbol,
		CurrencyName:       currencyName,
		CurrencyBalance:    r.Data.Plan.CurrencyBalance,
		Status:             r.Data.Plan.Status,
		Orders:             r.Data.Plan.Orders,
	}

	res := &ResponsePlanSuccess{
		Status: response.Success,
		Data: &UserPlanData{
			Plan: &data,
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
