package docs

import (
	"github.com/liustack/tifa"

	tifaweb "github.com/liustack/tifa-web"
)

// Layout gives every docs page the page tree for the sidebar, in the language of the page, and
// the version of Tifa the docs are for.
func Layout(c *tifa.Ctx) (Frame, error) {
	v, err := tifaweb.TifaVersion()
	if err != nil {
		return Frame{}, err
	}
	return Frame{Tree: BuildTree(c.Language()), Version: v}, nil
}

// Frame is what the docs layout shows around a page.
type Frame struct {
	Tree    Tree   `json:"tree"`
	Version string `json:"version"`
}
