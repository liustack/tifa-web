---
title: Ctx
description: The methods of *tifa.Ctx
source: docs/reference/ctx.md
group: Reference
groupOrder: 5
order: 5
---


`*tifa.Ctx` is the request and its response. `Load`, `Layout`, the method handlers, `OnRequest` and `OnError` get it.

| Method | What it does |
| --- | --- |
| `c.Param(name)` | A path parameter, decoded, such as `c.Param("order")` for the folder `order+`. In a `Layout`, `""` for a parameter of a folder below it |
| `c.Query(name)` | The first value of a query parameter, `""` when absent. For method handlers, `Validate` and `OnRequest`. `Load` and `Layout` take their query parameters as a second parameter instead, see [Query parameters](/docs/pages-and-data/loading-data#query-parameters), and calling this there fails the build |
| `c.QueryValues(name)` | Every value of a query parameter, such as `["a", "b"]` for `?tag=a&tag=b`, `nil` when absent. Not allowed in `Load` and `Layout`, as `c.Query` |
| `c.Header(name)`, `c.Cookie(name)` | A request header, a request cookie (`http.ErrNoCookie` when absent) |
| `c.Context()`, `c.WithContext(ctx)` | The request context, and a copy of `c` that carries another one, to store request-scoped values. `c` itself does not change |
| `c.Pattern()` | The path of the matched route, such as `/projects/{project}`, `""` when no route matched |
| `c.Path()` | The path of the page as the visitor sees it in the address bar: escaped, with the language prefix, without the query, such as `/en-US/store/items/a%20b`. The same for a page load and for the data request of a navigation. In a method handler it is the path of the handler. Not allowed in a `Layout`, see [How Load may use Ctx](/docs/pages-and-data/loading-data#how-load-may-use-ctx) |
| `c.Language()` | The language of the route being answered: the language in the URL, the default language when the URL has none, `""` on a site without languages. See [The language in code](/docs/guides/i18n#the-language-in-code) |
| `c.SetAlternates(list)` | Replaces the versions of the page in other languages, for `hreflang` and the sitemap. Call it in a `Load`. See [What search engines get](/docs/guides/i18n#what-search-engines-get) |
| `c.IsDataRequest()` | Whether this is a navigation data request of the browser router rather than a page load |
| `c.SetHeader(key, value)`, `c.SetCookie(cookie)` | A response header, a response cookie. Allowed in `Load` too |
| `c.JSON(code, v)`, `c.Text(code, s)`, `c.NoContent(code)` | Writes the response. Not allowed in `Load` |
| `c.Bind(&v)`, `c.BindQuery(&v)` | Decodes the body or the query string into a struct, see [Method handlers](/docs/pages-and-data/handlers). Not allowed in `Load` and `Layout` |
| `c.SetBodyLimit(n)` | Changes the body limit of this request (`server.body_limit` by default, see [Body limit and read time of a route](/docs/pages-and-data/handlers#body-limit-and-read-time-of-a-route)). Call it before anything reads the body, usually in `OnRequest`. Once reading has started it returns an error |
| `c.SetReadTimeout(d)` | Gives the body of this request `d` to arrive, in place of `server.read_timeout`, and the answer `server.write_timeout` after it. Call it before anything reads the body, usually in `OnRequest`. Once reading has started it returns an error |
| `c.Request()`, `c.Writer()` | The raw `*http.Request` and `http.ResponseWriter`. In `Load` and `Layout` the request is a copy with only the declared query parameters and no body, and writing to the response breaks the page, see [The raw request in Load](/docs/pages-and-data/loading-data#the-raw-request-in-load) |
