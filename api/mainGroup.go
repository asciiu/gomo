package api

import (
	"database/sql"

	"github.com/asciiu/gomo/api/handlers"
	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo, db *sql.DB) {

	main := &handlers.MainRoutes{DB: db}

	// Login route
	e.POST("/login", main.Login)
	e.POST("/signup", main.Signup)
}
