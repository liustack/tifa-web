package tifaweb

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Shiki loads only the languages vite.config.ts lists, and colors a block of any other language as
// plain text without a word. So every language a code block of the content names must be listed.
func TestCodeLanguagesAreLoaded(t *testing.T) {
	config, err := os.ReadFile("vite.config.ts")
	if err != nil {
		t.Fatal(err)
	}
	list := regexp.MustCompile(`langs:\s*\[([^\]]*)\]`).FindSubmatch(config)
	if list == nil {
		t.Fatal("vite.config.ts has no langs list for @shikijs/rehype")
	}
	loaded := map[string]bool{"text": true} // the fallback language, never loaded
	for _, l := range regexp.MustCompile(`'([^']+)'`).FindAllSubmatch(list[1], -1) {
		loaded[string(l[1])] = true
	}
	fence := regexp.MustCompile("(?m)^\\s*```([A-Za-z0-9_+-]+)")
	err = filepath.WalkDir("content", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		text, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, m := range fence.FindAllSubmatch(text, -1) {
			if lang := string(m[1]); !loaded[lang] {
				t.Errorf("%s has a %s code block, which langs in vite.config.ts does not list", path, lang)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
