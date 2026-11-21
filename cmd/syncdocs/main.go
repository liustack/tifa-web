// Command syncdocs writes the docs of the Tifa version this site depends on into content/docs.
//
// The site has no copy of the docs. Every run reads docs/ from the Tifa module in the Go module
// cache, keeps the docs written for people who build apps (front matter audience: app), and
// turns each into a page of the site: its own front matter, links that point at pages of the
// site or at the Tifa repository. content/docs is generated and committed, so a build needs no step of its own, and CI checks it matches. Run it with
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
	out := "content/docs"
	if len(os.Args) > 1 {
		out = os.Args[1]
	}
	dir, version, err := tifaModule()
	if err != nil {
		log.Fatalf("syncdocs: %v", err)
	}
	n, err := Sync(dir, version, out)
	if err != nil {
		log.Fatalf("syncdocs: %v", err)
	}
	fmt.Printf("syncdocs: %d docs of %s %s written to %s\n", n, modulePath, version, out)
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
