package routes

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/asciiu/gomo/api/controllers"
	asql "github.com/asciiu/gomo/api/db/sql"
	"github.com/asciiu/gomo/api/middlewares"
	"github.com/labstack/echo"
)

// clean up stage refresh tokens in DB every 30 minutes
const cleanUpInterval = 30 * time.Minute

// send 200 ok to ping requests
func health(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

// routine to clean up refresh tokens in DB
func cleanDatabase(db *sql.DB) {
	for {
		time.Sleep(cleanUpInterval)
		error := asql.DeleteStaleTokens(db, time.Now())
		if error != nil {
			log.Fatal(error)
		}
	}
}

func New(db *sql.DB) *echo.Echo {
	go cleanDatabase(db)

	e := echo.New()

	// controllers
	auth := &controllers.AuthController{DB: db}

	// api group
	openApi := e.Group("/api")

	protectedApi := e.Group("/api")
	protectedApi.Use(auth.RefreshAccess)

	middlewares.SetMainMiddlewares(e)
	// the protected api will require auth header
	middlewares.SetApiMiddlewares(protectedApi)

	//AuthRoutes(e.Group("/api"), db)
	// Login route
	openApi.POST("/login", auth.Login)
	openApi.POST("/signup", auth.Signup)

	SessionRoutes(protectedApi, db)
	OrderRoutes(protectedApi, db)

	// required for health checks
	e.GET("/index.html", health)

	return e
}
