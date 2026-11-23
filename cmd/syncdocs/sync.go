package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Languages of the site. English is read from <name>.md, Chinese from <name>.zh-CN.md, and each is
// written to the folder of its language.
const (
	english = "en-US"
	chinese = "zh-CN"
)

// siteURLs is the address of every doc of the site by its path in the Tifa repository, for both
// languages: docs/start/app.md and docs/start/app.zh-CN.md are both /docs/start/app. Set by Sync.
var siteURLs map[string]string

// titleSlugs is the anchor of the title of each doc, by its path in the Tifa repository, set by Sync.
var titleSlugs map[string]string

var nonSlug = regexp.MustCompile(`[^\p{L}\p{N}_ -]`)

// slugify is the id Tifa gives a heading: the text in lower case without punctuation, and - for spaces.
func slugify(text string) string {
	return strings.ReplaceAll(nonSlug.ReplaceAllString(strings.ToLower(text), ""), " ", "-")
}

// section is a folder of docs/ for apps: one group of the sidebar.
type section struct {
	slug           string
	order          int
	title, zhTitle string
}

// source is one doc in both languages.
type source struct {
	section        *section
	name           string // file name without .md
	order          int
	enSrc, zhSrc   string // paths in the Tifa repository
	en, zh         frontMatter
	enBody, zhBody []byte
}

func (s source) url() string { return "/docs/" + s.section.slug + "/" + s.name }

// Sync reads the docs of the module in dir and writes the pages into out/en-US and out/zh-CN,
// which it empties first. A folder of docs/ is a section when its index.md says audience: app.
// Its title and order come from index.md (index.zh-CN.md for Chinese), and each other doc of the
// folder with audience: app is a page, in the place its order says. Every page must have its
// Chinese file beside it. version is the Tifa version that links to other files of the repository
// name. It returns the number of docs, one per URL.
func Sync(dir, version, out string) (int, error) {
	indexes, err := filepath.Glob(filepath.Join(dir, "docs", "*", "index.md"))
	if err != nil {
		return 0, err
	}
	slices.Sort(indexes)
	var sections []*section
	var srcs []source
	for _, index := range indexes {
		folder := filepath.Dir(index)
		sec, err := readSection(dir, folder)
		if err != nil {
			return 0, err
		}
		if sec == nil {
			continue
		}
		sections = append(sections, sec)
		files, err := filepath.Glob(filepath.Join(folder, "*.md"))
		if err != nil {
			return 0, err
		}
		slices.Sort(files)
		for _, f := range files {
			base := filepath.Base(f)
			if base == "index.md" || strings.HasSuffix(base, ".zh-CN.md") {
				continue
			}
			s, ok, err := readSource(dir, sec, f)
			if err != nil {
				return 0, err
			}
			if ok {
				srcs = append(srcs, s)
			}
		}
	}
	if len(srcs) == 0 {
		return 0, fmt.Errorf("no docs for apps under %s/docs: does the module zip lack them?", dir)
	}
	slices.SortStableFunc(sections, func(a, b *section) int { return a.order - b.order })
	for i := 1; i < len(sections); i++ {
		if sections[i].order == sections[i-1].order {
			return 0, fmt.Errorf("the sections %s and %s both have order %d", sections[i-1].slug, sections[i].slug, sections[i].order)
		}
	}
	siteURLs, titleSlugs = map[string]string{}, map[string]string{}
	for _, s := range srcs {
		siteURLs[s.enSrc], siteURLs[s.zhSrc] = s.url(), s.url()
		title, id, _, err := takeTitle(s.enBody)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", s.enSrc, err)
		}
		titleSlugs[s.enSrc] = slugify(title)
		_, id, _, err = takeTitle(s.zhBody)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", s.zhSrc, err)
		}
		if id == "" {
			return 0, fmt.Errorf("%s: the title has no {#id}: it needs the id of the English title, %s", s.zhSrc, titleSlugs[s.enSrc])
		}
		titleSlugs[s.zhSrc] = id
	}
	for _, language := range []string{english, chinese} {
		if err := os.RemoveAll(filepath.Join(out, language)); err != nil {
			return 0, err
		}
	}
	for _, s := range srcs {
		for _, language := range []string{english, chinese} {
			src, front, body, group := s.enSrc, s.en, s.enBody, s.section.title
			if language == chinese {
				src, front, body, group = s.zhSrc, s.zh, s.zhBody, s.section.zhTitle
			}
			title, _, body, err := takeTitle(body)
			if err != nil {
				return 0, fmt.Errorf("%s: %w", src, err)
			}
			page := Page{
				Title:       title,
				Description: shorten(front.Summary),
				Source:      src,
				Group:       group,
				GroupOrder:  s.section.order,
				Order:       s.order,
			}
			text, err := page.render(rewriteLinks(body, path.Dir(src), dir, version))
			if err != nil {
				return 0, fmt.Errorf("%s: %w", src, err)
			}
			dst := filepath.Join(out, language, s.section.slug, s.name+".md")
			if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
				return 0, err
			}
			if err := os.WriteFile(dst, text, 0o644); err != nil {
				return 0, err
			}
		}
	}
	return len(srcs), nil
}

// readSection reads the index files of a folder. It returns nil for a folder that is not for apps.
func readSection(moduleDir, folder string) (*section, error) {
	en, _, err := readDoc(moduleDir, filepath.Join(folder, "index.md"))
	if err != nil {
		return nil, err
	}
	if en.Audience != "app" {
		return nil, nil
	}
	zhFile := filepath.Join(folder, "index.zh-CN.md")
	zh, _, err := readDoc(moduleDir, zhFile)
	if err != nil {
		return nil, err
	}
	rel := relPath(moduleDir, filepath.Join(folder, "index.md"))
	if en.Title == "" || zh.Title == "" || en.Order <= 0 {
		return nil, fmt.Errorf("%s and its Chinese file need a title, and %s an order above 0", rel, rel)
	}
	return &section{slug: filepath.Base(folder), order: en.Order, title: en.Title, zhTitle: zh.Title}, nil
}

// readSource reads a doc and its Chinese file. ok is false for a doc that is not for apps.
func readSource(moduleDir string, sec *section, file string) (source, bool, error) {
	en, enBody, err := readDoc(moduleDir, file)
	if err != nil {
		return source{}, false, err
	}
	if en.Audience != "app" {
		return source{}, false, nil
	}
	enSrc := relPath(moduleDir, file)
	if en.Order <= 0 {
		return source{}, false, fmt.Errorf("%s has no order in its front matter: a doc for apps needs one above 0 for its place in the section", enSrc)
	}
	zhFile := strings.TrimSuffix(file, ".md") + ".zh-CN.md"
	zh, zhBody, err := readDoc(moduleDir, zhFile)
	if err != nil {
		return source{}, false, fmt.Errorf("%s has no Chinese file: %w", enSrc, err)
	}
	return source{
		section: sec, name: strings.TrimSuffix(filepath.Base(file), ".md"), order: en.Order,
		enSrc: enSrc, zhSrc: relPath(moduleDir, zhFile), en: en, zh: zh, enBody: enBody, zhBody: zhBody,
	}, true, nil
}

func relPath(moduleDir, file string) string {
	rel, _ := filepath.Rel(moduleDir, file)
	return filepath.ToSlash(rel)
}

// readDoc reads a doc of the Tifa module: its front matter and its body.
func readDoc(moduleDir, file string) (frontMatter, []byte, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return frontMatter{}, nil, err
	}
	front, body, err := splitFront(b)
	if err != nil {
		return frontMatter{}, nil, fmt.Errorf("%s: %w", relPath(moduleDir, file), err)
	}
	return front, body, nil
}

// SyncChangelog writes the CHANGELOG.md of the module in dir as the page out/changelog.md, with
// its links rewritten as in the docs. Sync must have run first, for the places of the docs.
func SyncChangelog(dir, version, out string) error {
	b, err := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	if err != nil {
		return err
	}
	title, _, body, err := takeTitle(b)
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
	Title    string `yaml:"title"` // of a section index
	Summary  string `yaml:"summary"`
	Audience string `yaml:"audience"`
	Order    int    `yaml:"order"`
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
	// Only title, summary, audience and order are read. read_when has items that are not valid YAML
	// (docs/dev.md has one that starts with a backtick), so it is cut off first.
	head, _, _ = bytes.Cut(head, []byte("\nread_when:"))
	if err := yaml.Unmarshal(head, &f); err != nil {
		return f, nil, fmt.Errorf("front matter: %w", err)
	}
	return f, body, nil
}

var titleID = regexp.MustCompile(`\s*\{#([^}\s]+)\}$`)

// takeTitle returns the text of the first "# " heading, the id a Chinese heading ends in as
// {#id} (empty when it has none, and not part of the text), and the body without that line: the
// page shows the title itself.
func takeTitle(body []byte) (title, id string, rest []byte, err error) {
	lines := strings.SplitAfter(string(body), "\n")
	fenced := false
	for i, l := range lines {
		if strings.HasPrefix(l, "```") || strings.HasPrefix(l, "~~~") {
			fenced = !fenced
		}
		if !fenced && strings.HasPrefix(l, "# ") {
			title := strings.TrimSpace(strings.ReplaceAll(l[2:], "`", ""))
			if m := titleID.FindStringSubmatch(title); m != nil {
				id = m[1]
				title = strings.TrimSpace(strings.TrimSuffix(title, m[0]))
			}
			return title, id, []byte(strings.Join(append(lines[:i:i], lines[i+1:]...), "")), nil
		}
	}
	return "", "", nil, fmt.Errorf("no # heading")
}

// shorten makes the summary of a doc fit a meta description: the part before the first colon
// (": " in English, "：" in Chinese), the whole summary when it has none, cut near 160 characters
// of English or 80 of Chinese.
func shorten(summary string) string {
	s := strings.Join(strings.Fields(summary), " ")
	cjk := strings.IndexFunc(s, func(r rune) bool { return r >= 0x2E80 }) >= 0
	minHead, limit := 20, 160
	if cjk {
		minHead, limit = 4, 80
	}
	cut := len(s)
	for _, colon := range []string{": ", "："} {
		if i := strings.Index(s, colon); i >= 0 && i < cut && len([]rune(s[:i])) >= minHead {
			cut = i
		}
	}
	s = strings.ReplaceAll(s[:cut], "`", "")
	if r := []rune(s); len(r) > limit {
		s = string(r[:limit])
		if !cjk {
			if i := strings.LastIndexByte(s, ' '); i >= 0 {
				s = s[:i]
			}
		}
		s = strings.TrimRight(s, " ,;:(，；：（") + "..."
	}
	return s
}

// Page is the front matter of a page of the site. It is the schema of the Doc type in content.go.
type Page struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Source      string `yaml:"source,omitempty"` // the file in the Tifa repository the page comes from
	Group       string `yaml:"group,omitempty"`
	GroupOrder  int    `yaml:"groupOrder,omitempty"`
	Order       int    `yaml:"order,omitempty"`
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
	if url, ok := siteURLs[resolved]; ok {
		// The page shows its title itself, so its heading has no anchor: link to the page.
		if frag == "#"+titleSlugs[resolved] {
			frag = ""
		}
		return url + frag
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

// versionLink matches the part of a link to the Tifa repository that names a version.
var versionLink = regexp.MustCompile(`github\.com/liustack/tifa/(blob|tree)/[^/\s)"]+/`)

// PinLinks rewrites, in every file under the given folders, the version in links to the Tifa
// repository (github.com/liustack/tifa/blob/<ref>/ and .../tree/<ref>/) to the ref of version.
// Nothing else in a file changes. It returns the files it changed, in path order.
func PinLinks(roots []string, version string) ([]string, error) {
	ref := gitRef(version)
	var changed []string
	for _, root := range roots {
		err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			out := versionLink.ReplaceAllFunc(b, func(m []byte) []byte {
				kind := versionLink.FindSubmatch(m)[1]
				return []byte("github.com/liustack/tifa/" + string(kind) + "/" + ref + "/")
			})
			if bytes.Equal(out, b) {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return err
			}
			if err := os.WriteFile(p, out, info.Mode().Perm()); err != nil {
				return err
			}
			changed = append(changed, p)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return changed, nil
}
