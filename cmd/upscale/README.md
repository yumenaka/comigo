# 独立图片辅助程序（实验性）

复用项目已有的 Echo、imaging；页面使用原生 HTML/CSS、fetch、EventSource，无新增 Go 或前端依赖，不需要构建 Comigo 前端。程序不读取 Comigo 配置或书库，也尚未接入阅读器的转换链路。

## 启动

在仓库根目录运行：

```sh
go run ./cmd/upscale
```

打开 <http://127.0.0.1:1235>。默认无认证，仅监听本机。其他机器调用时可指定 `-listen 0.0.0.0:1235`。

普通缩放无需安装额外程序。AI 和 WebP 编码分别需要：

- [Real-ESRGAN ncnn Vulkan](https://github.com/xinntao/Real-ESRGAN-ncnn-vulkan)：下载适合系统的可执行文件和模型，保留同目录的 `models/`。本轮使用 [官方 20220424 macOS 包](https://github.com/xinntao/Real-ESRGAN/releases/tag/v0.2.5.0) 验证。
- [cwebp](https://developers.google.com/speed/webp/docs/precompiled)：官方 WebP 编码器；macOS 可用 `brew install webp`，其他系统可下载官方工具包。

将工具放入 PATH，或显式指定路径：

```sh
go run ./cmd/upscale \
  -engine ./engines/realesrgan-ncnn-vulkan \
  -models ./engines/models \
  -cwebp ./engines/cwebp
```

Windows 使用对应的 `.exe` 路径。`-models` 省略时使用引擎旁的 `models` 目录。程序启动时检测能力，安装新模型后重启；页面只显示找到 `.bin` 和 `.param` 的模型，发现文件并不代表驱动一定可用，推理失败会显示引擎诊断。

```sh
go build -o ./upscale ./cmd/upscale
./upscale -help
```

Go 辅助程序不依赖 CGO。AI 在 Windows/Linux 使用 Vulkan，在 macOS 由 ncnn 包的 MoltenVK 支持 Apple GPU；不要求 CUDA，也不限 NVIDIA。仍需在目标 AMD/Intel/NVIDIA 设备验证驱动与引擎兼容性。移动端当前作为 HTTP 客户端使用，不提供移动端本地引擎打包。

## 可用参数

| 参数 | 默认值 | 说明 |
|---|---|---|
| `algorithm` | `lanczos` | `ai` / `lanczos` / `cubic`（Catmull-Rom）/ `nearest` / `none`（仅转换格式） |
| `model` | 空 | AI 必填；名称及支持倍率从 capabilities 获取 |
| `scale` | `2` | 普通缩放 0.1–4，输出尺寸四舍五入且最小 1；`none` 必须为 1 |
| `format` | `png` | `png` / `jpeg` / `webp` |
| `quality` | `85` | WebP 0–100、JPEG 1–100；无损 WebP 时表示压缩投入；PNG 忽略 |
| `lossless` | `false` | WebP 无损；配合 `-exact` 保留完全透明像素的 RGB |
| `method` | `4` | WebP 编码投入 0–6，越高通常越慢 |
| `tile` | `0` | AI 分块：0 自动，或 32–4096；显存不足先尝试 128 / 64 |
| `gpu` | `-1` | 自动选择；非负值选择 ncnn 设备编号，支持 0–31；-1 **不是** CPU 模式 |
| `tta` | `false` | AI 测试时增强，明显增加耗时 |

| AI 模型 | 倍率 | 用途 |
|---|---|---|
| `realesr-animevideov3` | 2、3、4 | 速度优先的动画/漫画初步实验 |
| `realesrgan-x4plus-anime` | 4 | 动漫图像 |
| `realesrgan-x4plus` | 4 | 通用图像 |

AI 模型不能混用文件命名：animevideov3 需要 `realesr-animevideov3-x2/3/4.bin` 和对应 `.param`；另外两款使用模型名加 `.bin/.param`。

## HTTP API（供 Comigo 调用）

所有路径相对于辅助程序的地址；服务端调用无需 CORS。请求不接受本地文件路径、远程 URL、任意模型路径或命令参数。

| 方法与路径 | 行为 |
|---|---|
| `GET /api/capabilities` | 已安装模型、倍率和 WebP 可用状态 |
| `POST /api/jobs` | multipart：`file` 图片，`options` JSON；返回 202、任务状态及 Location |
| `GET /api/jobs/:id` | 查询状态 |
| `GET /api/jobs/:id/events` | SSE，每 250ms 发送 `data: {状态 JSON}`；终态后关闭 |
| `GET /api/jobs/:id/result` | 成功时返回图片字节；尚未成功返回 409 |
| `POST /api/jobs/:id/cancel` | 取消排队或运行中的任务，返回 202 |
| `DELETE /api/jobs/:id` | 删除已结束任务及临时文件，返回 204；运行中返回 409 |

先提交，例如 AI 2× 转有损 WebP：

```sh
curl -sS http://127.0.0.1:1235/api/jobs \
  -F 'file=@input.png' \
  -F 'options={"algorithm":"ai","model":"realesr-animevideov3","scale":2,"format":"webp","quality":85,"tile":128}'
```

仅转无损 WebP，将 `options` 改为：

```json
{"algorithm":"none","scale":1,"format":"webp","lossless":true,"quality":100,"method":6}
```

用响应中的 `id` 替换下面的 `JOB_ID`，监听并取结果：

```sh
curl -N http://127.0.0.1:1235/api/jobs/JOB_ID/events
curl -f http://127.0.0.1:1235/api/jobs/JOB_ID/result -o result.webp
curl -X DELETE http://127.0.0.1:1235/api/jobs/JOB_ID
```

状态 `state` 为 `queued`、`running`、`done`、`failed`、`canceled`；`stage` 是可展示文字，`progress` 为 0–100，`elapsed_seconds` 含排队时间，`format` 是输出格式。成功时带 `result_url`、`width`、`height`、`bytes`；失败或取消带 `error`。错误响应使用 Echo 的 `{"message":"..."}` 格式；非法请求返回 400，任务不存在 404，队列/保留量满返回 429。

Comigo 可用标准库 `net/http` 和 `mime/multipart` 上传已有图片字节，再轮询任务或消费 SSE；成功后读取 `result_url` 并按既有图片响应/缓存流程返回。客户端请求中断不会取消任务，需要显式调用 cancel，并在终态后 DELETE。当前没有修改 Comigo 主程序。

## 边界与验证

- 单 worker 串行转换，最多保留 16 个任务（含已完成任务）。API 调用者负责 DELETE；网页再次转换会先删除上一次结果，请先下载需要保留的结果。
- 每个任务从提交起最长 30 分钟。取消会终止 ncnn/cwebp；Go 内部解码、普通缩放、PNG/JPEG 编码只能在步骤间响应取消。
- 上传请求上限 32 MiB；输入最多 1600 万、输出最多 6400 万像素；输出单边最多 32768，WebP 最多 16383。高分辨率仍可能占用较多内存，需要在目标机器实测。
- 接受 PNG/JPEG/WebP 静态图。没有动画保帧功能；不要用于动画素材。EXIF 朝向归一化，元数据剥离；没有 ICC 色彩管理，图片颜色受解码器处理影响。PNG/WebP 保留透明度，JPEG 使用白色背景。
- AI 百分比读取 ncnn 分块日志，WebP 读取 cwebp 日志；总进度按阶段映射，不是精确剩余时间。PNG/JPEG 与普通缩放只有阶段进度，不伪造逐像素百分比。
- 文件放在系统临时目录，Ctrl-C 或 SIGTERM 正常退出时清理；任务不持久化，强制结束进程可能留下临时目录。
- 未自动下载模型，未打包引擎，未做缓存和批量上传。模型质量、显存和 5060 Ti 速度仍按调研计划继续实测；本机结果不能当作 5060 Ti 基准。

运行回归检查：

```sh
go test -race ./cmd/upscale
UPSCALE_TEST_ENGINE=./engines/realesrgan-ncnn-vulkan go test ./cmd/upscale -run TestRealAI -v
```

可另设 `UPSCALE_TEST_MODELS` 指定测试模型目录。未安装 cwebp 时 WebP 测试会明确跳过；默认不运行 GPU 测试。
