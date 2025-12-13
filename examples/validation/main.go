// Package main demonstrates request binding and validation.
package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"slices"

	"github.com/kolosys/helix"
)

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// CreateAccountRequest demonstrates binding from multiple sources with validation.
type CreateAccountRequest struct {
	// From JSON body
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`

	// From query parameters
	ReferralCode string `query:"referral"`
	Newsletter   bool   `query:"newsletter"`

	// From headers
	UserAgent string `header:"User-Agent"`
	Language  string `header:"Accept-Language"`
}

// Validate implements helix.Validatable using ValidationErrors for field-level errors.
// This produces RFC 7807 responses with an "errors" array for each field.
func (r *CreateAccountRequest) Validate() error {
	v := helix.NewValidationErrors()

	// Email validation
	if r.Email == "" {
		v.Add("email", "email is required")
	} else if !emailRegex.MatchString(r.Email) {
		v.Add("email", "email is invalid")
	}

	// Password validation
	if r.Password == "" {
		v.Add("password", "password is required")
	} else if len(r.Password) < 8 {
		v.Add("password", "password must be at least 8 characters")
	}

	// Name validation
	if r.FirstName == "" {
		v.Add("first_name", "first_name is required")
	}
	if r.LastName == "" {
		v.Add("last_name", "last_name is required")
	}

	// Age validation
	if r.Age < 13 {
		v.Add("age", "age must be at least 13")
	}
	if r.Age > 150 {
		v.Add("age", "age must be less than 150")
	}

	return v.Err() // Returns nil if no errors, or *ValidationErrors
}

// CreateAccountResponse is returned after successful account creation.
type CreateAccountResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Message   string `json:"message"`
}

// SearchRequest demonstrates query parameter binding with defaults.
type SearchRequest struct {
	Query    string   `query:"q,required"`
	Page     int      `query:"page"`
	Limit    int      `query:"limit"`
	Sort     string   `query:"sort"`
	Order    string   `query:"order"`
	Tags     []string `query:"tags"`
	MinPrice float64  `query:"min_price"`
	MaxPrice float64  `query:"max_price"`
	InStock  bool     `query:"in_stock"`
}

// Validate implements helix.Validatable.
func (r *SearchRequest) Validate() error {
	// Set defaults
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Limit <= 0 {
		r.Limit = 20
	}
	if r.Limit > 100 {
		return helix.BadRequestf("limit cannot exceed 100")
	}
	if r.Sort == "" {
		r.Sort = "relevance"
	}
	if r.Order == "" {
		r.Order = "desc"
	}

	// Validate order
	if r.Order != "asc" && r.Order != "desc" {
		return helix.BadRequestf("order must be 'asc' or 'desc'")
	}

	// Validate price range
	if r.MinPrice < 0 {
		return helix.BadRequestf("min_price cannot be negative")
	}
	if r.MaxPrice > 0 && r.MinPrice > r.MaxPrice {
		return helix.BadRequestf("min_price cannot be greater than max_price")
	}

	return nil
}

// SearchResult represents a search result item.
type SearchResult struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Price float64  `json:"price"`
	Tags  []string `json:"tags"`
}

// SearchResponse is the search results.
type SearchResponse struct {
	Query   string         `json:"query"`
	Page    int            `json:"page"`
	Limit   int            `json:"limit"`
	Total   int            `json:"total"`
	Results []SearchResult `json:"results"`
}

// PathParamRequest demonstrates path parameter binding with type conversion.
type PathParamRequest struct {
	UserID    int    `path:"userId"`
	PostID    string `path:"postId"` // UUID as string
	CommentID int64  `path:"commentId"`
}

func main() {
	s := helix.Default(
		helix.WithAddr(":8080"),
	)

	// Home page
	s.GET("/", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]string{
			"message": "Validation Examples API",
			"docs":    "Try POST /accounts, GET /search, or GET /users/{userId}/posts/{postId}/comments/{commentId}",
		})
	}))

	// Account creation with full validation
	// POST /accounts?referral=ABC123&newsletter=true
	// Body: {"email": "user@example.com", "password": "secret123", "first_name": "John", "last_name": "Doe", "age": 25}
	s.POST("/accounts", helix.Handle(func(ctx context.Context, req CreateAccountRequest) (CreateAccountResponse, error) {
		// At this point, the request has been validated
		log.Printf("Creating account for %s %s <%s>", req.FirstName, req.LastName, req.Email)
		log.Printf("Referral code: %s, Newsletter: %v", req.ReferralCode, req.Newsletter)
		log.Printf("User-Agent: %s, Language: %s", req.UserAgent, req.Language)

		return CreateAccountResponse{
			ID:        12345,
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Message:   "Account created successfully",
		}, nil
	}))

	// Search with query parameters and validation
	// GET /search?q=widget&page=1&limit=20&sort=price&order=asc&tags=electronics,gadgets&min_price=10&max_price=100&in_stock=true
	s.GET("/search", helix.Handle(func(ctx context.Context, req SearchRequest) (SearchResponse, error) {
		// Mock search results
		results := []SearchResult{
			{ID: 1, Name: "Widget A", Price: 29.99, Tags: []string{"electronics", "gadgets"}},
			{ID: 2, Name: "Widget B", Price: 49.99, Tags: []string{"electronics"}},
			{ID: 3, Name: "Widget C", Price: 19.99, Tags: []string{"gadgets", "tools"}},
		}

		// Filter by tags if provided
		if len(req.Tags) > 0 {
			filtered := make([]SearchResult, 0)
			for _, r := range results {
				for _, tag := range req.Tags {
					if slices.Contains(r.Tags, tag) {
						filtered = append(filtered, r)
					}
				}
			}
			results = filtered
		}

		// Filter by price range
		if req.MinPrice > 0 || req.MaxPrice > 0 {
			filtered := make([]SearchResult, 0)
			for _, r := range results {
				if req.MinPrice > 0 && r.Price < req.MinPrice {
					continue
				}
				if req.MaxPrice > 0 && r.Price > req.MaxPrice {
					continue
				}
				filtered = append(filtered, r)
			}
			results = filtered
		}

		return SearchResponse{
			Query:   req.Query,
			Page:    req.Page,
			Limit:   req.Limit,
			Total:   len(results),
			Results: results,
		}, nil
	}))

	// Path parameters with type conversion
	// GET /users/123/posts/abc-123-def/comments/456789
	s.GET("/users/{userId}/posts/{postId}/comments/{commentId}", helix.Handle(func(ctx context.Context, req PathParamRequest) (map[string]any, error) {
		return map[string]any{
			"user_id":    req.UserID,
			"post_id":    req.PostID,
			"comment_id": req.CommentID,
			"message":    fmt.Sprintf("Fetching comment %d on post %s by user %d", req.CommentID, req.PostID, req.UserID),
		}, nil
	}))

	// Manual binding example using Ctx
	s.POST("/manual", helix.HandleCtx(func(c *helix.Ctx) error {
		// Manual query parameter extraction with defaults
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 20)
		sortBy := c.QueryDefault("sort", "created_at")

		// Manual header extraction
		contentType := c.Header("Content-Type")
		authorization := c.Header("Authorization")

		// Manual body binding
		var body struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}
		if err := c.Bind(&body); err != nil {
			return c.BadRequest("invalid request body")
		}

		return c.OK(map[string]any{
			"page":         page,
			"limit":        limit,
			"sort":         sortBy,
			"content_type": contentType,
			"has_auth":     authorization != "",
			"body":         body,
		})
	}))

	log.Println("Server starting on :8080")
	log.Println("Routes:")
	s.PrintRoutes(log.Writer())

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
