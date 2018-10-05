package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	protoActivity "github.com/asciiu/gomo/activity-bulletin/proto"
	asql "github.com/asciiu/gomo/api/db/sql"
	constRes "github.com/asciiu/gomo/common/constants/response"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type PlanController struct {
	DB             *sql.DB
	PlanClient     protoPlan.PlanServiceClient
	BulletinClient protoActivity.ActivityBulletinClient
	// map of ticker symbol to full name
	currencies map[string]string
}

// A ResponsePlansSuccess will always contain a status of "successful".
// swagger:model ResponsePlansSuccess
type ResponsePlansSuccess struct {
	Status string     `json:"status"`
	Data   *PlansPage `json:"data"`
}

// A ResponsePlanSuccess will always contain a status of "successful".
// swagger:model ResponsePlanSuccess
type ResponsePlanSuccess struct {
	Status string `json:"status"`
	Data   *Plan  `json:"data"`
}

type PlansPage struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Total    uint32  `json:"total"`
	Plans    []*Plan `json:"plans"`
}

// This response should never return the key secret
type Plan struct {
	PlanID                     string   `json:"planID"`
	PlanTemplateID             string   `json:"planTemplateID"`
	PlanNumber                 uint64   `json:"planNumber"`
	Title                      string   `json:"title"`
	TotalDepth                 uint32   `json:"totalDepth"`
	Exchange                   string   `json:"exchange"`
	UserCurrencySymbol         string   `json:"userCurrencySymbol"`
	UserCurrencyBalanceAtInit  float64  `json:"userCurrencyBalanceAtInit"`
	InitialUserCurrencyBalance float64  `json:"initialUserCurrencyBalance"`
	InitialTimestamp           string   `json:"initialTimestamp"`
	ActiveCurrencySymbol       string   `json:"activeCurrencySymbol"`
	ActiveCurrencyName         string   `json:"activeCurrencyName"`
	ActiveCurrencyBalance      float64  `json:"activeCurrencyBalance"`
	InitialCurrencySymbol      string   `json:"initialCurrencySymbol"`
	InitialCurrencyName        string   `json:"initialCurrencyName"`
	InitialCurrencyBalance     float64  `json:"initialCurrencyBalance"`
	LastExecutedOrderID        string   `json:"lastExecutedOrderID"`
	LastExecutedPlanDepth      uint32   `json:"lastExecutedPlanDepth"`
	Status                     string   `json:"status"`
	CloseOnComplete            bool     `json:"closeOnComplete"`
	CreatedOn                  string   `json:"createdOn"`
	UpdatedOn                  string   `json:"updatedOn"`
	Orders                     []*Order `json:"orders"`
}

type PlanActivitySummary struct {
	Total  uint32                  `json:"total"`
	Recent *protoActivity.Activity `json:"recent"`
}

type Order struct {
	OrderID                  string                `json:"orderID,omitempty"`
	ParentOrderID            string                `json:"parentOrderID,omitempty"`
	PlanDepth                uint32                `json:"planDepth,omitempty"`
	OrderTemplateID          string                `json:"orderTemplateID,omitempty"`
	AccountID                string                `json:"accountID,omitempty"`
	KeyPublic                string                `json:"keyPublic,omitempty"`
	KeyDescription           string                `json:"keyDescription,omitempty"`
	OrderPriority            uint32                `json:"orderPriority,omitempty"`
	OrderType                string                `json:"orderType,omitempty"`
	Side                     string                `json:"side,omitempty"`
	LimitPrice               float64               `json:"limitPrice,omitempty"`
	Exchange                 string                `json:"exchange,omitempty"`
	ExchangeOrderID          string                `json:"exchangeOrderID"`
	ExchangePrice            float64               `json:"exchangePrice"`
	ExchangeTime             string                `json:"exchangeTime"`
	MarketName               string                `json:"marketName,omitempty"`
	BaseCurrencySymbol       string                `json:"baseCurrencySymbol"`
	BaseCurrencyName         string                `json:"baseCurrencyName"`
	MarketCurrencySymbol     string                `json:"marketCurrencySymbol"`
	MarketCurrencyName       string                `json:"marketCurrencyName"`
	InitialCurrencySymbol    string                `json:"initialCurrencySymbol"`
	InitialCurrencyName      string                `json:"initialCurrencyName"`
	InitialCurrencyBalance   float64               `json:"initialCurrencyBalance"`
	InitialCurrencyValue     float64               `json:"initialCurrencyValue"`
	InitialCurrencyTraded    float64               `json:"initialCurrencyTraded"`
	InitialCurrencyRemainder float64               `json:"initialCurrencyRemainder"`
	FinalCurrencySymbol      string                `json:"finalCurrencySymbol"`
	FinalCurrencyName        string                `json:"finalCurrencyName"`
	FinalCurrencyBalance     float64               `json:"finalCurrencyBalance"`
	FinalCurrencyValue       float64               `json:"finalCurrencyValue"`
	FeeCurrencySymbol        string                `json:"feeCurrencySymbol"`
	FeeCurrencyAmount        float64               `json:"feeCurrencyAmount"`
	Grupo                    string                `json:"grupo"`
	Status                   string                `json:"status,omitempty"`
	Errors                   string                `json:"errors"`
	CreatedOn                string                `json:"createdOn,omitempty"`
	UpdatedOn                string                `json:"updatedOn,omitempty"`
	Triggers                 []*protoOrder.Trigger `json:"triggers,omitempty"`
}

func fail(c echo.Context, msg string) error {
	res := &ResponseError{
		Status:  constRes.Fail,
		Message: msg,
	}

	return c.JSON(http.StatusBadRequest, res)
}

func NewPlanController(db *sql.DB, service micro.Service) *PlanController {
	controller := PlanController{
		DB:             db,
		PlanClient:     protoPlan.NewPlanServiceClient("plans", service.Client()),
		BulletinClient: protoActivity.NewActivityBulletinClient("bulletin", service.Client()),
		currencies:     make(map[string]string),
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

// swagger:route DELETE /plans/:planID plans DeletePlan
//
// deletes a plan (protected)
//
// You may delete a plan if it has not executed. That is, the plan has no filled protoOrder. Delete plan becomes cancel plan
// when an order has been filled. In theory, this should kill all active protoOrder for a plan and set the status for the plan
// as 'aborted'. Once a plan has been 'aborted' you cannot update or restart that plan.
//
// responses:
//  200: ResponsePlanSuccess "data" will contain plan summary.
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *PlanController) HandleDeletePlan(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	planID := c.Param("planID")

	delRequest := protoPlan.DeletePlanRequest{
		PlanID: planID,
		UserID: userID,
	}

	r, _ := controller.PlanClient.DeletePlan(context.Background(), &delRequest)
	if r.Status != constRes.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		switch {
		case r.Status == constRes.Nonentity:
			return c.JSON(http.StatusNotFound, res)
		case r.Status == constRes.Fail:
			return c.JSON(http.StatusBadRequest, res)
		default:
			return c.JSON(http.StatusInternalServerError, res)
		}
	}

	// res := &ResponsePlanWithOrderPageSuccess{
	// 	Status: constRes.Success,
	// 	Data: &PlanWithOrderPage{
	// 		PlanID:         r.Data.Plan.PlanID,
	// 		PlanTemplateID: r.Data.Plan.PlanTemplateID,
	// 		Exchange:       r.Data.Plan.Exchange,
	// 		Status:         r.Data.Plan.Status,
	// 		CreatedOn:      r.Data.Plan.CreatedOn,
	// 		UpdatedOn:      r.Data.Plan.UpdatedOn,
	// 	},
	// }

	return c.JSON(http.StatusOK, "")
}

// required for swaggered, otherwise never used
// swagger:parameters GetPlanParams
type GetPlanParams struct {
	PlanDepth  uint32 `json:"planDepth"`
	PlanLength uint32 `json:"planLength"`
}

// swagger:route GET /plans/:planID plans GetPlanParams
//
// get plan with planID (protected)
//
// Returns a plan with the currency and base currency names. Plan protoOrder will be retrieved
// based upon planDepth and planLength. The root order of a plan begins at planDepth=0. The
// planLength deteremines the length of the tree to retrieve.
//
// example: /plans/:some_plan_id?planDepth=0&planLength=10
// The example above will retrieve the plan from planDepth 0 to planDepth 10.
//
// responses:
//  200: ResponsePlanSuccess "data" will contain plan deets.
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *PlanController) HandleGetPlan(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	planID := c.Param("planID")

	planDepth := c.QueryParam("planDepth")
	planLength := c.QueryParam("planLength")

	// defaults for plan depth and plan length
	// ignore the errors and assume the values are int
	pd, _ := strconv.ParseUint(planDepth, 10, 32)
	pl, _ := strconv.ParseUint(planLength, 10, 32)
	if pl == 0 {
		pl = 10
	}

	getRequest := protoPlan.GetUserPlanRequest{
		PlanID:     planID,
		UserID:     userID,
		PlanDepth:  uint32(pd),
		PlanLength: uint32(pl)}

	r, _ := controller.PlanClient.GetUserPlan(context.Background(), &getRequest)
	if r.Status != constRes.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		switch {
		case r.Status == constRes.Nonentity:
			return c.JSON(http.StatusNotFound, res)
		case r.Status == constRes.Fail:
			return c.JSON(http.StatusBadRequest, res)
		default:
			return c.JSON(http.StatusInternalServerError, res)
		}
	}
	plan := r.Data.Plan

	// convert orders in data plan to api model orders
	newOrders := make([]*Order, 0)
	for _, o := range plan.Orders {
		names := strings.Split(o.MarketName, "-")
		baseCurrencySymbol := names[1]
		baseCurrencyName := controller.currencies[baseCurrencySymbol]
		marketCurrencySymbol := names[0]
		marketCurrencyName := controller.currencies[marketCurrencySymbol]
		newo := Order{
			OrderID:                  o.OrderID,
			ParentOrderID:            o.ParentOrderID,
			PlanDepth:                o.PlanDepth,
			OrderTemplateID:          o.OrderTemplateID,
			AccountID:                o.AccountID,
			KeyPublic:                o.KeyPublic,
			KeyDescription:           o.KeyDescription,
			OrderPriority:            o.OrderPriority,
			OrderType:                o.OrderType,
			Side:                     o.Side,
			LimitPrice:               o.LimitPrice,
			Exchange:                 o.Exchange,
			ExchangeOrderID:          o.ExchangeOrderID,
			ExchangePrice:            o.ExchangePrice,
			ExchangeTime:             o.ExchangeTime,
			MarketName:               o.MarketName,
			BaseCurrencySymbol:       baseCurrencySymbol,
			BaseCurrencyName:         baseCurrencyName,
			MarketCurrencySymbol:     marketCurrencySymbol,
			MarketCurrencyName:       marketCurrencyName,
			InitialCurrencySymbol:    o.InitialCurrencySymbol,
			InitialCurrencyName:      controller.currencies[o.InitialCurrencySymbol],
			InitialCurrencyBalance:   o.InitialCurrencyBalance,
			InitialCurrencyTraded:    o.InitialCurrencyTraded,
			InitialCurrencyRemainder: o.InitialCurrencyRemainder,
			FinalCurrencySymbol:      o.FinalCurrencySymbol,
			FinalCurrencyName:        controller.currencies[o.FinalCurrencySymbol],
			FinalCurrencyBalance:     o.FinalCurrencyBalance,
			FeeCurrencySymbol:        o.FeeCurrencySymbol,
			FeeCurrencyAmount:        o.FeeCurrencyAmount,
			Grupo:                    o.Grupo,
			Status:                   o.Status,
			Errors:                   o.Errors,
			CreatedOn:                o.CreatedOn,
			UpdatedOn:                o.UpdatedOn,
			Triggers:                 o.Triggers,
		}
		newOrders = append(newOrders, &newo)
	}

	res := &ResponsePlanSuccess{
		Status: constRes.Success,
		Data: &Plan{
			PlanID:                    plan.PlanID,
			PlanTemplateID:            plan.PlanTemplateID,
			PlanNumber:                plan.UserPlanNumber,
			Title:                     plan.Title,
			TotalDepth:                plan.TotalDepth,
			Exchange:                  plan.Exchange,
			InitialTimestamp:          plan.InitialTimestamp,
			UserCurrencySymbol:        plan.UserCurrencySymbol,
			UserCurrencyBalanceAtInit: plan.UserCurrencyBalanceAtInit,
			ActiveCurrencySymbol:      plan.ActiveCurrencySymbol,
			ActiveCurrencyName:        controller.currencies[plan.ActiveCurrencySymbol],
			ActiveCurrencyBalance:     plan.ActiveCurrencyBalance,
			InitialCurrencySymbol:     plan.InitialCurrencySymbol,
			InitialCurrencyName:       controller.currencies[plan.InitialCurrencySymbol],
			InitialCurrencyBalance:    plan.InitialCurrencyBalance,
			Status:                    plan.Status,
			CloseOnComplete:           plan.CloseOnComplete,
			LastExecutedOrderID:       plan.LastExecutedOrderID,
			LastExecutedPlanDepth:     plan.LastExecutedPlanDepth,
			Orders:                    newOrders,
			CreatedOn:                 plan.CreatedOn,
			UpdatedOn:                 plan.UpdatedOn,
		},
	}

	return c.JSON(http.StatusOK, res)
}

// required for swaggered, otherwise never used
// swagger:parameters GetUserPlansParams
type GetUserPlansParams struct {
	Exchange   string `json:"exchange"`
	MarketName string `json:"marketName"`
	Status     string `json:"status"`
	Page       uint32 `json:"page"`
	PageSize   uint32 `json:"pageSize"`
}

// swagger:route GET /plans plans GetUserPlansParams
//
// get user plans (protected)
//
// Returns a summary for each plan. The plan orders will contain the last executed order and the child orders of the executed order.
// Query Params: status, marketName, exchange, page, pageSize
//
// The defaults for the params are:
// status - active
// page - 0
// pageSize - 50
//
// example: /plans?exchange=binance
//
// responses:
//  200: ResponsePlansSuccess "data" will contain an array of plan summaries
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

	getRequest := protoPlan.GetUserPlansRequest{
		UserID:     userID,
		Page:       uint32(page),
		PageSize:   uint32(pageSize),
		Exchange:   exchange,
		MarketName: marketName,
		Status:     status,
	}

	r, _ := controller.PlanClient.GetUserPlans(context.Background(), &getRequest)
	if r.Status != constRes.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == constRes.Fail {
			return c.JSON(http.StatusBadRequest, res)
		}
		if r.Status == constRes.Error {
			return c.JSON(http.StatusInternalServerError, res)
		}
	}

	plans := make([]*Plan, 0)
	for _, plan := range r.Data.Plans {

		cOrders := make([]*Order, 0)
		for _, o := range plan.Orders {
			names := strings.Split(o.MarketName, "-")
			baseCurrencySymbol := names[1]
			baseCurrencyName := controller.currencies[baseCurrencySymbol]
			marketCurrencySymbol := names[0]
			marketCurrencyName := controller.currencies[marketCurrencySymbol]
			newo := Order{
				OrderID:                  o.OrderID,
				ParentOrderID:            o.ParentOrderID,
				PlanDepth:                o.PlanDepth,
				OrderTemplateID:          o.OrderTemplateID,
				AccountID:                o.AccountID,
				KeyPublic:                o.KeyPublic,
				KeyDescription:           o.KeyDescription,
				OrderPriority:            o.OrderPriority,
				OrderType:                o.OrderType,
				Side:                     o.Side,
				LimitPrice:               o.LimitPrice,
				Exchange:                 o.Exchange,
				ExchangeOrderID:          o.ExchangeOrderID,
				ExchangePrice:            o.ExchangePrice,
				ExchangeTime:             o.ExchangeTime,
				MarketName:               o.MarketName,
				BaseCurrencySymbol:       baseCurrencySymbol,
				BaseCurrencyName:         baseCurrencyName,
				MarketCurrencySymbol:     marketCurrencySymbol,
				MarketCurrencyName:       marketCurrencyName,
				InitialCurrencySymbol:    o.InitialCurrencySymbol,
				InitialCurrencyName:      controller.currencies[o.InitialCurrencySymbol],
				InitialCurrencyBalance:   o.InitialCurrencyBalance,
				InitialCurrencyTraded:    o.InitialCurrencyTraded,
				InitialCurrencyRemainder: o.InitialCurrencyRemainder,
				FinalCurrencySymbol:      o.FinalCurrencySymbol,
				FinalCurrencyName:        controller.currencies[o.FinalCurrencySymbol],
				FinalCurrencyBalance:     o.FinalCurrencyBalance,
				FeeCurrencySymbol:        o.FeeCurrencySymbol,
				FeeCurrencyAmount:        o.FeeCurrencyAmount,
				Grupo:                    o.Grupo,
				Status:                   o.Status,
				Errors:                   o.Errors,
				CreatedOn:                o.CreatedOn,
				UpdatedOn:                o.UpdatedOn,
				Triggers:                 o.Triggers,
			}
			cOrders = append(cOrders, &newo)
		}

		pln := Plan{
			PlanID:                    plan.PlanID,
			PlanTemplateID:            plan.PlanTemplateID,
			PlanNumber:                plan.UserPlanNumber,
			Title:                     plan.Title,
			TotalDepth:                plan.TotalDepth,
			Exchange:                  plan.Exchange,
			UserCurrencySymbol:        plan.UserCurrencySymbol,
			UserCurrencyBalanceAtInit: plan.UserCurrencyBalanceAtInit,
			InitialTimestamp:          plan.InitialTimestamp,
			ActiveCurrencySymbol:      plan.ActiveCurrencySymbol,
			ActiveCurrencyName:        controller.currencies[plan.ActiveCurrencySymbol],
			ActiveCurrencyBalance:     plan.ActiveCurrencyBalance,
			InitialCurrencySymbol:     plan.InitialCurrencySymbol,
			InitialCurrencyName:       controller.currencies[plan.InitialCurrencySymbol],
			InitialCurrencyBalance:    plan.InitialCurrencyBalance,
			Status:                    plan.Status,
			CloseOnComplete:           plan.CloseOnComplete,
			LastExecutedOrderID:       plan.LastExecutedOrderID,
			LastExecutedPlanDepth:     plan.LastExecutedPlanDepth,
			CreatedOn:                 plan.CreatedOn,
			UpdatedOn:                 plan.UpdatedOn,
			Orders:                    cOrders,
		}

		plans = append(plans, &pln)
	}

	res := &ResponsePlansSuccess{
		Status: constRes.Success,
		Data: &PlansPage{
			Page:     r.Data.Page,
			PageSize: r.Data.PageSize,
			Total:    r.Data.Total,
			Plans:    plans,
		},
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:parameters PostPlan
type PlanRequest struct {
	// Optional base currency from which plan currency will be measured with. e.g. USDT, BTC, ETH. Default to USDT.
	// in: body
	UserCurrencySymbol string `json:"userCurrencySymbol"`
	// Required plan title
	// in: body
	Title string `json:"title"`
	// Optional plan template ID. Leo wanted this for the templating system.
	// in: body
	PlanTemplateID string `json:"planTemplateID"`
	// Optional init timestamp for plan RFC3339 formatted (e.g. 2018-08-26T22:49:10.168652Z). This timestamp will be used to measure initial user currency balance (valuation in user preferred currency)
	// in: body
	InitialTimestamp string `json:"initialTimestamp"`
	// Optional defaults to 'active' status. Valid input status is 'active', 'inactive', or 'historic'
	// in: body
	Status string `json:"status"`
	// Required bool to indicate that you want the plan to be 'closed' when the last order for the plan finishes (note: order status fail will also close the plan)
	// in: body
	CloseOnComplete bool `json:"closeOnComplete"`
	// Required array of orders. The structure of the order tree will be dictated by the parentOrderID. All orders following the root order must have a parentOrderID. The root order must have a parentOrderID of "00000000-0000-0000-0000-000000000000". Use grupo (aka spanish for group) to assign a group label to the order.
	// in: body
	Orders []*NewOrderReq `json:"orders"`
}

type NewOrderReq struct {
	// Required the client assigns the order ID as a UUID, the format is 8-4-4-4-12.
	// in: body
	OrderID string `json:"orderID"`
	// Optional precedence of order when multiple orders are at the same depth: value of 1 is highest priority. E.g. depth 2 buy ADA (1) or buy EOS (2). ADA with higher priority 1 will execute and EOS will not execute.
	// in: body
	OrderPriority uint32 `json:"orderPriority"`
	// Required order types are "market", "limit", "paper". Orders not within these types will be rejected.
	// in: body
	OrderType string `json:"orderType"`
	// Optional order template ID. This is a Leo thing.
	// in: body
	OrderTemplateID string `json:"orderTemplateID"`
	// Deprecated this used to be our key ID (string uuid) assigned to the user's exchange key and secret. Use accountID instead. DO NOT USE THIS!
	// in: body
	KeyID string `json:"keyID"`
	// Required accountID to use for this order. The account defines the exchange keys and balances.
	// in: body
	AccountID string `json:"accountID"`
	// Required the root node of the decision tree should be assigned a parentOrderID of "00000000-0000-0000-0000-000000000000" .
	// in: body
	ParentOrderID string `json:"parentOrderID"`
	// Required e.g. ADA-BTC. Base pair should be the suffix.
	// in: body
	MarketName string `json:"marketName"`
	// Required "buy" or "sell"
	// in: body
	Side string `json:"side"`
	// Required for 'limit' protoOrder. Defines limit price.
	// in: body
	LimitPrice float64 `json:"limitPrice"`
	// Required for the root order of the tree. Child protoOrder for tree may or may not have a currencyBalance.
	// in: body
	InitialCurrencyBalance float64 `json:"initialCurrencyBalance"`
	// Required these are the conditions that trigger the order to execute: ???
	// in: body
	Triggers []*TriggerReq `json:"triggers"`
	// Optional group of order
	// in: body
	Grupo string `json:"grupo"`
}

type TriggerReq struct {
	// Optional trigger template ID.
	// in: body
	TriggerTemplateID string   `json:"triggerTemplateID"`
	Name              string   `json:"name"`
	Title             string   `json:"title"`
	Code              string   `json:"code"`
	Actions           []string `json:"actions"`
	Index             uint32   `json:"index"`
}

// swagger:route POST /plans plans PostPlan
//
// create a new plan (protected)
//
// This will create a new chain of orders for the user.
//
// responses:
//  200: ResponsePlanSuccess "data" will contain the order tree
//  400: responseError missing or incorrect params with "status": "fail"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *PlanController) HandlePostPlan(c echo.Context) error {
	//defer c.Request().Body.Close()
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	// read strategy from post body
	newPlan := new(PlanRequest)
	err := c.Bind(&newPlan)
	//err := json.NewDecoder(c.Request().Body).Decode(&newPlan)
	if err != nil {
		return fail(c, err.Error())
	}

	// assemble the order requests
	newOrderRequests := make([]*protoOrder.NewOrderRequest, 0)
	for _, order := range newPlan.Orders {

		or := protoOrder.NewOrderRequest{
			OrderID:                order.OrderID,
			OrderPriority:          order.OrderPriority,
			OrderType:              order.OrderType,
			OrderTemplateID:        order.OrderTemplateID,
			AccountID:              order.AccountID,
			ParentOrderID:          order.ParentOrderID,
			MarketName:             order.MarketName,
			Grupo:                  order.Grupo,
			Side:                   order.Side,
			LimitPrice:             order.LimitPrice,
			InitialCurrencyBalance: order.InitialCurrencyBalance}

		for _, cond := range order.Triggers {
			trigger := protoOrder.TriggerRequest{
				TriggerTemplateID: cond.TriggerTemplateID,
				Index:             cond.Index,
				Title:             cond.Title,
				Name:              cond.Name,
				Code:              cond.Code,
				Actions:           cond.Actions,
			}
			or.Triggers = append(or.Triggers, &trigger)
		}

		newOrderRequests = append(newOrderRequests, &or)
	}

	newPlanRequest := protoPlan.NewPlanRequest{
		UserID:             userID,
		Title:              newPlan.Title,
		UserCurrencySymbol: newPlan.UserCurrencySymbol,
		PlanTemplateID:     newPlan.PlanTemplateID,
		Status:             newPlan.Status,
		CloseOnComplete:    newPlan.CloseOnComplete,
		InitialTimestamp:   newPlan.InitialTimestamp,
		Orders:             newOrderRequests,
	}

	// add plan returns nil for error
	r, _ := controller.PlanClient.NewPlan(context.Background(), &newPlanRequest)
	if r.Status != constRes.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == constRes.Fail {
			return c.JSON(http.StatusBadRequest, res)
		}
		if r.Status == constRes.Error {
			return c.JSON(http.StatusInternalServerError, res)
		}
	}
	plan := r.Data.Plan

	newOrders := make([]*Order, 0)
	for _, o := range plan.Orders {
		names := strings.Split(o.MarketName, "-")
		baseCurrencySymbol := names[1]
		baseCurrencyName := controller.currencies[baseCurrencySymbol]
		marketCurrencySymbol := names[0]
		marketCurrencyName := controller.currencies[marketCurrencySymbol]
		newo := Order{
			OrderID:                  o.OrderID,
			ParentOrderID:            o.ParentOrderID,
			PlanDepth:                o.PlanDepth,
			OrderTemplateID:          o.OrderTemplateID,
			AccountID:                o.AccountID,
			KeyPublic:                o.KeyPublic,
			KeyDescription:           o.KeyDescription,
			OrderPriority:            o.OrderPriority,
			OrderType:                o.OrderType,
			Side:                     o.Side,
			LimitPrice:               o.LimitPrice,
			Exchange:                 o.Exchange,
			MarketName:               o.MarketName,
			BaseCurrencySymbol:       baseCurrencySymbol,
			BaseCurrencyName:         baseCurrencyName,
			MarketCurrencySymbol:     marketCurrencySymbol,
			MarketCurrencyName:       marketCurrencyName,
			InitialCurrencySymbol:    o.InitialCurrencySymbol,
			InitialCurrencyName:      controller.currencies[o.InitialCurrencySymbol],
			InitialCurrencyBalance:   o.InitialCurrencyBalance,
			InitialCurrencyTraded:    o.InitialCurrencyTraded,
			InitialCurrencyRemainder: o.InitialCurrencyRemainder,
			FinalCurrencySymbol:      o.FinalCurrencySymbol,
			FinalCurrencyName:        controller.currencies[o.FinalCurrencySymbol],
			FinalCurrencyBalance:     o.FinalCurrencyBalance,
			Grupo:                    o.Grupo,
			Status:                   o.Status,
			CreatedOn:                o.CreatedOn,
			UpdatedOn:                o.UpdatedOn,
			Triggers:                 o.Triggers,
		}
		newOrders = append(newOrders, &newo)
	}

	res := &ResponsePlanSuccess{
		Status: constRes.Success,
		Data: &Plan{
			PlanID:                    plan.PlanID,
			PlanTemplateID:            plan.PlanTemplateID,
			PlanNumber:                plan.UserPlanNumber,
			TotalDepth:                plan.TotalDepth,
			Title:                     plan.Title,
			Exchange:                  plan.Exchange,
			UserCurrencySymbol:        plan.UserCurrencySymbol,
			UserCurrencyBalanceAtInit: plan.UserCurrencyBalanceAtInit,
			ActiveCurrencySymbol:      plan.ActiveCurrencySymbol,
			ActiveCurrencyName:        controller.currencies[plan.ActiveCurrencySymbol],
			ActiveCurrencyBalance:     plan.ActiveCurrencyBalance,
			InitialCurrencySymbol:     plan.InitialCurrencySymbol,
			InitialCurrencyName:       controller.currencies[plan.InitialCurrencySymbol],
			InitialCurrencyBalance:    plan.InitialCurrencyBalance,
			InitialTimestamp:          plan.InitialTimestamp,
			Status:                    plan.Status,
			CloseOnComplete:           plan.CloseOnComplete,
			LastExecutedOrderID:       plan.LastExecutedOrderID,
			LastExecutedPlanDepth:     plan.LastExecutedPlanDepth,
			Orders:                    newOrders,
			CreatedOn:                 plan.CreatedOn,
			UpdatedOn:                 plan.UpdatedOn,
		},
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:parameters UpdatePlanParams
type UpdatePlanRequest struct {
	// Optional base currency from which plan currency will be measured with.
	// in: body
	UserCurrencySymbol string `json:"userCurrencySymbol"`
	// Optional plan title
	// in: body
	Title string `json:"title"`
	// Optional plan template ID.
	// in: body
	PlanTemplateID string `json:"planTemplateID"`
	// Optional init timestamp for plan RFC3339 formatted (e.g. 2018-08-26T22:49:10.168652Z). This timestamp will be used to measure initial user currency balance (valuation in user preferred currency)
	// in: body
	InitialTimestamp string `json:"initialTimestamp"`
	// Optional only needed to update the status of the plan to 'inactive', 'active'
	// in: body
	Status string `json:"status"`
	// Optional bool to indicate that you want the plan to be 'closed' when the last order for the plan finishes (note: order status fail will also close the plan)
	// in: body
	CloseOnComplete bool `json:"closeOnComplete"`
	// Required array of orders. You cannot update executed orders. The entire inactive chain is assumed to be in this array.
	// in: body
	Orders []*NewOrderReq `json:"orders"`
}

// swagger:route PUT /plans/:planID plans UpdatePlanParams
//
// update a plan (protected)
//
// You must send in the entire inactive chain that you want updated in a single call.
//
// responses:
//  200: responsePlanSuccess "data" will contain plan with inactive orders (all orders that have yet to be executed) with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error" "data" will contain order info with "status": "success"
func (controller *PlanController) HandleUpdatePlan(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	planID := c.Param("planID")

	// read strategy from post body
	updatePlan := new(UpdatePlanRequest)
	err := c.Bind(&updatePlan)
	if err != nil {
		return fail(c, err.Error())
	}

	// assemble the order requests
	orderRequests := make([]*protoOrder.NewOrderRequest, 0)
	for _, order := range updatePlan.Orders {

		or := protoOrder.NewOrderRequest{
			OrderID:                order.OrderID,
			OrderPriority:          order.OrderPriority,
			OrderType:              order.OrderType,
			OrderTemplateID:        order.OrderTemplateID,
			AccountID:              order.AccountID,
			ParentOrderID:          order.ParentOrderID,
			MarketName:             order.MarketName,
			Side:                   order.Side,
			LimitPrice:             order.LimitPrice,
			InitialCurrencyBalance: order.InitialCurrencyBalance}

		for _, cond := range order.Triggers {
			trigger := protoOrder.TriggerRequest{
				TriggerTemplateID: cond.TriggerTemplateID,
				Index:             cond.Index,
				Title:             cond.Title,
				Name:              cond.Name,
				Code:              cond.Code,
				Actions:           cond.Actions,
			}
			or.Triggers = append(or.Triggers, &trigger)
		}

		orderRequests = append(orderRequests, &or)
	}

	updatePlanRequest := protoPlan.UpdatePlanRequest{
		PlanID:             planID,
		UserID:             userID,
		Title:              updatePlan.Title,
		UserCurrencySymbol: updatePlan.UserCurrencySymbol,
		PlanTemplateID:     updatePlan.PlanTemplateID,
		InitialTimestamp:   updatePlan.InitialTimestamp,
		Status:             updatePlan.Status,
		CloseOnComplete:    updatePlan.CloseOnComplete,
		Orders:             orderRequests,
	}

	// add plan returns nil for error
	r, _ := controller.PlanClient.UpdatePlan(context.Background(), &updatePlanRequest)
	if r.Status != constRes.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		switch r.Status {
		case constRes.Fail:
			return c.JSON(http.StatusBadRequest, res)
		case constRes.Error:
			return c.JSON(http.StatusInternalServerError, res)
		case constRes.Nonentity:
			return c.JSON(http.StatusBadRequest, res)
		}
	}
	plan := r.Data.Plan

	newOrders := make([]*Order, 0)
	for _, o := range plan.Orders {
		names := strings.Split(o.MarketName, "-")
		baseCurrencySymbol := names[1]
		baseCurrencyName := controller.currencies[baseCurrencySymbol]
		marketCurrencySymbol := names[0]
		marketCurrencyName := controller.currencies[marketCurrencySymbol]
		newo := Order{
			OrderID:                  o.OrderID,
			ParentOrderID:            o.ParentOrderID,
			PlanDepth:                o.PlanDepth,
			OrderTemplateID:          o.OrderTemplateID,
			AccountID:                o.AccountID,
			KeyPublic:                o.KeyPublic,
			KeyDescription:           o.KeyDescription,
			OrderPriority:            o.OrderPriority,
			OrderType:                o.OrderType,
			Side:                     o.Side,
			LimitPrice:               o.LimitPrice,
			Exchange:                 o.Exchange,
			ExchangeOrderID:          o.ExchangeOrderID,
			MarketName:               o.MarketName,
			BaseCurrencySymbol:       baseCurrencySymbol,
			BaseCurrencyName:         baseCurrencyName,
			MarketCurrencySymbol:     marketCurrencySymbol,
			MarketCurrencyName:       marketCurrencyName,
			InitialCurrencySymbol:    o.InitialCurrencySymbol,
			InitialCurrencyName:      controller.currencies[o.InitialCurrencySymbol],
			InitialCurrencyBalance:   o.InitialCurrencyBalance,
			InitialCurrencyTraded:    o.InitialCurrencyTraded,
			InitialCurrencyRemainder: o.InitialCurrencyRemainder,
			FinalCurrencySymbol:      o.FinalCurrencySymbol,
			FinalCurrencyName:        controller.currencies[o.FinalCurrencySymbol],
			FinalCurrencyBalance:     o.FinalCurrencyBalance,
			Grupo:                    o.Grupo,
			Status:                   o.Status,
			CreatedOn:                o.CreatedOn,
			UpdatedOn:                o.UpdatedOn,
			Triggers:                 o.Triggers,
		}
		newOrders = append(newOrders, &newo)
	}

	res := &ResponsePlanSuccess{
		Status: constRes.Success,
		Data: &Plan{
			PlanID:                    plan.PlanID,
			PlanTemplateID:            plan.PlanTemplateID,
			PlanNumber:                plan.UserPlanNumber,
			Title:                     plan.Title,
			TotalDepth:                plan.TotalDepth,
			Exchange:                  plan.Exchange,
			UserCurrencySymbol:        plan.UserCurrencySymbol,
			UserCurrencyBalanceAtInit: plan.UserCurrencyBalanceAtInit,
			ActiveCurrencySymbol:      plan.ActiveCurrencySymbol,
			ActiveCurrencyName:        controller.currencies[plan.ActiveCurrencySymbol],
			ActiveCurrencyBalance:     plan.ActiveCurrencyBalance,
			InitialCurrencySymbol:     plan.InitialCurrencySymbol,
			InitialCurrencyName:       controller.currencies[plan.InitialCurrencySymbol],
			InitialCurrencyBalance:    plan.InitialCurrencyBalance,
			InitialTimestamp:          plan.InitialTimestamp,
			Status:                    plan.Status,
			CloseOnComplete:           plan.CloseOnComplete,
			LastExecutedOrderID:       plan.LastExecutedOrderID,
			LastExecutedPlanDepth:     plan.LastExecutedPlanDepth,
			Orders:                    newOrders,
			CreatedOn:                 plan.CreatedOn,
			UpdatedOn:                 plan.UpdatedOn,
		},
	}

	return c.JSON(http.StatusOK, res)
}
