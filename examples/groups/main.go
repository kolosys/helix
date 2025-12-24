// Package main demonstrates route grouping and API versioning.
package main

import (
	"context"
	"log"

	"github.com/kolosys/helix"
	"github.com/kolosys/helix/middleware"
)

// Product represents a product in the system.
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ProductV2 is the v2 representation with additional fields.
type ProductV2 struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

var products = []Product{
	{ID: 1, Name: "Widget", Price: 9.99},
	{ID: 2, Name: "Gadget", Price: 19.99},
}

var productsV2 = []ProductV2{
	{ID: 1, Name: "Widget", Price: 9.99, Description: "A useful widget", Tags: []string{"tool", "utility"}},
	{ID: 2, Name: "Gadget", Price: 19.99, Description: "A fancy gadget", Tags: []string{"electronics", "cool"}},
}

func main() {
	s := helix.Default(&helix.Options{
		Addr: ":8080",
	})

	// Public routes (no prefix)
	s.GET("/", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]string{
			"message": "Welcome to the API",
			"version": "Use /v1 or /v2 for API endpoints",
		})
	}))

	s.GET("/health", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]string{"status": "healthy"})
	}))

	// API v1 routes
	v1 := s.Group("/v1")

	v1.GET("/products", helix.HandleNoRequest(func(ctx context.Context) ([]Product, error) {
		return products, nil
	}))

	v1.GET("/products/{id}", helix.Handle(func(ctx context.Context, req struct {
		ID int `path:"id"`
	}) (Product, error) {
		for _, p := range products {
			if p.ID == req.ID {
				return p, nil
			}
		}
		return Product{}, helix.NotFoundf("product %d not found", req.ID)
	}))

	// API v2 routes with additional fields
	v2 := s.Group("/v2")

	v2.GET("/products", helix.HandleNoRequest(func(ctx context.Context) ([]ProductV2, error) {
		return productsV2, nil
	}))

	v2.GET("/products/{id}", helix.Handle(func(ctx context.Context, req struct {
		ID int `path:"id"`
	}) (ProductV2, error) {
		for _, p := range productsV2 {
			if p.ID == req.ID {
				return p, nil
			}
		}
		return ProductV2{}, helix.NotFoundf("product %d not found", req.ID)
	}))

	// Protected admin routes within v2 - no middleware casting needed!
	adminV2 := v2.Group("/admin", middleware.BasicAuth("admin", "secret"))

	adminV2.POST("/products", helix.HandleCreated(func(ctx context.Context, req struct {
		Name        string   `json:"name"`
		Price       float64  `json:"price"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}) (ProductV2, error) {
		newProduct := ProductV2{
			ID:          len(productsV2) + 1,
			Name:        req.Name,
			Price:       req.Price,
			Description: req.Description,
			Tags:        req.Tags,
		}
		productsV2 = append(productsV2, newProduct)
		return newProduct, nil
	}))

	adminV2.DELETE("/products/{id}", helix.HandleNoResponse(func(ctx context.Context, req struct {
		ID int `path:"id"`
	}) error {
		for i, p := range productsV2 {
			if p.ID == req.ID {
				productsV2 = append(productsV2[:i], productsV2[i+1:]...)
				return nil
			}
		}
		return helix.NotFoundf("product %d not found", req.ID)
	}))

	// Nested groups for organization
	internal := s.Group("/internal")
	metrics := internal.Group("/metrics")

	metrics.GET("/requests", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]any{
			"total_requests":      12345,
			"requests_per_second": 42.5,
		})
	}))

	metrics.GET("/errors", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]any{
			"total_errors": 42,
			"error_rate":   0.003,
		})
	}))

	log.Println("Server starting on :8080")
	log.Println("Routes:")
	s.PrintRoutes(log.Writer())

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
