package router

import (
	"strings"

	"github.com/go-amwk/core"
)

// node represents a node in the routing tree, containing the path segment, child nodes, and
// handlers
type node struct {
	// path is the path segment for this node (e.g., "users", ":id", etc.)
	path string
	// children maps path segments to their corresponding child nodes
	children map[string]*node
	// handlers is the list of handler functions to execute when this node is matched
	handlers []core.HandlerFunc
}

// newNode creates a new node with the given path segment.
func newNode(path string) *node {
	return &node{
		path:     path,
		children: make(map[string]*node),
	}
}

// addRoute adds a new route to the node with the given path and handlers.
func (n *node) addRoute(path string, handlers ...core.HandlerFunc) error {
	if path == "/" {
		// root node
		if len(n.handlers) > 0 {
			return ErrRouteConflict
		}
		n.handlers = handlers
		return nil
	}

	segments := splitPathToSegments(path)
	cur := n
	for _, segment := range segments {
		child, ok := cur.children[segment]
		if !ok {
			child = newNode(segment)
			cur.children[segment] = child
		}
		cur = child
	}

	if cur.handlers != nil {
		return ErrRouteConflict
	}
	cur.handlers = handlers

	return nil
}

// find searches for a matching route in the node based on the given path and returns the
// corresponding handlers, path parameters, and a boolean indicating whether a match was found.
func (n *node) find(path string) ([]core.HandlerFunc, map[string]string, bool) {
	if n == nil {
		return nil, nil, false
	}
	if path == "" || path == "/" {
		return n.handlers, nil, len(n.handlers) > 0
	}

	segments := splitPathToSegments(path)
	cur := n
	params := make(map[string]string)

	for _, segment := range segments {
		child, ok := cur.children[segment]
		if !ok {
			return nil, nil, false
		}
		cur = child
	}

	if len(cur.handlers) > 0 {
		return cur.handlers, params, true
	}

	return nil, nil, false
}

// splitPathToSegments splits the given path into segments and returns a slice of segments.
func splitPathToSegments(path string) []string {
	if path == "/" {
		return []string{}
	}

	if path[0] == '/' {
		path = path[1:]
	}

	if path != "" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}
