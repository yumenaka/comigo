<h2 align="center">
  <img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo: わかりやすい漫画リーダー" width="200">
  <br>ComiGo: わかりやすい漫画リーダー<br>
</h2>

![Windowssample](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows示例")  

### ダウンロード

```bash
#  curl：
bash <(curl -s https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

#  wget：
bash <(wget -qO- https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

# go 1.23 or higher：
go install github.com/yumenaka/comigo/cmd/comi@latest
```
或いは [github release](https://github.com/yumenaka/comigo/releases ) をご参照ください：　　

[https://github.com/yumenaka/comigo/releases](https://github.com/yumenaka/comigo/releases ) 

### Usage：
```
comi [flags] file_or_dir
```

### 詳細：
わかりやすい漫画リーダー。

1、画像フォルダ、zip、rar、cbr、cbz、epubファイルをサポート。

2、Windows、Linux、MacOSに対応しています。

3、同じLAN上にあるスマホやタブレットで、QRコードをスキャンして読み取ることができます。

4、カウントダウンの長さを設定することができます。sketch-66.exe の名前を sketch-99.exe にすると、カウントダウンが99秒になります。

5、コマンドライン環境で使用可能、config.yaml 設定ファイルをサポート。

### Special Thanks：

[mholt](https://github.com/mholt)  、[spf13](https://github.com/spf13)  [disintegration](https://github.com/disintegration)   、 [Baozisoftware ](https://github.com/Baozisoftware) 、 [markbates](github.com/markbates/pkger)  and more。

## License

This software is released under the MIT license.
