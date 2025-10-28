package service

//
// type BookRepositoryInterface interface {
// 	GetBook(ctx context.Context, bookID string) (*model.Book, error)
// 	GetBookByFilePath(ctx context.Context, filePath string) (*model.Book, error)
// 	ListBooks(ctx context.Context) ([]*model.Book, error)
// 	ListBooksByType(ctx context.Context, bookType string) ([]*model.Book, error)
// 	ListBooksByStorePath(ctx context.Context, storePath string) ([]*model.Book, error)
// 	SearchBooksByTitle(ctx context.Context, title string) ([]*model.Book, error)
// 	Create(ctx context.Context, book *model.Book) error
// 	Update(ctx context.Context, book *model.Book) error
// 	UpdateReadPercent(ctx context.Context, bookID string, readPercent float64) error
// 	MarkAsDeleted(ctx context.Context, bookID string) error
// 	Delete(ctx context.Context, bookID string) error
// 	GetMediaFiles(ctx context.Context, bookID string) ([]model.MediaFileInfo, error)
// 	GetMediaFile(ctx context.Context, bookID string, pageNum int) (*model.MediaFileInfo, error)
// 	GetCover(ctx context.Context, bookID string) (*model.MediaFileInfo, error)
// 	CreateMediaFile(ctx context.Context, mediaFile model.MediaFileInfo, bookID string) error
// 	UpdateMediaFile(ctx context.Context, mediaFile model.MediaFileInfo, bookID string) error
// 	DeleteMediaFiles(ctx context.Context, bookID string) error
// 	Count(ctx context.Context) (int64, error)
// 	CountByType(ctx context.Context, bookType string) (int64, error)
// 	CountMediaFiles(ctx context.Context, bookID string) (int64, error)
// 	GetTotalFileSize(ctx context.Context) (float64, error)
// }
//
// type StoreRepositoryInterface interface {
// 	GetByID(ctx context.Context, id int) (*model.StoreInfo, error)
// 	GetByName(ctx context.Context, name string) (*model.StoreInfo, error)
// 	List(ctx context.Context) ([]*model.StoreInfo, error)
// 	GetWithBackend(ctx context.Context, id int) (*model.StoreInfo, error)
// 	ListWithBackend(ctx context.Context) ([]*model.StoreInfo, error)
// 	Create(ctx context.Context, store *model.StoreInfo) error
// 	Update(ctx context.Context, store *model.StoreInfo) error
// 	Delete(ctx context.Context, id int) error
// 	Count(ctx context.Context) (int64, error)
// 	CountByType(ctx context.Context, storeType string) (int64, error)
// 	GetTotalFileSize(ctx context.Context) (float64, error)
// }
//
// // BookService 书籍业务逻辑层
// type BookService struct {
// 	bookRepo  BookRepositoryInterface
// 	storeRepo StoreRepositoryInterface
// }
//
// // NewBookService 创建新的BookService实例，bookRepo接受 *repository.BookRepository
// func NewBookService(bookRepo BookRepositoryInterface) *BookService {
// 	return &BookService{
// 		bookRepo: bookRepo,
// 	}
// }
//
// // GetBook 获取书籍信息（包含业务逻辑）
// func (s *BookService) GetBook(ctx context.Context, bookID string) (*model.Book, error) {
// 	// 1. 从数据库获取书籍基本信息
// 	book, err := s.bookRepo.GetBook(ctx, bookID)
// 	if err != nil {
// 		return nil, fmt.Errorf("获取书籍失败: %w", err)
// 	}
//
// 	// 2. 获取书籍的媒体文件信息
// 	mediaFiles, err := s.bookRepo.GetMediaFiles(ctx, bookID)
// 	if err != nil {
// 		return nil, fmt.Errorf("获取媒体文件失败: %w", err)
// 	}
//
// 	// 3. 设置书籍的页面信息
// 	book.PageInfos = mediaFiles
//
// 	// 4. 业务逻辑：设置封面
// 	if len(mediaFiles) > 0 {
// 		book.Cover = book.GuestCover()
// 	}
//
// 	return book, nil
// }
//
// // CreateBook 创建新书籍（包含业务逻辑）
// func (s *BookService) CreateBook(ctx context.Context, book *model.Book) error {
// 	// 1. 业务验证
// 	if book.Title == "" {
// 		return fmt.Errorf("书籍标题不能为空")
// 	}
// 	if book.BookID == "" {
// 		return fmt.Errorf("书籍ID不能为空")
// 	}
//
// 	// 2. 检查书籍是否已存在
// 	existingBook, err := s.bookRepo.GetBook(ctx, book.BookID)
// 	if err == nil && existingBook != nil {
// 		return fmt.Errorf("书籍已存在: %s", book.BookID)
// 	}
//
// 	// 3. 保存书籍基本信息
// 	err = s.bookRepo.Create(ctx, book)
// 	if err != nil {
// 		return fmt.Errorf("创建书籍失败: %w", err)
// 	}
//
// 	// 4. 保存媒体文件信息
// 	for _, mediaFile := range book.PageInfos {
// 		err = s.bookRepo.CreateMediaFile(ctx, mediaFile, book.BookID)
// 		if err != nil {
// 			return fmt.Errorf("创建媒体文件记录失败: %w", err)
// 		}
// 	}
//
// 	return nil
// }
//
// // UpdateReadProgress 更新阅读进度
// func (s *BookService) UpdateReadProgress(ctx context.Context, bookID string, readPercent float64) error {
// 	// 1. 业务验证
// 	if readPercent < 0 || readPercent > 100 {
// 		return fmt.Errorf("阅读进度必须在0-100之间")
// 	}
//
// 	// 2. 检查书籍是否存在
// 	_, err := s.bookRepo.GetBook(ctx, bookID)
// 	if err != nil {
// 		return fmt.Errorf("书籍不存在: %w", err)
// 	}
//
// 	// 3. 业务逻辑：如果阅读进度超过90%，标记为已完成
// 	if readPercent >= 90 {
// 		// 这里可以添加额外的业务逻辑，比如发送通知等
// 	}
//
// 	// 4. 更新阅读进度
// 	err = s.bookRepo.UpdateReadPercent(ctx, bookID, readPercent)
// 	if err != nil {
// 		return fmt.Errorf("更新阅读进度失败: %w", err)
// 	}
//
// 	return nil
// }
//
// // SearchBooks 搜索书籍
// func (s *BookService) SearchBooks(ctx context.Context, title string) ([]*model.Book, error) {
// 	// 1. 业务验证
// 	if title == "" {
// 		return nil, fmt.Errorf("搜索标题不能为空")
// 	}
//
// 	// 2. 执行搜索
// 	books, err := s.bookRepo.SearchBooksByTitle(ctx, title)
// 	if err != nil {
// 		return nil, fmt.Errorf("搜索书籍失败: %w", err)
// 	}
//
// 	// 3. 业务逻辑：为每个书籍设置封面
// 	for _, book := range books {
// 		mediaFiles, err := s.bookRepo.GetMediaFiles(ctx, book.BookID)
// 		if err == nil && len(mediaFiles) > 0 {
// 			book.PageInfos = mediaFiles
// 			book.Cover = book.GuestCover()
// 		}
// 	}
//
// 	return books, nil
// }
//
// // GetBookStatistics 获取书籍统计信息
// func (s *BookService) GetBookStatistics(ctx context.Context) (*BookStatistics, error) {
// 	// 1. 获取总书籍数
// 	totalBooks, err := s.bookRepo.Count(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("获取书籍总数失败: %w", err)
// 	}
//
// 	// 2. 获取总文件大小
// 	totalFileSize, err := s.bookRepo.GetTotalFileSize(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("获取总文件大小失败: %w", err)
// 	}
//
// 	// 3. 获取各类型书籍数量
// 	zipBooks, err := s.bookRepo.CountByType(ctx, "zip")
// 	if err != nil {
// 		return nil, fmt.Errorf("获取ZIP书籍数量失败: %w", err)
// 	}
//
// 	pdfBooks, err := s.bookRepo.CountByType(ctx, "pdf")
// 	if err != nil {
// 		return nil, fmt.Errorf("获取PDF书籍数量失败: %w", err)
// 	}
//
// 	return &BookStatistics{
// 		TotalBooks:    totalBooks,
// 		TotalFileSize: totalFileSize,
// 		ZipBooks:      zipBooks,
// 		PdfBooks:      pdfBooks,
// 	}, nil
// }
//
// // BookStatistics 书籍统计信息
// type BookStatistics struct {
// 	TotalBooks    int64   `json:"total_books"`
// 	TotalFileSize float64 `json:"total_file_size"`
// 	ZipBooks      int64   `json:"zip_books"`
// 	PdfBooks      int64   `json:"pdf_books"`
// }
