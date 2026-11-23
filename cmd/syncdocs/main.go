// Command syncdocs writes the docs and the changelog of the Tifa version this site depends on
// into content/docs/en-US, content/docs/zh-CN and content/site.
//
// The site has no copy of the docs and no list of them. Every run reads docs/ from the Tifa
// module in the Go module cache. Each folder there whose index.md says audience: app is a section
// of the sidebar, and each doc in it is a page, in English from <name>.md and in Chinese from the
// <name>.zh-CN.md that Tifa ships beside it. A page keeps its own front matter and its links, which
// point at pages of the site or at the Tifa repository. The folders of content/docs are generated
// and committed, so a build needs no step of its own, and CI checks they match. Run it with
// go generate, after pnpm install and whenever the Tifa version changes.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const modulePath = "github.com/liustack/tifa"

func main() {
	docs, siteOut := "content/docs", "content/site"
	if len(os.Args) > 1 {
		docs = os.Args[1]
	}
	if len(os.Args) > 2 {
		siteOut = os.Args[2]
	}
	dir, version, err := tifaModule()
	if err != nil {
		log.Fatalf("syncdocs: %v", err)
	}
	n, err := Sync(dir, version, docs)
	if err != nil {
		log.Fatalf("syncdocs: %v", err)
	}
	if err := SyncChangelog(dir, version, siteOut); err != nil {
		log.Fatalf("syncdocs: %v", err)
	}
	fmt.Printf("syncdocs: the changelog of %s written to %s\n", version, siteOut)
	fmt.Printf("syncdocs: %d docs of %s %s written to %s in %s and %s\n", n, modulePath, version, docs, english, chinese)
	// Links in the pages written by hand name the Tifa version of go.mod too.
	pinned, err := PinLinks([]string{"app/_locales", "content/pages"}, version)
	if err != nil {
		log.Fatalf("syncdocs: %v", err)
	}
	for _, p := range pinned {
		fmt.Printf("syncdocs: links to Tifa now name %s: %s\n", gitRef(version), p)
	}
}

// tifaModule finds the directory and the version of the Tifa module the app depends on.
func tifaModule() (dir, version string, err error) {
	b, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}\t{{.Version}}", modulePath).Output()
	if err != nil {
		return "", "", fmt.Errorf("go list -m %s: %w", modulePath, err)
	}
	dir, version, ok := strings.Cut(strings.TrimSpace(string(b)), "\t")
	if !ok || dir == "" || version == "" {
		return "", "", fmt.Errorf("go list -m %s gave %q: is the module required without a replace?", modulePath, b)
	}
	return dir, version, nil
}
