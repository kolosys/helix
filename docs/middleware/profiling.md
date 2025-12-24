# Profiling Middleware

Profiles middleware performance (requires `profile` build tag). Tracks duration, allocations, and memory usage.

## Usage

Build with profiling enabled:

```bash
go build -tags profile
```

Wrap middleware with profiling:

```go
s.Use(middleware.ProfileMiddleware("requestid", middleware.RequestID()))
s.Use(middleware.ProfileMiddleware("logger", middleware.Logger(...)))
```

Access profiles:

```go
profiles := middleware.GetProfiles()
for name, profile := range profiles {
    fmt.Printf("%s: %v, %d allocs, %d bytes\n",
        name, profile.Duration, profile.Allocs, profile.Bytes)
}

// Reset profiles
middleware.ResetProfiles()
```

## Features

- Tracks execution time
- Tracks memory allocations
- Tracks bytes allocated
- Thread-safe profiling data

## Build Tag

Profiling middleware is only available when built with the `profile` tag:

```bash
# Build with profiling
go build -tags profile

# Run tests with profiling
go test -tags profile ./...
```

## Profiling Middleware

Wrap any middleware to profile it:

```go
s.Use(middleware.ProfileMiddleware("requestid", middleware.RequestID()))
s.Use(middleware.ProfileMiddleware("logger", middleware.Logger(middleware.LogFormatJSON)))
s.Use(middleware.ProfileMiddleware("recover", middleware.Recover()))
```

## Accessing Profiles

Get all profiles:

```go
profiles := middleware.GetProfiles()

for name, profile := range profiles {
    fmt.Printf("Middleware: %s\n", name)
    fmt.Printf("  Total Duration: %v\n", profile.Duration)
    fmt.Printf("  Total Allocations: %d\n", profile.Allocs)
    fmt.Printf("  Total Bytes: %d\n", profile.Bytes)
}
```

## Profile Data

Each profile contains:

- **Name**: Middleware name
- **Duration**: Total execution time
- **Allocs**: Total number of allocations
- **Bytes**: Total bytes allocated

## Resetting Profiles

Clear all profiling data:

```go
middleware.ResetProfiles()
```

## Example

```go
package main

import (
    "fmt"
    "github.com/kolosys/helix"
    "github.com/kolosys/helix/middleware"
)

func main() {
    s := helix.New(nil)

    // Profile middleware
    s.Use(middleware.ProfileMiddleware("requestid", middleware.RequestID()))
    s.Use(middleware.ProfileMiddleware("logger", middleware.Logger(middleware.LogFormatJSON)))
    s.Use(middleware.ProfileMiddleware("recover", middleware.Recover()))

    // ... setup routes ...

    // Periodically print profiles
    go func() {
        for {
            time.Sleep(30 * time.Second)
            profiles := middleware.GetProfiles()
            for name, profile := range profiles {
                fmt.Printf("%s: %v, %d allocs, %d bytes\n",
                    name, profile.Duration, profile.Allocs, profile.Bytes)
            }
        }
    }()

    s.Start(":8080")
}
```

## Use Cases

- Performance analysis during development
- Identifying middleware bottlenecks
- Memory allocation tracking
- Performance regression testing

## Important Notes

- Profiling adds overhead - only use in development/testing
- Requires `profile` build tag
- Profile data accumulates over time - reset periodically
- Thread-safe but may impact performance under high load
