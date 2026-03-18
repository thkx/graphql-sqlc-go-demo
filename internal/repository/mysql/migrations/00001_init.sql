-- internal/repository/mysql/migrations/00001_init.sql
-- 重新初始化数据库结构，为所有表添加 created_at 和 updated_at
-- 修正：表选项（ENGINE/CHARSET/COLLATE）放在字段和约束之后，移除多余逗号

DROP TABLE IF EXISTS indicators;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;

-- 用户表（修正表选项位置，移除多余逗号）
CREATE TABLE users (
    id         VARCHAR(36) PRIMARY KEY,                -- UUID 或自定义字符串 ID
    name       VARCHAR(50) NOT NULL,                   -- 替换 TEXT 为 VARCHAR(50)
    email      VARCHAR(100) UNIQUE NOT NULL,           -- 邮箱设为 VARCHAR(100)
    password   VARCHAR(100) NOT NULL,                   -- 密码加密串设为 VARCHAR(100)
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL                            -- 最后一个字段后无逗号
) ENGINE = InnoDB                                      -- 表选项放在括号外，无前置逗号
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- Todo 表（修正表选项位置 + 外键约束格式）
CREATE TABLE todos (
    id         VARCHAR(36) PRIMARY KEY,
    text       VARCHAR(255) NOT NULL,                  -- 常规 Todo 文本长度 VARCHAR(255)
    done       TINYINT(1) NOT NULL DEFAULT 0,          -- MySQL 布尔等价类型 TINYINT(1)
    user_id    VARCHAR(36) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    -- 外键约束放在字段列表内（正确位置）
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 威胁情报指标表（修正表选项位置）
CREATE TABLE indicators (
    id             VARCHAR(36) PRIMARY KEY,
    indicator      VARCHAR(255) NOT NULL,               -- 常规指标长度 VARCHAR(255)
    indicator_type VARCHAR(50) NOT NULL,                -- 指标类型 VARCHAR(50)
    meta_source    VARCHAR(100) NOT NULL,               -- 数据源名称 VARCHAR(100)
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     DATETIME NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 为 updated_at 添加自动更新功能（MySQL 原生特性，替代 SQLite 触发器）
ALTER TABLE users MODIFY COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;
ALTER TABLE todos MODIFY COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;
ALTER TABLE indicators MODIFY COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;

-- 初始种子数据（无语法错误，直接执行）
INSERT INTO users (id, name, email, password)
VALUES ('1', 'test1', 'test1@qq.com', 'test1');

INSERT INTO todos (id, text, done, user_id)
VALUES ('1', 'First todo', 1, '1');