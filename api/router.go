package main

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

func NewRouter(db *sql.DB) *echo.Echo {
	go cleanDatabase(db)

	e := echo.New()
	middlewares.SetMainMiddlewares(e)

	// controllers
	authController := &controllers.AuthController{DB: db}
	sessionController := &controllers.SessionController{DB: db}
	userController := &controllers.NewUserController{DB: db}
	orderController := &controllers.OrderController{DB: db}

	// api group
	openApi := e.Group("/api")

	// open endpoints here
	openApi.POST("/login", authController.HandleLogin)
	openApi.POST("/signup", authController.HandleSignup)

	protectedApi := e.Group("/api")
	// set the auth middlewares
	protectedApi.Use(authController.RefreshAccess)
	middlewares.SetApiMiddlewares(protectedApi)

	// ###########################  protected endpoints here
	protectedApi.GET("/session", sessionController.HandleSession)
	protectedApi.GET("/logout", authController.HandleLogout)

	// user manangement endpoints
	protectedApi.PUT("/users/:id/changepassword", userController.ChangePassword)
	protectedApi.PUT("/users/:id", userController.UpdateUser)

	// order management endpoints
	protectedApi.GET("/orders", orderController.ListOrders)
	protectedApi.POST("/orders", orderController.PostOrder)
	protectedApi.GET("/orders/:id", orderController.GetOrder)
	protectedApi.PUT("/orders/:id", orderController.UpdateOrder)
	protectedApi.DELETE("/orders/:id", orderController.DeleteOrder)

	// required for health checks
	e.GET("/index.html", health)

	return e
}
