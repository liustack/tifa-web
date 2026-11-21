// Package examples is the page about the example shop.
package examples

import (
	"github.com/liustack/tifa"

	tifaweb "github.com/liustack/tifa-web"
)

func Load(c *tifa.Ctx) (tifaweb.Prose, error) { return tifaweb.LoadProse(c, "examples") }
