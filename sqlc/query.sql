-- Book related queries

-- Get single book by ID
-- name: GetBookByID :one
SELECT *
FROM books
WHERE book_id = ?
LIMIT 1;

-- Get book by file path
-- name: GetBookByBookPath :one
SELECT *
FROM books
WHERE book_path = ?
LIMIT 1;

-- List all books
-- name: ListBooks :many
SELECT *
FROM books
ORDER BY modified_time DESC;

-- List books by type
-- name: ListBooksByType :many
SELECT *
FROM books
WHERE type = ?
ORDER BY modified_time DESC;

-- get all store_url for books
-- name: ListAllBookStoreURLs :many
SELECT DISTINCT store_url
FROM books;

-- List books by store path
-- name: ListBooksByStorePath :many
SELECT *
FROM books
WHERE store_url = ?
ORDER BY modified_time DESC;

-- Search books by title (fuzzy search)
-- name: SearchBooksByTitle :many
SELECT *
FROM books
WHERE title LIKE '%' || ? || '%'
ORDER BY modified_time DESC;

-- Create new book
-- name: CreateBook :one
INSERT INTO books (title, book_id, owner, book_path, store_url, type,
                   child_books_num, child_books_id, depth, parent_folder, page_count, last_read_page, file_size,
                   author, isbn, press, published_at, extract_path, extract_num, book_complete,
                   init_complete, non_utf8zip, zip_text_encoding, created_by_version)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- Update book information
-- name: UpdateBook :exec
UPDATE books
SET title             = ?,
    owner             = ?,
    book_path         = ?,
    store_url         = ?,
    type              = ?,
    child_books_num   = ?,
    child_books_id    = ?,
    depth             = ?,
    parent_folder     = ?,
    page_count        = ?,
    last_read_page    = ?,
    file_size         = ?,
    author            = ?,
    isbn              = ?,
    press             = ?,
    published_at      = ?,
    extract_path      = ?,
    extract_num       = ?,
    book_complete     = ?,
    init_complete     = ?,
    non_utf8zip       = ?,
    zip_text_encoding = ?,
    modified_time     = CURRENT_TIMESTAMP
WHERE book_id = ?;

-- Update reading progress
-- name: UpdateLastReadPage :exec
UPDATE bookmarks
SET page_index = ?,
    type  = ?,
    updated_at  = CURRENT_TIMESTAMP
WHERE book_id = ?;

-- Mark book as deleted (soft delete)
-- name: MarkBookAsDeleted :exec
UPDATE books
SET deleted       = TRUE,
    modified_time = CURRENT_TIMESTAMP
WHERE book_id = ?;

-- Delete book
-- name: DeleteBook :exec
DELETE
FROM books
WHERE book_id = ?;

-- Media files related queries

-- Get all page information by book ID
-- name: GetPageInfosByBookID :many
SELECT *
FROM page_infos
WHERE book_id = ?
ORDER BY page_num;

-- Get specific page by book ID and page number
-- name: GetPageInfoByBookIDAndPage :one
SELECT *
FROM page_infos
WHERE book_id = ?
  AND page_num = ?
LIMIT 1;


-- Create media file record
-- name: CreatePageInfo :one
INSERT INTO page_infos (book_id, name, path, size, mod_time, url, page_num,
                         blurhash, height, width, img_type, insert_html)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- Update media file information
-- name: UpdatePageInfo :exec
UPDATE page_infos
SET name        = ?,
    path        = ?,
    size        = ?,
    mod_time    = ?,
    url         = ?,
    blurhash    = ?,
    height      = ?,
    width       = ?,
    img_type    = ?,
    insert_html = ?
WHERE book_id = ?
  AND page_num = ?;

-- Delete all media files for a book
-- name: DeletePageInfosByBookID :exec
DELETE
FROM page_infos
WHERE book_id = ?;

-- Bookmarks related queries

-- List bookmarks by book ID
-- name: ListBookmarksByBookID :many
SELECT *
FROM bookmarks
WHERE book_id = ?
ORDER BY created_at DESC;

-- Create a bookmark
-- name: CreateBookmark :one
INSERT INTO bookmarks (type, book_id, page_index, description, created_at, updated_at)
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;

-- Update a bookmark (by book_id, type)
-- name: UpdateBookmark :exec
UPDATE bookmarks
SET description = ?,
    page_index  = ?,
    description = ?,
    updated_at  = CURRENT_TIMESTAMP
WHERE book_id = ? and type = ?;

-- Delete a bookmark by (book_id, type)
-- name: DeleteBookmarkByBookIDAndType :exec
DELETE
FROM bookmarks
WHERE book_id = ?
  AND type = ?;

-- Delete all bookmarks for a book
-- name: DeleteBookmarksByBookID :exec
DELETE
FROM bookmarks
WHERE book_id = ?;


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
WHERE type = ?
  AND deleted = FALSE;

-- Count media files for a book
-- name: CountPageInfosByBookID :one
SELECT COUNT(*)
FROM page_infos
WHERE book_id = ?;

-- User related queries
-- Get user by ID
-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;

-- Get user by username
-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = ?
LIMIT 1;

-- Get user by email
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?
LIMIT 1;

-- List all users
-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC;

-- Create new user
-- name: CreateUser :one
INSERT INTO users (username, password, role, email, key, expires_at)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- Update user information
-- name: UpdateUser :exec
UPDATE users
SET username   = ?,
    password   = ?,
    role       = ?,
    email      = ?,
    key        = ?,
    expires_at = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- Update user password
-- name: UpdateUserPassword :exec
UPDATE users
SET password   = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- Update user key and expiration
-- name: UpdateUserKey :exec
UPDATE users
SET key        = ?,
    expires_at = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- Delete user
-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = ?;

-- Count total users
-- name: CountUsers :one
SELECT COUNT(*)
FROM users;

-- Count users by role
-- name: CountUsersByRole :one
SELECT COUNT(*)
FROM users
WHERE role = ?;