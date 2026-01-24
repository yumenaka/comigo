package data_api

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// ReadingHistoryResponse 阅读历史响应结构
type ReadingHistoryResponse struct {
	Items      []model.BookinfoWithBookMark `json:"items"`       // 阅读历史列表
	TotalCount int                          `json:"total_count"` // 总数量
	Page       int                          `json:"page"`        // 当前页码
	PageSize   int                          `json:"page_size"`   // 每页数量
	TotalPages int                          `json:"total_pages"` // 总页数
}

// GetReadingHistory 获取阅读历史的 API 处理函数
// 参数:
//   - limit: 限制返回数量 (drawer用，如10)，优先级高于分页参数
//   - page: 页码 (设置页面分页用，默认1)
//   - page_size: 每页数量 (默认20)
func GetReadingHistory(c echo.Context) error {
	// 解析查询参数
	limitStr := c.QueryParam("limit")
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("page_size")

	limit := 0 // 默认不限制
	if limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	page := 1 // 默认第一页
	if pageStr != "" {
		if val, err := strconv.Atoi(pageStr); err == nil && val > 0 {
			page = val
		}
	}

	pageSize := 20 // 默认每页20条
	if pageSizeStr != "" {
		if val, err := strconv.Atoi(pageSizeStr); err == nil && val > 0 {
			pageSize = val
		}
	}

	// 获取所有书籍
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}

	// 收集阅读历史（过滤 page_index > 0 的书签）
	readingHistory := []model.BookinfoWithBookMark{}
	for _, book := range allBooks {
		for _, mark := range book.BookMarks {
			// 只收集有阅读记录的书签（page_index > 0）
			if mark.PageIndex > 0 {
				book.BookInfo.Cover = book.GetCover()
				bookinfoWithMark := model.BookinfoWithBookMark{
					BookInfo: book.BookInfo,
					BookMark: mark,
				}
				readingHistory = append(readingHistory, bookinfoWithMark)
			}
		}
	}

	// 排序规则：手动书签（user）优先于自动书签（auto），同类型按更新时间降序
	sort.Slice(readingHistory, func(i, j int) bool {
		// 手动书签优先
		if readingHistory[i].BookMark.Type != readingHistory[j].BookMark.Type {
			return readingHistory[i].BookMark.Type == model.UserMark
		}
		// 同类型按更新时间降序
		return readingHistory[i].BookMark.UpdatedAt.After(readingHistory[j].BookMark.UpdatedAt)
	})

	totalCount := len(readingHistory)

	// 如果设置了 limit 参数，直接返回前 N 条（用于 drawer）
	if limit > 0 {
		if limit > totalCount {
			limit = totalCount
		}
		return c.JSON(http.StatusOK, ReadingHistoryResponse{
			Items:      readingHistory[:limit],
			TotalCount: totalCount,
			Page:       1,
			PageSize:   limit,
			TotalPages: 1,
		})
	}

	// 否则进行分页处理（用于设置页面）
	totalPages := (totalCount + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	// 确保页码有效
	if page > totalPages {
		page = totalPages
	}

	// 计算分页范围
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > totalCount {
		end = totalCount
	}

	// 获取当前页的数据
	var items []model.BookinfoWithBookMark
	if start < totalCount {
		items = readingHistory[start:end]
	} else {
		items = []model.BookinfoWithBookMark{}
	}

	return c.JSON(http.StatusOK, ReadingHistoryResponse{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}
