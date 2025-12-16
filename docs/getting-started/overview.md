# Helix 🧬

![GoVersion](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Zero Dependencies](https://img.shields.io/badge/Zero-Dependencies-green.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/kolosys/helix.svg)](https://pkg.go.dev/github.com/kolosys/helix)
[![Go Report Card](https://goreportcard.com/badge/github.com/kolosys/helix)](https://goreportcard.com/report/github.com/kolosys/helix)

```
    __ __    ___
   / // /__ / (_)_ __
  / _  / -_) / /\ \ /
 /_//_/\__/_/_//_\_\
 Developer friendly HTTP framework
```

**Helix** is a zero-dependency, high-performance HTTP web framework for Go with a focus on developer experience, type safety, and stdlib compatibility. Built by [Kolosys](https://github.com/kolosys) for enterprise-grade applications.

## Features

- **Zero Dependencies** - Built entirely on Go's standard library
- **High Performance** - Zero-allocation hot paths using `sync.Pool`
- **Type-Safe Handlers** - Generic handlers with automatic request binding and response encoding
- **RFC 7807 Problem Details** - Standardized error responses out of the box
- **Modular Architecture** - First-class support for organizing routes into modules
- **Fluent API** - Chainable context methods for clean handler code
- **Middleware Ecosystem** - Comprehensive built-in middleware suite
- **Dependency Injection** - Type-safe service registry with request-scoped support
- **Health Checks** - Built-in Kubernetes-ready liveness and readiness probes
- **Structured Logging** - High-performance logging with JSON and text formatters
- **Graceful Shutdown** - Context-aware shutdown with configurable grace period
- **stdlib Compatible** - Works with any `http.Handler` middleware
