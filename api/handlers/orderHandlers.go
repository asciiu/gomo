package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

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

func GetOrder(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("ERROR!")
	}

	log.Println("User name: ", claims["name"], "User ID: ", claims["jti"])
	// User ID from path `users/:id`
	id := c.Param("id")

	return c.JSON(http.StatusOK, map[string]string{
		"id":           id,
		"exchangeName": "amigonex",
	})
}

func ListOrders(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

func PostOrder(c echo.Context) error {
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

func UpdateOrder(c echo.Context) error {
	return c.String(http.StatusOK, "update it!")
}

func DeleteOrder(c echo.Context) error {
	return c.String(http.StatusOK, "delete it!")
}
