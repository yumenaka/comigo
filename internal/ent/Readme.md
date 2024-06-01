来自facebook，官方简介：
https://github.com/ent/ent/blob/master/README_zh.md
文档：
https://entgo.io/zh/docs/tutorial-setup/


Supported platforms and architectures
https://pkg.go.dev/modernc.org/sqlite#hdr-Supported_platforms_and_architectures
https://modern-c.appspot.com/-/builder/?importpath=modernc.org%2fsqlite
主要平台里面，也就windows 386不支持。

Adds support for Go fs.FS based SQLite virtual filesystems, see function New in modernc.org/sqlite/vfs and/or TestVFS in all_test.go
添加对 Go fs 的支持。基于 FS 的 SQLite 虚拟文件系统，请参阅函数 modernc.org/sqlite/vfs 中的新功能和/或 all_test.go 中的 TestVFS

https://gitlab.com/cznic/sqlite/-/blob/master/all_test.go
大约在2487行，搜索TestVFS，有示例代码。似乎可以内嵌数据库文件？可以存储默认配置，但无法保存的话，实际也没太大用处？



ent是一个简单而又功能强大的Go语言实体框架，ent易于构建和维护应用程序与大数据模型。
图就是代码 - 将任何数据库表建模为Go对象。
轻松地遍历任何图形 - 可以轻松地运行查询、聚合和遍历任何图形结构。
静态类型和显式API - 使用代码生成静态类型和显式API，查询数据更加便捷。
多存储驱动程序 - 支持MySQL, PostgreSQL, SQLite 和 Gremlin。
可扩展 - 简单地扩展和使用Go模板自定义。

100%型安全なgolangORM「ent」を使ってみた
https://future-architect.github.io/articles/20210728a/

```bash
# 在项目根目录执行，生成设计图（Schema）模板。
# 会生成 与对应的 schema/ent/book.go 与 schema/ent/user.go,编辑这些文件来定义实体的属性。
#go run entgo.io/ent/cmd/ent init User Book
因为我的目录不在根目录下，所以应该指定路径

go run entgo.io/ent/cmd/ent init  --target /internal/ent/schema/User Book

# 新建User与Book实体
go run -mod=mod entgo.io/ent/cmd/ent User Book

# 应该编辑 ent/schema/book.go 与 ent/schema/user.go。
# 不应编辑生成的文件（ent/book.go 与 ent/user.go等等）。重新生成时，修改将消失。
# 生成CRUD相关代码。每次添加或修改 fields 和 edges后, 都需要生成新的实体. 
# 在项目的根目录执行 ent generate或直接执行：
go generate ./internal/ent
```
