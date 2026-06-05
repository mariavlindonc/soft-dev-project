package main

import (
	"log"
	"net/http"
	"os"

	db "backend/dao"
	"backend/domain"

	"github.com/gin-gonic/gin"
)

func main() {
	dbase, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer dbase.Close()

	if err := dbase.Migrate(&domain.User{}, &domain.Event{}, &domain.Ticket{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	r := gin.Default()

	r.Use(corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	addr := os.Getenv("APP_PORT")
	if addr == "" {
		addr = "8080"
	}
	log.Printf("server starting on :%s", addr)
	if err := r.Run(":" + addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
