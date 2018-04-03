package controllers

import (
	"database/sql"
	"encoding/json"
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
	ExternalDeviceId string `json:"externalDeviceId"`
}

// A ResponseDeviceSuccess will always contain a status of "successful".
// swagger:model responseDeviceSuccess
type ResponseDeviceSuccess struct {
	Status string             `json:"status"`
	Data   *pb.UserDeviceData `json:"data"`
}

// A ResponseDevicesSuccess will always contain a status of "successful".
// swagger:model responseDevicesSuccess
type ResponseDevicesSuccess struct {
	Status string              `json:"status"`
	Data   *pb.UserDevicesData `json:"data"`
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
// Get a device by ID (protected)
//
// Get a user's device by the device's ID.
//
// responses:
//  200: responseDeviceSuccess "data" will contain device stuffs with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *DeviceController) HandleGetDevice(c echo.Context) error {
	deviceId := c.Param("deviceId")

	getRequest := pb.GetUserDeviceRequest{
		DeviceId: deviceId,
	}

	r, err := controller.Client.GetUserDevice(context.Background(), &getRequest)
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

// swagger:route GET /devices devices getAllDevices
//
// All registered devices. (protected)
//
// Returns a list of registered devices for logged in user.
//
// responses:
//  200: responseDevicesSuccess "data" will contain array of devices with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *DeviceController) HandleListDevices(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["jti"].(string)

	getRequest := pb.GetUserDevicesRequest{
		UserId: userId,
	}

	r, err := controller.Client.GetUserDevices(context.Background(), &getRequest)
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

	response := &ResponseDevicesSuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:route POST /devices devices addDevice
//
// Add new device. (protected)
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
// Update a registered device. (protected)
//
// Updates a user's device.
//
// responses:
//  200: responseDeviceSuccess "data" will contain updated device info with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *DeviceController) HandleUpdateDevice(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["jti"].(string)
	deviceId := c.Param("deviceId")

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
	if addDeviceRequest.DeviceToken == "" || addDeviceRequest.DeviceType == "" || addDeviceRequest.ExternalDeviceId == "" {
		response := &ResponseError{
			Status:  "fail",
			Message: "deviceType, deviceToken, and externalDeviceId are required!",
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	updateRequest := pb.UpdateDeviceRequest{
		DeviceId:         deviceId,
		UserId:           userId,
		DeviceType:       addDeviceRequest.DeviceType,
		DeviceToken:      addDeviceRequest.DeviceToken,
		ExternalDeviceId: addDeviceRequest.ExternalDeviceId,
	}

	r, err := controller.Client.UpdateDevice(context.Background(), &updateRequest)
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

// swagger:route DELETE /devices/:deviceId devices deleteDevice
//
// Removes a user's device. (protected)
//
// Removes device by ID.
//
// responses:
//  200: responseSuccess data will be null with "status": "success"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *DeviceController) HandleDeleteDevice(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["jti"].(string)
	deviceId := c.Param("deviceId")

	removeRequest := pb.RemoveDeviceRequest{
		DeviceId: deviceId,
		UserId:   userId,
	}

	r, err := controller.Client.RemoveDevice(context.Background(), &removeRequest)
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

	response := &ResponseSuccess{
		Status: "success",
	}

	return c.JSON(http.StatusOK, response)
}
