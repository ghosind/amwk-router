# router

A lightweight and fast HTTP router for the [amwk](https://github.com/go-amwk) web framework.

## Installation

**Requirements:** Go 1.22+

```bash
go get github.com/go-amwk/router
```

## Getting Started

There are a simple example demonstrating how to use the `router` package to create a basic web server with route handling.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-amwk/core"
	"github.com/go-amwk/router"
	"github.com/go-amwk/web"
)

func main() {
	app := web.Default()
	r := router.New()

	// Register routes
	r.GET("/hello", func(ctx core.Context) error {
		return ctx.Write([]byte("Hello, World!"))
	})

	r.POST("/users", func(ctx core.Context) error {
		return ctx.Write([]byte("User created"))
	})

	// Mount the router as middleware
	app.Use(r.Route())

	app.Start()
}
```

## Route Registration

### HTTP Method Helpers

Convenience methods are provided for all standard HTTP methods. Each returns the `*Router` instance for chaining:

```go
r := router.New()

r.GET("/users", listUsers)
r.POST("/users", createUser)
r.PUT("/users/:id", updateUser)
r.DELETE("/users/:id", deleteUser)
r.PATCH("/users/:id", patchUser)
r.OPTIONS("/users", optionsHandler)
r.HEAD("/users", headHandler)
```

### Catch-All Routes

`Any()` registers the same handlers for all standard HTTP methods (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD) in one call:

```go
r.Any("/ping", func(ctx core.Context) error {
	ctx.Write([]byte("pong"))
  return nil
})
```

### Chaining

All registration methods return `*Router`, so you can chain route definitions:

```go
r := router.New().
	GET("/api/v1/users", listUsers).
	POST("/api/v1/users", createUser).
	GET("/api/v1/users/:id", getUser).
	PUT("/api/v1/users/:id", updateUser)
```

## Handler Chains

Each route accepts multiple handler functions, executed in order via `ctx.Next()`:

```go
func authMiddleware(ctx core.Context) error {
	token := ctx.Header("Authorization")
	if token == "" {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return nil
	}
	return ctx.Next()
}

func logMiddleware(ctx core.Context) error {
	fmt.Printf("[%s] %s\n", ctx.Method(), ctx.Path())
	return ctx.Next()
}

r.GET("/protected", authMiddleware, logMiddleware, func(ctx core.Context) error {
	ctx.Write([]byte("Secret data"))
  return nil
})
```

## Custom Not-Found Handler

Set `NotFoundHandler` on the router to customize the 404 response:

```go
r.NotFoundHandler = func(ctx core.Context) error {
	ctx.Status(http.StatusNotFound)
	ctx.Write([]byte(`{"error": "route not found"}`))
  return nil
}
```

When not set, the router returns a default `404 Not Found` with an aborted context.

## Route Conflicts

Registering the same method + path combination twice panics with `ErrRouteConflict`. This is intentional — the framework treats duplicate routes as a programming error at startup time:

```go
r.GET("/api", handler1)
r.GET("/api", handler2) // panics: router: a route with the same path and method already exists
```

## API Reference

### `func New() *Router`

Creates a new `Router` instance with an empty routing tree.

### `Router` Type

| Method | Description |
|--------|-------------|
| `Handle(method, path string, handlers ...core.HandlerFunc) *Router` | Registers a route for an arbitrary HTTP method |
| `GET(path string, handlers ...core.HandlerFunc) *Router` | Registers a GET route |
| `POST(path string, handlers ...core.HandlerFunc) *Router` | Registers a POST route |
| `PUT(path string, handlers ...core.HandlerFunc) *Router` | Registers a PUT route |
| `DELETE(path string, handlers ...core.HandlerFunc) *Router` | Registers a DELETE route |
| `PATCH(path string, handlers ...core.HandlerFunc) *Router` | Registers a PATCH route |
| `OPTIONS(path string, handlers ...core.HandlerFunc) *Router` | Registers an OPTIONS route |
| `HEAD(path string, handlers ...core.HandlerFunc) *Router` | Registers a HEAD route |
| `Any(path string, handlers ...core.HandlerFunc) *Router` | Registers handlers for all standard methods |
| `Route() core.HandlerFunc` | Returns a middleware that performs request routing |

| Field | Description |
|-------|-------------|
| `NotFoundHandler core.HandlerFunc` | Optional custom handler for unmatched routes |

### `func DefaultNotFoundHandler() core.HandlerFunc`

Returns the default 404 handler (sets status to 404 and aborts the context).

## License

The project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.
