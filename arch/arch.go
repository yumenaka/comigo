package arch

import (
	"sync"
)

// 使用sync.Map代替map，保证并发情况下的读写安全
var mapBookFS sync.Map
