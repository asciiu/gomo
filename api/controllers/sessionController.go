package controllers

import (
	"database/sql"
	"net/http"

	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	constRes "github.com/asciiu/gomo/common/constants/response"
	protoUser "github.com/asciiu/gomo/user-service/proto/user"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type SessionController struct {
	DB            *sql.DB
	UserClient    protoUser.UserServiceClient
	AccountClient protoAccount.AccountServiceClient
}

type UserMetaData struct {
	UserMeta *UserMeta `json:"user"`
}

type UserMeta struct {
	UserID   string     `json:"userID"`
	First    string     `json:"first"`
	Last     string     `json:"last"`
	Email    string     `json:"email"`
	Accounts []*Account `json:"accounts"`
}

type KeyMeta struct {
	KeyID       string `json:"keyID"`
	Exchange    string `json:"exchange"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// A ResponseSessionSuccess will always contain a status of "successful".
// swagger:model ResponseSessionSuccess
type ResponseSessionSuccess struct {
	Status string        `json:"status"`
	Data   *UserMetaData `json:"data"`
}

func NewSessionController(db *sql.DB, service micro.Service) *SessionController {
	controller := SessionController{
		DB:            db,
		UserClient:    protoUser.NewUserServiceClient("users", service.Client()),
		AccountClient: protoAccount.NewAccountServiceClient("accounts", service.Client()),
	}
	return &controller
}

// swagger:route GET /session session sessionBegin
//
// create a new session for a user (protected)
//
// Creates a new session for an authenticated user. The session data will eventually contain
// whatever info you need to begin a new session. At the moment the response data mirrors
// login data. This endpoint depends on the user-service. If the user-service
// is unreachable, a 410 with a status of "error" will be returned.
//
// responses:
//  200: ResponseSessionSuccess data will be non null with status "success"
//  410: responseError the user-service is unreachable with status "error"
func (controller *SessionController) HandleSession(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	getRequest := protoUser.GetUserInfoRequest{
		UserID: userID,
	}
	r, _ := controller.UserClient.GetUserInfo(context.Background(), &getRequest)
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

	requestAccounts := protoAccount.AccountsRequest{UserID: userID}
	responseAccounts, _ := controller.AccountClient.ResyncAccounts(context.Background(), &requestAccounts)
	accounts := make([]*Account, 0)
	for _, a := range responseAccounts.Data.Accounts {

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

	response := &ResponseSessionSuccess{
		Status: constRes.Success,
		Data: &UserMetaData{
			UserMeta: &UserMeta{
				UserID:   r.Data.User.UserID,
				First:    r.Data.User.First,
				Last:     r.Data.User.Last,
				Email:    r.Data.User.Email,
				Accounts: accounts,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}
