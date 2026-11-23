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
	siteURLs = map[string]string{
		"docs/build-and-deploy/deploy.md":       "/docs/build-and-deploy/deploy",
		"docs/build-and-deploy/deploy.zh-CN.md": "/docs/build-and-deploy/deploy",
		"docs/reference/memo.md":                "/docs/reference/memo",
		"docs/reference/memo.zh-CN.md":          "/docs/reference/memo",
		"docs/guides/seo.md":                    "/docs/guides/seo",
		"docs/guides/seo.zh-CN.md":              "/docs/guides/seo",
	}
	titleSlugs = map[string]string{
		"docs/build-and-deploy/deploy.md":       "deploying",
		"docs/build-and-deploy/deploy.zh-CN.md": "deploying",
	}
	for _, c := range []struct{ target, srcDir, want string }{
		{"deploy.md#deploying", "docs/build-and-deploy", "/docs/build-and-deploy/deploy"},
		{"deploy.md#deploy-with-tifa-deploy", "docs/build-and-deploy", "/docs/build-and-deploy/deploy#deploy-with-tifa-deploy"},
		{"../reference/memo.md", "docs/guides", "/docs/reference/memo"},
		{"../guides/seo.md", "docs/reference", "/docs/guides/seo"},
		{"seo.md#x", "docs/guides", "/docs/guides/seo#x"},
		// A Chinese file links to the Chinese file, which is the page of the same URL.
		{"deploy.zh-CN.md#deploying", "docs/build-and-deploy", "/docs/build-and-deploy/deploy"},
		{"../reference/memo.zh-CN.md#a", "docs/guides", "/docs/reference/memo#a"},
		// An English file that is no page of the site is a link to the repository, also from a Chinese file.
		{"../internals/testing.md", "docs/guides", "https://github.com/liustack/tifa/blob/v0.1.12/docs/internals/testing.md"},
		{"../studies/capacity.md#a", "docs/reference", "https://github.com/liustack/tifa/blob/v0.1.12/docs/studies/capacity.md#a"},
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
	siteURLs = map[string]string{"docs/build-and-deploy/deploy.md": "/docs/build-and-deploy/deploy"}
	titleSlugs = map[string]string{}
	in := "See [a](deploy.md) and `[b](deploy.md)`.\n```md\n[c](deploy.md)\n```\n[d](deploy.md \"t\")\n"
	want := "See [a](/docs/build-and-deploy/deploy) and `[b](deploy.md)`.\n```md\n[c](deploy.md)\n```\n[d](/docs/build-and-deploy/deploy \"t\")\n"
	if got := string(rewriteLinks([]byte(in), "docs/build-and-deploy", t.TempDir(), "v0.1.12")); got != want {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
}

func TestShortenChinese(t *testing.T) {
	if got := shorten("怎样把 Tifa 应用放到服务器上：服务器在哪里指定、怎样回滚"); got != "怎样把 Tifa 应用放到服务器上" {
		t.Errorf("got %q", got)
	}
	if got := shorten(strings.Repeat("字", 200)); len([]rune(got)) > 85 || !strings.HasSuffix(got, "...") {
		t.Errorf("got %q", got)
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
// in both languages, leads to a page that exists, and a fragment to a heading of that page.
func TestLinksBetweenDocsAreNotBroken(t *testing.T) {
	dir, version, err := tifaModule()
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if _, err := Sync(dir, version, out); err != nil {
		t.Fatal(err)
	}
	link := regexp.MustCompile(`\]\((/docs/[^)#\s]*)(#[^)\s]*)?`)
	for _, language := range []string{"en-US", "zh-CN"} {
		root := filepath.Join(out, language)
		pages := map[string]map[string]bool{} // URL -> heading ids
		files := map[string]string{}          // URL -> text
		err = filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(root, p)
			url := "/docs/" + strings.TrimSuffix(filepath.ToSlash(rel), ".md")
			files[url] = string(b)
			pages[url] = headingIDs(string(b))
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		checked := 0
		for url, text := range files {
			for _, m := range link.FindAllStringSubmatch(stripCode(text), -1) {
				checked++
				ids, ok := pages[m[1]]
				if !ok {
					t.Errorf("%s: %s links to %s, which is no page", language, url, m[1])
					continue
				}
				if m[2] != "" && !ids[m[2][1:]] {
					t.Errorf("%s: %s links to %s%s, which has no such heading", language, url, m[1], m[2])
				}
			}
		}
		if checked == 0 {
			t.Fatalf("%s: no links found: the check looks at nothing", language)
		}
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
	mdLink     = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	headLine   = regexp.MustCompile(`^#{1,6} +(.+)$`)
	explicitID = regexp.MustCompile(`\{#([^}\s]+)\}$`)
)

// headingIDs is the ids Tifa gives the headings of a page: the {#id} a heading ends in, else the
// text in lower case without punctuation and with - for spaces, a repeat with -1, -2 and so on.
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
		if e := explicitID.FindStringSubmatch(h); e != nil {
			id = e[1]
		}
		if n := count[id]; n > 0 {
			ids[id+"-"+strconv.Itoa(n)] = true
		} else {
			ids[id] = true
		}
		count[id]++
	}
	return ids
}

func TestSyncChangelog(t *testing.T) {
	dir, out := t.TempDir(), t.TempDir()
	siteURLs = map[string]string{
		"docs/reference/performance.md":  "/docs/reference/performance",
		"docs/build-and-deploy/build.md": "/docs/build-and-deploy/build",
	}
	titleSlugs = map[string]string{}
	src := "# Changelog\n\nWhat changed in Tifa. Older changes are by date. See [Performance](docs/reference/performance.md).\n\n## [Unreleased]\n\n- A change, see [the build](docs/build-and-deploy/build.md#x) and [the shop](examples/shop/README.md).\n"
	if err := os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := SyncChangelog(dir, "v0.1.12", out); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(out, "changelog.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	for _, want := range []string{
		"title: Changelog\n",
		"description: What changed in Tifa.\n",
		"source: CHANGELOG.md\n",
		"](/docs/reference/performance)",
		"](/docs/build-and-deploy/build#x)",
		"](https://github.com/liustack/tifa/blob/v0.1.12/examples/shop/README.md)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the page lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "\n# Changelog") {
		t.Errorf("the page repeats its title:\n%s", got)
	}
}

func TestPinLinks(t *testing.T) {
	root := t.TempDir()
	old := "See [a](https://github.com/liustack/tifa/blob/v0.2.0/docs/a.md) and\n" +
		"[b](https://github.com/liustack/tifa/tree/v0.2.0/examples/shop#x), also v0.2.0 and\n" +
		"https://github.com/liustack/tifa/releases/tag/v0.2.0 stay.\r\n"
	same := "[c](https://github.com/liustack/tifa/blob/v0.3.0/docs/c.md)\n"
	none := "no links\n"
	for name, text := range map[string]string{"old.md": old, "sub/same.md": same, "none.json": none} {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	changed, err := PinLinks([]string{root}, "v0.3.0")
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) != 1 || changed[0] != filepath.Join(root, "old.md") {
		t.Fatalf("changed = %v, want only old.md", changed)
	}
	got, _ := os.ReadFile(filepath.Join(root, "old.md"))
	want := strings.NewReplacer("blob/v0.2.0/", "blob/v0.3.0/", "tree/v0.2.0/", "tree/v0.3.0/").Replace(old)
	if string(got) != want {
		t.Errorf("old.md = %q, want %q", got, want)
	}
	again, err := PinLinks([]string{root}, "v0.3.0")
	if err != nil || len(again) != 0 {
		t.Errorf("second run changed %v, err %v", again, err)
	}
	// A pseudo-version links at its commit.
	if _, err := PinLinks([]string{root}, "v0.3.1-0.20261120205010-e5449dd316a7"); err != nil {
		t.Fatal(err)
	}
	got, _ = os.ReadFile(filepath.Join(root, "sub/same.md"))
	if !strings.Contains(string(got), "blob/e5449dd316a7/docs/c.md") {
		t.Errorf("same.md = %q", got)
	}
}

// writeTree writes the files of a fake Tifa module under dir.
func writeTree(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for name, text := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func doc(summary string, order int, body string) string {
	return "---\nsummary: '" + summary + "'\naudience: app\norder: " + strconv.Itoa(order) + "\n---\n\n" + body
}

func indexPage(title string, order int, audience, body string) string {
	return "---\ntitle: " + title + "\norder: " + strconv.Itoa(order) + "\naudience: " + audience + "\n---\n\n# " + title + "\n\n" + body + "\n"
}

// fakeModule is a Tifa module with two sections whose folder names sort the other way than their
// order, a dev section, and a folder with no index.
func fakeModule(t *testing.T) string {
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{
		"docs/zeta/index.md":        indexPage("Zeta", 1, "app", "Zeta intro."),
		"docs/zeta/index.zh-CN.md":  indexPage("泽塔", 1, "app", "泽塔介绍。"),
		"docs/alpha/index.md":       indexPage("Alpha", 2, "app", "Alpha intro."),
		"docs/alpha/index.zh-CN.md": indexPage("阿尔法", 2, "app", "阿尔法介绍。"),
		"docs/zeta/second.md":       doc("Second doc about things: it is second", 2, "# Second\n\nSee [first](first.md#first) and [alpha](../alpha/x.md#sub) and [dev](../internals/dev.md).\n\n## Sub {#sub}\n"),
		"docs/zeta/second.zh-CN.md": doc("第二篇的内容：它排第二", 2, "# 第二篇 {#second}\n\n见[第一篇](first.zh-CN.md#first)和[阿尔法](../alpha/x.zh-CN.md#sub)和[开发](../internals/dev.md)。\n\n## 小节 {#sub}\n"),
		"docs/zeta/first.md":        doc("First doc", 1, "# First\n\nBody.\n"),
		"docs/zeta/first.zh-CN.md":  doc("第一篇", 1, "# 第一篇 {#first}\n\n正文。\n"),
		"docs/alpha/x.md":           doc("X doc", 1, "# X\n\n## Sub\n"),
		"docs/alpha/x.zh-CN.md":     doc("X 篇", 1, "# X 篇 {#x}\n\n## 小节 {#sub}\n"),
		"docs/internals/index.md":   indexPage("Internals", 3, "dev", "For Tifa developers."),
		"docs/internals/dev.md":     doc("Dev doc", 1, "# Dev\n"),
		"docs/studies/README.md":    "# Studies\n",
		"docs/studies/capacity.md":  "# Capacity\n",
	})
	return dir
}

func TestSyncMakesTheSidebarFromTheFolders(t *testing.T) {
	dir, out := fakeModule(t), t.TempDir()
	n, err := Sync(dir, "v0.1.12", out)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Errorf("Sync wrote %d docs, want 3", n)
	}
	for file, want := range map[string][]string{
		"en-US/zeta/first.md":  {"title: First\n", "description: First doc\n", "group: Zeta\n", "groupOrder: 1\n", "order: 1\n", "source: docs/zeta/first.md\n"},
		"en-US/zeta/second.md": {"description: Second doc about things\n", "group: Zeta\n", "order: 2\n"},
		"en-US/alpha/x.md":     {"group: Alpha\n", "groupOrder: 2\n", "order: 1\n"},
		"zh-CN/zeta/first.md":  {"title: 第一篇\n", "description: 第一篇\n", "group: 泽塔\n", "groupOrder: 1\n", "order: 1\n", "source: docs/zeta/first.zh-CN.md\n"},
		"zh-CN/zeta/second.md": {"description: 第二篇的内容\n", "group: 泽塔\n", "order: 2\n"},
		"zh-CN/alpha/x.md":     {"group: 阿尔法\n", "groupOrder: 2\n"},
	} {
		b, err := os.ReadFile(filepath.Join(out, filepath.FromSlash(file)))
		if err != nil {
			t.Fatal(err)
		}
		for _, w := range want {
			if !strings.Contains(string(b), w) {
				t.Errorf("%s lacks %q:\n%s", file, w, b)
			}
		}
		if strings.Contains(string(b), "sourceHash") || strings.Contains(string(b), "translatedFrom") {
			t.Errorf("%s has a hash key:\n%s", file, b)
		}
	}
	for _, none := range []string{"en-US/zeta/index.md", "zh-CN/zeta/index.md", "en-US/internals", "zh-CN/internals", "en-US/studies", "en-US/zeta/first.zh-CN.md"} {
		if _, err := os.Stat(filepath.Join(out, filepath.FromSlash(none))); err == nil {
			t.Errorf("%s was written", none)
		}
	}
}

func TestSyncChineseTitleLosesItsID(t *testing.T) {
	dir, out := fakeModule(t), t.TempDir()
	if _, err := Sync(dir, "v0.1.12", out); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(out, "zh-CN", "zeta", "first.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if strings.Contains(got, "{#first}") || strings.Contains(got, "\n# ") || !strings.Contains(got, "title: 第一篇\n") {
		t.Errorf("the title keeps its id or is repeated:\n%s", got)
	}
	// Headings below the title keep theirs, which gives them the id of the English heading.
	b, _ = os.ReadFile(filepath.Join(out, "zh-CN", "zeta", "second.md"))
	if !strings.Contains(string(b), "## 小节 {#sub}") {
		t.Errorf("a heading lost its id:\n%s", b)
	}
}

func TestSyncMapsLinksOfBothLanguages(t *testing.T) {
	dir, out := fakeModule(t), t.TempDir()
	if _, err := Sync(dir, "v0.1.12", out); err != nil {
		t.Fatal(err)
	}
	dev := "https://github.com/liustack/tifa/blob/v0.1.12/docs/internals/dev.md"
	for file, want := range map[string][]string{
		// first's title id is dropped: the page shows its title itself, with no heading of that id.
		"en-US/zeta/second.md": {"](/docs/zeta/first)", "](/docs/alpha/x#sub)", "](" + dev + ")"},
		"zh-CN/zeta/second.md": {"](/docs/zeta/first)", "](/docs/alpha/x#sub)", "](" + dev + ")"},
	} {
		b, err := os.ReadFile(filepath.Join(out, filepath.FromSlash(file)))
		if err != nil {
			t.Fatal(err)
		}
		for _, w := range want {
			if !strings.Contains(string(b), w) {
				t.Errorf("%s lacks %q:\n%s", file, w, b)
			}
		}
	}
}

func TestSyncFailsWithoutAChineseFile(t *testing.T) {
	dir := fakeModule(t)
	if err := os.Remove(filepath.Join(dir, "docs", "zeta", "first.zh-CN.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := Sync(dir, "v0.1.12", t.TempDir()); err == nil || !strings.Contains(err.Error(), "first.zh-CN.md") {
		t.Errorf("err = %v, want one that names first.zh-CN.md", err)
	}
}

func TestSyncFailsWithoutOrder(t *testing.T) {
	dir := fakeModule(t)
	writeTree(t, dir, map[string]string{"docs/zeta/first.md": "---\nsummary: First doc\naudience: app\n---\n\n# First\n"})
	if _, err := Sync(dir, "v0.1.12", t.TempDir()); err == nil || !strings.Contains(err.Error(), "order") {
		t.Errorf("err = %v, want one that says order is missing", err)
	}
}

// components/site.ts links the Docs button to the first doc of the sidebar.
func TestFirstDocIsTheTutorial(t *testing.T) {
	dir, version, err := tifaModule()
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if _, err := Sync(dir, version, out); err != nil {
		t.Fatal(err)
	}
	type key struct{ group, order int }
	groupRe, orderRe := regexp.MustCompile(`(?m)^groupOrder: (\d+)$`), regexp.MustCompile(`(?m)^order: (\d+)$`)
	read := func(file string) key {
		t.Helper()
		b, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		g, o := groupRe.FindSubmatch(b), orderRe.FindSubmatch(b)
		if g == nil || o == nil {
			t.Fatalf("%s has no groupOrder or order", file)
		}
		gi, _ := strconv.Atoi(string(g[1]))
		oi, _ := strconv.Atoi(string(o[1]))
		return key{gi, oi}
	}
	tutorial := read(filepath.Join(out, "en-US", "start", "tutorial.md"))
	err = filepath.WalkDir(filepath.Join(out, "en-US"), func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if k := read(p); k.group < tutorial.group || (k.group == tutorial.group && k.order < tutorial.order) {
			t.Errorf("%s comes before the tutorial: update links.docs in components/site.ts", p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
