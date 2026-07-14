package router

import (
	"net/http"
	"sync"

	"github.com/go-amwk/core"
)

// Router is the routing engine that manages the routing tree and handles incoming requests by
// matching them to the registered routes and executing the corresponding handlers.
type Router struct {
	// NotFoundHandler is an optional custom handler that will be executed when no matching route is
	// found.
	NotFoundHandler core.HandlerFunc
	// tree is the routing tree that maps HTTP methods to their corresponding root nodes.
	tree map[string]*node
	// mu is a mutex to protect concurrent access to the routing tree during route registration and
	// request handling.
	mu sync.RWMutex
}

// New creates a new Router instance.
func New() *Router {
	r := new(Router)
	r.tree = make(map[string]*node)

	return r
}

// Handle registers a new route with the given method, path, and handlers.
func (r *Router) Handle(method, path string, handlers ...core.HandlerFunc) *Router {
	r.handle(method, path, handlers...)

	return r
}

// GET registers a new GET route with the given path and handlers.
func (r *Router) GET(path string, handlers ...core.HandlerFunc) *Router {
	return r.Handle(http.MethodGet, path, handlers...)
}

// POST registers a new POST route with the given path and handlers.
func (r *Router) POST(path string, handlers ...core.HandlerFunc) *Router {
	return r.Handle(http.MethodPost, path, handlers...)
}

// PUT registers a new PUT route with the given path and handlers.
func (r *Router) PUT(path string, handlers ...core.HandlerFunc) *Router {
	return r.Handle(http.MethodPut, path, handlers...)
}

// DELETE registers a new DELETE route with the given path and handlers.
func (r *Router) DELETE(path string, handlers ...core.HandlerFunc) *Router {
	return r.Handle(http.MethodDelete, path, handlers...)
}

// PATCH registers a new PATCH route with the given path and handlers.
func (r *Router) PATCH(path string, handlers ...core.HandlerFunc) *Router {
	return r.Handle(http.MethodPatch, path, handlers...)
}

// OPTIONS registers a new OPTIONS route with the given path and handlers.
func (r *Router) OPTIONS(path string, handlers ...core.HandlerFunc) *Router {
	return r.Handle(http.MethodOptions, path, handlers...)
}

// HEAD registers a new HEAD route with the given path and handlers.
func (r *Router) HEAD(path string, handlers ...core.HandlerFunc) *Router {
	return r.Handle(http.MethodHead, path, handlers...)
}

// Any registers a new route for all supported HTTP methods with the given path and handlers.
func (r *Router) Any(path string, handlers ...core.HandlerFunc) *Router {
	for _, method := range []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
		http.MethodHead,
	} {
		r.Handle(method, path, handlers...)
	}

	return r
}

// handle is a helper method to register a new route with the given method, path, and handlers.
func (r *Router) handle(method, path string, handlers ...core.HandlerFunc) {
	if len(handlers) == 0 {
		// at least one handler must be provided for a route, otherwise panic with an error
		panic(ErrNoHandlerProvided)
	}

	if path == "" || path[0] != '/' {
		path = "/" + path
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	root, ok := r.tree[method]
	if !ok {
		root = newNode("")
		r.tree[method] = root
	}

	if err := root.addRoute(path, handlers...); err != nil {
		panic(err)
	}
}

// Route returns a handler middleware that matches the incoming request's method and path to the
// registered routes and executes the corresponding handlers.
func (r *Router) Route() core.HandlerFunc {
	return func(ctx core.Context) error {
		method := ctx.Method()
		path := ctx.Path()

		r.mu.RLock()
		defer r.mu.RUnlock()

		root, ok := r.tree[method]
		if ok {
			if handlers, params, found := root.find(path); found {
				for key, value := range params {
					ctx.Request().SetPathValue(key, value)
				}
				ctx.Use(handlers...)
				return nil
			}
		}

		ctx.Use(r.notFoundHandler())
		return nil
	}
}

// notFoundHandler returns the custom NotFoundHandler if set, otherwise it returns the default 404
// handler.
func (r *Router) notFoundHandler() core.HandlerFunc {
	if r.NotFoundHandler != nil {
		return r.NotFoundHandler
	}
	return DefaultNotFoundHandler()
}

// DefaultNotFoundHandler returns a default 404 Not Found handler.
func DefaultNotFoundHandler() core.HandlerFunc {
	return func(ctx core.Context) error {
		ctx.Status(http.StatusNotFound)
		ctx.Write([]byte("404 Not Found"))
		ctx.Abort()
		return nil
	}
}
