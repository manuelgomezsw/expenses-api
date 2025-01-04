package router

import (
	"expenses-api/internal/infraestructure/controller/paymentstype"
	"github.com/gin-gonic/gin"
)

func mapURLs(router *gin.Engine) {
	expensesUrls(router)
	paymentsTypeUrls(router)
	pocketsUrls(router)
}

func expensesUrls(router *gin.Engine) {
	router.GET("/payments/type", paymentstype.Get)
	router.GET("/payments/type/:payment_type_id", paymentstype.GetByID)
	router.POST("/payments/type", paymentstype.Create)
	router.PUT("/payments/type/:payment_type_id", paymentstype.Update)
	router.DELETE("/payments/type/:payment_type_id", paymentstype.Delete)
}

func paymentsTypeUrls(router *gin.Engine) {
	//router.POST("/quotes", quotesRegistry.Create)
}

func pocketsUrls(router *gin.Engine) {
	//router.PUT("/quotes", quotesRegistry.Create)
}
