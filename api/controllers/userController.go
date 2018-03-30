package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/asciiu/gomo/user-service/models"
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

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UpdateUserRequest struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Email string `json:"email"`
}

func NewUserController(db *sql.DB) *UserController {
	service := micro.NewService(micro.Name("user.client"))
	service.Init()

	controller := UserController{
		DB:     db,
		Client: pb.NewUserServiceClient("go.micro.srv.user", service.Client()),
	}
	return &controller
}

// Handle password change.
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
		return c.JSON(http.StatusBadRequest, response)
	}

	response := &ResponseSuccess{
		Status: r.Status,
	}

	return c.JSON(http.StatusOK, response)
}

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
			&models.UserInfo{
				Id:    r.Data.User.UserId,
				First: r.Data.User.First,
				Last:  r.Data.User.Last,
				Email: r.Data.User.Email,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}
