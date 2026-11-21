package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestRewriteTarget(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "examples", "shop"), 0o755); err != nil {
		t.Fatal(err)
	}
	titleSlugs = map[string]string{"docs/deploy.md": "deploying"}
	for _, c := range []struct{ target, srcDir, want string }{
		{"deploy.md#deploying", "docs", "/docs/ship/deploy"},
		{"deploy.md#deploy-with-tifa-deploy", "docs", "/docs/ship/deploy#deploy-with-tifa-deploy"},
		{"../memo.md", "docs/guides", "/docs/reference/memo"},
		{"guides/seo.md", "docs", "/docs/guides/seo"},
		{"seo.md#x", "docs/guides", "/docs/guides/seo#x"},
		{"testing.md", "docs", "https://github.com/liustack/tifa/blob/v0.1.12/docs/testing.md"},
		{"studies/capacity.md#a", "docs", "https://github.com/liustack/tifa/blob/v0.1.12/docs/studies/capacity.md#a"},
		{"../../examples/shop", "docs/guides", "https://github.com/liustack/tifa/tree/v0.1.12/examples/shop"},
		{"../test/apps/i18n", "docs", "https://github.com/liustack/tifa/blob/v0.1.12/test/apps/i18n"},
		{"#top", "docs", "#top"},
		{"https://example.com/a.md", "docs", "https://example.com/a.md"},
		{"mailto:a@b.c", "docs", "mailto:a@b.c"},
		{"/store/cart", "docs", "/store/cart"},
	} {
		if got := rewriteTarget(c.target, c.srcDir, dir, "v0.1.12"); got != c.want {
			t.Errorf("rewriteTarget(%q, %q) = %q, want %q", c.target, c.srcDir, got, c.want)
		}
	}
}

func TestGitRef(t *testing.T) {
	for in, want := range map[string]string{
		"v0.1.12":                               "v0.1.12",
		"v0.1.13-0.20261120205010-e5449dd316a7": "e5449dd316a7",
		"v0.1.13-rc.1":                          "v0.1.13-rc.1",
	} {
		if got := gitRef(in); got != want {
			t.Errorf("gitRef(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRewriteLinksSkipsCode(t *testing.T) {
	in := "See [a](deploy.md) and `[b](deploy.md)`.\n```md\n[c](deploy.md)\n```\n[d](deploy.md \"t\")\n"
	want := "See [a](/docs/ship/deploy) and `[b](deploy.md)`.\n```md\n[c](deploy.md)\n```\n[d](/docs/ship/deploy \"t\")\n"
	if got := string(rewriteLinks([]byte(in), "docs", t.TempDir(), "v0.1.12")); got != want {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
}

func TestShorten(t *testing.T) {
	if got := shorten("How to add CSS, Tailwind CSS and shadcn/ui to an app"); got != "How to add CSS, Tailwind CSS and shadcn/ui to an app" {
		t.Errorf("got %q", got)
	}
	if got := shorten("How to start a Tifa app: folder layout, how folders become URLs"); got != "How to start a Tifa app" {
		t.Errorf("got %q", got)
	}
	if got := shorten(strings.Repeat("word ", 60)); len(got) > 165 || !strings.HasSuffix(got, "...") {
		t.Errorf("got %q", got)
	}
}

// Every link to a page of the site, in the pages of the Tifa version the app depends on,
// leads to a page that exists, and a fragment to a heading of that page.
func TestLinksBetweenDocsAreNotBroken(t *testing.T) {
	dir, version, err := tifaModule()
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if _, err := Sync(dir, version, out); err != nil {
		t.Fatal(err)
	}
	pages := map[string]map[string]bool{} // URL -> heading ids
	files := map[string]string{}          // URL -> text
	err = filepath.WalkDir(out, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(out, p)
		url := "/docs/" + strings.TrimSuffix(filepath.ToSlash(rel), ".md")
		files[url] = string(b)
		pages[url] = headingIDs(string(b))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	link := regexp.MustCompile(`\]\((/docs/[^)#\s]*)(#[^)\s]*)?`)
	checked := 0
	for url, text := range files {
		for _, m := range link.FindAllStringSubmatch(stripCode(text), -1) {
			checked++
			ids, ok := pages[m[1]]
			if !ok {
				t.Errorf("%s links to %s, which is no page", url, m[1])
				continue
			}
			if m[2] != "" && !ids[m[2][1:]] {
				t.Errorf("%s links to %s%s, which has no such heading", url, m[1], m[2])
			}
		}
	}
	if checked == 0 {
		t.Fatal("no links found: the check looks at nothing")
	}
}

// stripCode removes fenced code blocks, so an example link in one is not checked.
func stripCode(text string) string {
	var out strings.Builder
	fenced := false
	for _, l := range strings.SplitAfter(text, "\n") {
		t := strings.TrimLeft(l, " ")
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			fenced = !fenced
			continue
		}
		if !fenced {
			out.WriteString(l)
		}
	}
	return out.String()
}

var (
	mdLink   = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	headLine = regexp.MustCompile(`^#{1,6} +(.+)$`)
)

// headingIDs is the ids Tifa gives the headings of a page: the text in lower case without
// punctuation and with - for spaces, a repeat with -1, -2 and so on.
func headingIDs(text string) map[string]bool {
	ids := map[string]bool{}
	count := map[string]int{}
	for _, l := range strings.Split(stripCode(text), "\n") {
		m := headLine.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		h := mdLink.ReplaceAllString(strings.ReplaceAll(strings.TrimSpace(m[1]), "`", ""), "$1")
		id := slugify(h)
		if n := count[id]; n > 0 {
			ids[id+"-"+strconv.Itoa(n)] = true
		} else {
			ids[id] = true
		}
		count[id]++
	}
	return ids
}

// components/site.ts links the Docs button to the first doc of the sidebar.
func TestFirstDocIsTheTutorial(t *testing.T) {
	first := Places["docs/tutorial.md"]
	for src, p := range Places {
		if p.Group.Order < first.Group.Order || (p.Group == first.Group && p.Order < first.Order) {
			t.Errorf("%s comes before the tutorial: update links.docs in components/site.ts", src)
		}
	}
	if got := siteURL("docs/tutorial.md"); got != "/docs/start/tutorial" {
		t.Errorf("the tutorial is at %s: update links.docs in components/site.ts", got)
	}
}
