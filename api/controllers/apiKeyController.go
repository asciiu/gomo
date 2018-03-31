package controllers

import (
	"database/sql"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type ApiKeyController struct {
	DB *sql.DB
}

func NewApiKeyController(db *sql.DB) *ApiKeyController {
	controller := ApiKeyController{
		DB: db,
	}
	return &controller
}

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

func (controller *ApiKeyController) HandleListKeys(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}

func (controller *ApiKeyController) HandlePostKey(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}

func (controller *ApiKeyController) HandleUpdateKey(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}

func (controller *ApiKeyController) HandleDeleteKey(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
}
