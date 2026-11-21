// Package app is the root of the site: the home page and the Layout every page has.
package app

import (
	"github.com/liustack/tifa"

	tifaweb "github.com/liustack/tifa-web"
)

// Home is what the home page shows besides its words: the version of Tifa.
type Home struct {
	Version string `json:"version"`
}

func Load(c *tifa.Ctx) (Home, error) {
	v, err := tifaweb.TifaVersion()
	if err != nil {
		return Home{}, err
	}
	return Home{Version: v}, nil
}
