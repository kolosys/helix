# main

This example demonstrates basic usage of the library.

## Source Code

```go
// Package main demonstrates advanced helix features:
// - Modular route organization
// - Service registration and dependency injection
// - Pagination helpers
// - Health check endpoints
package main

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/kolosys/helix"
)

// =============================================================================
// Domain Models
// =============================================================================

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// =============================================================================
// Services (Business Logic Layer)
// =============================================================================

// UserService handles user business logic.
type UserService struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users: map[int]User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
			3: {ID: 3, Name: "Charlie", Email: "charlie@example.com", CreatedAt: time.Now()},
		},
		nextID: 4,
	}
}

func (s *UserService) List(ctx context.Context, page, limit int) ([]User, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}

	total := len(users)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		return []User{}, total, nil
	}
	if end > total {
		end = total
	}

	return users[start:end], total, nil
}

func (s *UserService) Get(ctx context.Context, id int) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return User{}, helix.NotFoundf("user %d not found", id)
	}
	return user, nil
}

func (s *UserService) Create(ctx context.Context, name, email string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{
		ID:        s.nextID,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
	s.users[user.ID] = user
	s.nextID++
	return user, nil
}

// PostService handles post business logic.
type PostService struct {
	mu     sync.RWMutex
	posts  map[int]Post
	nextID int
}

func NewPostService() *PostService {
	return &PostService{
		posts: map[int]Post{
			1: {ID: 1, UserID: 1, Title: "Hello World", Content: "First post!", CreatedAt: time.Now()},
			2: {ID: 2, UserID: 1, Title: "Helix Guide", Content: "How to use helix...", CreatedAt: time.Now()},
		},
		nextID: 3,
	}
}

func (s *PostService) ListByUser(ctx context.Context, userID, page, limit int) ([]Post, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var posts []Post
	for _, p := range s.posts {
		if p.UserID == userID {
			posts = append(posts, p)
		}
	}

	total := len(posts)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		return []Post{}, total, nil
	}
	if end > total {
		end = total
	}

	return posts[start:end], total, nil
}

// =============================================================================
// Modules (Route Organization)
// =============================================================================

// UserModule handles user-related routes.
type UserModule struct{}

func (m *UserModule) Register(r helix.RouteRegistrar) {
	// GET /users - List users with pagination
	r.GET("/", helix.HandleCtx(func(c *helix.Ctx) error {
		userSvc := helix.MustGet[*UserService]()

		// Use built-in pagination binding
		p := c.BindPagination(20, 100)

		users, total, err := userSvc.List(c.Context(), p.GetPage(), p.GetLimit(20, 100))
		if err != nil {
			return err
		}

		// Use built-in paginated response helper
		return c.Paginated(users, total, p.GetPage(), p.GetLimit(20, 100))
	}))

	// GET /users/{id} - Get user by ID
	r.GET("/{id}", helix.Handle(func(ctx context.Context, req struct {
		ID int `path:"id"`
	}) (User, error) {
		userSvc := helix.MustGet[*UserService]()
		return userSvc.Get(ctx, req.ID)
	}))

	// POST /users - Create user
	r.POST("/", helix.HandleCreated(func(ctx context.Context, req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}) (User, error) {
		userSvc := helix.MustGet[*UserService]()
		return userSvc.Create(ctx, req.Name, req.Email)
	}))

	// Mount posts sub-module
	r.Group("/{userId}/posts").GET("/", helix.HandleCtx(func(c *helix.Ctx) error {
		postSvc := helix.MustGet[*PostService]()

		userID, err := c.ParamInt("userId")
		if err != nil {
			return c.BadRequest("invalid user ID")
		}

		p := c.BindPagination(10, 50)
		posts, total, err := postSvc.ListByUser(c.Context(), userID, p.GetPage(), p.GetLimit(10, 50))
		if err != nil {
			return err
		}

		return c.Paginated(posts, total, p.GetPage(), p.GetLimit(10, 50))
	}))
}

// =============================================================================
// Health Checks
// =============================================================================

func setupHealthChecks(userSvc *UserService) *helix.HealthBuilder {
	return helix.Health().
		Version("1.0.0").
		Timeout(5*time.Second).
		CheckFunc("database", func(ctx context.Context) error {
			// Simulate DB health check
			return nil
		}).
		CheckFunc("user_service", func(ctx context.Context) error {
			// Check if user service is healthy
			_, err := userSvc.Get(ctx, 1)
			if err != nil {
				return errors.New("user service unhealthy")
			}
			return nil
		}).
		CheckFunc("external_api", func(ctx context.Context) error {
			// Simulate external API check
			time.Sleep(10 * time.Millisecond)
			return nil
		})
}

// =============================================================================
// Main
// =============================================================================

func main() {
	// Create services
	userSvc := NewUserService()
	postSvc := NewPostService()

	// Register services in global registry (dependency injection)
	helix.Register(userSvc)
	helix.Register(postSvc)

	// Create server
	s := helix.Default(&helix.Options{
		Addr: ":8080",
	})

	// Root endpoint
	s.GET("/", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]any{
			"name":    "Modular API",
			"version": "1.0.0",
			"docs": map[string]string{
				"users":     "GET /users",
				"health":    "GET /health",
				"liveness":  "GET /health/live",
				"readiness": "GET /health/ready",
			},
		})
	}))

	// Mount user module
	s.Mount("/users", &UserModule{})

	// Health check endpoints
	health := setupHealthChecks(userSvc)
	s.GET("/health", health.Handler())
	s.GET("/health/live", helix.LivenessHandler())
	s.GET("/health/ready", helix.ReadinessHandler(
		func(ctx context.Context) error {
			// Check database connection
			return nil
		},
		func(ctx context.Context) error {
			// Check external services
			return nil
		},
	))

	// Alternative: Mount routes using a function
	s.MountFunc("/admin", func(r helix.RouteRegistrar) {
		r.GET("/stats", helix.HandleCtx(func(c *helix.Ctx) error {
			return c.OK(map[string]any{
				"total_users": 3,
				"total_posts": 2,
				"uptime":      time.Since(time.Now()).String(),
			})
		}))
	})

	// Pre-compile middleware chain for optimal performance
	s.Build()

	log.Println("Server starting on :8080")
	log.Println("Routes:")
	s.PrintRoutes(log.Writer())

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}

```

## Running the Example

To run this example:

```bash
cd modular
go run main.go
```

## Expected Output

```
Hello from Proton examples!
```
