package tifaweb

import (
	"bytes"
	"mime"
	"net/http"
	"path"
	"strings"
)

// markdownForAgents answers a GET or HEAD of a page whose Accept names text/markdown, as coding
// agents such as Claude Code send it, with the Markdown version of the page, which Tifa serves at
// the address of the page with .md. A page without one, such as the home page, answers as usual.
// Every page answer says Vary: Accept, since the same address has two versions.
func markdownForAgents(app http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead || path.Ext(r.URL.Path) != "" {
			app.ServeHTTP(w, r)
			return
		}
		w.Header().Add("Vary", "Accept")
		if !acceptsMarkdown(r.Header.Values("Accept")) {
			app.ServeHTTP(w, r)
			return
		}
		md := r.Clone(r.Context())
		md.URL.Path = markdownPath(r.URL.Path)
		md.URL.RawPath = ""
		md.RequestURI = md.URL.RequestURI()
		rec := &recorder{header: http.Header{}}
		app.ServeHTTP(rec, md)
		if rec.code != http.StatusOK {
			app.ServeHTTP(w, r)
			return
		}
		h := w.Header()
		for k, v := range rec.header {
			if k == "Vary" {
				for _, value := range v {
					h.Add(k, value)
				}
				continue
			}
			h[k] = v
		}
		w.WriteHeader(http.StatusOK)
		w.Write(rec.body.Bytes())
	})
}

// markdownPath is the address of the Markdown version of the page at p.
func markdownPath(p string) string {
	p = strings.TrimSuffix(p, "/")
	if p == "" {
		return "/index.md"
	}
	return p + ".md"
}

// acceptsMarkdown reports whether the Accept headers name text/markdown with a weight above 0.
// Browsers never name it, so the weight of HTML does not matter.
func acceptsMarkdown(accept []string) bool {
	for _, line := range accept {
		for _, part := range strings.Split(line, ",") {
			media, params, err := mime.ParseMediaType(strings.TrimSpace(part))
			if err != nil || media != "text/markdown" {
				continue
			}
			if q, ok := params["q"]; ok && strings.Trim(q, "0.") == "" {
				continue
			}
			return true
		}
	}
	return false
}

// recorder keeps an answer so that markdownForAgents can drop it when it is not the Markdown
// version.
type recorder struct {
	header http.Header
	code   int
	body   bytes.Buffer
}

func (r *recorder) Header() http.Header { return r.header }

func (r *recorder) WriteHeader(code int) {
	if r.code == 0 {
		r.code = code
	}
}

func (r *recorder) Write(b []byte) (int, error) {
	if r.code == 0 {
		r.code = http.StatusOK
	}
	return r.body.Write(b)
}
