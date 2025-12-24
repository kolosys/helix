# main

This example demonstrates basic usage of the library.

## Source Code

```go
// Package main demonstrates the most basic usage of the helix framework.
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kolosys/helix"
	"github.com/kolosys/helix/logs"
)

func main() {
	logs.SetDefaultLevel(logs.InfoLevel)
	logs.SetDefaultFormatter(&logs.JSONFormatter{})

	// Create a new server with default settings
	s := helix.Default(&helix.Options{
		Addr: ":8080",
	})

	// Simple handler using http.HandlerFunc
	s.GET("/", func(w http.ResponseWriter, r *http.Request) {
		helix.OK(w, map[string]string{
			"message": "Welcome to Helix!",
		})
	})

	// Handler using Ctx for a cleaner API
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
		logs.Infof("Server starting on %s", s.Addr())
	})

	s.OnStop(func(ctx context.Context, s *helix.Server) {
		logs.Info("Server shutting down...")
	})

	// Run server with graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Run(ctx); err != nil {
		logs.Fatal(err.Error())
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
