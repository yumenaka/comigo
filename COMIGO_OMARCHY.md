# comigo-omarchy

独立 Git 仓库，插件 ID `yumenaka.comigo`；父仓库忽略此目录。仅在用户要求时提交、推送。

## 实现

- 概览页本机启用对外服务且存在多个 IP 时，二维码两侧箭头循环切换阅读链接与二维码地址，并加粗当前 IP；默认沿用服务返回的出口地址，远程模式保持配置 URL。
- `Panel.qml` 管理概览、状态、服务、设置四页；概览提供二维码及下方阅读入口和 IP，状态页提供服务状态、速度、累计流量与书籍/连接统计，默认本机模式及概览页，本机缺少 CLI 时进入服务页；远程未配置地址时进入设置页。
- 本机与远程通过侧栏按钮切换，Shell 启动默认本机；`serverURL` 保存本机回环地址，`remoteURL` 保存远程完整地址。远程各页只显示对应 REST 数据，隐藏本地启停、日志、CLI/书库设置、对外服务与防火墙。
- `Service.qml` 共享 REST、登录会话、轮询和设置；`bin/comigo-ctl` 使用 Bash、curl 调用 CLI 及管理进程。
- JSON、剪贴板、浏览器使用 QML/Quickshell；平台文件工具与 pkexec 使用宿主提供的版本。
- 服务页本机控制块提供 `autoStart` 开关，默认关闭；插件加载后直接调用本地 CLI，失败间隔 10 秒、最多 3 次。运行目录保存次数与完成标记，重载不重置；成功或手动控制后不再自动拉起，关闭选项取消后续重试，不使用 systemd。
- CLI 进程以 PID 和内核启动时间识别，启停加锁；端口被占用时拒绝启动。日志在用户 state 目录。
- 未指定书库时，插件不填充目录、不传书库参数，由 Comigo 自身的默认规则处理；指定书库时校验目录存在。
- 插件不提供 Comigo 二进制下载或安装功能；服务页在本机和远程模式均显示 GitHub 项目地址与 comigo.xyz 官网地址（中国大陆推荐），支持浏览器打开和复制链接。pkexec 仅用于 UFW 规则变更。
- 关键函数写中文注释，文本同步 `I18n.js` 的中英日翻译。

## 接口与数据

- `/api/server` 提供地址、流量、`externalAccess` 和 `listenAddress`；完整状态每 30 秒及手动刷新，打开面板时每 2 秒刷新流量；设置页同期读取完整状态和运行配置。后台查询不禁用控件，相同快照复用对象。
- 防火墙使用宿主 UFW、ip；只在本机回环端点下提供操作，放行默认网卡的直连 RFC1918 网段和当前 TCP 端口。规则标记 `omarchy-comigo`，撤销仅删除本插件规则；等价现有规则保持归属。
- 对外服务通过 `PATCH /api/configs` 设置 `DisableLAN`，遵守认证、只读模式及 BasePath；更新后等待重连。配置由 Comigo 解析、保存并重启监听。
- 插件不提供更新检查、发布页或升级命令；保留当前版本展示及最低支持版本校验。
- 配置文件只读展示 `/api/configs/status` 的 `current`：实际路径、位置、运行类型、格式和存在状态。修改、保存、删除由网页管理；插件启动不传配置文件路径。
- `/api/configs` 返回脱敏配置，`POST /api/login` 登录，`/api/qrcode.png` 生成阅读二维码。
- 登录会话按模式及地址隔离，仅存内存；切换恢复对应会话并清空状态重查，退出或切换使在途旧响应失效。远程阅读、浏览器及二维码使用配置的完整远程地址。
- 密码提交后清空，token 只存内存；凭据通过 stdin/fd 传给 curl。设置仅保存非敏感 JSON，禁止 source/eval。
- 设置文件监听外部变更；未编辑字段自动同步，草稿保留，保存只合并已编辑字段。
- 命令使用参数数组。修改宿主安装副本只操作脚本拥有的目录。

## 验证

以下命令在 `comigo-omarchy/` 目录执行。

```bash
omarchy plugin validate .
bash -n install.sh bin/comigo-ctl tests/*.sh
bash tests/test-ctl.sh
bash tests/test-refresh.sh
bash tests/test-modes.sh
bash tests/test-reading-ip.sh
bash tests/test-autostart.sh
bash tests/test-download-links.sh
bash tests/test-version.sh
bash tests/test-firewall.sh
COMIGO_TEST_CLI=/path/to/comi bash tests/test-default-library.sh
COMIGO_TEST_CLI=/path/to/comi bash tests/smoke.sh
```

联调使用临时书库、配置和端口，退出时停止测试进程。检查启停、登录、二维码、剪贴板、流量、监听切换、默认书库与页面、语言；核心修改执行相关 Go 测试。

保持文档描述当前实现。不得提交凭据、私人配置或截图。安装到本机时使用 Omarchy 与系统档案技能，同步当前状态。`install.sh` 复制后重启 Shell，加载安装目录的 QML 组件；安装验证包含文件一致性、启用状态及实际面板。


最低支持 Comigo v1.3.6，桌面协议版本为 `1`。
