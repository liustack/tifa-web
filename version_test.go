package tifaweb

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The version the site shows is the one in go.mod, the one cmd/syncdocs reads the docs of.
func TestTifaVersionIsTheOneInGoMod(t *testing.T) {
	v, err := TifaVersion()
	if err != nil {
		t.Fatal(err)
	}
	mod, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(mod), tifaModule+" "+v+"\n") {
		t.Errorf("TifaVersion() = %q, which go.mod does not require:\n%s", v, mod)
	}
}

func TestSourceURL(t *testing.T) {
	for _, c := range []struct{ version, want string }{
		{"v0.1.13", "https://github.com/liustack/tifa/blob/v0.1.13/docs/app.md"},
		{"v0.1.13-0.20261120205010-e5449dd316a7", "https://github.com/liustack/tifa/blob/e5449dd316a7/docs/app.md"},
	} {
		if got := SourceURL(c.version, "docs/app.md"); got != c.want {
			t.Errorf("SourceURL(%q) = %q, want %q", c.version, got, c.want)
		}
	}
}

// The pages and the home page, written by hand, copy text and code from the Tifa repository at the version of go.mod
// and link to GitHub at that version. go generate rewrites the version in those links. This fails
// when a link was added by hand with another version and go generate was not run.
func TestHandWrittenPagesNameTheVersionOfGoMod(t *testing.T) {
	v, err := TifaVersion()
	if err != nil {
		t.Fatal(err)
	}
	ref := regexp.MustCompile(`github\.com/liustack/tifa/(?:blob|tree)/([^/]+)/`)
	found := 0
	walk := func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		for _, m := range ref.FindAllStringSubmatch(string(b), -1) {
			found++
			if m[1] != gitRef(v) {
				t.Errorf("%s links to Tifa at %s, but go.mod has %s", p, m[1], v)
			}
		}
		return nil
	}
	for _, dir := range []string{"content/pages", "app/_locales", "content/docs/zh-CN"} {
		if err := filepath.WalkDir(dir, walk); err != nil {
			t.Fatal(err)
		}
	}
	if found == 0 {
		t.Error("no link to the Tifa repository found: the check looks at nothing")
	}
}
