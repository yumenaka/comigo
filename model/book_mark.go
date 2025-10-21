package model

import "time"

type BookMark struct {
	BookID      string    `json:"book_id"`     // 书籍 ID
	PageIndex   int       `json:"page_index"`  // 书签页码，从 0 开始，不会超过 PageCount - 1
	Description string    `json:"description"` // 用户添加的备注
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
	Position    float64   `json:"position"`    // 位置百分比，0.0 - 100.0
}
type BookMarks []BookMark

func NewBookMark(bookID string, pageIndex int, pageCount int, description string) *BookMark {
	createdAt := time.Now()
	position := (float64(pageIndex) / float64(pageCount)) * 100.0
	return &BookMark{
		BookID:      bookID,
		PageIndex:   pageIndex,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
		Position:    position,
	}
}

// go
func (b *BookMarks) Add(bookMark BookMark) {
	*b = append(*b, bookMark)
	b.SortByCreatedAtDesc()
}

func (b *BookMarks) Len() int { return len(*b) }

func (b *BookMarks) Swap(i, j int) {
	(*b)[i], (*b)[j] = (*b)[j], (*b)[i]
}

func (b *BookMarks) Less(i, j int) bool {
	return (*b)[i].CreatedAt.Before((*b)[j].CreatedAt)
}

// SortByCreatedAtDesc 按 CreatedAt 降序排序
func (b *BookMarks) SortByCreatedAtDesc() {
	n := len(*b)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if (*b)[j].CreatedAt.Before((*b)[j+1].CreatedAt) {
				b.Swap(j, j+1)
			}
		}
	}
}

// FindByPageIndex 寻找指定页码的书签
func (b *BookMarks) FindByPageIndex(pageIndex int) *BookMark {
	for i := range *b {
		if (*b)[i].PageIndex == pageIndex {
			return &(*b)[i]
		}
	}
	return nil
}

// FindLastBookMark 寻找最新的书签
func (b *BookMarks) FindLastBookMark() *BookMark {
	if len(*b) == 0 {
		return nil
	}
	lastIdx := 0
	for i := range *b {
		if (*b)[i].CreatedAt.Before((*b)[lastIdx].CreatedAt) {
			lastIdx = i
		}
	}
	return &(*b)[lastIdx]
}

// FindFarthestBookMark 寻找位置百分比最大的书签
func (b *BookMarks) FindFarthestBookMark() *BookMark {
	if len(*b) == 0 {
		return nil
	}
	farthestIdx := 0
	for i := range *b {
		if (*b)[i].Position > (*b)[farthestIdx].Position {
			farthestIdx = i
		}
	}
	return &(*b)[farthestIdx]
}
