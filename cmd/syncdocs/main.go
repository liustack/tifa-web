// Command syncdocs writes the docs and the changelog of the Tifa version this site depends on
// into content/docs/en-US and content/site.
//
// The site has no copy of the docs. Every run reads docs/ from the Tifa module in the Go module
// cache, keeps the docs written for people who build apps (front matter audience: app), and
// turns each into a page of the site: its own front matter, links that point at pages of the
// site or at the Tifa repository. content/docs/en-US is generated and committed, so a build needs no step of its own, and CI checks it matches. The other folders of content/docs are translations written by hand, which syncdocs never changes: it prints the ones whose English doc changed. Run it with
// go generate, after pnpm install and whenever the Tifa version changes.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	// Only the English folder is written. The translations in the other language folders are
	// written by people, and syncdocs only reads them to find the ones that are out of date.
	out := filepath.Join(docs, "en-US")
	dir, version, err := tifaModule()
	if err != nil {
		log.Fatalf("syncdocs: %v", err)
	}
	n, err := Sync(dir, version, out)
	if err != nil {
		log.Fatalf("syncdocs: %v", err)
	}
	if err := SyncChangelog(dir, version, siteOut); err != nil {
		log.Fatalf("syncdocs: %v", err)
	}
	fmt.Printf("syncdocs: the changelog of %s written to %s\n", version, siteOut)
	fmt.Printf("syncdocs: %d docs of %s %s written to %s\n", n, modulePath, version, out)
	zh := filepath.Join(docs, "zh-CN")
	stale, err := Stale(out, zh)
	if err != nil {
		log.Fatalf("syncdocs: %v", err)
	}
	for _, l := range stale {
		fmt.Printf("syncdocs: out of date, translate again: %s/%s\n", zh, l)
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
