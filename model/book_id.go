package model

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// BookIdentity 区分书库、路径和类型；远程 Comigo 使用公开书库标识及原始 ID。
// 文件大小不用于入库时的身份核对，同一本书可以更新内容。
func (b *BookInfo) BookIdentity() [3]string {
	if b.RemoteBookID != "" {
		return [3]string{"comigo", b.RemoteStoreKey, b.RemoteBookID}
	}
	return [3]string{b.StoreUrl, b.BookPath, string(b.Type)}
}

// fullBookID 保留原有摘要字段；文件大小变化仍会更换 ID，使图片和解压缓存失效。
func (b *BookInfo) fullBookID() string {
	fields := []string{b.BookPath, strconv.FormatInt(b.FileSize, 10), string(b.Type), b.ParentFolder, b.StoreUrl}
	if b.RemoteBookID != "" {
		identity := b.BookIdentity()
		fields = identity[:]
	}
	// 路径不含 NUL，用它分隔字段并保留文件名的原始字节。
	digest := sha256.Sum256([]byte(strings.Join(fields, "\x00")))
	return base62.EncodeToString(digest[:])
}

// CheckBookID 拒绝不同书籍共用 ID；更新同一本书时允许覆盖旧内容。
func CheckBookID(existing, incoming *BookInfo) error {
	if existing.BookID == incoming.BookID && existing.BookIdentity() != incoming.BookIdentity() {
		return fmt.Errorf(locale.GetString("err_book_id_conflict"), incoming.BookID)
	}
	return nil
}

// BookIDIndex 在一次扫描或书组重建期间记录已入库和新分配的 ID。
// 由调用方串行使用，不保存全局状态，也不等待书籍入库后才登记占用。
type BookIDIndex struct {
	byID     map[string][3]string
	byDigest map[string]string
}

// NewBookIDIndex 从全库快照建立占用表，旧版 ID 只占位，不再用于新书索引。
func NewBookIDIndex(books []*Book) (*BookIDIndex, error) {
	ids := &BookIDIndex{byID: make(map[string][3]string), byDigest: make(map[string]string)}
	for _, book := range books {
		if err := ids.Reserve(&book.BookInfo); err != nil {
			return nil, err
		}
		if !book.Deleted && config.SupportsBookVersion(book.CreatedByVersion) {
			ids.byDigest[book.fullBookID()] = book.BookID
		}
	}
	return ids, nil
}

// Assign 将 SHA-256 原始摘要直接编码为 Base62，从 7 位开始，冲突时逐位加长。
func (ids *BookIDIndex) Assign(book *BookInfo) error {
	fullID := book.fullBookID()
	identity := book.BookIdentity()
	if id, exists := ids.byDigest[fullID]; exists && ids.byID[id] == identity {
		book.BookID = id
		return nil
	}
	for length := 7; length <= len(fullID); length++ {
		id := fullID[:length]
		if owner, exists := ids.byID[id]; exists && owner != identity {
			continue
		}
		book.BookID = id
		ids.byID[id] = identity
		ids.byDigest[fullID] = id
		return nil
	}
	return errors.New(locale.GetString("err_book_id_exhausted"))
}

// Reserve 只登记仍占用的 ID，不把待恢复或软删除记录加入摘要复用表。
func (ids *BookIDIndex) Reserve(book *BookInfo) error {
	identity := book.BookIdentity()
	if owner, exists := ids.byID[book.BookID]; exists && owner != identity {
		return fmt.Errorf(locale.GetString("err_book_id_conflict"), book.BookID)
	}
	ids.byID[book.BookID] = identity
	return nil
}

// NewStoreBookIDIndex 补齐书架列表不可见但尚未释放的 ID，供扫描和书组重建共用。
func NewStoreBookIDIndex(s StoreInterface, books []*Book) (*BookIDIndex, error) {
	ids, err := NewBookIDIndex(books)
	if err != nil {
		return nil, err
	}
	if reservations, ok := s.(interface{ ReserveBookIDs(*BookIDIndex) error }); ok {
		if err := reservations.ReserveBookIDs(ids); err != nil {
			return nil, err
		}
	}
	return ids, nil
}
