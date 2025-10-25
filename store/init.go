package store

import (
	"sync"
)

// RamStore 扫描后生成。可以有多个子书库。内部使用并发安全的 sync.Map 存储书籍和书组
var RamStore = &StoreInRam{
	StoreInfo: StoreInfo{
		BackendURL:  "comigo://main", // 主书库的 URL
		Name:        "Comigo StoreInfo",
		Description: "Comigo Main book store",
	},
	// 使用 sync.Map 存储书籍和子书库
	ChildStores: sync.Map{}, // 子书库，储存层级关系
}
