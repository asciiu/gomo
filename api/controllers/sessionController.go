package controllers

import (
	"database/sql"
	"net/http"

	keys "github.com/asciiu/gomo/key-service/proto/key"
	users "github.com/asciiu/gomo/user-service/proto/user"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type SessionController struct {
	DB    *sql.DB
	Users users.UserServiceClient
	Keys  keys.KeyServiceClient
}

type UserMetaData struct {
	UserMeta *UserMeta `json:"user"`
}

type UserMeta struct {
	UserID string     `json:"userID"`
	First  string     `json:"first"`
	Last   string     `json:"last"`
	Email  string     `json:"email"`
	Keys   []*KeyMeta `json:"keys"`
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
		DB:    db,
		Users: users.NewUserServiceClient("users", service.Client()),
		Keys:  keys.NewKeyServiceClient("keys", service.Client()),
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

	getRequest := users.GetUserInfoRequest{
		UserID: userID,
	}
	r, _ := controller.Users.GetUserInfo(context.Background(), &getRequest)
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

	getKeysRequest := keys.GetUserKeysRequest{
		UserID: userID,
	}
	r2, _ := controller.Keys.GetUserKeys(context.Background(), &getKeysRequest)
	lekeys := make([]*KeyMeta, 0)

	for _, k := range r2.Data.Keys {
		lekeys = append(lekeys,
			&KeyMeta{
				Exchange:    k.Exchange,
				Status:      k.Status,
				Description: k.Description,
				KeyID:       k.KeyID})
	}

	response := &ResponseSessionSuccess{
		Status: "success",
		Data: &UserMetaData{
			UserMeta: &UserMeta{
				UserID: r.Data.User.UserID,
				First:  r.Data.User.First,
				Last:   r.Data.User.Last,
				Email:  r.Data.User.Email,
				Keys:   lekeys}}}

	return c.JSON(http.StatusOK, response)
}
