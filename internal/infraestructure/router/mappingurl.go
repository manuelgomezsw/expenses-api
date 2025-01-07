package router

import (
	"expenses-api/internal/infraestructure/controller/paymentstype"
	"expenses-api/internal/infraestructure/controller/pockets"
	"github.com/gin-gonic/gin"
)

func mapURLs(router *gin.Engine) {
	expensesUrls(router)
	paymentsTypeUrls(router)
	pocketsUrls(router)
}

func expensesUrls(router *gin.Engine) {
	//router.GET("/payments/type", paymentstype.Get)
	//router.GET("/payments/type/:payment_type_id", paymentstype.GetByID)
	//router.POST("/payments/type", paymentstype.Create)
	//router.PUT("/payments/type/:payment_type_id", paymentstype.Update)
	//router.DELETE("/payments/type/:payment_type_id", paymentstype.Delete)
}

func paymentsTypeUrls(router *gin.Engine) {
	router.GET("/payments/type", paymentstype.Get)
	router.GET("/payments/type/:payment_type_id", paymentstype.GetByID)
	router.POST("/payments/type", paymentstype.Create)
	router.PUT("/payments/type/:payment_type_id", paymentstype.Update)
	router.DELETE("/payments/type/:payment_type_id", paymentstype.Delete)
}

func pocketsUrls(router *gin.Engine) {
	router.GET("/pockets", pockets.Get)
	router.GET("/pockets/:pocket_id", pockets.GetByID)
	router.POST("/pockets", pockets.Create)
	router.PUT("/pockets/:pocket_id", pockets.Update)
	router.DELETE("/pockets/:pocket_id", pockets.Delete)
}
