package tifaweb

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/liustack/tifa/tifatest"
)

// doc makes an entry of the docs in a language.
func doc(language, id, title string) Doc {
	d := Doc{Title: title}
	d.ID, d.Path, d.Slug, d.Language = id, id, id[len("start/"):], language
	return d
}

func setDocs(t *testing.T, en []Doc, zh []Doc) {
	t.Helper()
	var all []Doc
	all = append(all, en...)
	all = append(all, zh...)
	tifatest.SetContent(t, Docs, all...)
}

// A translation shows its own body in the place of the English doc in the sidebar, and a fresh
// one has no note.
func TestFreshTranslationReplacesTheEnglishDoc(t *testing.T) {
	a := doc(English, "start/a", "A")
	a.Group, a.Order, a.SourceHash = "Start", 3, "aaaaaaaaaaaa"
	z := doc(Chinese, "start/a", "甲")
	z.TranslatedFrom = "aaaaaaaaaaaa"
	setDocs(t, []Doc{a}, []Doc{z})
	got, err := FindDoc(Chinese, "start/a")
	if err != nil {
		t.Fatal(err)
	}
	if !got.Own || got.Outdated || got.Title != "甲" || got.Group != "Start" || got.Order != 3 {
		t.Errorf("Chinese doc = %+v", got)
	}
	if e, _ := FindDoc(English, "start/a"); e.Title != "A" || !e.Own || e.Outdated {
		t.Errorf("English doc = %+v", e)
	}
	if err := CheckTranslations(); err != nil {
		t.Error(err)
	}
}

// A translation made from an older English doc is marked, so that the page shows the note.
func TestStaleTranslationIsOutdated(t *testing.T) {
	a := doc(English, "start/a", "A")
	a.SourceHash = "bbbbbbbbbbbb"
	z := doc(Chinese, "start/a", "甲")
	z.TranslatedFrom = "aaaaaaaaaaaa"
	setDocs(t, []Doc{a}, []Doc{z})
	got, err := FindDoc(Chinese, "start/a")
	if err != nil {
		t.Fatal(err)
	}
	if !got.Own || !got.Outdated {
		t.Errorf("Chinese doc = %+v, want Own and Outdated", got)
	}
	// The English page never says it is outdated.
	if e, _ := FindDoc(English, "start/a"); e.Outdated {
		t.Error("the English doc is outdated")
	}
	if p := ProseOfDoc(got); !p.Outdated || p.EnglishURL != "/docs/start/a" {
		t.Errorf("Prose = %+v", p)
	}
}

// A doc with no Chinese file shows the English body in Chinese, marked as not Own, which is the
// englishOnly note, and it is not outdated.
func TestMissingTranslationFallsBackToEnglish(t *testing.T) {
	a := doc(English, "start/a", "A")
	a.SourceHash = "aaaaaaaaaaaa"
	setDocs(t, []Doc{a}, nil)
	got, err := FindDoc(Chinese, "start/a")
	if err != nil {
		t.Fatal(err)
	}
	if got.Own || got.Outdated || got.Title != "A" {
		t.Errorf("Chinese doc = %+v, want the English doc, not Own", got)
	}
	if p := ProseOfDoc(got); !p.English || p.Outdated {
		t.Errorf("Prose = %+v", p)
	}
	if len(SortedDocs(Chinese)) != 1 {
		t.Error("the sidebar of the Chinese site lacks the doc")
	}
}

func TestTranslationWithoutADocIsAnError(t *testing.T) {
	z := doc(Chinese, "start/nowhere", "x")
	setDocs(t, nil, []Doc{z})
	if err := CheckTranslations(); err == nil {
		t.Error("CheckTranslations found nothing wrong")
	}
}

// The files of the site: every translation has an English doc, and every English doc has the
// sourceHash that syncdocs writes, which the stale note compares with.
func TestDocsOfTheSiteAreConsistent(t *testing.T) {
	tifatest.LoadContent(t)
	if err := CheckTranslations(); err != nil {
		t.Error(err)
	}
	en := Docs.In(English)
	if len(en) == 0 {
		t.Fatal("no English docs")
	}
	for _, d := range en {
		if d.SourceHash == "" {
			t.Errorf("%s has no sourceHash: run go generate", d.ID)
		}
	}
}

// Tifa gives the root-relative links of a translation the prefix of its language when it converts
// the Markdown, so the Chinese Markdown keeps /docs/... and the page leads to /zh-CN/docs/....
// Tests have no build, so this reads the HTML that tifa build wrote, and is skipped without one.
func TestChineseLinksGetThePrefix(t *testing.T) {
	root := filepath.Join(".tifa", "dist", "content", "docs")
	b, err := os.ReadFile(filepath.Join(root, Chinese, "start", "tutorial.html"))
	if err != nil {
		t.Skip("no build: run go tool tifa build first")
	}
	links := regexp.MustCompile(`href="(/docs/[^"]*)"`).FindAllStringSubmatch(string(b), -1)
	if len(links) != 0 {
		t.Errorf("the Chinese tutorial links to %q without the language prefix", links[0][1])
	}
	if !strings.Contains(string(b), `href="/zh-CN/docs/`) {
		t.Error("the Chinese tutorial has no link with the /zh-CN prefix")
	}
	en, err := os.ReadFile(filepath.Join(root, English, "start", "tutorial.html"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(en), `href="/zh-CN/`) || !strings.Contains(string(en), `href="/docs/`) {
		t.Error("the English tutorial links should keep /docs/ with no prefix")
	}
}
