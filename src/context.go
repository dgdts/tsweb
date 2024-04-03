package tsweb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H is a shorthand for map[string]interface{} used to represent a generic map of data.
type H map[string]interface{}

// Context represents the context of an HTTP request.
type Context struct {
	Writer       http.ResponseWriter // Response writer for sending HTTP response
	Req          *http.Request       // HTTP request object
	Path         string              // Request path
	Method       string              // HTTP method (GET, POST, etc.)
	Params       map[string]string   // Parameters extracted from the request path
	StatusCode   int                 // HTTP status code to be sent in the response
	handle       HandlerFunc         // Handler function for processing the request
	middlewares  *[]HandlerFunc      // Slice of middleware functions to be executed
	processIndex int                 // Index to keep track of the current middleware being processed
	engine       *Engine             // Pointer to the Gee engine instance
}

// makeContext creates a new Context object.
func makeContext(w http.ResponseWriter, r *http.Request, engine *Engine) *Context {
	r.ParseForm()
	return &Context{
		Writer:       w,
		Req:          r,
		Path:         r.URL.Path,
		Method:       r.Method,
		processIndex: 0,
		engine:       engine,
	}
}

// HTML sends an HTML response with the specified status code, content, and data.
func (p *Context) HTML(status int, content string, data interface{}) {
	p.SetHeader("Content-Type", "text/html")
	p.Status(status)
	if err := p.engine.htmlTemplates.ExecuteTemplate(p.Writer, content, data); err != nil {
		p.String(500, err.Error())
	}
}

// String sends a plain text response with the specified status code and formatted content.
func (p *Context) String(status int, formatString string, contents ...interface{}) {
	p.SetHeader("Content-Type", "text/plain")
	p.Status(status)
	p.Writer.Write([]byte(fmt.Sprintf(formatString, contents...)))
}

// Query returns the value of the specified query parameter from the request URL.
func (p *Context) Query(key string) string {
	return p.Req.URL.Query().Get(key)
}

// PostForm returns the value of the specified form parameter from the HTTP POST body.
func (p *Context) PostForm(key string) string {
	return p.Req.FormValue(key)
}

// Status sets the HTTP status code for the response.
func (p *Context) Status(code int) {
	p.StatusCode = code
	p.Writer.WriteHeader(code)
}

// SetHeader sets the value of the specified HTTP header field.
func (p *Context) SetHeader(key string, value string) {
	p.Writer.Header().Set(key, value)
}

// Data sends raw data with the specified status code.
func (p *Context) Data(status int, data []byte) {
	p.Status(status)
	p.Writer.Write(data)
}

// Param returns the value of the specified path parameter from the request.
func (p *Context) Param(key string) string {
	return p.Params[key]
}

// Next proceeds to the next middleware in the chain.
func (p *Context) Next() {
	if p.processIndex >= len(*p.middlewares) {
		p.handle(p)
	} else {
		p.processIndex++
		(*p.middlewares)[p.processIndex-1](p)
	}
}

// JSON sends a JSON response with the specified status code and header data.
func (p *Context) JSON(status int, header map[string]interface{}) {
	p.SetHeader("Content-Type", "application/json")
	p.Status(status)
	value, err := json.Marshal(header)
	if err != nil {
		p.Error(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	_, err = p.Writer.Write(value)
	if err != nil {
		// Failed to write response
		p.Error(http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// Error sends an error response with the specified status code and message.
func (p *Context) Error(status int, message string) {
	http.Error(p.Writer, message, status)
}
