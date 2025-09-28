package router

import (
	"expenses-api/internal/infraestructure/controller/frontend"
	"github.com/gin-gonic/gin"
)

func mapURLs(router *gin.Engine) {
	// Solo rutas para el frontend Angular
	frontendUrls(router)
}

// frontendUrls define las rutas específicas para el frontend Angular
func frontendUrls(router *gin.Engine) {
	// Inicializar controladores
	summaryController := frontend.NewSummaryController()
	configController := frontend.NewConfigController()
	fixedExpensesController := frontend.NewFixedExpensesController()
	dailyExpensesController := frontend.NewDailyExpensesController()

	// Grupo de rutas API
	api := router.Group("/api")
	{
		// Resumen mensual
		api.GET("/summary/:month", summaryController.GetMonthlySummary)

		// Configuración
		api.GET("/config/income", configController.GetIncome)
		api.PUT("/config/income", configController.UpdateIncome)
		api.GET("/config/pockets", configController.GetPockets)
		api.POST("/config/pockets", configController.CreatePocket)
		api.PUT("/config/pockets/:id", configController.UpdatePocket)
		api.DELETE("/config/pockets/:id", configController.DeletePocket)

		// Gastos fijos
		api.GET("/fixed-expenses/:month", fixedExpensesController.GetByMonth)
		api.PUT("/fixed-expenses/:id/status", fixedExpensesController.UpdateStatus)

		// Gastos diarios
		api.GET("/daily-expenses/:month", dailyExpensesController.GetByMonth)
		api.POST("/daily-expenses", dailyExpensesController.Create)
		api.PUT("/daily-expenses/:id", dailyExpensesController.Update)
		api.DELETE("/daily-expenses/:id", dailyExpensesController.Delete)
	}
}