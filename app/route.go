// Package app is the root of the site: the home page and the Layout every page has.
package app

import (
	"github.com/liustack/tifa"

	tifaweb "github.com/liustack/tifa-web"
)

// File is a file of the demo on the home page, as it is in examples/shop of Tifa: its path, a link
// to it on GitHub and its code, colored by the build.
type File struct {
	Path string    `json:"path"`
	URL  string    `json:"url"`
	Body tifa.HTML `json:"body"`
}

// Home is what the home page shows besides its words: the version of Tifa and the demo files.
type Home struct {
	Version string `json:"version"`
	Demo    []File `json:"demo"`
}

func Load(c *tifa.Ctx) (Home, error) {
	v, err := tifaweb.TifaVersion()
	if err != nil {
		return Home{}, err
	}
	h := Home{Version: v}
	for _, d := range tifaweb.Demo.All() {
		h.Demo = append(h.Demo, File{Path: d.Title, URL: tifaweb.SourceURL(v, "examples/shop/"+d.Title), Body: d.Body()})
	}
	return h, nil
}
