package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	models "github.com/asciiu/gomo/user-service/models"
	pb "github.com/asciiu/gomo/user-service/proto/user"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type UserController struct {
	DB     *sql.DB
	Client pb.UserServiceClient
}

// swagger:parameters changePassword
type ChangePasswordRequest struct {
	// Required.
	// in: body
	OldPassword string `json:"oldPassword"`
	// Required.
	// in: body
	NewPassword string `json:"newPassword"`
}

// swagger:parameters updateUser
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

func NewUserController(db *sql.DB) *UserController {
	service := micro.NewService(micro.Name("user.client"))
	service.Init()

	controller := UserController{
		DB:     db,
		Client: pb.NewUserServiceClient("go.srv.user-service", service.Client()),
	}
	return &controller
}

// swagger:route PUT /users/:id/changepassword users changePassword
//
// change a user's password (protected)
//
// Allows an authenticated user to change their password. The url param is the user's id.
//
// responses:
//  200: responseSuccess the status will be "success" with data null.
//  400: responseError you did something wrong here with status "fail". Hopefully, the message is descriptive enough.
//  401: responseError the user Id in url param does not match with status "fail".
//  410: responseError the user-service is unreachable with status "error"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *UserController) HandleChangePassword(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	paramId := c.Param("id")
	userId := claims["jti"].(string)

	if paramId != userId {
		response := &ResponseError{
			Status:  "fail",
			Message: "unauthorized",
		}

		return c.JSON(http.StatusUnauthorized, response)
	}

	passwordRequest := new(ChangePasswordRequest)

	err := json.NewDecoder(c.Request().Body).Decode(&passwordRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	changeRequest := pb.ChangePasswordRequest{
		UserId:      userId,
		OldPassword: passwordRequest.OldPassword,
		NewPassword: passwordRequest.NewPassword,
	}

	r, err := controller.Client.ChangePassword(context.Background(), &changeRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: "change password service unavailable",
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

	response := &ResponseSuccess{
		Status: r.Status,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route PUT /users/:id users updateUser
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
	paramId := c.Param("id")
	userId := claims["jti"].(string)

	if paramId != userId {
		response := &ResponseError{
			Status:  "fail",
			Message: "unauthorized",
		}

		return c.JSON(http.StatusUnauthorized, response)
	}

	updateRequest := new(UpdateUserRequest)

	err := json.NewDecoder(c.Request().Body).Decode(&updateRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	changeRequest := pb.UpdateUserRequest{
		UserId: userId,
		First:  updateRequest.First,
		Last:   updateRequest.Last,
		Email:  updateRequest.Email,
	}

	r, err := controller.Client.UpdateUser(context.Background(), &changeRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: "update service unavailable",
		}

		return c.JSON(http.StatusGone, response)
	}

	response := &ResponseSuccess{
		Status: "success",
		Data: &UserData{
			User: &models.UserInfo{
				Id:    r.Data.User.UserId,
				First: r.Data.User.First,
				Last:  r.Data.User.Last,
				Email: r.Data.User.Email,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}
