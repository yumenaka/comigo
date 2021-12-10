<h2 align="center">
  <img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo: 简单粗暴的漫画阅读器" width="200">
  <br>ComiGo: 简单粗暴的漫画阅读器<br>
</h2>

![Windows示例](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows示例")

[English](https://github.com/yumenaka/comi/blob/master/README_EN.md)   [日本語](https://github.com/yumenaka/comi/blob/master/README_JP.md)

### Features：  
1. 支持Windows、Linux、MacOS命令行调用，支持卷轴、下拉、翻页等多种模式。
2. 支持图片文件夹与.rar、.zip、.tar、.cbz、.cbr、.epub压缩包。
3. 局域网的手机或平板设备，可扫描二维码阅读。  
4. windows支持拖拽压缩包到comi.exe（或快捷方式）上打开。
5. frp转发、webp-server压缩等扩展功能
### 安装：
在 [Releases页面](https://github.com/yumenaka/comi/releases ) ，下载最新版文件。

### 用法：
```
comi [flags] file_or_dir
# more
comi --help
```

### TODO：
- [ ] 引入多文件机制
- [ ] 网页端：简易书架
- [ ] 优化打开速度，压缩包后台解压
- [ ] 虚拟文件系统、zip不用解压 
- [ ] 固实压缩文件（7z）
- [ ] 扫描文件夹，生成书单或书架
- [ ] websocket通信（[参考](https://github.com/Unrud/remote-touchpad)）
- [ ] 访问权限设置，账号系统
- [ ] 程序log，阅读历史
- [ ] 处理损坏文件，扩展名错误的文件
- [ ] 崩溃后恢复
- [ ] rar2zip？
- [ ] 恶意存档处理?
- [ ] 编写测试程序
- [ ] 支持rar压缩包密码
- [ ] 更准确的文件类型判断
- [ ] 命令行交互
- [ ] 支持中间件，可整合frp与webpserver
- [ ] 网页端：分享功能
- [ ] 网页端：网页判断分辨率，去除服务器分析步骤
- [ ] 网页端：显示QRCode
- [ ] 网页端：配合后台解压，更新本地数据
- [ ] 网页端：设定展现效果
- [ ] 网页端：自定义页面样式，与保存设置
- [ ] 网页端：浏览器快捷键
- [ ] 网页端：HTTPS加密
- [ ] 网页端：显示服务器信息
- [ ] 阅读历史记录与统计
- [ ] 新一代图片格式支持（heic）
- [ ] 图片自动裁边，分割、拼接单双页。
- [ ] [pdf支持](https://github.com/pdfcpu/pdfcpu)，直接交给epub.js与浏览器？
- [ ] 第三方调用API
- [ ] 挂载smb、webdav虚拟文件系统？
- [ ] CPU、内存占用、状态监控
- [ ] 文件管理与书库
- [ ] 支持HTTPS
- [ ] 更新提示
- [ ] 用户系统、访问密码，流量限制等
- [ ] 移动客户端（Android，iOS）
- [ ] 跨平台gui界面（[gio](https://gioui.org/) flutter？）

### Special Thanks：
[mholt](https://github.com/mholt)  、[spf13](https://github.com/spf13)  [disintegration](https://github.com/disintegration)   、 [Baozisoftware ](https://github.com/Baozisoftware) 、 [markbates](github.com/markbates/pkger)  and more。

## License

This software is released under the GPL-3.0 license.
