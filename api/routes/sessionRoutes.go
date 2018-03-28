package routes

import (
	"database/sql"

	"github.com/asciiu/gomo/api/controllers"
	"github.com/labstack/echo"
)

func SessionRoutes(e *echo.Group, db *sql.DB) {

	auth := &controllers.SessionController{DB: db}

	// this needs to be protected
	e.GET("/session", auth.Session)
}
