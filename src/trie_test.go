package tsweb

import (
	"testing"
)

func TestNode_Insert(t *testing.T) {
	root := &node{part: "/"}

	// Test inserting a pattern with single segment
	root.insert("/hello", []string{"hello"}, 0)
	if root.children[0].pattern != "/hello" {
		t.Errorf("Insert failed, expected pattern: '/hello', got: '%s'", root.children[0].pattern)
	}

	// Test inserting a pattern with multiple segments
	root.insert("/hello/world", []string{"hello", "world"}, 0)
	if root.children[0].children[0].pattern != "/hello/world" {
		t.Errorf("Insert failed, expected pattern: '/hello/world', got: '%s'", root.children[0].children[0].pattern)
	}

	// Test inserting a pattern with wildcard segment
	root.insert("/:name", []string{":name"}, 0)
	if root.children[1].pattern != "/:name" {
		t.Errorf("Insert failed, expected pattern: '/:name', got: '%s'", root.children[0].children[1].pattern)
	}

	// Test inserting a pattern with wildcard segment followed by static segment
	root.insert("/:name/world", []string{":name", "world"}, 0)
	if root.children[1].children[0].pattern != "/:name/world" {
		t.Errorf("Insert failed, expected pattern: '/:name/world', got: '%s'", root.children[0].children[1].children[0].pattern)
	}
}

func TestNode_Search(t *testing.T) {
	root := &node{part: "/"}
	root.insert("/hello", []string{"hello"}, 0)
	root.insert("/hello/world", []string{"hello", "world"}, 0)
	root.insert("/:name", []string{":name"}, 0)
	root.insert("/:name/world", []string{":name", "world"}, 0)

	// Test searching for exact match
	found := root.search([]string{"hello"}, 0)
	if found == nil || found.pattern != "/hello" {
		t.Errorf("Search failed, expected pattern: '/hello', got: '%s'", found.pattern)
	}

	// Test searching for pattern with wildcard segment
	found = root.search([]string{"john"}, 0)
	if found == nil || found.pattern != "/:name" {
		t.Errorf("Search failed, expected pattern: '/:name', got: '%s'", found.pattern)
	}

	// Test searching for pattern with wildcard segment followed by static segment
	found = root.search([]string{"john", "world"}, 0)
	if found == nil || found.pattern != "/:name/world" {
		t.Errorf("Search failed, expected pattern: '/:name/world', got: '%s'", found.pattern)
	}
}
