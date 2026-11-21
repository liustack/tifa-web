package tifaweb

import (
	"fmt"

	"github.com/liustack/tifa"
)

// notes is the start of /llms.txt: how to read this site, as facts. It names the version of the
// docs in sitemap, since the version comes from the build.
const notes = "The docs on this site are the docs of Tifa %s. Every page of the docs has a Markdown version at its address with .md added, such as /docs/start/tutorial.md, and /llms-full.txt has all of them in one file. To read the docs of the Tifa version an app uses, run `go tool tifa docs` in the app."

// sitemap lists the docs for /sitemap.xml, /llms.txt, /llms-full.txt and the .md version of
// each page, in both languages. A Chinese page has a Markdown version only when it has a body of
// its own, since the others show the English one.
func sitemap(c *tifa.Ctx) (tifa.Sitemap, error) {
	version, err := TifaVersion()
	if err != nil {
		return tifa.Sitemap{}, err
	}
	sm := tifa.Sitemap{
		Name:    "Tifa",
		Summary: fmt.Sprintf(notes, version),
	}
	for _, language := range tifa.Languages() {
		sm.Pages = append(sm.Pages, tifa.SitemapPage{Language: language, Path: "/", Title: "Tifa"})
		for _, d := range SortedDocs(language) {
			page := tifa.SitemapPage{
				Language:    language,
				Path:        DocURL(d.Doc),
				Title:       d.Title,
				Description: d.Description,
			}
			if d.Own {
				page.Content = d.Doc
			}
			sm.Pages = append(sm.Pages, page)
		}
	}
	return sm, nil
}
