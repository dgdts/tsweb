package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
	tsweb "tsweb/src"
)

type student struct {
	Name string
	Age  int
}

func onlyForV2(c *tsweb.Context) {
	t := time.Now()
	c.Next()
	// c.Fail(500, "Internal Server Error")
	log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := tsweb.NewEngine()
	r.Use(tsweb.Logger())

	r.GET("/hello", func(c *tsweb.Context) {
		c.String(http.StatusOK, "Hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *tsweb.Context) {
		c.JSON(http.StatusOK, tsweb.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/hello/:name", func(c *tsweb.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Params["name"], c.Path)
	})

	v1 := r.Group("/v1")

	v1.GET("/hello", func(c *tsweb.Context) {
		c.String(http.StatusOK, "This is v1 hello")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2)
	v2.GET("/hello", func(c *tsweb.Context) {
		c.String(http.StatusOK, "This is v2 hello")
	})
	v2.GET("/hello/:name", func(c *tsweb.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("../templates/*")
	r.Static("/assets", "../static")

	stu1 := &student{
		Name: "N1",
		Age:  10,
	}

	stu2 := &student{
		Name: "N2",
		Age:  11,
	}

	r.GET("/", func(c *tsweb.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	r.GET("/students", func(c *tsweb.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", tsweb.H{
			"title":  "tsweb",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *tsweb.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", tsweb.H{
			"title": "tsweb",
			"now":   time.Now(),
		})
	})

	r.Use(tsweb.Recovery())
	r.GET("/panic", func(c *tsweb.Context) {
		names := []string{"test"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
