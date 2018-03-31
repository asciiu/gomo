package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type OrderController struct {
	DB *sql.DB
}

type Order struct {
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

type OrderRequest struct {
	ApiKeyId           string  `json:"apiKeyId"`
	ExchangeMarketName string  `json:"exchangeMarketName"`
	MarketName         string  `json:"marketName"`
	Side               string  `json:"side"`
	OrderType          string  `json:"orderType"`
	UnitPrice          float64 `json:"unitPrice"`
	Qauntity           float64 `json:"quantity"`
	Conditions         string  `json:"conditions"`
}

func NewOrderController(db *sql.DB) *OrderController {
	controller := OrderController{
		DB: db,
	}
	return &controller
}

// swagger:route GET /orders/:orderId orders getOrder
//
// not implemented (protected)
//
// ...
func (controller *OrderController) HandleGetOrder(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("ERROR!")
	}

	log.Println("User name: ", claims["name"], "User ID: ", claims["jti"])
	// User ID from path `users/:id`
	id := c.Param("orderId")

	return c.JSON(http.StatusOK, map[string]string{
		"id":           id,
		"exchangeName": "amigonex",
	})
}

// swagger:route GET /orders orders getAllOrders
//
// not implemented (protected)
//
// ...
func (controller *OrderController) HandleListOrders(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

// swagger:route POST /orders orders addOrder
//
// not implemented (protected)
//
// ...
func (controller *OrderController) HandlePostOrder(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)

	order := OrderRequest{}

	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&order)
	if err != nil {
		log.Printf("failed reading the request %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("this is your order: %#v", order)
	return c.String(http.StatusOK, "Welcome "+name+" your order has posted!")
}

// swagger:route PUT /orders/:orderId orders updateOrder
//
// not implemented (protected)
//
// ...
func (controller *OrderController) HandleUpdateOrder(c echo.Context) error {
	return c.String(http.StatusOK, "update it!")
}

// swagger:route DELETE /orders/:orderId orders deleteOrder
//
// not implemented (protected)
//
// ...
func (controller *OrderController) HandleDeleteOrder(c echo.Context) error {
	return c.String(http.StatusOK, "delete it!")
}
