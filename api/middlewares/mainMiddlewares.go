package middlewares

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetMainMiddlewares(e *echo.Echo) {

	e.Use(serverHeader)

	// this logs the server interaction
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]  ${status}  ${method}  ${host}${path} ${latency_human}` + "\n",
	}))

	// redirect http to https
	e.Pre(middleware.HTTPSRedirect())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// don't crash on exceptions
	//e.Use(middleware.Recover())
}

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "gomo/0.2.1")
		return next(c)
	}
}
