package router

import (
	"expenses-api/internal/infraestructure/controller/budget"
	"expenses-api/internal/infraestructure/controller/cycles"
	"expenses-api/internal/infraestructure/controller/expenses"
	"expenses-api/internal/infraestructure/controller/paymentstype"
	"expenses-api/internal/infraestructure/controller/pockets"
	"github.com/gin-gonic/gin"
)

func mapURLs(router *gin.Engine) {
	expensesUrls(router)
	paymentsTypeUrls(router)
	pocketsUrls(router)
	cyclesUrls(router)
	budgetsUrls(router)
}

func expensesUrls(router *gin.Engine) {
	router.GET("/expenses", expenses.GetByActiveCycles)
	router.POST("/expenses", expenses.Create)
	router.PUT("/expenses/:expense_id", expenses.Update)
	router.DELETE("/expenses/:expense_id", expenses.Delete)
}

func paymentsTypeUrls(router *gin.Engine) {
	router.GET("/payments/type", paymentstype.Get)
	router.GET("/payments/type/:payment_type_id", paymentstype.GetByID)
	router.POST("/payments/type", paymentstype.Create)
	router.PUT("/payments/type/:payment_type_id", paymentstype.Update)
	router.DELETE("/payments/type/:payment_type_id", paymentstype.Delete)
}

func pocketsUrls(router *gin.Engine) {
	router.GET("/pockets/all", pockets.Get)
	router.GET("/pockets/active", pockets.GetActives)
	router.GET("/pockets/:pocket_id", pockets.GetByID)
	router.POST("/pockets", pockets.Create)
	router.PUT("/pockets/:pocket_id", pockets.Update)
	router.DELETE("/pockets/:pocket_id", pockets.Delete)
}

func cyclesUrls(router *gin.Engine) {
	router.GET("/cycles", cycles.GetAll)
	router.GET("/cycles/active", cycles.GetActive)
	router.GET("/cycles/:cycle_id", cycles.GetByID)
	router.POST("/cycles", cycles.Create)
	router.PUT("/cycles/:cycle_id", cycles.Update)
	router.DELETE("/cycles/:cycle_id", cycles.Delete)
}

func budgetsUrls(router *gin.Engine) {
	router.POST("/budgets/:cycle_id", budget.Calculate)
}
