package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	constRes "github.com/asciiu/gomo/common/constants/response"
	user "github.com/asciiu/gomo/user-service/models"
	protoUser "github.com/asciiu/gomo/user-service/proto/user"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type UserController struct {
	DB         *sql.DB
	UserClient protoUser.UserServiceClient
}

// swagger:parameters ChangePassword
type ChangePasswordRequest struct {
	// Required.
	// in: body
	OldPassword string `json:"oldPassword"`
	// Required.
	// in: body
	NewPassword string `json:"newPassword"`
}

// swagger:parameters UpdateUser
type UpdateUserRequest struct {
	// Optional.
	// in: body
	First string `json:"first"`
	// Optional.
	// in: body
	Last string `json:"last"`
	// Optional. Note: we need to validate these!
	// in: body
	Email string `json:"email"`
}

func NewUserController(db *sql.DB, service micro.Service) *UserController {
	controller := UserController{
		DB:         db,
		UserClient: protoUser.NewUserServiceClient("users", service.Client()),
	}
	return &controller
}

// swagger:route PUT /users/:id/changepassword users ChangePassword
//
// change a user's password (protected)
//
// Allows an authenticated user to change their password. The url param is the user's id.
//
// responses:
//  200: responseSuccess the status will be "success" with data null.
//  400: responseError you did something wrong here with status "fail". Hopefully, the message is descriptive enough.
//  401: responseError the user ID in url param does not match with status "fail".
//  410: responseError the user-service is unreachable with status "error"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *UserController) HandleChangePassword(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	paramID := c.Param("id")
	userID := claims["jti"].(string)

	if paramID != userID {
		response := &ResponseError{
			Status:  constRes.Fail,
			Message: "unauthorized",
		}

		return c.JSON(http.StatusUnauthorized, response)
	}

	passwordRequest := new(ChangePasswordRequest)

	err := json.NewDecoder(c.Request().Body).Decode(&passwordRequest)
	if err != nil {
		response := &ResponseError{
			Status:  constRes.Fail,
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	changeRequest := protoUser.ChangePasswordRequest{
		UserID:      userID,
		OldPassword: passwordRequest.OldPassword,
		NewPassword: passwordRequest.NewPassword,
	}

	r, err := controller.UserClient.ChangePassword(context.Background(), &changeRequest)
	if err != nil {
		response := &ResponseError{
			Status:  constRes.Error,
			Message: "change password service unavailable",
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

	response := &ResponseSuccess{
		Status: r.Status,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route PUT /users/:id users UpdateUser
//
// updates user info (protected)
//
// You can change the user's first, last, or email. Note we need to implement a secure method of
// verifing the user's new email. This has yet to be implemented.
//
// responses:
//  200: responseSuccess "data" will contain updated user data with "status": "success"
//  400: responseError message in badrequest should be descriptive with "status": "fail"
//  401: responseError unauthorized user because of incorrect url param with "status": "fail"
//  410: responseError the user-service is unreachable with status "error"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *UserController) HandleUpdateUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	paramID := c.Param("id")
	userID := claims["jti"].(string)

	if paramID != userID {
		response := &ResponseError{
			Status:  constRes.Fail,
			Message: "unauthorized",
		}

		return c.JSON(http.StatusUnauthorized, response)
	}

	updateRequest := new(UpdateUserRequest)

	err := json.NewDecoder(c.Request().Body).Decode(&updateRequest)
	if err != nil {
		response := &ResponseError{
			Status:  constRes.Fail,
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	changeRequest := protoUser.UpdateUserRequest{
		UserID: userID,
		First:  updateRequest.First,
		Last:   updateRequest.Last,
		Email:  updateRequest.Email,
	}

	r, err := controller.UserClient.UpdateUser(context.Background(), &changeRequest)
	if err != nil {
		response := &ResponseError{
			Status:  constRes.Error,
			Message: "update service unavailable",
		}

		return c.JSON(http.StatusGone, response)
	}

	response := &ResponseSuccess{
		Status: constRes.Success,
		Data: &UserData{
			User: &user.UserInfo{
				UserID: r.Data.User.UserID,
				First:  r.Data.User.First,
				Last:   r.Data.User.Last,
				Email:  r.Data.User.Email,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}
