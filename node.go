package router

import "github.com/go-amwk/core"

type node struct {
	path     string
	children []*node
}

// addRoute adds a new route to the node with the given path and handlers.
func (n *node) addRoute(path string, handlers ...core.HandlerFunc) {
	// TODO
}

// find searches for a matching route in the node based on the given path and returns the
// corresponding handlers, path parameters, and a boolean indicating whether a match was found.
func (h *node) find(path string) ([]core.HandlerFunc, map[string]string, bool) {
	// TODO
	return nil, nil, false
}
