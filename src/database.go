package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS websites (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT UNIQUE NOT NULL,
			title TEXT,
			content_hash TEXT NOT NULL,
			last_recorded DATETIME NOT NULL,
			content_report TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func closeDB() {
	if db != nil {
		db.Close()
	}
}

func saveWebsite(url, title, htmlContent, report string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	hash := hashContent(htmlContent)
	now := time.Now()

	_, err := db.Exec(`
		INSERT INTO websites (url, title, content_hash, last_recorded, content_report, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(url) DO UPDATE SET
			title = excluded.title,
			content_hash = excluded.content_hash,
			last_recorded = excluded.last_recorded,
			content_report = excluded.content_report,
			updated_at = excluded.updated_at
	`, url, title, hash, now, report, now)

	if err != nil {
		return fmt.Errorf("failed to save website: %w", err)
	}

	return nil
}

func getWebsite(url string) (map[string]interface{}, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var id int
	var urlStr, title, hash, lastRecorded, report, createdAt, updatedAt string

	err := db.QueryRow(`
		SELECT id, url, title, content_hash, last_recorded, content_report, created_at, updated_at
		FROM websites WHERE url = ?
	`, url).Scan(&id, &urlStr, &title, &hash, &lastRecorded, &report, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query website: %w", err)
	}

	return map[string]interface{}{
		"id":             id,
		"url":            urlStr,
		"title":          title,
		"content_hash":   hash,
		"last_recorded":  lastRecorded,
		"content_report": report,
		"created_at":     createdAt,
		"updated_at":     updatedAt,
	}, nil
}

func listWebsites(limit int) ([]map[string]interface{}, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query(`
		SELECT id, url, title, content_hash, last_recorded, content_report, created_at, updated_at
		FROM websites ORDER BY updated_at DESC LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query websites: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var urlStr, title, hash, lastRecorded, report, createdAt, updatedAt string
		if err := rows.Scan(&id, &urlStr, &title, &hash, &lastRecorded, &report, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id":             id,
			"url":            urlStr,
			"title":          title,
			"content_hash":   hash,
			"last_recorded":  lastRecorded,
			"content_report": report,
			"created_at":     createdAt,
			"updated_at":     updatedAt,
		})
	}

	return results, nil
}

func deleteWebsite(url string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	result, err := db.Exec("DELETE FROM websites WHERE url = ?", url)
	if err != nil {
		return fmt.Errorf("failed to delete website: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("website not found")
	}

	return nil
}

func hashContent(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

func ensureDB() {
	if db == nil {
		home, _ := os.UserHomeDir()
		dir := home + "/.whats-new"
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to create directory: %v\n", err)
		}
		dbPath := dir + "/website.db"
		if err := initDB(dbPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to init database: %v\n", err)
		}
	}
}
