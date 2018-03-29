package controllers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
)

type UserController struct {
	DB *sql.DB
}

type UserChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (controller *UserController) ChangePassword(c echo.Context) error {

	response := &ResponseSuccess{
		Status: "success",
	}

	return c.JSON(http.StatusOK, response)
}

func (controller *UserController) UpdateUser(c echo.Context) error {

	response := &ResponseSuccess{
		Status: "success",
	}

	return c.JSON(http.StatusOK, response)
}
