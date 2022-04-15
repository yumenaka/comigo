<h2 align="center">
  <img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo: simple comic reader" width="200">
  <br>ComiGo: simple comic reader<br>
</h2>

![Windows Example]  

[中文文档][README CN]   [日本語][README JP]

## Supported platforms：

- Windows(386\64)
- macOS
- Linux（amd64、armv7\armv8）
- and More(You can compile it yourself).

### Install：
Download [the latest release for your system][Releases] 

### Usage：
```
comi [flags] file_or_dir
#more：
comi --help
```

### Special Thanks：
[mholt]、[spf13]、[disintegration]、[Baozisoftware]、[markbates] and more。

### Possible Problems：

#### MIME type ('text/plain') error。

[If you use Windows，and webpage cannot be loaded：][Windows Page Issue]

```bash
1、Win+R 
2、type "regedit".
3、into the registry \HKEY_CLASSES_ROOT.js > "Content Type" 4、 changed it from "text/plain" to "application/javascript" and i restart the Comigo and that fixed it
```

## License

This software is released under the GPL-3.0 license.


<!----------------------------------------------------------------------------->

[README JP]: README_JP.md
[README CN]: README_CN.md

[Windows Page Issue]: https://github.com/golang/go/issues/32350
[Releases]: https://github.com/yumenaka/comi/releases

[Windows Example]: https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif 'Windows Example'


[disintegration]: https://github.com/disintegration
[Baozisoftware]: https://github.com/Baozisoftware
[markbates]: github.com/markbates/pkger
[mholt]: https://github.com/mholt
[spf13]: https://github.com/spf13
