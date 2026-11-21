package tifaweb

import (
	"os"
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
