<h2 align="center">
  <img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo: simple comic reader" width="200">
  <br>ComiGo: simple comic reader<br>
</h2>

![Windows example](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows example")  

[中文文档](https://github.com/yumenaka/comi/blob/master/README_CN.md)   [日本語](https://github.com/yumenaka/comi/blob/master/README_JP.md)

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
#more：
comi --help
```

### Special Thanks：
[mholt](https://github.com/mholt)  、[spf13](https://github.com/spf13)  [disintegration](https://github.com/disintegration)   、 [Baozisoftware ](https://github.com/Baozisoftware) 、 [markbates](github.com/markbates/pkger)  and more。

### Possible Problems：

#### MIME type ('text/plain') error。

[If you use Windows，and webpage cannot be loaded：](https://github.com/golang/go/issues/32350)

```bash
1、Win+R 
2、type "regedit".
3、into the registry \HKEY_CLASSES_ROOT.js > "Content Type" 4、 changed it from "text/plain" to "application/javascript" and i restart the Comigo and that fixed it
```

## License

This software is released under the GPL-3.0 license.
