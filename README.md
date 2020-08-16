<h2 align="center">
  <img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo: 简单粗暴的漫画阅读器" width="200">
  <br>ComiGo: 简单粗暴的漫画阅读器<br>
</h2>

![Windows示例](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows示例")  

### 下载地址：

[https://github.com/yumenaka/comi/releases](https://github.com/yumenaka/comi/releases ) 

### Usage：
```

comi [flags] file_or_dir

设定网页服务端口（默认1234）：
comi -p 2345 bookdir

不打开浏览器（windows）：
comi -b=false book.rar

webp传输，需要webp-server配合：
comi -w book.zip

指定多个参数：
comi -lw -webp-command=C:\Users\test\Desktop\webp-server-windows-amd64.exe -p 3344 -q 45  test.zip

Flags:
  -b, --broswer               自动打开浏览器，windows=true
  -h, --help                  help for comi
  -i, --imagenum int          有几张图片，才认定为漫画 (default 3)
  -l, --local-only            禁用LAN分享
      --log                   记录log文件
  -m, --max-depth int         最大搜索深度 (default 2)
  -p, --port int              服务端口 (default 1234)
      --print-allip           打印所有可用ip
  -g, --usego                 启用并发
  -v, --version               输出版本号
  -w, --webp                  启用webp压缩
      --webp-command string   webp-server命令，或可执行文件路径 (default "webp-server")
  -q, --webp-quality string   webp压缩质量 (default "60")

```


### 详细介绍：
1. **支持Windows、Linux、MacOS命令行调用**（主要为了远程看，拖拽功能只是附带）。
2. **支持压缩包（rar、zip、tar等）与图片文件夹。**
3. 同一**局域网的手机或平板设备，可扫描二维码阅读。**  
4. 目前只支持下拉阅读，也不支持IE浏览器。
5. 图片文件，解压到系统临时文件夹，在正常退出（Ctrl+c）时清除。
6. 阅读页书名，就是原始文件的下载链接，如果出BUG，可以直接下载看。   
7. 默认端口被占用，比如重复打开了多个文件时，将使用随机端口号。  
8. **命令行调用，需要将ComiGo程序所在目录，添加到PATH环境变量**。  
- [如何设置或更改 PATH 系统变量](https://www.java.com/zh_CN/download/help/path.xml)
    - 如果实在看不懂，就扔到“C:\Windows”里面吧……
9. windows支持拖拽压缩包到comi.exe（或exe文件的快捷方式）上打开。  
10. 可以在快捷方式上，添加命令行参数（下图为“取消自动打开浏览器”。）
    - ![取消“自动打开浏览器”的功能](https://www.yumenaka.net/wp-content/uploads/2020/08/tips1.png "取消自动打开浏览器的功能")
11. **支持webp压缩，需要下载webp-server文件，启用后可节省30~60%流量。**
     - [配套文件下载地址](https://github.com/webp-sh/webp_server_go/releases/latest)
     - 下载后，请命名为“webp-server”（或webp-server.exe），并放入系统PATH。  
12. Ctrl+c 退出，会清空系统临时文件夹的漫画缓存。  
13. 直接关闭或kill的话，到下一次Ctrl+C退出时，才会清空缓存。  
15. 根据图片分辨率、屏幕可见区域比例，区分单双页，自适应手机或平板。  
15. PC浏览，可使用“Ctrl +”，“Ctrl -”缩放。  
16. 程序使用upx压缩，如有问题，可自行编译，或通知错误。
17. **毛病很多、勉强能用，欢迎提供建议、帮忙修改。**  
18. 解压速度取决于机器配置，有加速方案，但目前没写完。

### TODO：
- [x] 代码开源
- [ ] 多文件展示与切换
- [ ] 网页端：简易书架
- [ ] 网页端：下拉之外的翻页模式
- [ ] 本地配置文件
- [ ] 读取系统变量
- [ ] 优化打开速度，压缩包后台解压
- [ ] BUG：无法读取损坏文件，或扩展名错误的文件
- [ ] BUG：无法使用 go get 下载
- [ ] 支持非UTF-8编码的ZIP文件
- [ ] 扫描文件夹，生成书单或书架
- [ ] 文件排序（文件名、修改时间、路径）
- [ ] websocket通信（[参考](https://github.com/Unrud/remote-touchpad)）
- [ ] 访问权限设置，账号系统？
- [ ] 程序log，阅读历史
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
- [ ] 网页端：双页模式
- [ ] 网页端：切换模板文件
- [ ] 网页端：浏览器快捷键
- [ ] 网页端：阅读密码
- [ ] 网页端：HTTPS加密
- [ ] 网页端：显示服务器信息
- [ ] 阅读历史记录与统计
- [ ] 程序国际化
- [ ] 支持新一代图片格式（heic）
- [ ] 插件：frp端口映射
- [ ] 图片自动裁边，分割、拼接单双页。
- [ ] pdf、epub支持，直接交给epub.js与浏览器？
- [ ] 第三方调用API
- [ ] CPU、内存占用、状态监控
- [ ] 文件管理与书库
- [ ] 远程访问助手（内网穿透）
- [ ] 支持HTTPS
- [ ] 优化资源占用
- [ ] 更新提示
- [ ] 用户系统、访问密码，流量限制等
- [ ] 移动端支持（Android，iOS）
- [ ] 跨平台gui界面（[gio](https://gioui.org/)？）

### Special Thanks：
[mholt](https://github.com/mholt)  、[spf13](https://github.com/spf13)  [disintegration](https://github.com/disintegration)   、 [Baozisoftware ](https://github.com/Baozisoftware) 、 [markbates](github.com/markbates/pkger)  and more。

## License

This software is released under the GPL-3.0 license.
