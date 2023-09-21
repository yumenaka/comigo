来自facebook，官方简介：
https://github.com/ent/ent/blob/master/README_zh.md
文档：
https://entgo.io/zh/docs/tutorial-setup/

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
go run entgo.io/ent/cmd/ent init User Book

# 新建User与Book实体
go run -mod=mod entgo.io/ent/cmd/ent User Book

# 应该编辑 ent/schema/book.go 与 ent/schema/user.go。
# 不应编辑生成的文件（ent/book.go 与 ent/user.go等等）。重新生成时，修改将消失。
# 生成CRUD相关代码。每次添加或修改 fields 和 edges后, 都需要生成新的实体. 
# 在项目的根目录执行 ent generate或直接执行：
go generate ./ent
```
