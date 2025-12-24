# main

This example demonstrates basic usage of the library.

## Source Code

```go
// Package main demonstrates a full CRUD API using typed handlers.
package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/kolosys/helix"
	"github.com/kolosys/helix/logs"
)

// User represents a user in the system.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserStore provides thread-safe in-memory storage for users.
type UserStore struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

// NewUserStore creates a new UserStore with sample data.
func NewUserStore() *UserStore {
	return &UserStore{
		users: map[int]User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

// Request/Response types for typed handlers
type (
	// ListUsersRequest contains parameters for listing users.
	ListUsersRequest struct {
		Page  int `query:"page"`
		Limit int `query:"limit"`
	}

	// ListUsersResponse is the response for listing users.
	ListUsersResponse struct {
		Users []User `json:"users"`
		Total int    `json:"total"`
	}

	// GetUserRequest contains the user ID from the path.
	GetUserRequest struct {
		ID int `path:"id"`
	}

	// CreateUserRequest contains the data for creating a user.
	CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// UpdateUserRequest contains the data for updating a user.
	UpdateUserRequest struct {
		ID    int    `path:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// DeleteUserRequest contains the user ID to delete.
	DeleteUserRequest struct {
		ID int `path:"id"`
	}
)

// Validate implements helix.Validatable for CreateUserRequest.
// Uses ValidationErrors for collecting multiple errors.
func (r *CreateUserRequest) Validate() error {
	v := helix.NewValidationErrors()

	if r.Name == "" {
		v.Add("name", "name is required")
	}
	if r.Email == "" {
		v.Add("email", "email is required")
	}

	return v.Err() // Returns nil if no errors
}

func main() {
	store := NewUserStore()

	// Custom error handler that logs errors and provides custom formatting
	customErrorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		// Log the error for debugging
		logs.Errorf("Error handling request %s %s: %v", r.Method, r.URL.Path, err)

		// You can customize error handling here, e.g., different formats for different error types
		// For this example, we'll use the default Problem format but add custom logging
		// In production, you might want to:
		// - Send errors to a monitoring service
		// - Format errors differently for different clients
		// - Add request context to error responses
		// - Implement retry logic headers

		// Fall back to default error handling
		helix.HandleErrorDefault(w, r, err)
	}

	s := helix.Default(&helix.Options{
		Addr:         ":8080",
		ErrorHandler: customErrorHandler,
		// Uncomment to add a base path prefix to all routes:
		// BasePath: "/api/v1",
	})

	// List users - GET /users?page=1&limit=10
	s.GET("/users", helix.Handle(func(ctx context.Context, req ListUsersRequest) (ListUsersResponse, error) {
		store.mu.RLock()
		defer store.mu.RUnlock()

		users := make([]User, 0, len(store.users))
		for _, u := range store.users {
			users = append(users, u)
		}

		return ListUsersResponse{
			Users: users,
			Total: len(users),
		}, nil
	}))

	// Get user - GET /users/{id}
	s.GET("/users/{id}", helix.Handle(func(ctx context.Context, req GetUserRequest) (User, error) {
		store.mu.RLock()
		defer store.mu.RUnlock()

		user, ok := store.users[req.ID]
		if !ok {
			return User{}, helix.NotFoundf("user %d not found", req.ID)
		}

		return user, nil
	}))

	// Create user - POST /users (uses HandleCreated for 201 status)
	s.POST("/users", helix.HandleCreated(func(ctx context.Context, req CreateUserRequest) (User, error) {
		store.mu.Lock()
		defer store.mu.Unlock()

		user := User{
			ID:    store.nextID,
			Name:  req.Name,
			Email: req.Email,
		}
		store.users[user.ID] = user
		store.nextID++

		return user, nil
	}))

	// Update user - PUT /users/{id}
	s.PUT("/users/{id}", helix.Handle(func(ctx context.Context, req UpdateUserRequest) (User, error) {
		store.mu.Lock()
		defer store.mu.Unlock()

		if _, ok := store.users[req.ID]; !ok {
			return User{}, helix.NotFoundf("user %d not found", req.ID)
		}

		user := User(req)
		store.users[req.ID] = user

		return user, nil
	}))

	// Delete user - DELETE /users/{id}
	s.DELETE("/users/{id}", helix.HandleNoResponse(func(ctx context.Context, req DeleteUserRequest) error {
		store.mu.Lock()
		defer store.mu.Unlock()

		if _, ok := store.users[req.ID]; !ok {
			return helix.NotFoundf("user %d not found", req.ID)
		}

		delete(store.users, req.ID)
		return nil
	}))

	// Fallback route - handles all unmatched paths
	s.Any("/{path...}", helix.HandleCtx(func(c *helix.Ctx) error {
		path := c.Param("path")
		return c.Problem(helix.NotFoundf("route not found: /%s", path))
	}))

	// Print registered routes
	logs.Info("Registered routes:")
	s.PrintRoutes(log.Writer())

	logs.Println("Server starting on :8080")
	if err := s.Start(); err != nil {
		logs.Fatal(err.Error())
	}
}

```

## Running the Example

To run this example:

```bash
cd crud
go run main.go
```

## Expected Output

```
Hello from Proton examples!
```
