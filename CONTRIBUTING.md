# Contributing to tifa-web

First, the policy: **this repository does not accept pull requests.** It is the source of tifago.com, the website of [Tifa](https://github.com/liustack/tifa), and the same maintainer writes, reviews and owns every line, as in Tifa itself.

Two ways to help:

- **[Open an issue](https://github.com/liustack/tifa-web/issues/new).** A broken link, a page that renders wrong, a doc that is unclear. The docs themselves live in the Tifa repository, so for a mistake in the text of a doc, [open the issue there](https://github.com/liustack/tifa/issues/new/choose).
- **Fork it.** Under the MIT license your copy is yours to change and publish. No permission needed.

Everything below is for people working on a fork.

## Setup

You need Go 1.27, and Node 20.19 or 22.12 and later with pnpm.

```sh
pnpm install          # also writes the docs into content/docs, see below
go tool tifa dev      # http://127.0.0.1:9527
go tool tifa check    # routes, vet, tests and the TypeScript check
go tool tifa build    # bin/tifa-web
```

The docs are the docs of the Tifa version in `go.mod`, written into `content/docs/` by `go generate ./...` (`cmd/syncdocs`). Run it again after you change that version. See AGENTS.md for how the site is laid out.

## Commits

- `type(scope): imperative summary`, at most 72 characters, no period at the end.
- One commit does one thing. Keep a refactor or a reformat apart from a change of behavior.
