# Quick Start

This guide will help you get started with Helix quickly with practical examples.

## Hello World

Here's the simplest possible Helix server:

```go
package main

import (
    "net/http"
    "github.com/kolosys/helix"
)

func main() {
    s := helix.Default(nil)

    s.GET("/", func(w http.ResponseWriter, r *http.Request) {
        helix.OK(w, map[string]string{"message": "Hello, World!"})
    })

    s.Start(":8080")
}
```

Run it:

```bash
go run main.go
```

Visit `http://localhost:8080` to see your response.

## Using Ctx for Cleaner Code

The `Ctx` type provides a fluent API for handlers:

```go
package main

import (
    "github.com/kolosys/helix"
)

func main() {
    s := helix.Default(nil)

    s.GET("/hello", helix.HandleCtx(func(c *helix.Ctx) error {
        name := c.QueryDefault("name", "World")
        return c.OK(map[string]string{
            "message": "Hello, " + name + "!",
        })
    }))

    s.Start(":8080")
}
```

## Path Parameters

Extract dynamic values from URLs:

```go
s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    id := c.Param("id")
    return c.OK(map[string]string{
        "id":   id,
        "name": "John Doe",
    })
}))
```

## Typed Handlers with Automatic Binding

Use generic handlers for type-safe request binding:

```go
package main

import (
    "context"
    "github.com/kolosys/helix"
)

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    s := helix.Default(nil)

    s.POST("/users", helix.Handle(func(ctx context.Context, req CreateUserRequest) (User, error) {
        // req is automatically bound from JSON body
        return User{
            ID:    1,
            Name:  req.Name,
            Email: req.Email,
        }, nil
    }))

    s.Start(":8080")
}
```

## Complete Example

Here's a complete example with multiple routes and error handling:

```go
package main

import (
    "context"
    "net/http"
    "github.com/kolosys/helix"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    s := helix.Default(nil)

    // Simple handler
    s.GET("/", func(w http.ResponseWriter, r *http.Request) {
        helix.OK(w, map[string]string{"message": "Welcome to Helix!"})
    })

    // Handler with Ctx
    s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
        id := c.Param("id")
        if id == "" {
            return helix.NotFoundf("user not found")
        }
        return c.OK(User{ID: 1, Name: "John Doe"})
    }))

    // Typed handler
    s.POST("/users", helix.Handle(func(ctx context.Context, req struct {
        Name string `json:"name"`
    }) (User, error) {
        return User{ID: 1, Name: req.Name}, nil
    }))

    s.Start(":8080")
}
```
