// Package main demonstrates the use of middleware in helix.
package main

import (
	"log"
	"time"

	"github.com/kolosys/helix"
	"github.com/kolosys/helix/middleware"
)

func main() {
	// Create a server without default middleware
	s := helix.New(&helix.Options{
		Addr: ":8080",
	})

	// Add global middleware in order - no explicit casting needed!
	s.Use(middleware.RequestID()) // Generate unique request IDs
	s.Use(middleware.Logger())    // Log requests (dev format by default)
	s.Use(middleware.Recover())   // Recover from panics

	// CORS middleware for API endpoints
	s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://example.com", "https://app.example.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	// Compression middleware
	s.Use(middleware.Compress())

	// Public routes
	s.GET("/", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]string{
			"message": "Hello, World!",
		})
	}))

	s.GET("/health", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]any{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	}))

	// Protected routes with Basic Auth - no casting needed!
	admin := s.Group("/admin", middleware.BasicAuth("admin", "secret"))

	admin.GET("/dashboard", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]string{
			"message": "Welcome to the admin dashboard!",
		})
	}))

	admin.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]any{
			"users": []string{"alice", "bob", "charlie"},
		})
	}))

	// Rate-limited API routes - no casting needed!
	api := s.Group("/api", middleware.RateLimit(10, 5)) // 10 requests/second, burst of 5

	api.GET("/data", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]string{
			"data": "This endpoint is rate-limited",
		})
	}))

	// Routes with timeout middleware
	slow := s.Group("/slow", middleware.Timeout(2*time.Second))

	slow.GET("/process", helix.HandleCtx(func(c *helix.Ctx) error {
		// Simulate slow processing
		select {
		case <-time.After(1 * time.Second):
			return c.OK(map[string]string{
				"result": "Processing complete",
			})
		case <-c.Context().Done():
			return helix.ErrGatewayTimeout.WithDetail("Request timed out")
		}
	}))

	log.Println("Server starting on :8080")
	log.Println("Routes:")
	s.PrintRoutes(log.Writer())

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
