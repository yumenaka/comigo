<h2 align="center">
  <img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo: 简单粗暴的漫画阅读器" width="200">
  <br>ComiGo: 简单粗暴的漫画阅读器<br>
</h2>

![Windows示例](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows示例")

[English](https://github.com/yumenaka/comi/blob/master/README_EN.md)   
[日本語](https://github.com/yumenaka/comi/blob/master/README_JP.md)

### 安装：

在 [Releases页面](https://github.com/yumenaka/comi/releases ) ，下载最新版文件，并放入系统PATH。

### 用法：
```
comi [flags] file_or_dir

设定网页服务端口（默认为1234）：
comi -p 2345 book.zip

不打开浏览器（windows）：
comi -o=false book.zip

指定多个参数：
comi -w -q 70 --frpc  --token aX4457d3O -p 23455 --frps-addr sh.example.com test.zip

Flags:
      --check-image                在服务器端分析图片分辨率 (default true)
  -c, --config string              指定配置文件
      --debug                      Debug Mode
  -d, --disable-lan                禁用LAN分享
  -f, --frpc                       启用frp反向代理
      --frpc-command string        frpc命令,或frpc可执行文件路径 (default "frpc")
      --frps-addr string           frps-addr, frpc必须 (default "frps.example.com")
      --frps-port int              frps server_port, frpc必须 (default 7000)
      --frps-random-remote         frp 远程随机端口（40000~50000） (default true)
      --frps-remote-port int       frp 远程端口,如果设置为-1，端口与本地相同 (default 50000)
  -h, --help                       help for comi
      --host string                自定义域名
      --log                        记录log文件
  -m, --max-depth int              最大搜索深度 (default 1)
      --min-image-num int          至少有几个媒体文件，才认定为漫画压缩包 (default 1)
  -o, --open-browser               同时打开浏览器，windows=true
  -p, --port int                   服务端口 (default 1234)
      --print-all-ip               打印所有可用网卡ip
      --sketch_count_seconds int   速写倒计时秒数 (default 90)
      --sort string                图片重拍规则(none、name、time) (default "none")
  -t, --template string            默认页面模板(scroll,single,double,sketch) (default "scroll")
      --token string               token, frpc必须 (default "token_secretSAMPLE")
  -v, --version                    version for comi
  -w, --webp                       webp传输，需要webp-server
      --webp-command string        webp-server命令,或webp-server可执行文件路径 (default "webp-server")
  -q, --webp-quality int           webp压缩质量 (default 85)

```

### 介绍：
1. **支持Windows、Linux、MacOS命令行调用**（主要为了远程看，拖拽功能只是附带）。
2. **支持压缩包（rar、zip、tar等）与图片文件夹。**
3. 同一**局域网的手机或平板设备，可扫描二维码阅读。**  
4. 支持下拉阅读、翻页阅读模板（不支持IE浏览器）。
5. 图片文件，解压到系统临时文件夹，在正常退出（Ctrl+c）时清除。
6. 提供原始文件的下载链接，可以直接下载看。   
7. 默认端口被占用，比如重复打开了多个文件时，将使用随机端口号。  
8. **命令行调用，需要将Comi程序所在目录，添加到PATH环境变量**。  
- [如何设置或更改 PATH 系统变量](https://www.java.com/zh_CN/download/help/path.xml)
    - 如果实在看不懂，就扔到“C:\Windows”里面吧……
9. windows支持拖拽压缩包到comi.exe（或exe文件的快捷方式）上打开。  
10. 可以在快捷方式上，添加命令行参数（下图为“取消自动打开浏览器”。）
    - ![取消“自动打开浏览器”的功能](https://www.yumenaka.net/wp-content/uploads/2020/08/tips1-1.png "取消自动打开浏览器的功能")
11. **支持webp压缩，需要下载webp-server文件，启用后可节省30~60%流量。**
     - [配套文件下载地址](https://github.com/webp-sh/webp_server_go/releases/latest)
     - 下载后，请命名为“webp-server”（或webp-server.exe），并放入系统PATH。  
12. Ctrl+c 退出，会清空系统临时文件夹的漫画缓存。  
13. 直接关闭或kill的话，到下一次Ctrl+C退出时，才会清空缓存。  
15. 根据图片分辨率、屏幕可见区域比例，区分单双页，自适应手机或平板。  
15. PC浏览，可使用“Ctrl +”，“Ctrl -”缩放。  
16. 程序使用upx压缩，如有问题，可自行编译，或通知错误。
17. **毛病很多、勉强能用，欢迎提供建议、帮忙修改。**

### TODO：
- [x] 代码开源
- [ ] 多文件展示与切换
- [ ] 网页端：简易书架
- [x] 网页端：下拉之外的翻页模式
- [x] 本地配置文件
- [x] 读取系统变量
- [ ] 优化打开速度，压缩包后台解压
- [ ] 处理损坏文件，或扩展名错误的文件
- [ ] 一键安装脚本，go get 下载
- [ ] 支持非UTF-8编码的ZIP文件
- [ ] 扫描文件夹，生成书单或书架
- [x] 文件排序（文件名、修改时间、路径）
- [ ] websocket通信（[参考](https://github.com/Unrud/remote-touchpad)）
- [ ] 访问权限设置，账号系统？
- [x] 程序log，阅读历史
- [ ] 自定义临时文件夹
- [ ] 崩溃后恢复
- [ ] 编写测试程序
- [ ] 支持更多压缩包格式（.7z）
- [ ] 支持rar压缩包密码
- [ ] 更准确的文件类型判断
- [ ] 命令行交互
- [ ] 网页端：分享功能
- [ ] 网页端：显示QRCode
- [ ] 网页端：配合后台解压，更新本地数据
- [ ] 网页端：设定展现效果
- [ ] 网页端：自定义页面样式，与保存设置
- [x] 网页端：双页模式
- [x] 网页端：模板文件
- [ ] 网页端：浏览器快捷键
- [ ] 网页端：HTTPS加密
- [ ] 网页端：显示服务器信息
- [ ] 阅读历史记录与统计
- [x] 程序国际化
- [ ] 支持新一代图片格式（heic）
- [x] 插件：frp端口映射（内网穿透，远程访问）
- [ ] 图片自动裁边，分割、拼接单双页。
- [ ] pdf、epub支持，直接交给epub.js与浏览器？
- [ ] 第三方调用API
- [ ] CPU、内存占用、状态监控
- [ ] 文件管理与书库
- [ ] 支持HTTPS
- [x] 优化资源占用
- [ ] 更新提示
- [ ] 用户系统、访问密码，流量限制等
- [ ] 移动客户端（Android，iOS）
- [ ] 跨平台gui界面（[gio](https://gioui.org/)？）

### Special Thanks：
[mholt](https://github.com/mholt)  、[spf13](https://github.com/spf13)  [disintegration](https://github.com/disintegration)   、 [Baozisoftware ](https://github.com/Baozisoftware) 、 [markbates](github.com/markbates/pkger)  and more。

## License

This software is released under the GPL-3.0 license.
