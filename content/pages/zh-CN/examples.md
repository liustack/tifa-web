---
title: '山海工作台：Tifa 示例应用'
description: 一家小店的两面，放在同一个程序里。
---

一家小店的两面，放在同一个程序里。

- **工作台**，给店主用：两家店铺，每页 25 行、带商品图的订单列表，可以排序、搜索，按状态、日期、快递和金额筛选，也可以勾选订单批量导出。订单详情有操作记录、内部备注、付款与发票、库存检查、拣货单和物流轨迹，还有发货表单和修改收货信息的弹窗。销售统计有图表、带走势线的商品销售明细、下单时段热力图、快递表现、发货时效和常客。每家店铺都有自己的商品，可以改价格、库存和中英文的标题与介绍，第一家店铺卖的就是店面的商品。发货后列表和统计跟着更新。
- **店面**，给顾客用，有英文和中文：每页都有分类菜单和购物车抽屉，首页有轮播和按分类切换的热销，商品列表有搜索、带数量的筛选、常用筛选、已选条件、网格和列表两种视图和每件商品的快速查看，商品详情有图库、完整规格、和同类商品的对比、分页的评价（带图，店家会回复）、问答和搭配购买，购物车存在程序内存里、cookie 只放编号，周围有优惠券和推荐，结算页表单很长，还有用 Markdown 写、代码带高亮的手记文章。价格按语言用人民币或美元，日期和数字也按语言格式化。

两面共用一份数据。店面结账后，订单就成了第一家店铺的下一单。每张订单都从所在店铺的库存里扣掉相应件数，店面结账和工作台建单都一样，买的比剩下的多就拒绝下单。在工作台改了商品，店面的下一次请求就能看到。

页面在服务器上渲染成完整 HTML，浏览器接管后在页面之间跳转。数据放在进程内存里，重启恢复初始值。

界面照着印刷的商品目录来做：细线，商品编号，规格摆在明处。每件商品有三张图：商品本身、细节和使用场景，存成 AVIF，各有 160、480 和 1200 像素宽三份。换了颜色的款和原款共用照片。英文标题用 Newsreader，随应用打包，中文标题和正文用系统字体。界面用 Tailwind CSS 和 shadcn/ui，订单列表用 TanStack Table，表单用 react-hook-form 和 zod，图表用 recharts，数字补间用 react-spring。店面的翻译用 i18next，手记用 react-markdown 渲染，代码高亮用 highlight.js。

它用来演示一个 Tifa 应用长什么样：怎么按业务分目录，Go 怎么给页面准备数据，React 组件怎么读这些数据。

## 运行

在 `examples/shop` 目录执行，需要 Go 1.27、Node 20.19 或 22.12 及以上，以及 pnpm。构建时用 Vite 在 Node 上打包前端，构建出的 `bin/shop` 运行时不需要 Node：

```sh
pnpm install --frozen-lockfile
go tool tifa build
./bin/shop
```

打开 <http://127.0.0.1:9527>，会先到登录页，登录后跳到第一家店铺的订单列表。店面在 <http://127.0.0.1:9527/store>，中文版在 <http://127.0.0.1:9527/zh-CN/store>。`./bin/shop --server-addr 127.0.0.1:9528` 换端口。`bin/shop` 是单个可执行文件，前端产物已经嵌在里面，拷到别的机器上直接运行。

检查与测试，一条命令：

```sh
go tool tifa check
```

## 登录

工作台是给店主用的，`/projects` 下的每个页面都要登录，店面不用。演示账号：

| 邮箱 | 密码 |
| --- | --- |
| `owner@shanhai.example` | `shanhai-demo` |

没登录的访客会被 `hooks.go` 里的 `OnRequest` 带到 `/login?next=...`，登录后回到原来要看的页面。顶栏和侧栏显示名字（Layout 的数据对所有人一样，名字来自浏览器里的会话），顶栏有退出按钮或登录链接，都用 `tifa/client` 的 `SignedIn`、`SignedOut` 和 `useSession`。应用没有调用 `tifa.SetSessionStore`，会话放在进程内存里，重启后所有人要重新登录。密码和 `session.go` 里的 bcrypt 哈希比对，真实的店铺会从数据库里读它。

想不登录就看工作台，用 `SHOP_OPEN=1` 启动程序。和 Next.js、TanStack Start 的对比以及页面压测不能登录，用的就是它。它放行整个工作台，包括发货这类写操作，只给对比和压测用，不要在服务器上设置。

## GitHub 上的代码

- [examples/shop](https://github.com/liustack/tifa/tree/v0.1.13/examples/shop)
- [店面](https://github.com/liustack/tifa/tree/v0.1.13/examples/shop/app/store)
- [工作台](https://github.com/liustack/tifa/tree/v0.1.13/examples/shop/app/projects)
- [示例的 README](https://github.com/liustack/tifa/blob/v0.1.13/examples/shop/README.zh-CN.md)
