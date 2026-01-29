// Package main demonstrates versioned API endpoints using helix.HandleVersions.
// This pattern allows serving multiple API versions from the same endpoint,
// switching behavior based on a header value (e.g., API-Version, Accept-Version).
package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/kolosys/helix"
)

// V1 Models - Original API version
type (
	// ProductV1 is the original product model.
	ProductV1 struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	// GetProductRequest is shared across versions.
	GetProductRequest struct {
		ID int `path:"id"`
	}

	// ListProductsRequest is shared across versions.
	ListProductsRequest struct {
		Page  int `query:"page"`
		Limit int `query:"limit"`
	}

	// CreateProductV1Request for v1 API.
	CreateProductV1Request struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
)

// V2 Models - Enhanced API version with more fields
type (
	// ProductV2 adds category and metadata.
	ProductV2 struct {
		ID          int               `json:"id"`
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Price       float64           `json:"price"`
		Currency    string            `json:"currency"`
		Category    string            `json:"category"`
		Tags        []string          `json:"tags"`
		Metadata    map[string]string `json:"metadata,omitempty"`
		CreatedAt   time.Time         `json:"created_at"`
		UpdatedAt   time.Time         `json:"updated_at"`
	}

	// CreateProductV2Request for v2 API with more fields.
	CreateProductV2Request struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Price       float64           `json:"price"`
		Currency    string            `json:"currency"`
		Category    string            `json:"category"`
		Tags        []string          `json:"tags"`
		Metadata    map[string]string `json:"metadata,omitempty"`
	}
)

// ProductStore provides thread-safe storage.
type ProductStore struct {
	mu       sync.RWMutex
	products map[int]*ProductV2 // Store the full v2 model internally
	nextID   int
}

func NewProductStore() *ProductStore {
	now := time.Now()
	return &ProductStore{
		products: map[int]*ProductV2{
			1: {
				ID:          1,
				Name:        "Widget",
				Description: "A versatile widget for all your needs",
				Price:       29.99,
				Currency:    "USD",
				Category:    "gadgets",
				Tags:        []string{"popular", "sale"},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			2: {
				ID:          2,
				Name:        "Gadget",
				Description: "The latest and greatest gadget",
				Price:       49.99,
				Currency:    "USD",
				Category:    "electronics",
				Tags:        []string{"new", "featured"},
				Metadata:    map[string]string{"color": "blue"},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		nextID: 3,
	}
}

// toV1 converts a V2 product to V1 format (downgrades).
func (p *ProductV2) toV1() ProductV1 {
	return ProductV1{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
	}
}

// ProductService handles business logic for products.
type ProductService struct {
	store *ProductStore
}

func NewProductService(store *ProductStore) *ProductService {
	return &ProductService{store: store}
}

// V1 Handlers
func (svc *ProductService) GetProductV1(ctx context.Context, req GetProductRequest) (ProductV1, error) {
	svc.store.mu.RLock()
	defer svc.store.mu.RUnlock()

	product, ok := svc.store.products[req.ID]
	if !ok {
		return ProductV1{}, helix.NotFoundf("product %d not found", req.ID)
	}
	return product.toV1(), nil
}

func (svc *ProductService) ListProductsV1(ctx context.Context, req ListProductsRequest) ([]ProductV1, error) {
	svc.store.mu.RLock()
	defer svc.store.mu.RUnlock()

	products := make([]ProductV1, 0, len(svc.store.products))
	for _, p := range svc.store.products {
		products = append(products, p.toV1())
	}
	return products, nil
}

func (svc *ProductService) CreateProductV1(ctx context.Context, req CreateProductV1Request) (ProductV1, error) {
	svc.store.mu.Lock()
	defer svc.store.mu.Unlock()

	now := time.Now()
	product := &ProductV2{
		ID:        svc.store.nextID,
		Name:      req.Name,
		Price:     req.Price,
		Currency:  "USD",
		CreatedAt: now,
		UpdatedAt: now,
	}
	svc.store.products[product.ID] = product
	svc.store.nextID++

	return product.toV1(), nil
}

// V2 Handlers
func (svc *ProductService) GetProductV2(ctx context.Context, req GetProductRequest) (ProductV2, error) {
	svc.store.mu.RLock()
	defer svc.store.mu.RUnlock()

	product, ok := svc.store.products[req.ID]
	if !ok {
		return ProductV2{}, helix.NotFoundf("product %d not found", req.ID)
	}
	return *product, nil
}

func (svc *ProductService) ListProductsV2(ctx context.Context, req ListProductsRequest) ([]ProductV2, error) {
	svc.store.mu.RLock()
	defer svc.store.mu.RUnlock()

	products := make([]ProductV2, 0, len(svc.store.products))
	for _, p := range svc.store.products {
		products = append(products, *p)
	}
	return products, nil
}

func (svc *ProductService) CreateProductV2(ctx context.Context, req CreateProductV2Request) (ProductV2, error) {
	svc.store.mu.Lock()
	defer svc.store.mu.Unlock()

	now := time.Now()
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	product := &ProductV2{
		ID:          svc.store.nextID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    currency,
		Category:    req.Category,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	svc.store.products[product.ID] = product
	svc.store.nextID++

	return *product, nil
}

const (
	// HeaderAPIVersion is the header used for API version selection.
	HeaderAPIVersion = "API-Version"

	// DefaultVersion is used when no version header is provided.
	DefaultVersion = "v1"
)

func main() {
	store := NewProductStore()
	svc := NewProductService(store)

	s := helix.Default(&helix.Options{
		Addr: ":8080",
	})

	// Versioned GET endpoint - same request type, different response types
	// Uses an interface{} (any) as the response type to support different versions
	s.GET("/products/{id}", helix.HandleVersions(
		HeaderAPIVersion,
		DefaultVersion,
		helix.VersionHandlerMap[GetProductRequest, any]{
			"v1": func(ctx context.Context, req GetProductRequest) (any, error) {
				return svc.GetProductV1(ctx, req)
			},
			"v2": func(ctx context.Context, req GetProductRequest) (any, error) {
				return svc.GetProductV2(ctx, req)
			},
		},
	))

	// Versioned LIST endpoint
	s.GET("/products", helix.HandleVersions(
		HeaderAPIVersion,
		DefaultVersion,
		helix.VersionHandlerMap[ListProductsRequest, any]{
			"v1": func(ctx context.Context, req ListProductsRequest) (any, error) {
				return svc.ListProductsV1(ctx, req)
			},
			"v2": func(ctx context.Context, req ListProductsRequest) (any, error) {
				return svc.ListProductsV2(ctx, req)
			},
		},
	))

	// For POST endpoints with different request types per version,
	// you need to use a common request type or handle binding manually.
	// Here's an approach using a superset request that works for both versions:
	s.POST("/products", helix.HandleVersions(
		HeaderAPIVersion,
		DefaultVersion,
		helix.VersionHandlerMap[CreateProductV2Request, any]{
			"v1": func(ctx context.Context, req CreateProductV2Request) (any, error) {
				// Convert to v1 request (only uses name and price)
				v1Req := CreateProductV1Request{
					Name:  req.Name,
					Price: req.Price,
				}
				return svc.CreateProductV1(ctx, v1Req)
			},
			"v2": func(ctx context.Context, req CreateProductV2Request) (any, error) {
				return svc.CreateProductV2(ctx, req)
			},
		},
	))

	log.Println("Versioned API Example")
	log.Println("=====================")
	log.Println("Usage examples:")
	log.Println("")
	log.Println("  V1 API (default):")
	log.Println("    curl http://localhost:8080/products")
	log.Println("    curl http://localhost:8080/products/1")
	log.Println("    curl -X POST http://localhost:8080/products -d '{\"name\":\"New\",\"price\":19.99}'")
	log.Println("")
	log.Println("  V2 API (explicit header):")
	log.Println("    curl -H 'API-Version: v2' http://localhost:8080/products")
	log.Println("    curl -H 'API-Version: v2' http://localhost:8080/products/1")
	log.Println("    curl -H 'API-Version: v2' -X POST http://localhost:8080/products \\")
	log.Println("         -d '{\"name\":\"New\",\"description\":\"Desc\",\"price\":19.99,\"category\":\"misc\"}'")
	log.Println("")

	s.PrintRoutes(log.Writer())

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
