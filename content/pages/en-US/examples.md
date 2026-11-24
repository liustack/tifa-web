---
title: 'Shanhai Workbench: a Tifa example app'
description: Two sides of one small shop, in one program.
---

Two sides of one small shop, in one program.

- **The workbench**, for the shop owner: two shops, an order list of 25 rows a page with product photos, to sort, search and filter by status, date, carrier and amount, with rows to select and export, order details with a log of what happened, internal notes, payment and invoice, a stock check, a packing list, the way of the parcel, a shipping form and a dialog to change where an order goes, and a sales stats page with charts, a table of every product with sparklines, a heatmap of when orders come, carriers, dispatch times and the best customers. Each shop has its products, with a page to change the price, the stock and the title and description in both languages, and the first shop sells those of the storefront. After you ship an order the list and the stats update with it.
- **The storefront**, for customers, in English and Chinese: a category menu and a cart drawer on every page, a home page with a carousel and tabs of best sellers, a product list with search, facets with counts, quick filters, chosen conditions, a grid and list view and a quick view of each product, product pages with a photo gallery, the whole specification, a comparison with similar products, reviews in pages with photos and replies from the shop, questions and products to buy together, coupons and suggestions around a cart kept in the program's memory with its ID in a cookie, a checkout with a long form, and a journal of articles in Markdown with highlighted code. Prices are in yuan or dollars, and dates and numbers follow the language.

The two sides share their data. A checkout in the storefront becomes the next order of the first shop. Every order takes its pieces from the stock of its shop, a checkout and an order placed in the workbench alike, and an order of more than is left is refused. A change to a product in the workbench shows in the storefront on its next request.

Pages are rendered to complete HTML on the server. After the browser takes over, it moves between pages on its own. Data lives in process memory and resets on restart.

The interface looks like a printed catalogue: thin lines, numbered products and their details in plain sight. Every product has three photos, of itself, a close detail and in use, stored as AVIF at 160, 480 and 1200 pixels wide. A product in another colour shares the photos of the one it is made from. Latin headings use Newsreader, bundled with the app. Chinese headings and body text use the system fonts. The interface is built with Tailwind CSS and shadcn/ui, the order list with TanStack Table, the forms with react-hook-form and zod, the charts with recharts, and the numbers tween with react-spring. The storefront is translated with i18next, and the journal is rendered with react-markdown, with code highlighted by highlight.js.

It shows what a Tifa app looks like: how folders are grouped by business area, how Go prepares data for pages, and how React components read that data.

## Run it

Run in the `examples/shop` folder. You need Go 1.27, Node 20.19 or 22.12 and later, and pnpm. The build bundles the front end with Vite on Node, and the built `bin/shop` runs without Node:

```sh
pnpm install --frozen-lockfile
go tool tifa build
./bin/shop
```

Open <http://127.0.0.1:9527>. It sends you to the sign-in page, and after you sign in, to the order list of the first shop. The storefront is at <http://127.0.0.1:9527/store> and, in Chinese, <http://127.0.0.1:9527/zh-CN/store>. Use `./bin/shop --server-addr 127.0.0.1:9528` to change the port. `bin/shop` is a single executable with the frontend output embedded, so you can copy it to another machine and run it there.

Checks and tests, with one command:

```sh
go tool tifa check
```

## Sign in

The workbench is for the owner, so every page under `/projects` needs a sign-in. The storefront does not. The demo account is:

| E-mail | Password |
| --- | --- |
| `owner@shanhai.example` | `shanhai-demo` |

A visitor who is not signed in is sent to `/login?next=...` by `OnRequest` in `hooks.go`, and goes back to the page they asked for after signing in. The top bar and the sidebar show the name (the Layout data is the same for everyone, and the name comes from the session in the browser) and the top bar has a sign-out button or a link to sign in, with `SignedIn`, `SignedOut` and `useSession` from `tifa/client`. The app never calls `tifa.SetSessionStore`, so sessions live in the memory of the process and a restart signs everyone out. The password is checked against a bcrypt hash in `session.go`, where a real shop would read it from its database.

To look at the workbench without signing in, start the program with `SHOP_OPEN=1`. The comparison with Next.js and TanStack Start and the benchmarks of the pages use it, since they cannot sign in. It lets in the whole workbench, writes such as shipping an order included, so use it only for the comparison and the benchmarks and never set it on a server.

## Code on GitHub

- [examples/shop](https://github.com/liustack/tifa/tree/v0.5.1/examples/shop)
- [The storefront](https://github.com/liustack/tifa/tree/v0.5.1/examples/shop/app/store)
- [The workbench](https://github.com/liustack/tifa/tree/v0.5.1/examples/shop/app/projects)
- [README of the example](https://github.com/liustack/tifa/blob/v0.5.1/examples/shop/README.md)
