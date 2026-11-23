---
title: 为什么有 Tifa
description: 我想让 TanStack Start 在 1 核 1G 的 VPS 上跑起来，常见的那一套却卡成狗。
---

说真的，我都不知道怎么定义这个东西。叫它没有 Node.js 运行时的 React SSR 框架，还是用 Go 写后端的 React SSR 框架？听起来都挺神经的。它是一个基于 Go、[slowjs](https://github.com/liustack/slowjs)（纯 Go 写的 JavaScript 引擎）和 TanStack Router 实现的全栈框架。

我做它的初衷，是想让 TanStack Start 在 1 核 1G 的 VPS 上跑起来。我想做个小产品，但手上的 VPS 都跑满了，黑五还没到，我也不想去买新的，就找了台旧的 1 核 1G。我让 Codex 一把梭，部署了 nginx、Node.js、PM2、TanStack Start、Redis、BullMQ、cron 和 SQLite。

然后发现这些东西根本跑不起来，准确说是卡成狗。

所以我开始糊 Tifa 这个轮子。理论上糊完以后，我只需要在 VPS 上部署一个 Go 二进制文件就完了。

| 在 1 核 1G 的 VPS 上 | 常见的一套 | 用 Tifa |
| :-- | :-- | :-- |
| 反向代理与证书 | nginx + certbot | 程序里的 Go HTTP 服务器，自己申请 Let's Encrypt 证书，支持 HTTP/2 |
| 运行时与进程管理 | Node.js + PM2 | 一个 Go 可执行文件 |
| React SSR | TanStack Start | TanStack Router，由 Go 进程内的 [slowjs](https://github.com/liustack/slowjs) 渲染，它是渲染 React 最快的纯 Go JavaScript 引擎 |
| 缓存 | Redis | 进程内的渲染记忆，想要的话前面放 CDN，以及你放在 Go 内存里的数据 |
| 队列与定时任务 | BullMQ + cron | 在 `OnStart` 里启动的 goroutine |
| 数据库 | SQLite | SQLite，用[全网最快的纯 Go SQLite 驱动](https://github.com/liustack/sqlite)，CGO 保持关闭 |

## 对比

- [对比 Next.js 和 TanStack Start](/docs/reference/performance#against-nextjs-and-tanstack-start)
- [Next.js 和 TanStack Start 的对比研究（GitHub，英文）](https://github.com/liustack/tifa/blob/v0.2.0/docs/studies/nextjs-and-tanstack-start.md)
