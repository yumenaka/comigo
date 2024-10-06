<div align="center">

# ComiGo: 简单粗暴的漫画阅读器 
[![Go Report](https://goreportcard.com/badge/github.com/yumenaka/comi?style=flat-square)](https://goreportcard.com/report/github.com/yumenaka/comi)
[![License](https://img.shields.io/github/license/yumenaka/comi?style=flat-square&color=blue)](https://github.com/yumenaka/comigo/blob/main/LICENSE)

[//]: # ([![Downloads]&#40;https://img.shields.io/github/downloads/yumenaka/comi/total?style=flat-square&color=success&#41;]&#40;https://github.com/yumenaka/comigo/releases&#41;)
[//]: # (<img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo: 简单粗暴的漫画阅读器" width="200">)
</div>

![Windows示例](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows示例")

[English](https://github.com/yumenaka/comigo/blob/master/README_EN.md)   [日本語](https://github.com/yumenaka/comigo/blob/master/README_JP.md)

### Features：  
1. 支持Windows、Linux、MacOS。支持卷轴、下拉、翻页等多种模式。
2. 支持图片文件夹与.rar、.zip、.tar、.cbz、.cbr、.epub压缩包。
3. 局域网的手机或平板设备，可扫描二维码阅读。  
4. windows支持拖拽压缩包到comi.exe（或快捷方式）上打开。

### 安装：
手动下载  
在 [Releases页面](https://github.com/yumenaka/comigo/releases ) ，下载最新版文件，放到系统PATH。
Linux MacOS 一键安装脚本  
```bash
# 需要curl与tar，文件将安装到/usr/local/bin/ 
bash <(curl -s https://raw.githubusercontent.com/yumenaka/comi/master/get_comigo.sh)

#  如果你设置了golang环境，也可以这样安装：
go install github.com/yumenaka/comi@latest
```
### 用法：
```
comi [flags] file_or_dir
# more
comi --help
```

### 配置文件：
Comigo服务器设置，可选的保存位置：

**Home目录**  
一般是`C:\Users\用户名\.config\comigo.toml`，或者`/home/用户名/.config/comigo.toml`。  
程序启动时，默认读取这个文件。如果只是命令行使用，可以不使用配置文件。

**程序所在目录**
`comigo.toml`  
与可执行文件放在一起，同样也是启动时生效。当作绿色软件使用，可以保存到这个位置。

**当前运行目录**  
如果你想把配置文件放在**当前运行目录**。切换到这个目录以后，启动命令时生效。

**用户指定目录**  
在命令行中调用时，指定`--config`参数，也可以指定任意位置的配置文件。


### 特性与 Todo：
- [x] 多文件支持
- [x] 网页书架
- [x] 优化打开速度
- [x] 新一代图片格式支持（heic avif）。
- [x] 图片自动裁边，分割、拼接单双页。
- [x] 网页端：分享功能
- [x] 网页端：显示QRCode
- [x] 网页端：多种展示模式
- [x] 网页端：服务器设置
- [x] 网页端：HTTPS加密
- [x] 网页端：显示服务器信息
- [x] 网页端：上一章、下一章,快速跳转。
- [x] websocket通信（[参考](https://github.com/Unrud/remote-touchpad)）
- [x] 访问权限设置，账号系统
- [x] log记录
- [x] 设置中心，设置热重载
- [x] CPU、内存占用、状态监控
- [x] 网页端：卷轴模式分页。
- [ ] 画个示例漫画。
- [ ] PWA模式。
- [ ] 优化配置文件 （[参考](https://toml.io/cn/v1.0.0） (better config file formart).
- [ ] 嵌入html，防剧透效果。回忆模式，特殊背景，音乐etc
- [ ] 网页端：优化图片预加载，长图片支持。
- [x] 网页端：添加预定义主题与颜色。
- [ ] 静态绑定模式
- [ ] 网页端：内置帮助文档。
- [ ] 网页端：二维码界面文本显示链接
- [ ] 网页端：网页前端查看log
- [ ] 跨平台GUI（flutter+gomobile）
- [ ] 更新提示，自动更新。
- [ ] 文件夹监控，自动更新
- [ ] 文件持久化，meta文件，阅读历史与统计。
- [ ] 用户系统、访问密码，流量限制等
- [ ] 网页端：浏览器快捷键。
- [ ] shell 互动（https://github.com/rivo/tview ）
- [ ] 子命令，download rar2zip 
- [ ] 支持rar压缩包密码。处理损坏文件，扩展名错误的文件，固实压缩文件（7z）。更准确的文件类型判断。
- [ ] 崩溃后恢复，恶意存档处理。
- [ ] 编写测试
- [ ] 命令行交互
- [ ] 调用第三方API
- [ ] 挂载smb、webdav
- [ ] 文件管理，删除。
- [ ] 移动客户端（Android，iOS）
- [ ] Debian，RPM包（https://github.com/goreleaser/nfpm）
- [ ] 优化epub与PDF阅读体验，支持图文混排（pdf.js与epub.js）

### Special Thanks：
[mholt](https://github.com/mholt)  、[spf13](https://github.com/spf13)  [disintegration](https://github.com/disintegration)   、 [Baozisoftware ](https://github.com/Baozisoftware) 、 [markbates](github.com/markbates/pkger)  and more。

## License

This software is released under the MIT license.

## Stargazers over time
[![Stargazers over time](https://starchart.cc/yumenaka/comi.svg?variant=adaptive)](https://starchart.cc/yumenaka/comi)
