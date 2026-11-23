package tifaweb

import (
	"github.com/liustack/tifa"
	"github.com/liustack/tifa/content"
)

// Prose is a page of Markdown that a route shows: its title, its body, and whether the body is
// English in a site that is Chinese, which the page says.
type Prose struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Body        tifa.HTML         `json:"body"`
	Headings    []content.Heading `json:"headings"`
	English     bool              `json:"english"`
}

func proseOf(d Doc, english bool) Prose {
	return Prose{Title: d.Title, Description: d.Description, Body: d.Body(), Headings: d.Headings(), English: english}
}

// ProseOfDoc is a doc shown as a page of its own.
func ProseOfDoc(d Doc) Prose { return proseOf(d, false) }

// LoadProse reads the page written by hand with an ID, in the language of the request, such as
// why or examples.
func LoadProse(c *tifa.Ctx, id string) (Prose, error) {
	d, err := Pages.Get(c.Language(), id)
	if err != nil {
		return Prose{}, err
	}
	return proseOf(d, false), nil
}

// LoadChangelog reads the changelog of the Tifa version of the site. It is English, and a
// request in another language says so.
func LoadChangelog(c *tifa.Ctx) (Prose, error) {
	d, err := Site.Get("", "changelog")
	if err != nil {
		return Prose{}, err
	}
	return proseOf(d, c.Language() == Chinese), nil
}
