package routes

import (
	"database/sql"
	"net/http"

	"github.com/asciiu/gomo/api/controllers"
	"github.com/asciiu/gomo/api/middlewares"
	"github.com/labstack/echo"
)

func health(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

func New(db *sql.DB) *echo.Echo {
	e := echo.New()

	// controllers
	auth := &controllers.AuthController{DB: db}

	// api group
	protectedApi := e.Group("/api")
	protectedApi.Use(auth.RefreshAccess)

	middlewares.SetMainMiddlewares(e)
	// the protected api will require auth header
	middlewares.SetApiMiddlewares(protectedApi)

	//AuthRoutes(e.Group("/api"), db)
	// Login route
	e.POST("/login", auth.Login)
	e.POST("/signup", auth.Signup)

	SessionRoutes(protectedApi, db)
	OrderRoutes(protectedApi, db)

	// required for health checks
	e.GET("/index.html", health)

	return e
}
