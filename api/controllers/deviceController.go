package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	devices "github.com/asciiu/gomo/device-service/proto/device"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type DeviceController struct {
	DB      *sql.DB
	Devices devices.DeviceServiceClient
}

// swagger:parameters addDevice updateDevice
type DeviceRequest struct {
	// Required.
	// in: body
	DeviceType string `json:"deviceType"`
	// Required.
	// in: body
	DeviceToken string `json:"deviceToken"`
	// Required.
	// in: body
	ExternalDeviceID string `json:"externalDeviceID"`
}

// A ResponseDeviceSuccess will always contain a status of "successful".
// swagger:model responseDeviceSuccess
type ResponseDeviceSuccess struct {
	Status string          `json:"status"`
	Data   *UserDeviceData `json:"data"`
}

// A ResponseDevicesSuccess will always contain a status of "successful".
// swagger:model responseDevicesSuccess
type ResponseDevicesSuccess struct {
	Status string           `json:"status"`
	Data   *UserDevicesData `json:"data"`
}

type UserDeviceData struct {
	Device *ApiDevice `json:"device"`
}

type UserDevicesData struct {
	Devices []*ApiDevice `json:"devices"`
}

type ApiDevice struct {
	DeviceID         string `json:"deviceID"`
	ExternalDeviceID string `json:"externalDeviceID"`
	DeviceType       string `json:"deviceType"`
	DeviceToken      string `json:"deviceToken"`
}

func NewDeviceController(db *sql.DB) *DeviceController {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("device.client"))
	service.Init()

	controller := DeviceController{
		DB:      db,
		Devices: devices.NewDeviceServiceClient("devices", service.Client()),
	}
	return &controller
}

// swagger:route GET /devices/:deviceID devices getDevice
//
// get a device by ID (protected)
//
// Get a user's device by the device's ID.
//
// responses:
//  200: responseDeviceSuccess "data" will contain device stuffs with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *DeviceController) HandleGetDevice(c echo.Context) error {
	deviceID := c.Param("deviceID")

	getRequest := devices.GetUserDeviceRequest{
		DeviceID: deviceID,
	}

	r, _ := controller.Devices.GetUserDevice(context.Background(), &getRequest)
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
		Data: &UserDeviceData{
			Device: &ApiDevice{
				DeviceID:         r.Data.Device.DeviceID,
				ExternalDeviceID: r.Data.Device.ExternalDeviceID,
				DeviceType:       r.Data.Device.DeviceType,
				DeviceToken:      r.Data.Device.DeviceToken,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route GET /devices devices getAllDevices
//
// all registered devices (protected)
//
// Returns a list of registered devices for logged in user.
//
// responses:
//  200: responseDevicesSuccess "data" will contain array of devices with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *DeviceController) HandleListDevices(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	getRequest := devices.GetUserDevicesRequest{
		UserID: userID,
	}

	r, _ := controller.Devices.GetUserDevices(context.Background(), &getRequest)
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

	data := make([]*ApiDevice, len(r.Data.Devices))
	for i, device := range r.Data.Devices {
		data[i] = &ApiDevice{
			DeviceID:         device.DeviceID,
			ExternalDeviceID: device.ExternalDeviceID,
			DeviceType:       device.DeviceType,
			DeviceToken:      device.DeviceToken,
		}
	}

	response := &ResponseDevicesSuccess{
		Status: "success",
		Data: &UserDevicesData{
			Devices: data,
		},
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route POST /devices devices addDevice
//
// add new device (protected)
//
// Registers a new device for a user so they may receive push notifications.
//
// responses:
//  200: responseDeviceSuccess "data" will be non null with "status": "success"
//  400: responseError missing params
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *DeviceController) HandlePostDevice(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	addDeviceRequest := DeviceRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&addDeviceRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	// verify that all params are present
	if addDeviceRequest.DeviceToken == "" || addDeviceRequest.DeviceType == "" || addDeviceRequest.ExternalDeviceID == "" {
		response := &ResponseError{
			Status:  "fail",
			Message: "deviceType, deviceToken, and externalDeviceID are required!",
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	createRequest := devices.AddDeviceRequest{
		UserID:           userID,
		DeviceType:       addDeviceRequest.DeviceType,
		DeviceToken:      addDeviceRequest.DeviceToken,
		ExternalDeviceID: addDeviceRequest.ExternalDeviceID,
	}

	r, _ := controller.Devices.AddDevice(context.Background(), &createRequest)
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
		Data: &UserDeviceData{
			Device: &ApiDevice{
				DeviceID:         r.Data.Device.DeviceID,
				ExternalDeviceID: r.Data.Device.ExternalDeviceID,
				DeviceType:       r.Data.Device.DeviceType,
				DeviceToken:      r.Data.Device.DeviceToken,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route PUT /devices/:deviceID devices updateDevice
//
// update a registered device (protected)
//
// Updates a user's device.
//
// responses:
//  200: responseDeviceSuccess "data" will contain updated device info with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *DeviceController) HandleUpdateDevice(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	deviceID := c.Param("deviceID")

	addDeviceRequest := DeviceRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&addDeviceRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	// verify that all params are present
	if addDeviceRequest.DeviceToken == "" || addDeviceRequest.DeviceType == "" || addDeviceRequest.ExternalDeviceID == "" {
		response := &ResponseError{
			Status:  "fail",
			Message: "deviceType, deviceToken, and externalDeviceID are required!",
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	updateRequest := devices.UpdateDeviceRequest{
		DeviceID:         deviceID,
		UserID:           userID,
		DeviceType:       addDeviceRequest.DeviceType,
		DeviceToken:      addDeviceRequest.DeviceToken,
		ExternalDeviceID: addDeviceRequest.ExternalDeviceID,
	}

	r, _ := controller.Devices.UpdateDevice(context.Background(), &updateRequest)
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
		Data: &UserDeviceData{
			Device: &ApiDevice{
				DeviceID:         r.Data.Device.DeviceID,
				ExternalDeviceID: r.Data.Device.ExternalDeviceID,
				DeviceType:       r.Data.Device.DeviceType,
				DeviceToken:      r.Data.Device.DeviceToken,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route DELETE /devices/:deviceID devices deleteDevice
//
// removes a user's device (protected)
//
// Removes device by ID.
//
// responses:
//  200: responseSuccess data will be null with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *DeviceController) HandleDeleteDevice(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)
	deviceID := c.Param("deviceID")

	removeRequest := devices.RemoveDeviceRequest{
		DeviceID: deviceID,
		UserID:   userID,
	}

	r, _ := controller.Devices.RemoveDevice(context.Background(), &removeRequest)
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

	response := &ResponseSuccess{
		Status: "success",
	}

	return c.JSON(http.StatusOK, response)
}
