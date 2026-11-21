// The words of the home page and the links of the site, in one place. Every sentence is from
// the README of the Tifa repository or from docs/tutorial.md. [text](url) makes a link and
// `text` makes code, see inline.tsx.

export const site = {
  name: 'Tifa',
  url: 'https://tifago.com',
  tagline: 'The full-stack Go framework for Gophers.',
  intro: [
    'Write the backend in Go and the pages in React. Tifa compiles the whole site into one executable that you copy to a small VPS and run. No Node in production, no process manager, no reverse proxy, no Redis.',
    'It is for Go developers whose business logic and data already live in Go, who want to write pages with the React ecosystem, and who do not want to keep a Node service running in production just for server-side rendering.',
  ],
  vps: {
    header: ['On a 1 CPU / 1 GB VPS', 'The usual stack', 'With Tifa'],
    rows: [
      [
        'Reverse proxy and certificates',
        'nginx + certbot',
        'The Go HTTP server in the binary, which gets its own Let\'s Encrypt certificate and speaks HTTP/2',
      ],
      ['Runtime and process manager', 'Node.js + PM2', 'One Go executable'],
      [
        'React SSR',
        'TanStack Start',
        'TanStack Router, rendered by a [pure Go QuickJS](https://github.com/liustack/slowjs) inside the Go process',
      ],
      [
        'Cache',
        'Redis',
        'An in-process render memo for pages, a CDN in front if you want one, plus whatever you keep in Go memory',
      ],
      ['Queues and scheduled jobs', 'BullMQ + cron', 'Goroutines you start in `OnStart`'],
      ['Database', 'SQLite', 'SQLite, with a [pure Go driver](https://github.com/liustack/sqlite) so CGO stays off'],
    ],
  },
  quickstart: {
    title: 'Start an app',
    note: 'You need Go 1.27, Node 20.19 or 22.12 and later, and [pnpm](https://pnpm.io). Node only builds the front end. The executable you ship never needs it.',
    commands: `go run github.com/liustack/tifa/cmd/tifa@latest new tutorial-app
cd tutorial-app
go mod tidy
pnpm install
go tool tifa dev     # http://127.0.0.1:9527`,
  },
  links: {
    docs: '/docs/start/tutorial', // the first doc, which /docs leads to. TestFirstDocIsTheTutorial in cmd/syncdocs keeps it right
    github: 'https://github.com/liustack/tifa',
    x: 'https://x.com/liustack',
  },
}
