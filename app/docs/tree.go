// Package docs is the docs site: the page tree for the sidebar, the page of one entry and search.
package docs

import (
	"slices"
	"strings"

	tifaweb "github.com/liustack/tifa-web"
)

// Node is a node of Fumadocs' page tree: a page or a folder of pages.
type Node struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	URL         string `json:"url,omitzero"`
	DefaultOpen bool   `json:"defaultOpen,omitzero"`
	Children    []Node `json:"children,omitzero"`
}

// Tree is the root of Fumadocs' page tree.
type Tree struct {
	Name     string `json:"name"`
	Children []Node `json:"children"`
}

// URL is the address of a doc.
func URL(d tifaweb.Doc) string { return "/docs/" + d.Path }

// Sorted returns every doc in the order of the sidebar: by section, then by place in it.
func Sorted() []tifaweb.Doc {
	all := tifaweb.Docs.All()
	slices.SortStableFunc(all, func(a, b tifaweb.Doc) int {
		if a.GroupOrder != b.GroupOrder {
			return a.GroupOrder - b.GroupOrder
		}
		return a.Order - b.Order
	})
	return all
}

// BuildTree makes the page tree from the collection: one open folder per section.
func BuildTree() Tree {
	t := Tree{Name: "Tifa"}
	for _, d := range Sorted() {
		if n := len(t.Children); n == 0 || t.Children[n-1].Name != d.Group {
			t.Children = append(t.Children, Node{Type: "folder", Name: d.Group, DefaultOpen: true})
		}
		f := &t.Children[len(t.Children)-1]
		f.Children = append(f.Children, Node{Type: "page", Name: d.Title, URL: URL(d)})
	}
	return t
}

// First is the doc /docs leads to.
func First() tifaweb.Doc { return Sorted()[0] }

// Find returns the doc at a path below /docs.
func Find(path string) (tifaweb.Doc, error) {
	return tifaweb.Docs.Get("", strings.Trim(path, "/"))
}
