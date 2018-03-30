package controllers

import (
	"database/sql"
	"net/http"

	"github.com/asciiu/gomo/user-service/models"
	pb "github.com/asciiu/gomo/user-service/proto/user"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type SessionController struct {
	DB     *sql.DB
	Client pb.UserServiceClient
}

func NewSessionController(db *sql.DB) *SessionController {
	service := micro.NewService(micro.Name("user.client"))
	service.Init()

	controller := SessionController{
		DB:     db,
		Client: pb.NewUserServiceClient("go.micro.srv.user", service.Client()),
	}
	return &controller
}

func (controller *SessionController) HandleSession(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["jti"].(string)

	getRequest := pb.GetUserInfoRequest{
		UserId: userId,
	}
	r, err := controller.Client.GetUserInfo(context.Background(), &getRequest)
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
