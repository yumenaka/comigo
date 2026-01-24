-- User  information table
CREATE TABLE IF NOT EXISTS users
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    username   TEXT UNIQUE NOT NULL,               -- Username
    password   TEXT        NOT NULL,               -- Password (hashed)
    role       TEXT     DEFAULT 'user',            -- User role (admin/user)
    email      TEXT UNIQUE,                        -- Email address
    key        TEXT,                               -- User key (for API access)
    expires_at DATETIME,                           -- Key expiration time
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Created time
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- Updated time
);


-- Book information table
CREATE TABLE IF NOT EXISTS books
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    title             TEXT        NOT NULL,               -- Book title
    book_id           TEXT UNIQUE NOT NULL,               -- Book ID, unique identifier generated from book path
    owner             TEXT     DEFAULT 'admin',           -- Owner
    book_path         TEXT        NOT NULL,               -- Absolute book path
    store_url         TEXT        NOT NULL,               -- Bookstore url
    type              TEXT        NOT NULL,               -- Book type
    child_books_num   INTEGER  DEFAULT 0,                 -- Number of child books
    child_books_id    TEXT,                               -- Child book IDs (comma-separated)
    depth             INTEGER  DEFAULT 0,                 -- Book depth
    parent_folder     TEXT,                               -- Parent folder
    page_count        INTEGER  DEFAULT 0,                 -- Total page count
    last_read_page    INTEGER  DEFAULT 0,                 -- Last read position
    file_size         INTEGER  DEFAULT 0,                 -- File size
    author            TEXT,                               -- Author
    isbn              TEXT,                               -- ISBN
    press             TEXT,                               -- Publisher
    published_at      TEXT,                               -- Publication date
    extract_path      TEXT,                               -- Extract path
    modified_time     DATETIME DEFAULT CURRENT_TIMESTAMP, -- Modified time
    extract_num       INTEGER  DEFAULT 0,                 -- Extract number
    book_complete     BOOLEAN  DEFAULT FALSE,             -- Book complete flag
    init_complete     BOOLEAN  DEFAULT FALSE,             -- Initialization complete flag
    non_utf8zip       BOOLEAN  DEFAULT FALSE,             -- Non-UTF8 zip flag
    zip_text_encoding  TEXT,                               -- Zip text encoding
    created_by_version TEXT,                               -- Comigo version when data was created
    deleted            BOOLEAN  DEFAULT FALSE              -- Soft delete flag
);

-- Pages information table (for storing page image information)
CREATE TABLE IF NOT EXISTS page_infos
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id     TEXT NOT NULL,                                         -- Associated book ID
    name        TEXT NOT NULL,                                         -- File name (for compressed file path or image name)
    path        TEXT,                                                  -- File path
    size        INTEGER DEFAULT 0,                                     -- File size
    mod_time    DATETIME,                                              -- Modified time
    url         TEXT,                                                  -- Remote URL for reading images
    page_num    INTEGER DEFAULT 0,                                     -- Page number
    blurhash    TEXT,                                                  -- Blurhash placeholder
    height      INTEGER DEFAULT 0,                                     -- Image height
    width       INTEGER DEFAULT 0,                                     -- Image width
    img_type    TEXT,                                                  -- Image type
    insert_html TEXT,                                                  -- Insert HTML
    FOREIGN KEY (book_id) REFERENCES books (book_id) ON DELETE CASCADE -- 外键约束:当引用表（books）中某条记录被 删除 时，将会 自动删除 当前表中所有引用该记录的数据
);

-- Bookmarks table
CREATE TABLE IF NOT EXISTS bookmarks
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    type        TEXT NOT NULL,                     -- Bookmark type (like: "auto", "user")
    book_id     TEXT    NOT NULL,                   -- Associated book ID
    page_index  INTEGER NOT NULL,                   -- Page index, starts from 0
    description TEXT,                               -- User note
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP, -- Created time
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP -- Updated time
    -- FOREIGN KEY (book_id) REFERENCES books (book_id) ON DELETE CASCADE
);


-- Create indexes for better query performance / 创建索引以获得更好的查询性
CREATE INDEX IF NOT EXISTS idx_books_book_id ON books (book_id);
CREATE INDEX IF NOT EXISTS idx_books_file_path ON books (book_path);
CREATE INDEX IF NOT EXISTS idx_books_type ON books (type);
CREATE INDEX IF NOT EXISTS idx_books_modified_time ON books (modified_time);
CREATE INDEX IF NOT EXISTS idx_page_infos_book_id ON page_infos (book_id);
CREATE INDEX IF NOT EXISTS idx_page_infos_page_num ON page_infos (book_id, page_num);
CREATE INDEX IF NOT EXISTS idx_bookmarks_book_id ON bookmarks (book_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_book_page ON bookmarks (book_id, page_index);

