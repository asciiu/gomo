package controllers

import (
	"database/sql"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type DeviceController struct {
	DB *sql.DB
}

func NewDeviceController(db *sql.DB) *DeviceController {
	controller := DeviceController{
		DB: db,
	}
	return &controller
}

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

func (controller *DeviceController) HandleListDevices(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}

func (controller *DeviceController) HandlePostDevice(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}

func (controller *DeviceController) HandleUpdateDevice(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}

func (controller *DeviceController) HandleDeleteDevice(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}
