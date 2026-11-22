// Package path is the page of one doc, at /docs/{path...}.
package path

import (
	"github.com/liustack/tifa"
	"github.com/liustack/tifa/content"

	tifaweb "github.com/liustack/tifa-web"
	"github.com/liustack/tifa-web/app/docs"
)

// Page is what the docs page shows.
type Page struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Body        tifa.HTML         `json:"body"`
	Headings    []content.Heading `json:"headings"`
	// English is true when the page shows an English body in the Chinese site, which says so.
	English bool `json:"english"`
	// Outdated is true when the translation was made from an older English doc, which the page
	// says with a link to EnglishURL.
	Outdated   bool   `json:"outdated"`
	EnglishURL string `json:"englishUrl,omitempty"`
	// Markdown is the address of the Markdown version of the page: the one in the language of the
	// page when the language has a body of its own, else the English one.
	Markdown string `json:"markdown"`
	// Source is the file of the doc in the Tifa repository, at the version of the site.
	Source string `json:"source"`
}

// Load finds the doc of the URL. /docs itself goes to the first doc.
func Load(c *tifa.Ctx) (Page, error) {
	p := c.Param("path")
	if p == "" {
		return Page{}, tifa.Redirect(tifaweb.DocURL(docs.First()), 302)
	}
	d, err := docs.Find(c.Language(), p)
	if err != nil {
		return Page{}, err
	}
	md := tifaweb.DocURL(d.Doc) + ".md"
	if d.Own && c.Language() == tifaweb.Chinese {
		md = "/" + c.Language() + md
	}
	ref, err := tifaweb.TifaVersion()
	if err != nil {
		return Page{}, err
	}
	page := Page{
		Title:       d.Title,
		Description: d.Description,
		Body:        d.Body(),
		Headings:    d.Headings(),
		English:     !d.Own,
		Markdown:    md,
		Source:      tifaweb.SourceURL(ref, d.Source),
		Outdated:    d.Outdated,
	}
	if d.Outdated {
		page.EnglishURL = tifaweb.EnglishURL(d.Doc)
	}
	return page, nil
}
