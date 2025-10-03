package router

import (
	"expenses-api/internal/infrastructure/container"
	"log"

	"github.com/gin-gonic/gin"
)

func mapURLs(router *gin.Engine) {
	// Initialize dependency injection container
	c, err := container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	frontendUrls(router, c)
}

// frontendUrls define las rutas específicas para el frontend Angular usando handlers
func frontendUrls(router *gin.Engine, c *container.Container) {
	// Grupo de rutas API
	api := router.Group("/api")
	{
		// Resumen mensual
		api.GET("/summary/:month", c.SummaryHandler.GetMonthlySummary)

		// Configuración
		api.GET("/config/income", c.ConfigHandler.GetIncome)
		api.PUT("/config/income", c.ConfigHandler.UpdateIncome)
		api.GET("/config/pockets", c.ConfigHandler.GetPockets)
		api.POST("/config/pockets", c.ConfigHandler.CreatePocket)
		api.PUT("/config/pockets/:id", c.ConfigHandler.UpdatePocket)
		api.DELETE("/config/pockets/:id", c.ConfigHandler.DeletePocket)

		// Gastos fijos
		api.GET("/fixed-expenses/:month", c.FixedExpenseHandler.GetByMonth)
		api.POST("/fixed-expenses", c.FixedExpenseHandler.Create)
		api.PUT("/fixed-expenses/:id", c.FixedExpenseHandler.Update)
		api.PUT("/fixed-expenses/:id/status", c.FixedExpenseHandler.UpdateStatus)

		// Gastos diarios
		api.GET("/daily-expenses/:month", c.DailyExpenseHandler.GetByMonth)
		api.POST("/daily-expenses", c.DailyExpenseHandler.Create)
		api.PUT("/daily-expenses/:id", c.DailyExpenseHandler.Update)
		api.DELETE("/daily-expenses/:id", c.DailyExpenseHandler.Delete)
	}
}
