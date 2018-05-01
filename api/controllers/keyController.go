package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	kp "github.com/asciiu/gomo/key-service/proto/key"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type KeyController struct {
	DB     *sql.DB
	Client kp.KeyServiceClient
}

// swagger:parameters postKey
type KeyRequest struct {
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

// swagger:parameters updateKey
type UpdateKeyRequest struct {
	// Required.
	// in: body
	Description string `json:"description"`
}

// A ResponseKeySuccess will always contain a status of "successful".
// swagger:model responseKeySuccess
type ResponseKeySuccess struct {
	Status string          `json:"status"`
	Data   *kp.UserKeyData `json:"data"`
}

// A ResponseKeysSuccess will always contain a status of "successful".
// swagger:model responseKeysSuccess
type ResponseKeysSuccess struct {
	Status string           `json:"status"`
	Data   *kp.UserKeysData `json:"data"`
}

func NewKeyController(db *sql.DB) *KeyController {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("key.client"))
	service.Init()

	controller := KeyController{
		DB:     db,
		Client: kp.NewKeyServiceClient("go.srv.key-service", service.Client()),
	}
	return &controller
}

// swagger:route GET /keys/:keyID keys getKey
//
// get a key (protected)
//
// Gets a user's key by the key ID. The secret will not be returned in the response data.
//
// responses:
//  200: responseKeySuccess "data" will contain key stuffs with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *KeyController) HandleGetKey(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	keyID := c.Param("keyID")

	getRequest := kp.GetUserKeyRequest{
		KeyID:  keyID,
		UserID: userID,
	}

	r, err := controller.Client.GetUserKey(context.Background(), &getRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, response)
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
		Data: &kp.UserKeyData{
			Key: &kp.Key{
				KeyID:       r.Data.Key.KeyID,
				UserID:      r.Data.Key.UserID,
				Exchange:    r.Data.Key.Exchange,
				Key:         r.Data.Key.Key,
				Description: r.Data.Key.Description,
				Status:      r.Data.Key.Status,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route GET /keys keys getAllKey
//
// get all user keys (protected)
//
// Get all the user keys for this user. The api secrets will not be returned in the response data.
//
// responses:
//  200: responseKeysSuccess "data" will contain a list of key info with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *KeyController) HandleListKeys(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	getRequest := kp.GetUserKeysRequest{
		UserID: userID,
	}

	r, err := controller.Client.GetUserKeys(context.Background(), &getRequest)
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

	data := make([]*kp.Key, len(r.Data.Keys))
	for i, key := range data {
		// api removes the secret
		data[i] = &kp.Key{
			KeyID:       key.KeyID,
			UserID:      key.UserID,
			Exchange:    key.Exchange,
			Key:         key.Key,
			Description: key.Description,
			Status:      key.Status,
		}
	}

	response := &ResponseKeysSuccess{
		Status: "success",
		Data: &kp.UserKeysData{
			Keys: data,
		},
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route POST /keys keys postKey
//
// add an api key (protected)
//
// Associate a new exchange api key to a user's account. Secrets will not be returned in response data.
//
// responses:
//  200: responseKeySuccess "data" will contain key info with "status": "success"
//  400: responseError missing params with "status": "fail"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *KeyController) HandlePostKey(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	addKeyRequest := KeyRequest{}

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

	createRequest := kp.KeyRequest{
		UserID:      userID,
		Exchange:    addKeyRequest.Exchange,
		Key:         addKeyRequest.Key,
		Secret:      addKeyRequest.Secret,
		Description: addKeyRequest.Description,
	}

	r, err := controller.Client.AddKey(context.Background(), &createRequest)
	if err != nil {
		fmt.Println(err)
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

	response := &ResponseKeySuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route PUT /keys/:keyID keys updateKey
//
// update a user api key (protected)
//
// The user can only update the description of an added key. The secret will not be returned.
//
// responses:
//  200: responseKeySuccess "data" will contain key info with "status": "success"
//  400: responseError missing params with "status": "fail"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *KeyController) HandleUpdateKey(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	keyID := c.Param("keyID")

	keyRequest := UpdateKeyRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&keyRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	// client can only update description
	updateRequest := kp.KeyRequest{
		KeyID:       keyID,
		UserID:      userID,
		Description: keyRequest.Description,
	}

	r, err := controller.Client.UpdateKeyDescription(context.Background(), &updateRequest)
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

	response := &ResponseKeySuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route DELETE /keys/:keyID keys deleteKey
//
// remove user api key (protected)
//
// This will remove the api key from the system.
//
// responses:
//  200: responseKeySuccess data will be null with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *KeyController) HandleDeleteKey(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	keyID := c.Param("keyID")

	removeRequest := kp.RemoveKeyRequest{
		KeyID:  keyID,
		UserID: userID,
	}

	r, err := controller.Client.RemoveKey(context.Background(), &removeRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: "the key-service is not available",
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
	}

	return c.JSON(http.StatusOK, response)
}
