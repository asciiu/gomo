package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

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

func NewUserController(db *sql.DB) *UserController {
	service := micro.NewService(micro.Name("user.client"))
	service.Init()

	controller := UserController{
		DB:     db,
		Client: pb.NewUserServiceClient("go.micro.srv.user", service.Client()),
	}
	return &controller
}

func (controller *UserController) ChangePassword(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	paramId := c.Param("id")
	userId := claims["jti"].(string)

	if paramId != userId {
		response := &ResponseError{
			Status:  "fail",
			Message: "denied password change",
		}

		return c.JSON(http.StatusBadRequest, response)
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

	response := &ResponseSuccess{
		Status: r.Status,
	}

	return c.JSON(http.StatusOK, response)
}

func (controller *UserController) UpdateUser(c echo.Context) error {

	response := &ResponseSuccess{
		Status: "success",
	}

	return c.JSON(http.StatusOK, response)
}
