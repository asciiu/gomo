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
	apiKeyController := controllers.NewApiKeyController(db)
	authController := controllers.NewAuthController(db)
	balanceController := controllers.NewBalanceController(db)
	deviceController := controllers.NewDeviceController(db)
	orderController := controllers.NewOrderController(db)
	sessionController := controllers.NewSessionController(db)
	userController := controllers.NewUserController(db)

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

	// balance endpoints
	protectedApi.GET("/balances", balanceController.HandleGetBalances)

	// user manangement endpoints
	protectedApi.PUT("/users/:id/changepassword", userController.HandleChangePassword)
	protectedApi.PUT("/users/:id", userController.HandleUpdateUser)

	// api key endpoints
	protectedApi.GET("/keys", apiKeyController.HandleListKeys)
	protectedApi.POST("/keys", apiKeyController.HandlePostKey)
	protectedApi.GET("/keys/:keyId", apiKeyController.HandleGetKey)
	protectedApi.PUT("/keys/:keyId", apiKeyController.HandleUpdateKey)
	protectedApi.DELETE("/keys/:keyId", apiKeyController.HandleDeleteKey)

	// device manage endpoints
	protectedApi.GET("/devices", deviceController.HandleListDevices)
	protectedApi.POST("/devices", deviceController.HandlePostDevice)
	protectedApi.GET("/devices/:deviceId", deviceController.HandleGetDevice)
	protectedApi.PUT("/devices/:deviceId", deviceController.HandleUpdateDevice)
	protectedApi.DELETE("/devices/:deviceId", deviceController.HandleDeleteDevice)

	// order management endpoints
	protectedApi.GET("/orders", orderController.HandleListOrders)
	protectedApi.POST("/orders", orderController.HandlePostOrder)
	protectedApi.GET("/orders/:orderId", orderController.HandleGetOrder)
	protectedApi.PUT("/orders/:orderId", orderController.HandleUpdateOrder)
	protectedApi.DELETE("/orders/:orderId", orderController.HandleDeleteOrder)

	// required for health checks
	e.GET("/index.html", health)

	return e
}
