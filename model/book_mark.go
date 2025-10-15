package model

type BookMark struct {
	BookID      string  `json:"book_id"`     // 书籍 ID
	PageIndex   int     `json:"page_index"`  // 书签页码，从 0 开始，不会超过 PageCount - 1
	Description string  `json:"description"` // 用户添加的备注
	CreatedAt   int64   `json:"created_at"`  // 创建时间，Unix 时间戳
	UpdatedAt   int64   `json:"updated_at"`  // 更新时间，Unix 时间戳
	Position    float64 `json:"position"`    // 位置百分比，0.0 - 100.0
}
type BookMarks []BookMark

func (b BookMarks) Len() int           { return len(b) }
func (b BookMarks) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b BookMarks) Less(i, j int) bool { return b[i].CreatedAt < b[j].CreatedAt }

// SortByCreatedAtDesc 按 CreatedAt 降序排序
func (b BookMarks) SortByCreatedAtDesc() {
	for i := 0; i < len(b)-1; i++ {
		for j := 0; j < len(b)-i-1; j++ {
			if b[j].CreatedAt < b[j+1].CreatedAt {
				b.Swap(j, j+1)
			}
		}
	}
}

// FindByPageIndex 寻找指定页码的书签
func (b BookMarks) FindByPageIndex(pageIndex int) *BookMark {
	for i := range b {
		if b[i].PageIndex == pageIndex {
			return &b[i]
		}
	}
	return nil
}

// FindLastBookMark 寻找最新的书签
func (b BookMarks) FindLastBookMark() *BookMark {
	if len(b) == 0 {
		return nil
	}
	last := b[0]
	for i := range b {
		if b[i].CreatedAt > last.CreatedAt {
			last = b[i]
		}
	}
	return &last
}

// FindFarthestBookMark 寻找位置百分比最大的书签
func (b BookMarks) FindFarthestBookMark() *BookMark {
	if len(b) == 0 {
		return nil
	}
	farthest := b[0]
	for i := range b {
		if b[i].Position > farthest.Position {
			farthest = b[i]
		}
	}
	return &farthest
}
