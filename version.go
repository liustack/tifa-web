package tifaweb

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// tifaModule is the Tifa module the docs of the site come from.
const tifaModule = "github.com/liustack/tifa"

// TifaVersion is the version of Tifa this site is built with, such as v0.1.13. The docs of the
// site are the docs of that version (cmd/syncdocs reads the same one from go.mod).
func TifaVersion() (string, error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", fmt.Errorf("no build info: was the program built with module support?")
	}
	for _, m := range info.Deps {
		if m.Path == tifaModule {
			if m.Replace != nil {
				return m.Replace.Version, nil
			}
			return m.Version, nil
		}
	}
	return "", fmt.Errorf("%s is not among the dependencies of the program", tifaModule)
}

// SourceURL is the address on GitHub of a file of the Tifa repository at a version. A
// pseudo-version such as v0.1.13-0.20261120205010-e5449dd316a7 names the commit at its end.
func SourceURL(version, file string) string {
	ref := version
	if i := strings.LastIndexByte(version, '-'); i >= 0 && strings.Count(version, "-") >= 2 {
		ref = version[i+1:]
	}
	return "https://github.com/liustack/tifa/blob/" + ref + "/" + file
}
