# main

This example demonstrates a full CRUD API using typed handlers.

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
	ListUsersRequest struct {
		Page  int `query:"page"`
		Limit int `query:"limit"`
	}

	ListUsersResponse struct {
		Users []User `json:"users"`
		Total int    `json:"total"`
	}

	GetUserRequest struct {
		ID int `path:"id"`
	}

	CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	UpdateUserRequest struct {
		ID    int    `path:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	DeleteUserRequest struct {
		ID int `path:"id"`
	}
)

// Validate implements helix.Validatable for CreateUserRequest.
func (r *CreateUserRequest) Validate() error {
	v := helix.NewValidationErrors()

	if r.Name == "" {
		v.Add("name", "name is required")
	}
	if r.Email == "" {
		v.Add("email", "email is required")
	}

	return v.Err()
}

func main() {
	store := NewUserStore()

	s := helix.Default(&helix.Options{
		Addr: ":8080",
	})

	// List users - GET /users
	s.GET("/users", helix.Handle(func(ctx context.Context, req ListUsersRequest) (ListUsersResponse, error) {
		store.mu.RLock()
		defer store.mu.RUnlock()

		users := make([]User, 0, len(store.users))
		for _, u := range store.users {
			users = append(users, u)
		}

		return ListUsersResponse{Users: users, Total: len(users)}, nil
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

	// Create user - POST /users
	s.POST("/users", helix.HandleCreated(func(ctx context.Context, req CreateUserRequest) (User, error) {
		store.mu.Lock()
		defer store.mu.Unlock()

		user := User{ID: store.nextID, Name: req.Name, Email: req.Email}
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

	log.Println("Server starting on :8080")
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
```

## Running the Example

```bash
cd crud
go run main.go
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /users | List all users |
| GET | /users/{id} | Get user by ID |
| POST | /users | Create new user |
| PUT | /users/{id} | Update user |
| DELETE | /users/{id} | Delete user |
