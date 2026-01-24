package model

import (
	"time"
)

type MarkType string

const (
	AutoMark MarkType = "auto"
	UserMark MarkType = "user"
)

type BookMark struct {
	Type        MarkType  `json:"mark_type"`     // 书签类型，auto 表示自动书签，user 表示用户书签
	BookID      string    `json:"book_id"`       // 书籍 ID
	BookStoreID string    `json:"book_store_id"` // 书籍所属书库 ID
	PageIndex   int       `json:"page_index"`    // 书签页码，从 0 开始，理论上不会超过 PageCount - 1 ，但是现实中可能会有
	Description string    `json:"description"`   // 用户添加的备注
	CreatedAt   time.Time `json:"created_at"`    // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`    // 更新时间
}

type BookinfoWithBookMark struct {
	BookInfo BookInfo `json:"book_info"` // 书籍信息
	BookMark BookMark `json:"book_mark"` // 书签信息
}

func NewBookMark(markType MarkType, bookID string, bookStoreID string, pageIndex int, description string) *BookMark {
	createdAt := time.Now()
	return &BookMark{
		Type:        markType,
		BookID:      bookID,
		BookStoreID: bookStoreID,
		PageIndex:   pageIndex,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
}
