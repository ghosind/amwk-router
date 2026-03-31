package router

import "errors"

var (
	ErrRouteConflict = errors.New("router: a route with the same path and method already exists")
)
