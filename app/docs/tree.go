// Package docs is the docs site: the page tree for the sidebar, the page of one entry and search.
package docs

import (
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

// groupNames are the sections of the sidebar in Chinese. The titles of the docs follow their
// bodies, which are English until a file in content/docs/zh-CN says otherwise.
var groupNames = map[string]string{
	"Start":            "开始",
	"Pages and data":   "页面与数据",
	"Guides":           "做法",
	"Build and deploy": "构建与部署",
	"Reference":        "参考",
}

// GroupName is the name of a section in a language.
func GroupName(language, group string) string {
	if language == tifaweb.Chinese {
		if n, ok := groupNames[group]; ok {
			return n
		}
	}
	return group
}

// BuildTree makes the page tree from the collection: one open folder per section.
func BuildTree(language string) Tree {
	t := Tree{Name: "Tifa"}
	for _, d := range tifaweb.SortedDocs(language) {
		name := GroupName(language, d.Group)
		if n := len(t.Children); n == 0 || t.Children[n-1].Name != name {
			t.Children = append(t.Children, Node{Type: "folder", Name: name, DefaultOpen: true})
		}
		f := &t.Children[len(t.Children)-1]
		f.Children = append(f.Children, Node{Type: "page", Name: d.Title, URL: tifaweb.DocURL(d.Doc)})
	}
	return t
}

// First is the doc /docs leads to.
func First() tifaweb.Doc { return tifaweb.SortedDocs(tifaweb.English)[0].Doc }

// Find returns the doc at a path below /docs, as a language shows it.
func Find(language, path string) (tifaweb.DocView, error) {
	return tifaweb.FindDoc(language, strings.Trim(path, "/"))
}
