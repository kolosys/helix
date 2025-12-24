# BasicAuth Middleware

Provides HTTP Basic Authentication. Uses constant-time comparison to prevent timing attacks.

## Basic Usage

```go
// Single user
s.Use(middleware.BasicAuth("admin", "secret"))

// Multiple users
s.Use(middleware.BasicAuthUsers(map[string]string{
    "admin": "secret",
    "user":  "password",
}))
```

## Configuration

```go
s.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
    Validator: func(username, password string) bool {
        // Validate against database
        return validateUser(username, password)
    },
    Realm: "Restricted Area",
    SkipFunc: func(r *http.Request) bool {
        // Skip auth for public endpoints
        return r.URL.Path == "/public"
    },
}))
```

## Features

- Constant-time password comparison
- Configurable realm
- Custom validator function
- Skip function for public endpoints

## Single User

```go
s.Use(middleware.BasicAuth("admin", "secret"))
```

## Multiple Users

```go
s.Use(middleware.BasicAuthUsers(map[string]string{
    "admin":    "admin-secret",
    "user":     "user-password",
    "readonly": "readonly-pass",
}))
```

## Custom Validator

Validate against a database or external service:

```go
s.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
    Validator: func(username, password string) bool {
        user, err := userService.GetByUsername(username)
        if err != nil {
            return false
        }
        return userService.ValidatePassword(user, password)
    },
    Realm: "My Application",
}))
```

## Skip Authentication

Skip authentication for public endpoints:

```go
s.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
    Validator: validateUser,
    Realm:     "Restricted",
    SkipFunc: func(r *http.Request) bool {
        // Public endpoints
        publicPaths := []string{"/", "/health", "/public"}
        for _, path := range publicPaths {
            if r.URL.Path == path {
                return true
            }
        }
        return false
    },
}))
```

## Realm

The realm is displayed in the browser's authentication dialog:

```go
s.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
    Validator: validateUser,
    Realm:     "Admin Area - Please Login",
}))
```

## Security Notes

- The middleware uses constant-time comparison to prevent timing attacks
- Passwords are compared using SHA-256 hashing and `subtle.ConstantTimeCompare`
- Never log passwords or credentials
- Use HTTPS in production to protect credentials in transit

## Example

```go
s := helix.New(nil)

// Protect admin routes
admin := s.Group("/admin", middleware.BasicAuth("admin", "secret"))
admin.GET("/users", listUsers)
admin.POST("/users", createUser)

// Protect API with custom validator
api := s.Group("/api", middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
    Validator: func(username, password string) bool {
        return validateAPIKey(username, password)
    },
    Realm: "API Access",
}))
```
