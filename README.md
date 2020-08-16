## Comi Go: 简单粗暴的漫画阅读器  

![sample](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "sample")  

### 用法：
```
Usage:
  comi [flags] file_or_dir

Examples:

comi book.zip

设定网页服务端口（默认1234）：
comi -p 2345 bookdir

不打开浏览器（windows）：
comi -b=false book.rar

webp传输，需要webp-server配合：
comi -w book.zip

指定多个参数：
comi -lw -webp-command=C:\Users\test\Desktop\webp-server-windows-amd64.exe -p 3344 -q 45  test.zip

Flags:
  -b, --broswer               打开浏览器，windows=true
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

### 下载：
[https://github.com/yumenaka/comi/releases]（https://github.com/yumenaka/comi/releases）

### 补充：
1、支持压缩包（rar、zip）与图片文件夹。  
2、在同一局域网的设备，可使用链接，或扫描二维码阅读。  
3、可重复打开多个文件，此时将使用随机端口号。  
4、命令行使用，需要将程序路径添加到系统PATH。  
5、windows支持拖拽压缩包到exe（或exe文件的快捷方式）上打开。  
6、可以在修改快捷方式，来制定一些参数，比如这样就可以让windows程序，不再自动打开浏览器。  
7、支持第三方程序webp-server，压缩传输（[下载地址](https://github.com/webp-sh/webp_server_go/releases/latest)）。下载后，请命名为“webp-server”，并加入系统PATH。  
8、Ctrl+C退出，会清空系统临时文件夹的漫画缓存。  
9、直接关闭或kill，到下一次Ctrl+C退出时，才清空缓存。  
10、根据图片分辨率，区分单双页。   
11、根据屏幕可见区域比例，自适应手机或平板。   
12、PC单页状态，可使用“Ctrl +”，“Ctrl-”缩放。  
13、毛病很多、勉强能用，欢迎大佬建议与修改。  

### Special Thanks：
[mholt](https://github.com/mholt)  
[spf13](https://github.com/spf13)  
[disintegration](https://github.com/disintegration)  

