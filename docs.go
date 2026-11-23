package tifaweb

import "slices"

// Chinese is the language of the Chinese docs, in content/docs/zh-CN.
const Chinese = "zh-CN"

// English is the language the docs are written in, the default of the site.
const English = "en-US"

// DocURL is the address of a doc, without the language prefix, which Tifa adds.
func DocURL(d Doc) string { return "/docs/" + d.Path }

// SortedDocs returns every doc of a language in the order of the sidebar: by section, then by
// place in it. Both come from the front matter of the doc, which syncdocs takes from the folders
// of the Tifa docs, so each language has the section names of its own.
func SortedDocs(language string) []Doc {
	all := Docs.In(language)
	slices.SortStableFunc(all, func(a, b Doc) int {
		if a.GroupOrder != b.GroupOrder {
			return a.GroupOrder - b.GroupOrder
		}
		return a.Order - b.Order
	})
	return all
}

// FindDoc returns the doc at a path below /docs, in a language. It returns tifa.ErrNotFound when
// there is none.
func FindDoc(language, path string) (Doc, error) {
	return Docs.Get(language, path)
}
