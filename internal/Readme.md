internal是Go 1.4版本引入的机制，internal目录中的包是本地代码，不能导出到本module以外。
虽然internal目录没有限制位置，但一般都会将internal放在项目根目录。

需要导出到外部的代码，不要放在这里。这里的代码只能被本module内代码引用。

```