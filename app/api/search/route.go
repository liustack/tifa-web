// Package search answers the search box of the docs in the shape Fumadocs' client reads.
package search

import (
	"strings"

	"github.com/liustack/tifa"

	tifaweb "github.com/liustack/tifa-web"
)

// Result is one line of the search list.
type Result struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Type        string   `json:"type"`
	Content     string   `json:"content"`
	Breadcrumbs []string `json:"breadcrumbs,omitzero"`
}

// Input is the query string of the search.
type Input struct {
	Query string `json:"query"`
}

// GET finds docs, in the language of the URL (/zh-CN/api/search for Chinese), by title and headings that contain every word of the query.
func GET(c *tifa.Ctx, in Input) ([]Result, error) {
	words := strings.Fields(strings.ToLower(in.Query))
	out := []Result{}
	if len(words) == 0 {
		return out, nil
	}
	match := func(s string) bool {
		s = strings.ToLower(s)
		for _, w := range words {
			if !strings.Contains(s, w) {
				return false
			}
		}
		return true
	}
	for _, d := range tifaweb.SortedDocs(c.Language()) {
		url := tifaweb.DocURL(d.Doc)
		if match(d.Title) {
			out = append(out, Result{ID: d.ID, URL: url, Type: "page", Content: d.Title})
		}
		for _, h := range d.Headings() {
			if h.Depth >= 2 && match(h.Text) {
				out = append(out, Result{ID: d.ID + "#" + h.ID, URL: url + "#" + h.ID, Type: "heading", Content: h.Text, Breadcrumbs: []string{d.Title}})
			}
		}
		if len(out) >= 30 {
			break
		}
	}
	return out, nil
}
