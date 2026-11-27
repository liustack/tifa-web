---
title: Debugging
description: Debug an app under tifa dev
source: docs/build-and-deploy/debugging.md
group: Build and deploy
groupOrder: 4
order: 6
---


Three places run your code under [`tifa dev`](/docs/build-and-deploy/dev), and each has its own way in: the Go app, the React components in the browser, and the React components in the server render.

## Go: Delve

`tifa dev` runs the app as a process of its own and keeps its debug info, so [Delve](https://github.com/go-delve/delve) attaches to it. `tifa dev status` names the process ID as `app pid`, and `--json` as `appPid`. It is the app, not `tifa dev`, which is the `pid` above it:

```sh
go tool tifa dev --background
go tool tifa dev status --json   # "appPid": 41290
go run github.com/go-delve/delve/cmd/dlv@latest attach 41290
```

(Install it once with `go install github.com/go-delve/delve/cmd/dlv@latest` and run `dlv attach 41290`.) A breakpoint is a package path and a function, or a file and a line. Open the page, and the request stops there:

```text
(dlv) break demo/app/todos.ListTodos
(dlv) continue
> [Breakpoint 1] demo/app/todos.ListTodos() ./app/todos/todo.go:25
(dlv) bt 3
(dlv) print ctx
(dlv) continue
```

- The page waits while the app is stopped, and the browser shows nothing until you `continue`.
- `go build` optimizes, so Delve says `debugging optimized function` and may not show a variable. To build without optimizations, start `tifa dev` with the flags in `GOFLAGS`. The quotes keep the space in one flag:

  ```sh
  GOFLAGS="'-gcflags=all=-N -l'" go tool tifa dev --background
  ```

- A change to a Go file builds the app again and starts a new process, so the debugger loses it. Read `appPid` again and attach again. A change to a TSX file or a stylesheet leaves the process alone: the app loads the new bundle in place.
- `exit` in Delve asks whether to kill the process. Answer `n` to leave the app running.
- A program built with `tifa build` has no debug info unless `build.debug_info` is set. See [Debug info](/docs/build-and-deploy/build#debug-info).

## React in the browser

The browser gets the modules of your pages from Vite's dev server, with source maps, so the Sources panel of the browser shows your `.tsx` files and a breakpoint in a component stops where you set it. React Developer Tools shows the components and their props. `tifa dev` shows errors of the browser in an overlay on the page, and `tifa errors` prints them (see [Errors in the browser](/docs/build-and-deploy/dev-errors#errors-in-the-browser)).

## React on the server

The server render runs in a JavaScript VM with no debugger, so you cannot set a breakpoint in a component while it renders there. Under `tifa dev`, what a component prints with `console.log` goes to the log of the app with the page and the line of the call, and `tifa dev logs` prints it:

```text
2026/11/26 08:00:00 INFO console.log: rendering about {"visits":3} path=/about at=app/about/page.tsx:3:11
2026/11/26 08:00:00 INFO GET /about 200 28ms /about memo miss
```

The line below it is the request line of `tifa dev`: the route and whether the render memo had the page. See [Console in `tifa dev`](/docs/pages-and-data/server-render#console-in-tifa-dev) for what is printed and how, and [Log](/docs/build-and-deploy/dev-errors#log) for the request lines.
