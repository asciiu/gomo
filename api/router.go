package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/asciiu/gomo/api/controllers"
	repoToken "github.com/asciiu/gomo/api/db/sql"
	"github.com/asciiu/gomo/api/middlewares"
	constMessage "github.com/asciiu/gomo/common/constants/message"
	constEvt "github.com/asciiu/gomo/common/proto/events"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
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
		error := repoToken.DeleteStaleTokens(db, time.Now())
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
	e.AutoTLSManager.Cache = autocert.DirCache("/mnt/fomo/autocert")

	middlewares.SetMainMiddlewares(e)

	service := k8s.NewService(
		micro.Name("api"))

	service.Init()

	// controllers
	authController := controllers.NewAuthController(db, service)
	keyController := controllers.NewKeyController(db, service)
	balanceController := controllers.NewBalanceController(db, service)
	deviceController := controllers.NewDeviceController(db, service)
	activityController := controllers.NewActivityController(service)
	sessionController := controllers.NewSessionController(db, service)
	userController := controllers.NewUserController(db, service)
	socketController := controllers.NewWebsocketController()
	searchController := controllers.NewSearchController(db)
	planController := controllers.NewPlanController(db, service)

	// websocket ticker
	e.GET("/ticker", socketController.Connect)
	// required for health checks
	e.GET("/index.html", health)
	e.GET("/", health)

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

	// activity bulletin
	protectedApi.GET("/activity", activityController.HandleListActivity)
	protectedApi.PUT("/activity/:activityID", activityController.HandleUpdateActivity)

	// plan management endpoints
	protectedApi.GET("/plans", planController.HandleListPlans)
	protectedApi.POST("/plans", planController.HandlePostPlan)
	protectedApi.GET("/plans/:planID", planController.HandleGetPlan)
	protectedApi.PUT("/plans/:planID", planController.HandleUpdatePlan)
	protectedApi.DELETE("/plans/:planID", planController.HandleDeletePlan)

	// search endpoint
	protectedApi.GET("/search", searchController.Search)

	micro.RegisterSubscriber(constMessage.TopicAggTrade, service.Server(), func(ctx context.Context, tradeEvents *constEvt.TradeEvents) error {
		socketController.CacheEvents(tradeEvents)
		searchController.CacheEvents(tradeEvents)
		return nil
	})

	go socketController.Ticker()
	go func() {
		if err := service.Run(); err != nil {
			log.Println("nope! ", err)
		}
	}()

	return e
}
