// Package path is the page of one doc, at /docs/{path...}.
package path

import (
	"github.com/liustack/tifa"
	"github.com/liustack/tifa/content"

	"github.com/liustack/tifa-web/app/docs"
)

// Page is what the docs page shows.
type Page struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Body        tifa.HTML         `json:"body"`
	Headings    []content.Heading `json:"headings"`
}

// Load finds the doc of the URL. /docs itself goes to the first doc.
func Load(c *tifa.Ctx) (Page, error) {
	p := c.Param("path")
	if p == "" {
		return Page{}, tifa.Redirect(docs.URL(docs.First()), 302)
	}
	d, err := docs.Find(p)
	if err != nil {
		return Page{}, err
	}
	return Page{Title: d.Title, Description: d.Description, Body: d.Body(), Headings: d.Headings()}, nil
}
