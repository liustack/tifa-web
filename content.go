package tifaweb

import "github.com/liustack/tifa/content"

// The docs are written by cmd/syncdocs from the docs of the Tifa module, into content/docs.
// Run go generate after pnpm install and whenever the Tifa version changes.
//
//go:generate go run ./cmd/syncdocs

// Docs is the docs of Tifa, one entry per page.
var Docs = content.Define[Doc](content.Source{Dir: "docs"})

// Site is the pages that syncdocs writes from the Tifa repository besides the docs: the
// changelog.
var Site = content.Define[Doc](content.Source{Dir: "site"})

// Doc is the front matter of a doc. syncdocs writes it, see cmd/syncdocs/sync.go.
type Doc struct {
	content.Entry
	Title       string `json:"title"`
	Description string `json:"description"`
	Source      string `json:"source,omitempty"`     // the file in the Tifa repository the page comes from
	Group       string `json:"group,omitempty"`      // the section of the sidebar
	GroupOrder  int    `json:"groupOrder,omitempty"` // the place of the section, smaller first
	Order       int    `json:"order,omitempty"`      // the place in the section, smaller first
}

// DocsZH holds Chinese bodies of docs, in the same paths as content/docs. A doc with a file
// here shows it in Chinese. The others show the English body with a note, see app/docs.
var DocsZH = content.Define[Doc](content.Source{Dir: "docs-zh"})
