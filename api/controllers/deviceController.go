package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	pb "github.com/asciiu/gomo/device-service/proto/device"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type DeviceController struct {
	DB     *sql.DB
	Client pb.DeviceServiceClient
}

// swagger:parameters addDevice
type PostDeviceRequest struct {
	// Required.
	// in: body
	DeviceType string `json:"deviceType"`
	// Required.
	// in: body
	DeviceToken string `json:"deviceToken"`
	// Required.
	// in: body
	ExternalDeviceId string `json:"externalDeviceId"`
}

// A ResponseSuccess will always contain a status of "successful".
// swagger:model responseDeviceSuccess
type ResponseDeviceSuccess struct {
	Status string             `json:"status"`
	Data   *pb.UserDeviceData `json:"data"`
}

func NewDeviceController(db *sql.DB) *DeviceController {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("device.client"))
	service.Init()

	controller := DeviceController{
		DB:     db,
		Client: pb.NewDeviceServiceClient("go.srv.device-service", service.Client()),
	}
	return &controller
}

// swagger:route GET /devices/:deviceId devices getDevice
//
// not implemented (protected)
//
// ...
func (controller *DeviceController) HandleGetDevice(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("ERROR!")
	}

	deviceId := c.Param("deviceId")

	return c.JSON(http.StatusOK, map[string]string{
		"status":   "not implemented",
		"deviceId": deviceId,
	})
}

// swagger:route GET /devices devices getAllDevices
//
// not implemented (protected)
//
// ...
func (controller *DeviceController) HandleListDevices(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}

// swagger:route POST /devices devices addDevice
//
// Registers a new device for a user so they may receive push notifications. (protected)
//
// responses:
//  200: responseDeviceSuccess "data" will be non null with "status": "success"
//  400: responseError missing params
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *DeviceController) HandlePostDevice(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["jti"].(string)

	addDeviceRequest := PostDeviceRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&addDeviceRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	// verify that all params are present
	if addDeviceRequest.DeviceToken == "" || addDeviceRequest.DeviceType == "" || addDeviceRequest.ExternalDeviceId == "" {
		response := &ResponseError{
			Status:  "fail",
			Message: "deviceType, deviceToken, and externalDeviceId are required!",
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	createRequest := pb.AddDeviceRequest{
		UserId:           userId,
		DeviceType:       addDeviceRequest.DeviceType,
		DeviceToken:      addDeviceRequest.DeviceToken,
		ExternalDeviceId: addDeviceRequest.ExternalDeviceId,
	}

	r, err := controller.Client.AddDevice(context.Background(), &createRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: "the device-service is not available",
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

	response := &ResponseDeviceSuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route PUT /devices/:deviceId devices updateDevice
//
// not implemented (protected)
//
// ...
func (controller *DeviceController) HandleUpdateDevice(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}

// swagger:route DELETE /devices/:deviceId devices deleteDevice
//
// not implemented (protected)
//
// ...
func (controller *DeviceController) HandleDeleteDevice(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}
