// The links of the site and what is the same in every language: the commands, the terminal lines
// and the numbers. The words are in app/_locales, one file per language, and the home page is
// app/page.tsx. The hero (key hero in those files) is the first words a visitor reads. Every
// sentence is from the README of the Tifa repository, from its docs or from the changelog.

const tifa = 'https://github.com/liustack/tifa'

export const site = {
  name: 'Tifa',
  url: 'https://tifago.com',
  links: {
    docs: '/docs/start/tutorial', // the first doc, which /docs leads to. TestFirstDocIsTheTutorial in cmd/syncdocs keeps it right
    performance: '/docs/reference/performance',
    why: '/why',
    github: tifa,
    x: 'https://x.com/liustack',
    source: 'https://github.com/liustack/tifa-web',
  },
  // The command of the hero and of the end of the page.
  newApp: 'go run github.com/liustack/tifa/cmd/tifa@latest new myapp',
}

// A line of a terminal: what was typed, what the program printed, the line that says it went well,
// or a cut (…).
export type TerminalLine = { kind: 'cmd' | 'out' | 'ok' | 'dim'; text: string }

// What `tifa new myapp` and `tifa dev` print at v0.1.17, with what the other tools print cut to
// an ellipsis (cmd/tifa/new.go, cmd/tifa/dev.go).
export const heroTerminal: TerminalLine[] = [
  { kind: 'cmd', text: site.newApp },
  { kind: 'out', text: 'created myapp in myapp' },
  { kind: 'out', text: 'wrote AGENTS.md and the Tifa skill in .claude/skills/tifa for coding agents' },
  { kind: 'out', text: '' },
  { kind: 'out', text: 'Next:' },
  { kind: 'out', text: '  cd myapp' },
  { kind: 'out', text: '  go mod tidy' },
  { kind: 'out', text: '  pnpm install' },
  { kind: 'out', text: '  go tool tifa dev     # http://127.0.0.1:9527' },
  { kind: 'cmd', text: 'cd myapp' },
  { kind: 'cmd', text: 'go mod tidy' },
  { kind: 'cmd', text: 'pnpm install' },
  { kind: 'dim', text: '…' },
  { kind: 'cmd', text: 'go tool tifa dev' },
  { kind: 'ok', text: 'tifa dev: serving http://127.0.0.1:9527' },
]

// The deploy of tifa-web itself, commit b33189e, from the log of GitHub Actions on 2026-11-21. The
// server's name is masked there. The line that starts uploading is written out from its format in
// cmd/tifa/deploy.go: the release is that of the line below it, and the size is what the log says.
export const deployTerminal: TerminalLine[] = [
  { kind: 'cmd', text: 'git push' },
  { kind: 'out', text: 'deploying tifa-web to *** (linux/amd64, port 9527)' },
  { kind: 'out', text: 'built tifa-web (build cd2404c7cca3d607)' },
  { kind: 'out', text: 'uploading release 20261121T175155Z (40.2 MB, gzipped on the way)' },
  { kind: 'out', text: 'starting 20261121T175155Z' },
  { kind: 'out', text: 'health check: GET /healthz answered 404' },
  { kind: 'ok', text: 'deployed tifa-web release 20261121T175155Z (commit b33189ed3edc) in 25.7s' },
]

// Whole visits to the example shop, one core each, as docs/performance.md has them for
// 2026-11-18. Each row is the value of Tifa, TanStack Start (uncompressed) and Next.js. The width
// of a bar is the value over the largest of its row.
export const numbers = [
  { key: 'views', big: '118.8', unit: '', values: [118.8, 94.7, 28.3], text: ['118.8', '94.7', '28.3'] },
  { key: 'memory', big: '245', unit: 'MiB', values: [245, 451, 613], text: ['245 MiB', '451 MiB', '613 MiB'] },
  { key: 'cpu', big: '8.4', unit: 'ms', values: [8.4, 10.6, 35.3], text: ['8.4 ms', '10.6 ms', '35.3 ms'] },
] as const
