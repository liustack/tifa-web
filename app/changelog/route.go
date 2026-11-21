// Package changelog is the page of the changelog of Tifa.
package changelog

import (
	"github.com/liustack/tifa"

	tifaweb "github.com/liustack/tifa-web"
)

func Load(c *tifa.Ctx) (tifaweb.Prose, error) { return tifaweb.LoadChangelog(c) }
