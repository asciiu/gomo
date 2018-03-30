package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	pb "github.com/asciiu/gomo/user-service/proto/user"
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
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("user.client"))
	service.Init()

	controller := UserController{
		DB:     db,
		Client: pb.NewUserServiceClient("go.micro.srv.user", service.Client()),
	}
	return &controller
}

func (controller *UserController) ChangePassword(c echo.Context) error {
	passwordRequest := ChangePasswordRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&passwordRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	request2 := pb.ChangePasswordRequest{
		UserId:      "1234",
		OldPassword: "old",
		NewPassword: "new",
	}

	r, err := controller.Client.ChangePassword(context.Background(), &request2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)

	response := &ResponseSuccess{
		Status: "success",
	}

	return c.JSON(http.StatusOK, response)
}

func (controller *UserController) UpdateUser(c echo.Context) error {

	response := &ResponseSuccess{
		Status: "success",
	}

	return c.JSON(http.StatusOK, response)
}
