package tifaweb

import (
	"errors"
	"slices"
)

// Chinese is the language that may have bodies of its own, in content/docs-zh.
const Chinese = "zh-CN"

// DocView is a doc as a language shows it. Own is true when the language has a body of its own,
// so that no note about the language is needed.
type DocView struct {
	Doc
	Own bool
}

// DocURL is the address of a doc, without the language prefix, which Tifa adds.
func DocURL(d Doc) string { return "/docs/" + d.Path }

// viewOf is the doc of a language: the Chinese file when there is one, else the English doc. Its
// place in the sidebar is always the English doc's.
func viewOf(language string, d Doc) DocView {
	if language == Chinese {
		if z, err := DocsZH.Get("", d.Path); err == nil {
			z.Group, z.GroupOrder, z.Order, z.Source = d.Group, d.GroupOrder, d.Order, d.Source
			return DocView{Doc: z, Own: true}
		}
	}
	return DocView{Doc: d, Own: language != Chinese}
}

// SortedDocs returns every doc in the order of the sidebar: by section, then by place in it.
func SortedDocs(language string) []DocView {
	all := Docs.All()
	slices.SortStableFunc(all, func(a, b Doc) int {
		if a.GroupOrder != b.GroupOrder {
			return a.GroupOrder - b.GroupOrder
		}
		return a.Order - b.Order
	})
	out := make([]DocView, len(all))
	for i, d := range all {
		out[i] = viewOf(language, d)
	}
	return out
}

// FindDoc returns the doc at a path below /docs, as a language shows it. It returns
// tifa.ErrNotFound when there is none.
func FindDoc(language, path string) (DocView, error) {
	d, err := Docs.Get("", path)
	if err != nil {
		return DocView{}, err
	}
	return viewOf(language, d), nil
}

// CheckChinese reports a Chinese file that belongs to no doc: a file in content/docs-zh with no
// English doc at its path is a mistake, since no page would show it.
func CheckChinese() error {
	var errs []error
	for _, z := range DocsZH.All() {
		if _, err := Docs.Get("", z.Path); err != nil {
			errs = append(errs, errors.New("content/docs-zh/"+z.ID+".md has no doc at /docs/"+z.Path))
		}
	}
	return errors.Join(errs...)
}
