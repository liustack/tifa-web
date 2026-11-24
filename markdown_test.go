package tifaweb

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// site stands for the app: Markdown versions of /docs/a and /zh-CN/docs/a, pages at any other
// address, and a 404 for any other .md, as Tifa answers them.
var site = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Accept-Encoding")
	switch r.URL.Path {
	case "/docs/a.md", "/zh-CN/docs/a.md":
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		io.WriteString(w, "# A "+r.URL.RawQuery)
	default:
		if strings.HasSuffix(r.URL.Path, ".md") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, "<h1>"+r.Method+" "+r.URL.Path+"</h1>")
	}
})

func TestAgentsGetTheMarkdownOfAPage(t *testing.T) {
	h := markdownForAgents(site)
	for _, tc := range []struct {
		method, url, accept, body, vary string
	}{
		{"GET", "/docs/a", "text/markdown", "# A ", "Accept, Accept-Encoding"},
		{"GET", "/docs/a?x=1", "text/markdown, text/html;q=0.9, */*;q=0.8", "# A x=1", "Accept, Accept-Encoding"},
		{"GET", "/docs/a/", "TEXT/MARKDOWN", "# A ", "Accept, Accept-Encoding"},
		{"GET", "/zh-CN/docs/a", "text/html, text/markdown;q=0.5", "# A ", "Accept, Accept-Encoding"},
		{"HEAD", "/docs/a", "text/markdown", "# A ", "Accept, Accept-Encoding"},
		// A browser, and an Accept that refuses Markdown, get the page.
		{"GET", "/docs/a", "text/html,application/xhtml+xml,*/*;q=0.8", "<h1>GET /docs/a</h1>", "Accept, Accept-Encoding"},
		{"GET", "/docs/a", "text/markdown;q=0, text/html", "<h1>GET /docs/a</h1>", "Accept, Accept-Encoding"},
		{"GET", "/docs/a", "", "<h1>GET /docs/a</h1>", "Accept, Accept-Encoding"},
		// A page without a Markdown version.
		{"GET", "/", "text/markdown", "<h1>GET /</h1>", "Accept, Accept-Encoding"},
		{"GET", "/why", "text/markdown", "<h1>GET /why</h1>", "Accept, Accept-Encoding"},
		// Files and writes are left alone.
		{"GET", "/llms.txt", "text/markdown", "<h1>GET /llms.txt</h1>", "Accept-Encoding"},
		{"POST", "/docs/a", "text/markdown", "<h1>POST /docs/a</h1>", "Accept-Encoding"},
	} {
		r := httptest.NewRequest(tc.method, tc.url, nil)
		if tc.accept != "" {
			r.Header.Set("Accept", tc.accept)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		name := tc.method + " " + tc.url + " Accept: " + tc.accept
		if w.Code != 200 || w.Body.String() != tc.body {
			t.Errorf("%s: %d %q, want %q", name, w.Code, w.Body.String(), tc.body)
		}
		if got := strings.Join(w.Header().Values("Vary"), ", "); got != tc.vary {
			t.Errorf("%s: Vary %q, want %q", name, got, tc.vary)
		}
		if md := strings.HasPrefix(tc.body, "#"); md != strings.HasPrefix(w.Header().Get("Content-Type"), "text/markdown") {
			t.Errorf("%s: Content-Type %q", name, w.Header().Get("Content-Type"))
		}
	}
}
