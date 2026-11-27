# Notes for coding agents

<!-- tifa:begin -->
<!-- Written by Tifa v0.8.0 with go tool tifa agents. After an upgrade run it again: it refreshes the text between these markers and leaves the rest of this file alone. Put your own notes outside it. -->

This is a [Tifa](https://github.com/liustack/tifa) app. Go routes requests, loads data and runs the hooks. React pages render to HTML inside the Go process (beastjs, no Node at run time), and TanStack Router takes over in the browser. The build bundles the front end with Vite on Node into one executable that needs no Node.

Before you change code, read the docs that apply. [Docs](#docs) at the end of this block names them, and `go tool tifa docs` says what each one covers. They are those of the Tifa version in `go.mod`. Not those on GitHub, which are for main.

## Commands

| Command | What it does |
| --- | --- |
| `pnpm install` | Install the npm packages, then run `go tool tifa gen`. Keep `@tanstack/react-router` at the version in `package.json`: `tifa build` refuses another |
| `go mod tidy` | After you add or remove a Go import. Until then `go build` and `go tool tifa` stop with "missing go.sum entry" |
| `go tool tifa gen` | Write the generated types into `.tifa/`. Run it once after a batch of Go changes, before the TS that uses them. `tifa dev` does it itself |
| `go tool tifa dev --background` | Start the dev server detached (plain `tifa dev` holds the terminal) and return once it answers at http://127.0.0.1:9527. Then `go tool tifa dev status`, `logs` and `stop`. Stop it when you finish |
| `go tool tifa errors` | While `tifa dev` runs, print the errors since the last build, those of the browser too. Run it after a change and after you open the pages you changed. Exits with 1 when there are errors |
| `go tool tifa check` | The one command to run before you finish: the routing conventions, `.tifa/`, `go vet ./...`, `go test ./...` and the TypeScript check. It stops at the first step that fails |
| `go tool tifa build` | Generate `.tifa/`, bundle the front end and build `bin/` |
| `go tool tifa routes` | List the routes, their URLs and the addresses Tifa answers itself |
| `go tool tifa docs` | Print the docs folder of this Tifa version and what each doc covers |
| `go tool tifa agents` | Write this block and the skill in `.claude/skills/tifa/` again for the Tifa in `go.mod` |
| `go tool tifa upgrade [version]` | Move the app to a newer Tifa (not a hand-made `go get`) and print the `Upgrading:` notes to act on |

Never edit `.tifa/`. It is generated.

## Routes

- Every folder under `app/` is a URL segment, and `app/` is `/`. `name+` is a parameter named after the folder: `booking+` under `staff/` is `/staff/{booking}`, read with `c.Param("booking")`. For `/staff/booking/{id}`, make a folder `booking` with `id+` inside. `name++` is the rest of the path, last folder only, `_splat` in TypeScript. Name the Go package after the parameter: `todo+` holds `package todo`.
- A folder has at most three convention files. `route.go` holds every Go name of the folder: `Load` (the page data), `Layout` (the data of the Layout), `SSR`, `Memo` and the method handlers `GET`, `POST`, `PUT`, `PATCH`, `DELETE`. These names in another Go file of the package fail the build. `page.tsx` is the page. `layout.tsx` (with `<Outlet />`) wraps the page of its folder and every page below.
- A folder with `page.tsx` cannot define `GET`. A folder with `route.go` and no `page.tsx` is an API route when it has a method handler, and has no `Load`. `tifa.Use(c, pkg.Layout)` gets the data of an enclosing Layout, run once per request.
- Write Tifa, not another framework: no `next/link`, `'use client'`, `loader`, `getServerSideProps`, `createServerFn` or `[id]` folders. `go tool tifa check` reports them with the Tifa way. See "Coming from another framework" in `docs/start/app.md`.

## Go side

- Return `tifa.Redirect(url, status)` to redirect. Return `tifa.ErrNotFound`, `tifa.Error{Status, Message}` or `tifa.Fields{"field": "message"}` (422) for an error that a page shows in the nearest `ErrorBoundary` as a `TifaError` and a handler answers as JSON. Any other error or a panic goes to `OnError`, or a logged 500. See `docs/pages-and-data/errors.md`.
- A typed handler: `func POST(c *tifa.Ctx, in AddInput) (Todo, error)`. Tifa decodes `in` from a JSON or form body (the query string for `GET` and `DELETE`), runs `in.Validate()` when it exists and writes the result as JSON. `func POST(c *tifa.Ctx) error` writes its own answer.
- Data, inputs and results travel as JSON, and Tifa writes their TypeScript types. No interfaces, funcs or channels. A Go slice or map can be `null`, so write `?? []`. An integer above 9007199254740991 is refused where a page reads it: add `json:",string"` or send milliseconds.
- The build fails with the file and line when `Load` or `Layout`, or a helper it passes `c` to:
  - calls `c.Query`, `c.QueryValues`, `c.Bind` or `c.BindQuery` (take `q Filter`, a struct, as the second parameter),
  - calls `c.JSON`, `c.Text` or `c.NoContent` (return the data or an error),
  - calls `tifa.Use` with anything but a package-level `Layout` that wraps the caller,
  - calls `c.Path()` in a `Layout`.

  Do not keep `c` after the function returns, and wait for your goroutines, as in `net/http`. There `c.Request()` is a copy with only the declared query parameters and no body. See "How Load may use Ctx" in `docs/pages-and-data/loading-data.md`.

## React side

- Import from the generated `./route.gen` of the folder: `page.useLoaderData()`, `layout.useLoaderData()`, the `PageData` and `LayoutData` types, their Go structs by name (`Item`), and one function per handler, such as `POST({ params, input })`. A page reads the data of an enclosing Layout from the `route.gen` of that folder, such as `../route.gen`.
- A handler function throws `TifaError` (`status`, `message`, `fields`) for an error answer. It never reloads page data: after a write, call `router.invalidate()`. See `docs/pages-and-data/handlers.md#calling-method-handlers-from-the-browser`.
- Components that need the browser at import time go in `*.client.tsx` and load with `clientOnly` from `tifa/client`. See `docs/pages-and-data/server-render.md`.
- `page.tsx` may also export `head` and `ErrorBoundary`. Styles are CSS files imported from TSX, see `docs/pages-and-data/styles.md`.

## Hooks

`hooks.go` in the root package, all optional: `OnStart(ctx)` wires what the whole process shares (an error stops the app), `OnShutdown(ctx)`, `OnRequest(c, next)` and `OnError(c, err)`. See `docs/pages-and-data/hooks.md`.

## Tests

- Domain code: plain Go tests next to it. `useTestDB(t)` gives a test a fresh database by replacing `tifaweb.DB`, so these tests do not run in parallel.
- Route code: `tifatest.NewCtx(t, httptest.NewRequest(...))` gives a `*tifa.Ctx`. Set path parameters with `r.SetPathValue`, run `Load` through `tifa.Use(c, Load)`, call handlers directly and test `Validate` on its own.
- Whole requests: `tifatest.Start(t, tifaapp.App)`, from a test in a package of its own (`package app_test`) that imports `github.com/liustack/tifa-web/.tifa/tifaapp`. Set `tifaweb.DBPath` to a file in `t.TempDir()` first, and do not call `t.Parallel`. Pages answer their data as JSON there. `tifatest.StartPages` renders them, after `go tool tifa build`.
- See `docs/pages-and-data/testing-apps.md`.

## Settings and data

- `tifa.yaml` lists every setting at its default, and `options.go` holds the options of the app. See `docs/build-and-deploy/config.md`.
- `deploy.front` is `cloudflare` in `tifa.yaml` (needs a `site.url` domain). Others: `deploy.md`.
- `db.go` opens SQLite with `_txlock=immediate`, so a transaction takes the write lock at `BEGIN` and never fails with `database is locked`. Keep it, and start a transaction that only reads with `sql.TxOptions{ReadOnly: true}`. See `docs/guides/database.md`.
- The option `db.backup` (`DB_BACKUP`) copies every commit to S3 or R2 and restores a missing database at start. See "Back it up" in `docs/guides/database.md`.
- Under `tifa dev`, a page that fails answers with the whole error: `curl -i` the page, or run `go tool tifa errors`.

## Docs

`go tool tifa docs` prints the folder of the Tifa module in `go.mod` and one line on each doc. Open the files that apply. Each has a Chinese copy beside it, `<name>.zh-CN.md`.

- `docs/start/`: `tutorial.md`, `app.md`, `philosophy.md`, `structure.md`, `agents.md`
- `docs/pages-and-data/`: `routing.md`, `layouts.md`, `loading-data.md`, `handlers.md`, `navigation.md`, `head.md`, `errors.md`, `styles.md`, `content.md`, `sessions.md`, `hooks.md`, `server-render.md`, `ssr-and-cache.md`, `testing-apps.md`, `testing-pages.md`
- `docs/guides/`: `forms.md`, `database.md`, `sign-in.md`, `uploads.md`, `email.md`, `jobs.md`, `i18n.md`, `seo.md`, `streaming.md`, `environment.md`, `logs.md`, `security.md`, `sign-up.md`, `password-reset.md`, `oauth.md`, `flash-messages.md`, `cors.md`, `websocket.md`, `commands.md`
- `docs/build-and-deploy/`: `create-app.md`, `dev.md`, `dev-reloads.md`, `dev-errors.md`, `dev-scripts.md`, `debugging.md`, `check.md`, `build.md`, `upgrade.md`, `deploy.md`, `deploy-setup.md`, `deploy-github-actions.md`, `deploy-cloudflare.md`, `deploy-proxy.md`, `deploy-operate.md`, `deploy-by-hand.md`, `config.md`, `tuning.md`, `https.md`, `headers-and-origins.md`, `kubernetes.md`, `deploy-harden.md`, `deploy-several-servers.md`, `containers.md`
- `docs/reference/`: `conventions.md`, `memo.md`, `performance.md`, `compatibility.md`, `ctx.md`, `json.md`, `settings.md`, `generated-files.md`, `build-output.md`, `route-modules.md`, `route-tree.md`, `data-protocol.md`, `page-requests.md`, `server-program.md`, `i18next.md`

<!-- tifa:end -->

## This repository

This is the source of tifago.com, the website of [Tifa](https://github.com/liustack/tifa), written as a Tifa app. It has no database and no sessions. Like Tifa it is open source (MIT) and takes no pull requests, see CONTRIBUTING.md.

- **The docs are written in the Tifa repository, not here.** `content/docs/en-US/` and `content/docs/zh-CN/` are generated by `go generate ./...` and committed, so `tifa deploy` builds them with no step of its own, and CI fails when they do not match the Tifa version in `go.mod`. `cmd/syncdocs` also writes the changelog of that version to `content/site/changelog.md`. It reads `docs/` of the Tifa module this app depends on (`go list -m -f '{{.Dir}}' github.com/liustack/tifa`). A folder there whose `index.md` says `audience: app` is a section of the sidebar (its `title` and `order`, and the Chinese title from `index.zh-CN.md`), and each doc in it with `audience: app` is a page at `/docs/<folder>/<name>` with the place its `order` says. The English page is `<name>.md` and the Chinese page is the `<name>.zh-CN.md` that Tifa ships beside it. A doc without its Chinese file, or without `order`, fails the sync. Folders without such an index (`internals/`, `studies/`, `research/`) are not on the site. The title is the first `# ` heading (in Chinese without its `{#id}`), the description is the first part of `summary`. Links between docs become `/docs/<folder>/<name>` in both languages, and links to anything else in the Tifa repository become GitHub links at the version in `go.mod`. Never edit a file under `content/docs/`, and do not copy docs by hand. To change a doc or its translation, change it in the Tifa repository.
- **After `pnpm install` the docs are there**, since its `prepare` script runs `go generate ./...` (which runs `cmd/syncdocs`) and then `go tool tifa gen`. Run `go generate ./...` yourself after you change the Tifa version in `go.mod`, or when `tifa check` or `tifa build` says a folder under `content/` has no collection or a file is missing. A new doc or section in a newer Tifa needs no change here.
- **Languages.** English is the default with no prefix, Chinese is `/zh-CN` (`i18n.default` in `tifa.yaml`, one `app/_locales/<language>.json` each, read with `useTranslation` from `tifa/i18n`). The hero (`hero.eyebrow`, `hero.tagline` and `hero.sub`) and every other word of the home page are in those files. The links, the terminal lines and the numbers are in `components/site.ts`, and the code of the pages section in `components/landing-code.ts`. The terminals show lines the tools print: the hero one is what `tifa new` and `tifa dev` print (`cmd/tifa/new.go`, `dev.go` of the Tifa version in `go.mod`), the deploy one is the log of a real deploy of this site. Change a line only by cutting it or by copying a newer real one. Chinese text follows the style rules of the owner: no semicolons and no dashes as punctuation. `content/pages/<language>/` holds Why and Examples. `Docs` is one collection with the languages of the site, and both folders are generated. The sidebar, its section names and their order come from the entries of the language of the page. Chinese headings end in `{#english-id}`, which gives them the id of the English heading, so links to a heading work in both languages. `go generate` rewrites the version in the links to the Tifa repository (`github.com/liustack/tifa/blob/<version>/` and `.../tree/<version>/`) in `app/_locales` and `content/pages`, changes nothing else, and prints one line for each file it changed. `TestHandWrittenPagesNameTheVersionOfGoMod` fails when a link with another version was added by hand and `go generate` was not run. The Chinese site shows the changelog in English with the note `docs.englishOnly`.
- **Pages.** `app/page.tsx` is the home page (its styles are `app/landing.css`, all under `.lp`, and the fonts are IBM Plex from `@fontsource`), `app/why`, `app/examples` and `app/changelog` are text pages, and `components/layout-options.tsx` is the header. `app/docs/` is the docs: `route.go` gives the sidebar tree (built in `tree.go` from the collection in `content.go`), `path++/` is one doc, and `app/api/search/` is the search box (a Go handler that Fumadocs' client calls). Fumadocs UI gives the layouts, the search dialog and the dark mode.
- **Code colors** come from `@shikijs/rehype` in `vite.config.ts`, with a light and a dark theme. `app/styles.css` styles the result.
- **SEO** comes from `site.url` in `tifa.yaml` (canonical links, `/sitemap.xml`, `/robots.txt`, `/llms.txt`, the `.md` version of each doc). `sitemap.go` lists the docs for them.
- **Deploy.** `tifa.yaml` has no `deploy` block yet. Do not run `tifa deploy`.
- **No git worktree.** Use `git switch -c <branch>`. Use pnpm, never npm or npx (`pnpm dlx` for one-off tools).
- **Builds are heavy.** Run `go tool tifa build` and `go tool tifa check` as `~/.cache/mac-locks/heavy-run nice -n 10 <command>`, which waits when another agent holds the lock.
- **Commits.** `type(scope): imperative summary`, at most 72 characters, no period, no attribution lines. One thing per commit. Add files by name, never `git add -A` or `git add .`, and build with `go build -o` outside the checkout (`bin/` is ignored, but a 41 MB binary once went into a release).
- **Before you finish**, run `go tool tifa check`. It also runs the test that checks every link between docs: a link to a doc or a heading that does not exist fails it.
- English only.
