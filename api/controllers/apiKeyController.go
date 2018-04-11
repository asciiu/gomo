package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	keyProto "github.com/asciiu/gomo/apikey-service/proto/apikey"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type ApiKeyController struct {
	DB     *sql.DB
	Client keyProto.ApiKeyServiceClient
}

// swagger:parameters postKey
type ApiKeyRequest struct {
	// Required.
	// in: body
	Exchange string `json:"exchange"`
	// Required.
	// in: body
	Key string `json:"key"`
	// Required.
	// in: body
	Secret string `json:"secret"`
	// Optional.
	// in: body
	Description string `json:"description"`
}

// A ResponseApiKeySuccess will always contain a status of "successful".
// swagger:model responseKeySuccess
type ResponseKeySuccess struct {
	Status string                   `json:"status"`
	Data   *keyProto.UserApiKeyData `json:"data"`
}

// A ResponseApiKeysSuccess will always contain a status of "successful".
// swagger:model responseKeysSuccess
type ResponseKeysSuccess struct {
	Status string                    `json:"status"`
	Data   *keyProto.UserApiKeysData `json:"data"`
}

func NewApiKeyController(db *sql.DB) *ApiKeyController {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("apikey.client"))
	service.Init()

	controller := ApiKeyController{
		DB:     db,
		Client: keyProto.NewApiKeyServiceClient("go.srv.apikey-service", service.Client()),
	}
	return &controller
}

// swagger:route GET /keys/:keyId keys getKey
//
// not implemented (protected)
//
// ...
func (controller *ApiKeyController) HandleGetKey(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("ERROR!")
	}

	keyId := c.Param("keyId")

	return c.JSON(http.StatusOK, map[string]string{
		"status":   "not implemented",
		"deviceId": keyId,
	})
}

// swagger:route GET /keys keys getAllKey
//
// not implemented (protected)
//
// ...
func (controller *ApiKeyController) HandleListKeys(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}

// swagger:route POST /keys keys postKey
//
// not implemented (protected)
//
// ..
func (controller *ApiKeyController) HandlePostKey(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["jti"].(string)
	addKeyRequest := ApiKeyRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&addKeyRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	// verify that all params are present
	if addKeyRequest.Exchange == "" || addKeyRequest.Key == "" || addKeyRequest.Secret == "" {
		response := &ResponseError{
			Status:  "fail",
			Message: "exchange, key, and secret are required!",
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	createRequest := keyProto.ApiKeyRequest{
		UserId:      userId,
		Exchange:    addKeyRequest.Exchange,
		Key:         addKeyRequest.Key,
		Secret:      addKeyRequest.Secret,
		Description: addKeyRequest.Description,
	}

	r, err := controller.Client.AddApiKey(context.Background(), &createRequest)
	if err != nil {
		fmt.Println(err)
		response := &ResponseError{
			Status:  "error",
			Message: "the apikey-service is not available",
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

	response := &ResponseKeySuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route PUT /keys/:keyId keys updateKey
//
// not implemented (protected)
//
// ..
func (controller *ApiKeyController) HandleUpdateKey(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}

// swagger:route DELETE /keys/:keyId keys deleteKey
//
// not implemented (protected)
//
// ...
func (controller *ApiKeyController) HandleDeleteKey(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}
