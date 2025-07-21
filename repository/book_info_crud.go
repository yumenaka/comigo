package repository

import (
	"errors"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// GetCover 获取封面
func GetCover(b *model.BookInfo) model.MediaFileInfo {
	switch b.Type {
	// 书籍类型为书组的时候，遍历所有子书籍，然后获取第一个子书籍的封面
	case model.TypeBooksGroup:
		bookGroup, err := GetBookGroupByBookID(b.BookID)
		if err != nil {
			logger.Infof("Error getting book group: %s", err)
			return model.MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
		}
		for _, v := range bookGroup.ChildBook.Range {
			b := v.(*model.BookInfo)
			childBook, err := GetBookByID(b.BookID, "modify_time")
			if err != nil {
				return model.MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
			}
			// 递归调用
			return GetCover(&childBook.BookInfo)
		}
	case model.TypeDir, model.TypeZip, model.TypeRar, model.TypeCbz, model.TypeCbr, model.TypeTar, model.TypeEpub:
		tempBook, err := GetBookByID(b.BookID, "")
		if err != nil || len(tempBook.Pages.Images) == 0 {
			return model.MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
		}
		return tempBook.GuestCover()
	case model.TypePDF:
		return model.MediaFileInfo{Name: "1.jpg", Url: "/api/get_file?id=" + b.BookID + "&filename=" + "1.jpg"}
	case model.TypeVideo:
		return model.MediaFileInfo{Name: "video.png", Url: "/images/video.png"}
	case model.TypeAudio:
		return model.MediaFileInfo{Name: "audio.png", Url: "/images/audio.png"}
	case model.TypeUnknownFile:
		return model.MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
	}
	return model.MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
}

// GetBookInfoListByDepth 根据深度获取书籍列表
func GetBookInfoListByDepth(depth int, sortBy string) (*model.BookInfoList, error) {
	var infoList model.BookInfoList
	// 首先加上所有真实的书籍
	for _, bookValue := range MainStore.mapBooks.Range {
		b := bookValue.(*model.Book)
		if b.Depth == depth {
			info := model.NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	// 接下来还要加上扫描生成出来的书籍组
	for _, folderValue := range MainStore.SubStores.Range {
		bs := folderValue.(*subMemoryStore)
		for _, groupValue := range bs.BookGroupMap.Range {
			group := groupValue.(*model.BookInfo)
			if group.Depth == depth {
				infoList.BookInfos = append(infoList.BookInfos, *group)
			}
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error: cannot find bookshelf in GetBookInfoListByDepth")
}

// GetBookInfoListByMaxDepth 获取指定最大深度的书籍列表
func GetBookInfoListByMaxDepth(depth int, sortBy string) (*model.BookInfoList, error) {
	var infoList model.BookInfoList
	// 首先加上所有真实的书籍
	for _, bookValue := range MainStore.mapBooks.Range {
		b := bookValue.(*model.Book)
		if b.Depth <= depth {
			info := model.NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	// 扫描生成的书籍组
	for _, folderValue := range MainStore.SubStores.Range {
		bs := folderValue.(*subMemoryStore)
		for _, groupValue := range bs.BookGroupMap.Range {
			group := groupValue.(*model.BookInfo)
			if group.Depth <= depth {
				infoList.BookInfos = append(infoList.BookInfos, *group)
			}
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error: cannot find bookshelf in GetBookInfoListByMaxDepth")
}

// TopOfShelfInfo 获取顶层书架信息
func TopOfShelfInfo(sortBy string) (*model.BookInfoList, error) {
	// if len(*StoreSettings) == 0 {
	//	return nil, errors.New("error: cannot find book in TopOfShelfInfo")
	// }
	// if len(*StoreSettings) > 1 {
	//	// 有多个书库
	//	var infoList model.BookInfoList
	//	for _, localPath := range *StoreSettings {
	//		for _, groupValue := range mapBookGroup.Range {
	//			group := groupValue.(*BookGroup)
	//			if group.BookInfo.ParentFolder == localPath {
	//				infoList.BookInfos = append(infoList.BookInfos, group.BookInfo)
	//			}
	//		}
	//	}
	//	if len(infoList.BookInfos) > 0 {
	//		infoList.SortBooks(sortBy)
	//		return &infoList, nil
	//	}
	//	return nil, errors.New("error: cannot find book in TopOfShelfInfo")
	// }
	// 显示顶层书库的书籍
	var infoList model.BookInfoList
	for _, bookValue := range MainStore.mapBooks.Range {
		b := bookValue.(*model.Book)
		if b.Depth == 0 {
			info := model.NewBaseInfo(b)
			info.Cover = GetCover(info) // 设置封面图(为了兼容老版前端)TODO：升级新前端，去掉这部分
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	for _, groupValue := range MainStore.mapBookGroup.Range {
		group := groupValue.(*model.BookGroup)
		if group.BookInfo.Depth == 0 {
			group.BookInfo.Cover = GetCover(&group.BookInfo) // 设置封面图(为了兼容老版前端)TODO：升级新前端，去掉这部分
			infoList.BookInfos = append(infoList.BookInfos, group.BookInfo)
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	// 没找到任何书
	return nil, errors.New("error: cannot find book in TopOfShelfInfo")
}

// GetBookInfoListByID 根据 ID 获取书籍列表
func GetBookInfoListByID(BookID string, sortBy string) (*model.BookInfoList, error) {
	var infoList model.BookInfoList
	groupValue, ok := MainStore.mapBookGroup.Load(BookID)
	if ok {
		tempGroup := groupValue.(*model.BookGroup)
		for _, bookValue := range tempGroup.ChildBook.Range {
			b := bookValue.(*model.BookInfo)
			b.Cover = GetCover(b) // 设置封面图(为了兼容老版前端) TODO：升级前端，去掉这部分
			infoList.BookInfos = append(infoList.BookInfos, *b)
		}
		if len(infoList.BookInfos) > 0 {
			infoList.SortBooks(sortBy)
			return &infoList, nil
		}
	}
	return nil, errors.New("cannot find BookInfo，ID：" + BookID)
}

// GetBookInfoListByParentFolder 根据父文件夹获取书籍列表
func GetBookInfoListByParentFolder(parentFolder string, sortBy string) (*model.BookInfoList, error) {
	var infoList model.BookInfoList
	for _, bookValue := range MainStore.mapBooks.Range {
		b := bookValue.(*model.Book)
		if b.ParentFolder == parentFolder {
			info := model.NewBaseInfo(b)
			info.Cover = GetCover(info) // 设置封面图(为了兼容老版前端) TODO：升级前端，去掉这部分
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("cannot find book, parentFolder=" + parentFolder)
}
