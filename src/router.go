package tsweb

import (
	"net/http"
	"strings"
)

// HandlerFunc defines the function signature for request handlers.
type HandlerFunc func(*Context)

// Router defines a struct responsible for routing HTTP requests to appropriate handler functions.
type Router struct {
	roots                 map[string]*node        // roots stores the root nodes for each HTTP method
	handlerMap            map[string]HandlerFunc  // handlerMap stores the handler functions mapped to HTTP methods and patterns
	handlerRouterGroupMap map[string]*RouterGroup // handlerRouterGroupMap stores the router groups mapped to HTTP methods and patterns
}

// newRouter creates and returns a new Router instance.
func newRouter() *Router {
	return &Router{
		roots:                 make(map[string]*node),
		handlerMap:            make(map[string]HandlerFunc),
		handlerRouterGroupMap: make(map[string]*RouterGroup),
	}
}

// parsePattern parses the URL pattern and returns its segments.
func parsePattern(pattern string) []string {
	values := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range values {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute adds a route to the router for the specified HTTP method, pattern, handler, and router group.
func (p *Router) addRoute(method string, pattern string, handler HandlerFunc, routerGroup *RouterGroup) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := p.roots[method]
	if !ok {
		p.roots[method] = &node{}
	}
	p.roots[method].insert(pattern, parts, 0)
	p.handlerMap[key] = handler
	p.handlerRouterGroupMap[key] = routerGroup
}

// getRoute retrieves the route matching the HTTP method and path, and extracts any URL parameters.
func (p *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := p.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// handle processes the incoming HTTP request by matching the route and invoking the appropriate handler.
func (p *Router) handle(c *Context) {
	n, params := p.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handle = p.handlerMap[key]
		c.middlewares = &p.handlerRouterGroupMap[key].middlewares
		c.Next()
	} else {
		c.String(http.StatusNotFound, "404")
	}
}
