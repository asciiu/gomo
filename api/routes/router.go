package routes

import (
	"database/sql"
	"net/http"

	"github.com/asciiu/gomo/api/middlewares"
	"github.com/labstack/echo"
)

func health(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

func New(db *sql.DB) *echo.Echo {
	e := echo.New()

	// api group
	protectedApi := e.Group("/api")

	middlewares.SetMainMiddlewares(e)
	middlewares.SetApiMiddlewares(protectedApi)

	AuthRoutes(e.Group("/api"), db)
	SessionRoutes(protectedApi, db)
	OrderRoutes(protectedApi, db)

	// required for health checks
	e.GET("/index.html", health)

	return e
}
