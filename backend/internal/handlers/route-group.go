package handlers

import (
	"fmt"
	"net/http"
)

const (
	getMethod    = "GET"
	postMethod   = "POST"
	putMethod    = "PUT"
	deleteMethod = "DELETE"
)

// RouteGroup represents a collection of routes with a common prefix and shared middleware
type RouteGroup struct {
	mux         *http.ServeMux
	prefix      string
	middlewares []func(http.Handler) http.Handler
}

// CreateRouteGroup creates a new route group with the specified prefix
func CreateRouteGroup(mux *http.ServeMux, prefix string) *RouteGroup {
	return &RouteGroup{
		mux:         mux,
		prefix:      prefix,
		middlewares: []func(http.Handler) http.Handler{},
	}
}

// Use adds middleware to the route group
func (rg *RouteGroup) Use(middleware func(http.Handler) http.Handler) {
	rg.middlewares = append(rg.middlewares, middleware)
}

// Handle registers a handler for the specified pattern with applied middlewares
func (rg *RouteGroup) Handle(method, pattern string, handler http.Handler) {
	fullPattern := method + " " + rg.prefix + pattern

	// Apply middlewares in reverse order
	for i := len(rg.middlewares) - 1; i >= 0; i-- {
		handler = rg.middlewares[i](handler)
	}

	rg.mux.Handle(fullPattern, handler)

	fmt.Printf("Registered route: %s\n", fullPattern)
}

// HandleFunc registers a handler function for the specified pattern
func (rg *RouteGroup) HandleFunc(method, pattern string, handler http.HandlerFunc) {
	rg.Handle(method, pattern, handler)
}

// Get registers a GET handler for the specified pattern
func (rg *RouteGroup) Get(pattern string, handler http.HandlerFunc) {
	rg.Handle(getMethod, pattern, handler)
}

// Post registers a POST handler for the specified pattern
func (rg *RouteGroup) Post(pattern string, handler http.HandlerFunc) {
	rg.Handle(postMethod, pattern, handler)
}

// Put registers a PUT handler for the specified pattern
func (rg *RouteGroup) Put(pattern string, handler http.HandlerFunc) {
	rg.Handle(putMethod, pattern, handler)
}

// Delete registers a DELETE handler for the specified pattern
func (rg *RouteGroup) Delete(pattern string, handler http.HandlerFunc) {
	rg.Handle(deleteMethod, pattern, handler)
}

// Group returns a new child route group with a concatenated prefix
func (rg *RouteGroup) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		mux:         rg.mux,
		prefix:      rg.prefix + prefix,
		middlewares: rg.middlewares, // Inherit parent middlewares
	}
}
