package routes

import (
	"database/sql"

	"github.com/asciiu/gomo/api/middlewares"
	"github.com/labstack/echo"
)

func New(db *sql.DB) *echo.Echo {
	e := echo.New()

	// api group
	apiGroup := e.Group("/api")

	middlewares.SetMainMiddlewares(e)
	middlewares.SetApiMiddlewares(apiGroup)

	AuthRoutes(e, db)
	OrderRoutes(apiGroup, db)

	return e
}
