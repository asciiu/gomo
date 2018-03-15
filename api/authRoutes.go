package api

import (
	"database/sql"

	"github.com/asciiu/gomo/api/controllers"
	"github.com/labstack/echo"
)

func AuthRoutes(e *echo.Echo, db *sql.DB) {

	auth := &controllers.AuthController{DB: db}

	// Login route
	e.POST("/login", auth.Login)
	e.POST("/signup", auth.Signup)
}
