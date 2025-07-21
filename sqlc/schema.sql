-- User  information table
CREATE TABLE users
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    username   TEXT UNIQUE NOT NULL,               -- Username
    password   TEXT        NOT NULL,               -- Password (hashed)
    email      TEXT UNIQUE,                        -- Email address
    role       TEXT     DEFAULT 'user',            -- User role (admin/user)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Created time
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- Updated time
);


-- Book information table
CREATE TABLE books
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    title             TEXT        NOT NULL,               -- Book title
    book_id           TEXT UNIQUE NOT NULL,               -- Book ID, unique identifier generated from file path
    owner             TEXT     DEFAULT 'admin',           -- Owner
    file_path         TEXT        NOT NULL,               -- Absolute file path
    book_store_path   TEXT        NOT NULL,               -- Book store path
    type              TEXT        NOT NULL,               -- Book type
    child_book_num    INTEGER  DEFAULT 0,                 -- Number of child books
    depth             INTEGER  DEFAULT 0,                 -- Book depth
    parent_folder     TEXT,                               -- Parent folder
    page_count        INTEGER  DEFAULT 0,                 -- Total page count
    file_size         INTEGER  DEFAULT 0,                 -- File size
    author            TEXT,                               -- Author
    isbn              TEXT,                               -- ISBN
    press             TEXT,                               -- Publisher
    published_at      TEXT,                               -- Publication date
    extract_path      TEXT,                               -- Extract path
    modified_time     DATETIME DEFAULT CURRENT_TIMESTAMP, -- Modified time
    extract_num       INTEGER  DEFAULT 0,                 -- Extract number
    init_complete     BOOLEAN  DEFAULT FALSE,             -- Initialization complete flag
    read_percent      REAL     DEFAULT 0.0,               -- Reading progress
    non_utf8zip       BOOLEAN  DEFAULT FALSE,             -- Non-UTF8 zip flag
    zip_text_encoding TEXT,                               -- Zip text encoding
    deleted           BOOLEAN  DEFAULT FALSE              -- Soft delete flag
);

-- Media files information table (for storing page image information)
CREATE TABLE media_files
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

-- File backend storage configuration table
CREATE TABLE file_backends
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    type           INTEGER NOT NULL,                   -- File backend type (1: LocalDisk, 2: SMB, 3: SFTP, 4: WebDAV, 5: S3, 6: FTP)
    url            TEXT    NOT NULL,                   -- Storage URL
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

-- Book stores table
CREATE TABLE stores
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT    NOT NULL,                   -- Store name
    description     TEXT,                               -- Store description
    file_backend_id INTEGER NOT NULL,                   -- Associated file backend ID
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP, -- Created time
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP, -- Updated time
    FOREIGN KEY (file_backend_id) REFERENCES file_backends (id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_books_book_id ON books (book_id);
CREATE INDEX idx_books_file_path ON books (file_path);
CREATE INDEX idx_books_type ON books (type);
CREATE INDEX idx_books_modified_time ON books (modified_time);
CREATE INDEX idx_media_files_book_id ON media_files (book_id);
CREATE INDEX idx_media_files_page_num ON media_files (book_id, page_num);
CREATE INDEX idx_file_backends_type ON file_backends (type);
CREATE INDEX idx_stores_name ON stores (name);
CREATE INDEX idx_stores_file_backend_id ON stores (file_backend_id);

