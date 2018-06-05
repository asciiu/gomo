package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/asciiu/gomo/api/controllers"
	asql "github.com/asciiu/gomo/api/db/sql"
	"github.com/asciiu/gomo/api/middlewares"
	msg "github.com/asciiu/gomo/common/messages"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	k8s "github.com/micro/kubernetes/go/micro"
	"golang.org/x/crypto/acme/autocert"
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
	e.AutoTLSManager.Prompt = autocert.AcceptTOS
	e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("stage.fomo.exchange")
	e.AutoTLSManager.Cache = autocert.DirCache("/var/tmp/.cache")

	middlewares.SetMainMiddlewares(e)

	// controllers
	keyController := controllers.NewKeyController(db)
	authController := controllers.NewAuthController(db)
	balanceController := controllers.NewBalanceController(db)
	deviceController := controllers.NewDeviceController(db)
	notificationController := controllers.NewNotificationController()
	orderController := controllers.NewOrderController(db)
	sessionController := controllers.NewSessionController(db)
	userController := controllers.NewUserController(db)
	socketController := controllers.NewWebsocketController()
	searchController := controllers.NewSearchController(db)

	// websocket ticker
	e.GET("/ticker", socketController.Connect)

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
	protectedApi.PUT("/users/:userID/changepassword", userController.HandleChangePassword)
	protectedApi.PUT("/users/:userID", userController.HandleUpdateUser)

	// api key endpoints
	protectedApi.GET("/keys", keyController.HandleListKeys)
	protectedApi.POST("/keys", keyController.HandlePostKey)
	protectedApi.GET("/keys/:keyID", keyController.HandleGetKey)
	protectedApi.PUT("/keys/:keyID", keyController.HandleUpdateKey)
	protectedApi.DELETE("/keys/:keyID", keyController.HandleDeleteKey)

	// device manage endpoints
	protectedApi.GET("/devices", deviceController.HandleListDevices)
	protectedApi.POST("/devices", deviceController.HandlePostDevice)
	protectedApi.GET("/devices/:deviceID", deviceController.HandleGetDevice)
	protectedApi.PUT("/devices/:deviceID", deviceController.HandleUpdateDevice)
	protectedApi.DELETE("/devices/:deviceID", deviceController.HandleDeleteDevice)

	// notifications
	protectedApi.GET("/notifications", notificationController.HandleListNotifications)

	// order management endpoints
	protectedApi.GET("/orders", orderController.HandleListOrders)
	protectedApi.POST("/orders", orderController.HandlePostOrder)
	protectedApi.GET("/orders/:orderID", orderController.HandleGetOrder)
	protectedApi.PUT("/orders/:orderID", orderController.HandleUpdateOrder)
	protectedApi.DELETE("/orders/:orderID", orderController.HandleDeleteOrder)

	// search endpoint
	protectedApi.GET("/search", searchController.Search)

	// required for health checks
	e.GET("/index.html", health)
	e.GET("/", health)

	service := k8s.NewService(
		micro.Name("micro.fomo.api"),
	)
	service.Init()
	micro.RegisterSubscriber(msg.TopicAggTrade, service.Server(), socketController.ProcessEvent, server.SubscriberQueue("queue.pubsub"))
	micro.RegisterSubscriber(msg.TopicAggTrade, service.Server(), searchController.ProcessEvent, server.SubscriberQueue("queue.pubsub"))

	go socketController.Ticker()
	go service.Run()

	return e
}
