-- 删除旧表（按外键依赖逆序）
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS indicators;
DROP TABLE IF EXISTS users;

-- 启用 UUID 扩展（提前执行）
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 用户表
CREATE TABLE users (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email          VARCHAR(50) NOT NULL UNIQUE,
    name           VARCHAR(50) NOT NULL UNIQUE,
    password       VARCHAR(255) NOT NULL,
    avatar         VARCHAR(512) NOT NULL DEFAULT '',
    gender         CHAR(1) CHECK(gender IN ('m', 'f', 'x')) DEFAULT 'x',
    bio            TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_type   SMALLINT CHECK(deleted_type IN (0, 1, 2)) DEFAULT 0,
    deleted_at     TIMESTAMPTZ NULL
);

-- 帖子表
CREATE TABLE posts (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author         UUID NOT NULL,
    title          VARCHAR(255) NOT NULL,
    content        TEXT NOT NULL,
    pv             INTEGER NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_type   SMALLINT CHECK(deleted_type IN (0, 1, 2)) DEFAULT 0,
    deleted_at     TIMESTAMPTZ NULL,
    CONSTRAINT fk_posts_author FOREIGN KEY (author) REFERENCES users(id) ON DELETE CASCADE
);

-- 评论表
CREATE TABLE comments (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author         UUID NOT NULL,
    content        TEXT NOT NULL,
    post_id        UUID NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_type   SMALLINT CHECK(deleted_type IN (0, 1, 2)) DEFAULT 0,
    deleted_at     TIMESTAMPTZ NULL,
    CONSTRAINT fk_comments_author FOREIGN KEY (author) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- 威胁情报指标表
CREATE TABLE indicators (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    indicator      VARCHAR(512) NOT NULL,
    indicator_type VARCHAR(20) NOT NULL,
    meta_source    VARCHAR(100) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ NULL
);

-- 创建触发器函数（移除IF判断，全版本兼容）
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER 
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
    -- 直接更新updated_at，移除兼容性风险的IF判断
    NEW.updated_at = NOW();
    RETURN NEW;
END;  -- 确保分号存在，函数体闭合
$$;  -- 结尾$$单独一行，避免语法粘连

-- 绑定触发器
CREATE TRIGGER trigger_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trigger_posts_updated_at
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trigger_comments_updated_at
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trigger_indicators_updated_at
BEFORE UPDATE ON indicators
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- 创建索引（高效+无冗余）
-- 基础索引
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_posts_author ON posts(author);
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_indicators_type ON indicators(indicator_type);

-- 复合索引
CREATE INDEX idx_posts_author_id_desc ON posts(author, id DESC);
CREATE INDEX idx_comments_post_id_id_asc ON comments(post_id, id ASC);

-- 软删除部分索引（仅索引未删除数据，性能更优）
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_deleted_at ON posts(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_comments_deleted_at ON comments(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_indicators_deleted_at ON indicators(deleted_at) WHERE deleted_at IS NULL;

-- 软删除类型索引
CREATE INDEX idx_users_deleted_type ON users(deleted_type);
CREATE INDEX idx_posts_deleted_type ON posts(deleted_type);
CREATE INDEX idx_comments_deleted_type ON comments(deleted_type);