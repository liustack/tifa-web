# Notes for coding agents

<!-- tifa:begin -->
<!-- Written by Tifa v0.1.18 with go tool tifa agents. After an upgrade run it again: it refreshes the text between these markers and leaves the rest of this file alone. Put your own notes outside it. -->

This is a [tifa](https://github.com/liustack/tifa) app. Go routes requests, loads data and runs the hooks. React pages render to HTML inside the Go process (slowjs, no Node at run time), and TanStack Router takes over in the browser. The build bundles the front end with Vite, on Node, and makes one executable that needs no Node.

## Read the docs of this version first

Run `go tool tifa docs`. It prints the folder with the docs of the tifa version this app uses and one line for each doc about writing an app. Add `--when` to see when to read each, and `--all` to also list the docs about developing tifa, which an app does not need. Read the docs that apply before you change code. Do not read the docs on GitHub, which are those of the main branch and may differ from this version. The `docs/...` paths in these notes are in that folder.

## Commands

| Command | What it does |
| --- | --- |
| `pnpm install` | Install the npm packages (React, TanStack Router, TypeScript, Vite). Keep `@tanstack/react-router` at the version in `package.json`: `tifa build` refuses another one, since tifa is tested with that version only. It also runs `go tool tifa gen` |
| `go tool tifa gen` | Write the generated code into `.tifa/` and build nothing, so editors and `tsc` see its types. While `tifa dev` runs it keeps `.tifa/` and this command leaves it alone |
| `go tool tifa dev --background` | Start the dev server detached (plain `tifa dev` would hold your terminal) and return once it answers, at http://127.0.0.1:9527. Then `go tool tifa dev status` (`--json`, exit 1 when not running), `go tool tifa dev logs` (`-n 100`, `-f`, all output is in `.tifa/dev.log`) and `go tool tifa dev stop`. Stop it when you finish. A TSX or CSS edit updates the open page in place through Vite, a Go or `tifa.yaml` edit rebuilds and reloads it |
| `go tool tifa errors` | While `tifa dev` runs, print the errors since the last build: a failed build, unexpected request errors, errors of the browser (hydration mismatches, uncaught errors, unhandled rejections) and pages that render differently each time. Run it after you change code and after you open or reload the pages you changed. It exits with 1 when there are errors, 2 when no `tifa dev` runs. `--json` prints JSON |
| `go tool tifa build` | Generate `.tifa/`, bundle the front end with Vite and build `bin/` |
| `go tool tifa check` | The one command to run before you finish. It checks the routing conventions, generates `.tifa/`, then runs `go vet ./...`, `go test ./...` and the TypeScript type check with the app's own TypeScript. It stops at the first step that fails and says which |
| `go tool tifa routes` | List the routes |
| `go tool tifa docs` | Print the docs of the tifa version this app uses for writing an app (`--when` for when to read each, `--all` for the docs about developing tifa) |
| `go tool tifa agents` | Write these notes and the tifa skill in `.claude/skills/tifa/` again from the tifa version in `go.mod`. Run it after you upgrade tifa. `--check` changes nothing and exits with 1 when they are out of date |
| `go tool tifa upgrade [version]` | Move this app to a newer tifa (the latest without a version): `go get -tool` and `go mod tidy`, then the new tifa sets the versions in `package.json` that tifa new pins, installs, writes these notes and the skill again, prints the `Upgrading:` notes of the changelog between the two versions for you to act on, and runs `tifa check`. It changes no code of the app. Use it, not a hand-made `go get`, to upgrade tifa |

Never edit `.tifa/`. It is generated, and the Go tools skip it.

## Routes

- Every folder under `app/` is a URL segment. `app/` is `/`. `name+` is a parameter, and the folder name is the parameter name: `booking+` under `staff/` is `/staff/{booking}`, not `/staff/booking/{id}` (for that, make a literal folder `booking` with `id+` inside). It is read with `c.Param("name")`. `go tool tifa routes` prints the URLs that result. `name++` is the rest of the path, last folder only, `_splat` in TypeScript. Name the Go package after the parameter: `todo+` holds `package todo`.
- A folder has at most three convention files. `route.go` holds every Go name of the folder: `Load` (the page data), `Layout` (the data of the Layout), `LayoutSSR`, `SSR`, `Memo` and the method handlers `GET`, `POST`, `PUT`, `PATCH`, `DELETE`. These names in any other Go file of the package fail the build. `page.tsx` is the page. `layout.tsx` (with `<Outlet />`) wraps the page of its folder and every page below.
- A folder with `page.tsx` is a page, so it cannot define `GET`. A folder with `route.go` and no `page.tsx` is an API route when it has a method handler, and has no `Load`. Its `Layout` still gives the pages below it their layout data. `tifa.Use(c, pkg.Layout)` gets the data of an enclosing Layout, run once per request.
- Write Tifa, not another framework: no `next/link`, `'use client'`, `loader`, `getServerSideProps`, `createServerFn` or `[id]` folders. `go tool tifa check` reports them with the Tifa way. See "Coming from another framework" in `docs/app.md`.
- Details: `docs/app.md`, which links to the doc of each topic.

## Returned errors

| Returned | A page (`Load`, `Layout`) | A method handler |
| --- | --- | --- |
| `tifa.ErrNotFound`, or a `tifa.Error` with Status 404 | 404 shown by the first `ErrorBoundary` going up from the page (from the Layout one level up when a `Layout` returns it), which gets a `TifaError` with `status` 404 | 404 as JSON |
| `tifa.Error{Status, Message}`, Status 400 to 599 but not 404 | That status, shown by the nearest `ErrorBoundary` | That status as JSON `{"error": message}` |
| `tifa.Fields{"field": "message"}` | 422, shown by the nearest `ErrorBoundary` | 422 as JSON with `"fields"` |
| `tifa.Redirect(url, status)` | A redirect | Called from the browser, the browser goes there: through the router for a URL of this site, a full page load for another. Otherwise a 3xx |
| Any other error, a panic, or a `tifa.Error` with another status | `OnError`, or a logged 500 without one | `OnError`, or a logged 500 without one |

`Validate` on a handler input runs before the handler. A `tifa.Error` or `tifa.Fields` it returns is answered as in the table. Any other error, also `tifa.ErrNotFound` or a `tifa.Redirect`, is answered with 400 and `err.Error()` as the message. A URL that matches no page shows the `ErrorBoundary` of the deepest Layout it matches (the root `layout.tsx` here), inside that Layout, as a `TifaError` with `status` 404. `NotFound` is no export: exporting it fails the build.

## Load and Ctx

The build reads `Load` and `Layout` and the functions they call to find which parameters each page depends on. It fails with the file and line when they:

- call `c.Param`, `c.Query` or `c.QueryValues` with a string that is not a constant, or not directly (stored in a variable or passed as a value),
- call `tifa.Use` with anything but a package-level `Layout` that wraps the caller,
- convert, store, return or send `c` anywhere else than a plain function or method of this module or of tifa, or a new local variable of the same type (see the docs for the full list),
- call `c.Request()` or `c.Writer()`,
- call `c.BindQuery` or `c.Bind`: read query values with `c.Query("name")` or `c.QueryValues("name")`.

Writing the response is checked when the request runs instead: in `Load` and `Layout`, `c.JSON`, `c.Text` and `c.NoContent` write nothing and return an error. Returned, that error is an unexpected one (`OnError`, or a 500). Setting headers and cookies is fine. Pass `c.Context()` when a function only needs the context. See "How Load may use Ctx" in `docs/routes.md`.

## Handlers and the browser

- A typed handler takes an input after the Ctx, returns a result before the error, or both: `func POST(c *tifa.Ctx, in AddInput) (Todo, error)`. The framework decodes `in` from a JSON or form body (from the query string for `GET` and `DELETE`), runs `in.Validate()` if it exists (`func() error`, or `func(c *tifa.Ctx) error` to read the request, such as `c.Language()`), and writes the result as JSON (204 when there is no result). `func POST(c *tifa.Ctx) error` writes its own response.
- The data of `Load` and handler inputs and results travel as JSON, and tifa writes their TypeScript types. No interfaces, funcs or channels. A Go slice or map can be `null` in TypeScript, so write `?? []`. An integer a browser would round (above 9007199254740991, such as a snowflake ID or `time.UnixNano()`) is refused with an error when a page reads it: the results of `Load` and `Layout`, session data, and a handler result for a call that sends `X-Tifa-Call` (the generated function of a page). `tifatest.Start` refuses it in `Load` and `Layout` too. A plain HTTP client and curl get the raw number from a handler, and so does `tifatest.Start` unless the request sets `X-Tifa-Call: 1`, so in a test of a handler set that header to see the refusal. Add `json:",string"` to the field, or send milliseconds.
- In TSX, import from the generated `./route.gen` of the folder: `page.useLoaderData()` and `layout.useLoaderData()`, the `PageData` and `LayoutData` types, and one function per handler, such as `POST({ params, input })`. A page reads the data of an enclosing Layout from the `route.gen` of that Layout's folder, such as `../route.gen`. A call never reloads page data: after a write, call `router.invalidate()` yourself. See "Calling method handlers from the browser" in `docs/pages.md#calling-method-handlers-from-the-browser`.
- The function of a typed handler resolves to its result and throws `TifaError` (`status`, `message`, `fields`) for an error answer. The function of a handler that writes its own response resolves to the `Response`, also for a 4xx or 5xx, so check `response.ok`. When the handler or `OnRequest` redirects, the function resolves to `undefined` once the new page shows, so its result type is `Out | undefined`: check the result before using it. Redirect with `return tifa.Redirect(...)`, not by writing a 3xx yourself.
- Styles are CSS files imported from TSX, such as `app/styles.css` in the root `layout.tsx`. `vite.config.ts` adds the Vite plugins of tifa and React. For Tailwind CSS, add `@tailwindcss/vite` there. See "Styles" in `docs/styles.md`.
- A page may also export `head` and `ErrorBoundary`, which also shows its 404 as a `TifaError` with `status` 404. In the browser each `page.tsx` and `layout.tsx` but the root one loads with its route, also one that exports `head`: its top-level code runs when its route first loads, not when the app starts, so setup every page needs goes in the root `layout.tsx` or a module it imports. Components that need the browser at import time, or show nothing until the browser draws them, go in `*.client.tsx` and load with `clientOnly` from `tifa/client`, with a placeholder of the same size. See "Components that draw only in the browser" in `docs/server-render.md`.

## Hooks

`hooks.go` in the root package, all optional:

- `OnStart(ctx)` runs before the app listens. An error stops the app. It is where the app wires what the whole process shares, such as `tifa.SetSitemap`, which lists the pages with parameters for `/sitemap.xml` and `/llms.txt` once `site.url` is set in `tifa.yaml`. See `docs/guides/seo.md`.
- `OnShutdown(ctx)` runs after a stop signal, once the requests in flight have finished. Requests still running after `server.shutdown_timeout` are canceled and get one more `shutdown_timeout` to return. `OnShutdown` is skipped if they do not, and either way the app exits with status 1. A handler that streams with no end, such as server-sent events, returns when `tifa.ShuttingDown(c.Context())` closes, see `docs/guides/streaming.md`.
- `OnRequest(c, next)` runs before pages, navigation data requests, API routes and unmatched URLs, not static files. Call `next(c)` to go on, or `next(c.WithContext(ctx))` to pass a request-scoped value.
- `OnError(c, err)` gets unexpected errors and panics. Writing one replaces the default, which logs and answers 500.

## Tests

- Domain code: plain Go tests next to it, such as `app/todos/todo_test.go`. Each test gets a fresh database with `useTestDB(t)`, which replaces `tifaweb.DB`, so these tests do not run in parallel.
- Route code: `tifatest.NewCtx(t, httptest.NewRequest(...))` gives a `*tifa.Ctx`. Set path parameters with `r.SetPathValue`. Run `Load` and `Layout` through `tifa.Use(c, Load)`, call handlers directly, and test `Validate` on its own. These tests check what the function returns, not the HTTP answer. See "Testing" in `docs/testing-apps.md`.
- Whole requests: `tifatest.Start(t, tifaapp.App)` runs the app and returns an `*httptest.Server` whose client keeps cookies and does not follow redirects, for a flow such as posting a form and following the redirect. Import `github.com/liustack/tifa-web/.tifa/tifaapp` from a test in a package of its own (`package app_test`), set `tifaweb.DBPath` to a file in `t.TempDir()` before `Start`, and do not call `t.Parallel` in such a test. Pages answer the data of their Loads as JSON. `tifatest.StartPages` renders them, and needs `go tool tifa build` first. See "Testing a whole request" in `docs/testing-apps.md`.

## Settings and data

- `tifa.yaml` lists every setting commented out at its default, with an example for `ssr.vms`, whose default is GOMAXPROCS, at most 4. See `docs/config.md`. A setting `server.addr` is also the flag `--server-addr` and the environment variable `SERVER_ADDR`, and the app defines its own options in `options.go` with `tifa.String`, `tifa.Int` and the like.
- The database is `app.db`, or what the option `db.path` of `options.go` says (`DB_PATH`, `--db-path`). `schema.sql` runs at every start, so keep its statements safe to repeat (`CREATE TABLE IF NOT EXISTS`). A new domain adds its tables there. `tifa dev` does not watch `schema.sql`: after changing it, restart `tifa dev`.
- Under `tifa dev`, a page whose server render fails, or whose `OnRequest` returns an unexpected error, answers with its status and the whole error, a render error with its stack trace in the source files. `curl -i` on the page shows them. `go tool tifa errors` lists them too. A built app shows only the status and the public message.
- Build and generated files: `docs/build.md`. `tifa dev`: `docs/dev.md`.

<!-- tifa:end -->

## This repository

This is the source of tifago.com, the website of [Tifa](https://github.com/liustack/tifa), written as a Tifa app. It has no database and no sessions. Like Tifa it is open source (MIT) and takes no pull requests, see CONTRIBUTING.md.

- **The docs are written in the Tifa repository, not here.** `content/docs/en-US/` is generated by `go generate ./...` and committed, so `tifa deploy` builds it with no step of its own, and CI fails when it does not match the Tifa version in `go.mod`. `cmd/syncdocs` also writes the changelog of that version to `content/site/changelog.md`. It reads `docs/` of the Tifa module this app depends on (`go list -m -f '{{.Dir}}' github.com/liustack/tifa`), keeps the docs with `audience: app` and writes them as pages: the title is the first `# ` heading, the description is the first part of `summary`, and the section and the place in it come from the `Places` table in `cmd/syncdocs/sync.go`. Links between docs become `/docs/<section>/<name>`, and links to anything else in the Tifa repository become GitHub links at the version in `go.mod`. Never edit a file under `content/docs/en-US/`, and do not copy docs by hand. `syncdocs` also writes a `sourceHash` into each English doc: the first 12 hex digits of the SHA-256 of the file in the Tifa module, before it is changed. To change a doc, change it in the Tifa repository.
- **After `pnpm install` the docs are there**, since its `prepare` script runs `go generate ./...` (which runs `cmd/syncdocs`) and then `go tool tifa gen`. Run `go generate ./...` yourself after you change the Tifa version in `go.mod`, or when `tifa check` or `tifa build` says a folder under `content/` has no collection or a file is missing. A new doc for apps in a newer Tifa fails the sync until you add it to `Places`: that is on purpose.
- **Languages.** English is the default with no prefix, Chinese is `/zh-CN` (`i18n.default` in `tifa.yaml`, one `app/_locales/<language>.json` each, read with `useTranslation` from `tifa/i18n`). The hero (`hero.eyebrow`, `hero.tagline` and `hero.sub`) and every other word of the home page are in those files. The links, the terminal lines and the numbers are in `components/site.ts`, and the code of the pages section in `components/landing-code.ts`. The terminals show lines the tools print: the hero one is what `tifa new` and `tifa dev` print (`cmd/tifa/new.go`, `dev.go` of the Tifa version in `go.mod`), the deploy one is the log of a real deploy of this site. Change a line only by cutting it or by copying a newer real one. Chinese text follows the style rules of the owner: no semicolons and no dashes as punctuation. `content/pages/<language>/` holds Why and Examples. `Docs` is one collection with the languages of the site: `content/docs/en-US/` is generated and `content/docs/zh-CN/<path of the doc>.md` is the hand-written Chinese version, which `syncdocs` never touches. A Chinese file has the same front matter keys as the English one, a translated `title` and `description`, and `translatedFrom`, the `sourceHash` of the English doc it was translated from. `docs.go` shows the Chinese entry when there is one (Tifa itself does not fall back, so the site does) and otherwise the English body with the note `docs.englishOnly`. When `translatedFrom` differs from the English `sourceHash`, the Chinese page shows `docs.outdated` with a link to the English page, and `go generate` prints one line for each such file: translate it again, then set `translatedFrom` to the new `sourceHash`. The sidebar, its groups and their order always come from the English entry. Chinese headings end in `{#english-id}`, which Tifa does not support yet, so leave them as they are. Text and links copied by hand name the Tifa version, and `TestHandWrittenPagesNameTheVersionOfGoMod` fails after an upgrade until they are read again.
- **Pages.** `app/page.tsx` is the home page (its styles are `app/landing.css`, all under `.lp`, and the fonts are IBM Plex from `@fontsource`), `app/why`, `app/examples` and `app/changelog` are text pages, and `components/layout-options.tsx` is the header. `app/docs/` is the docs: `route.go` gives the sidebar tree (built in `tree.go` from the collection in `content.go`), `path++/` is one doc, and `app/api/search/` is the search box (a Go handler that Fumadocs' client calls). Fumadocs UI gives the layouts, the search dialog and the dark mode.
- **Code colors** come from `@shikijs/rehype` in `vite.config.ts`, with a light and a dark theme. `app/styles.css` styles the result.
- **SEO** comes from `site.url` in `tifa.yaml` (canonical links, `/sitemap.xml`, `/robots.txt`, `/llms.txt`, the `.md` version of each doc). `sitemap.go` lists the docs for them.
- **Deploy.** `tifa.yaml` has no `deploy` block yet. Do not run `tifa deploy`.
- **No git worktree.** Use `git switch -c <branch>`. Use pnpm, never npm or npx (`pnpm dlx` for one-off tools).
- **Builds are heavy.** Run `go tool tifa build` and `go tool tifa check` as `~/.cache/mac-locks/heavy-run nice -n 10 <command>`, which waits when another agent holds the lock.
- **Commits.** `type(scope): imperative summary`, at most 72 characters, no period, no attribution lines. One thing per commit. Add files by name, never `git add -A` or `git add .`, and build with `go build -o` outside the checkout (`bin/` is ignored, but a 41 MB binary once went into a release).
- **Before you finish**, run `go tool tifa check`. It also runs the test that checks every link between docs: a link to a doc or a heading that does not exist fails it.
- English only.
