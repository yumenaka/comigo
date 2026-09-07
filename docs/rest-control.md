# REST 控制接口（本地工作记录，不提交）

参考 Mihomo Meta 分支的资源路由和统一认证方式：
https://github.com/MetaCubeX/mihomo/blob/Meta/hub/route/server.go

保留 Comigo 的 `/api` 前缀，受 BasePath 影响。旧控制路由已移除，调用方包含设置页、书架、阅读历史、远程 Comigo 客户端。

| 方法 | 路径（省略 `/api`） | 用途 |
| --- | --- | --- |
| GET | `/server` | 服务版本、系统信息、书籍数、在线用户/设备/连接数量 |
| GET | `/connections` | 活跃 SSE/WebSocket 连接与去重统计 |
| POST | `/restart` | 重启 HTTP 服务，返回 202；宿主进程保持运行 |
| GET | `/stores` | 书库 ID、URL、名称、书籍数量、远程标志；远程 exists 为 null |
| GET | `/stores/:id` | 指定书库信息及书籍索引 |
| POST | `/stores/:id/refresh` | 刷新已配置的指定书库 |
| POST | `/stores/refresh` | 刷新全部书库 |
| DELETE | `/stores/:id` | 移除书库配置与索引，保留原始文件 |
| GET | `/books`、`/books/:id`、`/books/:id/parent` | 书籍索引、详情、父书组 |
| DELETE | `/books/:id/cache` | 清理指定书籍缓存 |
| GET/POST/DELETE | `/bookmarks` | 书签查询、写入、删除 |
| GET/PATCH | `/configs` | 读取配置（凭据脱敏）/局部更新 |
| PUT | `/configs/:name` | 更新单个字符串、布尔值或整数配置，JSON `{ "value": ... }` |
| POST/DELETE | `/configs/:name/items` | 增删数组配置项，JSON `{ "value": "..." }`；StoreUrls 用于添加书库 |
| PATCH | `/configs/login` | 更新账号密码，保留当前密码校验 |
| PATCH | `/configs/tailscale` | 更新 Tailscale 配置 |
| GET | `/configs/status` | 配置文件位置状态 |
| PUT/DELETE | `/configs/files/:location` | 保存/删除指定位置的配置文件 |
| PUT/DELETE | `/plugins/:name` | 启用/禁用插件 |

书库 ID 从 `/stores` 结果读取，为配置 URL 的 base64url 编码。解析后必须匹配已配置书库，相对与绝对本地路径按同一规范化规则匹配。未知资源返回 404，非法参数返回 400，未认证返回 401。

设置登录密码后，以上接口均要求 Cookie 或 `Authorization: Bearer <token>`。`POST /api/login` 接受浏览器表单，也接受 JSON `{ "username": "...", "password": "..." }`；JSON 登录响应包含 token、token_type、expires_in。修改密码后旧 token 立即失效。健康检查 `/healthz` 保持公开。

在线用户按登录账号去重，设备按浏览器 Cookie 去重，仅统计当前 SSE/WebSocket 连接；访客只计设备数，同一浏览器的多个标签页只有一个设备。Comigo 当前是单管理员账号模型，不能将设备数解释为自然人数。

验证：全量 Go 测试通过；Wails 标签下 routers、data_api、settings 测试通过。真实 Chromium 已验证单库刷新隔离、书库与书籍查询、在线统计、多标签去重、登录配置保存、表单/JSON 登录、Bearer 认证、错误密码/无效 token/匿名访问 401，以及服务重启后的恢复。新增书库通过设置页操作，删除书库通过浏览器调用接口验证。

手册修改保存在 `manual-before-rest-control-refactor` stash，未恢复。源码变化之外的 `_templ.go` 与 `assets/dist/main.js` 由 templ 与 bun 重新生成。原有 Makefile 与 docs 构建文件修改未主动改动；没有创建 commit。

最终生成产物检查：禁用 Parcel 缓存重建，确认只新增本轮 10 个 locale 键、书签接口路径及现有 package.json 锁定的 Alpine 3.16.3 版本字段（旧产物为 3.16.2），无手册逻辑混入。BasePath `/comics` 的页面、认证和重启也已通过浏览器验证。

最后一轮安全回归：浏览器验证改密返回 200、旧 Cookie 返回 401、新凭据登录返回 200、旧实时连接从连接集合移除；嵌入式无广播场景另有可运行回归测试。
