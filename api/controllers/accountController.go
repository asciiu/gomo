package controllers

import (
	"database/sql"
	"log"
	"net/http"

	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	asql "github.com/asciiu/gomo/api/db/sql"
	constRes "github.com/asciiu/gomo/common/constants/response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type AccountController struct {
	DB            *sql.DB
	AccountClient protoAccount.AccountServiceClient
	// map of ticker symbol to full name
	currencies map[string]string
}

type AccountList struct {
	Accounts []*Account `json:"accounts"`
}

// A ResponseAccountsSuccess will always contain a status of "successful".
// swagger:model ResponseAccountsSuccess
type ResponseAccountListSuccess struct {
	Status string      `json:"status"`
	Data   AccountList `json:"data"`
}

// A ResponseAccountSuccess will always contain a status of "successful".
// swagger:model ResponseAccountSuccess
type ResponseAccountSuccess struct {
	Status string   `json:"status"`
	Data   *Account `json:"data"`
}

// This response should never return the key secret
type Account struct {
	AccountID   string     `json:"accountID"`
	AccountType string     `json:"type"`
	Exchange    string     `json:"exchange"`
	KeyPublic   string     `json:"keyPublic"`
	Description string     `json:"description"`
	CreatedOn   string     `json:"createdOn"`
	UpdatedOn   string     `json:"updatedOn"`
	Status      string     `json:"status"`
	Balances    []*Balance `json:"balances"`
}

type Balance struct {
	CurrencySymbol    string  `json:"currencySymbol"`
	Available         float64 `json:"available"`
	Locked            float64 `json:"locked"`
	ExchangeTotal     float64 `json:"exchangeTotal"`
	ExchangeAvailable float64 `json:"exchangeAvailable"`
	ExchangeLocked    float64 `json:"exchangeLocked"`
	CreatedOn         string  `json:"createdOn"`
	UpdatedOn         string  `json:"updatedOn"`
}

func NewAccountController(db *sql.DB, service micro.Service) *AccountController {
	controller := AccountController{
		DB:            db,
		AccountClient: protoAccount.NewAccountServiceClient("accounts", service.Client()),
		currencies:    make(map[string]string),
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

// swagger:route DELETE /accounts/:accountID accounts DeleteAccount
//
// soft delete account (protected)
//
// This will set the status of an account to deleted.
//
// responses:
//  200: ResponseAccountSuccess "data" will contain account summary.
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *AccountController) HandleDeleteAccount(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	accountID := c.Param("accountID")

	delRequest := protoAccount.AccountRequest{AccountID: accountID, UserID: userID}
	r, _ := controller.AccountClient.DeleteAccount(context.Background(), &delRequest)
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

	res := &ResponseAccountSuccess{
		Status: constRes.Success,
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:route GET /accounts/:accountID accounts GetAccountParms
//
// get account by accountID (protected)
//
// Returns account deets with balances.
//
// responses:
//  200: ResponseAccountSuccess "data" will contain account deets.
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *AccountController) HandleGetAccount(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	accountID := c.Param("accountID")

	getRequest := protoAccount.AccountRequest{AccountID: accountID, UserID: userID}
	r, _ := controller.AccountClient.GetAccount(context.Background(), &getRequest)
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
	account := r.Data.Account
	balances := make([]*Balance, 0)
	for _, b := range account.Balances {
		balance := Balance{
			CurrencySymbol:    b.CurrencySymbol,
			Available:         b.Available,
			Locked:            b.Locked,
			ExchangeTotal:     b.ExchangeTotal,
			ExchangeLocked:    b.ExchangeLocked,
			ExchangeAvailable: b.ExchangeAvailable,
			CreatedOn:         b.CreatedOn,
			UpdatedOn:         b.UpdatedOn,
		}
		balances = append(balances, &balance)
	}

	res := &ResponseAccountSuccess{
		Status: constRes.Success,
		Data: &Account{
			AccountID:   account.AccountID,
			AccountType: account.AccountType,
			Exchange:    account.Exchange,
			KeyPublic:   account.KeyPublic,
			Description: account.Description,
			CreatedOn:   account.CreatedOn,
			UpdatedOn:   account.UpdatedOn,
			Status:      account.Status,
			Balances:    balances,
		},
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:route GET /accounts accounts GetUserAccountsParams
//
// get user accounts (protected)
//
// Returns all accounts with their balances.
//
// responses:
//  200: ResponseAccountListSuccess "data" will contain an array of accounts
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *AccountController) HandleListAccounts(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	getRequest := protoAccount.AccountsRequest{UserID: userID}
	r, _ := controller.AccountClient.GetAccounts(context.Background(), &getRequest)
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

	accounts := make([]*Account, 0)
	for _, a := range r.Data.Accounts {

		balances := make([]*Balance, 0)
		for _, b := range a.Balances {
			balance := Balance{
				CurrencySymbol:    b.CurrencySymbol,
				Available:         b.Available,
				Locked:            b.Locked,
				ExchangeTotal:     b.ExchangeTotal,
				ExchangeLocked:    b.ExchangeLocked,
				ExchangeAvailable: b.ExchangeAvailable,
				CreatedOn:         b.CreatedOn,
				UpdatedOn:         b.UpdatedOn,
			}
			balances = append(balances, &balance)
		}

		account := Account{
			AccountID:   a.AccountID,
			AccountType: a.AccountType,
			Exchange:    a.Exchange,
			KeyPublic:   a.KeyPublic,
			Description: a.Description,
			CreatedOn:   a.CreatedOn,
			UpdatedOn:   a.UpdatedOn,
			Status:      a.Status,
			Balances:    balances,
		}

		accounts = append(accounts, &account)
	}

	res := &ResponseAccountListSuccess{
		Status: constRes.Success,
		Data:   AccountList{Accounts: accounts},
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:parameters PostAccount
type AccountRequest struct {
	// Required exchange for account. Specify 'binance', 'binance paper', '... paper', etc
	// in: body
	Exchange string `json:"exchange"`
	// Optional public viewable key-secret pair. This value is required for non paper accounts.
	// in: body
	KeyPublic string `json:"keyPublic"`
	// Optional init timestamp for plan RFC3339 formatted (e.g. 2018-08-26T22:49:10.168652Z). This timestamp will be used to measure initial user currency balance (valuation in user preferred currency)
	// in: body
	KeySecret string `json:"keySecret"`
	// Optional defaults to 'active' status. Valid input status is 'active', 'inactive', or 'historic'
	// in: body
	Description string `json:"description"`
	// Optional balances for a paper account
	// in: body
	Balances []*NewBalanceReq `json:"balances"`
	// Required type
	// in: body
	AccountType string `json:"type"`
}

type NewBalanceReq struct {
	// Required examples: BTC, USDT.
	// in: body
	CurrencySymbol string `json:"currencySymbol"`
	// Required amount of currency.
	// in: body
	Available float64 `json:"available"`
}

// swagger:route POST /accounts accounts PostAccount
//
// add new exchange account (protected)
//
// An exchange account is associated with balances. Paper accounts may also have optional client specified balances.
// Real exchange accounts require the public/secret key pair to populate the account balances. All balances for each
// account will be wrapped in an account object.
//
// responses:
//  200: ResponseAccountSuccess "data" will contain account data
//  400: responseError missing or incorrect params with "status": "fail"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *AccountController) HandlePostAccount(c echo.Context) error {
	//defer c.Request().Body.Close()
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	// read account deets from post body
	newAccount := new(AccountRequest)
	err := c.Bind(&newAccount)
	if err != nil {
		return fail(c, err.Error())
	}

	newBalRequests := make([]*protoBalance.NewBalanceRequest, 0)
	for _, b := range newAccount.Balances {

		br := protoBalance.NewBalanceRequest{
			CurrencySymbol: b.CurrencySymbol,
			Available:      b.Available,
		}

		newBalRequests = append(newBalRequests, &br)
	}

	newAccountRequest := protoAccount.NewAccountRequest{
		UserID:      userID,
		Exchange:    newAccount.Exchange,
		KeyPublic:   newAccount.KeyPublic,
		KeySecret:   newAccount.KeySecret,
		Description: newAccount.Description,
		AccountType: newAccount.AccountType,
		Balances:    newBalRequests,
	}

	// add plan returns nil for error
	r, _ := controller.AccountClient.AddAccount(context.Background(), &newAccountRequest)
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
	account := r.Data.Account
	balances := make([]*Balance, 0)
	for _, b := range account.Balances {
		balance := Balance{
			CurrencySymbol:    b.CurrencySymbol,
			Available:         b.Available,
			Locked:            b.Locked,
			ExchangeTotal:     b.ExchangeTotal,
			ExchangeLocked:    b.ExchangeLocked,
			ExchangeAvailable: b.ExchangeAvailable,
			CreatedOn:         b.CreatedOn,
			UpdatedOn:         b.UpdatedOn,
		}
		balances = append(balances, &balance)
	}

	res := &ResponseAccountSuccess{
		Status: constRes.Success,
		Data: &Account{
			AccountID:   account.AccountID,
			AccountType: account.AccountType,
			Exchange:    account.Exchange,
			KeyPublic:   account.KeyPublic,
			Description: account.Description,
			CreatedOn:   account.CreatedOn,
			UpdatedOn:   account.UpdatedOn,
			Status:      account.Status,
			Balances:    balances,
		},
	}

	return c.JSON(http.StatusOK, res)
}

// swagger:parameters UpdateAccountParams
type UpdateAccountRequest struct {
	// Optional public viewable key-secret pair. This value is required for non paper accounts.
	// in: body
	KeyPublic string `json:"keyPublic"`
	// Optional init timestamp for plan RFC3339 formatted (e.g. 2018-08-26T22:49:10.168652Z). This timestamp will be used to measure initial user currency balance (valuation in user preferred currency)
	// in: body
	KeySecret string `json:"keySecret"`
	// Optional defaults to 'active' status. Valid input status is 'active', 'inactive', or 'historic'
	// in: body
	Description string `json:"description"`
}

// swagger:route PUT /accounts/:accountID accounts UpdateAccountParams
//
// update a account (protected)
//
// You can update the account's keys and description. Once an account exchange has been set you cannot change
// the exchange.
//
// responses:
//  200: responseAccountSuccess "data" will contain account deets"
//  500: responseError the message will state what the internal server error was with "status": "error" "data" will contain order info with "status": "success"
func (controller *AccountController) HandleUpdateAccount(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	accountID := c.Param("accountID")

	// read strategy from post body
	updateAccount := new(UpdateAccountRequest)
	err := c.Bind(&updateAccount)
	if err != nil {
		return fail(c, err.Error())
	}

	updateAccountRequest := protoAccount.UpdateAccountRequest{
		AccountID:   accountID,
		UserID:      userID,
		KeyPublic:   updateAccount.KeyPublic,
		KeySecret:   updateAccount.KeySecret,
		Description: updateAccount.Description,
	}

	r, _ := controller.AccountClient.UpdateAccount(context.Background(), &updateAccountRequest)
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
	account := r.Data.Account
	balances := make([]*Balance, 0)
	for _, b := range account.Balances {
		balance := Balance{
			CurrencySymbol:    b.CurrencySymbol,
			Available:         b.Available,
			Locked:            b.Locked,
			ExchangeTotal:     b.ExchangeTotal,
			ExchangeLocked:    b.ExchangeLocked,
			ExchangeAvailable: b.ExchangeAvailable,
			CreatedOn:         b.CreatedOn,
			UpdatedOn:         b.UpdatedOn,
		}
		balances = append(balances, &balance)
	}

	res := &ResponseAccountSuccess{
		Status: constRes.Success,
		Data: &Account{
			AccountID:   account.AccountID,
			AccountType: account.AccountType,
			Exchange:    account.Exchange,
			KeyPublic:   account.KeyPublic,
			Description: account.Description,
			CreatedOn:   account.CreatedOn,
			UpdatedOn:   account.UpdatedOn,
			Status:      account.Status,
			Balances:    balances,
		},
	}

	return c.JSON(http.StatusOK, res)
}
