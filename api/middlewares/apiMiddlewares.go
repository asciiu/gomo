package middlewares

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetApiMiddlewares(group *echo.Group) {
	group.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte("cuddlegang"),
		SigningMethod: "HS512",
	}))
}
