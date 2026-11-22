package tifaweb

import (
	"errors"
	"slices"
)

// Chinese is the language of the translations of the docs, in content/docs/zh-CN.
const Chinese = "zh-CN"

// English is the language the docs are written in, the default of the site. Its folder is
// written by syncdocs.
const English = "en-US"

// DocView is a doc as a language shows it. Own is true when the language has a body of its own,
// so that no note about the language is needed. Outdated is true when that body is a translation
// of an older version of the English doc.
type DocView struct {
	Doc
	Own      bool
	Outdated bool
}

// DocURL is the address of a doc, without the language prefix, which Tifa adds.
func DocURL(d Doc) string { return "/docs/" + d.Path }

// EnglishURL is the address of the English page of a doc. The default language has no prefix, and
// /en-US/... only redirects to it. It is for a plain link, since the page is in another language
// than the one that shows the link.
func EnglishURL(d Doc) string { return DocURL(d) }

// viewOf is the doc of a language: the translation when there is one, else the English doc,
// since Tifa does not fall back to another language itself. Its place in the sidebar is always
// the English doc's.
func viewOf(language string, en Doc) DocView {
	if language == English {
		return DocView{Doc: en, Own: true}
	}
	z, err := Docs.Get(language, en.ID)
	if err != nil {
		return DocView{Doc: en}
	}
	z.Group, z.GroupOrder, z.Order, z.Source, z.SourceHash = en.Group, en.GroupOrder, en.Order, en.Source, en.SourceHash
	return DocView{Doc: z, Own: true, Outdated: z.TranslatedFrom != en.SourceHash}
}

// SortedDocs returns every doc in the order of the sidebar: by section, then by place in it.
func SortedDocs(language string) []DocView {
	all := Docs.In(English)
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
	d, err := Docs.Get(English, path)
	if err != nil {
		return DocView{}, err
	}
	return viewOf(language, d), nil
}

// CheckTranslations reports a translation that belongs to no English doc: a file in
// content/docs/zh-CN with no English doc at its path is a mistake, since no page would show it.
func CheckTranslations() error {
	var errs []error
	for _, z := range Docs.In(Chinese) {
		if _, err := Docs.Get(English, z.ID); err != nil {
			errs = append(errs, errors.New("content/docs/"+Chinese+"/"+z.ID+".md has no doc at /docs/"+z.Path))
		}
	}
	return errors.Join(errs...)
}
