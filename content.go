package tifaweb

import "github.com/liustack/tifa/content"

// The docs are written by cmd/syncdocs from the docs of the Tifa module, into content/docs.
// Run go generate after pnpm install and whenever the Tifa version changes.
//
//go:generate go run ./cmd/syncdocs

// Docs is the docs of Tifa, one entry per page.
var Docs = content.Define[Doc](content.Source{Dir: "docs"})

// Doc is the front matter of a doc. syncdocs writes it, see cmd/syncdocs/sync.go.
type Doc struct {
	content.Entry
	Title       string `json:"title"`
	Description string `json:"description"`
	Group       string `json:"group"`      // the section of the sidebar
	GroupOrder  int    `json:"groupOrder"` // the place of the section, smaller first
	Order       int    `json:"order"`      // the place in the section, smaller first
}
