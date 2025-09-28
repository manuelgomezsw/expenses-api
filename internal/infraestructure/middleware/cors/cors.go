package cors

import (
	"expenses-api/internal/infraestructure/middleware"
	"expenses-api/internal/util/environment"
	"github.com/gin-gonic/gin"
)

type corsMiddleware struct{}

func (t corsMiddleware) Execute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", getCorsOrigin())
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

func NewCorsMiddleware() middleware.Middleware {
	return corsMiddleware{}
}

func getCorsOrigin() string {
	if environment.IsProductionEnv() {
		return "https://your-frontend-domain.com" // TODO: Configurar dominio de producción
	}
	return "http://localhost:4200" // Angular dev server
}
