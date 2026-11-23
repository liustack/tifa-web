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
func doc(language, id, title, group string, groupOrder, order int) Doc {
	d := Doc{Title: title, Group: group, GroupOrder: groupOrder, Order: order}
	d.ID, d.Path, d.Slug, d.Language = id, id, id[strings.Index(id, "/")+1:], language
	return d
}

// The sidebar of a language is ordered by the section order and the order in the section of its
// own docs, and shows the section names of its own files.
func TestSortedDocsFollowTheLanguage(t *testing.T) {
	tifatest.SetContent(t, Docs,
		doc(English, "guides/b", "B", "Guides", 2, 1),
		doc(English, "start/a2", "A2", "Start", 1, 2),
		doc(English, "start/a1", "A1", "Start", 1, 1),
		doc(Chinese, "start/a1", "甲一", "开始", 1, 1),
		doc(Chinese, "start/a2", "甲二", "开始", 1, 2),
		doc(Chinese, "guides/b", "乙", "做法", 2, 1),
	)
	for language, want := range map[string]string{English: "A1 A2 B", Chinese: "甲一 甲二 乙"} {
		var got []string
		for _, d := range SortedDocs(language) {
			got = append(got, d.Title)
		}
		if strings.Join(got, " ") != want {
			t.Errorf("%s: %v, want %s", language, got, want)
		}
	}
	if d, err := FindDoc(Chinese, "start/a1"); err != nil || d.Group != "开始" {
		t.Errorf("FindDoc = %+v, %v", d, err)
	}
}

// The files of the site: every doc has a Chinese file at the same path, with a section of its own
// language. syncdocs writes both, and a doc missing in one language would be a 404 there.
func TestEveryDocIsInBothLanguages(t *testing.T) {
	tifatest.LoadContent(t)
	en := Docs.In(English)
	if len(en) == 0 {
		t.Fatal("no English docs")
	}
	for _, d := range en {
		z, err := Docs.Get(Chinese, d.ID)
		if err != nil {
			t.Errorf("%s has no Chinese doc: run go generate", d.ID)
			continue
		}
		if z.GroupOrder != d.GroupOrder || z.Order != d.Order {
			t.Errorf("%s: the Chinese doc is at section %d place %d, the English one at %d and %d", d.ID, z.GroupOrder, z.Order, d.GroupOrder, d.Order)
		}
	}
	if len(Docs.In(Chinese)) != len(en) {
		t.Errorf("%d English docs and %d Chinese docs", len(en), len(Docs.In(Chinese)))
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
