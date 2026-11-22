package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Group is a section of the sidebar.
type Group struct {
	Slug  string
	Title string
	Order int
}

var (
	start   = Group{"start", "Start", 1}
	pages   = Group{"pages", "Pages and data", 2}
	guides  = Group{"guides", "Guides", 3}
	deploy  = Group{"ship", "Build and deploy", 4}
	details = Group{"reference", "Reference", 5}
)

// Place is where a doc of the Tifa repository goes on the site: its section, and its position
// in the section (smaller first). The URL is /docs/<section>/<file name without .md>.
type Place struct {
	Group Group
	Order int
}

// Places lists every doc of the site by its path in the Tifa repository. A doc with
// audience: app that is not in this table fails the sync, so a new doc is never left out by
// accident. Docs about developing Tifa itself (audience: dev) need no entry.
var Places = map[string]Place{
	"docs/tutorial.md":           {start, 1},
	"docs/app.md":                {start, 2},
	"docs/philosophy.md":         {start, 3},
	"docs/structure.md":          {start, 4},
	"docs/routes.md":             {pages, 1},
	"docs/pages.md":              {pages, 2},
	"docs/styles.md":             {pages, 3},
	"docs/content.md":            {pages, 4},
	"docs/sessions.md":           {pages, 5},
	"docs/hooks.md":              {pages, 6},
	"docs/server-render.md":      {pages, 7},
	"docs/ssr-and-cache.md":      {pages, 8},
	"docs/testing-apps.md":       {pages, 9},
	"docs/guides/forms.md":       {guides, 1},
	"docs/guides/database.md":    {guides, 2},
	"docs/guides/sign-in.md":     {guides, 3},
	"docs/guides/uploads.md":     {guides, 4},
	"docs/guides/email.md":       {guides, 5},
	"docs/guides/jobs.md":        {guides, 6},
	"docs/guides/i18n.md":        {guides, 7},
	"docs/guides/seo.md":         {guides, 8},
	"docs/guides/streaming.md":   {guides, 9},
	"docs/guides/environment.md": {guides, 10},
	"docs/dev.md":                {deploy, 1},
	"docs/build.md":              {deploy, 2},
	"docs/deploy.md":             {deploy, 3},
	"docs/config.md":             {deploy, 4},
	"docs/memo.md":               {details, 1},
	"docs/performance.md":        {details, 2},
	"docs/compatibility.md":      {details, 3},
}

// hashOf is the sourceHash of a file of the Tifa repository: the first 12 hex digits of the
// SHA-256 of its bytes.
func hashOf(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])[:12]
}

// titleSlugs is the anchor of the title of each doc, by its path in the Tifa repository, set by Sync.
var titleSlugs map[string]string

var nonSlug = regexp.MustCompile(`[^\p{L}\p{N}_ -]`)

// slugify is the id Tifa gives a heading: the text in lower case without punctuation, and - for spaces.
func slugify(text string) string {
	return strings.ReplaceAll(nonSlug.ReplaceAllString(strings.ToLower(text), ""), " ", "-")
}

// siteURL is the address of a doc on the site.
func siteURL(src string) string {
	p := Places[src]
	return "/docs/" + p.Group.Slug + "/" + strings.TrimSuffix(path.Base(src), ".md")
}

// Sync reads the docs of the module in dir and writes the pages into out, the folder of the
// English docs, which it empties first. version is the Tifa version that links to other files of the repository name.
func Sync(dir, version, out string) (int, error) {
	found, err := filepath.Glob(filepath.Join(dir, "docs", "*.md"))
	if err != nil {
		return 0, err
	}
	more, err := filepath.Glob(filepath.Join(dir, "docs", "*", "*.md"))
	if err != nil {
		return 0, err
	}
	var srcs []string
	for _, f := range append(found, more...) {
		rel, _ := filepath.Rel(dir, f)
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, "docs/research/") || strings.HasPrefix(rel, "docs/studies/") {
			continue
		}
		srcs = append(srcs, rel)
	}
	if len(srcs) == 0 {
		return 0, fmt.Errorf("no docs under %s/docs: does the module zip lack them?", dir)
	}
	slices.Sort(srcs)
	if err := os.RemoveAll(out); err != nil {
		return 0, err
	}
	titleSlugs = map[string]string{}
	for _, src := range srcs {
		if _, ok := Places[src]; !ok {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(src)))
		if err != nil {
			return 0, err
		}
		if _, body, err := splitFront(b); err == nil {
			if title, _, err := takeTitle(body); err == nil {
				titleSlugs[src] = slugify(title)
			}
		}
	}
	n := 0
	seen := map[string]string{}
	for _, src := range srcs {
		b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(src)))
		if err != nil {
			return 0, err
		}
		front, body, err := splitFront(b)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", src, err)
		}
		if front.Audience != "app" {
			continue
		}
		place, ok := Places[src]
		if !ok {
			return 0, fmt.Errorf("%s is a doc for apps and has no place on the site: add it to Places in cmd/syncdocs/sync.go", src)
		}
		url := siteURL(src)
		if other, dup := seen[url]; dup {
			return 0, fmt.Errorf("%s and %s would both be at %s", src, other, url)
		}
		seen[url] = src
		title, body, err := takeTitle(body)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", src, err)
		}
		page := Page{
			Title:       title,
			SourceHash:  hashOf(b),
			Description: shorten(front.Summary),
			Source:      src,
			Group:       place.Group.Title,
			GroupOrder:  place.Group.Order,
			Order:       place.Order,
		}
		text, err := page.render(rewriteLinks(body, path.Dir(src), dir, version))
		if err != nil {
			return 0, fmt.Errorf("%s: %w", src, err)
		}
		dst := filepath.Join(out, filepath.FromSlash(strings.TrimPrefix(url, "/docs/"))+".md")
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return 0, err
		}
		if err := os.WriteFile(dst, text, 0o644); err != nil {
			return 0, err
		}
		n++
	}
	for src := range Places {
		if _, ok := seen[siteURL(src)]; !ok {
			return 0, fmt.Errorf("Places lists %s, but this Tifa version has no such doc for apps: remove it or fix the table", src)
		}
	}
	return n, nil
}

// SyncChangelog writes the CHANGELOG.md of the module in dir as the page out/changelog.md, with
// its links rewritten as in the docs. Sync must have run first, for the places of the docs.
func SyncChangelog(dir, version, out string) error {
	b, err := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	if err != nil {
		return err
	}
	title, body, err := takeTitle(b)
	if err != nil {
		return fmt.Errorf("CHANGELOG.md: %w", err)
	}
	// The first paragraph says what the file is. Its first sentence is the description.
	intro, _, _ := strings.Cut(strings.TrimSpace(string(body)), "\n\n")
	desc, _, _ := strings.Cut(intro, ". ")
	page := Page{Title: title, Description: strings.TrimSuffix(desc, ".") + ".", Source: "CHANGELOG.md"}
	text, err := page.render(rewriteLinks(body, ".", dir, version))
	if err != nil {
		return err
	}
	if err := os.RemoveAll(out); err != nil {
		return err
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(out, "changelog.md"), text, 0o644)
}

// frontMatter is what the docs of the Tifa repository say about themselves.
type frontMatter struct {
	Summary  string `yaml:"summary"`
	Audience string `yaml:"audience"`
}

func splitFront(b []byte) (frontMatter, []byte, error) {
	var f frontMatter
	rest, ok := bytes.CutPrefix(b, []byte("---\n"))
	if !ok {
		return f, nil, fmt.Errorf("no front matter")
	}
	head, body, ok := bytes.Cut(rest, []byte("\n---\n"))
	if !ok {
		return f, nil, fmt.Errorf("front matter is not closed")
	}
	// Only summary and audience are read. read_when has items that are not valid YAML
	// (docs/dev.md has one that starts with a backtick), so it is cut off first.
	head, _, _ = bytes.Cut(head, []byte("\nread_when:"))
	if err := yaml.Unmarshal(head, &f); err != nil {
		return f, nil, fmt.Errorf("front matter: %w", err)
	}
	return f, body, nil
}

// takeTitle returns the text of the first "# " heading, and the body without that line: the
// page shows the title itself.
func takeTitle(body []byte) (string, []byte, error) {
	lines := strings.SplitAfter(string(body), "\n")
	fenced := false
	for i, l := range lines {
		if strings.HasPrefix(l, "```") || strings.HasPrefix(l, "~~~") {
			fenced = !fenced
		}
		if !fenced && strings.HasPrefix(l, "# ") {
			title := strings.TrimSpace(strings.ReplaceAll(l[2:], "`", ""))
			return title, []byte(strings.Join(append(lines[:i:i], lines[i+1:]...), "")), nil
		}
	}
	return "", nil, fmt.Errorf("no # heading")
}

// shorten makes the summary of a doc fit a meta description: the part before the first colon,
// the whole summary when it has none, cut at a word near 160 characters.
func shorten(summary string) string {
	s := strings.Join(strings.Fields(summary), " ")
	if head, _, ok := strings.Cut(s, ": "); ok && len(head) >= 20 {
		s = head
	}
	s = strings.ReplaceAll(s, "`", "")
	const limit = 160
	if len(s) <= limit {
		return s
	}
	cut := strings.LastIndexByte(s[:limit], ' ')
	return strings.TrimRight(s[:cut], " ,;:(") + "..."
}

// Page is the front matter of a page of the site. It is the schema of the Doc type in content.go.
type Page struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Source      string `yaml:"source,omitempty"` // the file in the Tifa repository the page comes from
	// SourceHash is the first 12 hex digits of the SHA-256 of that file, before it is changed
	// here. A translation records it as translatedFrom, and one that differs is out of date.
	SourceHash string `yaml:"sourceHash,omitempty"`
	Group      string `yaml:"group,omitempty"`
	GroupOrder int    `yaml:"groupOrder,omitempty"`
	Order      int    `yaml:"order,omitempty"`
}

func (p Page) render(body []byte) ([]byte, error) {
	head, err := yaml.Marshal(p)
	if err != nil {
		return nil, err
	}
	return slices.Concat([]byte("---\n"), head, []byte("---\n"), body), nil
}

var linkTarget = regexp.MustCompile(`\]\(([^)\s]*)((?:\s+"[^"]*")?)\)`)

// rewriteLinks points the relative links of a doc at the site, or at the Tifa repository for
// everything that is not a page of the site. Code blocks and inline code stay as written.
func rewriteLinks(body []byte, srcDir, moduleDir, version string) []byte {
	var out strings.Builder
	var fence string
	for _, line := range strings.SplitAfter(string(body), "\n") {
		t := strings.TrimLeft(line, " ")
		if fence != "" {
			if strings.HasPrefix(t, fence) {
				fence = ""
			}
			out.WriteString(line)
			continue
		}
		if m := fenceStart.FindString(t); m != "" {
			fence = m
			out.WriteString(line)
			continue
		}
		out.WriteString(mapOutsideCode(line, func(s string) string {
			return linkTarget.ReplaceAllStringFunc(s, func(m string) string {
				sub := linkTarget.FindStringSubmatch(m)
				return "](" + rewriteTarget(sub[1], srcDir, moduleDir, version) + sub[2] + ")"
			})
		}))
	}
	return []byte(out.String())
}

var fenceStart = regexp.MustCompile("^(`{3,}|~{3,})")

// mapOutsideCode applies f to the parts of a line that are not inline code.
func mapOutsideCode(line string, f func(string) string) string {
	var out strings.Builder
	start := 0
	for i := 0; i < len(line); {
		if line[i] != '`' {
			i++
			continue
		}
		n := 0
		for i+n < len(line) && line[i+n] == '`' {
			n++
		}
		ticks := strings.Repeat("`", n)
		end := -1
		for j := i + n; j < len(line); {
			k := strings.Index(line[j:], ticks)
			if k < 0 {
				break
			}
			k += j
			run := 0
			for k+run < len(line) && line[k+run] == '`' {
				run++
			}
			if run == n {
				end = k + run
				break
			}
			j = k + run
		}
		if end < 0 {
			i += n
			continue
		}
		out.WriteString(f(line[start:i]))
		out.WriteString(line[i:end])
		start, i = end, end
	}
	out.WriteString(f(line[start:]))
	return out.String()
}

// rewriteTarget maps one link target. Absolute links, other schemes and anchors stay.
func rewriteTarget(target, srcDir, moduleDir, version string) string {
	if target == "" || strings.HasPrefix(target, "#") || strings.HasPrefix(target, "/") || strings.Contains(strings.SplitN(target, "/", 2)[0], ":") {
		return target
	}
	file, frag, _ := strings.Cut(target, "#")
	if frag != "" {
		frag = "#" + frag
	}
	resolved := path.Join(srcDir, file)
	if _, ok := Places[resolved]; ok {
		// The page shows its title itself, so its heading has no anchor: link to the page.
		if frag == "#"+titleSlugs[resolved] {
			frag = ""
		}
		return siteURL(resolved) + frag
	}
	kind := "blob"
	if st, err := os.Stat(filepath.Join(moduleDir, filepath.FromSlash(resolved))); err == nil && st.IsDir() {
		kind = "tree"
	}
	return "https://github.com/liustack/tifa/" + kind + "/" + gitRef(version) + "/" + resolved + frag
}

// gitRef is the tag of a release, or the commit of a pseudo-version such as
// v0.1.13-0.20261120205010-e5449dd316a7.
func gitRef(version string) string {
	if i := strings.LastIndexByte(version, '-'); i >= 0 && strings.Count(version, "-") >= 2 {
		return version[i+1:]
	}
	return version
}

// Stale lists the translations in the folder of a language whose translatedFrom differs from the
// sourceHash of the English doc written by Sync into en, as one line each, in path order. A
// translation with no English doc is listed too. A folder that does not exist has none.
func Stale(en, translated string) ([]string, error) {
	var lines []string
	err := filepath.WalkDir(translated, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		rel, _ := filepath.Rel(translated, p)
		rel = filepath.ToSlash(rel)
		from, err := readKey(p, "translatedFrom")
		if err != nil {
			return err
		}
		hash, err := readKey(filepath.Join(en, filepath.FromSlash(rel)), "sourceHash")
		if os.IsNotExist(err) {
			lines = append(lines, rel+": there is no English doc at this path")
			return nil
		}
		if err != nil {
			return err
		}
		if from != hash {
			lines = append(lines, fmt.Sprintf("%s: translated from %s, the English doc is now %s", rel, from, hash))
		}
		return nil
	})
	if os.IsNotExist(err) {
		return nil, nil
	}
	return lines, err
}

// readKey reads one string key of the front matter of a page of the site.
func readKey(file, key string) (string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	rest, ok := bytes.CutPrefix(b, []byte("---\n"))
	if !ok {
		return "", fmt.Errorf("%s: no front matter", file)
	}
	head, _, ok := bytes.Cut(rest, []byte("\n---\n"))
	if !ok {
		return "", fmt.Errorf("%s: front matter is not closed", file)
	}
	var m map[string]any
	if err := yaml.Unmarshal(head, &m); err != nil {
		return "", fmt.Errorf("%s: %w", file, err)
	}
	s, _ := m[key].(string)
	return s, nil
}
