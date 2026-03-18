package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed "migrations/*.sql"
var migrationFS embed.FS

func migrate(db *sql.DB) error {
	// 检查是否已初始化
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if count > 0 {
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
