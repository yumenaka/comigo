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
    init_complete     BOOLEAN  DEFAULT FALSE,             -- Initialization complete flag
    non_utf8zip       BOOLEAN  DEFAULT FALSE,             -- Non-UTF8 zip flag
    zip_text_encoding TEXT,                               -- Zip text encoding
    deleted           BOOLEAN  DEFAULT FALSE              -- Soft delete flag
);

-- Media files information table (for storing page image information)
CREATE TABLE IF NOT EXISTS media_files
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id     TEXT NOT NULL,     -- Associated book ID
    name        TEXT NOT NULL,     -- File name (for compressed file path or image name)
    path        TEXT,              -- File path
    size        INTEGER DEFAULT 0, -- File size
    mod_time    DATETIME,          -- Modified time
    url         TEXT,              -- Remote URL for reading images
    page_num    INTEGER DEFAULT 0, -- Page number
    blurhash    TEXT,              -- Blurhash placeholder
    height      INTEGER DEFAULT 0, -- Image height
    width       INTEGER DEFAULT 0, -- Image width
    img_type    TEXT,              -- Image type
    insert_html TEXT,              -- Insert HTML
    FOREIGN KEY (book_id) REFERENCES books (book_id) ON DELETE CASCADE
);

-- Bookmarks table
CREATE TABLE IF NOT EXISTS bookmarks
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id     TEXT    NOT NULL,                   -- Associated book ID
    page_index  INTEGER NOT NULL,                   -- Page index, starts from 0
    description TEXT,                               -- User note
    position    REAL     DEFAULT 0.0,               -- Position percentage (0.0 - 100.0)
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP, -- Created time
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP, -- Updated time
    FOREIGN KEY (book_id) REFERENCES books (book_id) ON DELETE CASCADE
);

-- Book stores table
CREATE TABLE IF NOT EXISTS stores
(
    backend_url TEXT PRIMARY KEY NOT NULL,          -- Associated file backend ID
    name        TEXT             NOT NULL,          -- Store name
    description TEXT,                               -- Store description
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP, -- Created time
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP, -- Updated time
    FOREIGN KEY (backend_url) REFERENCES file_backends (url) ON DELETE CASCADE
);

-- File backend storage configuration table
CREATE TABLE IF NOT EXISTS file_backends
(
    url            TEXT PRIMARY KEY NOT NULL,          -- Store URL
    type           INTEGER          NOT NULL,          -- File backend type (1: LocalDisk, 2: SMB, 3: SFTP, 4: WebDAV, 5: S3, 6: FTP)
    server_host    TEXT,                               -- Server host address
    server_port    INTEGER  DEFAULT 0,                 -- Server port number
    need_auth      BOOLEAN  DEFAULT FALSE,             -- Whether authentication is required
    auth_username  TEXT,                               -- Authentication username
    auth_password  TEXT,                               -- Authentication password
    smb_share_name TEXT,                               -- SMB share name
    smb_path       TEXT,                               -- SMB share path
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP, -- Created time
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP  -- Updated time
);


-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_books_book_id ON books (book_id);
CREATE INDEX IF NOT EXISTS idx_books_file_path ON books (book_path);
CREATE INDEX IF NOT EXISTS idx_books_type ON books (type);
CREATE INDEX IF NOT EXISTS idx_books_modified_time ON books (modified_time);
CREATE INDEX IF NOT EXISTS idx_media_files_book_id ON media_files (book_id);
CREATE INDEX IF NOT EXISTS idx_media_files_page_num ON media_files (book_id, page_num);
CREATE INDEX IF NOT EXISTS idx_bookmarks_book_id ON bookmarks (book_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_book_page ON bookmarks (book_id, page_index);
CREATE INDEX IF NOT EXISTS idx_stores_url ON stores (backend_url);
CREATE INDEX IF NOT EXISTS idx_file_backends_url ON file_backends (url);

