package handlers

import (
	"net/http"
)

// RouteGroup allows grouping routes with a common prefix and middleware
type RouteGroup struct {
	prefix     string
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler
}

// CreateRouteGroup creates a new route group with the given prefix
func CreateRouteGroup(mux *http.ServeMux, prefix string) *RouteGroup {
	return &RouteGroup{
		prefix: prefix,
		mux:    mux,
	}
}

// Use adds middleware to the route group
func (rg *RouteGroup) Use(middleware func(http.Handler) http.Handler) {
	rg.middleware = append(rg.middleware, middleware)
}

// Get adds a GET route to the group
func (rg *RouteGroup) Get(path string, handler http.HandlerFunc) {
	rg.handle(http.MethodGet, path, handler)
}

// Post adds a POST route to the group
func (rg *RouteGroup) Post(path string, handler http.HandlerFunc) {
	rg.handle(http.MethodPost, path, handler)
}

// Put adds a PUT route to the group
func (rg *RouteGroup) Put(path string, handler http.HandlerFunc) {
	rg.handle(http.MethodPut, path, handler)
}

// Delete adds a DELETE route to the group
func (rg *RouteGroup) Delete(path string, handler http.HandlerFunc) {
	rg.handle(http.MethodDelete, path, handler)
}

// handle adds a route with the given method to the group
func (rg *RouteGroup) handle(method, path string, handler http.HandlerFunc) {
	// Apply middleware chain in reverse order
	wrappedHandler := http.Handler(handler)
	for i := len(rg.middleware) - 1; i >= 0; i-- {
		wrappedHandler = rg.middleware[i](wrappedHandler)
	}

	// Register the route with method and path
	fullPath := rg.prefix + path
	rg.mux.HandleFunc(fullPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		wrappedHandler.ServeHTTP(w, r)
	})
}
