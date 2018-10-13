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
	e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("admin.stage.fomo.exchange")
	e.AutoTLSManager.Cache = autocert.DirCache("/mnt/fomo/autocert")

	middlewares.SetMainMiddlewares(e)

	service := k8s.NewService(
		micro.Name("admin"))

	service.Init()

	// controllers
	authController := controllers.NewAuthController(db, service)
	sessionController := controllers.NewSessionController(db, service)
	socketController := controllers.NewWebsocketController()

	// websocket ticker
	e.GET("/ticker", socketController.Connect)
	// required for health checks
	e.GET("/index.html", health)
	e.GET("/", health)

	// api group
	openApi := e.Group("/admin")

	// open endpoints here
	openApi.POST("/login", authController.HandleLogin)

	protectedApi := e.Group("/admin")
	// set the auth middlewares
	protectedApi.Use(authController.RefreshAccess)
	middlewares.SetApiMiddlewares(protectedApi)

	// ###########################  protected endpoints here
	protectedApi.GET("/session", sessionController.HandleSession)
	protectedApi.GET("/logout", authController.HandleLogout)

	micro.RegisterSubscriber(constMessage.TopicAggTrade, service.Server(), func(ctx context.Context, tradeEvents *constEvt.TradeEvents) error {
		socketController.CacheEvents(tradeEvents)
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
