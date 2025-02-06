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
 - [ ] 服务器设置页面 v1.0
 - [ ]  [self update 功能](https://github.com/minio/selfupdate)
 - [ ] 显示服务器log：[web终端示例](https://zenn.dev/ikedam/articles/2e078bfc2a4cb6)，设置页面功能
 - [ ] 优化翻页模式 
 - [ ] 自定义js与css代码块功能。
 - [ ] [使用 Go1.24 的 os.Root 类型](https://antonz.org/go-1-24/)，将文件系统操作限制为特定目录，以防止用户访问系统文件
 - [ ] 等到go 1.24正式发布，[使用 Go1.24 的 go get -tool](https://antonz.org/go-1-24/)，添加工具依赖项
 - [ ] 等 [gowebly](https://github.com/gowebly/gowebly) 更新后，参考 gowebly 的模板，升级tailwindcss到4.0
 - [ ] 自动发版功能 [goreleaser](https://goreleaser.com/)  [github-action](https://dev.to/hadlow/how-to-release-to-homebrew-with-goreleaser-github-actions-and-semantic-release-2gbb)
 - [ ] 设置页面：显示当前用户状态、阅读书籍、阅读进度、阅读时间、服务器状态的地方。
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
├── assets
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

The `./templates` folder contains Templ templates that you can use in your frontend part. Also, the `./assets` folder contains the `styles.scss` (main styles) and `scripts.js` (main scripts) files.

The `./static` folder contains all the static files: icons, images, PWA (Progressive Web App) manifest and other builded/minified assets.

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
