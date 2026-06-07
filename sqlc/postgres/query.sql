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

-- Create new book
-- name: CreateBook :one
INSERT INTO books (title, book_id, owner, book_path, store_url, type,
                   child_books_num, child_books_id, depth, parent_folder, page_count, last_read_page, file_size,
                   author, isbn, press, published_at, extract_path, extract_num, book_complete,
                   init_complete, non_utf8zip, zip_text_encoding, created_by_version, is_remote, remote_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)
RETURNING *;

-- Update book information
-- name: UpdateBook :exec
UPDATE books
SET title              = $1,
    owner              = $2,
    book_path          = $3,
    store_url          = $4,
    type               = $5,
    child_books_num    = $6,
    child_books_id     = $7,
    depth              = $8,
    parent_folder      = $9,
    page_count         = $10,
    last_read_page     = $11,
    file_size          = $12,
    author             = $13,
    isbn               = $14,
    press              = $15,
    published_at       = $16,
    extract_path       = $17,
    extract_num        = $18,
    book_complete      = $19,
    init_complete      = $20,
    non_utf8zip        = $21,
    zip_text_encoding  = $22,
    created_by_version = $23,
    is_remote          = $24,
    remote_url         = $25,
    modified_time      = CURRENT_TIMESTAMP
WHERE book_id = $26;

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
ORDER BY page_num;

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

-- User related queries
-- Get user by ID
-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- Get user by username
-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- Get user by email
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- List all users
-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC;

-- Create new user
-- name: CreateUser :one
INSERT INTO users (username, password, role, email, key, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- Update user information
-- name: UpdateUser :exec
UPDATE users
SET username = $1,
    password = $2,
    role = $3,
    email = $4,
    key = $5,
    expires_at = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $7;

-- Update user password
-- name: UpdateUserPassword :exec
UPDATE users
SET password = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- Update user key and expiration
-- name: UpdateUserKey :exec
UPDATE users
SET key = $1,
    expires_at = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $3;

-- Delete user
-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;

-- Count total users
-- name: CountUsers :one
SELECT COUNT(*)
FROM users;

-- Count users by role
-- name: CountUsersByRole :one
SELECT COUNT(*)
FROM users
WHERE role = $1;
