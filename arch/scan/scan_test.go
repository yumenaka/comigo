package scan

import (
	"github.com/yumenaka/comi/types"
	"io/fs"
	"reflect"
	"testing"
)

func TestScanAndGetBookList(t *testing.T) {
	// bookList, err := ScanAndGetBookList("../test")
	// if err != nil {
	// 	t.Errorf("TestScanPath error")
	// }
	// if len(bookList) != 4 {
	// 	t.Errorf("书籍数量不正确" + strconv.Itoa(len(bookList)))
	// }
	// logger.Info(len(bookList))
	// t.Log("hello world")
}

func TestAddBooksToStore(t *testing.T) {
	type args struct {
		bookList    []*types.Book
		basePath    string
		MinImageNum int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddBooksToStore(tt.args.bookList, tt.args.basePath, tt.args.MinImageNum)
		})
	}
}

func TestClearDatabaseWhenExit(t *testing.T) {
	type args struct {
		ConfigPath string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ClearDatabaseWhenExit(tt.args.ConfigPath)
		})
	}
}

func TestNewScanOption(t *testing.T) {
	type args struct {
		reScanFile            bool
		storesPath            []string
		maxScanDepth          int
		minImageNum           int
		timeoutLimitForScan   int
		excludePath           []string
		supportMediaType      []string
		supportFileType       []string
		zipFileTextEncoding   string
		enableDatabase        bool
		clearDatabaseWhenExit bool
		debug                 bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewScanOption(tt.args.reScanFile, tt.args.storesPath, tt.args.maxScanDepth, tt.args.minImageNum, tt.args.timeoutLimitForScan, tt.args.excludePath, tt.args.supportMediaType, tt.args.supportFileType, tt.args.zipFileTextEncoding, tt.args.enableDatabase, tt.args.clearDatabaseWhenExit, tt.args.debug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewScanOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOption_IsSkipDir(t *testing.T) {
	type fields struct {
		ReScanFile            bool
		StoresPath            []string
		MaxScanDepth          int
		MinImageNum           int
		TimeoutLimitForScan   int
		ExcludePath           []string
		SupportMediaType      []string
		SupportFileType       []string
		ZipFileTextEncoding   string
		EnableDatabase        bool
		ClearDatabaseWhenExit bool
		Debug                 bool
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Option{
				ReScanFile:            tt.fields.ReScanFile,
				StoresPath:            tt.fields.StoresPath,
				MaxScanDepth:          tt.fields.MaxScanDepth,
				MinImageNum:           tt.fields.MinImageNum,
				TimeoutLimitForScan:   tt.fields.TimeoutLimitForScan,
				ExcludePath:           tt.fields.ExcludePath,
				SupportMediaType:      tt.fields.SupportMediaType,
				SupportFileType:       tt.fields.SupportFileType,
				ZipFileTextEncoding:   tt.fields.ZipFileTextEncoding,
				EnableDatabase:        tt.fields.EnableDatabase,
				ClearDatabaseWhenExit: tt.fields.ClearDatabaseWhenExit,
				Debug:                 tt.fields.Debug,
			}
			if got := o.IsSkipDir(tt.args.path); got != tt.want {
				t.Errorf("IsSkipDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOption_IsSupportArchiver(t *testing.T) {
	type fields struct {
		ReScanFile            bool
		StoresPath            []string
		MaxScanDepth          int
		MinImageNum           int
		TimeoutLimitForScan   int
		ExcludePath           []string
		SupportMediaType      []string
		SupportFileType       []string
		ZipFileTextEncoding   string
		EnableDatabase        bool
		ClearDatabaseWhenExit bool
		Debug                 bool
	}
	type args struct {
		checkPath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Option{
				ReScanFile:            tt.fields.ReScanFile,
				StoresPath:            tt.fields.StoresPath,
				MaxScanDepth:          tt.fields.MaxScanDepth,
				MinImageNum:           tt.fields.MinImageNum,
				TimeoutLimitForScan:   tt.fields.TimeoutLimitForScan,
				ExcludePath:           tt.fields.ExcludePath,
				SupportMediaType:      tt.fields.SupportMediaType,
				SupportFileType:       tt.fields.SupportFileType,
				ZipFileTextEncoding:   tt.fields.ZipFileTextEncoding,
				EnableDatabase:        tt.fields.EnableDatabase,
				ClearDatabaseWhenExit: tt.fields.ClearDatabaseWhenExit,
				Debug:                 tt.fields.Debug,
			}
			if got := o.IsSupportArchiver(tt.args.checkPath); got != tt.want {
				t.Errorf("IsSupportArchiver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOption_IsSupportMedia(t *testing.T) {
	type fields struct {
		ReScanFile            bool
		StoresPath            []string
		MaxScanDepth          int
		MinImageNum           int
		TimeoutLimitForScan   int
		ExcludePath           []string
		SupportMediaType      []string
		SupportFileType       []string
		ZipFileTextEncoding   string
		EnableDatabase        bool
		ClearDatabaseWhenExit bool
		Debug                 bool
	}
	type args struct {
		checkPath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Option{
				ReScanFile:            tt.fields.ReScanFile,
				StoresPath:            tt.fields.StoresPath,
				MaxScanDepth:          tt.fields.MaxScanDepth,
				MinImageNum:           tt.fields.MinImageNum,
				TimeoutLimitForScan:   tt.fields.TimeoutLimitForScan,
				ExcludePath:           tt.fields.ExcludePath,
				SupportMediaType:      tt.fields.SupportMediaType,
				SupportFileType:       tt.fields.SupportFileType,
				ZipFileTextEncoding:   tt.fields.ZipFileTextEncoding,
				EnableDatabase:        tt.fields.EnableDatabase,
				ClearDatabaseWhenExit: tt.fields.ClearDatabaseWhenExit,
				Debug:                 tt.fields.Debug,
			}
			if got := o.IsSupportMedia(tt.args.checkPath); got != tt.want {
				t.Errorf("IsSupportMedia() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveResultsToDatabase(t *testing.T) {
	type args struct {
		ConfigPath            string
		ClearDatabaseWhenExit bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveResultsToDatabase(tt.args.ConfigPath, tt.args.ClearDatabaseWhenExit); (err != nil) != tt.wantErr {
				t.Errorf("SaveResultsToDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScanAndGetBookList1(t *testing.T) {
	type args struct {
		storePath  string
		scanOption Option
	}
	tests := []struct {
		name            string
		args            args
		wantNewBookList []*types.Book
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewBookList, err := ScanAndGetBookList(tt.args.storePath, tt.args.scanOption)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanAndGetBookList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewBookList, tt.wantNewBookList) {
				t.Errorf("ScanAndGetBookList() gotNewBookList = %v, want %v", gotNewBookList, tt.wantNewBookList)
			}
		})
	}
}

func TestScanStorePath(t *testing.T) {
	type args struct {
		scanConfig Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ScanStorePath(tt.args.scanConfig); (err != nil) != tt.wantErr {
				t.Errorf("ScanStorePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_scanDirGetBook(t *testing.T) {
	type args struct {
		dirPath    string
		storePath  string
		depth      int
		scanOption Option
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Book
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := scanDirGetBook(tt.args.dirPath, tt.args.storePath, tt.args.depth, tt.args.scanOption)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanDirGetBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scanDirGetBook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scanFileGetBook(t *testing.T) {
	type args struct {
		filePath   string
		storePath  string
		depth      int
		scanOption Option
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Book
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := scanFileGetBook(tt.args.filePath, tt.args.storePath, tt.args.depth, tt.args.scanOption)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanFileGetBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scanFileGetBook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scanNonUTF8ZipFile(t *testing.T) {
	type args struct {
		filePath   string
		b          *types.Book
		scanOption Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := scanNonUTF8ZipFile(tt.args.filePath, tt.args.b, tt.args.scanOption); (err != nil) != tt.wantErr {
				t.Errorf("scanNonUTF8ZipFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_walkUTF8ZipFs(t *testing.T) {
	type args struct {
		fsys       fs.FS
		parent     string
		base       string
		b          *types.Book
		scanOption Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := walkUTF8ZipFs(tt.args.fsys, tt.args.parent, tt.args.base, tt.args.b, tt.args.scanOption); (err != nil) != tt.wantErr {
				t.Errorf("walkUTF8ZipFs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
