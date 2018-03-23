package middlewares

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetApiMiddlewares(group *echo.Group) {
	group.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(os.Getenv("GOMO_JWT")),
		SigningMethod: "HS512",
	}))
}
