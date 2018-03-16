package routes

import (
	"database/sql"

	"github.com/asciiu/gomo/api/controllers"
	"github.com/labstack/echo"
)

func OrderRoutes(g *echo.Group, db *sql.DB) {
	// api group
	g.GET("/orders", controllers.ListOrders)
	g.POST("/orders", controllers.PostOrder)
	g.GET("/orders/:id", controllers.GetOrder)
	g.PUT("/orders/:id", controllers.UpdateOrder)
	g.DELETE("/orders/:id", controllers.DeleteOrder)
}
