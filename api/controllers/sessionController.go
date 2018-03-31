package controllers

import (
	"database/sql"
	"net/http"

	models "github.com/asciiu/gomo/user-service/models"
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
//  200: responseSuccess data will be non null with status "success"
//  410: responseError the user-service is unreachable with status "error"
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
