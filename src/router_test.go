package tsweb

import (
	"reflect"
	"testing"
)

// newTestRouter creates and initializes a test router with predefined routes for testing purposes.
func newTestRouter() *Router {
	r := newRouter()
	r.addRoute("GET", "/", nil, nil)
	r.addRoute("GET", "/hello/:name", nil, nil)
	r.addRoute("GET", "/hello/b/c", nil, nil)
	r.addRoute("GET", "/hi/:name", nil, nil)
	r.addRoute("GET", "/assets/*filepath", nil, nil)
	return r
}

// TestParsePattern tests the functionality of the parsePattern function.
func TestParsePattern(t *testing.T) {
	// Test parsing a pattern with a wildcard segment
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	// Test parsing a pattern with a wildcard segment at the end
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	// Test parsing a pattern with multiple wildcard segments
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

// TestGetRoute tests the functionality of the getRoute method in retrieving routes from the router.
func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/test")

	// Check if route node is found
	if n == nil {
		t.Fatal("Test Failed")
	}

	// Check if the matched pattern is correct
	if n.pattern != "/hello/:name" {
		t.Fatal("Test Failed")
	}

	// Check if URL parameters are extracted correctly
	if ps["name"] != "test" {
		t.Fatal("Test Failed")
	}
	t.Log("Test success")
}

// TestRouter_AddRoute tests the functionality of adding routes to the router.
func TestRouter_AddRoute(t *testing.T) {
	router := newRouter()
	handler := func(c *Context) {}

	// Test adding a route with a simple pattern
	router.addRoute("GET", "/hello", handler, nil)
	if len(router.roots["GET"].children) != 1 {
		t.Error("Failed to add route with simple pattern")
	}

	// Test adding a route with wildcard segment
	router.addRoute("POST", "/users/:id", handler, nil)
	if len(router.roots["POST"].children) != 1 {
		t.Error("Failed to add route with wildcard segment")
	}

	// Test adding a route with wildcard segment followed by a static segment
	router.addRoute("PUT", "/files/*filepath", handler, nil)
	if len(router.roots["PUT"].children) != 1 {
		t.Error("Failed to add route with wildcard segment followed by a static segment")
	}
}

// TestRouter_GetRoute tests the functionality of retrieving routes from the router.
func TestRouter_GetRoute(t *testing.T) {
	router := newRouter()
	handler := func(c *Context) {}

	// Add routes for testing
	router.addRoute("GET", "/hello", handler, nil)
	router.addRoute("POST", "/users/:id", handler, nil)
	router.addRoute("PUT", "/files/*filepath", handler, nil)

	// Test getting a route with exact match
	node, params := router.getRoute("GET", "/hello")
	if node == nil || node.pattern != "/hello" || len(params) != 0 {
		t.Error("Failed to get route with exact match")
	}

	// Test getting a route with wildcard segment
	node, params = router.getRoute("POST", "/users/123")
	if node == nil || node.pattern != "/users/:id" || params["id"] != "123" {
		t.Error("Failed to get route with wildcard segment")
	}

	// Test getting a route with wildcard segment followed by a static segment
	node, params = router.getRoute("PUT", "/files/dir/subdir/file.txt")
	if node == nil || node.pattern != "/files/*filepath" || params["filepath"] != "dir/subdir/file.txt" {
		t.Error("Failed to get route with wildcard segment followed by a static segment")
	}

	// Test getting a non-existent route
	node, params = router.getRoute("GET", "/non-existent")
	if node != nil || params != nil {
		t.Error("Failed to handle non-existent route")
	}
}
