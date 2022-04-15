[![Badge License]][License]

<div align = center>

<img width = 200 alt = 'ComiGo Logo' align = right src = icon.ico>

<br>

# ComiGo

**A Simple Manga Reader**

[:jp:][README JP] [:cn:][README CN]

![Windows Example]  

</div>


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


<!----------------------------------------------------------------------------->

[Badge License]: https://img.shields.io/badge/License-GPLv3-blue.svg?style=for-the-badge

[README JP]: README_JP.md '日本語'
[README CN]: README_CN.md '中文文档'

[Windows Page Issue]: https://github.com/golang/go/issues/32350
[Releases]: https://github.com/yumenaka/comi/releases
[License]: LICENSE

[Windows Example]: https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif 'Windows Example'


[disintegration]: https://github.com/disintegration
[Baozisoftware]: https://github.com/Baozisoftware
[markbates]: github.com/markbates/pkger
[mholt]: https://github.com/mholt
[spf13]: https://github.com/spf13
