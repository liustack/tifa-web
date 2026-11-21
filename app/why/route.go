// Package why is the page about why Tifa exists, with the ideas behind how an app is organized.
package why

import (
	"github.com/liustack/tifa"

	tifaweb "github.com/liustack/tifa-web"
)

// Page is the page of the author, and the philosophy doc under it.
type Page struct {
	tifaweb.Prose
	Philosophy tifaweb.Prose `json:"philosophy"`
}

func Load(c *tifa.Ctx) (Page, error) {
	p, err := tifaweb.LoadProse(c, "why")
	if err != nil {
		return Page{}, err
	}
	d, err := tifaweb.FindDoc(c.Language(), "start/philosophy")
	if err != nil {
		return Page{}, err
	}
	return Page{Prose: p, Philosophy: tifaweb.ProseOfDoc(d)}, nil
}
