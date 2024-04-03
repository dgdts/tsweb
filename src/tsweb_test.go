package tsweb

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGETRequestHandler tests the GET request handler registration.
func TestGETRequestHandler(t *testing.T) {
	engine := NewEngine()
	engine.GET("/test", func(c *Context) {
		c.String(http.StatusOK, "GET request received")
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "GET request received"
	if recorder.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

// TestPOSTRequestHandler tests the POST request handler registration.
func TestPOSTRequestHandler(t *testing.T) {
	engine := NewEngine()
	engine.POST("/test", func(c *Context) {
		c.String(http.StatusOK, "POST request received")
	})

	req, err := http.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "POST request received"
	if recorder.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

// TestStaticFileServing tests serving static files.
func TestStaticFileServing(t *testing.T) {
	engine := NewEngine()
	engine.Static("/static", "../static")

	tests := []struct {
		name     string
		url      string
		expected int
	}{
		{"ValidFile", "/static/cssfile.css", http.StatusOK},
		{"InvalidFile", "/static/nonexistent.txt", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.expected {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expected)
			}
		})
	}
}
