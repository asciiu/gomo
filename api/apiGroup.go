package api

import (
	"github.com/asciiu/gomo/api/handlers"

	"github.com/labstack/echo"
)

func ApiGroup(g *echo.Group) {
	// api group
	g.GET("/orders", handlers.ListOrders)
	g.POST("/orders", handlers.PostOrder)
	g.GET("/orders/:id", handlers.GetOrder)
	g.PUT("/orders/:id", handlers.UpdateOrder)
	g.DELETE("/orders/:id", handlers.DeleteOrder)
}
