package mysql

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"sort"
	"strings"
)

//go:embed "migrations/*.sql"
var migrationFS embed.FS

// migrate 兼容 MySQL 的数据库迁移函数（替换 SQLite 专属逻辑）
func migrate(db *sql.DB) error {
	// 1. 检查 users 表是否已存在（MySQL 替代 sqlite_master 的查询逻辑）
	var tableCount int
	// information_schema.TABLES 是 MySQL 系统表，存储所有表元数据
	// TABLE_SCHEMA：数据库名（需替换为你的实际数据库名，或用 DATABASE() 获取当前数据库）
	// TABLE_NAME：表名
	query := `
		SELECT COUNT(*) 
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users'
	`
	err := db.QueryRow(query).Scan(&tableCount)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check table existence: %w", err)
	}
	// 若 users 表已存在，说明迁移已完成，直接返回
	if tableCount > 0 {
		return nil
	}

	// 2. 读取 migrations 目录下所有文件
	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migration dir: %w", err)
	}

	// 3. 过滤 + 排序迁移文件（确保按名称升序执行，避免顺序混乱）
	var sqlFiles []os.DirEntry
	for _, e := range entries {
		// 过滤非 .sql 后缀的文件
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".sql") {
			sqlFiles = append(sqlFiles, e)
		}
	}

	// 按文件名升序排序（保证 00001_init.sql 先于 00002_xxx.sql 执行）
	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	// 4. 开启事务执行迁移（原子操作，失败则回滚）
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin migration transaction: %w", err)
	}
	// 延迟回滚（若后续执行失败，自动触发回滚；若提交成功，回滚无效果）
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 5. 遍历执行所有迁移脚本
	for _, e := range sqlFiles {
		fileName := e.Name()
		// 读取迁移文件内容
		sqlBytes, err := migrationFS.ReadFile("migrations/" + fileName)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}
		sqlContent := string(sqlBytes)
		// 执行迁移 SQL
		_, err = tx.Exec(sqlContent)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("migration file %s execute failed: %w", fileName, err)
		}
	}

	// 6. 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration transaction: %w", err)
	}

	fmt.Println("database migration success")
	return nil
}
