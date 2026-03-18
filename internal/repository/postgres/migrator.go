package postgres

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func migrate(db *sql.DB) error {
	// PostgreSQL：检查 users 表是否存在
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
			  AND table_name = 'users'
		)
	`).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil // 已迁移
	}

	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, e := range entries {
		// 可选：跳过非 .sql 文件
		if e.IsDir() {
			continue
		}

		sqlBytes, err := migrationFS.ReadFile("migrations/" + e.Name())
		if err != nil {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", e.Name(), err)
		}
	}

	return tx.Commit()
}
