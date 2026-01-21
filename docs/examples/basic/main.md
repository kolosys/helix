# main

This example demonstrates basic usage of the library.

## Source Code

```go
// Package main demonstrates the most basic usage of the helix framework.
// The recommended pattern is HandleCtx, which provides a fluent API with automatic error handling.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kolosys/helix"
)

func main() {
	// Create a new server with default settings
	s := helix.Default(&helix.Options{
		Addr: ":8080",
	})

	// Recommended: HandleCtx provides a fluent API with automatic error handling
	s.GET("/", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]string{
			"message": "Welcome to Helix!",
		})
	}))

	// HandleCtx with query parameters
	s.GET("/hello", helix.HandleCtx(func(c *helix.Ctx) error {
		name := c.QueryDefault("name", "World")
		return c.OK(map[string]string{
			"message": fmt.Sprintf("Hello, %s!", name),
		})
	}))

	// Handler with path parameters
	s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
		id := c.Param("id")
		return c.OK(map[string]string{
			"id":   id,
			"name": "John Doe",
		})
	}))

	// Handler returning an error (automatically converted to RFC 7807 Problem)
	s.GET("/error", helix.HandleCtx(func(c *helix.Ctx) error {
		return helix.NotFoundf("resource not found")
	}))

	// Lifecycle hooks
	s.OnStart(func(s *helix.Server) {
		log.Printf("Server starting on %s", s.Addr())
	})

	s.OnStop(func(ctx context.Context, s *helix.Server) {
		log.Println("Server shutting down...")
	})

	// Run server with graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

```

## Running the Example

To run this example:

```bash
cd basic
go run main.go
```

## Expected Output

```
Hello from Proton examples!
```
