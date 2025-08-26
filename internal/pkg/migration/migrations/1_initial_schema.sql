-- +goose Up
-- u_user: ユーザー情報
CREATE TABLE u_user (
    user_id CHAR(36) NOT NULL PRIMARY KEY,
    user_name CHAR(32) NOT NULL UNIQUE,
    nickname VARCHAR(32) NOT NULL,
    status VARCHAR(4096) NOT NULL DEFAULT "",
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- u_user_private: ユーザーのプライベート情報
CREATE TABLE u_user_private (
    user_id CHAR(36) NOT NULL PRIMARY KEY,
    email VARCHAR(255) NOT NULL DEFAULT "" UNIQUE,
    password_hash CHAR(60) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES u_user(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS u_user;
DROP TABLE IF EXISTS u_user_private;
