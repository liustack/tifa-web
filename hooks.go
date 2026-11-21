package tifaweb

import (
	"context"

	"github.com/liustack/tifa"
)

// OnStart runs once before the app listens. A returned error stops the app.
func OnStart(context.Context) error {
	tifa.SetSitemap(sitemap)
	return nil
}
