package main

import (
	"log"
	"net/http"
	"os"

	"backend/clients"
	"backend/controllers"
	db "backend/dao"
	"backend/domain"
	"backend/services"

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

	// Initialize DAOs.
	conn := dbase.GetConnection()

	userDAO := db.NewUserDAO(conn)
	eventDAO := db.NewEventDAO(conn)
	ticketDAO := db.NewTicketDAO(conn)

	// Initialize external clients.
	emailClient := clients.NewEmailClient()

	// Initialize services.
	authService := services.NewAuthService(userDAO)
	eventService := services.NewEventService(eventDAO, ticketDAO)
	ticketService := services.NewTicketService(ticketDAO, eventDAO, userDAO, emailClient)
	reportService := services.NewReportService(eventDAO, ticketDAO, userDAO)

	// Initialize controllers.
	authController := controllers.NewAuthController(authService)
	eventController := controllers.NewEventController(eventService)
	ticketController := controllers.NewTicketController(ticketService)
	adminController := controllers.NewAdminController(reportService)

	r := gin.Default()

	r.Use(corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api")
	{
		// Public auth endpoints.
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Public event listing.
		events := v1.Group("/events")
		{
			events.GET("", eventController.GetAll)
			events.GET("/:id", eventController.GetByID)
		}

		// Authenticated ticket operations.
		tickets := v1.Group("/tickets")
		tickets.Use(controllers.AuthRequired())
		{
			tickets.POST("/purchase", ticketController.Purchase)
			tickets.GET("", ticketController.GetMyTickets)
			tickets.PATCH("/:id/cancel", ticketController.Cancel)
			tickets.PATCH("/:id/transfer", ticketController.Transfer)
		}

		// Admin-only endpoints.
		admin := v1.Group("/admin")
		admin.Use(controllers.AuthRequired(), controllers.AdminRequired())
		{
			admin.POST("/events", eventController.Create)
			admin.PUT("/events/:id", eventController.Update)
			admin.DELETE("/events/:id", eventController.Delete)
			admin.GET("/reports", adminController.GetReports)
		}
	}

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
