<h2 align="center">
  <img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo: simple comic reader" width="200">
  <br>ComiGo: simple comic reader<br>
</h2>

![Windows example](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows example")  

## Supported platforms：

- Windows(386\64)
- macOS
- Linux（amd64、armv7\armv8）
- and More(You can compile it yourself).

### Install：
Download [the latest release for your system](https://github.com/yumenaka/comi/releases ) 

### Usage：
```
comi [flags] file_or_dir

Set the web service port (default is 1234): 
comi -p 2345 book.zip 

Without opening a browser (windows):
comi -o=false book.zip
   
For local browsing, not available for LAN:
comi -l book.zip 

webp transfer, requiring webp-server to work with: 
comi -w book.zip 

Multiple parameters:
comi -w -q 70 --frpc  --token aX4457d3O -p 23455 --frps-addr sh.example.com test.zip
```

### Special Thanks：
[mholt](https://github.com/mholt)  、[spf13](https://github.com/spf13)  [disintegration](https://github.com/disintegration)   、 [Baozisoftware ](https://github.com/Baozisoftware) 、 [markbates](github.com/markbates/pkger)  and more。

## License

This software is released under the GPL-3.0 license.
