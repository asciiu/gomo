package controllers

import (
	"database/sql"
	"net/http"

	gsql "github.com/asciiu/gomo/common/db/sql"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type SessionController struct {
	DB *sql.DB
}

func (controller *SessionController) HandleSession(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["jti"].(string)

	user, err := gsql.FindUserById(controller.DB, userId)
	if err != nil {
		response := &ResponseError{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := &ResponseSuccess{
		Status: "success",
		Data:   &UserData{user.Info()},
	}

	return c.JSON(http.StatusOK, response)
}
