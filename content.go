package tifaweb

import (
	"github.com/liustack/tifa"
	"github.com/liustack/tifa/content"
)

// The docs are written by cmd/syncdocs from the docs of the Tifa module, into content/docs/en-US.
// Run go generate after pnpm install and whenever the Tifa version changes.
//
//go:generate go run ./cmd/syncdocs

// Docs is the docs of Tifa, one entry per page, in a folder per language. The English folder,
// content/docs/en-US, is written by cmd/syncdocs. The other languages are translations written
// by hand, and a doc with no file in a language shows the English body, see docs.go.
var Docs = content.Define[Doc](content.Source{Dir: "docs", Languages: tifa.Languages()})

// Site is the pages that syncdocs writes from the Tifa repository besides the docs: the
// changelog.
var Site = content.Define[Doc](content.Source{Dir: "site"})

// Doc is the front matter of a doc. syncdocs writes it, see cmd/syncdocs/sync.go.
type Doc struct {
	content.Entry
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Source      string `json:"source,omitempty"`     // the file in the Tifa repository the page comes from
	Group       string `json:"group,omitempty"`      // the section of the sidebar
	GroupOrder  int    `json:"groupOrder,omitempty"` // the place of the section, smaller first
	Order       int    `json:"order,omitempty"`      // the place in the section, smaller first
	// SourceHash is, in an English doc, the first 12 hex digits of the SHA-256 of the file in the
	// Tifa repository it comes from, as syncdocs found it.
	SourceHash string `json:"sourceHash,omitempty"`
	// TranslatedFrom is, in a translation, the SourceHash of the English doc it was translated
	// from. A translation whose TranslatedFrom differs from the SourceHash of its English doc is
	// out of date.
	TranslatedFrom string `json:"translatedFrom,omitempty"`
}

// Pages is the pages of the site that are written by hand, in each language of the site: Why and
// Examples.
var Pages = content.Define[Doc](content.Source{Dir: "pages", Languages: tifa.Languages()})
