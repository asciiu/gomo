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
	return c.JSON(http.StatusOK, map[string]string{
		"status": "not implemented",
	})
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
