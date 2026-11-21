# tifa-web

The source of tifago.com, the website of [Tifa](https://github.com/liustack/tifa). It is a Tifa app: Go serves the pages and the search, and React renders them.

The docs on the site are the docs of Tifa, from one source. `go generate ./...` reads `docs/` of the Tifa version in `go.mod`, keeps the docs for people who write apps, and writes them into `content/docs/`, which is generated and not in git.

You need Go 1.27, and Node 20.19 or 22.12 and later with [pnpm](https://pnpm.io).

| Command | What it does |
| --- | --- |
| `pnpm install` | Install the npm packages and write the docs into `content/docs/` |
| `go generate ./...` | Write the docs again, after the Tifa version in `go.mod` changed |
| `go tool tifa dev` | Run the site at http://127.0.0.1:9527, rebuilding on every change |
| `go tool tifa check` | Check the routes, vet and test the Go code, and type-check the TypeScript |
| `go tool tifa build` | Build the single executable into `bin/` |

The words of the home page are in `components/site.ts`. [AGENTS.md](AGENTS.md) says how the rest is laid out.

This repository is open source under the [MIT license](LICENSE) and takes no pull requests, like Tifa. See [CONTRIBUTING.md](CONTRIBUTING.md).
