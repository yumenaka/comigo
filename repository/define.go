package repository

import (
	"sync"
)

// MainStore 带有层级关系的总书组，用于前端展示
var MainStore = memoryStore{}

// memoryStore 本地总书库，扫描后生成。可以有多个子书库。
// 使用并发安全的 sync.Map 存储书籍和书组
type memoryStore struct {
	mapBooks     sync.Map // 实际存在的书 key: string (BookID), value: *Book
	mapBookGroup sync.Map // 虚拟书组    key: string (BookID), value: *BookGroup
	SubStores    sync.Map // key为路径 存储 *subMemoryStore
}

// 对应某个扫描路径的子书库
type subMemoryStore struct {
	Path         string   // 扫描路径
	SortBy       string   // 排序方式
	BookMap      sync.Map // 拥有的书籍,key是BookID,存储 *BookInfo
	BookGroupMap sync.Map // 拥有的书籍组,通过扫描书库生成，key是BookID,存储 *BookInfo。需要通过Depth从深到浅生成
}
