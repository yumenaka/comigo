[![Badge License]][License]

<div align = center>

<img width = 200 alt = 'ComiGo Logo' align = right src = icon.ico>

<br>

# ComiGo

**A Simple Manga Reader**

[:jp:][README JP] [:cn:][README CN]

![Windows Example]  

</div>

---

<div align = center>

<br>
  


<br>

---

<br>

## Install
```bash
#  curl：
bash <(curl -s https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

#  wget：
bash <(wget -qO- https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

# go 1.23 or higher：
go install github.com/yumenaka/comigo/cmd/comi@latest
```
Or Download [the latest release for your system][Releases] 


## Usage

<div align = left>

```sh
comi [flags] file_or_dir
```

```sh
comi --help
```

</div>

<br>

---

<br>

## Special Thanks

**[mholt]** **[spf13]** **[disintegration]** <br>
**[Baozisoftware]** **[markbates]** <br>
and more..

<br>

---

<br>

## Possible Problems

</div>

#### MIMEType / `text/plain` Error

*[If you use Windows，and webpage cannot be loaded.][Windows Page Issue]*

1. **<kbd>Win</kbd> + <kbd>R</kbd>**

2. Type `regedit`

3. In the registry navigate to `\HKEY_CLASSES_ROOT.js` **>** `Content Type`

4. Change `text/plain` to `application/javascript`

5. Restart **ComiGo**


<!----------------------------------------------------------------------------->

[Badge License]: https://img.shields.io/badge/License-GPLv3-blue.svg?style=for-the-badge
[Badge Windows]: https://img.shields.io/badge/Windows_-32_/_64-0078D6?style=for-the-badge&logo=windows&logoColor=white
[Badge Linux]: https://img.shields.io/badge/Linux-AMD64_/_ARMv7/8-10B981?style=for-the-badge&logo=linux&logoColor=white
[Badge MacOS]: https://img.shields.io/badge/MacOS-999999?style=for-the-badge&logo=apple&logoColor=white

[README JP]: README_JP.md '日本語'
[README CN]: README_CN.md '中文文档'

[Windows Page Issue]: https://github.com/golang/go/issues/32350
[Releases]: https://github.com/yumenaka/comigo/releases
[License]: LICENSE

[Windows Example]: https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif 'Windows Example'


[disintegration]: https://github.com/disintegration
[Baozisoftware]: https://github.com/Baozisoftware
[markbates]: github.com/markbates/pkger
[mholt]: https://github.com/mholt
[spf13]: https://github.com/spf13
