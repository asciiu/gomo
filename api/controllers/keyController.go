package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	constRes "github.com/asciiu/gomo/common/constants/response"
	protoKey "github.com/asciiu/gomo/key-service/proto/key"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type KeyController struct {
	DB        *sql.DB
	KeyClient protoKey.KeyServiceClient
}

// swagger:parameters AddKey
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

// swagger:parameters UpdateKey
type UpdateKeyRequest struct {
	// Required.
	// in: body
	Description string `json:"description"`
}

// A ResponseKeySuccess will always contain a status of "successful".
// swagger:model ResponseKeySuccess
type ResponseKeySuccess struct {
	Status string       `json:"status"`
	Data   *UserKeyData `json:"data"`
}

// A ResponseKeyListSuccess will always contain a status of "successful".
// swagger:model ResponseKeyListSuccess
type ResponseKeyListSuccess struct {
	Status string   `json:"status"`
	Data   *KeyList `json:"data"`
}

type UserKeyData struct {
	Key *Key `json:"key"`
}

type KeyList struct {
	Keys []*Key `json:"keys"`
}

type Key struct {
	KeyID       string `json:"keyID"`
	Exchange    string `json:"exchange"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func NewKeyController(db *sql.DB, service micro.Service) *KeyController {
	controller := KeyController{
		DB:        db,
		KeyClient: protoKey.NewKeyServiceClient("keys", service.Client()),
	}
	return &controller
}

func (controller *KeyController) HandleGetKey(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	keyID := c.Param("keyID")

	getRequest := protoKey.GetUserKeyRequest{
		KeyID:  keyID,
		UserID: userID,
	}

	r, _ := controller.KeyClient.GetUserKey(context.Background(), &getRequest)
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
		if r.Status == constRes.Nonentity {
			return c.JSON(http.StatusNotFound, response)
		}
	}

	response := &ResponseKeySuccess{
		Status: constRes.Success,
		Data: &UserKeyData{
			Key: &Key{
				KeyID:       r.Data.Key.KeyID,
				Exchange:    r.Data.Key.Exchange,
				Key:         r.Data.Key.Key,
				Description: r.Data.Key.Description,
				Status:      r.Data.Key.Status,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (controller *KeyController) HandleListKeys(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	getRequest := protoKey.GetUserKeysRequest{
		UserID: userID,
	}

	r, e := controller.KeyClient.GetUserKeys(context.Background(), &getRequest)
	fmt.Printf("error was %+v\n", e)
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

	keys := make([]*Key, len(r.Data.Keys))
	for i, key := range r.Data.Keys {
		// api removes the secret
		keys[i] = &Key{
			KeyID:       key.KeyID,
			Exchange:    key.Exchange,
			Key:         key.Key,
			Description: key.Description,
			Status:      key.Status,
		}
	}

	response := &ResponseKeyListSuccess{
		Status: constRes.Success,
		Data: &KeyList{
			Keys: keys,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (controller *KeyController) HandlePostKey(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	addKeyRequest := KeyRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&addKeyRequest)
	if err != nil {
		response := &ResponseError{
			Status:  constRes.Fail,
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	// verify that all params are present
	if addKeyRequest.Exchange == "" || addKeyRequest.Key == "" || addKeyRequest.Secret == "" {
		response := &ResponseError{
			Status:  constRes.Fail,
			Message: "exchange, key, and secret are required!",
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	createRequest := protoKey.KeyRequest{
		UserID:      userID,
		Exchange:    addKeyRequest.Exchange,
		Key:         addKeyRequest.Key,
		Secret:      addKeyRequest.Secret,
		Description: addKeyRequest.Description,
	}

	r, _ := controller.KeyClient.AddKey(context.Background(), &createRequest)
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

	response := &ResponseKeySuccess{
		Status: constRes.Success,
		Data: &UserKeyData{
			Key: &Key{
				KeyID:       r.Data.Key.KeyID,
				Exchange:    r.Data.Key.Exchange,
				Key:         r.Data.Key.Key,
				Description: r.Data.Key.Description,
				Status:      r.Data.Key.Status,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (controller *KeyController) HandleUpdateKey(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	keyID := c.Param("keyID")

	keyRequest := UpdateKeyRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&keyRequest)
	if err != nil {
		response := &ResponseError{
			Status:  constRes.Fail,
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	// client can only update description
	updateRequest := protoKey.KeyRequest{
		KeyID:       keyID,
		UserID:      userID,
		Description: keyRequest.Description,
	}

	r, _ := controller.KeyClient.UpdateKeyDescription(context.Background(), &updateRequest)
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

	response := &ResponseKeySuccess{
		Status: constRes.Success,
		Data: &UserKeyData{
			Key: &Key{
				KeyID:       r.Data.Key.KeyID,
				Exchange:    r.Data.Key.Exchange,
				Key:         r.Data.Key.Key,
				Description: r.Data.Key.Description,
				Status:      r.Data.Key.Status,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (controller *KeyController) HandleDeleteKey(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	keyID := c.Param("keyID")

	removeRequest := protoKey.RemoveKeyRequest{
		KeyID:  keyID,
		UserID: userID,
	}

	r, _ := controller.KeyClient.RemoveKey(context.Background(), &removeRequest)
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

	response := &ResponseKeySuccess{
		Status: constRes.Success,
	}

	return c.JSON(http.StatusOK, response)
}
