package main

import (
	"net/http"
	"os"
	"strings"

	"backend/clients"
	"backend/controllers"
	db "backend/dao"
	"backend/domain"
	"backend/logger"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	dbase, err := db.NewDatabase()
	if err != nil {
		logger.Fatal("failed to initialize database: %v", err)
	}
	defer dbase.Close()

	if err := dbase.Migrate(&domain.User{}, &domain.Event{}, &domain.Ticket{}); err != nil {
		logger.Fatal("failed to run migrations: %v", err)
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

	r.Use(securityHeadersMiddleware())
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
			events.GET("/:id/sale-status", eventController.GetSaleStatus)
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
			admin.GET("/reports/events/:id", adminController.GetEventReport)
		}
	}

	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")
	if certFile == "" || keyFile == "" {
		logger.Fatal("TLS_CERT_FILE and TLS_KEY_FILE must be set for HTTPS")
	}

	addr := os.Getenv("APP_PORT")
	if addr == "" {
		addr = "8443"
	}
	logger.Info("server starting on :%s (HTTPS)", addr)
	if err := r.RunTLS(":"+addr, certFile, keyFile); err != nil {
		logger.Fatal("failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	allowed := parseOrigins(raw)
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// If a specific origin is allowed and matches, echo it back.
		if origin != "" && isOriginAllowed(origin, allowed) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if len(allowed) == 1 && allowed[0] == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		}

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

func parseOrigins(raw string) []string {
	if raw == "" {
		return []string{"*"}
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return []string{"*"}
	}
	return result
}

func isOriginAllowed(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == "*" || a == origin {
			return true
		}
	}
	return false
}

func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}
