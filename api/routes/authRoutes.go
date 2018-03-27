package routes

import (
	"database/sql"
	"net/http"

	"github.com/asciiu/gomo/api/controllers"
	"github.com/labstack/echo"
)

func health(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

func AuthRoutes(e *echo.Echo, db *sql.DB) {

	auth := &controllers.AuthController{DB: db}

	// Login route
	e.POST("/login", auth.Login)
	e.POST("/signup", auth.Signup)

	// required for health checks
	e.GET("/index.html", health)
}
