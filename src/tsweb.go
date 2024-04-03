package tsweb

import (
	"log"
	"net/http"
	"path"
	"text/template"
)

// RouterGroup represents a group of routes with a common prefix and middleware.
type RouterGroup struct {
	prefix      string            // Common prefix for routes in this group.
	middlewares []HandlerFunc     // Middleware handlers for this group.
	parent      *RouterGroup      // Parent router group.
	engine      *Engine           // Associated engine.
	filePathMap map[string]string // Mapping of file paths.
}

// Engine is the web framework engine.
type Engine struct {
	*RouterGroup                     // Embedding RouterGroup for convenience.
	router        *Router            // Router for handling HTTP requests.
	htmlTemplates *template.Template // HTML template renderer.
	funcMap       template.FuncMap   // FuncMap for HTML templates.
}

// NewEngine creates a new Engine instance with an initialized router.
func NewEngine() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{
		engine:      engine,
		prefix:      "",
		parent:      nil,
		middlewares: make([]HandlerFunc, 0),
		filePathMap: make(map[string]string),
	}
	return engine
}

// Group creates a new RouterGroup with the given prefix.
func (r *RouterGroup) Group(prefix string) *RouterGroup {
	routerGroup := &RouterGroup{
		prefix:      r.prefix + prefix,
		parent:      r,
		engine:      r.engine,
		middlewares: make([]HandlerFunc, len(r.middlewares)),
		filePathMap: make(map[string]string),
	}

	if len(r.middlewares) > 0 {
		copy(routerGroup.middlewares, r.middlewares)
	}
	return routerGroup
}

// Use adds middleware handlers to the RouterGroup.
func (r *RouterGroup) Use(handlerFunc HandlerFunc) {
	r.middlewares = append(r.middlewares, handlerFunc)
}

// GET registers a GET request handler for the given URL pattern.
func (r *RouterGroup) GET(url string, handlerFunc HandlerFunc) {
	r.addRoute("GET", url, handlerFunc)
}

// POST registers a POST request handler for the given URL pattern.
func (r *RouterGroup) POST(url string, handlerFunc HandlerFunc) {
	r.addRoute("POST", url, handlerFunc)
}

// addRoute registers a request handler for the given HTTP method and URL pattern.
func (r *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := r.prefix + comp
	r.engine.router.addRoute(method, pattern, handler, r)
}

// createStaticHandler creates a handler function for serving static files.
func (r *RouterGroup) createStaticHandler(filePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(r.prefix, filePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static registers a handler for serving static files.
func (r *RouterGroup) Static(url string, filePath string) {
	handler := r.createStaticHandler(url, http.Dir(filePath))
	pattern := path.Join(url, "/*filepath")
	r.GET(pattern, handler)
}

// GET registers a GET request handler with the Engine's router.
func (p *Engine) GET(url string, handlerFunc HandlerFunc) {
	p.router.addRoute("GET", url, handlerFunc, p.RouterGroup)
}

// POST registers a POST request handler with the Engine's router.
func (p *Engine) POST(url string, handlerFunc HandlerFunc) {
	p.router.addRoute("POST", url, handlerFunc, p.RouterGroup)
}

// ServeHTTP handles HTTP requests by passing them to the router.
func (p *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := makeContext(w, req, p)
	p.router.handle(c)
}

// SetFuncMap sets the FuncMap for HTML templates.
func (p *Engine) SetFuncMap(funcMap template.FuncMap) {
	p.funcMap = funcMap
}

// LoadHTMLGlob loads HTML templates from the specified pattern.
func (p *Engine) LoadHTMLGlob(pattern string) {
	p.htmlTemplates = template.Must(template.New("").Funcs(p.funcMap).ParseGlob(pattern))
}

// Run starts the HTTP server and listens for incoming requests on the specified port.
func (p *Engine) Run(port string) {
	http.ListenAndServe(port, p)
}

// Logger is a middleware handler that logs the start and end of each request.
func Logger() HandlerFunc {
	return func(c *Context) {
		log.Printf("TSWeb start")
		c.Next()
		log.Printf("TSWeb end")
	}
}

// Recovery is a middleware handler that recovers from panics and returns an appropriate error response.
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				log.Printf("Panic: %v", err)

				// Set appropriate HTTP status code
				c.Status(http.StatusInternalServerError)

				// Return specific error message in response
				c.JSON(http.StatusInternalServerError, H{"error": "Internal Server Error"})
			}
		}()

		c.Next()
	}
}
