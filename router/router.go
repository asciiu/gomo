package router

import (
	"github.com/asciiu/gomo/api"
	"github.com/asciiu/gomo/api/middlewares"
	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	// api group
	apiGroup := e.Group("/api")

	middlewares.SetMainMiddlewares(e)
	middlewares.SetApiMiddlewares(apiGroup)

	api.MainGroup(e)
	api.ApiGroup(apiGroup)

	return e
}
