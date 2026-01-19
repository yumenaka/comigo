# Project Instructions

## Code Style
- 代码应该有必要的中文注释，尤其是函数和复杂逻辑部分

## 架构
- 本项目后端使用golang编写，前端使用bun+javascript+alpine.js(+persist插件)+flowbite+tailwindcss编写。虽然有htmx，但在计划逐步减少htmx的使用。
- templ/下面的 *_templ.go文件是templ生成的，不要手动修改。如果需要修改，请修改对应的*.templ文件，然后执行`templ fmt ./templ && templ generate`命令生成。
- 国际化文件，放在`assets/locale`目录下的 json 文件中（en_US.json ja_JP.json zh_CN.json）
- golang后端获取翻译后的字符串：locale.GetString("auth")
- javascript前端获取翻译后的字符串：i18next.t('auth')
- 前端构建命令是`bun run dev`
- golang使用 templ 进行 HTML 模板渲染，模板文件放在`templ`目录下。
- 内置插件放在`templ/plugins`目录下。用户可以设置启用插件，插件会自动加载到特定页面。
- 运行指令：templ fmt ./templ && templ generate && go run main.go