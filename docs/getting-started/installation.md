# Installation

This guide will help you install and set up Helix in your Go project.

## Prerequisites

Before installing Helix, ensure you have:

- **Go 1.24+** or later installed
- A Go module initialized in your project (run `go mod init` if needed)
- Access to the GitHub repository (for private repositories)

## Installation Steps

### Step 1: Install the Package

Use `go get` to install Helix:

```bash
go get github.com/kolosys/helix
```

This will download the package and add it to your `go.mod` file.

### Step 2: Import in Your Code

Import the package in your Go source files:

```go
import "github.com/kolosys/helix"
```

### Available Packages

Helix includes several packages. Import the ones you need:

#### Main Package

The main `helix` package provides the HTTP web framework:

```go
// Package helix provides a zero-dependency, context-aware, high-performance
// HTTP web framework for Go with stdlib compatibility.
import "github.com/kolosys/helix"
```

#### Middleware Package

The `middleware` package provides HTTP middleware:

```go
// Package middleware provides HTTP middleware for the Helix framework.
import "github.com/kolosys/helix/middleware"
```

### Step 3: Verify Installation

Create a simple test file to verify the installation:

```go
package main

import (
    "fmt"
    "github.com/kolosys/helix"
)

func main() {
    fmt.Println("Helix installed successfully!")
    fmt.Printf("Helix version: %s\n", helix.Version)
}
```

Run the test:

```bash
go run main.go
```

## Updating the Package

To update to the latest version:

```bash
go get -u github.com/kolosys/helix
```

To update to a specific version:

```bash
go get github.com/kolosys/helix@v1.2.3
```

## Installing a Specific Version

To install a specific version of the package:

```bash
go get github.com/kolosys/helix@v1.0.0
```

Check available versions on the [GitHub releases page](https://github.com/kolosys/helix/releases).

## Development Setup

If you want to contribute or modify the library:

### Clone the Repository

```bash
git clone https://github.com/kolosys/helix.git
cd helix
```

### Install Dependencies

```bash
go mod download
```

### Run Tests

```bash
go test ./...
```

Run tests with race detection:

```bash
go test -race ./...
```

## Troubleshooting

### Module Not Found

If you encounter a "module not found" error:

1. Ensure you're using Go modules (run `go mod init` if needed)
2. Check that you have network access to GitHub
3. Try running `go clean -modcache` and reinstall:
   ```bash
   go clean -modcache
   go get github.com/kolosys/helix
   ```

### Private Repository Access

For private repositories, configure Git to use SSH or a personal access token:

```bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

Or set up `GOPRIVATE`:

```bash
export GOPRIVATE=github.com/kolosys/helix
```

### Version Compatibility

If you encounter compatibility issues:

1. Check your Go version: `go version`
2. Ensure you're using Go 1.24 or later
3. Update Go if necessary: Visit [golang.org](https://golang.org/dl/)
