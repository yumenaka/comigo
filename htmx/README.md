## Starting your project

> ❗️ 下面是开发必须工具和库:
>
> - `Templ（模板引擎）`: [https://templ.guide](https://templ.guide/)
> - `Flowbite（组件库）`: [https://flowbite.com/docs/components/gallery/](https://flowbite.com/docs/components/gallery/)
> - `Gin`: [https://gin-gonic.com/zh-cn/docs/examples/](https://gin-gonic.com/zh-cn/docs/examples/)
> - `htmx`: [https://htmx.org/examples](https://htmx.org/examples/)
> - `Alpine.js`: [https://alpinejs.dev](https://alpinejs.dev/)
> - `Tailwind CSS`: [https://tailwindcss.com/docs/flex](https://tailwindcss.com/docs/flex)
> - `Air`: [https://github.com/air-verse/air](https://github.com/air-verse/air)
> - `Bun`: [https://github.com/oven-sh/bun](https://github.com/oven-sh/bun)
> - `i18next`: [https://www.i18next.com](https://www.i18next.com/)
> - `golangci-lint`: [https://github.com/golangci/golangci-lint](https://github.com/golangci/golangci-lint)
> - `Prettier`: [https://prettier.io/docs/en/index.html](https://prettier.io/docs/en/index.html)
> - `Icon`: [https://www.xicons.org/#/](https://www.xicons.org/#/)
> - `htmx-go`: [https://github.com/angelofallars/htmx-go](https://github.com/angelofallars/htmx-go)

## TODO
 - [x] 服务器设置页面 v1.0
 - [ ]  [self update 功能](https://github.com/minio/selfupdate)
 - [ ] 显示服务器log：[web终端示例](https://zenn.dev/ikedam/articles/2e078bfc2a4cb6)，设置页面功能
 - [ ] 优化翻页模式 
 - [ ] 自定义js与css代码块功能。
 - [ ] [使用 Go1.24 的 os.Root 类型](https://antonz.org/go-1-24/)，将[文件操作限制在特定目录](https://go.dev/blog/osroot)，以防止攻击者通过转义或相对路径非法访问文件
 - [x] go 1.24正式发布，[使用 Go1.24 的 go get -tool](https://antonz.org/go-1-24/)，添加工具依赖项
 - [x] 等 [gowebly](https://github.com/gowebly/gowebly) 更新后，参考 gowebly 的模板，升级tailwindcss到4.0
 - [ ] 自动发版功能 [goreleaser](https://goreleaser.com/)  [github-action](https://dev.to/hadlow/how-to-release-to-homebrew-with-goreleaser-github-actions-and-semantic-release-2gbb)
 - [ ] comigo后台：有几台设备在线，阅读文件，阅读页数，当前用户状态、阅读书籍、阅读进度、阅读时间、服务器状态 注册，阅读记录，
 - [ ] 可能的BUG：当前文件夹顶层无书籍的时候，首页为空或跳过规则不对？上一版有这个问题，目前有无问题还不确定。
 - [x] gin -> chi或gin -> echo。chi无外部依赖,echo似乎与wails兼容性更好。此外还可以试试不用框架全原生。
 - [ ] 尝试wails3, https://v3alpha.wails.io/getting-started/installation/
 - [ ] 合并htmx代码，参考[pagoda](https://github.com/mikestefanello/pagoda)，重新规划项目结构。我用的许多组件，最终都换成和这个模板一样的了，估计从这个项目里可以学到很多东西。[go-blueprint](https://docs.go-blueprint.dev/) 也是一个不错的参考，可以看看怎么集成websockets与templ。
 - [ ] 添加[数据验证](https://dev.to/leapcell/validator-complex-structs-arrays-and-maps-validation-for-go-34ni)。
 - [ ] 优化打开浏览器与扫描逻辑，减少等待时间。可以使用[端口检测包](https://github.com/wait4x/wait4x)。
 - [x] check：无参数的逻辑。epub页数解析似乎出现了问题。
 - [ ] SteamDeck支持（网页支持手柄操作）鼠标滚轮对应
 - [ ] 同步翻页 -> 全局多端同步跟踪页面状态，除了不同id的书籍，其他页面状态都可以同步。
## 提示

<https://github.com/angelofallars/htmx-go#triggers>  
Alpine.js 可以监听并接收由 htmx-go 触发的事件详细信息，这使得服务器端触发的事件在事件驱动的应用程序中非常方便！

对于 Alpine.js，你可以注册一个 x-on:<EventName>.window 监听器。.window 修饰符很重要，因为 HTMX 会从根窗口对象调度事件。要接收由 htmx.TriggerDetail 和 htmx.TriggerObject 发送的值，你可以使用 $event.detail.value。

要开始您的项目，在终端中运行 **Gowebly** CLI命令：

```console
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/gowebly/gowebly/v2@latest
go install github.com/air-verse/air@latest
curl -fsSL https://bun.sh/install | bash
gowebly run
```

或者：

```console
air
```

## Project overview

Backend:

- Module name in the go.mod file: `github.com/yumenaka/comigo/htmx`
- Go web framework/router: `Gin`
- Server port: `1234`

Frontend:

- Package name in the package.json file: `comigo`
- Reactivity library: `htmx with Alpine.js`
- CSS framework: `Tailwind CSS with daisyUI components`

Tools:

- Air tool to live-reloading: ✓
- Bun as a frontend runtime: ✓
- Templ to generate HTML: ✓
- Config for golangci-lint: ✓

## Folders structure

```console
.
├── script
│   ├── scripts.js
│   └── styles.scss
├── static
│   ├── images
│   │   └── gowebly.svg
│   ├── apple-touch-icon.png
│   ├── favicon.ico
│   ├── favicon.png
│   ├── favicon.svg
│   ├── manifest-desktop-screenshot.jpeg
│   ├── manifest-mobile-screenshot.jpeg
│   ├── manifest-touch-icon.svg
│   └── manifest.webmanifest
├── templates
│   ├── pages
│   │   └── index.templ
│   └── main.templ
├── .gitignore
├── .dockerignore
├── .prettierignore
├── .air.toml
├── golangci.yml
├── Dockerfile
├── docker-compose.yml
├── prettier.config.js
├── package.json
├── go.mod
├── go.sum
├── handlers.go
├── server.go
├── main.go
└── README.md
```

## Developing your project

The backend part is located in the `*.go` files in your project folder.

The `./templates` folder contains Templ templates that you can use in your frontend part. Also, the `./script` folder contains the `styles.scss` (main styles) and `scripts.js` (main scripts) files.

The `./static` folder contains all the static files: icons, images, PWA (Progressive Web App) manifest and other builded/minified script.

## Deploying your project

All deploy settings are located in the `Dockerfile` and `docker-compose.yml` files in your project folder.

To deploy your project to a remote server, follow these steps:

1. Go to your hosting/cloud provider and create a new VDS/VPS.
2. Update all OS packages on the server and install Docker, Docker Compose and Git packages.
3. Use `git clone` command to clone the repository with your project to the server and navigate to its folder.
4. Run the `docker-compose up` command to start your project on your server.

> ❗️ Don't forget to generate Go files from `*.templ` templates before run the `docker-compose up` command.

## About the Gowebly CLI

The [**Gowebly**](https://github.com/gowebly/gowebly) CLI is a next-generation CLI tool that makes it easy to create amazing web applications with **Go** on the backend, using **htmx**, **hyperscript** or **Alpine.js**, and the most popular **CSS frameworks** on the frontend.

It's highly recommended to start exploring the Gowebly CLI with short articles "[**What is Gowebly CLI?**](https://gowebly.org/getting-started)" and "[**How does it work?**](https://gowebly.org/getting-started/how-does-it-work)" to understand the basic principle and the main components built into the **Gowebly** CLI.
