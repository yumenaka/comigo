-- Book related queries

-- Get single book by ID
-- name: GetBookByID :one
SELECT *
FROM books
WHERE book_id = $1
LIMIT 1;

-- Get book by file path
-- name: GetBookByBookPath :one
SELECT *
FROM books
WHERE book_path = $1
LIMIT 1;

-- List all books
-- name: ListBooks :many
SELECT *
FROM books
WHERE deleted = FALSE
ORDER BY modified_time DESC;

-- List books by type
-- name: ListBooksByType :many
SELECT *
FROM books
WHERE type = $1
  AND deleted = FALSE
ORDER BY modified_time DESC;

-- get all store_url for books
-- name: ListAllBookStoreURLs :many
SELECT DISTINCT store_url
FROM books
WHERE deleted = FALSE;

-- List books by store path
-- name: ListBooksByStorePath :many
SELECT *
FROM books
WHERE store_url = $1
  AND deleted = FALSE
ORDER BY modified_time DESC;

-- Search books by title (fuzzy search)
-- name: SearchBooksByTitle :many
SELECT *
FROM books
WHERE title ILIKE '%' || $1 || '%'
  AND deleted = FALSE
ORDER BY modified_time DESC;

-- Insert or update the complete book metadata atomically.
-- name: UpsertBook :exec
INSERT INTO books (title, book_id, owner, book_path, store_url, type,
                   child_books_num, child_books_id, depth, parent_folder, page_count, last_read_page, file_size,
                   author, isbn, press, published_at, extract_path, extract_num, book_complete,
                   init_complete, non_utf8zip, zip_text_encoding, created_by_version, is_remote, remote_url, modified_time, deleted,
                   remote_book_id, remote_store_key, remote_shelf_key, remote_shelf_name,
                   cover_name, cover_path, cover_size, cover_mod_time, cover_url, cover_page_num,
                   cover_blurhash, cover_height, cover_width, cover_img_type, cover_insert_html)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43)
ON CONFLICT (book_id) DO UPDATE SET
    title = excluded.title,
    owner = excluded.owner,
    book_path = excluded.book_path,
    store_url = excluded.store_url,
    type = excluded.type,
    child_books_num = excluded.child_books_num,
    child_books_id = excluded.child_books_id,
    depth = excluded.depth,
    parent_folder = excluded.parent_folder,
    page_count = excluded.page_count,
    last_read_page = excluded.last_read_page,
    file_size = excluded.file_size,
    author = excluded.author,
    isbn = excluded.isbn,
    press = excluded.press,
    published_at = excluded.published_at,
    extract_path = excluded.extract_path,
    extract_num = excluded.extract_num,
    book_complete = excluded.book_complete,
    init_complete = excluded.init_complete,
    non_utf8zip = excluded.non_utf8zip,
    zip_text_encoding = excluded.zip_text_encoding,
    created_by_version = excluded.created_by_version,
    is_remote = excluded.is_remote,
    remote_url = excluded.remote_url,
    modified_time = excluded.modified_time,
    deleted = excluded.deleted,
    remote_book_id = excluded.remote_book_id,
    remote_store_key = excluded.remote_store_key,
    remote_shelf_key = excluded.remote_shelf_key,
    remote_shelf_name = excluded.remote_shelf_name,
    cover_name = excluded.cover_name,
    cover_path = excluded.cover_path,
    cover_size = excluded.cover_size,
    cover_mod_time = excluded.cover_mod_time,
    cover_url = excluded.cover_url,
    cover_page_num = excluded.cover_page_num,
    cover_blurhash = excluded.cover_blurhash,
    cover_height = excluded.cover_height,
    cover_width = excluded.cover_width,
    cover_img_type = excluded.cover_img_type,
    cover_insert_html = excluded.cover_insert_html;

-- Update reading progress
-- name: UpdateLastReadPage :exec
UPDATE bookmarks
SET page_index = $1,
    type = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE book_id = $3;

-- Mark book as deleted (soft delete)
-- name: MarkBookAsDeleted :exec
UPDATE books
SET deleted = TRUE,
    modified_time = CURRENT_TIMESTAMP
WHERE book_id = $1;

-- Delete book
-- name: DeleteBook :exec
DELETE
FROM books
WHERE book_id = $1;

-- Media files related queries

-- Get all page information by book ID
-- name: GetPageInfosByBookID :many
SELECT *
FROM page_infos
WHERE book_id = $1
ORDER BY id;

-- Get specific page by book ID and page number
-- name: GetPageInfoByBookIDAndPage :one
SELECT *
FROM page_infos
WHERE book_id = $1
  AND page_num = $2
LIMIT 1;

-- Create media file record
-- name: CreatePageInfo :one
INSERT INTO page_infos (book_id, name, path, size, mod_time, url, page_num,
                        blurhash, height, width, img_type, insert_html)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- Update media file information
-- name: UpdatePageInfo :exec
UPDATE page_infos
SET name        = $1,
    path        = $2,
    size        = $3,
    mod_time    = $4,
    url         = $5,
    blurhash    = $6,
    height      = $7,
    width       = $8,
    img_type    = $9,
    insert_html = $10
WHERE book_id = $11
  AND page_num = $12;

-- Delete all media files for a book
-- name: DeletePageInfosByBookID :exec
DELETE
FROM page_infos
WHERE book_id = $1;

-- Bookmarks related queries

-- List bookmarks by book ID
-- name: ListBookmarksByBookID :many
SELECT *
FROM bookmarks
WHERE book_id = $1
ORDER BY created_at DESC;

-- Create a bookmark
-- name: CreateBookmark :one
INSERT INTO bookmarks (type, book_id, book_store_id, page_index, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- Update a bookmark (by book_id, type)
-- name: UpdateBookmark :exec
UPDATE bookmarks
SET description = $1,
    page_index = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE book_id = $3 and type = $4;

-- Delete a bookmark by (book_id, type)
-- name: DeleteBookmarkByBookIDAndType :exec
DELETE
FROM bookmarks
WHERE book_id = $1
  AND type = $2;

-- Delete all bookmarks for a book
-- name: DeleteBookmarksByBookID :exec
DELETE
FROM bookmarks
WHERE book_id = $1;

-- Statistics queries

-- Count total books
-- name: CountBooks :one
SELECT COUNT(*)
FROM books
WHERE deleted = FALSE;

-- Count books by type
-- name: CountBooksByType :one
SELECT COUNT(*)
FROM books
WHERE type = $1
  AND deleted = FALSE;

-- Count media files for a book
-- name: CountPageInfosByBookID :one
SELECT COUNT(*)
FROM page_infos
WHERE book_id = $1;
