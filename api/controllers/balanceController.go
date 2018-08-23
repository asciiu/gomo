package controllers

import (
	"database/sql"
	"net/http"

	protoBalance "github.com/asciiu/gomo/balance-service/proto/balance"
	constRes "github.com/asciiu/gomo/common/constants/response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

// A ResponseBalancesSuccess will always contain a status of "successful".
// swagger:model responseBalancesSuccess
type ResponseBalancesSuccess struct {
	Status string           `json:"status"`
	Data   *AccountBalances `json:"data"`
}

// This struct is used in the generated swagger docs,
// and it is not used anywhere.
// swagger:parameters getAllBalances
type SearchSymbol struct {
	// Required: false
	// In: query
	Symbol string `json:"symbol"`
}

type AccountBalances struct {
	Balances []*Balance `json:"protoBalance"`
}

type Balance struct {
	BalanceID         string  `json:"balanceID"`
	KeyID             string  `json:"keyID"`
	ExchangeName      string  `json:"exchange"`
	CurrencyName      string  `json:"currencyName"`
	Available         float64 `json:"available"`
	Locked            float64 `json:"locked"`
	ExchangeTotal     float64 `json:"exchangeTotal"`
	ExchangeAvailable float64 `json:"exchangeAvailable"`
	ExchangeLocked    float64 `json:"exchangeLocked"`
}

type BalanceController struct {
	DB       *sql.DB
	Balances protoBalance.BalanceServiceClient
}

func NewBalanceController(db *sql.DB, service micro.Service) *BalanceController {
	controller := BalanceController{
		DB:       db,
		Balances: protoBalance.NewBalanceServiceClient("protoBalance", service.Client()),
	}
	return &controller
}

// swagger:route GET /protoBalance protoBalance getAllBalances
//
// get all protoBalance (protected)
//
// Returns all protoBalance for user. Use optional query param 'symbol' as lowercase ticker symbol - e.g. ada.
//
// responses:
//  200: responseBalancesSuccess "data" will contain array of protoBalance with "status": constRes.Success
//  500: responseError the message will state what the internal server error was with "status": constRes.Error
func (controller *BalanceController) HandleGetBalances(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	symbol := c.QueryParam("symbol")

	getRequest := protoBalance.GetUserBalancesRequest{
		UserID: userID,
		Symbol: symbol,
	}

	r, err := controller.Balances.GetUserBalances(context.Background(), &getRequest)
	if err != nil {
		response := &ResponseError{
			Status:  constRes.Error,
			Message: err.Error(),
		}

		return c.JSON(http.StatusGone, response)
	}

	if r.Status != constRes.Success {
		response := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == constRes.Fail {
			return c.JSON(http.StatusBadRequest, response)
		}
		if r.Status == constRes.Error {
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	data := make([]*Balance, len(r.Data.Balances))
	for i, balance := range r.Data.Balances {
		// api removes the secret
		data[i] = &Balance{
			BalanceID:         balance.ID,
			KeyID:             balance.KeyID,
			ExchangeName:      balance.ExchangeName,
			CurrencyName:      balance.CurrencyName,
			Available:         balance.Available,
			Locked:            balance.Locked,
			ExchangeTotal:     balance.ExchangeTotal,
			ExchangeAvailable: balance.ExchangeAvailable,
			ExchangeLocked:    balance.ExchangeLocked,
		}
	}

	response := &ResponseBalancesSuccess{
		Status: constRes.Success,
		Data: &AccountBalances{
			Balances: data,
		},
	}

	return c.JSON(http.StatusOK, response)
}
