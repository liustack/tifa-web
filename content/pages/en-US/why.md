---
title: Why it exists
description: I wanted to ship a small product on TanStack Start, and the usual stack crawled on a 1 CPU, 1 GB VPS.
---

I wanted to ship a small product on TanStack Start. My other VPSes were full, Black Friday was still weeks away and I did not feel like buying another one, so I dug out an old box with 1 CPU and 1 GB of RAM. I had Codex set up the usual stack in one go: nginx, Node.js, PM2, TanStack Start, Redis, BullMQ, cron and SQLite.

It did not really run. It crawled.

So I started hacking on Tifa. Once it is done, deploying should mean copying one Go binary to the VPS, and that is it.

| On a 1 CPU / 1 GB VPS | The usual stack | With Tifa |
| :-- | :-- | :-- |
| Reverse proxy and certificates | nginx + certbot | The Go HTTP server in the binary, which gets its own Let's Encrypt certificate and speaks HTTP/2 |
| Runtime and process manager | Node.js + PM2 | One Go executable |
| React SSR | TanStack Start | TanStack Router, rendered inside the Go process by [slowjs](https://github.com/liustack/slowjs), the fastest pure Go JavaScript engine for rendering React |
| Cache | Redis | An in-process render memo for pages, a CDN in front if you want one, plus whatever you keep in Go memory |
| Queues and scheduled jobs | BullMQ + cron | Goroutines you start in `OnStart` |
| Database | SQLite | SQLite, with [the fastest pure Go SQLite driver](https://github.com/liustack/sqlite), so CGO stays off |

I still do not know what to call it. A React SSR framework with no Node runtime? A React SSR framework with a Go backend? Both sound a little unhinged. Under the hood it is Go, [slowjs](https://github.com/liustack/slowjs) (a pure Go JavaScript engine), and TanStack Router.

## Comparisons

- [Against Next.js and TanStack Start](/docs/reference/performance#against-nextjs-and-tanstack-start)
- [Next.js and TanStack Start, the study on GitHub](https://github.com/liustack/tifa/blob/v0.2.1/docs/studies/nextjs-and-tanstack-start.md)
