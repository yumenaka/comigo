-- Book related queries

-- Get single book by ID
-- name: GetBookByID :one
SELECT * FROM books 
WHERE book_id = ? LIMIT 1;

-- Get book by file path
-- name: GetBookByFilePath :one
SELECT * FROM books 
WHERE file_path = ? LIMIT 1;

-- List all books
-- name: ListBooks :many
SELECT * FROM books 
ORDER BY modified_time DESC;

-- List books by type
-- name: ListBooksByType :many
SELECT * FROM books 
WHERE type = ? 
ORDER BY modified_time DESC;

-- List books by store path
-- name: ListBooksByStorePath :many
SELECT * FROM books 
WHERE book_store_path = ? 
ORDER BY modified_time DESC;

-- Search books by title (fuzzy search)
-- name: SearchBooksByTitle :many
SELECT * FROM books 
WHERE title LIKE '%' || ? || '%' 
ORDER BY modified_time DESC;

-- Create new book
-- name: CreateBook :one
INSERT INTO books (
    title, book_id, owner, file_path, book_store_path, type,
    child_books_num, child_books_id,depth, parent_folder, page_count, file_size,
    author, isbn, press, published_at, extract_path, extract_num,
    init_complete, read_percent, non_utf8zip, zip_text_encoding
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- Update book information
-- name: UpdateBook :exec
UPDATE books SET
    title = ?, owner = ?, file_path = ?, book_store_path = ?, type = ?,
    child_books_num = ?, child_books_id = ?, depth = ?, parent_folder = ?, page_count = ?, file_size = ?,
    author = ?, isbn = ?, press = ?, published_at = ?, extract_path = ?, extract_num = ?,
    init_complete = ?, read_percent = ?, non_utf8zip = ?, zip_text_encoding = ?,
    modified_time = CURRENT_TIMESTAMP
WHERE book_id = ?;

-- Update reading progress
-- name: UpdateReadPercent :exec
UPDATE books SET
    read_percent = ?,
    modified_time = CURRENT_TIMESTAMP
WHERE book_id = ?;

-- Mark book as deleted (soft delete)
-- name: MarkBookAsDeleted :exec
UPDATE books SET
    deleted = TRUE,
    modified_time = CURRENT_TIMESTAMP
WHERE book_id = ?;

-- Delete book
-- name: DeleteBook :exec
DELETE FROM books WHERE book_id = ?;

-- Media files related queries

-- Get all page information by book ID
-- name: GetMediaFilesByBookID :many
SELECT * FROM media_files 
WHERE book_id = ? 
ORDER BY page_num ASC;

-- Get specific page by book ID and page number
-- name: GetMediaFileByBookIDAndPage :one
SELECT * FROM media_files 
WHERE book_id = ? AND page_num = ? 
LIMIT 1;

-- Get book cover (usually page 0 or 1)
-- name: GetBookCover :one
SELECT * FROM media_files 
WHERE book_id = ? AND (page_num = 0 OR page_num = 1) 
LIMIT 1;

-- Create media file record
-- name: CreateMediaFile :one
INSERT INTO media_files (
    book_id, name, path, size, mod_time, url, page_num,
    blurhash, height, width, img_type, insert_html
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- Update media file information
-- name: UpdateMediaFile :exec
UPDATE media_files SET
    name = ?, path = ?, size = ?, mod_time = ?, url = ?,
    blurhash = ?, height = ?, width = ?, img_type = ?, insert_html = ?
WHERE book_id = ? AND page_num = ?;

-- Delete all media files for a book
-- name: DeleteMediaFilesByBookID :exec
DELETE FROM media_files WHERE book_id = ?;

-- File backend related queries

-- Get file backend by url
-- name: GetFileBackendByID :one
SELECT * FROM file_backends 
WHERE url = ? LIMIT 1;

-- List all file backends
-- name: ListFileBackends :many
SELECT * FROM file_backends 
ORDER BY created_at DESC;

-- List file backends by type
-- name: ListFileBackendsByType :many
SELECT * FROM file_backends 
WHERE type = ? 
ORDER BY created_at DESC;

-- Create file backend
-- name: CreateFileBackend :one
INSERT INTO file_backends (
    type, url, server_host, server_port, need_auth, auth_username,
    auth_password, smb_share_name, smb_path
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- Update file backend
-- name: UpdateFileBackend :exec
UPDATE file_backends SET
    type = ?, url = ?, server_host = ?, server_port = ?, need_auth = ?,
    auth_username = ?, auth_password = ?, smb_share_name = ?, smb_path = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE url = ?;

-- Delete file backend
-- name: DeleteFileBackend :exec
DELETE FROM file_backends WHERE url = ?;

-- Store related queries

-- Get store by URL
-- name: GetStoreByURL :one
SELECT * FROM stores 
WHERE url = ? LIMIT 1;

-- Get store by name
-- name: GetStoreByName :one
SELECT * FROM stores 
WHERE name = ? LIMIT 1;

-- List all stores
-- name: ListStores :many
SELECT * FROM stores 
ORDER BY created_at DESC;

-- Create store
-- name: CreateStore :one
INSERT INTO stores (
    name, description, url
) VALUES (
    ?, ?, ?
) RETURNING *;

-- Update store
-- name: UpdateStore :exec
UPDATE stores SET
    name = ?, description = ?, url = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE url = ?;

-- Delete store
-- name: DeleteStore :exec
DELETE FROM stores WHERE url = ?;

-- Get store with file backend information
-- name: GetStoreWithBackend :one
SELECT 
    s.url, s.name, s.description, s.created_at, s.updated_at,
    fb.type, fb.url, fb.server_host, fb.server_port,
    fb.need_auth, fb.auth_username, fb.auth_password, fb.smb_share_name, fb.smb_path
FROM stores s
JOIN file_backends fb ON s.url = fb.url
WHERE s.url = ? LIMIT 1;

-- List stores with file backend information
-- name: ListStoresWithBackend :many
SELECT 
    s.url, s.name, s.description, s.created_at, s.updated_at,
    fb.type, fb.url, fb.server_host, fb.server_port,
    fb.need_auth, fb.auth_username, fb.auth_password, fb.smb_share_name, fb.smb_path
FROM stores s
JOIN file_backends fb ON s.url = fb.url
ORDER BY s.created_at DESC;

-- Statistics queries

-- Count total books
-- name: CountBooks :one
SELECT COUNT(*) FROM books WHERE deleted = FALSE;

-- Count books by type
-- name: CountBooksByType :one
SELECT COUNT(*) FROM books WHERE type = ? AND deleted = FALSE;

-- Count media files for a book
-- name: CountMediaFilesByBookID :one
SELECT COUNT(*) FROM media_files WHERE book_id = ?;

-- Count total stores
-- name: CountStores :one
SELECT COUNT(*) FROM stores;

-- Count file backends by type
-- name: CountFileBackendsByType :one
SELECT COUNT(*) FROM file_backends WHERE type = ?;

-- Get total file size
-- name: GetTotalFileSize :one
SELECT SUM(file_size) FROM books WHERE deleted = FALSE;