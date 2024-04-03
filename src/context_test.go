package tsweb

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContext_String(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	engine := NewEngine()
	c := makeContext(w, req, engine)

	// Test String method with status code 404 and formatted content
	c.String(http.StatusNotFound, "Page %s not found", "/test")
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestContext_Query(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test?key=value", nil)
	engine := NewEngine()
	c := makeContext(nil, req, engine)

	// Test Query method with an existing query parameter
	value := c.Query("key")
	if value != "value" {
		t.Errorf("Expected value 'value', got '%s'", value)
	}

	// Test Query method with a non-existing query parameter
	value = c.Query("nonexistent")
	if value != "" {
		t.Errorf("Expected empty value, got '%s'", value)
	}
}

func TestContext_PostForm(t *testing.T) {
	req, _ := http.NewRequest("POST", "/test", nil)
	req.PostForm = map[string][]string{"key": {"value"}}
	engine := NewEngine()
	c := makeContext(nil, req, engine)

	// Test PostForm method with an existing form parameter
	value := c.PostForm("key")
	if value != "value" {
		t.Errorf("Expected value 'value', got '%s'", value)
	}

	// Test PostForm method with a non-existing form parameter
	value = c.PostForm("nonexistent")
	if value != "" {
		t.Errorf("Expected empty value, got '%s'", value)
	}
}

func TestContext_Next(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	engine := NewEngine()
	c := makeContext(nil, req, engine)

	// Test Next method with no middleware
	c.handle = func(c *Context) {}
	c.middlewares = &[]HandlerFunc{}
	c.Next()
}

func TestContext_JSON(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	engine := NewEngine()
	c := makeContext(w, req, engine)

	// Test JSON method with valid JSON data
	c.JSON(http.StatusOK, map[string]interface{}{"key": "value"})
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestContext_Error(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	engine := NewEngine()
	c := makeContext(w, req, engine)

	// Test Error method with status code 500 and error message
	c.Error(http.StatusInternalServerError, "Internal Server Error")
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}
