<h2 align="center">
  <img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo: 简单粗暴的漫画阅读器" width="200">
  <br>ComiGo: 简单粗暴的漫画阅读器<br>
</h2>

![Windows示例](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows示例")

[English](https://github.com/yumenaka/comi/blob/master/README_EN.md)   [日本語](https://github.com/yumenaka/comi/blob/master/README_JP.md)

### Features：  
1. 支持Windows、Linux、MacOS。支持卷轴、下拉、翻页等多种模式。
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
- [x] 引入多文件机制
- [x] 网页端：简易书架
- [x] 优化打开速度
- [x] 支持插件，可整合frp与webpserver
- [x] 网页端：分享功能
- [x] 网页端：网页判断分辨率，去除服务器分析步骤
- [x] 网页端：显示QRCode
- [x] 网页端：设定展现效果
- [x] 网页端：自定义页面样式，与保存设置
- [x] 网页端：浏览器快捷键
- [x] 网页端：HTTPS加密
- [x] 网页端：显示服务器信息
- [ ] 固实压缩文件（7z）
- [x] websocket通信（[参考](https://github.com/Unrud/remote-touchpad)）
- [ ] 访问权限设置，账号系统
- [ ] 程序log，阅读历史
- [ ] 子命令，download rar2zip 
- [ ] 处理损坏文件，扩展名错误的文件
- [ ] 崩溃后恢复
- [ ] 恶意存档处理?
- [ ] 编写测试程序
- [ ] 支持rar压缩包密码
- [ ] 更准确的文件类型判断
- [ ] 命令行交互
- [ ] 阅读历史记录与统计
- [x] 新一代图片格式支持（heic）
- [x] 图片自动裁边，分割、拼接单双页。
- [ ] 第三方调用API
- [ ] 挂载smb、webdav虚拟文件系统？
- [x] CPU、内存占用、状态监控
- [ ] 文件管理
- [ ] 更新提示
- [ ] 用户系统、访问密码，流量限制等
- [ ] 移动客户端（Android，iOS）
- [ ] 跨平台GUI（[gio](https://gioui.org/) flutter？）

### Special Thanks：
[mholt](https://github.com/mholt)  、[spf13](https://github.com/spf13)  [disintegration](https://github.com/disintegration)   、 [Baozisoftware ](https://github.com/Baozisoftware) 、 [markbates](github.com/markbates/pkger)  and more。

## License

This software is released under the GPL-3.0 license.
