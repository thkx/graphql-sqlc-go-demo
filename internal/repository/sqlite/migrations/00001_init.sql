-- internal/repository/sqlite/migrations/00001_init.sql
-- 重新初始化数据库结构，为所有表添加 created_at 和 updated_at

DROP TABLE IF EXISTS indicators;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS comments;

-- 用户表
CREATE TABLE users (
    -- UUID 或自定义字符串 ID
    id TEXT PRIMARY KEY,
    -- 邮箱：字符串类型，非空，唯一
    email TEXT UNIQUE NOT NULL,
    -- 用户名：字符串类型，非空，唯一（对应原代码的 unique: true）
    name TEXT NOT NULL UNIQUE,
    -- 密码：字符串类型，非空
    password TEXT NOT NULL,
    -- 头像：字符串类型，非空
    avatar TEXT NOT NULL,
    -- 性别：只能是 m/f/x，默认值 x
    gender TEXT CHECK(gender IN ('m', 'f', 'x')) DEFAULT 'x',
    -- 个人简介：字符串类型，非空
    bio TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    deleted_at TEXT NULL
);

-- 为用户名创建索引（对应原代码的 index({ email: 1 })）
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
-- 为用户名创建索引（对应原代码的 index({ name: 1 })）
CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);

-- Post 表
CREATE TABLE posts (
    id         TEXT PRIMARY KEY,
    title      TEXT NOT NULL, -- 文章标题：字符串类型，非空
    content    TEXT NOT NULL, -- 文章内容：字符串类型，非空
    pv         INTEGER DEFAULT 0, -- 浏览量：数字类型，默认 0
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- 作者 ID：关联 users 表的 _id，非空
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    deleted_at TEXT NULL
)

-- 为作者和创建时间创建复合索引（按创建时间降序）
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id, id DESC);

-- Comment 表
CREATE TABLE comments (
    id         TEXT PRIMARY KEY,
    -- 评论作者 ID：关联 users 表
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- 评论内容：字符串类型，非空
    content TEXT NOT NULL,
    -- 关联的文章 ID
    post_id TEXT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    deleted_at TEXT NULL
)

-- 威胁情报指标表
CREATE TABLE indicators (
    id             TEXT PRIMARY KEY,
    indicator      TEXT NOT NULL,
    indicator_type TEXT NOT NULL,
    meta_source    TEXT NOT NULL,
    created_at     TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at     TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    deleted_at     TEXT NULL
);

-- 可选：为 updated_at 自动更新创建触发器（SQLite 支持）
CREATE TRIGGER IF NOT EXISTS update_users_timestamp
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = datetime('now', 'localtime') WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS update_posts_timestamp
AFTER UPDATE ON posts
FOR EACH ROW
BEGIN
    UPDATE posts SET updated_at = datetime('now', 'localtime') WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS update_comments_timestamp
AFTER UPDATE ON comments
FOR EACH ROW
BEGIN
    UPDATE comments SET updated_at = datetime('now', 'localtime') WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS update_indicators_timestamp
AFTER UPDATE ON indicators
FOR EACH ROW
BEGIN
    UPDATE indicators SET updated_at = datetime('now', 'localtime') WHERE id = OLD.id;
END;

-- 初始种子数据（保持原有测试数据）
INSERT INTO users (id, name, email, password, avatar, gender, bio)
VALUES ('1', 'test1', 'test1@qq.com', 'test1', '/avatar/default.png', 'm', '这是测试用户');

INSERT INTO posts (id, title, content, pv, user_id)
VALUES ('1', 'First post', 'This is the content of the first post.', 0, '1');

INSERT INTO comments (id, user_id, content, post_id)
VALUES ('1', '1', '这是一条测试评论', '1');