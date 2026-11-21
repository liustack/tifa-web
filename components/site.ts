// The links of the site and what is the same in every language: the commands, the numbers and
// the files of the demo. The words are in app/_locales, one file per language, and the home page
// is app/page.tsx. The hero (key hero in those files) is the first words a visitor reads. Every
// sentence is from the README of the Tifa repository, from its docs or from the changelog.

const tifa = 'https://github.com/liustack/tifa'

export const site = {
  name: 'Tifa',
  url: 'https://tifago.com',
  links: {
    docs: '/docs/start/tutorial', // the first doc, which /docs leads to. TestFirstDocIsTheTutorial in cmd/syncdocs keeps it right
    github: tifa,
    issues: `${tifa}/issues/new/choose`,
    x: 'https://x.com/liustack',
    source: 'https://github.com/liustack/tifa-web',
  },
  // The command of the hero.
  newApp: 'go run github.com/liustack/tifa/cmd/tifa@latest new myapp',
  // The commands that close the home page, from docs/tutorial.md.
  quickstart: `go run github.com/liustack/tifa/cmd/tifa@latest new tutorial-app
cd tutorial-app
go mod tidy
pnpm install
go tool tifa dev     # http://127.0.0.1:9527`,
  // Whole visits to the example shop, one core each (docs/performance.md, 2026-11-18): page views,
  // full loads and client navigations a second. The name of each row is in the language files.
  throughput: [
    ['tifa', 118.8, 38.3, 80.5],
    ['next', 28.3, 9.1, 19.2],
    ['tanstack', 94.7, 30.5, 64.2],
    ['tifaNoMemo', 78.0, 25.1, 52.8],
    ['tanstackGzip', 64.5, 20.8, 43.7],
    ['nextNoPrefetch', 61.7, 19.9, 41.8],
  ] as const,
}

// Where a table row is one of the reference rows, which are not the frameworks as shipped.
export const referenceRows = new Set(['tifaNoMemo', 'tanstackGzip', 'nextNoPrefetch'])
