package scan

import (
	"errors"
	"strconv"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// 处理 PDF 文件
func handlePdfFiles(filePath string, newBook *model.Book) error {
	pageCount, err := file.CountPagesOfPDF(filePath)
	if err != nil {
		return err
	}
	if pageCount < 1 {
		return errors.New(locale.GetString("no_pages_in_pdf") + filePath)
	}
	logger.Infof(locale.GetString("scan_pdf")+" %s: %d pages", filePath, pageCount)
	newBook.PageCount = pageCount
	newBook.InitComplete = true
	for i := 1; i <= pageCount; i++ {
		tempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + strconv.Itoa(i) + ".jpg"
		newBook.PageInfos = append(newBook.PageInfos, model.PageInfo{
			Name:    strconv.Itoa(i),
			Url:     tempURL,
			PageNum: i,
		})
	}
	return nil
}
