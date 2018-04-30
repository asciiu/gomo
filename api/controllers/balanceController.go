package controllers

import (
	"database/sql"
	"net/http"

	bpb "github.com/asciiu/gomo/balance-service/proto/balance"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type BalanceController struct {
	DB     *sql.DB
	Client bpb.BalanceServiceClient
}

// swagger:parameters addDevice updateDevice
// type DeviceRequest struct {
// 	// Required.
// 	// in: body
// 	DeviceType string `json:"deviceType"`
// 	// Required.
// 	// in: body
// 	DeviceToken string `json:"deviceToken"`
// 	// Required.
// 	// in: body
// 	ExternalDeviceId string `json:"externalDeviceId"`
// }

// A ResponseDeviceSuccess will always contain a status of "successful".
// swagger:model responseDeviceSuccess
// type ResponseDeviceSuccess struct {
// 	Status string             `json:"status"`
// 	Data   *pb.UserDeviceData `json:"data"`
// }

// A ResponseBalancesSuccess will always contain a status of "successful".
// swagger:model responseBalancesSuccess
type ResponseBalancesSuccess struct {
	Status string               `json:"status"`
	Data   *bpb.AccountBalances `json:"data"`
}

func NewBalanceController(db *sql.DB) *BalanceController {
	service := micro.NewService(micro.Name("balance.client"))
	service.Init()

	controller := BalanceController{
		DB:     db,
		Client: bpb.NewBalanceServiceClient("go.micro.srv.balance", service.Client()),
	}
	return &controller
}

// swagger:route GET /balances balances getAllBalances
//
// get all balances (protected)
//
// Returns all balances for user.
//
// responses:
//  200: responseBalancesSuccess "data" will contain array of balances with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *BalanceController) HandleGetBalances(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["jti"].(string)

	getRequest := bpb.GetUserBalancesRequest{
		UserId: userId,
	}

	r, err := controller.Client.GetUserBalances(context.Background(), &getRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: err.Error(),
		}

		return c.JSON(http.StatusGone, response)
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

	response := &ResponseBalancesSuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}
