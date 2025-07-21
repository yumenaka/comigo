package repository

import (
	"errors"

	"github.com/yumenaka/comigo/model"
)

// GetAllBookInfoList 获取所有 BookInfo，并根据 sortBy 参数进行排序
func GetAllBookInfoList(sortBy string) (*model.BookInfoList, error) {
	var infoList model.BookInfoList
	// 添加所有真实的书籍
	for _, value := range MainStore.mapBooks.Range {
		b := value.(*model.Book)
		info := model.NewBaseInfo(b)
		infoList.BookInfos = append(infoList.BookInfos, *info)
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error: cannot find bookshelf in GetAllBookInfoList")
}
