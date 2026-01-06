# Go Types Library - AI Coding Agent Instructions

This document guides AI agents contributing to the Go types library, providing essential context for immediate productivity.

## Project Overview

This is a Go library providing generic data structures and a `Nullable` type for handling optional values. The goal is to offer robust, well-tested, and idiomatic building blocks for Go applications.

**Module Path**: `types` (Go 1.25+)

## Architecture & Key Components

### Core Packages
- **`collections/`**: Thread-safe generic data structures, each in its own subdirectory
  - Basic: `List`, `Queue`, `Stack`, `Set`, `Dictionary`
  - Advanced: `SortedSet`, `BitSet`, `BloomFilter`, `LRUCache`, `DisjointSet`, `SegmentTree`
  - Specialized variants: `queue/deque`, `queue/priority`, `queue/ring`, `dictionary/sorted`
- **`nullable/`**: Generic `Type[T]` for optional values with JSON/database integration
- **`cast/`**: Type conversion utilities with reflection-based casting
- **`sort/`**: Generic sorting helpers for slices
- **`internal/common/`**: Shared utilities like thread-safe `Vector[T]`

### Concurrency Pattern
All collection types are **thread-safe by default** using embedded mutexes:
- Use `sync.Mutex` for write-heavy operations (e.g., `DisjointSet`)
- Use `sync.RWMutex` for read-heavy operations (e.g., `BitSet`, `LRUCache`)
- Mutex is always embedded as `mutex` or `mu` field

## Critical Developer Workflows

### Testing
- **Per-package tests**: Each structure has `_test.go` files (e.g., `disjointset_test.go`)
- **Run all tests**: `go test ./...` or `make test`
- **Test structure**: Focus on edge cases, concurrency safety, and method contracts

### Linting & Formatting
- **Lint**: `make lint` or `golangci-lint run ./...`
- **Auto-fix**: `make lint-fix` 
- **Format**: Code must pass `gofmt` (handled automatically by IDE)

### Dependencies
- **Add/update deps**: Run `go mod tidy` after changes
- **No external deps**: Core collections avoid third-party dependencies for minimal footprint

## Project-Specific Patterns

### Generics Usage
- **All collections use generics** with appropriate constraints:
  - `comparable` for keys/maps/sets
  - `any` for values when no constraints needed
- **Example**:
  ```go
  // DisjointSet with comparable elements
  ds := disjointset.New[string]()
  
  // LRUCache with comparable keys, any values
  cache := lrucache.New[string, int](100)
  ```

### Nullable Type Integration
Use `nullable.Type[T]` for optional fields, especially with JSON/database:
```go
// Database model example
type User struct {
    ID    int                    `json:"id"`
    Name  nullable.Type[string]  `json:"name"`
    Email nullable.Type[string]  `json:"email"`
}

// Usage
user := User{
    ID: 1,
    Name: nullable.New("John"),
    Email: nullable.New[string](), // null value
}
```

### Error Handling
- **Panic on invalid usage**: Methods panic for programmer errors (e.g., `nullable.New()` with >1 arg)
- **Return errors for runtime issues**: Operations that can fail at runtime return `(result, error)`

### Internal Package Usage
- **`internal/common/mutex.go`** provides `RWLocker` interface and `Vector[T]` for shared mutex patterns
- Only use internal packages when implementing new collections that need consistent concurrency patterns

## Implementation Guidelines

### When Adding New Collections
1. Create subdirectory under `collections/` (e.g., `collections/fenwicktree/`)
2. Implement thread-safe methods with appropriate mutex strategy
3. Add comprehensive tests in `{name}_test.go`
4. Update `collections/README.md` with new structure
5. Follow existing naming conventions (`New()`, `Size()`, `IsEmpty()`, etc.)

### Code Style
- **Method names**: Use clear, descriptive names (e.g., `MakeSet`, `Union`, `FindRoot`)
- **Documentation**: Every public type/method must have Go doc comments
- **Performance**: Optimize for common cases while maintaining correctness
- **Memory**: Minimize allocations in hot paths; reuse buffers when possible

### Testing Requirements
- **Concurrency tests**: Include goroutine-based tests for thread safety
- **Edge cases**: Test empty states, boundary conditions, and error scenarios
- **Benchmark tests**: Add benchmarks for performance-critical operations

## Key Files for Reference
- **Concurrency pattern**: `internal/common/mutex.go`
- **Nullable implementation**: `nullable/type.go` 
- **Generic constraints**: Check individual collection files for constraint examples
- **Test patterns**: `collections/disjointset/disjointset_test.go`
- **Build workflow**: `Makefile` and `scripts/lint.ps1`

## Anti-Patterns to Avoid
- ❌ Non-thread-safe collections (all must be safe by default)
- ❌ External dependencies in core collections
- ❌ Inconsistent method naming across similar structures
- ❌ Missing test coverage for concurrent access scenarios
- ❌ Unnecessary memory allocations in performance-critical paths
