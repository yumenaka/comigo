internal是Go 1.4版本引入的机制，internal目录中的包是本地代码，不能导出到本module以外。

需要导出到外部的代码，不要放在这里。internal目录代码只能被本module内代码引用。

```