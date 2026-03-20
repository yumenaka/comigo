package file

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/yumenaka/archives"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/encoding"
	"github.com/yumenaka/comigo/tools/logger"
)

// UnArchiveAuto 根据 archives.Identify 自动识别格式并解压（支持 tar / tar.gz、zip、rar 等 Extractor）。
// zipTextEncoding 为空则不对 ZIP 设置文本编码；非空时传给 archives.Zip。
func UnArchiveAuto(filePath string, extractPath string, zipTextEncoding string) error {
	extractPath = tools.GetAbsPath(extractPath)
	if err := os.MkdirAll(extractPath, 0o755); err != nil {
		logger.Infof(locale.GetString("log_failed_to_create_extract_path"), err)
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_open_file_unarchive"), err)
		return err
	}
	defer f.Close()

	format, stream, err := archives.Identify(context.Background(), filePath, f)
	if err != nil {
		if errors.Is(err, archives.NoMatch) {
			msg := fmt.Sprintf(locale.GetString("upgrade_archive_unsupported"), filePath)
			logger.Infof("%s", msg)
			return errors.New(msg)
		}
		logger.Infof(locale.GetString("log_failed_to_identify_file_format"), err)
		return err
	}

	ctx := context.WithValue(context.Background(), "extractPath", extractPath)

	switch z := format.(type) {
	case archives.Zip:
		if zipTextEncoding != "" {
			z.TextEncoding = encoding.ByName(zipTextEncoding)
		}
		err = z.Extract(ctx, stream, extractFileHandler)
	case archives.Rar:
		err = z.Extract(ctx, stream, extractFileHandler)
	case archives.CompressedArchive:
		err = z.Extract(ctx, stream, extractFileHandler)
	case archives.Tar:
		err = z.Extract(ctx, stream, extractFileHandler)
	default:
		if ex, ok := format.(archives.Extractor); ok {
			err = ex.Extract(ctx, stream, extractFileHandler)
		} else {
			msg := fmt.Sprintf(locale.GetString("upgrade_archive_unsupported"), fmt.Sprintf("%T", format))
			logger.Infof("%s", msg)
			return errors.New(msg)
		}
	}
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_extract_zip_file"), err)
		return err
	}
	logger.Infof("%s %s -> %s", locale.GetString("completed_extract"), tools.GetAbsPath(filePath), extractPath)
	return nil
}
