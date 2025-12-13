# Overview

<p align="center">

## About helix

This documentation provides comprehensive guidance for using helix, a Go library designed to help you build better software.

## Project Information

- **Repository**: [https://github.com/kolosys/helix](https://github.com/kolosys/helix)
- **Import Path**: `github.com/kolosys/helix`
- **License**: MIT
- **Version**: latest

## What You'll Find Here

This documentation is organized into several sections to help you find what you need:

- **[Getting Started](../getting-started/)** - Installation instructions and quick start guides
- **[Core Concepts](../core-concepts/)** - Fundamental concepts and architecture details
- **[Advanced Topics](../advanced/)** - Performance tuning and advanced usage patterns
- **[API Reference](../api-reference/)** - Complete API reference documentation
- **[Examples](../examples/)** - Working code examples and tutorials

## Project Features

helix provides:
- **helix** - Package helix provides a zero-dependency, context-aware, high-performance
HTTP web framework for Go with stdlib compatibility.

- **main** - Package main demonstrates the most basic usage of the helix framework.

- **main** - Package main demonstrates a full CRUD API using typed handlers.

- **main** - Package main demonstrates route grouping and API versioning.

- **main** - Package main demonstrates the use of middleware in helix.

- **main** - Package main demonstrates advanced helix features:
- Modular route organization
- Service registration and dependency injection
- Pagination helpers
- Health check endpoints

- **main** - Package main demonstrates the RESTful resource builder pattern.

- **main** - Package main demonstrates request binding and validation.

- **logs** - Package logs provides a high-performance, context-aware structured logging library.

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

- **middleware** - Package middleware provides HTTP middleware for the Helix framework.


## Quick Links

- [Installation Guide](installation.md)
- [Quick Start Guide](quick-start.md)
- [API Reference](../api-reference/)
- [Examples](../examples/README.md)

## Community & Support

- **GitHub Issues**: [https://github.com/kolosys/helix/issues](https://github.com/kolosys/helix/issues)
- **Discussions**: [https://github.com/kolosys/helix/discussions](https://github.com/kolosys/helix/discussions)
- **Repository Owner**: [kolosys](https://github.com/kolosys)

## Getting Help

If you encounter any issues or have questions:

1. Check the [API Reference](../api-reference/) for detailed documentation
2. Browse the [Examples](../examples/README.md) for common use cases
3. Search existing [GitHub Issues](https://github.com/kolosys/helix/issues)
4. Open a new issue if you've found a bug or have a feature request

## Next Steps

Ready to get started? Head over to the [Installation Guide](installation.md) to begin using helix.

