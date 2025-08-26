-- +goose Up
-- u_channel: チャンネル情報
CREATE TABLE u_channel (
    channel_id CHAR(36) NOT NULL PRIMARY KEY,
    channel_name CHAR(32) NOT NULL UNIQUE,
    display_name VARCHAR(32) NOT NULL,
    description VARCHAR(256) NOT NULL DEFAULT "",
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- u_message: メッセージ情報
CREATE TABLE u_message (
    message_id CHAR(36) NOT NULL PRIMARY KEY,
    channel_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (channel_id) REFERENCES u_channel(channel_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES u_user(user_id) ON DELETE CASCADE,
    INDEX idx_channel_id (channel_id),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- u_pinned_message: ピン留めされたメッセージ情報
CREATE TABLE u_pinned_message (
    message_id CHAR(36) NOT NULL,
    channel_id CHAR(36) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (message_id) REFERENCES u_message(message_id) ON DELETE CASCADE,
    FOREIGN KEY (channel_id) REFERENCES u_channel(channel_id) ON DELETE CASCADE,
    UNIQUE KEY unique_message_channel (message_id, channel_id),
    INDEX idx_channel_id (channel_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- u_user: ユーザー情報
CREATE TABLE u_user (
    user_id CHAR(36) NOT NULL PRIMARY KEY,
    user_name CHAR(32) NOT NULL UNIQUE,
    nickname VARCHAR(32) NOT NULL,
    status VARCHAR(4096) NOT NULL DEFAULT "",
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_name (user_name)
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

-- u_tags: タグ情報
CREATE TABLE u_tags (
    hashed_tag CHAR(32) NOT NULL PRIMARY KEY,
    raw_tag VARCHAR(32) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- u_tag_relation: メッセージとタグの関連
CREATE TABLE u_tag_relation (
    message_id CHAR(36) NOT NULL,
    channel_id CHAR(36) NOT NULL,
    hashed_tag CHAR(32) NOT NULL,
    FOREIGN KEY (message_id) REFERENCES u_message(message_id) ON DELETE CASCADE,
    FOREIGN KEY (channel_id) REFERENCES u_channel(channel_id) ON DELETE CASCADE,
    FOREIGN KEY (hashed_tag) REFERENCES u_tags(hashed_tag) ON DELETE CASCADE,
    UNIQUE KEY unique_message_tag (message_id, hashed_tag),
    INDEX idx_channel_id_hashed_tag (channel_id, hashed_tag)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS u_tag_relation;
DROP TABLE IF EXISTS u_tags;
DROP TABLE IF EXISTS u_pinned_message;
DROP TABLE IF EXISTS u_message;
DROP TABLE IF EXISTS u_user_private;
DROP TABLE IF EXISTS u_user;
DROP TABLE IF EXISTS u_channel;
