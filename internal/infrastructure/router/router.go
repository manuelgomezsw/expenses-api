package router

import (
	"expenses-api/internal/infrastructure/middleware/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func StartApp() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	router := gin.Default()
	router.Use(cors.NewCorsMiddleware().Execute())
	mapURLs(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return router.Run(":" + port)
}
