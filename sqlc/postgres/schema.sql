-- User information table
CREATE TABLE IF NOT EXISTS users
(
    id         BIGSERIAL PRIMARY KEY,
    username   TEXT UNIQUE NOT NULL,
    password   TEXT        NOT NULL,
    role       TEXT      DEFAULT 'user',
    email      TEXT UNIQUE,
    key        TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Book information table
CREATE TABLE IF NOT EXISTS books
(
    id                 BIGSERIAL PRIMARY KEY,
    title              TEXT        NOT NULL,
    book_id            TEXT UNIQUE NOT NULL,
    owner              TEXT      DEFAULT 'admin',
    book_path          TEXT        NOT NULL,
    store_url          TEXT        NOT NULL,
    type               TEXT        NOT NULL,
    child_books_num    BIGINT    DEFAULT 0,
    child_books_id     TEXT,
    depth              BIGINT    DEFAULT 0,
    parent_folder      TEXT,
    page_count         BIGINT    DEFAULT 0,
    last_read_page     BIGINT    DEFAULT 0,
    file_size          BIGINT    DEFAULT 0,
    author             TEXT,
    isbn               TEXT,
    press              TEXT,
    published_at       TEXT,
    extract_path       TEXT,
    modified_time      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    extract_num        BIGINT    DEFAULT 0,
    book_complete      BOOLEAN   DEFAULT FALSE,
    init_complete      BOOLEAN   DEFAULT FALSE,
    non_utf8zip        BOOLEAN   DEFAULT FALSE,
    zip_text_encoding  TEXT,
    created_by_version TEXT,
    is_remote          BOOLEAN   DEFAULT FALSE,
    remote_url         TEXT,
    remote_book_id TEXT,
    remote_store_key TEXT,
    remote_shelf_key TEXT,
    remote_shelf_name TEXT,
    cover_name TEXT,
    cover_path TEXT,
    cover_size BIGINT,
    cover_mod_time TIMESTAMP,
    cover_url TEXT,
    cover_page_num BIGINT,
    cover_blurhash TEXT,
    cover_height BIGINT,
    cover_width BIGINT,
    cover_img_type TEXT,
    cover_insert_html TEXT,
    deleted            BOOLEAN   DEFAULT FALSE
);

-- Pages information table (for storing page image information)
CREATE TABLE IF NOT EXISTS page_infos
(
    id          BIGSERIAL PRIMARY KEY,
    book_id     TEXT NOT NULL,
    name        TEXT NOT NULL,
    path        TEXT,
    size        BIGINT DEFAULT 0,
    mod_time    TIMESTAMP,
    url         TEXT,
    page_num    BIGINT DEFAULT 0,
    blurhash    TEXT,
    height      BIGINT DEFAULT 0,
    width       BIGINT DEFAULT 0,
    img_type    TEXT,
    insert_html TEXT,
    FOREIGN KEY (book_id) REFERENCES books (book_id) ON DELETE CASCADE
);

-- Bookmarks table
CREATE TABLE IF NOT EXISTS bookmarks
(
    id            BIGSERIAL PRIMARY KEY,
    type          TEXT    NOT NULL,
    book_id       TEXT    NOT NULL,
    book_store_id TEXT,
    page_index    BIGINT NOT NULL,
    description   TEXT,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_books_book_id ON books (book_id);
CREATE INDEX IF NOT EXISTS idx_books_file_path ON books (book_path);
CREATE INDEX IF NOT EXISTS idx_books_type ON books (type);
CREATE INDEX IF NOT EXISTS idx_books_modified_time ON books (modified_time);
CREATE INDEX IF NOT EXISTS idx_page_infos_book_id ON page_infos (book_id);
CREATE INDEX IF NOT EXISTS idx_page_infos_page_num ON page_infos (book_id, page_num);
CREATE INDEX IF NOT EXISTS idx_bookmarks_book_id ON bookmarks (book_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_book_page ON bookmarks (book_id, page_index);

-- 扫描失败指纹也由数据库持久化，不依赖 metadata JSON 文件。
CREATE TABLE IF NOT EXISTS scan_failures (
    key TEXT PRIMARY KEY,
    store_url TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    modified_unix_nano BIGINT NOT NULL,
    created_by_version TEXT NOT NULL,
    failed_at TIMESTAMP NOT NULL,
    error TEXT NOT NULL,
    is_remote BOOLEAN NOT NULL
);
