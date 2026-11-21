package tifaweb

import "github.com/liustack/tifa"

// sitemap lists the docs for /sitemap.xml, /llms.txt, /llms-full.txt and the .md version of
// each page. The pages without parameters (the home page) Tifa lists by itself.
func sitemap(c *tifa.Ctx) (tifa.Sitemap, error) {
	sm := tifa.Sitemap{
		Name: "Tifa",
	}
	sm.Pages = append(sm.Pages, tifa.SitemapPage{Path: "/", Title: "Tifa"})
	for _, d := range Docs.All() {
		sm.Pages = append(sm.Pages, tifa.SitemapPage{
			Path:        "/docs/" + d.Path,
			Title:       d.Title,
			Description: d.Description,
			Content:     d,
		})
	}
	return sm, nil
}
