package controllers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	asql "github.com/asciiu/gomo/api/db/sql"
	"github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/response"
	keys "github.com/asciiu/gomo/key-service/proto/key"
	orders "github.com/asciiu/gomo/plan-service/proto/order"
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
// swagger:model ResponsePlansSuccess
type ResponsePlansSuccess struct {
	Status string     `json:"status"`
	Data   *PlansPage `json:"data"`
}

// A ResponsePlanWithOrderPageSuccess will always contain a status of "successful".
// swagger:model ResponsePlanWithOrderPageSuccess
type ResponsePlanWithOrderPageSuccess struct {
	Status string             `json:"status"`
	Data   *PlanWithOrderPage `json:"data"`
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

type PlanWithOrderPage struct {
	PlanID             string             `json:"planID"`
	PlanTemplateID     string             `json:"planTemplateID"`
	KeyID              string             `json:"keyID"`
	Exchange           string             `json:"exchange"`
	ExchangeMarketName string             `json:"exchangeMarketName"`
	MarketName         string             `json:"marketName"`
	BaseCurrencySymbol string             `json:"baseCurrencySymbol"`
	BaseCurrencyName   string             `json:"baseCurrencyName"`
	BaseBalance        float64            `json:"baseBalance"`
	CurrencySymbol     string             `json:"currencySymbol"`
	CurrencyName       string             `json:"currencyName"`
	CurrencyBalance    float64            `json:"currencyBalance"`
	Status             string             `json:"status"`
	CreatedOn          string             `json:"createdOn"`
	UpdatedOn          string             `json:"updatedOn"`
	OrdersPage         *orders.OrdersPage `json:"ordersPage"`
}

type OrdersPage struct {
	Page     uint32          `json:"page"`
	PageSize uint32          `json:"pageSize"`
	Total    uint32          `json:"total"`
	Orders   []*orders.Order `json:"orders"`
}

type UserPlansData struct {
	PLans []*Plan `json:"plans"`
}

// This response should never return the key secret
type Plan struct {
	PlanID                string   `json:"planID"`
	PlanTemplateID        string   `json:"planTemplateID"`
	Exchange              string   `json:"exchange"`
	ExchangeMarketName    string   `json:"exchangeMarketName"`
	MarketName            string   `json:"marketName"`
	CurrencySymbol        string   `json:"activeCurrencySymbol"`
	CurrencyName          string   `json:"activeCurrencyName"`
	CurrencyBalance       float64  `json:"activeCurrencyBalance"`
	LastExecutedOrderID   string   `json:"lastExecutedOrderID"`
	LastExecutedPlanDepth uint32   `json:"lastExecutedPlanDepth"`
	Status                string   `json:"status"`
	CloseOnComplete       bool     `json:"closeOnComplete"`
	CreatedOn             string   `json:"createdOn"`
	UpdatedOn             string   `json:"updatedOn"`
	Orders                []*Order `json:"orders,omitempty"`
}

type Order struct {
	OrderID            string            `name=orderID" json:"orderID,omitempty"`
	ParentOrderID      string            `name=parentOrderID" json:"parentOrderID,omitempty"`
	PlanDepth          uint32            `name=planDepth" json:"planDepth,omitempty"`
	OrderTemplateID    string            `name=orderTemplateID" json:"orderTemplateID,omitempty"`
	KeyID              string            `name=keyID" json:"keyID,omitempty"`
	KeyPublic          string            `name=keyPublic" json:"keyPublic,omitempty"`
	KeyDescription     string            `name=keyDescription" json:"keyDescription,omitempty"`
	OrderPriority      uint32            `name=orderPriority" json:"orderPriority,omitempty"`
	OrderType          string            `name=orderType" json:"orderType,omitempty"`
	Side               string            `name=side" json:"side,omitempty"`
	LimitPrice         float64           `name=limitPrice" json:"limitPrice,omitempty"`
	Exchange           string            `name=exchange" json:"exchange,omitempty"`
	ExchangeMarketName string            `name=exchangeMarketName" json:"exchangeMarketName,omitempty"`
	MarketName         string            `name=marketName" json:"marketName,omitempty"`
	BaseCurrencySymbol string            `json:"baseCurrencySymbol"`
	BaseCurrencyName   string            `json:"baseCurrencyName"`
	CurrencySymbol     string            `name=currencySymbol" json:"currencySymbol,omitempty"`
	CurrencyName       string            `json:"currencyName"`
	CurrencyBalance    float64           `name=currencyBalance" json:"currencyBalance,omitempty"`
	CurrencyTraded     float64           `name=currencyTraded" json:"currencyTraded,omitempty"`
	Status             string            `name=status" json:"status,omitempty"`
	CreatedOn          string            `name=createdOn" json:"createdOn,omitempty"`
	UpdatedOn          string            `name=updatedOn" json:"updatedOn,omitempty"`
	Triggers           []*orders.Trigger `name=triggers" json:"triggers,omitempty"`
}

func fail(c echo.Context, msg string) error {
	res := &ResponseError{
		Status:  response.Fail,
		Message: msg,
	}

	return c.JSON(http.StatusBadRequest, res)
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

// swagger:route DELETE /plans/:planID plans DeletePlan
//
// deletes a plan (protected)
//
// You may delete a plan if it has not executed. That is, the plan has no filled orders. Delete plan becomes cancel plan
// when an order has been filled. In theory, this should kill all active orders for a plan and set the status for the plan
// as 'aborted'. Once a plan has been 'aborted' you cannot update or restart that plan.
//
// responses:
//  200: ResponsePlanSuccess "data" will contain plan summary with null orders.
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *PlanController) HandleDeletePlan(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	planID := c.Param("planID")

	delRequest := plans.DeletePlanRequest{
		PlanID: planID,
		UserID: userID,
	}

	r, _ := controller.Plans.DeletePlan(context.Background(), &delRequest)
	if r.Status != response.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		switch {
		case r.Status == response.Nonentity:
			return c.JSON(http.StatusNotFound, res)
		case r.Status == response.Fail:
			return c.JSON(http.StatusBadRequest, res)
		default:
			return c.JSON(http.StatusInternalServerError, res)
		}
	}

	names := strings.Split(r.Data.Plan.MarketName, "-")
	baseCurrencySymbol := names[1]
	baseCurrencyName := controller.currencies[baseCurrencySymbol]
	currencySymbol := names[0]
	currencyName := controller.currencies[currencySymbol]

	res := &ResponsePlanWithOrderPageSuccess{
		Status: response.Success,
		Data: &PlanWithOrderPage{
			PlanID:         r.Data.Plan.PlanID,
			PlanTemplateID: r.Data.Plan.PlanTemplateID,
			//KeyID:              r.Data.Plan.KeyID,
			Exchange:           r.Data.Plan.Exchange,
			MarketName:         r.Data.Plan.MarketName,
			BaseCurrencySymbol: baseCurrencySymbol,
			BaseCurrencyName:   baseCurrencyName,
			//BaseBalance:        r.Data.Plan.BaseBalance,
			CurrencySymbol:  currencySymbol,
			CurrencyName:    currencyName,
			CurrencyBalance: r.Data.Plan.CurrencyBalance,
			Status:          r.Data.Plan.Status,
			CreatedOn:       r.Data.Plan.CreatedOn,
			UpdatedOn:       r.Data.Plan.UpdatedOn,
		},
	}

	return c.JSON(http.StatusOK, res)
}

// required for swaggered, otherwise never used
// swagger:parameters GetPlanParams
type GetPlanParams struct {
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"pageSize"`
}

// swagger:route GET /plans/:planID plans GetPlanParams
//
// get plan with planID (protected)
//
// Returns a plan with completd currency and base currency names. There is no limit on planned orders - these are the
// orders that belong to a trade plan. Therefore, the response data will contain a paged structure for the plan orders.
//
// example: /plans/:some_plan_id?page0&pageSize=50
//
// responses:
//  200: ResponsePlanWithOrderPageSuccess "data" will contain paged orders
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *PlanController) HandleGetPlan(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	planID := c.Param("planID")

	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("pageSize")
	// defaults for page and page size here
	// ignore the errors and assume the values are int
	page, _ := strconv.ParseUint(pageStr, 10, 32)
	pageSize, _ := strconv.ParseUint(pageSizeStr, 10, 32)
	if pageSize == 0 {
		pageSize = 20
	}

	getRequest := plans.GetUserPlanRequest{
		PlanID:   planID,
		UserID:   userID,
		Page:     uint32(page),
		PageSize: uint32(pageSize),
	}

	r, _ := controller.Plans.GetUserPlan(context.Background(), &getRequest)
	if r.Status != response.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		switch {
		case r.Status == response.Nonentity:
			return c.JSON(http.StatusNotFound, res)
		case r.Status == response.Fail:
			return c.JSON(http.StatusBadRequest, res)
		default:
			return c.JSON(http.StatusInternalServerError, res)
		}
	}

	names := strings.Split(r.Data.MarketName, "-")
	baseCurrencySymbol := names[1]
	baseCurrencyName := controller.currencies[baseCurrencySymbol]
	currencySymbol := names[0]
	currencyName := controller.currencies[currencySymbol]

	res := &ResponsePlanWithOrderPageSuccess{
		Status: response.Success,
		Data: &PlanWithOrderPage{
			PlanID:             r.Data.PlanID,
			PlanTemplateID:     r.Data.PlanTemplateID,
			KeyID:              r.Data.KeyID,
			Exchange:           r.Data.Exchange,
			MarketName:         r.Data.MarketName,
			BaseCurrencySymbol: baseCurrencySymbol,
			BaseCurrencyName:   baseCurrencyName,
			BaseBalance:        r.Data.BaseBalance,
			CurrencySymbol:     currencySymbol,
			CurrencyName:       currencyName,
			CurrencyBalance:    r.Data.CurrencyBalance,
			Status:             r.Data.Status,
			OrdersPage:         r.Data.OrdersPage,
			CreatedOn:          r.Data.CreatedOn,
			UpdatedOn:          r.Data.UpdatedOn,
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
// Returns a summary of plans. The plan orders will not be returned and will be null.
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
		//baseCurrencySymbol := names[1]
		//baseCurrencyName := controller.currencies[baseCurrencySymbol]
		currencySymbol := names[0]
		currencyName := controller.currencies[currencySymbol]

		pln := Plan{
			PlanTemplateID: plan.PlanTemplateID,
			PlanID:         plan.PlanID,
			//KeyID:                 plan.KeyID,
			Exchange: plan.Exchange,
			//ExchangeMarketName:    plan.ExchangeMarketName,
			MarketName: plan.MarketName,
			//BaseCurrencySymbol: baseCurrencySymbol,
			//BaseCurrencyName:   baseCurrencyName,
			//BaseBalance:           plan.BaseBalance,
			CurrencySymbol:        currencySymbol,
			CurrencyName:          currencyName,
			CurrencyBalance:       plan.CurrencyBalance,
			Status:                plan.Status,
			LastExecutedOrderID:   plan.LastExecutedOrderID,
			LastExecutedPlanDepth: plan.LastExecutedPlanDepth,
			CreatedOn:             plan.CreatedOn,
			UpdatedOn:             plan.UpdatedOn,
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

// swagger:parameters PostPlan
type PlanRequest struct {
	// Optional plan template ID.
	// in: body
	PlanTemplateID string `json:"planTemplateID"`
	// Optional defaults to 'active' status. Valid input status is 'active', 'inactive', or 'historic'
	// in: body
	Status string `json:"status"`
	// Required bool to indicate that you want the plan to be 'closed' when the last order for the plan finishes (note: order status fail will also close the plan)
	// in: body
	CloseOnComplete bool `json:"closeOnComplete"`
	// Required array of orders. The structure of the order tree will be dictated by the orderNumber and parentOrderNumber properties of each order.
	// in: body
	Orders []*NewOrderReq `json:"orders"`
}

type NewOrderReq struct {
	// Required the client assigns the order ID as a UUID, the format is 8-4-4-4-12.
	// in: body
	OrderID string `json:"orderID"`
	// Optional precedence of orders when multiple orders are at the same depth: value of 1 is highest priority. Example: depth 2 buy ADA (1) or buy EOS (2). ADA with higher priority 1 will execute and EOS will not execute.
	// in: body
	OrderPriority uint32 `json:"orderPriority"`
	// Required order types are "market", "limit", "paper". Orders not within these types will be rejected.
	// in: body
	OrderType string `json:"orderType"`
	// Optional order template ID.
	// in: body
	OrderTemplateID string `json:"orderTemplateID"`
	// Required this is our api key ID (string uuid) assigned to the user's exchange key and secret.
	// in: body
	KeyID string `json:"keyID"`
	// Required the root node of the decision tree should be assigned a parentOrderID of "00000000-0000-0000-0000-000000000000" .
	// in: body
	ParentOrderID string `json:"parentOrderID"`
	// Required e.g. ADA-BTC. Base pair should be the suffix.
	// in: body
	MarketName string `json:"marketName"`
	// Required "buy" or "sell"
	// in: body
	Side string `json:"side"`
	// Required for 'limit' orders. Defines limit price.
	// in: body
	LimitPrice float64 `json:"limitPrice"`
	// Required for the root order of the tree. Child orders for tree may or may not have a currencyBalance.
	// in: body
	CurrencyBalance float64 `json:"currencyBalance"`
	// Required these are the conditions that trigger the order to execute: ???
	// in: body
	Triggers []*TriggerReq `json:"triggers"`
}

type TriggerReq struct {
	// Optional trigger template ID.
	// in: body
	TriggerTemplateID string `json:"triggerTemplateID"`

	Name    string
	Code    string
	Actions []string
}

// swagger:route POST /plans plans PostPlan
//
// create a new plan (protected)
//
// This will create a new chain of orders for the user. All orders are encapsulated within a plan.
//
// responses:
//  200: ResponsePlanWithOrderPageSuccess "data" will contain paged orders
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
	newOrderRequests := make([]*orders.NewOrderRequest, 0)
	for _, order := range newPlan.Orders {

		or := orders.NewOrderRequest{
			OrderID:         order.OrderID,
			OrderPriority:   order.OrderPriority,
			OrderType:       order.OrderType,
			OrderTemplateID: order.OrderTemplateID,
			KeyID:           order.KeyID,
			ParentOrderID:   order.ParentOrderID,
			MarketName:      order.MarketName,
			Side:            order.Side,
			LimitPrice:      order.LimitPrice,
			CurrencyBalance: order.CurrencyBalance}

		for _, cond := range order.Triggers {
			trigger := orders.TriggerRequest{
				TriggerTemplateID: cond.TriggerTemplateID,
				Name:              cond.Name,
				Code:              cond.Code,
				Actions:           cond.Actions,
			}
			or.Triggers = append(or.Triggers, &trigger)
		}

		newOrderRequests = append(newOrderRequests, &or)
	}

	newPlanRequest := plans.NewPlanRequest{
		UserID:          userID,
		PlanTemplateID:  newPlan.PlanTemplateID,
		Status:          newPlan.Status,
		CloseOnComplete: newPlan.CloseOnComplete,
		Orders:          newOrderRequests,
	}

	// add plan returns nil for error
	r, _ := controller.Plans.NewPlan(context.Background(), &newPlanRequest)
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

	newOrders := make([]*Order, 0)
	for _, o := range r.Data.Plan.Orders {
		names := strings.Split(o.MarketName, "-")
		baseCurrencySymbol := names[1]
		baseCurrencyName := controller.currencies[baseCurrencySymbol]
		currencySymbol := names[0]
		currencyName := controller.currencies[currencySymbol]
		newo := Order{
			OrderID:            o.OrderID,
			ParentOrderID:      o.ParentOrderID,
			PlanDepth:          o.PlanDepth,
			OrderTemplateID:    o.OrderTemplateID,
			KeyID:              o.KeyID,
			KeyPublic:          o.KeyPublic,
			KeyDescription:     o.KeyDescription,
			OrderPriority:      o.OrderPriority,
			OrderType:          o.OrderType,
			Side:               o.Side,
			LimitPrice:         o.LimitPrice,
			Exchange:           o.Exchange,
			ExchangeMarketName: o.ExchangeMarketName,
			MarketName:         o.MarketName,
			BaseCurrencySymbol: baseCurrencySymbol,
			BaseCurrencyName:   baseCurrencyName,
			CurrencySymbol:     o.CurrencySymbol,
			CurrencyName:       currencyName,
			CurrencyBalance:    o.CurrencyBalance,
			CurrencyTraded:     o.CurrencyTraded,
			Status:             o.Status,
			CreatedOn:          o.CreatedOn,
			UpdatedOn:          o.UpdatedOn,
			Triggers:           o.Triggers,
		}
		newOrders = append(newOrders, &newo)
	}

	res := &ResponsePlanSuccess{
		Status: response.Success,
		Data: &Plan{
			PlanID:                r.Data.Plan.PlanID,
			PlanTemplateID:        r.Data.Plan.PlanTemplateID,
			Exchange:              r.Data.Plan.Exchange,
			MarketName:            r.Data.Plan.MarketName,
			CurrencySymbol:        r.Data.Plan.CurrencySymbol,
			CurrencyName:          controller.currencies[r.Data.Plan.CurrencySymbol],
			CurrencyBalance:       r.Data.Plan.CurrencyBalance,
			Status:                r.Data.Plan.Status,
			CloseOnComplete:       r.Data.Plan.CloseOnComplete,
			LastExecutedOrderID:   r.Data.Plan.LastExecutedOrderID,
			LastExecutedPlanDepth: r.Data.Plan.LastExecutedPlanDepth,
			Orders:                newOrders,
			CreatedOn:             r.Data.Plan.CreatedOn,
			UpdatedOn:             r.Data.Plan.UpdatedOn,
		},
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:parameters UpdatePlanParams
type UpdatePlanRequest struct {
	// Optional change base balance of unexecuted plan
	// in: body
	BaseBalance float64 `json:"baseBalance"`
	// Optional change currency balance of unexecuted plan
	// in: body
	CurrencyBalance float64 `json:"currencyBalance"`
	// Optional send 'inactive' to pause and 'active' to unpause
	// in: body
	Status string `json:"status"`

	// Required array of orders. The structure of the order tree will be dictated by the orderNumber and parentOrderNumber properties of each order.
	// in: body
	Orders []*UpdateOrderReq `json:"orders"`
}

type UpdateOrderReq struct {
	// Optional order ID
	OrderID string `json:"orderID"`
	// Required number for order. Order number should begin at 1.
	// in: body
	OrderNumber uint32 `json:"orderNumber"`
	// Required "buy" or "sell"
	// in: body
	Side string `json:"side"`
	// Optional order template ID.
	// in: body
	OrderTemplateID string `json:"orderTemplateID"`
	// Required order types are "market", "limit", "paper". Orders not within these types will be rejected.
	// in: body
	OrderType string `json:"orderType"`
	// Required for the precent of your plan's balance to use for the order - buy (percent of base balance) - sell (percent of currency balance)
	// in: body
	BalancePercent float64 `json:"balancePercent"`
	// Required for 'limit' orders. Defines limit price.
	// in: body
	LimitPrice float64 `json:"limitPrice"`
	// Required these are the conditions that trigger the order to execute: ???
	// in: body
	Triggers []*TriggerReq `json:"triggers"`
	// Optional the root node of the decision tree will have a parent order number of 0. 0 delineates no parent.
	// in: body
	ParentOrderNumber uint32 `json:"parentOrderNumber"`
	// Required action new, update, delete
	Action string `json:"action"`
}

// swagger:route PUT /plans/:planID plans UpdatePlanParams
//
// update a plan (protected)
//
// You may update the base balance and currency balance before the plan has executed its first order. Once a
// plan order has been executed you cannot change the balances for the plan.
// Other use case for this endpoint is to pause the plan by sending status='inactive'. Use DELETE to abort the plan.
//
// responses:
//  200: responsePlanSuccess "data" will contain plan summary with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error" "data" will contain order info with "status": "success"
func (controller *PlanController) HandleUpdatePlan(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	planID := c.Param("planID")

	requestBody, _ := ioutil.ReadAll(c.Request().Body)

	// there's got to be a better way to do this validation
	if !strings.Contains(string(requestBody), "baseBalance") ||
		!strings.Contains(string(requestBody), "currencyBalance") ||
		!strings.Contains(string(requestBody), "status") {
		res := &ResponseError{
			Status:  response.Fail,
			Message: "baseBalance, currencyBalance, and status are required",
		}
		return c.JSON(http.StatusBadRequest, res)
	}

	var updateParams UpdatePlanRequest

	err := json.Unmarshal([]byte(requestBody), &updateParams)
	if err != nil {
		res := &ResponseError{
			Status:  response.Fail,
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	updateRequest := plans.UpdatePlanRequest{
		PlanID: planID,
		UserID: userID,
		Status: updateParams.Status,
		//BaseBalance:     updateParams.BaseBalance,
		//CurrencyBalance: updateParams.CurrencyBalance,
	}

	r, _ := controller.Plans.UpdatePlan(context.Background(), &updateRequest)
	if r.Status != response.Success {
		res := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		switch {
		case r.Status == response.Nonentity:
			return c.JSON(http.StatusNotFound, res)
		case r.Status == response.Fail:
			return c.JSON(http.StatusBadRequest, res)
		default:
			return c.JSON(http.StatusInternalServerError, res)
		}
	}

	names := strings.Split(r.Data.Plan.MarketName, "-")
	baseCurrencySymbol := names[1]
	baseCurrencyName := controller.currencies[baseCurrencySymbol]
	currencySymbol := names[0]
	currencyName := controller.currencies[currencySymbol]

	res := &ResponsePlanWithOrderPageSuccess{
		Status: response.Success,
		Data: &PlanWithOrderPage{
			PlanID:         r.Data.Plan.PlanID,
			PlanTemplateID: r.Data.Plan.PlanTemplateID,
			//KeyID:              r.Data.Plan.KeyID,
			Exchange:           r.Data.Plan.Exchange,
			MarketName:         r.Data.Plan.MarketName,
			BaseCurrencySymbol: baseCurrencySymbol,
			BaseCurrencyName:   baseCurrencyName,
			//BaseBalance:        r.Data.Plan.BaseBalance,
			CurrencySymbol:  currencySymbol,
			CurrencyName:    currencyName,
			CurrencyBalance: r.Data.Plan.CurrencyBalance,
			Status:          r.Data.Plan.Status,
			CreatedOn:       r.Data.Plan.CreatedOn,
			UpdatedOn:       r.Data.Plan.UpdatedOn,
		},
	}

	return c.JSON(http.StatusOK, res)
}
