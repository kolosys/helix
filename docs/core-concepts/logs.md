# logs

> **Note**: This is a developer-maintained documentation page. The content here is not auto-generated and should be updated manually to explain the core concepts and architecture of the logs package.

## About This Package

**Import Path:** `github.com/kolosys/helix/logs`

Package logs provides a high-performance, context-aware structured logging library.

Features:
  - Zero-allocation hot paths using sync.Pool
  - Context-aware logging with context.Context
  - Type-safe field builders
  - Multiple output formats (text, JSON, pretty)
  - Sampling for high-volume logs
  - Async logging option
  - Hook system for extensibility
  - Built-in caller information
  - Chained/fluent API

Basic usage:

	log := logs.New()
	log.Info("server started", logs.Int("port", 8080))

With context:

	log.InfoContext(ctx, "request processed", logs.Duration("latency", time.Since(start)))


## Architecture Overview

<!-- Add information about the package architecture here -->

This section should explain:
- The main design patterns used in this package
- How components interact with each other
- The data flow through the package
- Key interfaces and their purposes

## Core Concepts

<!-- Document the fundamental concepts developers need to understand -->

### Concept 1

Explain the first core concept here.

### Concept 2

Explain the second core concept here.

## Design Decisions

<!-- Explain important design decisions and trade-offs -->

Document why certain approaches were chosen:
- Performance considerations
- API design choices
- Backward compatibility decisions

## Usage Patterns

<!-- Show common usage patterns and idioms -->

### Pattern 1: Basic Usage

```go
// Example code here
```

### Pattern 2: Advanced Usage

```go
// Example code here
```

## Common Pitfalls

<!-- Document common mistakes and how to avoid them -->

- Pitfall 1: Description and solution
- Pitfall 2: Description and solution

## Integration Guide

<!-- How this package integrates with other packages or systems -->

Explain how this package works with:
- Other packages in this library
- External dependencies
- Common frameworks or tools

## Further Reading

- [API Reference](../api-reference/logs.md) - Complete API documentation
- [Examples](../examples/README.md) - Practical examples
- [Best Practices](../advanced/best-practices.md) - Recommended patterns

---

*This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used.*

