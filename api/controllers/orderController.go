package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	orderProto "github.com/asciiu/gomo/order-service/proto/order"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type OrderController struct {
	DB     *sql.DB
	Client orderProto.OrderServiceClient
}

type OrderTemp struct {
	orderID            string
	apiKeyId           string //Key id used for the order? Remember why we have this?
	exchangeOrderID    string
	baseCurrency       string // "BTC",
	baseCurrencyLong   string // "Bitcoin", //As above
	marketCurrency     string // "LTC",
	marketCurrencyLong string // "Litecoin", //Only bittrex seems to have this, pass the short name if doesn't exist
	minTradeSize       string //"0.001", //string
	marketName         string // "LTCBTC", //Convention is market+base this is our name
	//marketPrice: "0.41231231", //String Last price from socket for the pair in the exchange
	//?btcPrice: "0.41231231", //String This is a shortcut for me not to calculate we can discuss it
	//?fiatPrice: "1.35",  //Stting This is a shortcut for me not to calculate we can discuss it
	exchange           string // "binance"
	exchangeMarketName string // "LTC-BTC", //Some exchanges put dash others reverse them i.e. BTCLTC,
	orderType          string // limit, market, stop, fake_market, see above.
	rate               string //String
	baseQuantity       float64
	quantity           float64 // baseQuantity / rate
	quantityRemaining  float64 // how many
	side               string  // buy, sell
	conditions         string
	status             string //open, draft, closed,
	createdAt          int64  //integer
}

// swagger:parameters addOrder
type OrderRequest struct {
	// Required.
	// in: body
	ApiKeyId string `json:"apiKeyId"`
	// Required.
	// in: body
	MarketName string `json:"marketName"`
	// Required.
	// in: body
	Side string `json:"side"`
	// Optional.
	// in: body
	OrderType string `json:"orderType"`
	// Required.
	// in: body
	BaseQuantity float64 `json:"baseQuantity"`
	// Required.
	// in: body
	Conditions string `json:"conditions"`
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
	Status string                    `json:"status"`
	Data   *orderProto.UserOrderData `json:"data"`
}

// A ResponseKeysSuccess will always contain a status of "successful".
// swagger:model responseOrdersSuccess
type ResponseOrdersSuccess struct {
	Status string                     `json:"status"`
	Data   *orderProto.UserOrdersData `json:"data"`
}

func NewOrderController(db *sql.DB) *OrderController {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("apikey.client"))
	service.Init()

	controller := OrderController{
		DB:     db,
		Client: orderProto.NewOrderServiceClient("go.srv.order-service", service.Client()),
	}
	return &controller
}

// swagger:route GET /orders/:orderId orders getOrder
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
	userId := claims["jti"].(string)
	orderId := c.Param("orderId")

	getRequest := orderProto.GetUserOrderRequest{
		OrderId: orderId,
		UserId:  userId,
	}

	r, err := controller.Client.GetUserOrder(context.Background(), &getRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, response)
	}

	if r.Status != "success" {
		response := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == "fail" {
			return c.JSON(http.StatusBadRequest, response)
		}
		if r.Status == "error" {
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	response := &ResponseOrderSuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
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
	userId := claims["jti"].(string)

	getRequest := orderProto.GetUserOrdersRequest{
		UserId: userId,
	}

	r, err := controller.Client.GetUserOrders(context.Background(), &getRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, response)
	}

	if r.Status != "success" {
		response := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == "fail" {
			return c.JSON(http.StatusBadRequest, response)
		}
		if r.Status == "error" {
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	response := &ResponseOrdersSuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route POST /orders orders addOrder
//
// create a new order  (protected)
//
// This will create a new order in the system.
//
// responses:
//  200: responseOrderSuccess "data" will contain order info with "status": "success"
//  400: responseError missing params with "status": "fail"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *OrderController) HandlePostOrder(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["jti"].(string)

	order := OrderRequest{}

	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&order)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	createRequest := orderProto.OrderRequest{
		UserId:       userId,
		ApiKeyId:     order.ApiKeyId,
		MarketName:   order.MarketName,
		Side:         order.Side,
		Conditions:   order.Conditions,
		OrderType:    order.OrderType,
		BaseQuantity: order.BaseQuantity,
	}

	r, err := controller.Client.AddOrder(context.Background(), &createRequest)
	if err != nil {
		fmt.Println(err)
		response := &ResponseError{
			Status:  "error",
			Message: r.Message,
		}

		return c.JSON(http.StatusInternalServerError, response)
	}

	if r.Status != "success" {
		response := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == "fail" {
			return c.JSON(http.StatusBadRequest, response)
		}
		if r.Status == "error" {
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	response := &ResponseOrderSuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route PUT /orders/:orderId orders updateOrder
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
	userId := claims["jti"].(string)
	orderId := c.Param("orderId")

	orderRequest := UpdateOrderRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&orderRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	// client can only update description
	updateRequest := orderProto.OrderRequest{
		OrderId:      orderId,
		UserId:       userId,
		Conditions:   orderRequest.Conditions,
		BaseQuantity: orderRequest.BaseQuantity,
	}

	r, err := controller.Client.UpdateOrder(context.Background(), &updateRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, response)
	}

	if r.Status != "success" {
		response := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == "fail" {
			return c.JSON(http.StatusBadRequest, response)
		}
		if r.Status == "error" {
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	response := &ResponseOrderSuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route DELETE /orders/:orderId orders deleteOrder
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
	userId := claims["jti"].(string)
	orderId := c.Param("orderId")

	removeRequest := orderProto.RemoveOrderRequest{
		OrderId: orderId,
		UserId:  userId,
	}

	r, err := controller.Client.RemoveOrder(context.Background(), &removeRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, response)
	}

	if r.Status != "success" {
		response := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == "fail" {
			return c.JSON(http.StatusBadRequest, response)
		}
		if r.Status == "error" {
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	response := &ResponseOrderSuccess{
		Status: "success",
	}

	return c.JSON(http.StatusOK, response)
}
