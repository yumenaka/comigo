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

### 安装：
1. 手动下载
在 [Releases页面](https://github.com/yumenaka/comi/releases ) ，下载最新版文件，放到系统PATH。
2. Linux MacOS 一键安装脚本  
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

### 特性与 Todo：
- [x] 多文件支持
- [x] 网页书架
- [x] 优化打开速度
- [x] 可配合frp与webpserver
- [x] 网页端：分享功能
- [x] 网页端：显示QRCode
- [x] 网页端：多种展示模式
- [x] 网页端：服务器设置
- [ ] 网页端：浏览器快捷键
- [x] 网页端：HTTPS加密
- [x] 网页端：显示服务器信息
- [ ] 网页端：上一章、下一章
- [ ] web shell like ttyd
- [x] websocket通信（[参考](https://github.com/Unrud/remote-touchpad)）
- [x] 访问权限设置，账号系统
- [x] log记录
- [x] 设置中心，设置热重载
- [ ] 固实压缩文件（7z）
- [ ] 子命令，download rar2zip 
- [ ] 处理损坏文件，扩展名错误的文件
- [ ] 崩溃后恢复
- [ ] 恶意存档处理
- [ ] 编写测试程序
- [ ] 支持rar压缩包密码
- [ ] 更准确的文件类型判断
- [ ] 命令行交互
- [ ] 阅读历史与统计
- [x] 新一代图片格式支持（heic）
- [x] 图片自动裁边，分割、拼接单双页。
- [ ] 调用第三方API
- [ ] 挂载smb、webdav
- [x] CPU、内存占用、状态监控
- [ ] 文件管理
- [ ] 更新提示
- [ ] 用户系统、访问密码，流量限制等
- [ ] 移动客户端（Android，iOS）
- [ ] 跨平台GUI（react naitive or flutter？）

### Special Thanks：
[mholt](https://github.com/mholt)  、[spf13](https://github.com/spf13)  [disintegration](https://github.com/disintegration)   、 [Baozisoftware ](https://github.com/Baozisoftware) 、 [markbates](github.com/markbates/pkger)  and more。

## License

This software is released under the MIT license.
