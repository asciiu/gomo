package api

import (
	"net/http"

	"github.com/asciiu/gomo/api/handlers"
	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	// Login route
	e.POST("/login", handlers.Login)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
