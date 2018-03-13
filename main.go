package main

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Order struct {
	orderID string 
	apiKeyId string //Key id used for the order? Remember why we have this?
	exchangeOrderID string
	baseCurrency string // "BTC",
    baseCurrencyLong string // "Bitcoin", //As above
	marketCurrency string // "LTC",
    marketCurrencyLong string // "Litecoin", //Only bittrex seems to have this, pass the short name if doesn't exist
    minTradeSize string  //"0.001", //string
    marketName string // "LTCBTC", //Convention is market+base this is our name
    //marketPrice: "0.41231231", //String Last price from socket for the pair in the exchange
    //?btcPrice: "0.41231231", //String This is a shortcut for me not to calculate we can discuss it
    //?fiatPrice: "1.35",  //Stting This is a shortcut for me not to calculate we can discuss it          
    exchange string // "binance"
    exchangeMarketName string // "LTC-BTC", //Some exchanges put dash others reverse them i.e. BTCLTC, 
	orderType string // limit, market, stop, fake_market, see above.
	rate string //String
	baseQuantity float64
	quantity float64 // baseQuantity / rate
	quantityRemaining float64 // how many 
	side string // buy, sell
	conditions string
	status string //open, draft, closed,
	createdAt int64 //integer
}

type OrderRequest struct {
	ApiKeyId string `json:"apiKeyId"`
	ExchangeMarketName string `json:"exchangeMarketName"`
	MarketName string `json:"marketName"`
	Side string `json:"side"`
	OrderType string `json:"orderType"`
	UnitPrice float64 `json:"unitPrice"`
	Qauntity float64 `json:"quantity"`
	Conditions string `json:"conditions"`
}

// e.GET("/users/:id", getUser)
func getOrder(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")

	return c.JSON(http.StatusOK, map[string]string{
		"id": id,
		"exchangeName": "amigonex",
	})
}

func listOrders(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:" + team + ", member:" + member)
}

func postOrder(c echo.Context) error {
	order := OrderRequest{}

	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&order)
	if err != nil {
		log.Printf("failed reading the request %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("this is your order: %#v", order)
	return c.String(http.StatusOK, "post order!")
} 

func updateOrder(c echo.Context) error {
	return c.String(http.StatusOK, "update it!")
}

func deleteOrder(c echo.Context) error {
	return c.String(http.StatusOK, "delete it!")
}

func mainHandler(c echo.Context) error {
	return c.String(http.StatusOK, "main")
}

func main() {
	e := echo.New()

	g := e.Group("/api")
	// this logs the server interaction
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]  ${status}  ${method}  ${host}${path} ${latency_human}` + "\n",
	}))
	g.GET("/main", mainHandler)


	e.GET("/orders", listOrders)
	e.POST("/orders", postOrder)
	e.GET("/orders/:id", getOrder)
	e.PUT("/orders/:id", updateOrder)
	e.DELETE("/orders/:id", deleteOrder)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":5000"))
}